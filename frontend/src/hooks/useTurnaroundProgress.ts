import { useMemo } from "react";
import type { FlightTurnaroundDetail } from "../types/FlightTurnaround";
import type { GroundTask, TaskBlocker } from "../types/GroundTask";
import type { DelayEvent, DelayBlocker } from "../types/DelayEvent";
import type { ResourceBooking, BookingBlocker } from "../types/ResourceBooking";
import type { GroundResource } from "../types/GroundResource";
import type { ReleaseBlockers } from "../types/Release";

type Input = {
  flight: FlightTurnaroundDetail;
  tasks: GroundTask[];
  delays: DelayEvent[];
  bookings: ResourceBooking[];
  resources?: GroundResource[];
};

export function useTurnaroundProgress({ flight, tasks, delays, bookings, resources = [] }: Input) {
  return useMemo(() => {
    const flightTasks = tasks.filter((task) => task.turnaround_id === flight.id);
    const flightDelays = delays.filter((delay) => delay.turnaround_id === flight.id && !delay.resolved_at);
    const flightBookings = bookings.filter((booking) => booking.turnaround_id === flight.id);
    const resourceNames = new Map(resources.map((resource) => [resource.id, resource.resource_code]));
    const blockers: ReleaseBlockers = {
      tasks: flightTasks
        .filter((task) => task.status !== "SIGNED")
        .map((task): TaskBlocker => ({ ...task, reason: "任务尚未 SIGNED" })),
      delays: flightDelays
        .map((delay): DelayBlocker => ({ ...delay, reason: "延误事件尚未关闭" })),
      bookings: flightBookings
        .filter((booking) => booking.booking_status === "CONFLICT")
        .map((booking): BookingBlocker => ({
          ...booking,
          resource_code: resourceNames.get(booking.resource_id) ?? `RESOURCE-${booking.resource_id}`,
          reason: "资源预约存在未处理冲突"
        }))
    };
    const activeBookings = flightBookings.filter((booking) => booking.booking_status === "ACTIVE");
    const signedCount = flightTasks.filter((task) => task.status === "SIGNED").length;
    const ready = blockers.tasks.length === 0 && blockers.delays.length === 0 && blockers.bookings.length === 0;
    return {
      blockers,
      ready,
      activeBookings,
      activeBookingCount: activeBookings.length,
      taskCompletionRate: flightTasks.length === 0 ? 100 : Math.round(signedCount * 100 / flightTasks.length),
      timelineSynced: Boolean(flight.timeline_synced_at)
    };
  }, [flight, tasks, delays, bookings, resources]);
}
