import { router } from "expo-router";
import { useEffect, useState } from "react";
import { Text, View } from "react-native";

import { ActionButton } from "../../src/components/ActionButton";
import { BrandHeader } from "../../src/components/BrandHeader";
import { LoadCard } from "../../src/components/LoadCard";
import { LoadingState, Screen, SectionTitle } from "../../src/components/Screen";
import { useAuth } from "../../src/context/auth";
import { listLoads } from "../../src/lib/loads";
import { Load } from "../../src/lib/types";

export default function DashboardScreen() {
  const { session } = useAuth();
  const [loads, setLoads] = useState<Load[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const loadDashboard = async () => {
    if (!session) {
      return;
    }

    setLoading(true);
    setError(null);
    try {
      const data = await listLoads(session.token);
      setLoads(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to fetch loads");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    void loadDashboard();
  }, [session?.token]);

  const counts = loads.reduce(
    (acc, load) => {
      acc[load.status] = (acc[load.status] || 0) + 1;
      return acc;
    },
    {} as Record<string, number>,
  );

  const recentLoads = [...loads].sort((a, b) => {
    const priorityDiff = b.priority - a.priority;
    if (priorityDiff !== 0) {
      return priorityDiff;
    }
    return new Date(b.created_at).getTime() - new Date(a.created_at).getTime();
  });

  return (
    <Screen scroll refreshable={{ refreshing: loading, onRefresh: loadDashboard }}>
      <BrandHeader
        eyebrow="Customer dashboard"
        title="Heavy-goods loads, matched fast."
        subtitle="Post trailers, see auto-matching progress, and keep every move visible."
      />

      <View style={{ flexDirection: "row", gap: 12, flexWrap: "wrap" }}>
        <MiniStat label="Posted" value={String(counts.posted || 0)} />
        <MiniStat label="Matched" value={String(counts.matched || 0)} />
        <MiniStat label="In transit" value={String(counts.in_transit || 0)} />
        <MiniStat label="Delivered" value={String(counts.delivered || 0)} />
      </View>

      <View style={{ gap: 12, marginTop: 18 }}>
        <ActionButton title="Create heavy-goods load" onPress={() => router.push("/create-load")} />
        <ActionButton title="Open profile" variant="secondary" onPress={() => router.push("/profile")} />
      </View>

      <SectionTitle title="Recent loads" subtitle="Sorted by priority and freshness so urgent freight is easy to spot." />

      {loading ? <LoadingState label="Loading loads" /> : null}
      {error ? <Text style={{ color: "#B91C1C" }}>{error}</Text> : null}

      <View style={{ gap: 14 }}>
        {recentLoads.map((load) => (
          <LoadCard key={load.id} load={load} onPress={() => router.push({ pathname: "/loads/[id]", params: { id: load.id } })} />
        ))}
        {!loading && !error && recentLoads.length === 0 ? (
          <Text style={{ color: "#64748B" }}>No loads yet. Create the first heavy-goods shipment.</Text>
        ) : null}
      </View>
    </Screen>
  );
}

function MiniStat({ label, value }: { label: string; value: string }) {
  return (
    <View style={{ minWidth: 140, flexGrow: 1, backgroundColor: "#FFFFFF", borderRadius: 18, padding: 14, borderWidth: 1, borderColor: "#E2E8F0" }}>
      <Text style={{ color: "#64748B", fontSize: 12, textTransform: "uppercase", letterSpacing: 0.8 }}>{label}</Text>
      <Text style={{ color: "#0F172A", fontSize: 28, fontWeight: "800", marginTop: 6 }}>{value}</Text>
    </View>
  );
}
