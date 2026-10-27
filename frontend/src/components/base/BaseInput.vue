<template>
	<div v-if="label" class="field">
		<label :for="inputId" class="field-label">{{ label }}</label>
		<input
			:id="inputId"
			v-model="model"
			v-bind="$attrs"
			class="input"
			:class="[size, { invalid: !!error }]"
			:aria-invalid="error ? 'true' : undefined"
			:aria-describedby="error || hint ? messageId : undefined"
		/>
		<span v-if="error" :id="messageId" class="field-error">{{ error }}</span>
		<span v-else-if="hint" :id="messageId" class="field-hint">{{ hint }}</span>
	</div>
	<input
		v-else
		v-model="model"
		v-bind="$attrs"
		class="input"
		:class="[size, { invalid: !!error }]"
		:aria-invalid="error ? 'true' : undefined"
	/>
</template>

<script setup lang="ts">
import { useId } from 'vue';

defineOptions({ inheritAttrs: false });

withDefaults(
	defineProps<{
		label?: string;
		size?: 'sm' | 'md';
		error?: string;
		hint?: string;
	}>(),
	{ size: 'md' },
);

const model = defineModel<string>({ default: '' });

const inputId = useId();
const messageId = `${inputId}-message`;
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
		box-shadow: 0 0 0 3px rgba(70, 72, 212, 0.18);
	}

	.input.invalid,
	.input.invalid:focus {
		border-color: var(--color-error);
	}

	.input.invalid:focus {
		box-shadow: 0 0 0 3px rgba(209, 53, 43, 0.15);
	}

	.field-error,
	.field-hint {
		font-size: 12px;
		line-height: 1.4;
	}

	.field-error {
		font-weight: 600;
		color: var(--color-error);
	}

	.field-hint {
		color: var(--color-outline);
	}

	.md {
		height: 47px;
	}

	.sm {
		height: 44px;
	}
</style>
