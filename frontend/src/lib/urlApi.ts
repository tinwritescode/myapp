import { api } from "./api";
import type {
  CreateURLRequest,
  UpdateURLRequest,
  URLResponse,
  URLsResponse,
  URLStatsResponse,
} from "@/types/url";

export const urlApi = {
  // Create a new URL (authenticated)
  createURL: async (data: CreateURLRequest): Promise<URLResponse> => {
    const response = await api.post("/urls", data);
    return response.data;
  },

  // Create a new URL (public, no authentication)
  createPublicURL: async (data: CreateURLRequest): Promise<URLResponse> => {
    const response = await api.post("/urls/public", data);
    return response.data;
  },

  // Get URLs with pagination and filtering
  getURLs: async (
    params: {
      page?: number;
      limit?: number;
      search?: string;
      is_active?: boolean;
      sort_by?: string;
      sort_dir?: "asc" | "desc";
    } = {}
  ): Promise<URLsResponse> => {
    const response = await api.get("/urls", { params });
    return response.data;
  },

  // Get a specific URL by ID
  getURL: async (id: number): Promise<URLResponse> => {
    const response = await api.get(`/urls/${id}`);
    return response.data;
  },

  // Update a URL
  updateURL: async (
    id: number,
    data: UpdateURLRequest
  ): Promise<URLResponse> => {
    const response = await api.put(`/urls/${id}`, data);
    return response.data;
  },

  // Delete a URL
  deleteURL: async (
    id: number
  ): Promise<{ success: boolean; message: string }> => {
    const response = await api.delete(`/urls/${id}`);
    return response.data;
  },

  // Get URL statistics
  getURLStats: async (id: number): Promise<URLStatsResponse> => {
    const response = await api.get(`/urls/${id}/stats`);
    return response.data;
  },
};
