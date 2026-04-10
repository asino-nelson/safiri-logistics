export type Role = "customer" | "driver" | "admin";

export interface AuthUser {
  id: string;
  name: string;
  email: string;
  role: Role;
}

export interface AuthResponse {
  token: string;
  user: AuthUser;
}

export interface KycProfile {
  user_id: string;
  national_id: string;
  license_number: string;
  truck_registration: string;
  status: "pending" | "approved" | "rejected";
  rejection_reason?: string | null;
  submitted_at: string;
  reviewed_at?: string | null;
  verified_at?: string | null;
  years_experience: number;
  max_load_kg: number;
  equipment_type: string;
  is_online: boolean;
  is_available: boolean;
  current_latitude?: number | null;
  current_longitude?: number | null;
  last_location_at?: string | null;
}

export interface Load {
  id: string;
  poster_id: string;
  assigned_driver_id?: string | null;
  title: string;
  description: string;
  origin: string;
  destination: string;
  weight_kg: number;
  status: "posted" | "matched" | "picked" | "in_transit" | "delivered";
  pickup_latitude: number;
  pickup_longitude: number;
  dropoff_latitude: number;
  dropoff_longitude: number;
  pickup_at: string;
  priority: number;
  cargo_type: string;
  equipment_type: string;
  quoted_price_kes: number;
  matched_at?: string | null;
  matching_score?: number | null;
  assignment_source: string;
  created_at: string;
  updated_at: string;
  picked_at?: string | null;
  delivered_at?: string | null;
}

export interface TrackingEvent {
  id: string;
  load_id: string;
  driver_id: string;
  latitude: number;
  longitude: number;
  speed_kph: number;
  heading_degrees: number;
  recorded_at: string;
}

export interface LoginInput {
  email: string;
  password: string;
}

export interface RegisterInput {
  name: string;
  email: string;
  password: string;
  role: "driver";
}

export interface SubmitKycInput {
  national_id: string;
  license_number: string;
  truck_registration: string;
}

export interface UpdateOperationsInput {
  years_experience: number;
  max_load_kg: number;
  equipment_type: string;
  is_online: boolean;
  is_available: boolean;
  latitude?: number | null;
  longitude?: number | null;
}

export interface CreateTrackingEventInput {
  latitude: number;
  longitude: number;
  speed_kph: number;
  heading_degrees: number;
}

export interface ApiListResponse<T> {
  data: T[];
}
