import type { TaskBlocker, GroundTask } from "./GroundTask";
import type { DelayBlocker } from "./DelayEvent";
import type { BookingBlocker, ResourceBooking } from "./ResourceBooking";
import type { FlightTurnaroundDetail } from "./FlightTurnaround";

export interface ReleaseBlockers {
  tasks: TaskBlocker[];
  delays: DelayBlocker[];
  bookings: BookingBlocker[];
}

export interface ReleaseResponse {
  message: string;
  turnaround: FlightTurnaroundDetail;
  tasks: GroundTask[];
  bookings: ResourceBooking[];
  released_count: number;
}

export interface ReleaseErrorPayload {
  code: string;
  message: string;
  blockers?: ReleaseBlockers;
}
