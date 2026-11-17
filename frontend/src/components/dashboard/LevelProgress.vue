<template>
	<BaseCard size="sm" class="level-card">
		<template v-if="level">
			<div class="level-top">
				<span class="level-badge" aria-hidden="true">{{ level.level }}</span>
				<div class="level-text">
					<span class="level-label">{{ t('level.label', { level: level.level }) }}</span>
					<span class="level-sub">{{ t('level.toNext', { n: formatNumber(level.xpForNext - level.xpIntoLevel), level: level.level + 1 }) }}</span>
				</div>
				<span class="level-xp">{{ t('level.progress', { xp: formatNumber(level.xpIntoLevel), needed: formatNumber(level.xpForNext) }) }}</span>
			</div>
			<div
				class="bar"
				role="progressbar"
				:aria-label="t('level.progressLabel', { level: level.level + 1 })"
				aria-valuemin="0"
				:aria-valuemax="level.xpForNext"
				:aria-valuenow="level.xpIntoLevel"
				:aria-valuetext="t('level.progress', { xp: level.xpIntoLevel, needed: level.xpForNext })"
			>
				<div class="bar-fill" :style="{ width: `${percent}%` }"></div>
			</div>
		</template>
		<template v-else>
			<BaseSkeleton width="140px" height="1.2em" />
			<BaseSkeleton class="bar-skeleton" width="100%" height="10px" shape="block" />
		</template>
	</BaseCard>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import type { LevelInfo } from '../../api/users';
import { t } from '../../i18n';
import { formatNumber } from '../../utils/format';
import BaseCard from '../base/BaseCard.vue';
import BaseSkeleton from '../base/BaseSkeleton.vue';

// level null: still loading.
const props = defineProps<{ level: LevelInfo | null | undefined }>();

const percent = computed(() => {
	const l = props.level;
	return l && l.xpForNext > 0 ? Math.min(100, (l.xpIntoLevel / l.xpForNext) * 100) : 0;
});
</script>

<style scoped>
	.level-card {
		display: grid;
		gap: 12px;
	}

	.level-top {
		display: flex;
		align-items: center;
		gap: 12px;
	}

	.level-badge {
		display: grid;
		place-items: center;
		width: 40px;
		height: 40px;
		flex-shrink: 0;
		border-radius: var(--radius-md);
		background: linear-gradient(135deg, var(--color-primary), var(--color-primary-container));
		color: var(--color-on-primary);
		font-family: var(--font-mono);
		font-weight: 800;
		font-size: var(--text-lg);
	}

	.level-text {
		flex: 1;
		min-width: 0;
		display: flex;
		flex-direction: column;
	}

	.level-label {
		font-weight: 800;
		font-size: 15px;
	}

	.level-sub,
	.level-xp {
		font-size: 12.5px;
		color: var(--color-outline);
	}

	.level-xp {
		font-family: var(--font-mono);
		font-weight: 700;
		white-space: nowrap;
	}

	.bar {
		height: 10px;
		border-radius: var(--radius-pill);
		background: var(--color-fill-muted);
		overflow: hidden;
	}

	.bar-fill {
		height: 100%;
		border-radius: inherit;
		background: linear-gradient(90deg, var(--color-primary), var(--color-primary-container));
		transition: width 0.4s ease;
	}

	.bar-skeleton {
		margin-top: 4px;
	}

	@media (prefers-reduced-motion: reduce) {
		.bar-fill {
			transition: none;
		}
	}
</style>
