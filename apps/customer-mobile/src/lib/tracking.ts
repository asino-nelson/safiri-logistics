import { request } from "./api";
import { TrackingEvent } from "./types";

export async function postTrackingEvent(token: string, loadId: string, input: { latitude: number; longitude: number; speed_kph?: number; heading_degrees?: number }) {
  return request<TrackingEvent>(`/tracking/loads/${loadId}/events`, {
    method: "POST",
    token,
    body: input,
  });
}
