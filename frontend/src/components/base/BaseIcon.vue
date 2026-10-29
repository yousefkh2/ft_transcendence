<template>
	<svg
		class="icon"
		:width="size"
		:height="size"
		viewBox="0 0 24 24"
		fill="none"
		stroke="currentColor"
		stroke-width="2"
		stroke-linecap="round"
		stroke-linejoin="round"
		:role="label ? 'img' : undefined"
		:aria-label="label"
		:aria-hidden="label ? undefined : 'true'"
		focusable="false"
	>
		<path v-for="(d, i) in ICONS[name]" :key="i" :d="d" />
	</svg>
</template>

<script setup lang="ts">
// LingoDash's own line icons on a 24 x 24 grid, drawn with the current text color.
// Each icon is a list of SVG path commands.
const ICONS = {
	close: ['M6 6l12 12', 'M18 6L6 18'],
	check: ['M5 12.5l4.5 4.5L19 7'],
	plus: ['M12 5v14', 'M5 12h14'],
	'alert-circle': ['M12 3a9 9 0 1 0 0 18a9 9 0 1 0 0-18z', 'M12 7.5v5.5', 'M12 16.5v.01'],
	'check-circle': ['M12 3a9 9 0 1 0 0 18a9 9 0 1 0 0-18z', 'M8 12.5l2.8 2.8L16 10'],
	info: ['M12 3a9 9 0 1 0 0 18a9 9 0 1 0 0-18z', 'M12 11v5.5', 'M12 7.5v.01'],
	play: ['M7 4.5v15l12-7.5z'],
	'arrow-left': ['M19 12H5', 'M11 6l-6 6l6 6'],
	maximize: ['M4 9V4h5', 'M20 9V4h-5', 'M4 15v5h5', 'M20 15v5h-5'],
	logout: ['M9 4H6a2 2 0 0 0-2 2v12a2 2 0 0 0 2 2h3', 'M14 8l4 4l-4 4', 'M18 12H8'],
	upload: ['M12 15V4', 'M7 9l5-5l5 5', 'M4 15v3a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2v-3'],
	trash: ['M4 7h16', 'M9 7V4h6v3', 'M6 7l1 13h10l1-13', 'M10 11v6', 'M14 11v6'],
	'user-plus': ['M9 4a4 4 0 1 0 0 8a4 4 0 1 0 0-8z', 'M2 20a7 7 0 0 1 14 0', 'M19 8v6', 'M16 11h6'],
	lock: ['M6 11h12v9H6z', 'M8.5 11V8a3.5 3.5 0 0 1 7 0v3'],
	globe: ['M12 3a9 9 0 1 0 0 18a9 9 0 1 0 0-18z', 'M3 12h18', 'M12 3c2.5 2.5 3.8 5.5 3.8 9s-1.3 6.5-3.8 9', 'M12 3c-2.5 2.5-3.8 5.5-3.8 9s1.3 6.5 3.8 9'],
} as const;

export type IconName = keyof typeof ICONS;

withDefaults(
	defineProps<{
		name: IconName;
		size?: number;
		// Without a label the icon is decorative and hidden from screen readers.
		label?: string;
	}>(),
	{ size: 18, label: undefined },
);
</script>

<style scoped>
	.icon {
		flex: none;
		display: inline-block;
		vertical-align: middle;
	}
</style>
