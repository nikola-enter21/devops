"use client";

import { useEffect, useState } from "react";

export default function Home() {
  const [status, setStatus] = useState<"loading" | "online" | "offline">(
    "loading"
  );

  useEffect(() => {
    const checkHealth = async () => {
      try {
        const res = await fetch(
          `${process.env.NEXT_PUBLIC_BACKEND_URL}/healthz`,
          {
            cache: "no-store",
          }
        );
        if (res.ok) {
          setStatus("online");
        } else {
          setStatus("offline");
        }
      } catch {
        setStatus("offline");
      }
    };

    checkHealth();
  }, []);

  return (
    <div className="flex min-h-screen items-center justify-center bg-zinc-950 px-6 py-24 text-white">
      <div className="w-full max-w-md rounded-xl bg-zinc-900 p-8 shadow-lg ring-1 ring-zinc-800">
        <h1 className="text-3xl font-bold mb-6 text-center">Backend Status</h1>

        <div className="flex items-center justify-between text-lg">
          <span>API</span>
          <span
            className={`inline-block rounded-full px-3 py-1 text-sm font-medium ${
              status === "loading"
                ? "bg-yellow-600"
                : status === "online"
                ? "bg-green-600"
                : "bg-red-600"
            }`}
          >
            {status === "loading"
              ? "Checking..."
              : status === "online"
              ? "Online"
              : "Offline"}
          </span>
        </div>

        <div className="mt-6 text-center text-sm text-zinc-400">
          <p>
            Pinged: <code>{process.env.NEXT_PUBLIC_BACKEND_URL}/healthz</code>
          </p>
        </div>
      </div>
    </div>
  );
}
