<template>
	<div class="alert" :class="[variant, size]" :role="variant === 'error' ? 'alert' : 'status'">
		<p class="message">
			<BaseIcon :name="ICONS[variant]" :size="size === 'md' ? 18 : 16" class="message-icon" />
			<span><slot /></span>
		</p>
		<button v-if="onRetry" type="button" class="retry" @click="onRetry">Try again</button>
	</div>
</template>

<script setup lang="ts">
import BaseIcon from './BaseIcon.vue';

const ICONS = { error: 'alert-circle', success: 'check-circle', info: 'info' } as const;

withDefaults(
	defineProps<{
		variant?: 'error' | 'success' | 'info';
		size?: 'sm' | 'md';
		// Listening to @retry shows a "Try again" button.
		onRetry?: () => void;
	}>(),
	{ variant: 'error', size: 'sm', onRetry: undefined },
);
</script>

<style scoped>
	.alert {
		display: flex;
		align-items: center;
		justify-content: space-between;
		flex-wrap: wrap;
		gap: 6px 14px;
		font-weight: 600;
		line-height: 1.4;
	}

	.message {
		display: flex;
		align-items: flex-start;
		gap: 8px;
		margin: 0;
	}

	/* Lines the icon up with the first line of text. */
	.message-icon {
		margin-top: 1px;
	}

	.sm {
		padding: 10px 14px;
		border-radius: 10px;
		font-size: var(--text-sm);
	}

	.md {
		padding: 16px 20px;
		border-radius: 12px;
		font-size: var(--text-md);
	}

	.error {
		background: var(--color-error-container);
		color: var(--color-error);
	}

	.success {
		background: var(--color-success-container);
		color: var(--color-on-success-container);
	}

	.info {
		background: var(--color-fill-muted);
		color: var(--color-on-surface-variant);
	}

	.retry {
		padding: 4px 12px;
		border: 1px solid currentColor;
		border-radius: var(--radius-pill);
		background: none;
		color: inherit;
		font: inherit;
		font-size: 12.5px;
		font-weight: 700;
		cursor: pointer;
	}

	.retry:hover {
		background: var(--color-error-container);
	}

	.retry:focus-visible {
		outline: 2px solid currentColor;
		outline-offset: 2px;
	}
</style>
