<template>
	<div class="play">
		<header class="play-header">
			<RouterLink to="/dashboard" class="brand">
				<img :src="logo" alt="" class="brand-logo" />
				<span class="brand-name">LingoDash</span>
			</RouterLink>
			<div class="actions">
				<BaseButton variant="ghost" pill @click="enterFullscreen">Fullscreen</BaseButton>
				<BaseButton to="/dashboard" variant="outline" pill>Back to dashboard</BaseButton>
			</div>
		</header>

		<main class="stage">
			<iframe
				ref="gameFrame"
				src="/game.html"
				title="LingoDash game"
				allow="autoplay; fullscreen; microphone"
				@load="loaded = true"
			></iframe>
			<div v-if="!loaded" class="loading" role="status">
				<BrandBolt />
				<span>Loading game...</span>
			</div>
		</main>
	</div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import logo from '../../images/logo-transparent.png';
import BaseButton from '../base/BaseButton.vue';
import BrandBolt from '../base/BrandBolt.vue';

const gameFrame = ref<HTMLIFrameElement | null>(null);
const loaded = ref(false);

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
		border-bottom: 1px solid rgba(118, 117, 134, 0.15);
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
		font-size: 17px;
	}

	.actions {
		display: flex;
		gap: 10px;
	}

	.stage {
		position: relative;
		flex: 1;
		min-height: 0;
		background: #000;
	}

	iframe {
		display: block;
		width: 100%;
		height: 100%;
		border: 0;
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
		font-size: 14px;
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

	@media (max-width: 400px) {
		.actions {
			gap: 6px;
		}

		.actions .btn {
			padding: 0 12px;
		}
	}
</style>
