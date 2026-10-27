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
		font-size: 13px;
		font-weight: 600;
		color: var(--color-on-surface-variant);
	}

	.select {
		height: 47px;
		padding: 0 14px;
		border-radius: 12px;
		border: 1px solid rgba(118, 117, 134, 0.18);
		background: var(--color-surface-raised);
		font-family: var(--font-display);
		font-size: 14px;
		color: var(--color-on-surface);
		transition: border-color 0.2s ease;
	}

	.select:focus {
		outline: none;
		border-color: var(--color-primary);
		box-shadow: 0 0 0 3px rgba(70, 72, 212, 0.18);
	}
</style>
