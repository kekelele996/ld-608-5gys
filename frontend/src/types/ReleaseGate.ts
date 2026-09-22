export interface TaskBlocker {
  id: number;
  task_type: string;
  team_id: number;
  status: string;
  deadline: string;
  blocker_note: string;
  reason: string;
}

export interface DelayBlocker {
  id: number;
  delay_type: string;
  minutes: number;
  responsibility_team: string;
  root_cause: string;
  reason: string;
}

export interface BookingBlocker {
  id: number;
  resource_id: number;
  resource_code: string;
  task_id: number;
  booking_status: string;
  conflict_reason: string;
  reason: string;
}

// ReleaseGate 放行门禁快照：看板阻塞明细直接消费。
export interface ReleaseGate {
  releasable: boolean;
  unsigned_tasks: TaskBlocker[];
  open_delays: DelayBlocker[];
  blocking_bookings: BookingBlocker[];
}

export interface TurnaroundProgress {
  total_tasks: number;
  signed_tasks: number;
  open_delays: number;
  active_bookings: number;
}
