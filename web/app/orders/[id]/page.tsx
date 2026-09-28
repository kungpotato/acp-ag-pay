import OrderStatus from "@/components/OrderStatus";

export default async function OrderConfirmationPage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const { id } = await params;

  return (
    <main className="mx-auto max-w-lg px-6 py-12">
      <h1 className="text-2xl font-semibold text-zinc-900 dark:text-zinc-50">
        สถานะคำสั่งซื้อ
      </h1>
      <p className="mt-2 text-zinc-600 dark:text-zinc-400">
        บทที่ 6: หน้ายืนยันคำสั่งซื้อ — อัปเดตอัตโนมัติทุก 3 วินาที
      </p>
      <div className="mt-6">
        <OrderStatus orderId={id} />
      </div>
    </main>
  );
}
