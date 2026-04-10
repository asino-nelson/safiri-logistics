const rawBaseUrl = process.env.EXPO_PUBLIC_API_BASE_URL ?? "http://localhost:8080";

export const API_BASE_URL = rawBaseUrl.replace(/\/+$/, "");

export function apiUrl(path: string) {
  return `${API_BASE_URL}${path.startsWith("/") ? path : `/${path}`}`;
}

export function wsUrl(path: string) {
  return apiUrl(path).replace(/^http/i, "ws");
}
