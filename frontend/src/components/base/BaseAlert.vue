<template>
	<p class="alert" :class="[variant, size]" :role="variant === 'error' ? 'alert' : 'status'">
		<slot />
	</p>
</template>

<script setup lang="ts">
withDefaults(
	defineProps<{
		variant?: 'error' | 'success' | 'info';
		size?: 'sm' | 'md';
	}>(),
	{ variant: 'error', size: 'sm' },
);
</script>

<style scoped>
	.alert {
		margin: 0;
		font-weight: 600;
		line-height: 1.4;
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
</style>
