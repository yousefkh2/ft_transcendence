<template>
	<div class="field">
		<label :for="selectId" class="field-label">{{ label }}</label>
		<select :id="selectId" v-model="model" v-bind="$attrs" class="select">
			<option v-for="option in options" :key="option.value" :value="option.value" :disabled="option.disabled">
				{{ option.label }}
			</option>
		</select>
	</div>
</template>

<script setup lang="ts">
import { useId } from 'vue';

defineOptions({ inheritAttrs: false });

defineProps<{
	label: string;
	options: { value: string; label: string; disabled?: boolean }[];
}>();

const model = defineModel<string>({ default: '' });

const selectId = useId();
</script>

<style scoped>
	.field {
		display: flex;
		flex-direction: column;
		gap: 6px;
	}

	.field-label {
		font-size: var(--text-sm);
		font-weight: 600;
		color: var(--color-on-surface-variant);
	}

	.select {
		height: 47px;
		padding: 0 14px;
		border-radius: 12px;
		border: 1px solid var(--color-border-strong);
		background: var(--color-surface-raised);
		font-family: var(--font-display);
		font-size: var(--text-md);
		color: var(--color-on-surface);
		transition: border-color 0.2s ease;
	}

	.select:focus {
		outline: none;
		border-color: var(--color-primary);
		box-shadow: 0 0 0 3px var(--color-focus-ring);
	}
</style>
