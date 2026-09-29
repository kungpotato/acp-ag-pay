"use client";

import { useState } from "react";
import { createLink, devSeedLink } from "@/lib/link";
import { useLink } from "@/lib/link-context";

const MAX_AMOUNT_CENTS = 500000; // 5,000 บาท: the mandate the shopper is granting

export default function AgentWallet() {
  const { link, setLink } = useLink();
  const [message, setMessage] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  async function handleRealLink() {
    setLoading(true);
    setMessage(null);
    try {
      const result = await createLink(MAX_AMOUNT_CENTS);
      if (result.ok) {
        setLink(result.link);
        setMessage("สร้าง SetupIntent สำเร็จ — ในระบบจริงจะแสดง Stripe Elements ให้ยืนยันตรงนี้");
      } else if (result.status === 503) {
        setMessage(
          "ยังไม่ได้ตั้งค่า STRIPE_SECRET_KEY — ระบบตอบ 503 อย่างสุภาพเหมือนบทที่ 5",
        );
      } else {
        setMessage(`เชื่อมต่อไม่สำเร็จ: ${result.error}`);
      }
    } finally {
      setLoading(false);
    }
  }

  async function handleSimulateLink() {
    setLoading(true);
    setMessage(null);
    try {
      const seeded = await devSeedLink(MAX_AMOUNT_CENTS);
      setLink(seeded);
      setMessage("จำลองว่าผู้ใช้ยืนยันตัวตนกับ Stripe Link สำเร็จแล้ว (ครั้งเดียว)");
    } catch (e) {
      setMessage(e instanceof Error ? e.message : "จำลองไม่สำเร็จ");
    } finally {
      setLoading(false);
    }
  }

  return (
    <section className="mt-8 rounded-lg border border-zinc-200 p-4 dark:border-zinc-800">
      <h2 className="font-medium text-zinc-900 dark:text-zinc-50">
        เชื่อมต่อวิธีชำระเงิน (ครั้งเดียว)
      </h2>
      <p className="mt-1 text-sm text-zinc-500">
        บทที่ 9: ยืนยันตัวตนครั้งนี้ครั้งเดียว agent จะจ่ายเงินแทนได้เองในทุกคำสั่งซื้อถัดไป
        โดยไม่ต้องกดปุ่มจ่ายเงินอีก (วงเงินไม่เกิน {(MAX_AMOUNT_CENTS / 100).toLocaleString("th-TH")}{" "}
        บาท)
      </p>

      {link ? (
        <div className="mt-4 rounded border border-green-600 p-3 text-xs">
          <p className="font-medium text-green-700 dark:text-green-400">
            เชื่อมต่อแล้ว ✓ agent พร้อมจ่ายเงินอัตโนมัติ
          </p>
          <p className="mt-1">
            link: <code>{link.id}</code>
          </p>
          <p>วงเงิน: {(link.max_amount_cents / 100).toLocaleString("th-TH")} บาท</p>
        </div>
      ) : (
        <div className="mt-4 flex flex-col gap-2 sm:flex-row">
          <button
            onClick={handleRealLink}
            disabled={loading}
            className="flex-1 rounded bg-zinc-900 px-4 py-2 text-sm text-zinc-50 disabled:opacity-50 dark:bg-zinc-100 dark:text-zinc-900"
          >
            เชื่อมต่อผ่าน Stripe จริง
          </button>
          <button
            onClick={handleSimulateLink}
            disabled={loading}
            className="flex-1 rounded border border-zinc-300 px-4 py-2 text-sm disabled:opacity-50 dark:border-zinc-700"
          >
            (dev) จำลองว่าเชื่อมต่อสำเร็จแล้ว
          </button>
        </div>
      )}

      {message && <p className="mt-2 text-sm text-zinc-600 dark:text-zinc-400">{message}</p>}
    </section>
  );
}
