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

export function countSince(matches: MatchHistoryEntry[], days: number) {
	const since = Date.now() - days * DAY_MS;
	return matches.filter((m) => new Date(m.startedAt).getTime() >= since).length;
}
