export function formatDate(iso: string) {
	return new Date(iso).toLocaleDateString(undefined, { month: 'short', day: 'numeric', year: 'numeric' });
}

export function formatLastSeen(online: boolean, lastSeenAt?: string) {
	if (online) return 'Online now';
	if (!lastSeenAt) return 'Offline';
	const minutes = Math.floor((Date.now() - new Date(lastSeenAt).getTime()) / 60000);
	if (minutes < 1) return 'Last seen just now';
	if (minutes < 60) return `Last seen ${minutes} min ago`;
	const hours = Math.floor(minutes / 60);
	if (hours < 24) return `Last seen ${hours} h ago`;
	const days = Math.floor(hours / 24);
	return `Last seen ${days} ${days === 1 ? 'day' : 'days'} ago`;
}
