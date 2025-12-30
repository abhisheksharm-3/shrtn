/** URL entity returned from the API */
export type UrlType = {
  id: string;
  shortCode: string;
  originalUrl: string;
  createdAt: string;
  updatedAt: string;
  clicks: number;
  userId?: string;
};

/** Input for creating a shortened URL */
export type UrlInputType = {
  originalUrl: string;
  customCode?: string;
};

/** Paginated list of URLs */
export type UrlListResponseType = {
  urls: UrlType[];
  total: number;
  limit: number;
  offset: number;
};

/** Link preview metadata from Open Graph */
export type LinkPreviewType = {
  url: string;
  title: string;
  description: string;
  image: string;
  siteName: string;
  favicon: string;
};

/** API error response */
export type ApiErrorType = {
  error: string;
  code: string;
};

export type RequestOptionsType = {
    method?: 'GET' | 'POST' | 'DELETE';
    body?: unknown;
};

export type ShortenerStateType = {
    url: UrlType | null;
    error: string | null;
    isPending: boolean;
};