import { request } from './client';

export interface Friend {
	friendId: string;
	friendUsername: string;
	createdAt: string;
	online: boolean;
	lastSeenAt?: string;
}

export interface FriendRequest {
	id: string;
	requesterId: string;
	requesterUsername: string;
	createdAt: string;
}

export function fetchFriends(): Promise<Friend[]> {
	return request<Friend[]>('/api/friends');
}

export function fetchFriendRequests(): Promise<FriendRequest[]> {
	return request<FriendRequest[]>('/api/friends/requests');
}

export function sendFriendRequest(username: string): Promise<{ id: string }> {
	return request<{ id: string }>('/api/friends/requests', { method: 'POST', body: { username } });
}

export function acceptFriendRequest(id: string): Promise<void> {
	return request<void>(`/api/friends/requests/${id}/accept`, { method: 'POST' });
}

export function declineFriendRequest(id: string): Promise<void> {
	return request<void>(`/api/friends/requests/${id}/decline`, { method: 'POST' });
}

export function removeFriend(friendId: string): Promise<void> {
	return request<void>(`/api/friends/${encodeURIComponent(friendId)}`, { method: 'DELETE' });
}
