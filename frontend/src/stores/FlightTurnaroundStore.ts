import { create } from "zustand";
import { listFlightTurnaround } from "../api/FlightTurnaround";
import type { FlightTurnaround } from "../types/FlightTurnaround";

type State = {
  rows: FlightTurnaround[];
  loading: boolean;
  lastError: string | null;
  load: () => Promise<void>;
  setRows: (rows: FlightTurnaround[]) => void;
};

export const useFlightTurnaroundStore = create<State>((set) => ({
  rows: [],
  loading: false,
  lastError: null,
  async load() {
    set({ loading: true, lastError: null });
    try {
      set({ rows: await listFlightTurnaround(), loading: false });
    } catch (error) {
      set({ loading: false, lastError: (error as Error).message });
    }
  },
  setRows(rows) {
    set({ rows });
  }
}));
