import { fetchBooks } from "@/lib/api";
import BookGrid from "@/components/BookGrid";
import CartPanel from "@/components/CartPanel";
import ChatPanel from "@/components/ChatPanel";

export default async function Home() {
  const books = await fetchBooks().catch(() => []);

  return (
    <main className="mx-auto max-w-3xl px-6 py-12">
      <h1 className="text-2xl font-semibold text-zinc-900 dark:text-zinc-50">
        ร้านขายหนังสือ
      </h1>
      <p className="mt-2 text-zinc-600 dark:text-zinc-400">
        บทที่ 1: แคตตาล็อกหนังสือจาก Go backend · บทที่ 4: ตะกร้าและ checkout session
      </p>

      <BookGrid books={books} />
      <CartPanel />
      <ChatPanel />
    </main>
  );
}
