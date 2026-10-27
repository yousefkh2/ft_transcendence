import { sendHeartbeat } from '../api/users';
import { usePolling } from './usePolling';

// How often the app reports itself as open and refreshes the online status of friends.
// Must stay below OnlineWindow (90 s) in backend/internal/db/presence.go.
export const PRESENCE_INTERVAL_MS = 30_000;

// Tells the backend regularly that the app is open, so friends see the user as online.
export function useHeartbeat() {
	usePolling(() => sendHeartbeat().catch(() => {}), PRESENCE_INTERVAL_MS);
}
