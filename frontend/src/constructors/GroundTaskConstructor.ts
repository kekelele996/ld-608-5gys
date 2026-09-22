import type { GroundTask } from "../types/GroundTask";

export const createDefaultGroundTask = (overrides: Partial<GroundTask> = {}): GroundTask => ({
  id: 0,
  turnaround_id: 0,
  task_type: "CLEANING",
  team_id: 0,
  planned_start: "",
  deadline: "",
  actual_finish: null,
  status: "DISPATCHED",
  blocker_note: "",
  timeline_synced_at: null,
  ...overrides
});

export const createGroundTaskForm = createDefaultGroundTask;
export const createGroundTaskResponse = createDefaultGroundTask;
