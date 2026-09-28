import type { Book } from "./api";

export type ChatReply = {
  message: string;
  books?: Book[];
  intent: { kind: string; query?: string };
};

const API_BASE_URL =
  process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:3001";

export async function sendChatMessage(message: string): Promise<ChatReply> {
  const res = await fetch(new URL("/api/agent/chat", API_BASE_URL), {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ message }),
  });
  if (!res.ok) throw new Error(`agent chat failed: ${res.status}`);
  return res.json();
}
