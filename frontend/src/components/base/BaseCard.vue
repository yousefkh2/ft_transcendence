<template>
	<component :is="tag" class="card" :class="size">
		<slot />
	</component>
</template>

<script setup lang="ts">
withDefaults(
	defineProps<{
		tag?: string;
		size?: 'sm' | 'md' | 'lg';
	}>(),
	{ tag: 'div', size: 'md' },
);
</script>

<style scoped>
	.card {
		background: var(--color-surface-raised);
		border: 1px solid rgba(118, 117, 134, 0.18);
	}

	.sm {
		padding: 22px;
		border-radius: 16px;
	}

	.md {
		padding: 26px;
		border-radius: 18px;
	}

	.lg {
		padding: 28px;
		border-radius: 16px;
	}
</style>
