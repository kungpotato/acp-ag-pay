import type { CheckoutSession } from "./checkout";

const API_BASE_URL =
  process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:3001";

export type PayResult =
  | { ok: true; session: CheckoutSession }
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
  return { ok: true, session: body as CheckoutSession };
}
