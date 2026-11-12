<template>
	<div class="play">
		<header class="play-header">
			<RouterLink to="/dashboard" class="brand" aria-label="LingoDash">
				<img :src="logo" alt="" class="brand-logo" />
				<span class="brand-name">LingoDash</span>
			</RouterLink>
			<div class="actions">
				<span v-if="voice.status.value !== 'off'" class="voice-status" :class="voice.status.value" role="status">
					<BaseIcon :name="voice.status.value === 'mic-denied' || voice.muted.value ? 'mic-off' : 'mic'" :size="16" />
					<span class="btn-text">{{ t(VOICE_STATUS_KEYS[voice.status.value]) }}</span>
				</span>
				<BaseButton
					v-if="voice.status.value === 'connected' || voice.status.value === 'waiting'"
					variant="ghost"
					pill
					:aria-pressed="voice.muted.value"
					@click="voice.toggleMute"
				>
					<BaseIcon :name="voice.muted.value ? 'mic-off' : 'mic'" :size="16" />
					<span class="btn-text">{{ voice.muted.value ? t('play.voice.unmute') : t('play.voice.mute') }}</span>
				</BaseButton>
				<BaseButton variant="ghost" pill @click="enterFullscreen"><BaseIcon name="maximize" :size="16" /><span class="btn-text">{{ t('play.fullscreen') }}</span></BaseButton>
				<BaseButton to="/dashboard" variant="outline" pill><BaseIcon name="arrow-left" :size="16" /><span class="btn-text">{{ t('play.back') }}</span></BaseButton>
			</div>
		</header>

		<main id="main" class="stage" tabindex="-1">
			<h1 class="visually-hidden">{{ t('play.heading') }}</h1>
			<iframe
				ref="gameFrame"
				src="/game.html"
				:title="t('play.gameTitle')"
				allow="autoplay; fullscreen; microphone"
				@load="loaded = true"
			></iframe>
			<div v-if="!loaded" class="loading" role="status">
				<BrandBolt />
				<span>{{ t('play.loading') }}</span>
			</div>
		</main>
	</div>
</template>

<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref, watch } from 'vue';
import logo from '../../images/logo-transparent.png';
import { useHeartbeat } from '../../composables/useHeartbeat';
import { useVoiceChat, type VoiceStatus } from '../../composables/useVoiceChat';
import { getToken, hasValidSession } from '../../auth/session';
import BaseButton from '../base/BaseButton.vue';
import BaseIcon from '../base/BaseIcon.vue';
import { t } from '../../i18n';
import BrandBolt from '../base/BrandBolt.vue';

useHeartbeat();

const gameFrame = ref<HTMLIFrameElement | null>(null);
const loaded = ref(false);
const voice = useVoiceChat();

const VOICE_STATUS_KEYS = {
	off: 'play.voice.off',
	'requesting-mic': 'play.voice.requestingMic',
	'mic-denied': 'play.voice.micDenied',
	waiting: 'play.voice.waiting',
	connected: 'play.voice.connected',
	error: 'play.voice.error',
} as const satisfies Record<VoiceStatus, string>;

// Messages to the game go only to our own origin.
function postToGame(message: Record<string, unknown>) {
	gameFrame.value?.contentWindow?.postMessage(JSON.stringify(message), window.location.origin);
}

// The game shows the voice state in its lobby; the backend blocks "ready" until it is 'connected'.
function sendVoiceStatus() {
	postToGame({ type: 'lingodash:voice-status', status: voice.status.value, muted: voice.muted.value });
}
watch([voice.status, voice.muted], sendVoiceStatus);

// The Godot game asks for the login token once it has started. Only our own iframe, on our own
// origin, gets an answer, and the reply is addressed to that origin only. The game gets the signed
// session token (never the password) and the backend checks it on every request.
function onGameMessage(event: MessageEvent) {
	const frameWindow = gameFrame.value?.contentWindow;
	if (!frameWindow || event.source !== frameWindow || event.origin !== window.location.origin) return;
	if (typeof event.data !== 'string') return;

	let message: { type?: unknown; roomCode?: unknown };
	try {
		message = JSON.parse(event.data);
	} catch {
		return;
	}

	switch (message?.type) {
		case 'lingodash:auth-request':
			postToGame({ type: 'lingodash:auth', token: hasValidSession() ? getToken() : null });
			break;
		// Sent by the game when it enters a lobby, and when it goes back to the lobby list.
		case 'lingodash:voice-join':
			if (typeof message.roomCode === 'string') voice.join(message.roomCode);
			sendVoiceStatus();
			break;
		case 'lingodash:voice-leave':
			voice.leave();
			break;
	}
}

onMounted(() => window.addEventListener('message', onGameMessage));
onBeforeUnmount(() => window.removeEventListener('message', onGameMessage));

function enterFullscreen() {
	gameFrame.value?.requestFullscreen();
}
</script>

<style scoped>
	.play {
		display: flex;
		flex-direction: column;
		height: 100vh;
		height: 100dvh;
		background: var(--color-surface-inverse);
		font-family: var(--font-display);
	}

	.play-header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 16px;
		padding: 12px 32px;
		background: var(--color-surface);
		border-bottom: 1px solid var(--color-border);
	}

	.brand {
		display: flex;
		align-items: center;
		gap: 10px;
		color: var(--color-on-surface);
		text-decoration: none;
	}

	.brand-logo {
		width: 34px;
		height: 34px;
		object-fit: contain;
	}

	.brand-name {
		font-weight: 800;
		font-size: var(--text-lg);
	}

	.actions {
		display: flex;
		gap: 10px;
	}

	/* The game is designed for 16:9. The frame keeps that shape and is as large as fits, so Godot
	   never adds black bars; leftover space shows the page background instead. */
	.voice-status {
		display: inline-flex;
		align-items: center;
		gap: 6px;
		padding: 6px 12px;
		border-radius: var(--radius-pill);
		font-size: var(--text-sm);
		font-weight: 700;
		color: var(--color-on-surface);
		background: var(--color-surface);
		border: 1px solid var(--color-border);
	}

	.voice-status.connected {
		color: var(--color-on-primary);
		background: var(--color-success);
		border-color: transparent;
	}

	.voice-status.mic-denied,
	.voice-status.error {
		color: var(--color-on-primary);
		background: var(--color-tertiary);
		border-color: transparent;
	}

	.stage {
		position: relative;
		flex: 1;
		min-height: 0;
		display: grid;
		place-items: center;
		container-type: size;
		background: var(--color-surface);
	}

	iframe {
		display: block;
		width: min(100cqw, 100cqh * 16 / 9);
		aspect-ratio: 16 / 9;
		border: 0;
		background: #000;
	}

	/* In fullscreen the frame takes the whole screen again. */
	iframe:fullscreen {
		width: 100%;
		height: 100%;
	}

	.loading {
		position: absolute;
		inset: 0;
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		gap: 14px;
		color: var(--color-on-primary);
		font-size: var(--text-md);
		font-weight: 600;
		background: var(--color-surface-inverse);
	}

	@media (max-width: 640px) {
		.play-header {
			padding: 10px 16px;
		}

		.brand-name {
			display: none;
		}
	}

	/* Small phones: icon-only buttons; the text stays for screen readers. */
	@media (max-width: 480px) {
		.actions {
			gap: 6px;
		}

		.actions .btn {
			padding: 0 12px;
		}

		.btn-text {
			position: absolute;
			width: 1px;
			height: 1px;
			overflow: hidden;
			clip: rect(0 0 0 0);
			white-space: nowrap;
		}
	}
</style>
