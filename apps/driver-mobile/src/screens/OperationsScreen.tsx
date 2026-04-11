import { useEffect, useState } from "react";
import { Alert, ScrollView, StyleSheet, Text, View } from "react-native";

import { Button } from "@/components/Button";
import { Card } from "@/components/Card";
import { Field } from "@/components/Field";
import { Screen } from "@/components/Screen";
import { SectionHeader } from "@/components/SectionHeader";
import { StatusPill } from "@/components/StatusPill";
import { api, ApiError } from "@/lib/api";
import { colors } from "@/theme/colors";
import { useSession } from "@/context/session";

export function OperationsScreen() {
  const { token, kycProfile, refreshProfile } = useSession();
  const [yearsExperience, setYearsExperience] = useState("0");
  const [maxLoadKg, setMaxLoadKg] = useState("0");
  const [equipmentType, setEquipmentType] = useState("flatbed");
  const [isOnline, setIsOnline] = useState(true);
  const [isAvailable, setIsAvailable] = useState(true);
  const [latitude, setLatitude] = useState("");
  const [longitude, setLongitude] = useState("");
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    if (!kycProfile) {
      return;
    }

    setYearsExperience(String(kycProfile.years_experience ?? 0));
    setMaxLoadKg(String(kycProfile.max_load_kg ?? 0));
    setEquipmentType(kycProfile.equipment_type ?? "flatbed");
    setIsOnline(Boolean(kycProfile.is_online));
    setIsAvailable(Boolean(kycProfile.is_available));
    setLatitude(kycProfile.current_latitude?.toString() ?? "");
    setLongitude(kycProfile.current_longitude?.toString() ?? "");
  }, [kycProfile]);

  const save = async () => {
    if (!token) {
      return;
    }

    setSaving(true);
    try {
      await api.updateOperations(token, {
        years_experience: Number(yearsExperience) || 0,
        max_load_kg: Number(maxLoadKg) || 0,
        equipment_type: equipmentType,
        is_online: isOnline,
        is_available: isAvailable,
        latitude: latitude.trim() ? Number(latitude) : null,
        longitude: longitude.trim() ? Number(longitude) : null,
      });
      await refreshProfile();
      Alert.alert("Operations updated", "Driver readiness is now synced.");
    } catch (err) {
      const message = err instanceof ApiError || err instanceof Error ? err.message : "Failed to update operations.";
      Alert.alert("Operations error", message);
    } finally {
      setSaving(false);
    }
  };

  return (
    <Screen
      title="Driver operations"
      subtitle="Stay online, advertise your capacity, and keep your live coordinates current for matching."
    >
      <ScrollView style={styles.scroll}>
        <View style={styles.stack}>
          <Card>
            <SectionHeader title="Operational status" />
            <View style={styles.row}>
              <Button
                label={isOnline ? "Online" : "Offline"}
                onPress={() => setIsOnline(!isOnline)}
                tone={isOnline ? "primary" : "secondary"}
                style={styles.flex}
              />
              <Button
                label={isAvailable ? "Available" : "Busy"}
                onPress={() => setIsAvailable(!isAvailable)}
                tone={isAvailable ? "primary" : "secondary"}
                style={styles.flex}
              />
            </View>
            <StatusPill label={kycProfile?.status ?? "kyc pending"} tone={kycProfile?.status === "approved" ? "success" : "warning"} />
          </Card>

          <Card>
            <SectionHeader title="Capacity and experience" />
            <Field label="Years of experience" keyboardType="numeric" value={yearsExperience} onChangeText={setYearsExperience} />
            <Field label="Max load capacity (kg)" keyboardType="numeric" value={maxLoadKg} onChangeText={setMaxLoadKg} />
            <Field label="Equipment type" value={equipmentType} onChangeText={setEquipmentType} helperText="Examples: flatbed, low-bed, trailer, container" />
          </Card>

          <Card>
            <SectionHeader title="Current coordinates" />
            <Field label="Latitude" keyboardType="numeric" value={latitude} onChangeText={setLatitude} />
            <Field label="Longitude" keyboardType="numeric" value={longitude} onChangeText={setLongitude} />
            <Text style={styles.note}>
              Keep these coordinates fresh so the matching engine can rank you by proximity.
            </Text>
            <Button label={saving ? "Saving..." : "Save operations"} onPress={save} disabled={saving} />
          </Card>
        </View>
      </ScrollView>
    </Screen>
  );
}

const styles = StyleSheet.create({
  scroll: {
    flex: 1,
  },
  stack: {
    gap: 14,
  },
  row: {
    flexDirection: "row",
    gap: 12,
  },
  flex: {
    flex: 1,
  },
  note: {
    color: colors.muted,
    lineHeight: 20,
  },
});
