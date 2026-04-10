import { request } from "./api";
import { AuthResponse, UserPayload } from "./types";

export type SignUpInput = {
  name: string;
  email: string;
  password: string;
};

export async function login(email: string, password: string) {
  return request<AuthResponse>("/auth/login", {
    method: "POST",
    body: { email, password },
  });
}

export async function register(input: SignUpInput) {
  return request<AuthResponse>("/auth/register", {
    method: "POST",
    body: { ...input, role: "customer" },
  });
}

export async function me(token: string) {
  return request<UserPayload>("/auth/me", { token });
}
