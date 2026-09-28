# ภาพรวม Workshop: สร้างร้านขายหนังสือด้วย Agentic Commerce Protocol (ACP) + Stripe

## เป้าหมายของ workshop นี้

ผู้เรียนจะสร้างระบบ "ร้านขายหนังสือที่คุยกับ AI agent ได้" ตั้งแต่ศูนย์ จนถึงจุดที่ agent
สามารถค้นหาหนังสือ สร้างตะกร้าสินค้า และ **settle การชำระเงินจริงผ่าน Stripe** ได้ตาม
มาตรฐาน Agentic Commerce Protocol (ACP) ทุกบทเรียนออกแบบให้ตอบคำถามได้ว่า
**"เรียนบทนี้ไปทำไม"** ก่อนเริ่มลงมือเขียนโค้ด

## สถาปัตยกรรมโดยรวม

- **Backend**: Go (net/http มาตรฐาน ไม่ใช้ framework ใหญ่ เพื่อให้เห็น ACP ชัดเจน)
- **Frontend**: Next.js (หน้าจอแชทคุยกับ shopping agent)
- **Payment settlement**: Stripe (test mode) ผ่านขั้นตอน ACP checkout → settle
- **เอกสารประกอบ**: ภาษาไทยทั้งหมด อยู่ใน `docs/`

ดูภาพสถาปัตยกรรมระดับสูง: [00-architecture.png](diagrams/00-architecture.png)

## โครงสร้างบทเรียน

| บท | ชื่อบท | เรียนไปทำไม |
|----|--------|--------------|
| 1 | [แคตตาล็อกหนังสือ (Go API + Next.js)](01-catalog.md) | ปูพื้นฐาน REST API และการต่อ frontend-backend ก่อนใส่ความฉลาดของ agent |
| 2 | [Shopping agent: แปลความต้องการผู้ใช้](02-shopping-agent.md) | เข้าใจว่า agent "ฟังภาษาคน" แล้วแปลงเป็น action ได้อย่างไร ทั้งแบบ rule-based (ฟรี) และแบบ LLM (OpenRouter) |
| 3 | [ACP product feed & discovery](03-acp-product-feed.md) | เข้าใจว่ามาตรฐาน ACP ให้ agent ภายนอก "ค้นพบ" สินค้าของร้านได้อย่างไร โดยไม่ต้อง hardcode |
| 4 | [ตะกร้าสินค้าและ ACP checkout session](04-checkout-session.md) | เข้าใจ lifecycle ของ checkout session ตาม ACP ก่อนจะไปถึงขั้นจ่ายเงินจริง |
| 5 | [Settle การชำระเงินผ่าน Stripe](05-stripe-settlement.md) | หัวใจของ workshop: agent สั่ง "จ่ายเงิน" แล้วระบบ settle ผ่าน Stripe PaymentIntent จริง |
| 6 | [สถานะคำสั่งซื้อและหน้ายืนยัน](06-order-status.md) | ปิด loop ฝั่งผู้ใช้ ให้เห็นสถานะคำสั่งซื้อเปลี่ยนแบบ real-time |
| 7 | [Stripe webhook: ยืนยันการชำระเงินแบบ async](07-stripe-webhook.md) | เข้าใจว่าทำไมจะเชื่อ response ตอน checkout อย่างเดียวไม่พอ ต้องรอ webhook ยืนยันจริง |
| 8 | [สรุปรวม: เดโมแบบ end-to-end](08-summary.md) | ต่อทุกชิ้นส่วนเข้าด้วยกัน ทบทวนสิ่งที่ได้เรียนรู้ทั้งหมด |

## หลักการทำงานระหว่าง workshop

- **Commit เล็กและถี่** ตามแนวทางหนังสือ *Accelerate*: เปลี่ยนแปลงทีละก้อนเล็ก ๆ ที่ build ผ่านและทดสอบได้
  แทนที่จะรวมงานทั้งบทไว้ commit เดียว เพื่อให้ trunk พร้อม deploy ได้ตลอดเวลา
- **แต่ละบท** จะมีโฟลเดอร์เอกสารของตัวเองใน `docs/` พร้อม diagram (png) และ GIF สาธิตการใช้งาน
- **ทดสอบด้วย Stripe test mode เท่านั้น** ห้ามใช้ live key ใน workshop นี้

## สิ่งที่ต้องเตรียมก่อนเริ่ม

1. Go >= 1.22
2. Node.js >= 18 และ pnpm/npm
3. บัญชี Stripe (test mode) และ Stripe CLI สำหรับ `stripe listen`
4. (ไม่บังคับ) OpenRouter API key ถ้าอยากลองใช้ LLM แทน rule-based parser ในบทที่ 2

## แผนที่เอกสาร

- `docs/00-overview.md` (ไฟล์นี้)
- `docs/01-catalog.md` ถึง `docs/08-summary.md` — เอกสารแต่ละบท
- `docs/diagrams/` — ไฟล์ png ของทุก diagram
- `docs/demos/` — GIF สาธิตการทำงานของแต่ละบท
