import { create } from "zustand";
import { listGroundTask } from "../api/GroundTask";
import type { GroundTask } from "../types/GroundTask";

type State = {
  rows: GroundTask[];
  loading: boolean;
  load: () => Promise<void>;
  setRows: (rows: GroundTask[]) => void;
};

export const useGroundTaskStore = create<State>((set) => ({
  rows: [],
  loading: false,
  async load() {
    set({ loading: true });
    set({ rows: await listGroundTask(), loading: false });
  },
  setRows(rows) {
    set({ rows });
  }
}));
