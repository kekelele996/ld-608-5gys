import { mockData } from "../mocks/seedData";
import type { GroundResource } from "../types/GroundResource";
import { request } from "./request";

const endpoint = "/api/ground-resource";

export async function listGroundResource(): Promise<GroundResource[]> {
  try {
    return await request<GroundResource[]>(endpoint);
  } catch {
    return [...mockData.groundResource];
  }
}
