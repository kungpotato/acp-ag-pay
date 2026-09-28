"use client";

import { useEffect, useState } from "react";
import { fetchOrder, type Order } from "@/lib/pay";

const STATUS_LABEL: Record<Order["status"], string> = {
  pending_confirmation: "รอการยืนยันจาก Stripe",
  paid: "ชำระเงินสำเร็จ",
  failed: "ชำระเงินไม่สำเร็จ",
};

const POLL_INTERVAL_MS = 3000;

export default function OrderStatus({ orderId }: { orderId: string }) {
  const [order, setOrder] = useState<Order | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let cancelled = false;
    let timer: ReturnType<typeof setTimeout>;

    async function poll() {
      try {
        const latest = await fetchOrder(orderId);
        if (cancelled) return;
        setOrder(latest);
        setError(null);
        // Lesson 7's webhook is what eventually moves this out of
        // pending_confirmation; keep polling until it does.
        if (latest.status === "pending_confirmation") {
          timer = setTimeout(poll, POLL_INTERVAL_MS);
        }
      } catch {
        if (!cancelled) setError("ไม่สามารถตรวจสอบสถานะคำสั่งซื้อได้");
      }
    }

    poll();
    return () => {
      cancelled = true;
      clearTimeout(timer);
    };
  }, [orderId]);

  if (error) return <p className="text-sm text-red-600">{error}</p>;
  if (!order) return <p className="text-sm text-zinc-500">กำลังโหลดสถานะคำสั่งซื้อ...</p>;

  return (
    <div className="mt-2 rounded border border-zinc-200 p-3 text-xs dark:border-zinc-800">
      <p>
        order: <code>{order.id}</code>
      </p>
      <p className="mt-1 font-medium text-zinc-900 dark:text-zinc-50">
        {STATUS_LABEL[order.status]}
        {order.status === "pending_confirmation" && (
          <span className="ml-1 animate-pulse text-zinc-400">●</span>
        )}
      </p>
      <p className="mt-1 text-zinc-500">
        ยอดรวม {(order.total_cents / 100).toLocaleString("th-TH")} บาท
      </p>
    </div>
  );
}
