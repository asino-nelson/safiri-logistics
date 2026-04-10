import { PropsWithChildren } from "react";
import { StyleSheet, View } from "react-native";

import { colors } from "@/theme/colors";
import { spacing } from "@/theme/spacing";

export function Card({ children }: PropsWithChildren) {
  return <View style={styles.card}>{children}</View>;
}

const styles = StyleSheet.create({
  card: {
    backgroundColor: colors.surface,
    borderColor: colors.border,
    borderWidth: 1,
    borderRadius: 22,
    padding: spacing.md,
    gap: spacing.sm,
  },
});
