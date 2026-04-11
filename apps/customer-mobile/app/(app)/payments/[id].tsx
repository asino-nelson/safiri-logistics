import { router, useLocalSearchParams } from "expo-router";
import { useEffect, useState } from "react";
import { Text, View } from "react-native";

import { ActionButton } from "../../../src/components/ActionButton";
import { AppField } from "../../../src/components/AppField";
import { BrandHeader } from "../../../src/components/BrandHeader";
import { Screen } from "../../../src/components/Screen";
import { useAuth } from "../../../src/context/auth";
import { getLoadById } from "../../../src/lib/loads";
import { initiateMpesaCheckout } from "../../../src/lib/payment";
import { Load, Payment } from "../../../src/lib/types";

export default function PaymentScreen() {
  const { id } = useLocalSearchParams<{ id: string }>();
  const { session } = useAuth();
  const [load, setLoad] = useState<Load | null>(null);
  const [phoneNumber, setPhoneNumber] = useState("+254700000000");
  const [payment, setPayment] = useState<Payment | null>(null);
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const run = async () => {
      if (!session || !id) {
        return;
      }

      setLoading(true);
      try {
        const currentLoad = await getLoadById(session.token, id);
        setLoad(currentLoad);
      } catch (err) {
        setError(err instanceof Error ? err.message : "Failed to fetch payment target");
      } finally {
        setLoading(false);
      }
    };

    void run();
  }, [session?.token, id]);

  const onSubmit = async () => {
    if (!session || !id) {
      return;
    }

    setSubmitting(true);
    setError(null);
    try {
      const response = await initiateMpesaCheckout(session.token, id, phoneNumber);
      setPayment(response);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to initiate payment");
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <Screen scroll>
      <ActionButton title="Back" variant="secondary" onPress={() => router.back()} />
      <BrandHeader
        eyebrow="Payments"
        title={load ? load.title : "M-Pesa checkout"}
        subtitle="Initiate payment against the quoted amount for this heavy-goods shipment."
      />

      {loading ? <Text style={{ color: "#64748B" }}>Loading payment details...</Text> : null}
      {error ? <Text style={{ color: "#B91C1C" }}>{error}</Text> : null}

      {load ? (
        <View style={{ gap: 14 }}>
          <View style={{ backgroundColor: "#FFFFFF", borderRadius: 18, borderWidth: 1, borderColor: "#E2E8F0", padding: 14 }}>
            <Text style={{ color: "#64748B", textTransform: "uppercase", fontSize: 12, letterSpacing: 0.8 }}>Quoted amount</Text>
            <Text style={{ color: "#0F172A", fontSize: 28, fontWeight: "800", marginTop: 6 }}>KES {load.quoted_price_kes.toLocaleString()}</Text>
          </View>

          <AppField label="Phone number" value={phoneNumber} onChangeText={setPhoneNumber} keyboardType="phone-pad" />
          <ActionButton title="Initiate M-Pesa checkout" onPress={onSubmit} loading={submitting} />

          {payment ? (
            <View style={{ backgroundColor: "#FFFFFF", borderRadius: 18, borderWidth: 1, borderColor: "#E2E8F0", padding: 14, gap: 8 }}>
              <Text style={{ color: "#0F172A", fontSize: 18, fontWeight: "800" }}>Payment initiated</Text>
              <Text style={{ color: "#64748B" }}>Provider: {payment.provider}</Text>
              <Text style={{ color: "#64748B" }}>Reference: {payment.provider_reference}</Text>
              <Text style={{ color: "#64748B" }}>Status: {payment.status}</Text>
            </View>
          ) : null}
        </View>
      ) : null}
    </Screen>
  );
}
