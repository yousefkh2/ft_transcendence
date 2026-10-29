<template>
	<span class="badge" :class="variant">
		<BaseIcon v-if="icon" :name="icon" :size="13" />
		<slot />
	</span>
</template>

<script setup lang="ts">
import BaseIcon, { type IconName } from './BaseIcon.vue';

// A small pill for a status or label next to other content, like "Friends" or "Only visible to you".
withDefaults(
	defineProps<{
		variant?: 'neutral' | 'success' | 'primary';
		icon?: IconName;
	}>(),
	{ variant: 'neutral', icon: undefined },
);
</script>

<style scoped>
	.badge {
		display: inline-flex;
		align-items: center;
		gap: 5px;
		padding: 3px 10px;
		border-radius: var(--radius-pill);
		font-size: var(--text-xs);
		font-weight: 700;
		line-height: 1.4;
		white-space: nowrap;
	}

	.neutral {
		background: var(--color-fill-muted);
		color: var(--color-outline);
	}

	.success {
		background: var(--color-success-container);
		color: var(--color-on-success-container);
	}

	.primary {
		background: var(--color-focus-ring);
		color: var(--color-primary);
	}
</style>
