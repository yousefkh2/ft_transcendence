<template>
	<span class="skeleton" :class="shape" :style="{ width, height }" aria-hidden="true"></span>
</template>

<script setup lang="ts">
// A gray placeholder in the shape of content that is still loading. It is hidden from screen readers:
// the surrounding loading area announces "Loading ..." once instead.
withDefaults(
	defineProps<{
		shape?: 'text' | 'circle' | 'block';
		width?: string;
		height?: string;
	}>(),
	{ shape: 'text', width: '100%', height: undefined },
);
</script>

<style scoped>
	.skeleton {
		display: block;
		flex: none;
		background: rgba(118, 117, 134, 0.14);
		animation: pulse 1.4s ease-in-out infinite;
	}

	.text {
		height: 0.8em;
		border-radius: 6px;
	}

	.circle {
		border-radius: var(--radius-pill);
	}

	.block {
		border-radius: 13px;
	}

	@keyframes pulse {
		50% {
			opacity: 0.45;
		}
	}

	@media (prefers-reduced-motion: reduce) {
		.skeleton {
			animation: none;
		}
	}
</style>
