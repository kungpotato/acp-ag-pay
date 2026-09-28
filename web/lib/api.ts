export type Book = {
  id: string;
  title: string;
  author: string;
  genre: string;
  price_cents: number;
  currency: string;
  description: string;
  stock: number;
};

const API_BASE_URL =
  process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:3001";

export async function fetchBooks(query?: string): Promise<Book[]> {
  const url = new URL("/api/books", API_BASE_URL);
  if (query) url.searchParams.set("q", query);
  const res = await fetch(url, { cache: "no-store" });
  if (!res.ok) throw new Error(`failed to fetch books: ${res.status}`);
  return res.json();
}

export function formatPrice(book: Pick<Book, "price_cents" | "currency">) {
  return new Intl.NumberFormat("th-TH", {
    style: "currency",
    currency: book.currency.toUpperCase(),
  }).format(book.price_cents / 100);
}
