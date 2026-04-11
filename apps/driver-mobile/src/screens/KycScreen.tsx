import { useRouter } from "expo-router";
import { useEffect, useState } from "react";
import { Alert, ScrollView, StyleSheet, Text, View } from "react-native";

import { Button } from "@/components/Button";
import { Card } from "@/components/Card";
import { EmptyState } from "@/components/EmptyState";
import { Field } from "@/components/Field";
import { Screen } from "@/components/Screen";
import { SectionHeader } from "@/components/SectionHeader";
import { StatusPill } from "@/components/StatusPill";
import { api, ApiError } from "@/lib/api";
import { colors } from "@/theme/colors";
import { useSession } from "@/context/session";

export function KycScreen() {
  const router = useRouter();
  const { token, kycProfile, refreshProfile, user } = useSession();
  const [nationalId, setNationalId] = useState("");
  const [licenseNumber, setLicenseNumber] = useState("");
  const [truckRegistration, setTruckRegistration] = useState("");
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    if (kycProfile) {
      setNationalId(kycProfile.national_id ?? "");
      setLicenseNumber(kycProfile.license_number ?? "");
      setTruckRegistration(kycProfile.truck_registration ?? "");
    }
  }, [kycProfile]);

  const submit = async () => {
    if (!token) {
      return;
    }

    setSaving(true);
    try {
      await api.submitKyc(token, {
        national_id: nationalId,
        license_number: licenseNumber,
        truck_registration: truckRegistration,
      });
      await refreshProfile();
      Alert.alert("KYC submitted", "Your documents are now pending review.");
    } catch (err) {
      const message = err instanceof ApiError || err instanceof Error ? err.message : "Failed to submit KYC.";
      Alert.alert("KYC error", message);
    } finally {
      setSaving(false);
    }
  };

  return (
    <Screen
      title="KYC status"
      subtitle="Submit your driver credentials once and use this page to track approval state."
      footer={
        <Button
          label="Go to dashboard"
          onPress={() => router.replace("/(tabs)/loads")}
          tone="secondary"
        />
      }
    >
      <ScrollView style={styles.scroll}>
        <View style={styles.stack}>
          <Card>
            <SectionHeader title="Current review" />
            {kycProfile ? (
              <View style={styles.stack}>
                <StatusPill label={kycProfile.status} tone={kycTone(kycProfile.status)} />
                <Text style={styles.meta}>Driver: {user?.name ?? "Unknown"}</Text>
                <Text style={styles.meta}>National ID: {kycProfile.national_id}</Text>
                <Text style={styles.meta}>License: {kycProfile.license_number}</Text>
                <Text style={styles.meta}>Truck: {kycProfile.truck_registration}</Text>
                <Text style={styles.meta}>Equipment: {kycProfile.equipment_type}</Text>
                <Text style={styles.meta}>Experience: {kycProfile.years_experience} years</Text>
                <Text style={styles.meta}>Capacity: {kycProfile.max_load_kg} kg</Text>
                {kycProfile.rejection_reason ? (
                  <Text style={styles.rejection}>Reason: {kycProfile.rejection_reason}</Text>
                ) : null}
              </View>
            ) : (
              <EmptyState
                title="No KYC record yet"
                message="Complete the form below so the operations team can approve you for matching."
              />
            )}
          </Card>

          <Card>
            <SectionHeader title="Submit KYC" />
            <Field label="National ID" value={nationalId} onChangeText={setNationalId} />
            <Field label="License number" value={licenseNumber} onChangeText={setLicenseNumber} />
            <Field label="Truck registration" value={truckRegistration} onChangeText={setTruckRegistration} />
            <Button label={saving ? "Submitting..." : "Submit KYC"} onPress={submit} disabled={saving} />
          </Card>
        </View>
      </ScrollView>
    </Screen>
  );
}

function kycTone(status: string): "success" | "warning" | "danger" {
  if (status === "approved") {
    return "success";
  }
  if (status === "rejected") {
    return "danger";
  }
  return "warning";
}

const styles = StyleSheet.create({
  scroll: {
    flex: 1,
  },
  stack: {
    gap: 8,
  },
  meta: {
    color: colors.text,
    lineHeight: 20,
  },
  rejection: {
    color: colors.warning,
    lineHeight: 20,
  },
});
