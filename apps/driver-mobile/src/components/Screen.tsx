import { LinearGradient } from "expo-linear-gradient";
import type { PropsWithChildren, ReactNode } from "react";
import { StyleSheet, Text, View } from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";

import { colors } from "@/theme/colors";
import { spacing } from "@/theme/spacing";

type ScreenProps = PropsWithChildren<{
  title: string;
  subtitle?: string;
  footer?: ReactNode;
}>;

export function Screen({ title, subtitle, footer, children }: ScreenProps) {
  return (
    <LinearGradient colors={[colors.background, "#09192e", "#0d1a2d"]} style={styles.flex}>
      <SafeAreaView style={styles.flex}>
        <View style={styles.container}>
          <View style={styles.hero}>
            <Text style={styles.title}>{title}</Text>
            {subtitle ? <Text style={styles.subtitle}>{subtitle}</Text> : null}
          </View>
          <View style={styles.body}>{children}</View>
          {footer ? <View style={styles.footer}>{footer}</View> : null}
        </View>
      </SafeAreaView>
    </LinearGradient>
  );
}

const styles = StyleSheet.create({
  flex: {
    flex: 1,
  },
  container: {
    padding: spacing.lg,
    gap: spacing.md,
    flex: 1,
  },
  body: {
    flex: 1,
    gap: spacing.md,
  },
  hero: {
    gap: spacing.sm,
  },
  title: {
    color: colors.text,
    fontSize: 34,
    lineHeight: 40,
    fontWeight: "800",
    letterSpacing: -0.8,
  },
  subtitle: {
    color: colors.muted,
    fontSize: 15,
    lineHeight: 22,
  },
  footer: {
    paddingBottom: spacing.xl,
  },
});
