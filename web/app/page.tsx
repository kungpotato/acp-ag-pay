import { fetchBooks, formatPrice } from "@/lib/api";
import ChatPanel from "@/components/ChatPanel";

export default async function Home() {
  const books = await fetchBooks().catch(() => []);

  return (
    <main className="mx-auto max-w-3xl px-6 py-12">
      <h1 className="text-2xl font-semibold text-zinc-900 dark:text-zinc-50">
        ร้านขายหนังสือ
      </h1>
      <p className="mt-2 text-zinc-600 dark:text-zinc-400">
        บทที่ 1: แคตตาล็อกหนังสือจาก Go backend
      </p>

      {books.length === 0 ? (
        <p className="mt-8 text-sm text-zinc-500">
          ไม่พบข้อมูลหนังสือ — ตรวจสอบว่า Go server รันอยู่ที่ port 3001
        </p>
      ) : (
        <ul className="mt-8 grid grid-cols-1 gap-4 sm:grid-cols-2">
          {books.map((book) => (
            <li
              key={book.id}
              className="rounded-lg border border-zinc-200 p-4 dark:border-zinc-800"
            >
              <p className="text-xs uppercase tracking-wide text-zinc-500">
                {book.genre}
              </p>
              <h2 className="mt-1 font-medium text-zinc-900 dark:text-zinc-50">
                {book.title}
              </h2>
              <p className="text-sm text-zinc-600 dark:text-zinc-400">
                {book.author}
              </p>
              <p className="mt-2 text-sm text-zinc-500">{book.description}</p>
              <div className="mt-3 flex items-center justify-between">
                <span className="font-medium">{formatPrice(book)}</span>
                <span className="text-xs text-zinc-500">
                  คงเหลือ {book.stock} เล่ม
                </span>
              </div>
            </li>
          ))}
        </ul>
      )}

      <ChatPanel />
    </main>
  );
}
