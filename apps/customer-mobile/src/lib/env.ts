const defaultApiUrl = "http://localhost:8080/api/v1";

function read(value: string | undefined, fallback: string) {
  const trimmed = value?.trim();
  return trimmed && trimmed.length > 0 ? trimmed : fallback;
}

export const apiUrl = read(process.env.EXPO_PUBLIC_API_URL, defaultApiUrl);
export const wsUrl = read(process.env.EXPO_PUBLIC_WS_URL, apiUrl.replace(/^http/, "ws"));
