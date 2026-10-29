import { locale, t } from '../i18n';

export function formatDate(iso: string) {
	return new Date(iso).toLocaleDateString(locale.value, { month: 'short', day: 'numeric', year: 'numeric' });
}

export function formatNumber(value: number) {
	return value.toLocaleString(locale.value);
}

// "5 minutes ago", "vor 5 Minuten", ... in the current language.
function timeAgo(iso: string) {
	const rtf = new Intl.RelativeTimeFormat(locale.value, { numeric: 'auto' });
	const minutes = Math.floor((Date.now() - new Date(iso).getTime()) / 60000);
	if (minutes < 60) return rtf.format(-minutes, 'minute');
	const hours = Math.floor(minutes / 60);
	if (hours < 24) return rtf.format(-hours, 'hour');
	return rtf.format(-Math.floor(hours / 24), 'day');
}

export function formatLastSeen(online: boolean, lastSeenAt?: string) {
	if (online) return t('presence.online');
	if (!lastSeenAt) return t('presence.offline');
	return t('presence.lastSeen', { time: timeAgo(lastSeenAt) });
}

export function formatPercent(ratio: number) {
	return new Intl.NumberFormat(locale.value, { style: 'percent', maximumFractionDigits: 0 }).format(ratio);
}
