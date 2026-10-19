<template>
	<div class="bolt" :class="size" aria-hidden="true"></div>
</template>

<script setup lang="ts">
withDefaults(defineProps<{ size?: 'sm' | 'lg' }>(), { size: 'sm' });
</script>

<style scoped>
	.bolt {
		background: linear-gradient(160deg, var(--color-secondary-container) 0%, var(--color-tertiary) 45%, var(--color-primary-container) 100%);
		clip-path: polygon(55% 0%, 25% 55%, 45% 55%, 35% 100%, 80% 40%, 55% 40%, 70% 0%);
	}

	.sm {
		width: 30px;
		height: 48px;
	}

	.lg {
		width: 170px;
		height: 270px;
	}
</style>
