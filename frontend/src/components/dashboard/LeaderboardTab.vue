<template>
	<div class="leaderboard">
		<div class="toolbar">
			<div>
				<h1>{{ t('leaderboard.title') }}</h1>
				<p class="subtitle">{{ t('leaderboard.subtitle') }}</p>
			</div>
			<div class="scope" role="group" :aria-label="t('leaderboard.scopeLabel')">
				<button
					v-for="option in SCOPES"
					:key="option"
					type="button"
					class="scope-btn"
					:class="{ active: scope === option }"
					:aria-pressed="scope === option"
					@click="scope = option"
				>
					{{ t(`leaderboard.${option}`) }}
				</button>
			</div>
		</div>

		<BaseCard v-if="status === 'ready' && board.me && !meInList" size="sm" class="my-rank">
			<span class="my-rank-label">{{ t('leaderboard.yourRank') }}</span>
			<span class="rank">#{{ board.me.rank }}</span>
			<span class="my-rank-meta">{{ t('level.label', { level: board.me.level.level }) }} · {{ t('common.xp', { xp: formatNumber(board.me.xp) }) }}</span>
		</BaseCard>

		<BaseCard tag="section">
			<ListSkeleton v-if="status === 'loading'" :label="t('leaderboard.loading')" avatar />
			<BaseAlert v-else-if="status === 'error'" @retry="load">{{ t('leaderboard.loadError') }}</BaseAlert>
			<p v-else-if="!board.entries.length" class="empty-hint">{{ t('leaderboard.empty') }}</p>
			<table v-else class="table">
				<caption class="visually-hidden">{{ t(`leaderboard.${scope}`) }}</caption>
				<thead>
					<tr>
						<th scope="col" class="col-rank">{{ t('leaderboard.rank') }}</th>
						<th scope="col">{{ t('leaderboard.player') }}</th>
						<th scope="col" class="col-num">{{ t('leaderboard.level') }}</th>
						<th scope="col" class="col-num">{{ t('leaderboard.xp') }}</th>
						<th scope="col" class="col-num hide-sm">{{ t('leaderboard.wins') }}</th>
						<th scope="col" class="col-num hide-sm">{{ t('leaderboard.rounds') }}</th>
						<th scope="col" class="col-num hide-md">{{ t('leaderboard.achievements') }}</th>
					</tr>
				</thead>
				<tbody>
					<tr v-for="entry in board.entries" :key="entry.userId" :class="{ me: entry.isMe, podium: entry.rank <= 3 }">
						<td class="col-rank">
							<span class="rank" :class="`rank-${entry.rank}`">
								<BaseIcon v-if="entry.rank === 1" name="trophy" :size="14" />{{ entry.rank }}
							</span>
						</td>
						<td>
							<div class="player">
								<BaseAvatar :name="entry.username" :src="apiAssetUrl(entry.avatarUrl)" />
								<RouterLink v-if="!entry.isMe" :to="{ name: 'user-profile', params: { username: entry.username } }" class="player-name">{{ entry.username }}</RouterLink>
								<RouterLink v-else :to="{ name: 'profile' }" class="player-name">{{ entry.username }}</RouterLink>
								<BaseBadge v-if="entry.isMe" variant="primary">{{ t('leaderboard.you') }}</BaseBadge>
							</div>
						</td>
						<td class="col-num"><span class="level-pill">{{ entry.level.level }}</span></td>
						<td class="col-num mono">{{ formatNumber(entry.xp) }}</td>
						<td class="col-num mono hide-sm">{{ formatNumber(entry.wins) }}</td>
						<td class="col-num mono hide-sm">{{ formatNumber(entry.rounds) }}</td>
						<td class="col-num mono hide-md">{{ entry.achievements }}</td>
					</tr>
				</tbody>
			</table>
		</BaseCard>
	</div>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { apiAssetUrl } from '../../api/client';
import { fetchLeaderboard, type Leaderboard, type LeaderboardScope } from '../../api/progress';
import { useLoader } from '../../composables/useLoader';
import { t } from '../../i18n';
import { formatNumber } from '../../utils/format';
import BaseAlert from '../base/BaseAlert.vue';
import BaseAvatar from '../base/BaseAvatar.vue';
import BaseBadge from '../base/BaseBadge.vue';
import BaseCard from '../base/BaseCard.vue';
import BaseIcon from '../base/BaseIcon.vue';
import ListSkeleton from './ListSkeleton.vue';

const SCOPES: LeaderboardScope[] = ['global', 'friends'];
const scope = ref<LeaderboardScope>('global');

const { data: board, status, load } = useLoader<Leaderboard>(() => fetchLeaderboard(scope.value), { entries: [] });

const meInList = computed(() => board.value.entries.some((e) => e.isMe));

watch(
	scope,
	() => {
		status.value = 'loading';
		load();
	},
	{ immediate: true },
);
</script>

<style scoped>
	.leaderboard {
		display: grid;
		grid-template-columns: minmax(0, 1fr);
		gap: 22px;
	}

	.toolbar {
		display: flex;
		justify-content: space-between;
		align-items: flex-end;
		gap: 20px;
		flex-wrap: wrap;
	}

	.toolbar h1 {
		font-size: 26px;
		line-height: 0.95;
		font-weight: 800;
		margin: 0 0 6px;
		letter-spacing: -0.02em;
	}

	.subtitle {
		margin: 0;
		color: var(--color-outline);
		font-size: 14.5px;
	}

	.scope {
		display: flex;
		gap: 4px;
		padding: 4px;
		border-radius: 13px;
		background: var(--color-fill-muted);
	}

	.scope-btn {
		padding: 7px 16px;
		border: none;
		border-radius: 10px;
		background: transparent;
		color: var(--color-outline);
		font-family: inherit;
		font-weight: 700;
		font-size: var(--text-sm);
		cursor: pointer;
	}

	.scope-btn.active {
		background: var(--color-surface-raised);
		color: var(--color-primary);
		box-shadow: var(--shadow-md);
	}

	.scope-btn:focus-visible,
	.player-name:focus-visible {
		outline: 2px solid var(--color-primary);
		outline-offset: 2px;
	}

	.my-rank {
		display: flex;
		align-items: center;
		gap: 12px;
		flex-wrap: wrap;
	}

	.my-rank-label {
		font-weight: 800;
	}

	.my-rank-meta {
		font-size: var(--text-sm);
		color: var(--color-outline);
	}

	.table {
		width: 100%;
		border-collapse: collapse;
	}

	th {
		text-align: left;
		padding: 0 10px 10px;
		font-family: var(--font-mono);
		font-size: var(--text-2xs);
		font-weight: 700;
		letter-spacing: 0.08em;
		text-transform: uppercase;
		color: var(--color-outline);
	}

	td {
		padding: 10px;
		border-top: 1px solid var(--color-border-subtle);
		font-size: var(--text-md);
	}

	tr.me td {
		background: var(--color-focus-ring);
	}

	.col-rank {
		width: 64px;
	}

	.col-num {
		text-align: right;
		white-space: nowrap;
	}

	.mono {
		font-family: var(--font-mono);
		font-weight: 700;
	}

	.rank {
		display: inline-flex;
		align-items: center;
		gap: 4px;
		font-family: var(--font-mono);
		font-weight: 800;
	}

	.rank-1 {
		color: var(--color-primary);
	}

	.player {
		display: flex;
		align-items: center;
		gap: 10px;
		min-width: 0;
	}

	.player-name {
		color: inherit;
		text-decoration: none;
		font-weight: 700;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}

	.player-name:hover {
		text-decoration: underline;
	}

	.level-pill {
		display: inline-block;
		min-width: 28px;
		padding: 2px 8px;
		border-radius: var(--radius-pill);
		background: var(--color-fill-muted);
		font-family: var(--font-mono);
		font-weight: 800;
		font-size: var(--text-xs);
		text-align: center;
	}

	.empty-hint {
		margin: 0;
		font-size: var(--text-sm);
		color: var(--color-outline);
	}

	.visually-hidden {
		position: absolute;
		width: 1px;
		height: 1px;
		overflow: hidden;
		clip: rect(0 0 0 0);
		white-space: nowrap;
	}

	@media (max-width: 860px) {
		.hide-md {
			display: none;
		}
	}

	@media (max-width: 600px) {
		.hide-sm {
			display: none;
		}

		td,
		th {
			padding-left: 6px;
			padding-right: 6px;
		}

		.col-rank {
			width: 44px;
		}
	}
</style>
