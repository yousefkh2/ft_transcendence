<template>
	<BaseCard class="profile-head">
		<template v-if="name">
			<BaseAvatar :name="name" :src="src" size="lg" variant="gradient" :status="status" />
			<div class="identity">
				<h1>{{ name }}</h1>
				<slot />
			</div>
		</template>
		<template v-else>
			<BaseSkeleton shape="circle" width="72px" height="72px" />
			<div class="identity" role="status">
				<h1 class="visually-hidden">{{ t('profile.loading') }}</h1>
				<BaseSkeleton width="180px" height="26px" class="name-skeleton" />
				<BaseSkeleton width="240px" height="0.8em" />
			</div>
		</template>
		<div v-if="name && $slots.actions" class="actions">
			<slot name="actions" />
		</div>
	</BaseCard>
</template>

<script setup lang="ts">
import BaseAvatar from '../base/BaseAvatar.vue';
import BaseCard from '../base/BaseCard.vue';
import BaseSkeleton from '../base/BaseSkeleton.vue';
import { t } from '../../i18n';

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

	.name-skeleton {
		margin-bottom: 10px;
		max-width: 100%;
	}

	.actions {
		flex: none;
	}

	@media (max-width: 720px) {
		.profile-head {
			flex-wrap: wrap;
		}

		/* Buttons get their own line, so name and email keep the full width in every language. */
		.actions {
			flex-basis: 100%;
		}

		h1 {
			font-size: 22px;
		}
	}
</style>
