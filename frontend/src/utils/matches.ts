import type { MatchHistoryEntry } from '../api/users';

const DAY_MS = 24 * 60 * 60 * 1000;

export function formatLabel(value: string) {
	return value.replace(/_/g, ' ').replace(/\b\w/g, (c) => c.toUpperCase());
}

export function formatWhen(iso: string) {
	const diffDays = Math.floor((Date.now() - new Date(iso).getTime()) / DAY_MS);
	if (diffDays <= 0) return 'Today';
	if (diffDays === 1) return 'Yesterday';
	return `${diffDays} days ago`;
}

export function roundResult(match: MatchHistoryEntry): 'won' | 'lost' | 'pending' {
	if (match.isWinner === true) return 'won';
	if (match.isWinner === false) return 'lost';
	return 'pending';
}

export function countSince(matches: MatchHistoryEntry[], days: number) {
	const since = Date.now() - days * DAY_MS;
	return matches.filter((m) => new Date(m.startedAt).getTime() >= since).length;
}
