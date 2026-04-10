import { Redirect } from "expo-router";
import { ActivityIndicator, StyleSheet, View } from "react-native";

import { useAuth } from "../src/context/auth";
import { colors } from "../src/theme/colors";

export default function IndexScreen() {
  const { isHydrating, session } = useAuth();

  if (isHydrating) {
    return (
      <View style={styles.loading}>
        <ActivityIndicator size="large" color={colors.accent} />
      </View>
    );
  }

  return <Redirect href={session ? "/dashboard" : "/login"} />;
}

const styles = StyleSheet.create({
  loading: {
    flex: 1,
    alignItems: "center",
    justifyContent: "center",
    backgroundColor: colors.background,
  },
});
