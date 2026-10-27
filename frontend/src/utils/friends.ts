import type { Friend } from '../api/friends';

// Online friends first, then the ones seen most recently, so you see who you can play with right now.
export function sortByPresence(friends: Friend[]): Friend[] {
	const lastSeen = (f: Friend) => (f.lastSeenAt ? new Date(f.lastSeenAt).getTime() : 0);
	return [...friends].sort((a, b) => Number(b.online) - Number(a.online) || lastSeen(b) - lastSeen(a));
}
