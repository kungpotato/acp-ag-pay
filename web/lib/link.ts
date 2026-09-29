const API_BASE_URL =
  process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:3001";

export type Link = {
  id: string;
  status: "pending" | "active" | "revoked";
  max_amount_cents: number;
  currency: string;
  customer_id: string;
  setup_intent_id: string;
  client_secret?: string;
  payment_method_id?: string;
  created_at: string;
  activated_at?: string;
};

export type LinkResult =
  | { ok: true; link: Link }
  | { ok: false; status: number; error: string };

// createLink starts the real, one-time linking flow (Lesson 9): a Stripe
// Customer + SetupIntent. In a real deployment the returned client_secret
// would be handed to Stripe.js/Elements for the shopper to confirm ONCE.
export async function createLink(maxAmountCents: number, currency = "thb"): Promise<LinkResult> {
  const res = await fetch(new URL("/api/agent/link", API_BASE_URL), {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ max_amount_cents: maxAmountCents, currency }),
  });
  const body = await res.json().catch(() => ({}));
  if (!res.ok) {
    return { ok: false, status: res.status, error: body.error ?? `HTTP ${res.status}` };
  }
  return { ok: true, link: body as Link };
}

// devSeedLink stands in for a shopper completing real Stripe linking, since
// this workshop has no live Stripe account. See server/internal/acp/dev_handler.go.
export async function devSeedLink(maxAmountCents: number, currency = "thb"): Promise<Link> {
  const res = await fetch(new URL("/api/dev/seed_link", API_BASE_URL), {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ max_amount_cents: maxAmountCents, currency }),
  });
  if (!res.ok) throw new Error(`failed to seed link: ${res.status}`);
  return res.json();
}
