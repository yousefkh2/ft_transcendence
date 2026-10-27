<template>
	<div class="alert" :class="[variant, size]" :role="variant === 'error' ? 'alert' : 'status'">
		<p class="message"><slot /></p>
		<button v-if="onRetry" type="button" class="retry" @click="onRetry">Try again</button>
	</div>
</template>

<script setup lang="ts">
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
		margin: 0;
	}

	.sm {
		padding: 10px 14px;
		border-radius: 10px;
		font-size: 13px;
	}

	.md {
		padding: 16px 20px;
		border-radius: 12px;
		font-size: 14px;
	}

	.error {
		background: rgba(209, 53, 43, 0.08);
		color: var(--color-error);
	}

	.success {
		background: rgba(34, 197, 94, 0.1);
		color: var(--color-on-success-container);
	}

	.info {
		background: rgba(118, 117, 134, 0.1);
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
		background: rgba(209, 53, 43, 0.08);
	}

	.retry:focus-visible {
		outline: 2px solid currentColor;
		outline-offset: 2px;
	}
</style>
