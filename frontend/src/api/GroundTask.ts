import { mockData } from "../mocks/seedData";
import type { GroundTask } from "../types/GroundTask";
import { request } from "./request";

const endpoint = "/api/ground-task";

export async function listGroundTask(): Promise<GroundTask[]> {
  try {
    return await request<GroundTask[]>(endpoint);
  } catch {
    return [...mockData.groundTask];
  }
}

// signGroundTask 任务签收（SIGNED），与过站时间线同步刷新。
export async function signGroundTask(id: number, actor = "dispatcher"): Promise<GroundTask> {
  return request<GroundTask>(`${endpoint}/${id}/sign`, {
    method: "POST",
    body: JSON.stringify({ actor })
  });
}
