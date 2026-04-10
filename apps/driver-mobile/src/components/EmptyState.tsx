import { StyleSheet, Text, View } from "react-native";

import { colors } from "@/theme/colors";

export function EmptyState({
  title,
  message,
}: {
  title: string;
  message: string;
}) {
  return (
    <View style={styles.box}>
      <Text style={styles.title}>{title}</Text>
      <Text style={styles.message}>{message}</Text>
    </View>
  );
}

const styles = StyleSheet.create({
  box: {
    padding: 20,
    borderRadius: 18,
    backgroundColor: colors.surface,
    borderColor: colors.border,
    borderWidth: 1,
    gap: 8,
  },
  title: {
    color: colors.text,
    fontSize: 16,
    fontWeight: "800",
  },
  message: {
    color: colors.muted,
    lineHeight: 20,
  },
});
