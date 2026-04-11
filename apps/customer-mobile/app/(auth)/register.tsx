import { Link, router } from "expo-router";
import { useState } from "react";
import { Text, View } from "react-native";

import { ActionButton } from "../../src/components/ActionButton";
import { AppField } from "../../src/components/AppField";
import { BrandHeader } from "../../src/components/BrandHeader";
import { Screen } from "../../src/components/Screen";
import { useAuth } from "../../src/context/auth";

export default function RegisterScreen() {
  const { signUp } = useAuth();
  const [name, setName] = useState("Customer One");
  const [email, setEmail] = useState("customer@example.com");
  const [password, setPassword] = useState("Password123");
  const [error, setError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);

  const onSubmit = async () => {
    setSubmitting(true);
    setError(null);
    try {
      await signUp({ name, email, password });
      router.replace("/dashboard");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Unable to create account");
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <Screen scroll>
      <BrandHeader
        eyebrow="New customer"
        title="Register for freight operations."
        subtitle="Create a customer account and start posting loads immediately."
      />

      <View style={{ gap: 14 }}>
        <AppField label="Full name" value={name} onChangeText={setName} />
        <AppField label="Email" value={email} onChangeText={setEmail} keyboardType="email-address" autoCapitalize="none" />
        <AppField label="Password" value={password} onChangeText={setPassword} secureTextEntry />
        {error ? <Text style={{ color: "#B91C1C" }}>{error}</Text> : null}
        <ActionButton title="Create account" onPress={onSubmit} loading={submitting} />
        <Text style={{ color: "#64748B" }}>
          Already registered? <Link href="/login" style={{ color: "#0F766E", fontWeight: "700" }}>Sign in</Link>
        </Text>
      </View>
    </Screen>
  );
}
