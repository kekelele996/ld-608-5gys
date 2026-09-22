import { useMemo } from "react";

import type { ResourceBooking } from "../types/ResourceBooking";

// useResourceConflict 统计 ACTIVE/CONFLICT 预约，资源页与看板共用。
export function useResourceConflict(bookings: ResourceBooking[] = []) {
  return useMemo(() => {
    const active = bookings.filter((booking) => booking.booking_status === "ACTIVE");
    const conflicts = bookings.filter((booking) => booking.booking_status === "CONFLICT");
    const released = bookings.filter((booking) => booking.booking_status === "RELEASED");
    return {
      active,
      conflicts,
      released,
      activeCount: active.length,
      conflictCount: conflicts.length,
      releasedCount: released.length
    };
  }, [bookings]);
}
