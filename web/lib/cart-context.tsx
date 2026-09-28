"use client";

import { createContext, useContext, useMemo, useState } from "react";
import type { Book } from "./api";

export type CartItem = {
  productId: string;
  title: string;
  unitPriceCents: number;
  currency: string;
  quantity: number;
};

type CartContextValue = {
  items: CartItem[];
  addBook: (book: Book) => void;
  removeItem: (productId: string) => void;
  clear: () => void;
  totalCents: number;
};

const CartContext = createContext<CartContextValue | null>(null);

export function CartProvider({ children }: { children: React.ReactNode }) {
  const [items, setItems] = useState<CartItem[]>([]);

  function addBook(book: Book) {
    setItems((prev) => {
      const existing = prev.find((i) => i.productId === book.id);
      if (existing) {
        return prev.map((i) =>
          i.productId === book.id ? { ...i, quantity: i.quantity + 1 } : i,
        );
      }
      return [
        ...prev,
        {
          productId: book.id,
          title: book.title,
          unitPriceCents: book.price_cents,
          currency: book.currency,
          quantity: 1,
        },
      ];
    });
  }

  function removeItem(productId: string) {
    setItems((prev) => prev.filter((i) => i.productId !== productId));
  }

  function clear() {
    setItems([]);
  }

  const totalCents = useMemo(
    () => items.reduce((sum, i) => sum + i.unitPriceCents * i.quantity, 0),
    [items],
  );

  return (
    <CartContext.Provider value={{ items, addBook, removeItem, clear, totalCents }}>
      {children}
    </CartContext.Provider>
  );
}

export function useCart() {
  const ctx = useContext(CartContext);
  if (!ctx) throw new Error("useCart must be used within CartProvider");
  return ctx;
}
