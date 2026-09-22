import { create } from "zustand";
import { listFlightTurnaround, releaseFlightTurnaround, ReleaseRequestError } from "../api/FlightTurnaround";
import { useGroundTaskStore } from "./GroundTaskStore";
import { useResourceBookingStore } from "./ResourceBookingStore";
import { useGroundResourceStore } from "./GroundResourceStore";
import type { FlightTurnaroundDetail } from "../types/FlightTurnaround";
import type { ReleaseBlockers } from "../types/Release";

type State = {
  rows: FlightTurnaroundDetail[];
  loading: boolean;
  releasingId: number | null;
  blockersByFlight: Record<number, ReleaseBlockers>;
  notice: { type: "success" | "error"; message: string } | null;
  load: () => Promise<void>;
  release: (id: number) => Promise<boolean>;
  clearNotice: () => void;
};

const emptyBlockers = (): ReleaseBlockers => ({ tasks: [], delays: [], bookings: [] });

export const useFlightTurnaroundStore = create<State>((set, get) => ({
  rows: [],
  loading: false,
  releasingId: null,
  blockersByFlight: {},
  notice: null,
  async load() {
    set({ loading: true });
    const rows = await listFlightTurnaround();
    set({ rows, loading: false });
  },
  async release(id) {
    if (get().releasingId !== null) {
      set({ notice: { type: "error", message: "放行请求正在处理，请勿重复提交" } });
      return false;
    }
    set({ releasingId: id, notice: null });
    try {
      const result = await releaseFlightTurnaround(id);
      set((state) => ({
        rows: state.rows.map((row) => row.id === id ? result.turnaround : row),
        blockersByFlight: { ...state.blockersByFlight, [id]: emptyBlockers() },
        notice: { type: "success", message: `${result.turnaround.flight_no} 已进入 READY，释放 ${result.released_count} 条 ACTIVE 预约` },
        releasingId: null
      }));
      useGroundTaskStore.setState((state) => ({
        rows: state.rows.map((task) => result.tasks.find((updated) => updated.id === task.id) ?? task)
      }));
      useResourceBookingStore.setState((state) => ({
        rows: state.rows.map((booking) => result.bookings.find((updated) => updated.id === booking.id) ?? booking)
      }));
      void useGroundResourceStore.getState().load();
      return true;
    } catch (error) {
      const requestError = error instanceof ReleaseRequestError ? error : null;
      set((state) => ({
        blockersByFlight: { ...state.blockersByFlight, [id]: requestError?.payload.blockers ?? emptyBlockers() },
        notice: { type: "error", message: requestError?.payload.message ?? "放行失败，请稍后重试" },
        releasingId: null
      }));
      return false;
    }
  },
  clearNotice() {
    set({ notice: null });
  }
}));
