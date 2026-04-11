import { router, useLocalSearchParams } from "expo-router";
import { useEffect, useRef, useState } from "react";
import { Text, View } from "react-native";

import { ActionButton } from "../../../src/components/ActionButton";
import { BrandHeader } from "../../../src/components/BrandHeader";
import { Screen } from "../../../src/components/Screen";
import { useAuth } from "../../../src/context/auth";
import { wsUrl } from "../../../src/lib/env";
import { getLoadById } from "../../../src/lib/loads";
import { Load, TrackingEvent } from "../../../src/lib/types";

export default function TrackingScreen() {
  const { id } = useLocalSearchParams<{ id: string }>();
  const { session } = useAuth();
  const [load, setLoad] = useState<Load | null>(null);
  const [events, setEvents] = useState<TrackingEvent[]>([]);
  const [connectionState, setConnectionState] = useState<"connecting" | "live" | "closed" | "error">("connecting");
  const socketRef = useRef<WebSocket | null>(null);

  useEffect(() => {
    let mounted = true;
    const run = async () => {
      if (!session || !id) {
        return;
      }

      try {
        const currentLoad = await getLoadById(session.token, id);
        if (mounted) {
          setLoad(currentLoad);
        }
      } catch {
        if (mounted) {
          setLoad(null);
        }
      }
    };

    void run();
    return () => {
      mounted = false;
    };
  }, [session?.token, id]);

  useEffect(() => {
    if (!session || !id) {
      return;
    }

    const socket = new (WebSocket as any)(`${wsUrl}/tracking/ws?load_id=${encodeURIComponent(id)}`, undefined, {
      headers: {
        Authorization: `Bearer ${session.token}`,
      },
    }) as WebSocket;

    socketRef.current = socket;
    setConnectionState("connecting");

    socket.onopen = () => setConnectionState("live");
    socket.onerror = () => setConnectionState("error");
    socket.onclose = () => setConnectionState((current) => (current === "error" ? current : "closed"));
    socket.onmessage = (event) => {
      try {
        const parsed = JSON.parse(String(event.data)) as TrackingEvent;
        setEvents((current) => [parsed, ...current].slice(0, 20));
      } catch {
        // ignore malformed messages
      }
    };

    return () => {
      socket.close();
      socketRef.current = null;
    };
  }, [session?.token, id]);

  const latest = events[0];

  return (
    <Screen scroll>
      <ActionButton title="Back" variant="secondary" onPress={() => router.back()} />
      <BrandHeader
        eyebrow="Live tracking"
        title={load ? load.title : "Tracking stream"}
        subtitle="WebSocket updates for the assigned shipment, with the latest coordinate feed at the top."
      />

      <View style={{ backgroundColor: "#FFFFFF", borderRadius: 20, padding: 16, borderWidth: 1, borderColor: "#E2E8F0", gap: 10 }}>
        <Text style={{ color: "#64748B", textTransform: "uppercase", letterSpacing: 0.8, fontSize: 12 }}>Connection</Text>
        <Text style={{ color: "#0F172A", fontSize: 18, fontWeight: "800" }}>{connectionState}</Text>
        <Text style={{ color: "#64748B" }}>
          If you are on a device and localhost is not reachable, point EXPO_PUBLIC_WS_URL to your machine IP.
        </Text>
      </View>

      <View style={{ gap: 14, marginTop: 16 }}>
        <Row label="Latest position" value={latest ? `${latest.latitude.toFixed(5)}, ${latest.longitude.toFixed(5)}` : "Waiting for updates"} />
        <Row label="Speed" value={latest ? `${latest.speed_kph.toFixed(1)} kph` : "-"} />
        <Row label="Heading" value={latest ? `${latest.heading_degrees.toFixed(1)} deg` : "-"} />
        <Row label="Recorded at" value={latest ? new Date(latest.recorded_at).toLocaleString() : "-"} />
      </View>

      <View style={{ marginTop: 18 }}>
        <Text style={{ color: "#0F172A", fontSize: 18, fontWeight: "800", marginBottom: 10 }}>Recent events</Text>
        <View style={{ gap: 12 }}>
          {events.map((event) => (
            <View key={event.id} style={{ backgroundColor: "#FFFFFF", borderWidth: 1, borderColor: "#E2E8F0", borderRadius: 18, padding: 14 }}>
              <Text style={{ color: "#0F172A", fontWeight: "700" }}>{new Date(event.recorded_at).toLocaleString()}</Text>
              <Text style={{ color: "#64748B", marginTop: 6 }}>
                {event.latitude.toFixed(5)}, {event.longitude.toFixed(5)}
              </Text>
              <Text style={{ color: "#64748B" }}>
                Speed {event.speed_kph.toFixed(1)} kph, heading {event.heading_degrees.toFixed(1)} deg
              </Text>
            </View>
          ))}
          {events.length === 0 ? <Text style={{ color: "#64748B" }}>No tracking events yet. The latest pings will appear here live.</Text> : null}
        </View>
      </View>
    </Screen>
  );
}

function Row({ label, value }: { label: string; value: string }) {
  return (
    <View style={{ backgroundColor: "#FFFFFF", borderWidth: 1, borderColor: "#E2E8F0", borderRadius: 18, padding: 14 }}>
      <Text style={{ color: "#64748B", textTransform: "uppercase", letterSpacing: 0.8, fontSize: 12 }}>{label}</Text>
      <Text style={{ color: "#0F172A", fontSize: 16, fontWeight: "700", marginTop: 6 }}>{value}</Text>
    </View>
  );
}
