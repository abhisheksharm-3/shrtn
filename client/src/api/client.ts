import type { UrlType, UrlInputType, ApiErrorType, LinkPreviewType, RequestOptionsType } from '@/types/api';

const API_BASE_URL = import.meta.env.VITE_API_URL || '';
const API_KEY = import.meta.env.VITE_API_KEY || '';

class ApiError extends Error {
    code: string;

    constructor(message: string, code: string) {
        super(message);
        this.code = code;
        this.name = 'ApiError';
    }
}

async function request<T>(endpoint: string, options: RequestOptionsType = {}): Promise<T> {
    const headers: Record<string, string> = {
        'Content-Type': 'application/json',
    };

    if (API_KEY) {
        headers['X-API-Key'] = API_KEY;
    }

    const response = await fetch(`${API_BASE_URL}${endpoint}`, {
        method: options.method || 'GET',
        headers,
        body: options.body ? JSON.stringify(options.body) : undefined,
    });

    if (!response.ok) {
        const errorData = await response.json() as ApiErrorType;
        throw new ApiError(errorData.error, errorData.code);
    }

    return response.json() as Promise<T>;
}

/** Shorten a URL */
export async function shortenUrl(input: UrlInputType): Promise<UrlType> {
    return request<UrlType>('/api/shorten', {
        method: 'POST',
        body: input,
    });
}

/** Get URL info by short code */
export async function getUrlByCode(shortCode: string): Promise<UrlType> {
    return request<UrlType>(`/api/${shortCode}`);
}

/** Get redirect URL for a short code */
export function getRedirectUrl(shortCode: string): string {
    return `${API_BASE_URL}/${shortCode}`;
}

/** Fetch link preview metadata */
export async function fetchLinkPreview(url: string): Promise<LinkPreviewType> {
    const encodedUrl = encodeURIComponent(url);
    return request<LinkPreviewType>(`/api/preview?url=${encodedUrl}`);
}

export { ApiError };
