import { View, Text, StyleSheet } from "react-native";

import { colors } from "../theme/colors";

export function BrandHeader({
  eyebrow,
  title,
  subtitle,
}: {
  eyebrow: string;
  title: string;
  subtitle?: string;
}) {
  return (
    <View style={styles.wrapper}>
      <Text style={styles.eyebrow}>{eyebrow}</Text>
      <Text style={styles.title}>{title}</Text>
      {subtitle ? <Text style={styles.subtitle}>{subtitle}</Text> : null}
    </View>
  );
}

const styles = StyleSheet.create({
  wrapper: {
    paddingVertical: 6,
  },
  eyebrow: {
    color: colors.accent,
    fontSize: 12,
    textTransform: "uppercase",
    letterSpacing: 1.2,
    fontWeight: "800",
    marginBottom: 8,
  },
  title: {
    color: colors.ink,
    fontSize: 30,
    fontWeight: "900",
    lineHeight: 36,
  },
  subtitle: {
    color: colors.muted,
    marginTop: 12,
    fontSize: 15,
    lineHeight: 22,
  },
});
