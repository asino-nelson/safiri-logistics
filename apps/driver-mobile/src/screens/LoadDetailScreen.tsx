import { useLocalSearchParams, useRouter } from "expo-router";
import { useEffect, useState } from "react";
import { Alert, ScrollView, StyleSheet, Text, View } from "react-native";

import { Button } from "@/components/Button";
import { Card } from "@/components/Card";
import { EmptyState } from "@/components/EmptyState";
import { Screen } from "@/components/Screen";
import { SectionHeader } from "@/components/SectionHeader";
import { StatusPill } from "@/components/StatusPill";
import { api } from "@/lib/api";
import { colors } from "@/theme/colors";
import { useSession } from "@/context/session";
import { Load } from "@/types/api";

export function LoadDetailScreen() {
  const router = useRouter();
  const params = useLocalSearchParams<{ loadId: string }>();
  const { token, user, refreshProfile } = useSession();
  const [load, setLoad] = useState<Load | null>(null);
  const [loading, setLoading] = useState(true);
  const [busy, setBusy] = useState(false);

  useEffect(() => {
    const fetchLoad = async () => {
      if (!token) {
        setLoading(false);
        return;
      }

      setLoading(true);
      try {
        const response = await api.listLoads(token);
        setLoad(response.data.find((item) => item.id === params.loadId) ?? null);
      } finally {
        setLoading(false);
      }
    };

    void fetchLoad();
  }, [params.loadId, token]);

  if (loading) {
    return (
      <Screen title="Load details" subtitle="Fetching the latest assignment record...">
        <EmptyState title="Loading load" message="Please wait while we sync the latest assignment details." />
      </Screen>
    );
  }

  if (!load) {
    return (
      <Screen title="Load details" subtitle="Fetching the latest assignment record...">
        <EmptyState title="Load not found" message="The current driver feed does not include this load yet." />
      </Screen>
    );
  }

  const canAct = load.assigned_driver_id === user?.id || load.status === "posted";

  const pick = async () => {
    if (!token) {
      return;
    }

    setBusy(true);
    try {
      const updated = await api.pickLoad(token, load.id);
      setLoad(updated);
      await refreshProfile();
      Alert.alert("Load picked", "You are now assigned to this load.");
    } catch (err) {
      Alert.alert("Pick failed", err instanceof Error ? err.message : "Unable to pick this load.");
    } finally {
      setBusy(false);
    }
  };

  const updateStatus = async (status: Load["status"]) => {
    if (!token) {
      return;
    }

    setBusy(true);
    try {
      const updated = await api.updateLoadStatus(token, load.id, status);
      setLoad(updated);
      Alert.alert("Status updated", `Load moved to ${status}.`);
    } catch (err) {
      Alert.alert("Update failed", err instanceof Error ? err.message : "Unable to change load status.");
    } finally {
      setBusy(false);
    }
  };

  return (
    <Screen
      title="Load details"
      subtitle="Pick the load, move it to transit, and complete delivery from one place."
      footer={<Button label="Back to loads" onPress={() => router.back()} tone="secondary" />}
    >
      <ScrollView style={styles.scroll}>
        <View style={styles.stack}>
          <Card>
            <View style={styles.topRow}>
              <View style={{ flex: 1 }}>
                <Text style={styles.title}>{load.title}</Text>
                <Text style={styles.meta}>{load.origin} -> {load.destination}</Text>
              </View>
              <StatusPill label={load.status} tone={statusTone(load.status)} />
            </View>
            <Text style={styles.meta}>{load.description}</Text>
            <Text style={styles.meta}>Poster: {load.poster_id}</Text>
            <Text style={styles.meta}>Assigned driver: {load.assigned_driver_id ?? "unassigned"}</Text>
            <Text style={styles.meta}>Quoted price: KES {load.quoted_price_kes.toFixed(2)}</Text>
            <Text style={styles.meta}>Matching score: {load.matching_score?.toFixed(2) ?? "n/a"}</Text>
            <Text style={styles.meta}>Assignment source: {load.assignment_source}</Text>
          </Card>

          <Card>
            <SectionHeader title="Cargo profile" />
            <Text style={styles.meta}>Weight: {load.weight_kg} kg</Text>
            <Text style={styles.meta}>Cargo type: {load.cargo_type}</Text>
            <Text style={styles.meta}>Equipment: {load.equipment_type}</Text>
            <Text style={styles.meta}>Priority: {load.priority}</Text>
          </Card>

          <Card>
            <SectionHeader title="Route" />
            <Text style={styles.meta}>Pickup at: {new Date(load.pickup_at).toLocaleString()}</Text>
            <Text style={styles.meta}>Pickup coordinates: {load.pickup_latitude.toFixed(4)}, {load.pickup_longitude.toFixed(4)}</Text>
            <Text style={styles.meta}>Dropoff coordinates: {load.dropoff_latitude.toFixed(4)}, {load.dropoff_longitude.toFixed(4)}</Text>
          </Card>

          {canAct ? (
            <Card>
              <SectionHeader title="Actions" />
              <Button label="Pick load" onPress={pick} disabled={busy || load.status === "picked"} />
              <Button
                label="Move to in transit"
                onPress={() => updateStatus("in_transit")}
                tone="secondary"
                disabled={busy || load.status === "delivered"}
              />
              <Button
                label="Mark delivered"
                onPress={() => updateStatus("delivered")}
                tone="primary"
                disabled={busy || load.status !== "in_transit"}
              />
              <Button
                label="Open tracking"
                onPress={() => router.push({ pathname: "/(tabs)/tracking", params: { loadId: load.id } })}
                tone="secondary"
              />
            </Card>
          ) : null}
        </View>
      </ScrollView>
    </Screen>
  );
}

function statusTone(status: Load["status"]) {
  switch (status) {
    case "delivered":
      return "success";
    case "in_transit":
      return "info";
    case "picked":
    case "matched":
      return "warning";
    default:
      return "default";
  }
}

const styles = StyleSheet.create({
  stack: {
    gap: 14,
  },
  scroll: {
    flex: 1,
  },
  topRow: {
    flexDirection: "row",
    alignItems: "flex-start",
    gap: 12,
  },
  title: {
    color: colors.text,
    fontSize: 20,
    fontWeight: "800",
    marginBottom: 4,
  },
  meta: {
    color: colors.muted,
    lineHeight: 20,
  },
});
