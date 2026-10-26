import { request } from './client';

export interface Profile {
	id: string;
	username: string;
	email: string;
	xp: number;
	createdAt: string;
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
