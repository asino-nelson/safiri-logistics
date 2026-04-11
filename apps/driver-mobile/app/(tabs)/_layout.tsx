import { Tabs } from "expo-router";

import { colors } from "@/theme/colors";

export default function TabsLayout() {
  return (
    <Tabs
      screenOptions={{
        headerShown: false,
        tabBarActiveTintColor: colors.accent,
        tabBarInactiveTintColor: colors.muted,
        tabBarStyle: {
          backgroundColor: colors.surface,
          borderTopColor: colors.border,
        },
      }}
    >
      <Tabs.Screen name="loads" options={{ title: "Loads" }} />
      <Tabs.Screen name="operations" options={{ title: "Ops" }} />
      <Tabs.Screen name="tracking" options={{ title: "Track" }} />
      <Tabs.Screen name="account" options={{ title: "Account" }} />
    </Tabs>
  );
}
