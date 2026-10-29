<template>
	<main id="main" class="not-found" tabindex="-1">
		<BrandBolt class="bolt-mark" />
		<span class="code">404</span>
		<h1>{{ t('notFound.title') }}</h1>
		<p>{{ t('notFound.text') }}</p>
		<BaseButton to="/" class="home-btn">{{ t('notFound.home') }}</BaseButton>
	</main>
</template>

<script setup lang="ts">
import BaseButton from './base/BaseButton.vue';
import BrandBolt from './base/BrandBolt.vue';
import { t } from '../i18n';
</script>

<style scoped>
	.not-found {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		gap: 12px;
		min-height: 100vh;
		padding: 32px 16px;
		text-align: center;
		background: var(--color-surface);
		color: var(--color-on-surface);
		font-family: var(--font-display);
	}

	.bolt-mark {
		margin-bottom: 8px;
	}

	.code {
		font-family: var(--font-mono);
		font-size: var(--text-sm);
		font-weight: 700;
		letter-spacing: 0.1em;
		color: var(--color-primary);
	}

	h1 {
		margin: 0;
		font-size: 28px;
		line-height: 0.95;
		font-weight: 800;
	}

	p {
		margin: 0;
		color: var(--color-on-surface-variant);
	}

	.home-btn {
		margin-top: 16px;
	}
</style>
