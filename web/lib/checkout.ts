export type CheckoutSession = {
  id: string;
  status: string;
  line_items: {
    product_id: string;
    title: string;
    quantity: number;
    unit_price_cents: number;
  }[];
  total_cents: number;
  currency: string;
  created_at: string;
  payment_intent_id?: string;
  client_secret?: string;
};

const API_BASE_URL =
  process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:3001";

export async function createCheckoutSession(
  items: { product_id: string; quantity: number }[],
): Promise<CheckoutSession> {
  const res = await fetch(new URL("/api/acp/checkout_sessions", API_BASE_URL), {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ items }),
  });
  if (!res.ok) {
    const body = await res.json().catch(() => ({}));
    throw new Error(body.error ?? `checkout session failed: ${res.status}`);
  }
  return res.json();
}
