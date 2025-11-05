"use client";

import { useEffect, useState } from "react";

type Status = "loading" | "online" | "offline";

export default function Home() {
  const [apiStatus, setApiStatus] = useState<Status>("loading");
  const [dbStatus, setDbStatus] = useState<Status>("loading");

  useEffect(() => {
    const check = async () => {
      try {
        const [apiRes, dbRes] = await Promise.all([
          fetch(`${process.env.NEXT_PUBLIC_BACKEND_URL}/api/v1/healthz`, {
            cache: "no-store",
          }),
          fetch(`${process.env.NEXT_PUBLIC_BACKEND_URL}/api/v1/checkDatabase`, {
            cache: "no-store",
          }),
        ]);

        setApiStatus(apiRes.ok ? "online" : "offline");
        setDbStatus(dbRes.ok ? "online" : "offline");
      } catch {
        setApiStatus("offline");
        setDbStatus("offline");
      }
    };

    check();
  }, []);

  const color = (status: Status) =>
    status === "loading"
      ? "bg-yellow-600 animate-pulse"
      : status === "online"
      ? "bg-green-600"
      : "bg-red-600";

  const text = (status: Status) =>
    status === "loading"
      ? "Checking..."
      : status === "online"
      ? "Online"
      : "Offline";

  return (
    <div className="flex min-h-screen items-center justify-center bg-zinc-950 px-4 sm:px-6 md:px-8 py-16 sm:py-24 text-white">
      <div className="w-full max-w-sm sm:max-w-md md:max-w-lg rounded-xl bg-zinc-900 p-6 sm:p-8 shadow-lg ring-1 ring-zinc-800 opacity-0 translate-y-3 animate-fade-in transition-all duration-700 ease-out">
        <h1 className="text-2xl sm:text-3xl font-bold mb-6 sm:mb-8 text-center tracking-tight">
          System Health
        </h1>

        <div className="space-y-3 sm:space-y-4 text-base sm:text-lg">
          <div className="flex items-center justify-between">
            <span>API</span>
            <span
              className={`inline-block rounded-full px-3 py-1 text-sm font-medium transition-colors duration-500 ${color(
                apiStatus
              )}`}
            >
              {text(apiStatus)}
            </span>
          </div>

          <div className="flex items-center justify-between">
            <span>Database</span>
            <span
              className={`inline-block rounded-full px-3 py-1 text-sm font-medium transition-colors duration-500 ${color(
                dbStatus
              )}`}
            >
              {text(dbStatus)}
            </span>
          </div>
        </div>

        <div className="mt-8 text-center text-xs sm:text-sm text-zinc-400 space-y-1 break-all">
          <p>
            API:{" "}
            <code className="text-zinc-500">
              {process.env.NEXT_PUBLIC_BACKEND_URL}/api/v1/healthz
            </code>
          </p>
          <p>
            DB:{" "}
            <code className="text-zinc-500">
              {process.env.NEXT_PUBLIC_BACKEND_URL}/api/v1/checkDatabase
            </code>
          </p>
        </div>
      </div>
    </div>
  );
}
