const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080';

export interface Profile {
	id: string;
	username: string;
	email: string;
	xp: number;
}

export interface MatchHistoryEntry {
	sessionId: string;
	gameMode: string;
	status: string;
	startedAt: string;
	endedAt: string | null;
	role: string;
	isWinner: boolean | null;
}

export interface Friend {
	friendId: string;
	friendUsername: string;
	createdAt: string;
}

export interface FriendRequest {
	id: string;
	requesterId: string;
	requesterUsername: string;
	createdAt: string;
}

class ApiError extends Error {}

function authHeaders(): HeadersInit {
	const token = localStorage.getItem('lingodash_token');
	return token ? { Authorization: `Bearer ${token}` } : {};
}

async function request<T>(path: string, options: RequestInit = {}): Promise<T> {
	const response = await fetch(`${apiUrl}${path}`, {
		...options,
		headers: { ...authHeaders(), ...(options.headers || {}) },
	});

	if (!response.ok) {
		const body = await response.json().catch(() => null);
		throw new ApiError(typeof body?.message === 'string' ? body.message : `Request failed: ${response.status}`);
	}

	if (response.status === 204) {
		return undefined as T;
	}

	return response.json();
}

export function fetchProfile(): Promise<Profile> {
	return request<Profile>('/api/users/me');
}

export function fetchMatchHistory(): Promise<MatchHistoryEntry[]> {
	return request<MatchHistoryEntry[]>('/api/users/me/matches');
}

export function fetchFriends(): Promise<Friend[]> {
	return request<Friend[]>('/api/friends');
}

export function fetchFriendRequests(): Promise<FriendRequest[]> {
	return request<FriendRequest[]>('/api/friends/requests');
}

export function sendFriendRequest(username: string): Promise<{ id: string }> {
	return request<{ id: string }>('/api/friends/requests', {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ username }),
	});
}

export function acceptFriendRequest(id: string): Promise<void> {
	return request<void>(`/api/friends/requests/${id}/accept`, { method: 'POST' });
}

export function declineFriendRequest(id: string): Promise<void> {
	return request<void>(`/api/friends/requests/${id}/decline`, { method: 'POST' });
}
