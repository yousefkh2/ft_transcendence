<template>
	<div class="switcher">
		<BaseIcon name="globe" :size="16" />
		<!-- Named directly: a wrapping <label> would add all option names to the select's name. -->
		<select :value="locale" :aria-label="t('common.language')" @change="onChange">
			<option v-for="option in LOCALES" :key="option.code" :value="option.code" :lang="option.code">
				{{ option.name }}
			</option>
		</select>
	</div>
</template>

<script setup lang="ts">
import { LOCALES, locale, setLocale, t, type Locale } from '../i18n';
import BaseIcon from './base/BaseIcon.vue';

// Each language is listed in its own name, so people find theirs whatever is selected.
function onChange(event: Event) {
	setLocale((event.target as HTMLSelectElement).value as Locale);
}
</script>

<style scoped>
	.switcher {
		position: relative;
		display: inline-flex;
		align-items: center;
		gap: 6px;
		padding: 0 10px;
		height: 34px;
		border-radius: var(--radius-pill);
		border: 1px solid var(--color-border-strong);
		background: var(--color-surface-raised);
		color: var(--color-on-surface-variant);
	}

	select {
		appearance: none;
		border: none;
		background: none;
		color: inherit;
		font: inherit;
		font-size: var(--text-sm);
		font-weight: 700;
		cursor: pointer;
	}

	select:focus {
		outline: none;
	}

	.switcher:focus-within {
		outline: 2px solid var(--color-primary);
		outline-offset: 2px;
	}
</style>
