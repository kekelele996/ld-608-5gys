import { useMemo } from "react";

import type { FlightTurnaround } from "../types/FlightTurnaround";

// useTurnaroundProgress 计算航班任务签收进度与放行门禁阻塞数，看板与过站页共用。
export function useTurnaroundProgress(turnaround: FlightTurnaround | null | undefined) {
  return useMemo(() => {
    if (!turnaround) {
      return {
        total: 0,
        signed: 0,
        percent: 0,
        openDelays: 0,
        activeBookings: 0,
        unsignedTasks: 0,
        blockingBookings: 0,
        releasable: false,
        isReady: false
      };
    }
    const { progress, gate } = turnaround;
    const percent = progress.total_tasks === 0
      ? 0
      : Math.round((progress.signed_tasks / progress.total_tasks) * 100);
    return {
      total: progress.total_tasks,
      signed: progress.signed_tasks,
      percent,
      openDelays: progress.open_delays,
      activeBookings: progress.active_bookings,
      unsignedTasks: gate.unsigned_tasks.length,
      blockingBookings: gate.blocking_bookings.length,
      releasable: gate.releasable,
      isReady: turnaround.turnaround_status === "READY" || turnaround.turnaround_status === "DEPARTED"
    };
  }, [turnaround]);
}
