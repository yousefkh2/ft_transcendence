<template>
	<RouterLink v-if="to" :to="to" class="btn" :class="classes">
		<slot />
	</RouterLink>
	<button v-else :type="type" class="btn" :class="classes">
		<slot />
	</button>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import type { RouteLocationRaw } from 'vue-router';

const props = withDefaults(
	defineProps<{
		variant?: 'primary' | 'gradient' | 'outline' | 'ghost' | 'danger';
		size?: 'sm' | 'md' | 'lg';
		pill?: boolean;
		block?: boolean;
		to?: RouteLocationRaw;
		type?: 'button' | 'submit';
	}>(),
	{ variant: 'primary', size: 'md', pill: false, block: false, type: 'button' },
);

const classes = computed(() => [props.variant, props.size, { pill: props.pill, block: props.block }]);
</script>

<style scoped>
	.btn {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		gap: 8px;
		border: none;
		border-radius: 12px;
		font-family: var(--font-display);
		font-weight: 700;
		text-decoration: none;
		white-space: nowrap;
		cursor: pointer;
		transition: transform 0.15s ease, box-shadow 0.15s ease, filter 0.15s ease, background 0.15s ease;
	}

	.btn:focus-visible {
		outline: 2px solid var(--color-primary);
		outline-offset: 2px;
	}

	.btn:disabled {
		opacity: 0.6;
		cursor: not-allowed;
	}

	.sm {
		height: 30px;
		padding: 0 14px;
		font-size: 12px;
	}

	.md {
		height: 44px;
		padding: 0 20px;
		font-size: 14px;
	}

	.lg {
		height: 48px;
		padding: 0 24px;
		font-size: 15px;
	}

	.pill {
		border-radius: var(--radius-pill);
	}

	.block {
		width: 100%;
	}

	.primary {
		background: var(--color-primary);
		color: var(--color-on-primary);
	}

	.gradient {
		background: linear-gradient(135deg, var(--color-primary), var(--color-primary-container));
		color: var(--color-on-primary);
		box-shadow: 0 10px 24px -8px rgba(70, 72, 212, 0.55);
	}

	.primary:not(:disabled):hover,
	.gradient:not(:disabled):hover {
		transform: translateY(-2px);
		box-shadow: 0 10px 24px -6px rgba(70, 72, 212, 0.5);
	}

	.gradient:not(:disabled):hover {
		filter: brightness(1.08);
	}

	.primary:not(:disabled):active,
	.gradient:not(:disabled):active {
		transform: translateY(0);
	}

	.outline {
		background: transparent;
		color: var(--color-on-surface);
		border: 1.5px solid rgba(118, 117, 134, 0.35);
	}

	.ghost {
		background: transparent;
		color: var(--color-on-surface);
	}

	.danger {
		background: var(--color-error);
		color: var(--color-on-primary);
	}

	.danger:not(:disabled):hover {
		filter: brightness(1.08);
	}

	.outline:not(:disabled):hover,
	.ghost:not(:disabled):hover {
		background: rgba(17, 28, 45, 0.06);
	}
</style>
