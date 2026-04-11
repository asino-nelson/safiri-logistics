import { Link, router } from "expo-router";
import { useState } from "react";
import { Text, View } from "react-native";

import { ActionButton } from "../../src/components/ActionButton";
import { AppField } from "../../src/components/AppField";
import { BrandHeader } from "../../src/components/BrandHeader";
import { Screen } from "../../src/components/Screen";
import { useAuth } from "../../src/context/auth";

export default function LoginScreen() {
  const { signIn } = useAuth();
  const [email, setEmail] = useState("customer@example.com");
  const [password, setPassword] = useState("Password123");
  const [error, setError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);

  const onSubmit = async () => {
    setSubmitting(true);
    setError(null);
    try {
      await signIn(email, password);
      router.replace("/dashboard");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Unable to sign in");
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <Screen scroll>
      <BrandHeader
        eyebrow="Customer access"
        title="Move heavy goods with clarity."
        subtitle="Log in to post loads, track shipments, and trigger M-Pesa payment initiation."
      />

      <View style={{ gap: 14 }}>
        <AppField label="Email" value={email} onChangeText={setEmail} keyboardType="email-address" autoCapitalize="none" />
        <AppField label="Password" value={password} onChangeText={setPassword} secureTextEntry />
        {error ? <Text style={{ color: "#B91C1C" }}>{error}</Text> : null}
        <ActionButton title="Sign in" onPress={onSubmit} loading={submitting} />
        <Text style={{ color: "#64748B" }}>
          No account yet? <Link href="/register" style={{ color: "#0F766E", fontWeight: "700" }}>Create one</Link>
        </Text>
      </View>
    </Screen>
  );
}
