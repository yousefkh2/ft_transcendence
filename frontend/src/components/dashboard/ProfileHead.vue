<template>
	<BaseCard class="profile-head">
		<BaseAvatar :name="name" :src="src" size="lg" variant="gradient" :status="status" />
		<div class="identity">
			<h1>{{ name ?? '…' }}</h1>
			<slot />
		</div>
		<div v-if="$slots.actions" class="actions">
			<slot name="actions" />
		</div>
	</BaseCard>
</template>

<script setup lang="ts">
import BaseAvatar from '../base/BaseAvatar.vue';
import BaseCard from '../base/BaseCard.vue';

defineProps<{ name?: string; src?: string; status?: 'online' | 'offline' }>();
</script>

<style scoped>
	.profile-head {
		display: flex;
		align-items: center;
		gap: 20px;
	}

	.identity {
		flex: 1;
		min-width: 0;
	}

	h1 {
		margin: 0 0 6px;
		font-size: 26px;
		font-weight: 800;
		line-height: 1.1;
		letter-spacing: -0.02em;
		overflow: hidden;
		text-overflow: ellipsis;
	}

	.actions {
		flex: none;
	}

	@media (max-width: 720px) {
		.profile-head {
			flex-wrap: wrap;
		}

		h1 {
			font-size: 22px;
		}
	}
</style>
