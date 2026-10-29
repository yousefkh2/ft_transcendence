<template>
	<div class="list-skeleton" role="status">
		<span class="visually-hidden">{{ label }}</span>
		<div v-for="n in rows" :key="n" class="row">
			<BaseSkeleton :shape="avatar ? 'circle' : 'block'" width="34px" height="34px" />
			<div class="lines">
				<BaseSkeleton :width="n % 2 ? '45%' : '60%'" />
				<BaseSkeleton width="30%" height="0.65em" />
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import BaseSkeleton from '../base/BaseSkeleton.vue';

// Placeholder rows in the shape of the friend and round lists while they load.
withDefaults(defineProps<{ label: string; rows?: number; avatar?: boolean }>(), { rows: 3, avatar: false });
</script>

<style scoped>
	.list-skeleton {
		display: grid;
		gap: 10px;
	}

	.row {
		display: flex;
		align-items: center;
		gap: 12px;
		padding: 12px 14px;
		border-radius: 13px;
		border: 1px solid var(--color-border-subtle);
	}

	.lines {
		flex: 1;
		display: grid;
		gap: 7px;
	}
</style>
