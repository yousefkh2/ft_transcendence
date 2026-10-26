<template>
	<div class="legal">
		<header class="legal-header">
			<RouterLink to="/" class="brand" aria-label="LingoDash home">
				<img :src="logo" alt="" class="brand-logo" />
				<span class="brand-name">LingoDash</span>
			</RouterLink>
			<BaseButton to="/" variant="ghost" pill>Back to LingoDash</BaseButton>
		</header>

		<main class="legal-main">
			<BaseEyebrow>Legal</BaseEyebrow>
			<h1>{{ title }}</h1>
			<p class="updated">Last updated: {{ updated }}</p>
			<div class="legal-content">
				<slot />
			</div>
		</main>

		<footer class="legal-footer">
			<nav class="legal-links" aria-label="Legal">
				<RouterLink to="/privacy">Privacy Policy</RouterLink>
				<RouterLink to="/terms">Terms of Service</RouterLink>
			</nav>
			<span class="copyright">© 2026 LingoDash. A student project of the 42 curriculum.</span>
		</footer>
	</div>
</template>

<script setup lang="ts">
import logo from '../../images/logo-transparent.png';
import BaseButton from '../base/BaseButton.vue';
import BaseEyebrow from '../base/BaseEyebrow.vue';

defineProps<{ title: string; updated: string }>();
</script>

<style scoped>
	.legal {
		min-height: 100vh;
		display: flex;
		flex-direction: column;
		background: var(--color-surface);
		color: var(--color-on-surface);
		font-family: var(--font-display);
	}

	.legal-header {
		position: sticky;
		top: 0;
		z-index: 10;
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 16px;
		padding: 14px 32px;
		background: rgba(249, 249, 255, 0.85);
		backdrop-filter: blur(12px);
		-webkit-backdrop-filter: blur(12px);
		border-bottom: 1px solid rgba(118, 117, 134, 0.15);
	}

	.brand {
		display: flex;
		align-items: center;
		gap: 10px;
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

	.legal-main {
		flex: 1;
		width: 100%;
		max-width: 760px;
		margin: 0 auto;
		padding: 48px 24px 64px;
	}

	h1 {
		margin: 0 0 8px;
		font-size: clamp(28px, 5vw, 38px);
		font-weight: 800;
		letter-spacing: -0.02em;
		line-height: 1.15;
	}

	.updated {
		margin: 0 0 32px;
		font-size: 13.5px;
		color: var(--color-outline);
	}

	.legal-content {
		font-size: 15px;
		line-height: 1.7;
		color: var(--color-on-surface-variant);
	}

	.legal-content :slotted(h2) {
		margin: 36px 0 10px;
		font-size: 19px;
		font-weight: 800;
		line-height: 1.3;
		color: var(--color-on-surface);
	}

	.legal-content :slotted(h3) {
		margin: 20px 0 6px;
		font-size: 15.5px;
		font-weight: 700;
		color: var(--color-on-surface);
	}

	.legal-content :slotted(p) {
		margin: 0 0 12px;
	}

	.legal-content :slotted(ul) {
		margin: 0 0 12px;
		padding-left: 22px;
	}

	.legal-content :slotted(li) {
		margin: 4px 0;
	}

	.legal-content :slotted(strong) {
		color: var(--color-on-surface);
	}

	.legal-content :slotted(a) {
		color: var(--color-primary);
		font-weight: 600;
		overflow-wrap: anywhere;
	}

	.legal-footer {
		display: flex;
		flex-wrap: wrap;
		align-items: center;
		justify-content: space-between;
		gap: 12px;
		padding: 24px 32px;
		border-top: 1px solid rgba(118, 117, 134, 0.15);
	}

	.legal-links {
		display: flex;
		gap: 20px;
	}

	.legal-links a {
		font-size: 13.5px;
		font-weight: 600;
		color: var(--color-on-surface-variant);
		text-decoration: none;
	}

	.legal-links a:hover,
	.legal-links a.router-link-active {
		color: var(--color-primary);
	}

	.copyright {
		font-size: 12.5px;
		color: var(--color-outline);
	}

	@media (max-width: 640px) {
		.legal-header,
		.legal-footer {
			padding-left: 16px;
			padding-right: 16px;
		}

		.legal-main {
			padding: 32px 16px 48px;
		}
	}

	@media (max-width: 400px) {
		.brand-name {
			display: none;
		}
	}
</style>
