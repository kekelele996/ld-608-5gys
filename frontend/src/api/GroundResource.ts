import { mockData } from "../mocks/seedData";
import type { GroundResourceView } from "../types/GroundResource";

const endpoint = "/api/ground-resource";

export async function listGroundResource(): Promise<GroundResourceView[]> {
  if (typeof fetch !== "undefined" && endpoint.startsWith("/api")) {
    try {
      const res = await fetch(endpoint);
      if (res.ok) return await res.json() as Promise<GroundResourceView[]>;
    } catch {
      // Local mock fallback keeps the UI available during offline review.
    }
  }
  return [...mockData.groundResource];
}

export async function saveGroundResource(payload: GroundResourceView) {
  console.info("save GroundResource", payload);
  return payload;
}
