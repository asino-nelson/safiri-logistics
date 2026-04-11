import { useRouter } from "expo-router";
import { useEffect, useState } from "react";
import { RefreshControl, ScrollView, StyleSheet, Text, View } from "react-native";

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

export function LoadsScreen() {
  const router = useRouter();
  const { token, user, kycProfile, refreshProfile } = useSession();
  const [loads, setLoads] = useState<Load[]>([]);
  const [refreshing, setRefreshing] = useState(false);

  const loadData = async () => {
    if (!token) {
      return;
    }

    const response = await api.listLoads(token);
    setLoads(response.data);
  };

  useEffect(() => {
    void loadData();
  }, [token]);

  const refresh = async () => {
    setRefreshing(true);
    try {
      await Promise.all([loadData(), refreshProfile()]);
    } finally {
      setRefreshing(false);
    }
  };

  const assignedLoads = loads.filter((load) => load.assigned_driver_id === user?.id);
  const activeLoads = loads.filter((load) => load.status !== "delivered");

  return (
    <Screen
      title="Assigned loads"
      subtitle="Heavy goods first, ranked by proximity, availability, capacity, and experience."
      footer={<Button label="Refresh" onPress={refresh} tone="secondary" />}
    >
      <ScrollView
        style={styles.scroll}
        refreshControl={<RefreshControl refreshing={refreshing} onRefresh={refresh} tintColor={colors.accent} />}
      >
        <View style={styles.stack}>
          <Card>
            <SectionHeader title="Driver snapshot" />
            <Text style={styles.line}>Name: {user?.name ?? "-"}</Text>
            <Text style={styles.line}>KYC: {kycProfile?.status ?? "missing"}</Text>
            <Text style={styles.line}>Online: {kycProfile?.is_online ? "Yes" : "No"}</Text>
            <Text style={styles.line}>Available: {kycProfile?.is_available ? "Yes" : "No"}</Text>
          </Card>

          <Card>
            <SectionHeader title={`Assigned loads (${assignedLoads.length})`} />
            {assignedLoads.length === 0 ? (
              <EmptyState
                title="Nothing assigned yet"
                message="Once the engine matches you with a heavy load, it will show up here."
              />
            ) : (
              assignedLoads.map((load) => <LoadCard key={load.id} load={load} onPress={() => router.push(`/load/${load.id}`)} />)
            )}
          </Card>

          <Card>
            <SectionHeader title={`All active loads (${activeLoads.length})`} />
            {activeLoads.map((load) => (
              <LoadCard key={load.id} load={load} onPress={() => router.push(`/load/${load.id}`)} />
            ))}
          </Card>
        </View>
      </ScrollView>
    </Screen>
  );
}

function LoadCard({ load, onPress }: { load: Load; onPress: () => void }) {
  return (
    <View style={styles.loadCard}>
      <View style={styles.cardTop}>
        <View style={{ flex: 1 }}>
          <Text style={styles.title}>{load.title}</Text>
          <Text style={styles.meta}>
            {load.origin} -> {load.destination}
          </Text>
        </View>
        <StatusPill label={load.status} tone={statusTone(load.status)} />
      </View>
      <Text style={styles.meta}>Priority {load.priority} | {load.weight_kg} kg | {load.cargo_type}</Text>
      <Text style={styles.meta}>Equipment: {load.equipment_type}</Text>
      <Text style={styles.meta}>Pickup: {new Date(load.pickup_at).toLocaleString()}</Text>
      <Button label="Open details" onPress={onPress} tone="secondary" />
    </View>
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
  line: {
    color: colors.text,
    lineHeight: 20,
  },
  loadCard: {
    marginTop: 12,
    padding: 14,
    borderRadius: 18,
    backgroundColor: colors.surfaceElevated,
    borderWidth: 1,
    borderColor: colors.border,
    gap: 8,
  },
  cardTop: {
    flexDirection: "row",
    alignItems: "flex-start",
    gap: 10,
  },
  title: {
    color: colors.text,
    fontSize: 16,
    fontWeight: "800",
  },
  meta: {
    color: colors.muted,
    lineHeight: 19,
  },
});
