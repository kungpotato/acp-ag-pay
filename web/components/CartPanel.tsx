"use client";

import { useState } from "react";
import { useCart } from "@/lib/cart-context";
import { createCheckoutSession, type CheckoutSession } from "@/lib/checkout";
import { settlePayment } from "@/lib/pay";

export default function CartPanel() {
  const { items, removeItem, totalCents } = useCart();
  const [session, setSession] = useState<CheckoutSession | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [payMessage, setPayMessage] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [paying, setPaying] = useState(false);

  async function handleCheckout() {
    setLoading(true);
    setError(null);
    setPayMessage(null);
    try {
      const created = await createCheckoutSession(
        items.map((i) => ({ product_id: i.productId, quantity: i.quantity })),
      );
      setSession(created);
    } catch (e) {
      setError(e instanceof Error ? e.message : "checkout ไม่สำเร็จ");
    } finally {
      setLoading(false);
    }
  }

  async function handlePay() {
    if (!session) return;
    setPaying(true);
    setPayMessage(null);
    try {
      const result = await settlePayment(session.id);
      if (result.ok) {
        setSession(result.session);
        setPayMessage("สร้าง Stripe PaymentIntent สำเร็จ พร้อมชำระเงินจริง");
      } else if (result.status === 503) {
        setPayMessage(
          "ยังไม่ได้ตั้งค่า STRIPE_SECRET_KEY — ระบบตอบ 503 อย่างสุภาพแทนการ crash (ดูบทที่ 5)",
        );
      } else {
        setPayMessage(`ชำระเงินไม่สำเร็จ: ${result.error}`);
      }
    } finally {
      setPaying(false);
    }
  }

  if (items.length === 0 && !session) return null;

  return (
    <section className="mt-8 rounded-lg border border-zinc-200 p-4 dark:border-zinc-800">
      <h2 className="font-medium text-zinc-900 dark:text-zinc-50">ตะกร้าสินค้า</h2>

      {items.length > 0 && (
        <ul className="mt-3 flex flex-col gap-2">
          {items.map((item) => (
            <li
              key={item.productId}
              className="flex items-center justify-between text-sm"
            >
              <span>
                {item.title} × {item.quantity}
              </span>
              <button
                onClick={() => removeItem(item.productId)}
                className="text-xs text-zinc-500 underline"
              >
                นำออก
              </button>
            </li>
          ))}
        </ul>
      )}

      {items.length > 0 && (
        <div className="mt-3 flex items-center justify-between border-t border-zinc-200 pt-3 dark:border-zinc-800">
          <span className="text-sm text-zinc-600 dark:text-zinc-400">รวม</span>
          <span className="font-medium">
            {(totalCents / 100).toLocaleString("th-TH")} บาท
          </span>
        </div>
      )}

      {items.length > 0 && (
        <button
          onClick={handleCheckout}
          disabled={loading}
          className="mt-4 w-full rounded bg-zinc-900 px-4 py-2 text-sm text-zinc-50 disabled:opacity-50 dark:bg-zinc-100 dark:text-zinc-900"
        >
          {loading ? "กำลังสร้าง checkout session..." : "สร้าง checkout session"}
        </button>
      )}

      {error && <p className="mt-2 text-sm text-red-600">{error}</p>}

      {session && (
        <div className="mt-4 rounded border border-zinc-200 p-3 text-xs dark:border-zinc-800">
          <p>
            session: <code>{session.id}</code>
          </p>
          <p>status: {session.status}</p>
          <p>total: {(session.total_cents / 100).toLocaleString("th-TH")} บาท</p>
          {session.payment_intent_id && (
            <p>
              payment_intent: <code>{session.payment_intent_id}</code>
            </p>
          )}
        </div>
      )}

      {session && session.status === "pending" && (
        <button
          onClick={handlePay}
          disabled={paying}
          className="mt-3 w-full rounded border border-zinc-900 px-4 py-2 text-sm disabled:opacity-50 dark:border-zinc-100"
        >
          {paying ? "กำลัง settle กับ Stripe..." : "ชำระเงินผ่าน Stripe (settle)"}
        </button>
      )}

      {payMessage && (
        <p className="mt-2 text-sm text-zinc-600 dark:text-zinc-400">{payMessage}</p>
      )}
    </section>
  );
}
