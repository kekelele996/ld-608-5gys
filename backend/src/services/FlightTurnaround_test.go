package services

import (
	"net/http"
	"testing"

	"groundTurn/src/constants"
	"groundTurn/src/repositories"
)

func TestReleaseSucceedsAtomicallyAndRejectsDuplicate(t *testing.T) {
	store := repositories.NewStore(repositories.SeedData())
	service := NewFlightTurnaroundService(store)

	result, err := service.Release(1)
	if err != nil {
		t.Fatalf("expected release to succeed, got %v", err)
	}
	if result.Turnaround.TurnaroundStatus != constants.TurnaroundReady {
		t.Fatalf("expected READY, got %s", result.Turnaround.TurnaroundStatus)
	}
	if result.ReleasedCount != 2 {
		t.Fatalf("expected two active bookings released, got %d", result.ReleasedCount)
	}
	for _, booking := range result.Bookings {
		if booking.BookingStatus != constants.BookingReleased {
			t.Fatalf("booking %d remains %s", booking.ID, booking.BookingStatus)
		}
		if booking.ReleasedAt == nil {
			t.Fatalf("booking %d missing released_at", booking.ID)
		}
	}
	for _, task := range result.Tasks {
		if task.TimelineSyncedAt == nil {
			t.Fatalf("task %d timeline was not synced", task.ID)
		}
	}
	for _, resource := range store.Snapshot().Resources[:2] {
		if resource.AvailabilityStatus != constants.ResourceAvailable {
			t.Fatalf("resource %s status = %s", resource.ResourceCode, resource.AvailabilityStatus)
		}
	}

	_, err = service.Release(1)
	releaseErr, ok := err.(*ReleaseError)
	if !ok {
		t.Fatalf("expected ReleaseError, got %T", err)
	}
	if releaseErr.Code != constants.TurnaroundAlreadyReleased || releaseErr.HTTPCode != http.StatusConflict {
		t.Fatalf("unexpected duplicate release error: %+v", releaseErr)
	}
}

func TestReleaseRejectsAndReturnsAllBlockerLists(t *testing.T) {
	store := repositories.NewStore(repositories.SeedData())
	service := NewFlightTurnaroundService(store)
	before := store.Snapshot()

	_, err := service.Release(2)
	releaseErr, ok := err.(*ReleaseError)
	if !ok {
		t.Fatalf("expected ReleaseError, got %T", err)
	}
	if releaseErr.Code != constants.ReleaseBlocked || releaseErr.HTTPCode != http.StatusConflict {
		t.Fatalf("unexpected release error: %+v", releaseErr)
	}
	if releaseErr.Blockers == nil {
		t.Fatal("expected blocker lists")
	}
	if len(releaseErr.Blockers.Tasks) != 1 || releaseErr.Blockers.Tasks[0].ID != 202 {
		t.Fatalf("unexpected task blockers: %+v", releaseErr.Blockers.Tasks)
	}
	if len(releaseErr.Blockers.Delays) != 1 || releaseErr.Blockers.Delays[0].ID != 1 {
		t.Fatalf("unexpected delay blockers: %+v", releaseErr.Blockers.Delays)
	}
	if len(releaseErr.Blockers.Bookings) != 1 || releaseErr.Blockers.Bookings[0].ID != 2001 {
		t.Fatalf("unexpected booking blockers: %+v", releaseErr.Blockers.Bookings)
	}

	after := store.Snapshot()
	if after.Flights[1].TurnaroundStatus != before.Flights[1].TurnaroundStatus {
		t.Fatal("flight changed after failed release")
	}
	if after.Tasks[3].Status != constants.GroundTaskBlocked {
		t.Fatal("task changed after failed release")
	}
	if after.Bookings[2].BookingStatus != constants.BookingConflict {
		t.Fatal("booking changed after failed release")
	}
	if len(after.Logs) != len(before.Logs) {
		t.Fatal("audit log changed after failed release")
	}
}

func TestConcurrentReleasesApplyExactlyOnce(t *testing.T) {
	store := repositories.NewStore(repositories.SeedData())
	service := NewFlightTurnaroundService(store)
	start := make(chan struct{})
	results := make(chan error, 8)

	for range 8 {
		go func() {
			<-start
			_, err := service.Release(3)
			results <- err
		}()
	}
	close(start)

	successes := 0
	failures := 0
	for range 8 {
		if err := <-results; err == nil {
			successes++
		} else {
			failures++
		}
	}
	if successes != 1 || failures != 7 {
		t.Fatalf("expected exactly one success, got successes=%d failures=%d", successes, failures)
	}
	flight, _ := store.GetFlight(3)
	if flight.TurnaroundStatus != constants.TurnaroundReady {
		t.Fatalf("expected READY, got %s", flight.TurnaroundStatus)
	}
	if len(store.Snapshot().Logs) != 1 {
		t.Fatalf("expected one audit log, got %d", len(store.Snapshot().Logs))
	}
}
