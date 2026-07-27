<template>
	<footer class="footer">
		<div class="content">
			<div class="grid">
				<div class="brand-col">
					<div class="brand">
						<img :src="logo" alt="LingoDash Logo" class="brand-logo" />
						<span class="brand-name">LingoDash</span>
					</div>
					<p class="brand-desc">
						Real conversation, real pressure. Available now: German. Coming
						soon: Spanish, French, Japanese.
					</p>
				</div>

				<nav class="links">
					<a href="#how-it-works">How it Works</a>
					<a href="#game-modes">Game Modes</a>
					<a href="#gamification">Gamification</a>
				</nav>
			</div>

			<div class="bottom">
				<span class="copyright">© 2026 LingoDash. All rights reserved.</span>
			</div>
		</div>
	</footer>
</template>

<script setup lang="ts">
import logo from '../../images/logo-transparent.png';
</script>

<style scoped>
	.footer {
		border-top: 1px solid rgba(118, 117, 134, 0.15);
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
		margin-top: 32px;
		padding-top: 20px;
		border-top: 1px solid rgba(118, 117, 134, 0.12);
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
