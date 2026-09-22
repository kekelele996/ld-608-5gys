export interface FlightTurnaround {
  id: number;
  flight_no: string;
  aircraft_reg: string;
  stand_no: string;
  arrival_time: string;
  departure_time: string;
  turnaround_status: string;
  delay_reason: string;
  ready_at?: string | null;
  timeline_synced_at?: string | null;
}

export interface TimelineEvent {
  code: string;
  label: string;
  occurred_at: string;
  status: string;
}

export interface FlightTurnaroundDetail extends FlightTurnaround {
  task_completion_rate: number;
  timeline: TimelineEvent[];
}
