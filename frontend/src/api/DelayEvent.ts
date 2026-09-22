import { mockData } from "../mocks/seedData";
import type { DelayEvent } from "../types/DelayEvent";
import { request } from "./request";

const endpoint = "/api/delay-event";

export async function listDelayEvent(): Promise<DelayEvent[]> {
  try {
    return await request<DelayEvent[]>(endpoint);
  } catch {
    return [...mockData.delayEvent];
  }
}

// resolveDelayEvent 关闭延误事件；关闭后未关闭延误归零。
export async function resolveDelayEvent(id: number, actor = "dispatcher"): Promise<DelayEvent> {
  return request<DelayEvent>(`${endpoint}/${id}/resolve`, {
    method: "POST",
    body: JSON.stringify({ actor })
  });
}
