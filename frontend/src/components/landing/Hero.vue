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
				<BaseButton size="lg" pill @click="emit('open-auth', 'register')">Play Now</BaseButton>
				<BaseButton variant="outline" size="lg" pill @click="scrollToHowItWorks">How It Works</BaseButton>
			</div>
		</div>
		<HeroVisual class="hero-visual" />
	</section>
</template>

<script setup lang="ts">
import BaseButton from '../base/BaseButton.vue';
import HeroVisual from './HeroVisual.vue';

const emit = defineEmits<{ 'open-auth': [tab: 'login' | 'register'] }>();

function scrollToHowItWorks() {
	document.getElementById('how-it-works')?.scrollIntoView();
}
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

	.headline {
		max-width: 780px;
		margin: 0;
		font-size: clamp(42px, 9vw, 84px);
		line-height: 0.95;
		letter-spacing: 0;
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