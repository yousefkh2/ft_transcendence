<template>
	<main class="landing">
		<Navbar @open-auth="openAuth" />
		<Hero @open-auth="openAuth" />
		<HowItWorks />
		<GameModesPreview />
		<Gamification />
		<Footer />
		<AuthModal v-model:open="authOpen" v-model:tab="authTab" />
	</main>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import AuthModal from './landing/AuthModal.vue';
import Footer from './landing/Footer.vue';
import GameModesPreview from './landing/GameModesPreview.vue';
import Gamification from './landing/Gamification.vue';
import Hero from './landing/Hero.vue';
import HowItWorks from './landing/HowItWorks.vue';
import Navbar from './landing/Navbar.vue';

const authOpen = ref(false);
const authTab = ref<'login' | 'register'>('login');

function openAuth(tab: 'login' | 'register') {
	authTab.value = tab;
	authOpen.value = true;
}
</script>

<style scoped>
	.landing {
		background: var(--color-surface);
		color: var(--color-on-surface);
		font-family: var(--font-display);
		min-height: 100vh;
	}

</style>