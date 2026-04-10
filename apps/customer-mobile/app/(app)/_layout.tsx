import { Redirect, Stack } from "expo-router";

import { LoadingScreen } from "../../src/components/Screen";
import { useAuth } from "../../src/context/auth";

export default function AppLayout() {
  const { isHydrating, session } = useAuth();

  if (isHydrating) {
    return <LoadingScreen />;
  }

  if (!session) {
    return <Redirect href="/login" />;
  }

  return <Stack screenOptions={{ headerShown: false }} />;
}
