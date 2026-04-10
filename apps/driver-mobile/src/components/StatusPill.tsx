import { StyleSheet, Text, View } from "react-native";

import { colors } from "@/theme/colors";

type StatusPillProps = {
  label: string;
  tone?: "default" | "success" | "warning" | "danger" | "info";
};

export function StatusPill({ label, tone = "default" }: StatusPillProps) {
  return (
    <View style={[styles.base, toneStyles[tone]]}>
      <Text style={styles.label}>{label}</Text>
    </View>
  );
}

const styles = StyleSheet.create({
  base: {
    alignSelf: "flex-start",
    paddingHorizontal: 12,
    paddingVertical: 7,
    borderRadius: 999,
    backgroundColor: colors.surfaceElevated,
  },
  label: {
    color: colors.text,
    fontSize: 12,
    fontWeight: "800",
    textTransform: "uppercase",
    letterSpacing: 0.4,
  },
});

const toneStyles = StyleSheet.create({
  default: {
    backgroundColor: colors.surfaceElevated,
  },
  success: {
    backgroundColor: colors.accentSoft,
  },
  warning: {
    backgroundColor: "rgba(246,198,91,0.16)",
  },
  danger: {
    backgroundColor: "rgba(255,122,122,0.16)",
  },
  info: {
    backgroundColor: "rgba(103,168,255,0.16)",
  },
});
