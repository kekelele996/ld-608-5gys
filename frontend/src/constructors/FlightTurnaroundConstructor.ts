import type { FlightTurnaround } from "../types/FlightTurnaround";

export const createDefaultFlightTurnaround = (overrides: Partial<FlightTurnaround> = {}): FlightTurnaround => ({
  id: 0,
  flight_no: "",
  aircraft_reg: "",
  stand_no: "",
  arrival_time: "",
  departure_time: "",
  turnaround_status: "ON_STAND",
  delay_reason: "",
  ready_at: null,
  timeline_synced_at: null,
  ...overrides
});

export const createFlightTurnaroundForm = createDefaultFlightTurnaround;
export const createFlightTurnaroundResponse = createDefaultFlightTurnaround;
