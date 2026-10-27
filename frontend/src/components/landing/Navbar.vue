<template>
	<header class="navbar">
		<div class="brand">
			<img :src="logo" alt="LingoDash Logo" class="brand-logo" />
			<span class="brand-name">LingoDash</span>
		</div>

		<nav class="nav-links" :aria-label="t('landing.nav.main')">
			<a href="#how-it-works">{{ t('landing.nav.howItWorks') }}</a>
			<a href="#game-modes">{{ t('landing.nav.games') }}</a>
		</nav>

		<div class="auth-buttons">
			<LanguageSwitcher />
			<BaseButton variant="ghost" pill @click="emit('open-auth', 'login')">{{ t('landing.nav.login') }}</BaseButton>
			<BaseButton pill @click="emit('open-auth', 'register')">{{ t('landing.nav.signUp') }}</BaseButton>
		</div>
	</header>
</template>

<script setup lang="ts">
import BaseButton from '../base/BaseButton.vue';
import LanguageSwitcher from '../LanguageSwitcher.vue';
import { t } from '../../i18n';
import logo from "../../images/logo-transparent.png"

const emit = defineEmits<{ 'open-auth': [tab: 'login' | 'register'] }>();
</script>

<style scoped>
	.navbar {
		position: sticky;
		top: 0;
		z-index: 50;
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 24px;
		padding: 16px 40px;
		background: rgba(249, 249, 255, 0.72);
		backdrop-filter: blur(12px);
		-webkit-backdrop-filter: blur(12px);
		border-bottom: 1px solid var(--color-border);
		width: 80vw;
		margin: auto;
	}

	.brand {
		display: flex;
		align-items: center;
		gap: 10px;
	}

	.brand-logo {
		width: 45px;
		height: 42px;
		border-radius: var(--radius-md);
	}

	.brand-name {
		font-family: var(--font-display);
		font-weight: 800;
		font-size: 18px;
		letter-spacing: -0.01em;
		color: var(--color-on-surface);
	}

	.nav-links {
		display: flex;
		align-items: center;
		gap: 32px;
	}

	.nav-links a {
		position: relative;
		color: var(--color-on-surface-variant);
		text-decoration: none;
		font-weight: 600;
		font-size: var(--text-md);
		padding: 4px 0;
		transition: color 0.2s ease;
	}

	.nav-links a::after {
		content: "";
		position: absolute;
		left: 0;
		bottom: 0;
		width: 0%;
		height: 2px;
		background: var(--color-primary);
		transition: width 0.25s ease;
	}

	.nav-links a:hover {
		color: var(--color-on-surface);
	}

	.nav-links a:hover::after {
		width: 100%;
	}

	.auth-buttons {
		display: flex;
		align-items: center;
		gap: 8px;
	}

	@media (max-width: 720px) {
		.navbar {
			width: 100%;
			gap: 12px;
			padding: 12px 16px;
		}

		.nav-links {
			display: none;
		}

		.auth-buttons .btn {
			padding: 0 14px;
		}
	}

	/* Language switcher and both buttons no longer fit next to the logo: they move to their own line. */
	@media (max-width: 560px) {
		.navbar {
			flex-wrap: wrap;
		}

		.auth-buttons {
			flex-wrap: wrap;
			justify-content: flex-end;
			margin-left: auto;
		}
	}

	@media (max-width: 400px) {
		.brand-name {
			display: none;
		}
	}

</style>