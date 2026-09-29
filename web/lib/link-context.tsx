"use client";

import { createContext, useContext, useState } from "react";
import type { Link } from "./link";

type LinkContextValue = {
  link: Link | null;
  setLink: (link: Link | null) => void;
};

const LinkContext = createContext<LinkContextValue | null>(null);

export function LinkProvider({ children }: { children: React.ReactNode }) {
  const [link, setLink] = useState<Link | null>(null);
  return (
    <LinkContext.Provider value={{ link, setLink }}>{children}</LinkContext.Provider>
  );
}

export function useLink() {
  const ctx = useContext(LinkContext);
  if (!ctx) throw new Error("useLink must be used within LinkProvider");
  return ctx;
}
