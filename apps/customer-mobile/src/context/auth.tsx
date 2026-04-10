import { createContext, PropsWithChildren, useContext, useEffect, useState } from "react";

import { login, me, register, SignUpInput } from "../lib/auth";
import { clearSession, readSession, writeSession } from "../lib/storage";
import { AuthResponse } from "../lib/types";

type Session = AuthResponse;

type AuthContextValue = {
  session: Session | null;
  isHydrating: boolean;
  signIn: (email: string, password: string) => Promise<void>;
  signUp: (input: SignUpInput) => Promise<void>;
  signOut: () => Promise<void>;
  refreshProfile: () => Promise<void>;
};

const AuthContext = createContext<AuthContextValue | undefined>(undefined);

export function AuthProvider({ children }: PropsWithChildren) {
  const [session, setSession] = useState<Session | null>(null);
  const [isHydrating, setIsHydrating] = useState(true);

  useEffect(() => {
    const hydrate = async () => {
      const stored = await readSession<Session>();
      if (stored) {
        setSession(stored);
      }
      setIsHydrating(false);
    };

    void hydrate();
  }, []);

  const persist = async (nextSession: Session | null) => {
    setSession(nextSession);
    if (nextSession) {
      await writeSession(nextSession);
    } else {
      await clearSession();
    }
  };

  const signIn = async (email: string, password: string) => {
    const response = await login(email, password);
    await persist(response);
  };

  const signUp = async (input: SignUpInput) => {
    const response = await register(input);
    await persist(response);
  };

  const signOut = async () => {
    await persist(null);
  };

  const refreshProfile = async () => {
    if (!session) {
      return;
    }

    const current = await me(session.token);
    await persist({ ...session, user: current });
  };

  return (
    <AuthContext.Provider value={{ session, isHydrating, signIn, signUp, signOut, refreshProfile }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const value = useContext(AuthContext);
  if (!value) {
    throw new Error("useAuth must be used inside AuthProvider");
  }

  return value;
}
