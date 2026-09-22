import type { ResourceBooking } from "../types/ResourceBooking";

export const createDefaultResourceBooking = (overrides: Partial<ResourceBooking> = {}): ResourceBooking => ({
  id: 0,
  resource_id: 0,
  resource_code: "",
  turnaround_id: 0,
  task_id: 0,
  start_time: "",
  end_time: "",
  booking_status: "ACTIVE",
  conflict_reason: "",
  released_at: null,
  ...overrides
});

export const createResourceBookingForm = createDefaultResourceBooking;
export const createResourceBookingResponse = createDefaultResourceBooking;
