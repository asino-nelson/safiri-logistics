import { apiUrl } from "./env";

export class ApiError extends Error {
  constructor(message: string, public readonly status?: number) {
    super(message);
    this.name = "ApiError";
  }
}

type RequestOptions = {
  method?: "GET" | "POST" | "PATCH";
  token?: string;
  body?: unknown;
};

export async function request<T>(path: string, options: RequestOptions = {}): Promise<T> {
  const response = await fetch(`${apiUrl}${path}`, {
    method: options.method || "GET",
    headers: {
      Accept: "application/json",
      ...(options.body ? { "Content-Type": "application/json" } : {}),
      ...(options.token ? { Authorization: `Bearer ${options.token}` } : {}),
    },
    body: options.body ? JSON.stringify(options.body) : undefined,
  });

  const text = await response.text();
  const parsed = text ? tryParseJson(text) : undefined;

  if (!response.ok) {
    const message = (parsed && (parsed.error || parsed.message)) || `Request failed with status ${response.status}`;
    throw new ApiError(String(message), response.status);
  }

  return parsed as T;
}

function tryParseJson(text: string) {
  try {
    return JSON.parse(text);
  } catch {
    return text;
  }
}
