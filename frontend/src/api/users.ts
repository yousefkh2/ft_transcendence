import { request } from './client';

export interface Profile {
	id: string;
	username: string;
	email: string;
	xp: number;
	createdAt: string;
	learningLanguage: string;
	avatarUrl?: string;
}

export interface MatchHistoryEntry {
	sessionId: string;
	gameMode: string;
	status: string;
	startedAt: string;
	endedAt: string | null;
	role?: string;
	isWinner: boolean | null;
}

export function fetchProfile(): Promise<Profile> {
	return request<Profile>('/api/users/me');
}

export function fetchMatchHistory(): Promise<MatchHistoryEntry[]> {
	return request<MatchHistoryEntry[]>('/api/users/me/matches');
}

export interface PublicProfile {
	username: string;
	xp: number;
	createdAt: string;
	learningLanguage: string;
	avatarUrl?: string;
	online: boolean;
	lastSeenAt?: string;
	roundsPlayed: number;
	wins: number;
	losses: number;
}

export function fetchPublicProfile(username: string): Promise<PublicProfile> {
	return request<PublicProfile>(`/api/users/${encodeURIComponent(username)}`);
}

export function deleteAccount(password: string): Promise<void> {
	return request<void>('/api/users/me', { method: 'DELETE', body: { password } });
}

export interface ProfileChanges {
	username?: string;
	email?: string;
	learningLanguage?: string;
	currentPassword?: string;
}

export function updateProfile(changes: ProfileChanges): Promise<Profile> {
	return request<Profile>('/api/users/me', { method: 'PATCH', body: changes });
}

export function changePassword(currentPassword: string, newPassword: string): Promise<void> {
	return request<void>('/api/users/me/password', { method: 'PATCH', body: { currentPassword, newPassword } });
}

export function sendHeartbeat(): Promise<void> {
	return request<void>('/api/users/me/heartbeat', { method: 'POST' });
}

export function uploadAvatar(picture: Blob): Promise<Profile> {
	const form = new FormData();
	form.append('avatar', picture, 'avatar.jpg');
	return request<Profile>('/api/users/me/avatar', { method: 'PUT', body: form });
}

export function deleteAvatar(): Promise<Profile> {
	return request<Profile>('/api/users/me/avatar', { method: 'DELETE' });
}
