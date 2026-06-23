/**
 * HTTP client for the Call Booking API.
 *
 * Responsibilities (kept out of components):
 *  - resolve the base URL from the environment (VITE_API_BASE_URL);
 *  - perform requests and parse JSON;
 *  - translate HTTP/error responses into a typed `ApiError`.
 *
 * Components must use the resource modules (eventTypes, slots, bookings, ...),
 * never call `fetch` or build URLs directly.
 */
import type { ApiErrorBody, ApiErrorCode } from './types';

/**
 * Base URL of the API. Sourced from the environment so the same build can
 * target the Prism mock or a real backend. Falls back to the Prism mock's
 * default port for a zero-config dev experience.
 */
const BASE_URL = (
  import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:4010'
).replace(/\/+$/, '');

/** A failed API call: transport error, or a non-2xx HTTP response. */
export class ApiError extends Error {
  /** HTTP status code, or 0 for a transport/network error. */
  readonly status: number;
  /** Machine-readable error code from the contract, when present. */
  readonly code?: ApiErrorCode;
  /** Optional field-level validation details from the contract. */
  readonly details?: string[];

  constructor(
    status: number,
    message: string,
    code?: ApiErrorCode,
    details?: string[],
  ) {
    super(message);
    this.name = 'ApiError';
    this.status = status;
    this.code = code;
    this.details = details;
  }
}

function isErrorBody(value: unknown): value is ApiErrorBody {
  return (
    typeof value === 'object' &&
    value !== null &&
    typeof (value as Record<string, unknown>).message === 'string'
  );
}

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  let response: Response;
  try {
    response = await fetch(`${BASE_URL}${path}`, {
      ...init,
      headers: {
        Accept: 'application/json',
        ...(init?.body ? { 'Content-Type': 'application/json' } : {}),
        ...init?.headers,
      },
    });
  } catch {
    throw new ApiError(
      0,
      'Could not reach the API. Is the server (or Prism mock) running?',
    );
  }

  const hasJson = response.headers
    .get('content-type')
    ?.includes('application/json');
  const body: unknown = hasJson
    ? await response.json().catch(() => undefined)
    : undefined;

  if (!response.ok) {
    if (isErrorBody(body)) {
      throw new ApiError(response.status, body.message, body.code, body.details);
    }
    throw new ApiError(
      response.status,
      `Request failed with status ${response.status}.`,
    );
  }

  return body as T;
}

export const apiClient = {
  get: <T>(path: string): Promise<T> => request<T>(path),
  post: <T>(path: string, data: unknown): Promise<T> =>
    request<T>(path, { method: 'POST', body: JSON.stringify(data) }),
};

/** Base URL the client is currently using (handy for debugging / display). */
export const apiBaseUrl = BASE_URL;
