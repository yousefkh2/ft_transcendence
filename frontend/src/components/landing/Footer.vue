<template>
	<footer class="footer">
		<div class="content">
			<div class="grid">
				<div class="brand-col">
					<div class="brand">
						<img :src="logo" alt="LingoDash Logo" class="brand-logo" />
						<span class="brand-name">LingoDash</span>
					</div>
					<p class="brand-desc">{{ t('landing.footer.text') }}</p>
				</div>

				<nav class="links" :aria-label="t('landing.nav.footer')">
					<a href="#how-it-works">{{ t('landing.nav.howItWorks') }}</a>
					<a href="#game-modes">{{ t('landing.nav.gameModes') }}</a>
					<a href="#gamification">{{ t('landing.nav.gamification') }}</a>
				</nav>
			</div>

			<div class="bottom">
				<span class="copyright">{{ t('landing.footer.rights') }}</span>
				<nav class="legal-links" :aria-label="t('landing.nav.legal')">
					<RouterLink to="/privacy">{{ t('common.privacyPolicy') }}</RouterLink>
					<RouterLink to="/terms">{{ t('common.termsOfService') }}</RouterLink>
				</nav>
			</div>
		</div>
	</footer>
</template>

<script setup lang="ts">
import logo from '../../images/logo-transparent.png';
import { t } from '../../i18n';
</script>

<style scoped>
	.footer {
		border-top: 1px solid var(--color-border);
		padding: 48px 40px;
	}

	.content {
		max-width: 1200px;
		margin: 0 auto;
	}

	.grid {
		display: grid;
		grid-template-columns: minmax(240px, 1.3fr) minmax(200px, 1fr);
		gap: 32px;
	}

	.brand {
		display: flex;
		align-items: center;
		gap: 10px;
		margin-bottom: 12px;
	}

	.brand-logo {
		width: 30px;
		height: 30px;
	}

	.brand-name {
		font-family: var(--font-display);
		font-weight: 700;
		font-size: 16px;
		color: var(--color-on-surface);
	}

	.brand-desc {
		max-width: 420px;
		font-size: 13.5px;
		line-height: 1.5;
		color: var(--color-outline);
		margin: 0;
	}

	.links {
		display: flex;
		flex-wrap: wrap;
		align-items: flex-start;
		gap: 20px;
	}

	.links a {
		font-size: 13.5px;
		font-weight: 600;
		color: var(--color-on-surface-variant);
		text-decoration: none;
	}

	.links a:hover {
		color: var(--color-on-surface);
	}

	.bottom {
		display: flex;
		flex-wrap: wrap;
		align-items: center;
		justify-content: space-between;
		gap: 12px;
		margin-top: 32px;
		padding-top: 20px;
		border-top: 1px solid var(--color-border-subtle);
	}

	.legal-links {
		display: flex;
		flex-wrap: wrap;
		gap: 8px 20px;
	}

	.legal-links a {
		font-size: 12.5px;
		font-weight: 600;
		color: var(--color-outline);
		text-decoration: none;
	}

	.legal-links a:hover {
		color: var(--color-on-surface);
	}

	.copyright {
		font-size: 12.5px;
		color: var(--color-outline);
	}

	@media (max-width: 640px) {
		.grid {
			grid-template-columns: 1fr;
		}
	}
</style>
