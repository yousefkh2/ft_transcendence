<template>
	<div class="avatar" :class="[size, variant]" :title="name">
		{{ initial }}
		<span v-if="status" class="status" :class="status" aria-hidden="true"></span>
	</div>
</template>

<script setup lang="ts">
import { computed } from 'vue';

const props = withDefaults(
	defineProps<{
		name?: string;
		size?: 'sm' | 'md' | 'lg';
		variant?: 'solid' | 'gradient';
		// Only a visual hint: show the status as text next to the avatar for screen readers.
		status?: 'online' | 'offline';
	}>(),
	{ size: 'sm', variant: 'solid' },
);

const initial = computed(() => props.name?.[0]?.toUpperCase() ?? '?');
</script>

<style scoped>
	.avatar {
		position: relative;
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

	.status {
		position: absolute;
		right: -1px;
		bottom: -1px;
		width: 11px;
		height: 11px;
		border-radius: var(--radius-pill);
		border: 2px solid var(--color-surface-raised);
	}

	.lg .status {
		right: 2px;
		bottom: 2px;
		width: 16px;
		height: 16px;
		border-width: 3px;
	}

	.status.online {
		background: var(--color-success);
	}

	.status.offline {
		background: var(--color-outline);
	}
</style>
