import type { MatchHistoryEntry } from '../api/users';
import { locale } from '../i18n';

const DAY_MS = 24 * 60 * 60 * 1000;

export function formatLabel(value: string) {
	return value.replace(/_/g, ' ').replace(/\b\w/g, (c) => c.toUpperCase());
}

// "Today", "Yesterday", "3 days ago" in the current language.
export function formatWhen(iso: string) {
	const diffDays = Math.max(0, Math.floor((Date.now() - new Date(iso).getTime()) / DAY_MS));
	const text = new Intl.RelativeTimeFormat(locale.value, { numeric: 'auto' }).format(-diffDays, 'day');
	return text.charAt(0).toLocaleUpperCase(locale.value) + text.slice(1);
}

export function roundResult(match: MatchHistoryEntry): 'won' | 'lost' | 'pending' {
	if (match.isWinner === true) return 'won';
	if (match.isWinner === false) return 'lost';
	return 'pending';
}

// Time between start and end as "m:ss", or null while a round has no end time.
export function formatDuration(startedAt: string, endedAt: string | null) {
	if (!endedAt) return null;
	const seconds = Math.max(0, Math.round((new Date(endedAt).getTime() - new Date(startedAt).getTime()) / 1000));
	return `${Math.floor(seconds / 60)}:${String(seconds % 60).padStart(2, '0')}`;
}

export function countSince(matches: MatchHistoryEntry[], days: number) {
	const since = Date.now() - days * DAY_MS;
	return matches.filter((m) => new Date(m.startedAt).getTime() >= since).length;
}
