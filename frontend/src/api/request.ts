import type { ApiEnvelope, ApiError } from "../types/Api";

// request 统一解包后端 {ok,data,error} 外层结构；失败抛出携带明细的 ApiError。
export async function request<T>(input: string, init?: RequestInit): Promise<T> {
  const res = await fetch(input, {
    headers: { "Content-Type": "application/json", ...(init?.headers ?? {}) },
    ...init
  });
  let envelope: ApiEnvelope<T> | null = null;
  try {
    envelope = (await res.json()) as ApiEnvelope<T>;
  } catch {
    envelope = null;
  }
  if (!res.ok || !envelope?.ok || envelope.error) {
    const error: ApiError = envelope?.error ?? {
      code: "INTERNAL_ERROR",
      message: `请求失败：HTTP ${res.status}`
    };
    throw error;
  }
  return envelope.data as T;
}
