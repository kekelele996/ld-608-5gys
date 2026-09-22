import { mockData } from "../mocks/seedData";
import type { FlightTurnaround, ReleaseResult } from "../types/FlightTurnaround";
import type { ReleaseGate } from "../types/ReleaseGate";
import { request } from "./request";

const endpoint = "/api/flight-turnaround";

export async function listFlightTurnaround(): Promise<FlightTurnaround[]> {
  try {
    return await request<FlightTurnaround[]>(endpoint);
  } catch {
    // Local mock fallback keeps the UI available during offline review.
    return [...mockData.flightTurnaround];
  }
}

export async function getFlightTurnaround(id: number): Promise<FlightTurnaround> {
  return request<FlightTurnaround>(`${endpoint}/${id}`);
}

export async function getReleaseGate(id: number): Promise<ReleaseGate> {
  return request<ReleaseGate>(`${endpoint}/${id}/gate`);
}

// releaseTurnaround 过站放行联动入口。
// 被阻塞时后端返回 409 + 任务/延误/预约三类阻塞清单（错误对象 details）。
export async function releaseTurnaround(id: number, actor = "dispatcher"): Promise<ReleaseResult> {
  return request<ReleaseResult>(`${endpoint}/${id}/release`, {
    method: "POST",
    body: JSON.stringify({ actor })
  });
}
