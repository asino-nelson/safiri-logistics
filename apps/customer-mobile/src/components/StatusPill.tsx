import { StyleSheet, Text, View } from "react-native";

import { colors } from "../theme/colors";

const palette: Record<string, { backgroundColor: string; color: string }> = {
  posted: { backgroundColor: "#E0F2FE", color: "#075985" },
  matched: { backgroundColor: "#DCFCE7", color: "#166534" },
  picked: { backgroundColor: "#FEF3C7", color: "#92400E" },
  in_transit: { backgroundColor: "#E9D5FF", color: "#6B21A8" },
  delivered: { backgroundColor: "#DCFCE7", color: "#166534" },
};

export function StatusPill({ status }: { status: string }) {
  const theme = palette[status] || { backgroundColor: colors.surfaceMuted, color: colors.ink };

  return (
    <View style={[styles.pill, { backgroundColor: theme.backgroundColor }]}>
      <Text style={[styles.text, { color: theme.color }]}>{status.replaceAll("_", " ")}</Text>
    </View>
  );
}

const styles = StyleSheet.create({
  pill: {
    paddingHorizontal: 12,
    paddingVertical: 8,
    borderRadius: 999,
  },
  text: {
    fontSize: 12,
    fontWeight: "800",
    textTransform: "uppercase",
    letterSpacing: 0.6,
  },
});
