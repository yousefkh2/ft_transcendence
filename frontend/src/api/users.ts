import { request } from './client';

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

export function fetchProfile(): Promise<Profile> {
	return request<Profile>('/api/users/me');
}

export function fetchMatchHistory(): Promise<MatchHistoryEntry[]> {
	return request<MatchHistoryEntry[]>('/api/users/me/matches');
}
