<template>
	<div class="round-list">
		<div v-for="round in rows" :key="round.sessionId" class="round-row">
			<div class="round-badge" :class="round.result" :title="resultLabel(round.result)">
				<BaseIcon v-if="round.result !== 'pending'" :name="round.result === 'won' ? 'check' : 'close'" :size="16" />
				<span v-else aria-hidden="true">-</span>
				<span class="visually-hidden">{{ resultLabel(round.result) }}</span>
			</div>
			<div class="round-info">
				<div class="round-mode">{{ round.mode }}</div>
				<div class="round-when">{{ round.details }}</div>
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import type { MatchHistoryEntry } from '../../api/users';
import { formatLabel, formatWhen, roundResult } from '../../utils/matches';
import { t } from '../../i18n';
import BaseIcon from '../base/BaseIcon.vue';

const props = defineProps<{ matches: MatchHistoryEntry[]; showRole?: boolean }>();

// The badge shows a language-neutral icon; the word is for screen readers and the tooltip.
function resultLabel(result: 'won' | 'lost' | 'pending') {
	return result === 'pending' ? t('rounds.pending') : t(result === 'won' ? 'rounds.won' : 'rounds.lost');
}

const rows = computed(() =>
	props.matches.map((m) => ({
		sessionId: m.sessionId,
		result: roundResult(m),
		mode: formatLabel(m.gameMode),
		details: props.showRole && m.role ? `${formatWhen(m.startedAt)} · ${formatLabel(m.role)}` : formatWhen(m.startedAt),
	})),
);
</script>

<style scoped>
	.round-list {
		display: grid;
		gap: 10px;
	}

	.round-row {
		display: flex;
		align-items: center;
		gap: 14px;
		padding: 13px 15px;
		border-radius: 13px;
		background: var(--color-surface);
		border: 1px solid var(--color-border-subtle);
	}

	.round-badge {
		width: 36px;
		height: 36px;
		flex: none;
		border-radius: 10px;
		display: grid;
		place-items: center;
		font-family: var(--font-mono);
		font-weight: 700;
		font-size: var(--text-2xs);
		color: var(--color-on-primary);
	}

	.round-badge.won {
		background: var(--color-success);
	}

	.round-badge.lost {
		background: var(--color-tertiary);
	}

	.round-badge.pending {
		background: var(--color-outline);
	}

	.round-info {
		flex: 1;
		min-width: 0;
	}

	.round-mode {
		font-size: var(--text-md);
		font-weight: 700;
	}

	.round-when {
		font-size: 12.5px;
		color: var(--color-outline);
	}
</style>
