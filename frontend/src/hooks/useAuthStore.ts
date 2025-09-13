import { useCallback } from "react";
import { create } from "zustand";
import { createComputed } from "zustand-computed";
import { createJSONStorage, persist } from "zustand/middleware";
import { api } from "@/lib/api";

interface User {
  id: number;
  email: string;
  username: string;
  full_name: string;
  is_active: boolean;
}

interface RegisterData {
  email: string;
  username: string;
  password: string;
  full_name: string;
}

interface AuthTokens {
  accessToken: string;
  refreshToken: string;
  expiresAt: string;
}

interface AuthStore {
  accessToken: string | null;
  refreshToken: string | null;
  expiresAt: string | null;
  user: User | null;
  setTokens: (tokens: AuthTokens) => void;
  setUser: (user: User) => void;
  login: (email: string, password: string) => Promise<void>;
  register: (data: RegisterData) => Promise<void>;
  refreshAccessToken: () => Promise<boolean>;
  logout: () => void;
}

interface AuthStoreComputed {
  isLoggedIn: boolean;
}

const computed = createComputed<AuthStore, AuthStoreComputed>((state) => ({
  isLoggedIn: !!state.accessToken,
}));

export const authStore = create<AuthStore>()(
  computed(
    persist(
      (set, get) => ({
        accessToken: null,
        refreshToken: null,
        expiresAt: null,
        user: null,
        setTokens: (tokens: AuthTokens) =>
          set({
            accessToken: tokens.accessToken,
            refreshToken: tokens.refreshToken,
            expiresAt: tokens.expiresAt,
          }),
        setUser: (user: User) => set({ user }),
        login: async (email: string, password: string) => {
          const response = await api.post("/auth/login", {
            email,
            password,
          });

          const { token, refresh_token, expires_at, user } = response.data;
          set({
            accessToken: token,
            refreshToken: refresh_token,
            expiresAt: expires_at,
            user,
          });
        },
        register: async (data: RegisterData) => {
          const response = await api.post("/auth/register", data);

          const { token, refresh_token, expires_at, user } = response.data;
          set({
            accessToken: token,
            refreshToken: refresh_token,
            expiresAt: expires_at,
            user,
          });
        },
        refreshAccessToken: async () => {
          const { refreshToken } = get();
          if (!refreshToken) {
            return false;
          }

          try {
            const response = await api.post("/auth/refresh", {
              refresh_token: refreshToken,
            });

            const { token, refresh_token, expires_at, user } = response.data;
            set({
              accessToken: token,
              refreshToken: refresh_token,
              expiresAt: expires_at,
              user,
            });
            return true;
          } catch (error: unknown) {
            console.error("Token refresh failed:", error);
            // Refresh failed, logout user
            set({
              accessToken: null,
              refreshToken: null,
              expiresAt: null,
              user: null,
            });
            return false;
          }
        },
        logout: () =>
          set({
            accessToken: null,
            refreshToken: null,
            expiresAt: null,
            user: null,
          }),
      }),
      {
        name: "auth",
        storage: createJSONStorage(() => localStorage),
      }
    )
  )
);

export const useAuthStore = () => {
  const store = authStore();

  const login = useCallback(
    (email: string, password: string) => {
      return store.login(email, password);
    },
    [store]
  );

  const register = useCallback(
    (data: RegisterData) => {
      return store.register(data);
    },
    [store]
  );

  const refreshAccessToken = useCallback(() => {
    return store.refreshAccessToken();
  }, [store]);

  const logout = useCallback(() => {
    store.logout();
  }, [store]);

  return {
    ...store,
    login,
    register,
    refreshAccessToken,
    logout,
  };
};
