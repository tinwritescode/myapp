import axios from "axios";

// Create axios instance
export const api = axios.create({
  baseURL: "http://localhost:8080/api/v1",
  headers: {
    "Content-Type": "application/json",
  },
});

// Function to get auth store instance
const getAuthStore = (): ReturnType<typeof useAuthStore> => {
  // Import dynamically to avoid circular dependency
  // eslint-disable-next-line @typescript-eslint/no-require-imports
  const { useAuthStore } = require("@/hooks/useAuthStore");
  return useAuthStore.getState();
};

// Function to clear auth data
const clearAuthData = () => {
  const authStore = getAuthStore();
  authStore.logout();
  window.location.href = "/";
};

// Function to refresh access token
const refreshAccessToken = async () => {
  const authStore = getAuthStore();
  const { refreshToken } = authStore;

  if (!refreshToken) {
    return false;
  }

  try {
    const response = await axios.post(
      "http://localhost:8080/api/v1/auth/refresh",
      {
        refresh_token: refreshToken,
      }
    );

    const { token, refresh_token, expires_at, user } = response.data;

    // Update auth store with new tokens
    authStore.setTokens({
      accessToken: token,
      refreshToken: refresh_token,
      expiresAt: expires_at,
    });
    authStore.setUser(user);

    return true;
  } catch (error) {
    console.error("Token refresh failed:", error);
    clearAuthData();
    return false;
  }
};

// Request interceptor to add auth token and handle token refresh
api.interceptors.request.use(
  async (config) => {
    // Skip token refresh for auth endpoints
    if (config.url?.includes("/auth/")) {
      return config;
    }

    const authStore = getAuthStore();
    const { accessToken, expiresAt } = authStore;

    if (accessToken) {
      // Check if token is expired and refresh if needed
      if (expiresAt && new Date(expiresAt) <= new Date()) {
        const refreshed = await refreshAccessToken();
        if (refreshed) {
          // Get the new token after refresh
          const updatedAuthStore = getAuthStore();
          if (updatedAuthStore.accessToken) {
            config.headers.Authorization = `Bearer ${updatedAuthStore.accessToken}`;
          }
        } else {
          // Refresh failed, redirect to login
          clearAuthData();
          return Promise.reject(new Error("Token refresh failed"));
        }
      } else {
        config.headers.Authorization = `Bearer ${accessToken}`;
      }
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response interceptor to handle auth errors and automatic token refresh
api.interceptors.response.use(
  (response) => {
    return response;
  },
  async (error) => {
    const originalRequest = error.config;

    // If we get a 401 and haven't already tried to refresh
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;

      // Try to refresh the token
      const refreshed = await refreshAccessToken();
      if (refreshed) {
        // Retry the original request with new token
        const authStore = getAuthStore();
        if (authStore.accessToken) {
          originalRequest.headers.Authorization = `Bearer ${authStore.accessToken}`;
          return api(originalRequest);
        }
      } else {
        // Refresh failed, clear auth and redirect
        clearAuthData();
      }
    }

    return Promise.reject(error);
  }
);
