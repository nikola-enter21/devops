export default async function Page() {
  const apiBase = process.env.NEXT_PUBLIC_API_BASE ?? "http://localhost:8080";

  let ok = false;

  try {
    const res = await fetch(`${apiBase}/healthz`, {
      cache: "no-store",
    });
    if (!res.ok) throw new Error(`HTTP ${res.status}`);
    const data = (await res.json()) as { ok?: boolean };
    ok = Boolean(data?.ok);
  } catch (e) {
    console.error(e);
  }

  return (
    <main style={{ fontFamily: "system-ui, sans-serif", padding: 24 }}>
      <h1>Very nice frontend app</h1>
      <p>
        Backend ({apiBase}/healthz):{" "}
        <span style={{ fontWeight: 700, color: ok ? "green" : "crimson" }}>
          {ok ? "UP" : "DOWN"}
        </span>
      </p>
      {!ok && (
        <pre
          style={{
            background: "#111",
            color: "#eee",
            padding: 12,
            borderRadius: 8,
            whiteSpace: "pre-wrap",
          }}
        >
          Error: Could not reach backend
        </pre>
      )}
    </main>
  );
}
