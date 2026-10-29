<template>
	<div class="avatar" :class="[size, variant]" :title="name">
		<img v-if="src && !failed" :src="src" alt="" class="picture" @error="failed = true" />
		<template v-else>{{ initial }}</template>
		<span v-if="status" class="status" :class="status" aria-hidden="true"></span>
	</div>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue';

const props = withDefaults(
	defineProps<{
		name?: string;
		// Profile picture URL. The initial is shown without one or when it fails to load.
		src?: string;
		size?: 'sm' | 'md' | 'lg';
		variant?: 'solid' | 'gradient';
		// Only a visual hint: show the status as text next to the avatar for screen readers.
		status?: 'online' | 'offline';
	}>(),
	{ size: 'sm', variant: 'solid' },
);

const initial = computed(() => props.name?.[0]?.toUpperCase() ?? '?');

const failed = ref(false);
watch(
	() => props.src,
	() => (failed.value = false),
);
</script>

<style scoped>
	.avatar {
		position: relative;
		flex: none;
		display: grid;
		place-items: center;
		border-radius: var(--radius-pill);
		color: var(--color-on-primary);
		font-weight: 800;
	}

	.sm {
		width: 34px;
		height: 34px;
		font-size: var(--text-sm);
	}

	.md {
		width: 36px;
		height: 36px;
		font-size: var(--text-md);
	}

	.lg {
		width: 72px;
		height: 72px;
		font-size: 28px;
	}

	.solid {
		background: var(--color-primary);
	}

	.gradient {
		background: linear-gradient(135deg, var(--color-primary), var(--color-primary-container));
	}

	.picture {
		width: 100%;
		height: 100%;
		border-radius: inherit;
		object-fit: cover;
	}

	.status {
		position: absolute;
		right: -1px;
		bottom: -1px;
		width: 11px;
		height: 11px;
		border-radius: var(--radius-pill);
		border: 2px solid var(--color-surface-raised);
	}

	.lg .status {
		right: 2px;
		bottom: 2px;
		width: 16px;
		height: 16px;
		border-width: 3px;
	}

	.status.online {
		background: var(--color-success);
	}

	.status.offline {
		background: var(--color-outline);
	}
</style>
