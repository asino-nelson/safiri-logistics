import { request } from "./api";
import { Load } from "./types";

export type CreateLoadInput = {
  title: string;
  description: string;
  origin: string;
  destination: string;
  weight_kg: number;
  pickup_latitude: number;
  pickup_longitude: number;
  dropoff_latitude: number;
  dropoff_longitude: number;
  pickup_at?: string;
  priority?: number;
  cargo_type?: string;
  equipment_type?: string;
};

type LoadListResponse = {
  data: Load[];
};

export async function listLoads(token: string) {
  const response = await request<LoadListResponse>("/loads", { token });
  return response.data;
}

export async function getLoadById(token: string, id: string) {
  const loads = await listLoads(token);
  const load = loads.find((entry) => entry.id === id);
  if (!load) {
    throw new Error("Load not found");
  }
  return load;
}

export async function createLoad(token: string, input: CreateLoadInput) {
  return request<Load>("/loads", {
    method: "POST",
    token,
    body: input,
  });
}
