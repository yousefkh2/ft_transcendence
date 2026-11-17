<template>
	<BaseCard tag="section" class="achievements">
		<div class="panel-header">
			<h2>{{ t('achievements.title') }}</h2>
			<span v-if="status === 'ready'" class="count">{{ t('achievements.unlockedCount', { n: unlockedCount, total: achievements.length }) }}</span>
		</div>
		<ListSkeleton v-if="status === 'loading'" :label="t('achievements.loading')" />
		<BaseAlert v-else-if="status === 'error'" @retry="emit('retry')">{{ t('achievements.loadError') }}</BaseAlert>
		<ul v-else class="grid">
			<li v-for="a in achievements" :key="a.id" class="item" :class="{ unlocked: a.unlocked }">
				<span class="icon" aria-hidden="true"><BaseIcon :name="ACHIEVEMENT_ICONS[a.id]" :size="20" /></span>
				<div class="body">
					<h3 class="title">{{ t(`achievements.items.${a.id}.title`) }}</h3>
					<p class="text">{{ t(`achievements.items.${a.id}.text`) }}</p>
					<p v-if="a.unlocked && a.unlockedAt" class="meta done">
						<BaseIcon name="check" :size="12" />{{ t('achievements.unlockedOn', { date: formatDate(a.unlockedAt) }) }}
					</p>
					<template v-else>
						<span class="visually-hidden">{{ t('achievements.locked') }}</span>
						<div v-if="a.target > 1" class="progress">
							<div
								class="bar"
								role="progressbar"
								:aria-label="t(`achievements.items.${a.id}.title`)"
								aria-valuemin="0"
								:aria-valuemax="a.target"
								:aria-valuenow="a.progress"
							>
								<div class="bar-fill" :style="{ width: `${(a.progress / a.target) * 100}%` }"></div>
							</div>
							<span class="meta">{{ t('achievements.progress', { n: a.progress, target: a.target }) }}</span>
						</div>
					</template>
				</div>
			</li>
		</ul>
	</BaseCard>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import type { Achievement } from '../../api/progress';
import type { LoadStatus } from '../../composables/useLoader';
import { ACHIEVEMENT_ICONS } from '../../constants/achievements';
import { t } from '../../i18n';
import { formatDate } from '../../utils/format';
import BaseAlert from '../base/BaseAlert.vue';
import BaseCard from '../base/BaseCard.vue';
import BaseIcon from '../base/BaseIcon.vue';
import ListSkeleton from './ListSkeleton.vue';

const props = defineProps<{ achievements: Achievement[]; status: LoadStatus }>();
const emit = defineEmits<{ retry: [] }>();

const unlockedCount = computed(() => props.achievements.filter((a) => a.unlocked).length);
</script>

<style scoped>
	.panel-header {
		display: flex;
		justify-content: space-between;
		align-items: baseline;
		gap: 12px;
		margin-bottom: 16px;
	}

	h2 {
		margin: 0;
		font-size: var(--text-lg);
		font-weight: 800;
	}

	.count {
		font-family: var(--font-mono);
		font-size: var(--text-xs);
		font-weight: 700;
		color: var(--color-outline);
	}

	.grid {
		list-style: none;
		margin: 0;
		padding: 0;
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
		gap: 10px;
	}

	.item {
		display: flex;
		gap: 12px;
		padding: 12px;
		border-radius: var(--radius-md);
		border: 1px solid var(--color-border-subtle);
		background: var(--color-surface);
	}

	.icon {
		display: grid;
		place-items: center;
		width: 38px;
		height: 38px;
		flex-shrink: 0;
		border-radius: var(--radius-sm);
		background: var(--color-fill-muted);
		color: var(--color-outline);
	}

	.unlocked .icon {
		background: linear-gradient(135deg, var(--color-primary), var(--color-primary-container));
		color: var(--color-on-primary);
	}

	.body {
		flex: 1;
		min-width: 0;
	}

	.title {
		margin: 0;
		font-size: 13.5px;
		font-weight: 800;
	}

	.item:not(.unlocked) .title {
		color: var(--color-on-surface-variant);
	}

	.text {
		margin: 2px 0 6px;
		font-size: 12px;
		color: var(--color-outline);
	}

	.meta {
		display: inline-flex;
		align-items: center;
		gap: 4px;
		margin: 0;
		font-size: var(--text-2xs);
		font-weight: 700;
		color: var(--color-outline);
	}

	.meta.done {
		color: var(--color-on-success-container);
	}

	.progress {
		display: flex;
		align-items: center;
		gap: 8px;
	}

	.bar {
		flex: 1;
		height: 6px;
		border-radius: var(--radius-pill);
		background: var(--color-fill-muted);
		overflow: hidden;
	}

	.bar-fill {
		height: 100%;
		background: var(--color-primary);
	}

	.visually-hidden {
		position: absolute;
		width: 1px;
		height: 1px;
		overflow: hidden;
		clip: rect(0 0 0 0);
		white-space: nowrap;
	}
</style>
