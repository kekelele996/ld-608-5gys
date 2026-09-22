import { mockData } from "../mocks/seedData";
import type { ResourceBooking } from "../types/ResourceBooking";
import { request } from "./request";

const endpoint = "/api/resource-booking";

export async function listResourceBooking(): Promise<ResourceBooking[]> {
  try {
    return await request<ResourceBooking[]>(endpoint);
  } catch {
    return [...mockData.resourceBooking];
  }
}

// resolveBookingConflict 解决预约冲突（CONFLICT -> ACTIVE），ACTIVE 由放行事务一次性释放。
export async function resolveBookingConflict(id: number, actor = "dispatcher"): Promise<ResourceBooking> {
  return request<ResourceBooking>(`${endpoint}/${id}/resolve-conflict`, {
    method: "POST",
    body: JSON.stringify({ actor })
  });
}
