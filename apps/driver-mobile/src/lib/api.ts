import { apiUrl, wsUrl } from "@/lib/env";
import {
  ApiListResponse,
  AuthResponse,
  AuthUser,
  CreateTrackingEventInput,
  KycProfile,
  Load,
  LoginInput,
  SubmitKycInput,
  TrackingEvent,
  UpdateOperationsInput,
} from "@/types/api";

class ApiError extends Error {
  constructor(
    message: string,
    public status: number,
    public payload?: unknown,
  ) {
    super(message);
    this.name = "ApiError";
  }
}

async function request<T>(
  path: string,
  options: RequestInit = {},
  token?: string | null,
): Promise<T> {
  const headers = new Headers(options.headers);
  headers.set("Accept", "application/json");

  if (options.body && !headers.has("Content-Type")) {
    headers.set("Content-Type", "application/json");
  }

  if (token) {
    headers.set("Authorization", `Bearer ${token}`);
  }

  const response = await fetch(apiUrl(path), {
    ...options,
    headers,
  });

  const text = await response.text();
  const payload = text ? JSON.parse(text) : null;

  if (!response.ok) {
    throw new ApiError(
      payload?.error ?? `Request failed with status ${response.status}`,
      response.status,
      payload,
    );
  }

  return payload as T;
}

export const api = {
  login(input: LoginInput) {
    return request<AuthResponse>("/api/v1/auth/login", {
      method: "POST",
      body: JSON.stringify(input),
    });
  },
  me(token: string) {
    return request<AuthUser>("/api/v1/auth/me", { method: "GET" }, token);
  },
  getKycProfile(token: string) {
    return request<KycProfile>("/api/v1/drivers/kyc/me", { method: "GET" }, token);
  },
  submitKyc(token: string, input: SubmitKycInput) {
    return request<KycProfile>("/api/v1/drivers/kyc", {
      method: "POST",
      body: JSON.stringify(input),
    }, token);
  },
  updateOperations(token: string, input: UpdateOperationsInput) {
    return request<KycProfile>("/api/v1/drivers/me/operations", {
      method: "PATCH",
      body: JSON.stringify(input),
    }, token);
  },
  listLoads(token: string) {
    return request<ApiListResponse<Load>>("/api/v1/loads", { method: "GET" }, token);
  },
  pickLoad(token: string, loadId: string) {
    return request<Load>(`/api/v1/loads/${loadId}/pick`, { method: "POST" }, token);
  },
  updateLoadStatus(token: string, loadId: string, status: Load["status"]) {
    return request<Load>(`/api/v1/loads/${loadId}/status`, {
      method: "PATCH",
      body: JSON.stringify({ status }),
    }, token);
  },
  createTrackingEvent(token: string, loadId: string, input: CreateTrackingEventInput) {
    return request<TrackingEvent>(`/api/v1/tracking/loads/${loadId}/events`, {
      method: "POST",
      body: JSON.stringify(input),
    }, token);
  },
  createCheckout(token: string, loadId: string, phoneNumber: string) {
    return request(`/api/v1/payments/loads/${loadId}/mpesa-checkout`, {
      method: "POST",
      body: JSON.stringify({ phone_number: phoneNumber }),
    }, token);
  },
  buildTrackingSocketUrl(loadId: string) {
    return wsUrl(`/api/v1/tracking/ws?load_id=${encodeURIComponent(loadId)}`);
  },
};

export { ApiError };
