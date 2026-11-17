import { request } from './client';
import type { LevelInfo } from './users';

// Ids from the backend (internal/game/achievements.go). Each has a translated title and text.
export const ACHIEVEMENT_IDS = [
	'first_round',
	'first_win',
	'first_loss',
	'quick_solve',
	'lightning',
	'flawless',
	'streak_3',
	'rounds_10',
	'wins_10',
	'level_5',
	'level_10',
	'polyglot',
	'team_player',
] as const;
export type AchievementId = (typeof ACHIEVEMENT_IDS)[number];

export interface Achievement {
	id: AchievementId;
	progress: number;
	target: number;
	unlocked: boolean;
	unlockedAt?: string;
	// Unlocked but not shown to the user yet.
	new: boolean;
}

export function fetchMyAchievements(): Promise<Achievement[]> {
	return request<Achievement[]>('/api/users/me/achievements');
}

export function fetchUserAchievements(username: string): Promise<Achievement[]> {
	return request<Achievement[]>(`/api/users/${encodeURIComponent(username)}/achievements`);
}

export function markAchievementsSeen(): Promise<void> {
	return request<void>('/api/users/me/achievements/seen', { method: 'POST' });
}

export interface LeaderboardEntry {
	rank: number;
	userId: string;
	username: string;
	avatarUrl?: string;
	xp: number;
	level: LevelInfo;
	wins: number;
	rounds: number;
	achievements: number;
	isMe: boolean;
}

export interface Leaderboard {
	entries: LeaderboardEntry[];
	me?: LeaderboardEntry;
}

export type LeaderboardScope = 'global' | 'friends';

export function fetchLeaderboard(scope: LeaderboardScope): Promise<Leaderboard> {
	return request<Leaderboard>(`/api/leaderboard?scope=${scope}`);
}
