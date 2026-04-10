import AsyncStorage from "@react-native-async-storage/async-storage";

const sessionKey = "safiri.customer.session";

export async function readSession<T>() {
  const raw = await AsyncStorage.getItem(sessionKey);
  if (!raw) {
    return null;
  }

  try {
    return JSON.parse(raw) as T;
  } catch {
    return null;
  }
}

export async function writeSession(value: unknown) {
  await AsyncStorage.setItem(sessionKey, JSON.stringify(value));
}

export async function clearSession() {
  await AsyncStorage.removeItem(sessionKey);
}
