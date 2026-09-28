"use client";

import { useState } from "react";
import { sendChatMessage, type ChatReply } from "@/lib/agent";
import { formatPrice } from "@/lib/api";
import { useCart } from "@/lib/cart-context";

type ChatEntry = { role: "user" | "agent"; reply?: ChatReply; text?: string };

export default function ChatPanel() {
  const { addBook } = useCart();
  const [input, setInput] = useState("");
  const [entries, setEntries] = useState<ChatEntry[]>([]);
  const [loading, setLoading] = useState(false);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    const text = input.trim();
    if (!text || loading) return;

    setEntries((prev) => [...prev, { role: "user", text }]);
    setInput("");
    setLoading(true);
    try {
      const reply = await sendChatMessage(text);
      setEntries((prev) => [...prev, { role: "agent", reply }]);
    } catch {
      setEntries((prev) => [
        ...prev,
        {
          role: "agent",
          reply: {
            message: "เชื่อมต่อกับ shopping agent ไม่สำเร็จ ลองใหม่อีกครั้งค่ะ",
            intent: { kind: "unknown" },
          },
        },
      ]);
    } finally {
      setLoading(false);
    }
  }

  return (
    <section className="mt-12 rounded-lg border border-zinc-200 p-4 dark:border-zinc-800">
      <h2 className="font-medium text-zinc-900 dark:text-zinc-50">
        คุยกับ shopping agent
      </h2>
      <p className="mt-1 text-sm text-zinc-500">
        บทที่ 2: ลองพิมพ์ &quot;หาแฮร์รี่ พอตเตอร์&quot; หรือ &quot;find clean code&quot;
      </p>

      <div className="mt-4 flex flex-col gap-3">
        {entries.map((entry, i) =>
          entry.role === "user" ? (
            <p key={i} className="self-end rounded-lg bg-zinc-900 px-3 py-2 text-sm text-zinc-50 dark:bg-zinc-100 dark:text-zinc-900">
              {entry.text}
            </p>
          ) : (
            <div key={i} className="self-start max-w-full">
              <p className="rounded-lg bg-zinc-100 px-3 py-2 text-sm text-zinc-900 dark:bg-zinc-800 dark:text-zinc-50">
                {entry.reply?.message}
              </p>
              {entry.reply?.books && entry.reply.books.length > 0 && (
                <ul className="mt-2 grid grid-cols-1 gap-2 sm:grid-cols-2">
                  {entry.reply.books.map((book) => (
                    <li
                      key={book.id}
                      className="rounded border border-zinc-200 p-2 text-xs dark:border-zinc-800"
                    >
                      <p className="font-medium text-zinc-900 dark:text-zinc-50">
                        {book.title}
                      </p>
                      <p className="text-zinc-500">{book.author}</p>
                      <p className="mt-1 text-zinc-700 dark:text-zinc-300">
                        {formatPrice(book)}
                      </p>
                      <button
                        onClick={() => addBook(book)}
                        disabled={book.stock <= 0}
                        className="mt-2 w-full rounded border border-zinc-300 px-2 py-1 text-xs disabled:opacity-40 dark:border-zinc-700"
                      >
                        เพิ่มลงตะกร้า
                      </button>
                    </li>
                  ))}
                </ul>
              )}
            </div>
          ),
        )}
        {loading && (
          <p className="self-start text-sm text-zinc-500">agent กำลังพิมพ์...</p>
        )}
      </div>

      <form onSubmit={handleSubmit} className="mt-4 flex gap-2">
        <input
          value={input}
          onChange={(e) => setInput(e.target.value)}
          placeholder="พิมพ์สิ่งที่อยากได้..."
          className="flex-1 rounded border border-zinc-300 px-3 py-2 text-sm dark:border-zinc-700 dark:bg-zinc-900"
        />
        <button
          type="submit"
          disabled={loading}
          className="rounded bg-zinc-900 px-4 py-2 text-sm text-zinc-50 disabled:opacity-50 dark:bg-zinc-100 dark:text-zinc-900"
        >
          ส่ง
        </button>
      </form>
    </section>
  );
}
