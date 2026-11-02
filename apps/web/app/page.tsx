import { Suspense } from "react";
import { connection } from "next/server";
import { fetchApiStatus } from "@/lib/api-status";

async function ApiStatusText() {
  await connection();
  const status = await fetchApiStatus(process.env.API_URL);
  return <p className="text-muted">API: {status}</p>;
}

export default function Home() {
  return (
    <main className="mx-auto flex min-h-screen max-w-3xl flex-col justify-center gap-4 px-6">
      <h1 className="text-3xl font-semibold">Beatroom</h1>
      <Suspense fallback={<p className="text-muted">API: checking...</p>}>
        <ApiStatusText />
      </Suspense>
    </main>
  );
}
