import type { CheckoutSession } from "./checkout";

const API_BASE_URL =
  process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:3001";

export type Order = {
  id: string;
  session_id: string;
  status: "pending_confirmation" | "paid" | "failed";
  total_cents: number;
  currency: string;
  payment_intent_id: string;
  created_at: string;
  confirmed_at?: string;
};

export type PayResult =
  | { ok: true; session: CheckoutSession; order: Order }
  | { ok: false; status: number; error: string };

export async function settlePayment(sessionId: string): Promise<PayResult> {
  const res = await fetch(new URL("/api/agent/pay", API_BASE_URL), {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ session_id: sessionId }),
  });
  const body = await res.json().catch(() => ({}));
  if (!res.ok) {
    return { ok: false, status: res.status, error: body.error ?? `HTTP ${res.status}` };
  }
  return { ok: true, session: body.session as CheckoutSession, order: body.order as Order };
}

export async function fetchOrder(orderId: string): Promise<Order> {
  const res = await fetch(new URL(`/api/orders/${orderId}`, API_BASE_URL), {
    cache: "no-store",
  });
  if (!res.ok) throw new Error(`failed to fetch order: ${res.status}`);
  return res.json();
}
