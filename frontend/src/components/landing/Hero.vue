<template>
	<section class="hero">
		<div class="hero-copy">
			<h1 class="headline">
				Learn Languages
				<span class="accent">Through Play.</span>
			</h1>
			<p class="subtext">
				Jump into fast-paced multiplayer games built to get you speaking from day one. Describe, argue, guess, and race the clock - whichever language you're learning.
			</p>
			<div class="cta-row">
				<button class="btn btn-primary">Play Now</button>
				<button class="btn btn-ghost">How It Works</button>
			</div>
		</div>
		<HeroVisual class="hero-visual" />
	</section>
</template>

<script setup lang="ts">
import HeroVisual from './HeroVisual.vue';
</script>

<style scoped>
	.hero {
		display: grid;
		grid-template-columns: 1fr 1fr;
		align-items: center;
		gap: 48px;
		max-width: 1200px;
		margin: 0 auto;
		padding: 96px 40px;
	}

	.headling {
		font-size: clamp(32px, 5vw, 48px);
		font-weight: 800;
		line-height: 1.1;
		letter-spacing: -0.2em;
		margin: 0 0 20px;
	}

	.accent {
		color: var(--color-primary);
	}

	.subtext {
		max-width: 480px;
		color: var(--color-on-surface-variant);
		font-size: 16px;
		line-height: 1.6;
		margin: 32px 0 32px;
	}

	.cta-row {
		display: flex;
		gap: 12px;
		flex-wrap: wrap;
	}

	.btn {
		font-family: var(--font-display);
		font-weight: 700;
		font-size: 15px;
		padding: 14px 24px;
		border-radius: var(--radius-pill);
		border: none;
		cursor: pointer;
	}

	.btn-primary {
		background: var(--color-primary);
		color: var(--color-on-primary);
		border-bottom: 4px solid #2f2ebe;
	}

	.btn-primary:active {
		border-bottom-width: 0;
		transform: translateY(4px);
	}

	.btn-ghost {
		background: transparent;
		color: var(--color-on-surface);
		border: 2px solid var(--color-outline);
	}

	@media (max-width: 900px) {
		.hero {
			grid-template-columns: 1fr;
			text-align: center;
			padding: 48px 16px;
		}

		.subtext {
			margin-left: auto;
			margin-right: auto;
		}

		.hero-visual {
			order: -1;
		}

		.cta-row {
			justify-content: center;
		}
	}


</style>