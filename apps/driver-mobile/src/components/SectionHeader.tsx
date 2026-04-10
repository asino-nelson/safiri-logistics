import { ReactNode } from "react";
import { StyleSheet, Text, View } from "react-native";

import { colors } from "@/theme/colors";

type SectionHeaderProps = {
  title: string;
  right?: ReactNode;
};

export function SectionHeader({ title, right }: SectionHeaderProps) {
  return (
    <View style={styles.row}>
      <Text style={styles.title}>{title}</Text>
      {right}
    </View>
  );
}

const styles = StyleSheet.create({
  row: {
    flexDirection: "row",
    alignItems: "center",
    justifyContent: "space-between",
  },
  title: {
    color: colors.text,
    fontSize: 18,
    fontWeight: "800",
  },
});
