<template>
	<span class="eyebrow" :class="variant"><slot /></span>
</template>

<script setup lang="ts">
withDefaults(defineProps<{ variant?: 'plain' | 'chip' }>(), { variant: 'plain' });
</script>

<style scoped>
	.eyebrow {
		font-family: var(--font-mono);
		font-weight: 700;
		letter-spacing: 0.1em;
		text-transform: uppercase;
	}

	.plain {
		display: inline-block;
		margin: 0 0 10px;
		font-size: 13px;
		color: var(--color-primary);
	}

	.chip {
		display: inline-flex;
		margin-bottom: 14px;
		padding: 5px 12px;
		border-radius: var(--radius-pill);
		background: rgba(255, 255, 255, 0.14);
		font-size: 11px;
		color: var(--color-on-primary);
	}
</style>
