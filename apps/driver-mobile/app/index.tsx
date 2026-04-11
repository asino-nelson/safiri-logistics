import { Redirect } from "expo-router";

import { LoadingScreen } from "@/screens/LoadingScreen";
import { useSession } from "@/context/session";

export default function IndexPage() {
  const { ready, token, user, kycProfile } = useSession();

  if (!ready) {
    return <LoadingScreen message="Preparing your driver workspace..." />;
  }

  if (!token || !user || user.role !== "driver") {
    return <Redirect href="/(auth)/login" />;
  }

  if (kycProfile?.status !== "approved") {
    return <Redirect href="/(auth)/kyc" />;
  }

  return <Redirect href="/(tabs)/loads" />;
}
