import { Pressable, StyleSheet, Text, ViewStyle } from "react-native";

import { colors } from "@/theme/colors";
import { spacing } from "@/theme/spacing";

type ButtonProps = {
  label: string;
  onPress: () => void;
  tone?: "primary" | "secondary" | "danger";
  style?: ViewStyle;
  disabled?: boolean;
};

export function Button({ label, onPress, tone = "primary", style, disabled }: ButtonProps) {
  return (
    <Pressable
      onPress={onPress}
      disabled={disabled}
      style={({ pressed }) => [
        styles.base,
        tone === "primary" && styles.primary,
        tone === "secondary" && styles.secondary,
        tone === "danger" && styles.danger,
        pressed && !disabled && styles.pressed,
        disabled && styles.disabled,
        style,
      ]}
    >
      <Text style={[styles.label, tone === "secondary" && styles.secondaryLabel]}>{label}</Text>
    </Pressable>
  );
}

const styles = StyleSheet.create({
  base: {
    minHeight: 50,
    borderRadius: 16,
    paddingHorizontal: spacing.md,
    alignItems: "center",
    justifyContent: "center",
  },
  primary: {
    backgroundColor: colors.accent,
  },
  secondary: {
    backgroundColor: colors.surfaceElevated,
    borderColor: colors.border,
    borderWidth: 1,
  },
  danger: {
    backgroundColor: colors.danger,
  },
  pressed: {
    opacity: 0.82,
    transform: [{ scale: 0.99 }],
  },
  disabled: {
    opacity: 0.55,
  },
  label: {
    color: "#08111d",
    fontSize: 15,
    fontWeight: "800",
    letterSpacing: 0.2,
  },
  secondaryLabel: {
    color: colors.text,
  },
});
