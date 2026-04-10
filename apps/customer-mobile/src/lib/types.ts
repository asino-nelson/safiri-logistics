export type Role = "customer" | "driver" | "admin";

export interface UserPayload {
  id: string;
  name: string;
  email: string;
  role: Role;
}

export interface AuthResponse {
  token: string;
  user: UserPayload;
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
  status: string;
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

export interface Payment {
  id: string;
  load_id: string;
  customer_id: string;
  provider: string;
  provider_reference: string;
  phone_number: string;
  amount_kes: number;
  currency: string;
  status: string;
  created_at: string;
  updated_at: string;
}
