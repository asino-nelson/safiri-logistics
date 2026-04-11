import { useLocalSearchParams } from "expo-router";
import { useEffect, useRef, useState } from "react";
import { Alert, ScrollView, StyleSheet, Text, View } from "react-native";

import { Button } from "@/components/Button";
import { Card } from "@/components/Card";
import { EmptyState } from "@/components/EmptyState";
import { Field } from "@/components/Field";
import { Screen } from "@/components/Screen";
import { SectionHeader } from "@/components/SectionHeader";
import { StatusPill } from "@/components/StatusPill";
import { api, ApiError } from "@/lib/api";
import { colors } from "@/theme/colors";
import { useSession } from "@/context/session";
import { Load, TrackingEvent } from "@/types/api";

export function TrackingScreen() {
  const params = useLocalSearchParams<{ loadId?: string }>();
  const { token, user } = useSession();
  const [loads, setLoads] = useState<Load[]>([]);
  const [selectedLoadId, setSelectedLoadId] = useState<string | null>(params.loadId ?? null);
  const [events, setEvents] = useState<TrackingEvent[]>([]);
  const [latitude, setLatitude] = useState("");
  const [longitude, setLongitude] = useState("");
  const [speedKph, setSpeedKph] = useState("0");
  const [headingDegrees, setHeadingDegrees] = useState("0");
  const [connected, setConnected] = useState(false);
  const socketRef = useRef<WebSocket | null>(null);

  const loadData = async () => {
    if (!token) {
      return;
    }

    const response = await api.listLoads(token);
    const assigned = response.data.filter((load) => load.assigned_driver_id === user?.id);
    setLoads(assigned);

    if (!selectedLoadId && assigned[0]) {
      setSelectedLoadId(assigned[0].id);
    }
  };

  useEffect(() => {
    void loadData();
  }, [token, user?.id]);

  useEffect(() => {
    if (!token || !selectedLoadId) {
      return;
    }

    socketRef.current?.close();
    // React Native accepts authenticated websocket headers in the third constructor argument.
    // @ts-expect-error React Native websocket header options are not in the DOM typings.
    const socket = new WebSocket(api.buildTrackingSocketUrl(selectedLoadId), undefined, {
      headers: { Authorization: `Bearer ${token}` },
    });

    socketRef.current = socket;
    setConnected(false);

    socket.onopen = () => setConnected(true);
    socket.onclose = () => setConnected(false);
    socket.onerror = () => setConnected(false);
    socket.onmessage = (event) => {
      try {
        const payload = JSON.parse(String(event.data)) as TrackingEvent;
        setEvents((current) => [payload, ...current].slice(0, 20));
      } catch {
        // ignore malformed socket payloads
      }
    };

    return () => {
      socket.close();
    };
  }, [token, selectedLoadId]);

  const sendPing = async () => {
    if (!token || !selectedLoadId) {
      return;
    }

    try {
      const event = await api.createTrackingEvent(token, selectedLoadId, {
        latitude: Number(latitude),
        longitude: Number(longitude),
        speed_kph: Number(speedKph) || 0,
        heading_degrees: Number(headingDegrees) || 0,
      });
      setEvents((current) => [event, ...current].slice(0, 20));
      Alert.alert("Tracking ping sent", "Your live location is now recorded.");
    } catch (err) {
      const message = err instanceof ApiError || err instanceof Error ? err.message : "Tracking ping failed.";
      Alert.alert("Tracking error", message);
    }
  };

  const currentLoad = loads.find((load) => load.id === selectedLoadId) ?? null;

  return (
    <Screen
      title="Live tracking"
      subtitle="Publish pings over HTTP and stream updates over WebSockets for live visibility."
    >
      <ScrollView style={styles.scroll}>
        <View style={styles.stack}>
          <Card>
            <SectionHeader title="Connected load" />
            {loads.length === 0 ? (
              <EmptyState title="No assigned loads" message="Tracking becomes active when a load is matched or picked." />
            ) : (
              <View style={styles.chips}>
                {loads.map((load) => (
                  <Button
                    key={load.id}
                    label={load.title}
                    onPress={() => setSelectedLoadId(load.id)}
                    tone={selectedLoadId === load.id ? "primary" : "secondary"}
                    style={styles.chip}
                  />
                ))}
              </View>
            )}
            {currentLoad ? (
              <View style={styles.summary}>
                <Text style={styles.meta}>{currentLoad.origin} -> {currentLoad.destination}</Text>
                <StatusPill label={currentLoad.status} tone={currentLoad.status === "in_transit" ? "info" : "warning"} />
              </View>
            ) : null}
            <Text style={styles.meta}>Socket: {connected ? "live" : "waiting"}</Text>
          </Card>

          <Card>
            <SectionHeader title="Send tracking ping" />
            <Field label="Latitude" value={latitude} onChangeText={setLatitude} keyboardType="numeric" />
            <Field label="Longitude" value={longitude} onChangeText={setLongitude} keyboardType="numeric" />
            <Field label="Speed (kph)" value={speedKph} onChangeText={setSpeedKph} keyboardType="numeric" />
            <Field label="Heading (degrees)" value={headingDegrees} onChangeText={setHeadingDegrees} keyboardType="numeric" />
            <Button label="Send ping" onPress={sendPing} disabled={!selectedLoadId} />
          </Card>

          <Card>
            <SectionHeader title="Recent events" />
            {events.length === 0 ? (
              <EmptyState title="No events yet" message="The newest tracking updates will appear here as they arrive." />
            ) : (
              <View>
                {events.map((item) => (
                  <View key={item.id} style={styles.event}>
                    <Text style={styles.eventTitle}>{new Date(item.recorded_at).toLocaleString()}</Text>
                    <Text style={styles.meta}>
                      {item.latitude.toFixed(5)}, {item.longitude.toFixed(5)}
                    </Text>
                    <Text style={styles.meta}>Speed {item.speed_kph} kph | Heading {item.heading_degrees} deg</Text>
                  </View>
                ))}
              </View>
            )}
          </Card>
        </View>
      </ScrollView>
    </Screen>
  );
}

const styles = StyleSheet.create({
  scroll: {
    flex: 1,
  },
  stack: {
    gap: 14,
  },
  chips: {
    flexDirection: "row",
    flexWrap: "wrap",
    gap: 10,
  },
  chip: {
    flexGrow: 0,
  },
  summary: {
    gap: 8,
    marginTop: 12,
  },
  meta: {
    color: colors.muted,
    lineHeight: 20,
  },
  event: {
    paddingVertical: 10,
    borderBottomWidth: 1,
    borderBottomColor: colors.border,
    gap: 4,
  },
  eventTitle: {
    color: colors.text,
    fontWeight: "800",
  },
});
