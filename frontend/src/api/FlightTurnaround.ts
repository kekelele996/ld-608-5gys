import { mockData } from "../mocks/seedData";
import type { FlightTurnaroundDetail } from "../types/FlightTurnaround";
import type { ReleaseErrorPayload, ReleaseResponse } from "../types/Release";

const endpoint = "/api/flight-turnaround";

export class ReleaseRequestError extends Error {
  payload: ReleaseErrorPayload;
  status: number;

  constructor(payload: ReleaseErrorPayload, status: number) {
    super(payload.message);
    this.name = "ReleaseRequestError";
    this.payload = payload;
    this.status = status;
  }
}

async function parseReleaseFailure(response: Response): Promise<ReleaseRequestError> {
  let payload: ReleaseErrorPayload = { code: "INTERNAL_ERROR", message: "放行失败，请稍后重试" };
  try {
    payload = await response.json() as ReleaseErrorPayload;
  } catch {
    // Keep a safe payload when the backend returns a non-JSON error page.
  }
  return new ReleaseRequestError(payload, response.status);
}

export async function listFlightTurnaround(): Promise<FlightTurnaroundDetail[]> {
  if (typeof fetch !== "undefined" && endpoint.startsWith("/api")) {
    try {
      const res = await fetch(endpoint);
      if (res.ok) return await res.json() as Promise<FlightTurnaroundDetail[]>;
    } catch {
      // Local mock fallback keeps the UI available during offline review.
    }
  }
  return [...(mockData.flightTurnaround as unknown as FlightTurnaroundDetail[])];
}

export async function releaseFlightTurnaround(id: number): Promise<ReleaseResponse> {
  const res = await fetch(`${endpoint}/${id}/release`, {
    method: "POST",
    headers: { "Content-Type": "application/json" }
  });
  if (!res.ok) throw await parseReleaseFailure(res);
  return await res.json() as Promise<ReleaseResponse>;
}

export async function saveFlightTurnaround(payload: FlightTurnaroundDetail) {
  console.info("save FlightTurnaround", payload);
  return payload;
}
