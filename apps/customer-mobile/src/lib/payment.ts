import { request } from "./api";
import { Payment } from "./types";

export async function initiateMpesaCheckout(token: string, loadId: string, phone_number: string) {
  return request<Payment>(`/payments/loads/${loadId}/mpesa-checkout`, {
    method: "POST",
    token,
    body: { phone_number },
  });
}
