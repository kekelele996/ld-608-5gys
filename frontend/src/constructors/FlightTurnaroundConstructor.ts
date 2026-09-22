import type { FlightTurnaround } from "../types/FlightTurnaround";

export const createDefaultFlightTurnaround = (overrides: Partial<FlightTurnaround> = {}): FlightTurnaround => ({
  id: 0,
  flight_no: "",
  aircraft_reg: "",
  stand_no: "",
  arrival_time: "",
  departure_time: "",
  turnaround_status: "ARRIVING",
  delay_reason: "",
  ready_at: null,
  progress: { total_tasks: 0, signed_tasks: 0, open_delays: 0, active_bookings: 0 },
  gate: { releasable: false, unsigned_tasks: [], open_delays: [], blocking_bookings: [] },
  ...overrides
});

export const createFlightTurnaroundForm = createDefaultFlightTurnaround;
export const createFlightTurnaroundResponse = createDefaultFlightTurnaround;
