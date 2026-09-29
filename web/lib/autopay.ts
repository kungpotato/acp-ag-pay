import type { CheckoutSession } from "./checkout";
import type { Order } from "./pay";

const API_BASE_URL =
  process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:3001";

export type AutopayResult =
  | { ok: true; session: CheckoutSession; order: Order }
  | { ok: false; status: number; error: string };

// autopay is the point of Lesson 9: the agent settles a session using an
// already-linked payment method. No client_secret, no confirm-on-client
// step, no button the shopper has to click per order.
export async function autopay(sessionId: string, linkId: string): Promise<AutopayResult> {
  const res = await fetch(new URL("/api/agent/autopay", API_BASE_URL), {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ session_id: sessionId, link_id: linkId }),
  });
  const body = await res.json().catch(() => ({}));
  if (!res.ok) {
    return { ok: false, status: res.status, error: body.error ?? `HTTP ${res.status}` };
  }
  return { ok: true, session: body.session as CheckoutSession, order: body.order as Order };
}
