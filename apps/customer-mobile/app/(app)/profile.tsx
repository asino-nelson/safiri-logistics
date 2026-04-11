import { router } from "expo-router";
import { useEffect, useState } from "react";
import { Text, View } from "react-native";

import { ActionButton } from "../../src/components/ActionButton";
import { BrandHeader } from "../../src/components/BrandHeader";
import { Screen } from "../../src/components/Screen";
import { useAuth } from "../../src/context/auth";

export default function ProfileScreen() {
  const { session, signOut, refreshProfile } = useAuth();
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    void refreshProfile();
  }, [session?.token]);

  return (
    <Screen scroll>
      <BrandHeader
        eyebrow="Account"
        title={session?.user.name || "Profile"}
        subtitle="Manage your account session and verify the current signed-in user."
      />

      <View style={{ gap: 14 }}>
        <InfoCard label="Name" value={session?.user.name || "-"} />
        <InfoCard label="Email" value={session?.user.email || "-"} />
        <InfoCard label="Role" value={session?.user.role || "-"} />
        <InfoCard label="API base" value="http://localhost:8080/api/v1" />
      </View>

      <View style={{ gap: 12, marginTop: 18 }}>
        <ActionButton
          title="Sign out"
          variant="danger"
          loading={loading}
          onPress={async () => {
            setLoading(true);
            try {
              await signOut();
              router.replace("/login");
            } finally {
              setLoading(false);
            }
          }}
        />
      </View>
    </Screen>
  );
}

function InfoCard({ label, value }: { label: string; value: string }) {
  return (
    <View style={{ backgroundColor: "#FFFFFF", borderWidth: 1, borderColor: "#E2E8F0", borderRadius: 18, padding: 14 }}>
      <Text style={{ color: "#64748B", textTransform: "uppercase", letterSpacing: 0.8, fontSize: 12 }}>{label}</Text>
      <Text style={{ color: "#0F172A", fontSize: 16, fontWeight: "700", marginTop: 6 }}>{value}</Text>
    </View>
  );
}
