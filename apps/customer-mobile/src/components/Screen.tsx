import { PropsWithChildren } from "react";
import { RefreshControl, ScrollView, StyleSheet, Text, View } from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";

import { colors } from "../theme/colors";

type ScreenProps = PropsWithChildren<{
  scroll?: boolean;
  refreshable?: {
    refreshing: boolean;
    onRefresh: () => void;
  };
}>;

export function Screen({ children, scroll = false, refreshable }: ScreenProps) {
  return (
    <SafeAreaView style={styles.safeArea}>
      <View style={styles.bgBlobTop} />
      <View style={styles.bgBlobBottom} />
      {scroll ? (
        <ScrollView
          contentContainerStyle={styles.content}
          refreshControl={
            refreshable ? <RefreshControl refreshing={refreshable.refreshing} onRefresh={refreshable.onRefresh} tintColor={colors.accent} /> : undefined
          }
        >
          {children}
        </ScrollView>
      ) : (
        <View style={styles.content}>{children}</View>
      )}
    </SafeAreaView>
  );
}

export function LoadingScreen() {
  return (
    <SafeAreaView style={styles.safeArea}>
      <View style={[styles.content, { alignItems: "center", justifyContent: "center" }]}>
        <Text style={{ color: colors.muted }}>Loading Safiri Customer...</Text>
      </View>
    </SafeAreaView>
  );
}

export function LoadingState({ label }: { label: string }) {
  return (
    <View style={{ paddingVertical: 18 }}>
      <Text style={{ color: colors.muted }}>{label}</Text>
    </View>
  );
}

export function SectionTitle({ title, subtitle }: { title: string; subtitle?: string }) {
  return (
    <View style={{ marginTop: 24, marginBottom: 12 }}>
      <Text style={styles.sectionTitle}>{title}</Text>
      {subtitle ? <Text style={styles.sectionSubtitle}>{subtitle}</Text> : null}
    </View>
  );
}

const styles = StyleSheet.create({
  safeArea: {
    flex: 1,
    backgroundColor: colors.background,
  },
  content: {
    flexGrow: 1,
    padding: 20,
    gap: 18,
  },
  bgBlobTop: {
    position: "absolute",
    top: -40,
    right: -30,
    width: 180,
    height: 180,
    borderRadius: 90,
    backgroundColor: "#DCEFEA",
    opacity: 0.7,
  },
  bgBlobBottom: {
    position: "absolute",
    bottom: 80,
    left: -50,
    width: 160,
    height: 160,
    borderRadius: 80,
    backgroundColor: "#FDEAD6",
    opacity: 0.55,
  },
  sectionTitle: {
    color: colors.ink,
    fontSize: 22,
    fontWeight: "800",
  },
  sectionSubtitle: {
    color: colors.muted,
    marginTop: 6,
    lineHeight: 20,
  },
});
