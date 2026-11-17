import { onMounted, onUnmounted, ref } from 'vue';
import { apiUrl } from '../api/client';
import type { ChatMessage } from '../api/chat';
import { getToken } from '../auth/session';

const RETRY_MIN_MS = 1000;
const RETRY_MAX_MS = 15000;

// Live chat messages from /ws/chat. Sending goes through the REST API; this socket only receives.
// It reconnects on its own with a growing delay, and onReconnect lets the page reload what it missed.
export function useChatSocket(onMessage: (message: ChatMessage) => void, onReconnect?: () => void) {
	const connected = ref(false);
	let socket: WebSocket | null = null;
	let retryTimer: ReturnType<typeof setTimeout> | undefined;
	let retryDelay = RETRY_MIN_MS;
	let wasConnected = false;
	let stopped = false;

	function connect() {
		const token = getToken();
		if (stopped || !token) return;

		const ws = new WebSocket(`${apiUrl.replace(/^http/, 'ws')}/ws/chat`);
		socket = ws;
		ws.onopen = () => ws.send(JSON.stringify({ type: 'chat.auth', token }));
		ws.onmessage = (event) => {
			let data: { type?: string; payload?: ChatMessage };
			try {
				data = JSON.parse(event.data);
			} catch {
				return;
			}
			if (data.type === 'chat.ready') {
				connected.value = true;
				retryDelay = RETRY_MIN_MS;
				if (wasConnected) onReconnect?.();
				wasConnected = true;
			} else if (data.type === 'chat.message' && data.payload) {
				onMessage(data.payload);
			}
		};
		ws.onclose = () => {
			connected.value = false;
			if (socket === ws) socket = null;
			if (stopped) return;
			retryTimer = setTimeout(connect, retryDelay);
			retryDelay = Math.min(retryDelay * 2, RETRY_MAX_MS);
		};
	}

	onMounted(connect);
	onUnmounted(() => {
		stopped = true;
		clearTimeout(retryTimer);
		socket?.close();
	});

	return { connected };
}
