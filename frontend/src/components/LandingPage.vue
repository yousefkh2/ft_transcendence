<template>
	<main class="landing">
		<Navbar />
		<Hero />
		<HowItWorks />
		<GameModesPreview />
		<Gamification />
	</main>
</template>

<script setup lang="ts">
import GameModesPreview from './landing/GameModesPreview.vue';
import Gamification from './landing/Gamification.vue';
import Hero from './landing/Hero.vue';
import HowItWorks from './landing/HowItWorks.vue';
import Navbar from './landing/Navbar.vue';
</script>

<style scoped>
	.landing {
		--color-primary: #4648d4;
		--color-primary-container: #6063ee;
		--color-on-primary: #ffffff;
		--color-secondary-container: #fea619;
		--color-on-secondary-container: #684000;
		--color-surface: #f9f9ff;
		--color-on-surface: #111c2d;
		--color-on-surface-variant: #464554;
		--color-outline: #767586;

		--font-display: "Plus Jakarta Sans", sans-serif;
		--font-mono: "Space Grotesk", monospace;

		--radius-lg: 1rem;
		--radius-pill: 9999px;

		background: var(--color-surface);
		color: var(--color-on-surface);
		font-family: var(--font-display);
		min-height: 100vh;
	}

</style>