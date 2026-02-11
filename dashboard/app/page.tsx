cat > app/page.tsx <<'EOF'
"use client";

import { useEffect, useMemo, useState } from "react";

type Host = { host_id: string; hostname: string };
type ReportMeta = { id: number; timestamp?: string; ts_utc?: string; created_at: string };
type CheckResult = {
  check_id: string;
  title: string;
  status: "PASS" | "FAIL" | "ERROR";
  severity: "LOW" | "MEDIUM" | "HIGH";
  evidence: string;
};

export default function Page() {
  const BACKEND = process.env.NEXT_PUBLIC_BACKEND_URL;

  const [hosts, setHosts] = useState<Host[]>([]);
  const [selectedHost, setSelectedHost] = useState<Host | null>(null);

  const [reports, setReports] = useState<ReportMeta[]>([]);
  const [selectedReportId, setSelectedReportId] = useState<number | null>(null);

  const [checks, setChecks] = useState<CheckResult[]>([]);
  const [loading, setLoading] = useState(false);
  const [err, setErr] = useState<string>("");

  // Load hosts
  useEffect(() => {
    (async () => {
      try {
        setErr("");
        const r = await fetch(`${BACKEND}/hosts`);
        if (!r.ok) throw new Error(`hosts fetch failed: ${r.status}`);
        const data = await r.json();
        setHosts(data);
        if (data?.length) setSelectedHost(data[0]);
      } catch (e: any) {
        setErr(e?.message ?? "Failed to load hosts");
      }
    })();
  }, [BACKEND]);

  // Load reports list for host
  useEffect(() => {
    if (!selectedHost) return;
    (async () => {
      try {
        setErr("");
        const r = await fetch(`${BACKEND}/reports/${encodeURIComponent(selectedHost.host_id)}?limit=10`);
        if (!r.ok) throw new Error(`reports fetch failed: ${r.status}`);
        const data = await r.json();
        setReports(data);
        if (data?.length) setSelectedReportId(data[0].id); // latest
      } catch (e: any) {
        setErr(e?.message ?? "Failed to load reports");
      }
    })();
  }, [BACKEND, selectedHost]);

  // Load full report payload
  useEffect(() => {
    if (!selectedReportId) return;
    (async () => {
      try {
        setLoading(true);
        setErr("");
        const r = await fetch(`${BACKEND}/report/${selectedReportId}`);
        if (!r.ok) throw new Error(`report fetch failed: ${r.status}`);
        const data = await r.json();
        const payload = data.payload;
        setChecks(payload?.cis_results ?? []);
      } catch (e: any) {
        setErr(e?.message ?? "Failed to load report payload");
      } finally {
        setLoading(false);
      }
    })();
  }, [BACKEND, selectedReportId]);

  const stats = useMemo(() => {
    const pass = checks.filter(c => c.status === "PASS").length;
    const fail = checks.filter(c => c.status === "FAIL").length;
    const error = checks.filter(c => c.status === "ERROR").length;
    const total = checks.length || 1;
    const score = Math.round((pass / total) * 100);
    return { pass, fail, error, total: checks.length, score };
  }, [checks]);

  return (
    <main className="min-h-screen bg-gray-50 text-gray-900">
      <div className="mx-auto max-w-6xl p-6">
        <div className="flex items-end justify-between gap-4">
          <div>
            <h1 className="text-2xl font-semibold">CIS Compliance Dashboard</h1>
            <p className="text-sm text-gray-600">Live agent reports from your VM backend</p>
          </div>

          <div className="rounded-xl bg-white px-4 py-3 shadow-sm border">
            <div className="text-xs text-gray-500">Compliance Score</div>
            <div className="text-2xl font-semibold">{stats.score}%</div>
          </div>
        </div>

        {err ? (
          <div className="mt-4 rounded-xl border border-red-200 bg-red-50 p-4 text-sm text-red-700">
            {err}
          </div>
        ) : null}

        <div className="mt-6 grid grid-cols-1 gap-6 md:grid-cols-3">
          <div className="rounded-2xl bg-white shadow-sm border p-4">
            <div className="text-sm font-semibold">Hosts</div>
            <div className="mt-3 space-y-2">
              {hosts.map(h => (
                <button
                  key={h.host_id}
                  onClick={() => setSelectedHost(h)}
                  className={`w-full rounded-xl border px-3 py-2 text-left text-sm ${
                    selectedHost?.host_id === h.host_id ? "bg-gray-100" : "bg-white hover:bg-gray-50"
                  }`}
                >
                  <div className="font-medium">{h.hostname || "unknown-host"}</div>
                  <div className="text-xs text-gray-500">{h.host_id}</div>
                </button>
              ))}
              {!hosts.length ? <div className="text-sm text-gray-500">No hosts yet.</div> : null}
            </div>
          </div>

          <div className="rounded-2xl bg-white shadow-sm border p-4 md:col-span-2">
            <div className="flex flex-wrap items-center justify-between gap-3">
              <div>
                <div className="text-sm font-semibold">Latest Report</div>
                <div className="text-xs text-gray-500">
                  Host: {selectedHost?.hostname ?? "-"} ({selectedHost?.host_id ?? "-"})
                </div>
              </div>

              <select
                className="rounded-xl border bg-white px-3 py-2 text-sm"
                value={selectedReportId ?? ""}
                onChange={(e) => setSelectedReportId(Number(e.target.value))}
              >
                {reports.map(r => (
                  <option key={r.id} value={r.id}>
                    Report #{r.id} — {r.created_at}
                  </option>
                ))}
              </select>
            </div>

            <div className="mt-4 grid grid-cols-2 gap-3 md:grid-cols-4">
              <StatCard label="PASS" value={stats.pass} />
              <StatCard label="FAIL" value={stats.fail} />
              <StatCard label="ERROR" value={stats.error} />
              <StatCard label="TOTAL" value={stats.total} />
            </div>

            <div className="mt-5 overflow-auto rounded-xl border">
              <table className="w-full text-sm">
                <thead className="bg-gray-50">
                  <tr>
                    <th className="p-3 text-left">Check</th>
                    <th className="p-3 text-left">Status</th>
                    <th className="p-3 text-left">Severity</th>
                    <th className="p-3 text-left">Evidence</th>
                  </tr>
                </thead>
                <tbody>
                  {loading ? (
                    <tr><td className="p-3 text-gray-500" colSpan={4}>Loading report…</td></tr>
                  ) : checks.length ? (
                    checks.map((c) => (
                      <tr key={c.check_id} className="border-t">
                        <td className="p-3 font-medium">{c.check_id}</td>
                        <td className="p-3">{badge(c.status)}</td>
                        <td className="p-3">{c.severity}</td>
                        <td className="p-3 text-gray-600">{c.evidence}</td>
                      </tr>
                    ))
                  ) : (
                    <tr><td className="p-3 text-gray-500" colSpan={4}>No checks found.</td></tr>
                  )}
                </tbody>
              </table>
            </div>

            <p className="mt-3 text-xs text-gray-500">
              Tip: Run the agent again to generate a fresh report, then refresh this page.
            </p>
          </div>
        </div>
      </div>
    </main>
  );
}

function StatCard({ label, value }: { label: string; value: number }) {
  return (
    <div className="rounded-xl border bg-white p-3">
      <div className="text-xs text-gray-500">{label}</div>
      <div className="text-xl font-semibold">{value}</div>
    </div>
  );
}

function badge(s: string) {
  const base = "inline-flex rounded-full px-2 py-1 text-xs font-semibold";
  if (s === "PASS") return <span className={`${base} bg-green-100 text-green-800`}>PASS</span>;
  if (s === "FAIL") return <span className={`${base} bg-red-100 text-red-800`}>FAIL</span>;
  return <span className={`${base} bg-yellow-100 text-yellow-800`}>ERROR</span>;
}
EOF
