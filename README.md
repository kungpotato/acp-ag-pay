# ร้านขายหนังสือ ACP + Stripe Workshop

Workshop สร้างร้านขายหนังสือที่คุยกับ AI shopping agent ได้ และ settle การชำระเงินจริง
ผ่าน Stripe ตามมาตรฐาน Agentic Commerce Protocol (ACP)

เอกสารฉบับเต็มอยู่ที่ [docs/00-overview.md](docs/00-overview.md)

## โครงสร้างโปรเจกต์

- `server/` — Go backend (catalog, shopping agent, ACP checkout, Stripe settlement)
- `web/` — Next.js frontend (หน้าจอแชทซื้อหนังสือ)
- `docs/` — เอกสารประกอบแต่ละบทเรียน (ภาษาไทย), diagram และ demo GIF

## เริ่มต้นอย่างเร็ว

```bash
cp .env.example .env          # ตั้งค่า Stripe key เมื่อถึงบทที่ 5+
cd server && go run ./cmd/server &
cd web && npm install && npm run dev
```

เปิด http://localhost:3000
