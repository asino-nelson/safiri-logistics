import { useRouter } from "expo-router";
import { useState } from "react";
import { ActivityIndicator, StyleSheet, Text } from "react-native";

import { Button } from "@/components/Button";
import { Card } from "@/components/Card";
import { Field } from "@/components/Field";
import { Screen } from "@/components/Screen";
import { ApiError } from "@/lib/api";
import { colors } from "@/theme/colors";
import { useSession } from "@/context/session";

export function LoginScreen() {
  const router = useRouter();
  const { signIn } = useSession();
  const [email, setEmail] = useState("driver@example.com");
  const [password, setPassword] = useState("Password123");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleLogin = async () => {
    setLoading(true);
    setError(null);

    try {
      await signIn({ email, password });
      router.replace("/");
    } catch (err) {
      if (err instanceof ApiError || err instanceof Error) {
        setError(err.message);
      } else {
        setError("Unable to sign in right now.");
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <Screen
      title="Driver login"
      subtitle="Sign in, confirm your KYC, then move loads from pickup to delivery with live tracking."
      footer={
        <Card>
          <Text style={styles.footerText}>Backend</Text>
          <Text style={styles.footerValue}>{process.env.EXPO_PUBLIC_API_BASE_URL ?? "http://localhost:8080"}</Text>
        </Card>
      }
    >
      <Card>
        <Field label="Email" value={email} autoCapitalize="none" keyboardType="email-address" onChangeText={setEmail} />
        <Field label="Password" value={password} secureTextEntry onChangeText={setPassword} />
        {error ? <Text style={styles.error}>{error}</Text> : null}
        <Button label={loading ? "Signing in..." : "Sign in"} onPress={handleLogin} disabled={loading} />
        {loading ? <ActivityIndicator color={colors.accent} /> : null}
      </Card>
    </Screen>
  );
}

const styles = StyleSheet.create({
  error: {
    color: colors.danger,
    fontSize: 13,
    lineHeight: 18,
  },
  footerText: {
    color: colors.muted,
    fontSize: 12,
    textTransform: "uppercase",
    letterSpacing: 0.8,
  },
  footerValue: {
    color: colors.text,
    fontSize: 14,
    fontWeight: "700",
    marginTop: 4,
  },
});
