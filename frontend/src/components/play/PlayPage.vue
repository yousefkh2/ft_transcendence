<template>
	<div class="play">
		<header class="play-header">
			<RouterLink to="/dashboard" class="brand" aria-label="LingoDash">
				<img :src="logo" alt="" class="brand-logo" />
				<span class="brand-name">LingoDash</span>
			</RouterLink>
			<div class="actions">
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
import { onBeforeUnmount, onMounted, ref } from 'vue';
import logo from '../../images/logo-transparent.png';
import { useHeartbeat } from '../../composables/useHeartbeat';
import { getToken, hasValidSession } from '../../auth/session';
import BaseButton from '../base/BaseButton.vue';
import BaseIcon from '../base/BaseIcon.vue';
import { t } from '../../i18n';
import BrandBolt from '../base/BrandBolt.vue';

useHeartbeat();

const gameFrame = ref<HTMLIFrameElement | null>(null);
const loaded = ref(false);

// The Godot game asks for the login token once it has started. Only our own iframe, on our own
// origin, gets an answer, and the reply is addressed to that origin only. The game gets the signed
// session token (never the password) and the backend checks it on every request.
function onGameMessage(event: MessageEvent) {
	const frameWindow = gameFrame.value?.contentWindow;
	if (!frameWindow || event.source !== frameWindow || event.origin !== window.location.origin) return;
	if (typeof event.data !== 'string') return;

	let message: { type?: unknown };
	try {
		message = JSON.parse(event.data);
	} catch {
		return;
	}
	if (message?.type !== 'lingodash:auth-request') return;

	const token = hasValidSession() ? getToken() : null;
	frameWindow.postMessage(JSON.stringify({ type: 'lingodash:auth', token }), window.location.origin);
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
