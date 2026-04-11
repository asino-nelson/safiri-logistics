import { router } from "expo-router";
import { useState } from "react";
import { Text, View } from "react-native";

import { ActionButton } from "../../src/components/ActionButton";
import { AppField } from "../../src/components/AppField";
import { BrandHeader } from "../../src/components/BrandHeader";
import { Screen } from "../../src/components/Screen";
import { useAuth } from "../../src/context/auth";
import { createLoad } from "../../src/lib/loads";

const tomorrow = new Date(Date.now() + 24 * 60 * 60 * 1000).toISOString();

export default function CreateLoadScreen() {
  const { session } = useAuth();
  const [form, setForm] = useState({
    title: "Heavy machinery to Nairobi",
    description: "Transport a generator and earthmoving equipment on a flatbed trailer.",
    origin: "Eldoret",
    destination: "Nairobi",
    weight_kg: "32000",
    pickup_latitude: "-0.5143",
    pickup_longitude: "35.2698",
    dropoff_latitude: "-1.2864",
    dropoff_longitude: "36.8172",
    pickup_at: tomorrow,
    priority: "5",
    cargo_type: "machinery",
    equipment_type: "flatbed",
  });
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const update = (key: keyof typeof form, value: string) => setForm((current) => ({ ...current, [key]: value }));

  const onSubmit = async () => {
    if (!session) {
      return;
    }

    setSubmitting(true);
    setError(null);
    try {
      const load = await createLoad(session.token, {
        title: form.title,
        description: form.description,
        origin: form.origin,
        destination: form.destination,
        weight_kg: Number(form.weight_kg),
        pickup_latitude: Number(form.pickup_latitude),
        pickup_longitude: Number(form.pickup_longitude),
        dropoff_latitude: Number(form.dropoff_latitude),
        dropoff_longitude: Number(form.dropoff_longitude),
        pickup_at: form.pickup_at ? new Date(form.pickup_at).toISOString() : undefined,
        priority: Number(form.priority),
        cargo_type: form.cargo_type,
        equipment_type: form.equipment_type,
      });
      router.replace({ pathname: "/loads/[id]", params: { id: load.id } });
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to create load");
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <Screen scroll>
      <BrandHeader
        eyebrow="New freight"
        title="Post a heavy-goods load."
        subtitle="Capture the route, equipment, and pickup timing so the matching engine can work faster."
      />

      <View style={{ gap: 14 }}>
        <AppField label="Title" value={form.title} onChangeText={(value) => update("title", value)} />
        <AppField label="Description" value={form.description} onChangeText={(value) => update("description", value)} multiline />
        <AppField label="Origin" value={form.origin} onChangeText={(value) => update("origin", value)} />
        <AppField label="Destination" value={form.destination} onChangeText={(value) => update("destination", value)} />
        <AppField label="Weight (kg)" value={form.weight_kg} onChangeText={(value) => update("weight_kg", value)} keyboardType="numeric" />
        <AppField label="Pickup latitude" value={form.pickup_latitude} onChangeText={(value) => update("pickup_latitude", value)} keyboardType="numeric" />
        <AppField label="Pickup longitude" value={form.pickup_longitude} onChangeText={(value) => update("pickup_longitude", value)} keyboardType="numeric" />
        <AppField label="Dropoff latitude" value={form.dropoff_latitude} onChangeText={(value) => update("dropoff_latitude", value)} keyboardType="numeric" />
        <AppField label="Dropoff longitude" value={form.dropoff_longitude} onChangeText={(value) => update("dropoff_longitude", value)} keyboardType="numeric" />
        <AppField label="Pickup time ISO" value={form.pickup_at} onChangeText={(value) => update("pickup_at", value)} />
        <AppField label="Priority" value={form.priority} onChangeText={(value) => update("priority", value)} keyboardType="numeric" />
        <AppField label="Cargo type" value={form.cargo_type} onChangeText={(value) => update("cargo_type", value)} />
        <AppField label="Equipment type" value={form.equipment_type} onChangeText={(value) => update("equipment_type", value)} />
        {error ? <Text style={{ color: "#B91C1C" }}>{error}</Text> : null}
        <ActionButton title="Create load" onPress={onSubmit} loading={submitting} />
      </View>
    </Screen>
  );
}
