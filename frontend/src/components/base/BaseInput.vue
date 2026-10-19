<template>
	<label v-if="label" class="field">
		<span class="field-label">{{ label }}</span>
		<input v-model="model" v-bind="$attrs" class="input" :class="size" />
	</label>
	<input v-else v-model="model" v-bind="$attrs" class="input" :class="size" />
</template>

<script setup lang="ts">
defineOptions({ inheritAttrs: false });

withDefaults(
	defineProps<{
		label?: string;
		size?: 'sm' | 'md';
	}>(),
	{ size: 'md' },
);

const model = defineModel<string>({ default: '' });
</script>

<style scoped>
	.field {
		display: flex;
		flex-direction: column;
		gap: 6px;
	}

	.field-label {
		font-size: 13px;
		font-weight: 600;
		color: var(--color-on-surface-variant);
	}

	.input {
		padding: 0 14px;
		border-radius: 12px;
		border: 1px solid rgba(118, 117, 134, 0.18);
		background: var(--color-surface-raised);
		font-family: var(--font-display);
		font-size: 14px;
		color: var(--color-on-surface);
		transition: border-color 0.2s ease;
	}

	.input:focus {
		outline: none;
		border-color: var(--color-primary);
	}

	.md {
		height: 47px;
	}

	.sm {
		height: 44px;
	}
</style>
