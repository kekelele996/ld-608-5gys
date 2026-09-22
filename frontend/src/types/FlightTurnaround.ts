import type { ReleaseGate, TurnaroundProgress } from "./ReleaseGate";

export interface FlightTurnaround {
  id: number;
  flight_no: string;
  aircraft_reg: string;
  stand_no: string;
  arrival_time: string;
  departure_time: string;
  turnaround_status: string;
  delay_reason: string;
  ready_at: string | null;
  progress: TurnaroundProgress;
  gate: ReleaseGate;
}

export interface ReleaseResult {
  turnaround_id: number;
  flight_no: string;
  turnaround_status: string;
  ready_at: string;
  released_bookings: number;
  gate: ReleaseGate;
}
