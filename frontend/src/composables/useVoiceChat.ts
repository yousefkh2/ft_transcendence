import { onBeforeUnmount, ref } from 'vue';
import { apiUrl } from '../api/client';
import { getToken } from '../auth/session';

// Voice chat between the players of a lobby. The audio goes directly from browser to browser
// (WebRTC); the backend (/ws/voice, backend/internal/hub/voice.go) only passes the connection
// setup messages along and remembers whose microphone works, because only players with a
// connected microphone may ready up.
//
// It runs on the website instead of inside Godot: a Godot web export cannot send microphone
// audio, and here the call keeps going while the game switches from the lobby to the match.

export type VoiceStatus =
	| 'off' // not in a lobby
	| 'requesting-mic' // waiting for the browser's microphone permission
	| 'mic-denied' // no microphone or permission refused
	| 'waiting' // microphone ready, not connected to the other player yet
	| 'connected' // talking to the other player
	| 'error'; // voice server unreachable, retrying

const RECONNECT_MS = 2000;

// Optional STUN servers for players on different networks, e.g. VITE_STUN_URLS=stun:stun.example.org:3478.
// Without them, players on the same machine or network still connect.
const iceServers: RTCIceServer[] = (import.meta.env.VITE_STUN_URLS ?? '')
	.split(',')
	.map((url: string) => url.trim())
	.filter(Boolean)
	.map((urls: string) => ({ urls }));

interface Peer {
	connection: RTCPeerConnection;
	audio: HTMLAudioElement;
	// ICE candidates that arrived before the offer/answer; they can only be added after it.
	pendingCandidates: RTCIceCandidateInit[];
}

interface SignalMessage {
	type: string;
	message?: string;
	peers?: string[];
	fromUser?: string;
	sdp?: string;
	candidate?: string;
}

export function useVoiceChat() {
	const status = ref<VoiceStatus>('off');
	const muted = ref(false);

	let roomCode = '';
	let socket: WebSocket | null = null;
	let microphone: MediaStream | null = null;
	let reconnectTimer: ReturnType<typeof setTimeout> | undefined;
	let reportedConnected: boolean | null = null;
	const peers = new Map<string, Peer>();

	function send(message: Record<string, unknown>) {
		if (socket?.readyState === WebSocket.OPEN) socket.send(JSON.stringify(message));
	}

	// Connected means: our microphone is live and at least one other player hears us.
	function updateStatus() {
		if (!roomCode || !microphone) return;
		const micLive = microphone.getAudioTracks().some((track) => track.readyState === 'live');
		const anyPeerConnected = [...peers.values()].some((p) => p.connection.connectionState === 'connected');
		const connected = micLive && anyPeerConnected;
		if (!micLive) status.value = 'mic-denied';
		else if (status.value !== 'error') status.value = connected ? 'connected' : 'waiting';

		if (connected !== reportedConnected) {
			reportedConnected = connected;
			send({ type: 'voice.status', connected });
		}
	}

	function createPeer(userId: string): Peer {
		closePeer(userId);
		const connection = new RTCPeerConnection({ iceServers });
		const audio = new Audio();
		audio.autoplay = true;
		const peer: Peer = { connection, audio, pendingCandidates: [] };
		peers.set(userId, peer);

		microphone?.getTracks().forEach((track) => connection.addTrack(track, microphone!));
		connection.ontrack = (event) => {
			audio.srcObject = event.streams[0] ?? new MediaStream([event.track]);
			audio.play().catch(() => {});
		};
		connection.onicecandidate = (event) => {
			if (event.candidate) {
				send({ type: 'webrtc.ice_candidate', targetUser: userId, candidate: JSON.stringify(event.candidate) });
			}
		};
		connection.onconnectionstatechange = updateStatus;
		return peer;
	}

	function closePeer(userId: string) {
		const peer = peers.get(userId);
		if (!peer) return;
		peer.connection.close();
		peer.audio.srcObject = null;
		peers.delete(userId);
	}

	async function flushCandidates(peer: Peer) {
		for (const candidate of peer.pendingCandidates.splice(0)) {
			await peer.connection.addIceCandidate(candidate).catch(() => {});
		}
	}

	// The player who joins later sends the offers, so two players never offer to each other at once.
	async function offerTo(userId: string) {
		const peer = createPeer(userId);
		const offer = await peer.connection.createOffer();
		await peer.connection.setLocalDescription(offer);
		send({ type: 'webrtc.offer', targetUser: userId, sdp: offer.sdp });
	}

	async function handleMessage(message: SignalMessage) {
		const from = message.fromUser ?? '';
		switch (message.type) {
			case 'voice.joined':
				for (const userId of message.peers ?? []) await offerTo(userId);
				updateStatus();
				break;
			case 'voice.peer_left':
				closePeer(from);
				updateStatus();
				break;
			case 'webrtc.offer': {
				const peer = createPeer(from);
				await peer.connection.setRemoteDescription({ type: 'offer', sdp: message.sdp });
				await flushCandidates(peer);
				const answer = await peer.connection.createAnswer();
				await peer.connection.setLocalDescription(answer);
				send({ type: 'webrtc.answer', targetUser: from, sdp: answer.sdp });
				break;
			}
			case 'webrtc.answer': {
				const peer = peers.get(from);
				if (!peer) return;
				await peer.connection.setRemoteDescription({ type: 'answer', sdp: message.sdp });
				await flushCandidates(peer);
				break;
			}
			case 'webrtc.ice_candidate': {
				const peer = peers.get(from);
				if (!peer || !message.candidate) return;
				const candidate = JSON.parse(message.candidate) as RTCIceCandidateInit;
				if (peer.connection.remoteDescription) await peer.connection.addIceCandidate(candidate).catch(() => {});
				else peer.pendingCandidates.push(candidate);
				break;
			}
			case 'error':
				console.warn('Voice chat:', message.message);
				break;
		}
	}

	function openSocket() {
		const url = `${apiUrl.replace(/^http/, 'ws')}/ws/voice`;
		const ws = new WebSocket(url);
		socket = ws;
		reportedConnected = null;

		ws.onopen = () => {
			if (status.value === 'error') status.value = 'waiting';
			ws.send(JSON.stringify({ type: 'voice.join', roomCode, token: getToken() }));
		};
		ws.onmessage = (event) => {
			let message: SignalMessage;
			try {
				message = JSON.parse(event.data);
			} catch {
				return;
			}
			handleMessage(message).catch((error) => console.warn('Voice chat signal failed:', error));
		};
		ws.onclose = () => {
			if (socket !== ws) return;
			socket = null;
			// Still in the lobby: the server or network dropped, so try again.
			if (roomCode) {
				for (const userId of [...peers.keys()]) closePeer(userId);
				status.value = 'error';
				reconnectTimer = setTimeout(openSocket, RECONNECT_MS);
			}
		};
	}

	async function join(code: string) {
		code = code.trim().toUpperCase();
		// Same lobby again: only retry when the microphone was refused before.
		if (!code || (code === roomCode && status.value !== 'mic-denied')) return;
		leave();
		roomCode = code;

		status.value = 'requesting-mic';
		try {
			microphone = await navigator.mediaDevices.getUserMedia({
				audio: { echoCancellation: true, noiseSuppression: true, autoGainControl: true },
			});
		} catch {
			if (roomCode === code) status.value = 'mic-denied';
			return;
		}
		if (roomCode !== code) {
			// Left again while the permission prompt was open.
			microphone.getTracks().forEach((track) => track.stop());
			microphone = null;
			return;
		}
		microphone.getAudioTracks().forEach((track) => {
			track.enabled = !muted.value;
			track.onended = updateStatus; // microphone unplugged
		});
		status.value = 'waiting';
		openSocket();
	}

	function leave() {
		roomCode = '';
		clearTimeout(reconnectTimer);
		const ws = socket;
		socket = null;
		ws?.close();
		for (const userId of [...peers.keys()]) closePeer(userId);
		microphone?.getTracks().forEach((track) => track.stop());
		microphone = null;
		status.value = 'off';
	}

	// Muting keeps the connection (and the ready state); the other player just hears nothing.
	function toggleMute() {
		muted.value = !muted.value;
		microphone?.getAudioTracks().forEach((track) => (track.enabled = !muted.value));
	}

	onBeforeUnmount(leave);

	return { status, muted, join, leave, toggleMute };
}
