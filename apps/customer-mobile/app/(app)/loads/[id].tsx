import { router, useLocalSearchParams } from "expo-router";
import { useEffect, useState } from "react";
import { Text, View } from "react-native";

import { ActionButton } from "../../../src/components/ActionButton";
import { BrandHeader } from "../../../src/components/BrandHeader";
import { LoadCard } from "../../../src/components/LoadCard";
import { Screen, SectionTitle } from "../../../src/components/Screen";
import { useAuth } from "../../../src/context/auth";
import { getLoadById } from "../../../src/lib/loads";
import { Load } from "../../../src/lib/types";

export default function LoadDetailScreen() {
  const { id } = useLocalSearchParams<{ id: string }>();
  const { session } = useAuth();
  const [load, setLoad] = useState<Load | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchLoad = async () => {
    if (!session || !id) {
      return;
    }

    setLoading(true);
    setError(null);
    try {
      const current = await getLoadById(session.token, id);
      setLoad(current);
    } catch (err) {
      setLoad(null);
      setError(err instanceof Error ? err.message : "Failed to fetch load");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    void fetchLoad();
  }, [session?.token, id]);

  return (
    <Screen scroll>
      <ActionButton title="Back to dashboard" variant="secondary" onPress={() => router.back()} />
      <BrandHeader
        eyebrow="Load detail"
        title={load ? load.title : "Load details"}
        subtitle="Inspect route, quote, assignment state, and delivery progress."
      />

      {loading ? <Text style={{ color: "#64748B" }}>Loading load detail...</Text> : null}
      {error ? <Text style={{ color: "#B91C1C" }}>{error}</Text> : null}

      {load ? (
        <View style={{ gap: 14 }}>
          <LoadCard load={load} compact />
          <DetailRow label="Route" value={`${load.origin} -> ${load.destination}`} />
          <DetailRow label="Cargo" value={`${load.cargo_type} | ${load.equipment_type}`} />
          <DetailRow label="Weight" value={`${load.weight_kg.toLocaleString()} kg`} />
          <DetailRow label="Priority" value={String(load.priority)} />
          <DetailRow label="Quoted price" value={`KES ${load.quoted_price_kes.toLocaleString()}`} />
          <DetailRow label="Assignment source" value={load.assignment_source} />
          <DetailRow label="Assigned driver" value={load.assigned_driver_id || "Not assigned yet"} />
          <DetailRow label="Pickup time" value={new Date(load.pickup_at).toLocaleString()} />

          <SectionTitle title="Actions" subtitle="Track the shipment or trigger payment once the quote is ready." />
          <View style={{ gap: 12 }}>
            <ActionButton title="Open tracking" onPress={() => router.push(`/tracking/${load.id}`)} />
            <ActionButton title="Initiate payment" variant="secondary" onPress={() => router.push(`/payments/${load.id}`)} />
            <ActionButton title="Refresh" variant="ghost" onPress={fetchLoad} />
          </View>
        </View>
      ) : null}
    </Screen>
  );
}

function DetailRow({ label, value }: { label: string; value: string }) {
  return (
    <View style={{ backgroundColor: "#FFFFFF", borderWidth: 1, borderColor: "#E2E8F0", borderRadius: 18, padding: 14 }}>
      <Text style={{ color: "#64748B", fontSize: 12, textTransform: "uppercase", letterSpacing: 0.8 }}>{label}</Text>
      <Text style={{ color: "#0F172A", fontSize: 16, fontWeight: "700", marginTop: 6 }}>{value}</Text>
    </View>
  );
}
