export interface URL {
  id: number;
  original_url: string;
  short_code: string;
  user_id?: number;
  expires_at?: string;
  click_count: number;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface CreateURLRequest {
  original_url: string;
  short_code?: string;
  expires_at?: string;
}

export interface UpdateURLRequest {
  original_url?: string;
  expires_at?: string;
  is_active?: boolean;
}

export interface URLResponse {
  success: boolean;
  message: string;
  data: URL;
}

export interface URLsResponse {
  success: boolean;
  message: string;
  data: URL[];
  pagination: {
    page: number;
    limit: number;
    total: number;
    total_pages: number;
  };
}

export interface URLStats {
  url_response: URL;
  recent_clicks?: ClickEvent[];
}

export interface ClickEvent {
  id: number;
  url_id: number;
  ip_address: string;
  user_agent: string;
  referer?: string;
  clicked_at: string;
}

export interface URLStatsResponse {
  success: boolean;
  message: string;
  data: URLStats;
}
