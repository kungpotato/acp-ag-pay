"use client";

import type { Book } from "@/lib/api";
import { formatPrice } from "@/lib/api";
import { useCart } from "@/lib/cart-context";

export default function BookGrid({ books }: { books: Book[] }) {
  const { addBook } = useCart();

  if (books.length === 0) {
    return (
      <p className="mt-8 text-sm text-zinc-500">
        ไม่พบข้อมูลหนังสือ — ตรวจสอบว่า Go server รันอยู่ที่ port 3001
      </p>
    );
  }

  return (
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
          <button
            onClick={() => addBook(book)}
            disabled={book.stock <= 0}
            className="mt-3 w-full rounded border border-zinc-300 px-3 py-1.5 text-sm disabled:opacity-40 dark:border-zinc-700"
          >
            เพิ่มลงตะกร้า
          </button>
        </li>
      ))}
    </ul>
  );
}
