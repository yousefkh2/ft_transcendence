import { logout } from '../auth/logout';
import { getToken } from '../auth/session';
import { t } from '../i18n';
import { translateApiError } from '../i18n/apiErrors';

const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080';

// Turns a path the API returns (like an avatarUrl) into a full URL for <img src>.
export function apiAssetUrl(path?: string): string | undefined {
	return path ? `${apiUrl}${path}` : undefined;
}

// message is translated for the user; serverMessage is the backend's original English text, so code
// can still tell errors apart (status 0 means the server could not be reached at all).
export class ApiError extends Error {
	constructor(
		message: string,
		readonly status: number,
		readonly serverMessage = '',
	) {
		super(message);
	}
}

interface RequestOptions {
	method?: string;
	body?: unknown;
	auth?: boolean;
}

export async function request<T>(path: string, { method = 'GET', body, auth = true }: RequestOptions = {}): Promise<T> {
	const headers: Record<string, string> = {};
	const token = auth ? getToken() : null;
	if (token) headers.Authorization = `Bearer ${token}`;
	// FormData (file uploads) sets its own multipart Content-Type with the boundary.
	const isForm = body instanceof FormData;
	if (body !== undefined && !isForm) headers['Content-Type'] = 'application/json';

	let response: Response;
	try {
		response = await fetch(`${apiUrl}${path}`, {
			method,
			headers,
			body: body === undefined ? undefined : isForm ? body : JSON.stringify(body),
		});
	} catch {
		throw new ApiError(t('apiErrors.network'), 0);
	}

	if (!response.ok) {
		// Only an authenticated request can mean "session no longer valid"; a failed login is also a 401.
		if (response.status === 401 && token) {
			await logout();
		}
		const data = await response.json().catch(() => null);
		const serverMessage = typeof data?.message === 'string' ? data.message : '';
		throw new ApiError(translateApiError(serverMessage), response.status, serverMessage);
	}

	if (response.status === 204) {
		return undefined as T;
	}
	return response.json();
}
