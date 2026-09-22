import type { ReleaseGate } from "./ReleaseGate";

// 放行被拒绝时后端返回的错误体，details 携带三类阻塞清单。
export interface ApiError {
  code: string;
  message: string;
  details?: {
    turnaround_id?: number;
    unsigned_tasks?: ReleaseGate["unsigned_tasks"];
    open_delays?: ReleaseGate["open_delays"];
    blocking_bookings?: ReleaseGate["blocking_bookings"];
  };
}

export interface ApiEnvelope<T> {
  ok: boolean;
  data?: T;
  error?: ApiError;
}
