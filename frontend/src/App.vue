<template>
	<a href="#main" class="skip-link" @click.prevent="skipToContent">{{ t('common.skipToContent') }}</a>
	<RouterView />
	<!-- Single-page navigation does not load a new page, so screen readers would not notice it otherwise. -->
	<p class="visually-hidden" aria-live="polite">{{ announcement }}</p>
</template>

<script setup lang="ts">
import { nextTick, ref, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { locale, t } from './i18n';
import { updateTitle } from './router';

const announcement = ref('');

// Every page marks its main content with id="main" and tabindex="-1", so it can take focus.
function skipToContent() {
	document.getElementById('main')?.focus();
}

// The browser tab title follows a language switch too.
const route = useRoute();
watch(locale, () => updateTitle(route));

useRouter().afterEach((_to, from) => {
	// The first page load is announced by the browser itself.
	if (!from.matched.length) return;
	nextTick(() => (announcement.value = document.title));
});
</script>

<style scoped>
	.skip-link {
		position: fixed;
		top: 12px;
		left: 12px;
		z-index: 100;
		padding: 10px 16px;
		border-radius: 10px;
		background: var(--color-surface-inverse);
		color: var(--color-on-primary);
		font-weight: 700;
		font-size: 14px;
		text-decoration: none;
		transform: translateY(-200%);
	}

	.skip-link:focus {
		transform: none;
	}
</style>
