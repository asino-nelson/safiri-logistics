import { forwardRef } from "react";
import { StyleSheet, Text, TextInput, TextInputProps, View } from "react-native";

import { colors } from "@/theme/colors";
import { spacing } from "@/theme/spacing";

type FieldProps = TextInputProps & {
  label: string;
  helperText?: string;
};

export const Field = forwardRef<TextInput, FieldProps>(function Field(
  { label, helperText, style, ...props },
  ref,
) {
  return (
    <View style={styles.wrapper}>
      <Text style={styles.label}>{label}</Text>
      <TextInput ref={ref} placeholderTextColor={colors.muted} style={[styles.input, style]} {...props} />
      {helperText ? <Text style={styles.helper}>{helperText}</Text> : null}
    </View>
  );
});

const styles = StyleSheet.create({
  wrapper: {
    gap: 8,
  },
  label: {
    color: colors.text,
    fontSize: 13,
    fontWeight: "700",
  },
  input: {
    minHeight: 50,
    borderRadius: 14,
    paddingHorizontal: spacing.md,
    backgroundColor: colors.surfaceElevated,
    borderWidth: 1,
    borderColor: colors.border,
    color: colors.text,
  },
  helper: {
    color: colors.muted,
    fontSize: 12,
    lineHeight: 17,
  },
});
