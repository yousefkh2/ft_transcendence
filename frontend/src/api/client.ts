import { logout } from '../auth/logout';
import { getToken } from '../auth/session';

const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080';

export class ApiError extends Error {
	constructor(
		message: string,
		readonly status: number,
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
	if (body !== undefined) headers['Content-Type'] = 'application/json';

	const response = await fetch(`${apiUrl}${path}`, {
		method,
		headers,
		body: body === undefined ? undefined : JSON.stringify(body),
	});

	if (!response.ok) {
		// Only an authenticated request can mean "session no longer valid"; a failed login is also a 401.
		if (response.status === 401 && token) {
			await logout();
		}
		const data = await response.json().catch(() => null);
		throw new ApiError(
			typeof data?.message === 'string' ? data.message : `Request failed (${response.status})`,
			response.status,
		);
	}

	if (response.status === 204) {
		return undefined as T;
	}
	return response.json();
}
