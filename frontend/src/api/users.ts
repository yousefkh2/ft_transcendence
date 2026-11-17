import { request } from './client';

// Level and progress towards the next one, computed by the backend from the XP.
export interface LevelInfo {
	level: number;
	xpIntoLevel: number;
	xpForNext: number;
}

export interface Profile {
	id: string;
	username: string;
	email: string;
	xp: number;
	createdAt: string;
	learningLanguage: string;
	avatarUrl?: string;
	level: LevelInfo;
}

export interface MatchHistoryEntry {
	sessionId: string;
	gameMode: string;
	status: string;
	startedAt: string;
	endedAt: string | null;
	role?: string;
	isWinner: boolean | null;
	lang: string;
	// Missing for rounds saved before the backend recorded it.
	endReason?: 'solved' | 'time_up' | 'out_of_checks';
	// XP this player got for the round (0 for rounds saved before XP was awarded).
	xpEarned: number;
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
	level: LevelInfo;
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
