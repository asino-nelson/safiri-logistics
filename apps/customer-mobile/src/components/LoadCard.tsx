import { Pressable, StyleSheet, Text, View } from "react-native";

import { Load } from "../lib/types";
import { colors } from "../theme/colors";
import { StatusPill } from "./StatusPill";

export function LoadCard({
  load,
  onPress,
  compact = false,
}: {
  load: Load;
  onPress?: () => void;
  compact?: boolean;
}) {
  if (onPress) {
    return (
      <Pressable onPress={onPress} style={styles.card}>
        <View style={styles.topRow}>
          <View style={{ flex: 1 }}>
            <Text style={styles.title}>{load.title}</Text>
            <Text style={styles.route}>
              {load.origin} -> {load.destination}
            </Text>
          </View>
          <StatusPill status={load.status} />
        </View>

        <Text style={styles.body} numberOfLines={compact ? 2 : 4}>
          {load.description}
        </Text>

        <View style={styles.metaRow}>
          <Text style={styles.meta}>Priority {load.priority}</Text>
          <Text style={styles.meta}>KES {load.quoted_price_kes.toLocaleString()}</Text>
        </View>
      </Pressable>
    );
  }

  return (
    <View style={styles.card}>
      <View style={styles.topRow}>
        <View style={{ flex: 1 }}>
          <Text style={styles.title}>{load.title}</Text>
          <Text style={styles.route}>
            {load.origin} -> {load.destination}
          </Text>
        </View>
        <StatusPill status={load.status} />
      </View>

      <Text style={styles.body} numberOfLines={compact ? 2 : 4}>
        {load.description}
      </Text>

      <View style={styles.metaRow}>
        <Text style={styles.meta}>Priority {load.priority}</Text>
        <Text style={styles.meta}>KES {load.quoted_price_kes.toLocaleString()}</Text>
      </View>
    </View>
  );
}

const styles = StyleSheet.create({
  card: {
    backgroundColor: colors.surface,
    borderRadius: 22,
    borderWidth: 1,
    borderColor: colors.border,
    padding: 16,
    gap: 12,
    shadowColor: "#0F172A",
    shadowOpacity: 0.06,
    shadowRadius: 18,
    shadowOffset: { width: 0, height: 8 },
    elevation: 2,
  },
  topRow: {
    flexDirection: "row",
    gap: 12,
    alignItems: "flex-start",
  },
  title: {
    color: colors.ink,
    fontSize: 18,
    fontWeight: "800",
  },
  route: {
    color: colors.accent,
    marginTop: 4,
    fontWeight: "700",
  },
  body: {
    color: colors.muted,
    lineHeight: 20,
  },
  metaRow: {
    flexDirection: "row",
    justifyContent: "space-between",
    gap: 12,
  },
  meta: {
    color: colors.ink,
    fontWeight: "700",
  },
});
