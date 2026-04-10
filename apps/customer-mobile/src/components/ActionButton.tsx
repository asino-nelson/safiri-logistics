import { Pressable, StyleSheet, Text } from "react-native";

import { colors } from "../theme/colors";

type Variant = "primary" | "secondary" | "ghost" | "danger";

export function ActionButton({
  title,
  onPress,
  loading = false,
  variant = "primary",
}: {
  title: string;
  onPress: () => void;
  loading?: boolean;
  variant?: Variant;
}) {
  const palette = {
    primary: { backgroundColor: colors.accent, color: "#FFFFFF" },
    secondary: { backgroundColor: "#FFFFFF", color: colors.ink, borderColor: colors.border },
    ghost: { backgroundColor: "transparent", color: colors.accent, borderColor: "transparent" },
    danger: { backgroundColor: colors.danger, color: "#FFFFFF" },
  }[variant];

  return (
    <Pressable onPress={onPress} style={[styles.button, { backgroundColor: palette.backgroundColor, borderColor: palette.borderColor || "transparent" }]}>
      <Text style={[styles.text, { color: palette.color }]}>{loading ? "Please wait..." : title}</Text>
    </Pressable>
  );
}

const styles = StyleSheet.create({
  button: {
    borderRadius: 16,
    paddingVertical: 14,
    paddingHorizontal: 16,
    alignItems: "center",
    justifyContent: "center",
    borderWidth: 1,
  },
  text: {
    fontSize: 16,
    fontWeight: "800",
  },
});
