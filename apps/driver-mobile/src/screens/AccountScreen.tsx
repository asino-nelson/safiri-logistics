import { useRouter } from "expo-router";
import { useEffect } from "react";
import { StyleSheet, Text } from "react-native";

import { Button } from "@/components/Button";
import { Card } from "@/components/Card";
import { Screen } from "@/components/Screen";
import { SectionHeader } from "@/components/SectionHeader";
import { StatusPill } from "@/components/StatusPill";
import { colors } from "@/theme/colors";
import { useSession } from "@/context/session";

export function AccountScreen() {
  const router = useRouter();
  const { user, kycProfile, signOut, refreshProfile } = useSession();

  useEffect(() => {
    void refreshProfile();
  }, []);

  return (
    <Screen
      title="Account"
      subtitle="Driver identity, approval state, and session controls."
      footer={<Button label="Sign out" onPress={() => void signOut()} tone="danger" />}
    >
      <Card>
        <SectionHeader title="Identity" />
        <Text style={styles.meta}>Name: {user?.name ?? "-"}</Text>
        <Text style={styles.meta}>Email: {user?.email ?? "-"}</Text>
        <Text style={styles.meta}>Role: {user?.role ?? "-"}</Text>
      </Card>

      <Card>
        <SectionHeader title="KYC" />
        <StatusPill label={kycProfile?.status ?? "missing"} tone={kycTone(kycProfile?.status)} />
        <Text style={styles.meta}>Experience: {kycProfile?.years_experience ?? 0} years</Text>
        <Text style={styles.meta}>Capacity: {kycProfile?.max_load_kg ?? 0} kg</Text>
        <Text style={styles.meta}>Equipment: {kycProfile?.equipment_type ?? "flatbed"}</Text>
        <Button label="Review or submit KYC" onPress={() => router.push("/(auth)/kyc")} tone="secondary" />
      </Card>
    </Screen>
  );
}

function kycTone(status?: string): "success" | "warning" | "danger" | "default" {
  if (status === "approved") {
    return "success";
  }
  if (status === "rejected") {
    return "danger";
  }
  if (status === "pending") {
    return "warning";
  }
  return "default";
}

const styles = StyleSheet.create({
  meta: {
    color: colors.text,
    lineHeight: 20,
  },
});
