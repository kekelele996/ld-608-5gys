import { create } from "zustand";
import {
  listFlightTurnaround,
  releaseTurnaround
} from "../api/FlightTurnaround";
import { signGroundTask } from "../api/GroundTask";
import { resolveDelayEvent } from "../api/DelayEvent";
import {
  listResourceBooking,
  resolveBookingConflict
} from "../api/ResourceBooking";
import { listGroundTask } from "../api/GroundTask";
import { listDelayEvent } from "../api/DelayEvent";
import { listGroundResource } from "../api/GroundResource";
import { useFlightTurnaroundStore } from "./FlightTurnaroundStore";
import { useGroundTaskStore } from "./GroundTaskStore";
import { useDelayEventStore } from "./DelayEventStore";
import { useResourceBookingStore } from "./ResourceBookingStore";
import { useGroundResourceStore } from "./GroundResourceStore";
import type { ApiError } from "../types/Api";
import type { ReleaseResult } from "../types/FlightTurnaround";

// 过站放行联动编排：任一写操作成功后刷新航班/任务/延误/预约/资源，看板同步。
type ReleaseState = {
  releasingIds: Record<number, boolean>;
  signingTaskIds: Record<number, boolean>;
  resolvingDelayIds: Record<number, boolean>;
  resolvingBookingIds: Record<number, boolean>;
  // 最近一次被阻塞的放行错误（含三类阻塞清单），看板弹窗展示。
  lastBlocked: { turnaroundId: number; error: ApiError } | null;
  notice: string | null;
  refreshAll: () => Promise<void>;
  release: (id: number) => Promise<ReleaseResult | null>;
  signTask: (id: number) => Promise<void>;
  resolveDelay: (id: number) => Promise<void>;
  resolveBooking: (id: number) => Promise<void>;
  clearBlocked: () => void;
  clearNotice: () => void;
};

async function refreshAll() {
  const [turnarounds, tasks, delays, bookings, resources] = await Promise.all([
    listFlightTurnaround(),
    listGroundTask(),
    listDelayEvent(),
    listResourceBooking(),
    listGroundResource()
  ]);
  useFlightTurnaroundStore.getState().setRows(turnarounds);
  useGroundTaskStore.getState().setRows(tasks);
  useDelayEventStore.getState().setRows(delays);
  useResourceBookingStore.getState().setRows(bookings);
  useGroundResourceStore.setState({ rows: resources });
}

export const useReleaseStore = create<ReleaseState>((set, get) => ({
  releasingIds: {},
  signingTaskIds: {},
  resolvingDelayIds: {},
  resolvingBookingIds: {},
  lastBlocked: null,
  notice: null,

  async refreshAll() {
    await refreshAll();
  },

  // release：重复/并发提交由后端门禁保证只生效一次；前端对 in-flight 的放行按钮置灰。
  async release(id) {
    if (get().releasingIds[id]) return null;
    set({ releasingIds: { ...get().releasingIds, [id]: true }, lastBlocked: null, notice: null });
    try {
      const result = await releaseTurnaround(id);
      await refreshAll();
      set({
        notice: `航班 ${result.flight_no} 已放行（READY），一次性释放 ACTIVE 预约 ${result.released_bookings} 份`
      });
      return result;
    } catch (error) {
      const apiError = error as ApiError;
      // 阻塞拒绝 / 重复放行 / 并发落败都把错误抛给看板展示明细。
      set({ lastBlocked: { turnaroundId: id, error: apiError }, notice: apiError.message });
      throw apiError;
    } finally {
      set({ releasingIds: { ...get().releasingIds, [id]: false } });
    }
  },

  async signTask(id) {
    if (get().signingTaskIds[id]) return;
    set({ signingTaskIds: { ...get().signingTaskIds, [id]: true } });
    try {
      await signGroundTask(id);
      await refreshAll();
      set({ notice: `任务 #${id} 已签收，任务与过站时间线已同步刷新` });
    } finally {
      set({ signingTaskIds: { ...get().signingTaskIds, [id]: false } });
    }
  },

  async resolveDelay(id) {
    if (get().resolvingDelayIds[id]) return;
    set({ resolvingDelayIds: { ...get().resolvingDelayIds, [id]: true } });
    try {
      await resolveDelayEvent(id);
      await refreshAll();
      set({ notice: `延误事件 #${id} 已关闭，未关闭延误计数已刷新` });
    } finally {
      set({ resolvingDelayIds: { ...get().resolvingDelayIds, [id]: false } });
    }
  },

  async resolveBooking(id) {
    if (get().resolvingBookingIds[id]) return;
    set({ resolvingBookingIds: { ...get().resolvingBookingIds, [id]: true } });
    try {
      await resolveBookingConflict(id);
      await refreshAll();
      set({ notice: `预约 #${id} 冲突已解决，将在放行时随 ACTIVE 预约一并释放` });
    } finally {
      set({ resolvingBookingIds: { ...get().resolvingBookingIds, [id]: false } });
    }
  },

  clearBlocked() {
    set({ lastBlocked: null });
  },
  clearNotice() {
    set({ notice: null });
  }
}));
