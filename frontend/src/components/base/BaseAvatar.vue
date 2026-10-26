<template>
	<div class="avatar" :class="[size, variant]" :title="name">{{ initial }}</div>
</template>

<script setup lang="ts">
import { computed } from 'vue';

const props = withDefaults(
	defineProps<{
		name?: string;
		size?: 'sm' | 'md' | 'lg';
		variant?: 'solid' | 'gradient';
	}>(),
	{ size: 'sm', variant: 'solid' },
);

const initial = computed(() => props.name?.[0]?.toUpperCase() ?? '?');
</script>

<style scoped>
	.avatar {
		flex: none;
		display: grid;
		place-items: center;
		border-radius: var(--radius-pill);
		color: var(--color-on-primary);
		font-weight: 800;
	}

	.sm {
		width: 34px;
		height: 34px;
		font-size: 13px;
	}

	.md {
		width: 36px;
		height: 36px;
		font-size: 14px;
	}

	.lg {
		width: 72px;
		height: 72px;
		font-size: 28px;
	}

	.solid {
		background: var(--color-primary);
	}

	.gradient {
		background: linear-gradient(135deg, var(--color-primary), var(--color-primary-container));
	}
</style>
