import AsyncStorage from "@react-native-async-storage/async-storage";
import type { ReactNode } from "react";
import { createContext, useContext, useEffect, useState } from "react";

import { api } from "@/lib/api";
import { storageKeys } from "@/lib/storage";
import { AuthUser, KycProfile, LoginInput } from "@/types/api";

type SessionContextValue = {
  ready: boolean;
  token: string | null;
  user: AuthUser | null;
  kycProfile: KycProfile | null;
  signIn: (input: LoginInput) => Promise<void>;
  signOut: () => Promise<void>;
  refreshProfile: () => Promise<void>;
};

const SessionContext = createContext<SessionContextValue | null>(null);

export function SessionProvider({ children }: { children: ReactNode }) {
  const [ready, setReady] = useState(false);
  const [token, setToken] = useState<string | null>(null);
  const [user, setUser] = useState<AuthUser | null>(null);
  const [kycProfile, setKycProfile] = useState<KycProfile | null>(null);

  const refreshProfileInternal = async (currentToken: string) => {
    try {
      const profile = await api.getKycProfile(currentToken);
      setKycProfile(profile);
      return profile;
    } catch {
      setKycProfile(null);
      return null;
    }
  };

  useEffect(() => {
    let mounted = true;

    const hydrate = async () => {
      try {
        const storedToken = await AsyncStorage.getItem(storageKeys.token);
        const storedUser = await AsyncStorage.getItem(storageKeys.user);

        if (!mounted) {
          return;
        }

        if (!storedToken || !storedUser) {
          setReady(true);
          return;
        }

        let parsedUser: AuthUser;
        try {
          parsedUser = JSON.parse(storedUser) as AuthUser;
        } catch {
          await AsyncStorage.multiRemove([storageKeys.token, storageKeys.user]);
          setReady(true);
          return;
        }
        if (parsedUser.role !== "driver") {
          await AsyncStorage.multiRemove([storageKeys.token, storageKeys.user]);
          setReady(true);
          return;
        }

        setToken(storedToken);
        setUser(parsedUser);

        try {
          await refreshProfileInternal(storedToken);
        } catch {
          setKycProfile(null);
        }
      } finally {
        if (mounted) {
          setReady(true);
        }
      }
    };

    void hydrate();

    return () => {
      mounted = false;
    };
  }, []);

  const signIn = async (input: LoginInput) => {
    const response = await api.login(input);

    if (response.user.role !== "driver") {
      throw new Error("This app is for drivers only.");
    }

    await AsyncStorage.multiSet([
      [storageKeys.token, response.token],
      [storageKeys.user, JSON.stringify(response.user)],
    ]);

    setToken(response.token);
    setUser(response.user);
    await refreshProfileInternal(response.token);
  };

  const signOut = async () => {
    await AsyncStorage.multiRemove([storageKeys.token, storageKeys.user]);
    setToken(null);
    setUser(null);
    setKycProfile(null);
  };

  const refreshProfile = async () => {
    if (!token) {
      return;
    }

    await refreshProfileInternal(token);
  };

  return (
    <SessionContext.Provider
      value={{
        ready,
        token,
        user,
        kycProfile,
        signIn,
        signOut,
        refreshProfile,
      }}
    >
      {children}
    </SessionContext.Provider>
  );
}

export function useSession() {
  const value = useContext(SessionContext);
  if (!value) {
    throw new Error("useSession must be used within SessionProvider");
  }

  return value;
}
