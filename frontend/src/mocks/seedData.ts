import type { FlightTurnaroundDetail } from "../types/FlightTurnaround";
import type { GroundTask } from "../types/GroundTask";
import type { GroundResourceView } from "../types/GroundResource";
import type { ResourceBooking } from "../types/ResourceBooking";
import type { DelayEvent } from "../types/DelayEvent";

const timeline = (status: string) => [
  { code: "ARRIVAL", label: "航班到达", occurred_at: "2026-09-22T08:00:00Z", status: "ARRIVING" },
  { code: "RESOURCE_OCCUPANCY", label: "资源占用", occurred_at: "2026-09-22T08:05:00Z", status }
];

export const mockData: {
  flightTurnaround: FlightTurnaroundDetail[];
  groundTask: GroundTask[];
  groundResource: GroundResourceView[];
  resourceBooking: ResourceBooking[];
  delayEvent: DelayEvent[];
} = {
  flightTurnaround: [
    {
      id: 1,
      flight_no: "CA1801",
      aircraft_reg: "B-2026",
      stand_no: "A12",
      arrival_time: "2026-09-22T08:00:00Z",
      departure_time: "2026-09-22T09:10:00Z",
      turnaround_status: "IN_SERVICE",
      delay_reason: "",
      ready_at: null,
      timeline_synced_at: null,
      task_completion_rate: 100,
      timeline: timeline("ACTIVE")
    },
    {
      id: 2,
      flight_no: "MU5102",
      aircraft_reg: "B-5088",
      stand_no: "B07",
      arrival_time: "2026-09-22T10:00:00Z",
      departure_time: "2026-09-22T11:10:00Z",
      turnaround_status: "IN_SERVICE",
      delay_reason: "行李等待",
      ready_at: null,
      timeline_synced_at: null,
      task_completion_rate: 50,
      timeline: timeline("CONFLICT")
    },
    {
      id: 3,
      flight_no: "CZ3120",
      aircraft_reg: "B-6310",
      stand_no: "C03",
      arrival_time: "2026-09-22T12:00:00Z",
      departure_time: "2026-09-22T13:10:00Z",
      turnaround_status: "ON_STAND",
      delay_reason: "",
      ready_at: null,
      timeline_synced_at: null,
      task_completion_rate: 100,
      timeline: timeline("ACTIVE")
    }
  ],
  groundTask: [
    { id: 101, turnaround_id: 1, task_type: "CLEANING", team_id: 1, planned_start: "2026-09-22T08:05:00Z", deadline: "2026-09-22T08:40:00Z", actual_finish: "2026-09-22T09:00:00Z", status: "SIGNED", blocker_note: "", timeline_synced_at: null },
    { id: 102, turnaround_id: 1, task_type: "BAGGAGE", team_id: 2, planned_start: "2026-09-22T08:05:00Z", deadline: "2026-09-22T08:50:00Z", actual_finish: "2026-09-22T09:00:00Z", status: "SIGNED", blocker_note: "", timeline_synced_at: null },
    { id: 201, turnaround_id: 2, task_type: "CLEANING", team_id: 1, planned_start: "2026-09-22T10:05:00Z", deadline: "2026-09-22T10:40:00Z", actual_finish: null, status: "SIGNED", blocker_note: "", timeline_synced_at: null },
    { id: 202, turnaround_id: 2, task_type: "REFUEL", team_id: 3, planned_start: "2026-09-22T10:10:00Z", deadline: "2026-09-22T10:50:00Z", actual_finish: null, status: "BLOCKED", blocker_note: "加油栓井占用", timeline_synced_at: null },
    { id: 301, turnaround_id: 3, task_type: "CATERING", team_id: 4, planned_start: "2026-09-22T12:05:00Z", deadline: "2026-09-22T12:40:00Z", actual_finish: null, status: "SIGNED", blocker_note: "", timeline_synced_at: null }
  ],
  groundResource: [
    { id: 1, resource_code: "CLN-A12", resource_type: "CLEANING", location: "A12", availability_status: "BOOKED", maintenance_due_at: "2026-10-01T00:00:00Z", owner_team: "清洁一组", active_booking_count: 1 },
    { id: 2, resource_code: "BAG-07", resource_type: "BAGGAGE", location: "B07", availability_status: "BOOKED", maintenance_due_at: "2026-10-02T00:00:00Z", owner_team: "行李班组", active_booking_count: 1 },
    { id: 3, resource_code: "FUEL-B07", resource_type: "REFUEL", location: "B07", availability_status: "BOOKED", maintenance_due_at: "2026-10-03T00:00:00Z", owner_team: "加油班组", active_booking_count: 0 },
    { id: 4, resource_code: "CAT-C03", resource_type: "CATERING", location: "C03", availability_status: "BOOKED", maintenance_due_at: "2026-10-04T00:00:00Z", owner_team: "配餐班组", active_booking_count: 1 }
  ],
  resourceBooking: [
    { id: 1001, resource_id: 1, turnaround_id: 1, task_id: 101, start_time: "2026-09-22T08:05:00Z", end_time: "2026-09-22T08:40:00Z", booking_status: "ACTIVE", conflict_reason: "", released_at: null },
    { id: 1002, resource_id: 2, turnaround_id: 1, task_id: 102, start_time: "2026-09-22T08:05:00Z", end_time: "2026-09-22T08:50:00Z", booking_status: "ACTIVE", conflict_reason: "", released_at: null },
    { id: 2001, resource_id: 3, turnaround_id: 2, task_id: 202, start_time: "2026-09-22T10:10:00Z", end_time: "2026-09-22T10:50:00Z", booking_status: "CONFLICT", conflict_reason: "加油栓井与邻机位冲突", released_at: null },
    { id: 3001, resource_id: 4, turnaround_id: 3, task_id: 301, start_time: "2026-09-22T12:05:00Z", end_time: "2026-09-22T12:40:00Z", booking_status: "ACTIVE", conflict_reason: "", released_at: null }
  ],
  delayEvent: [
    { id: 1, turnaround_id: 2, delay_type: "GROUND_HANDLING", minutes: 25, root_cause: "加油资源冲突", responsibility_team: "加油班组", resolved_at: null }
  ]
};
