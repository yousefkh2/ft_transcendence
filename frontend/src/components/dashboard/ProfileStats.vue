<template>
	<section class="stat-grid" aria-label="Stats" :aria-busy="loading || xp === null">
		<BaseCard v-for="stat in stats" :key="stat.label" size="sm">
			<span class="stat-label">{{ stat.label }}</span>
			<template v-if="stat.value === null">
				<BaseSkeleton class="value-skeleton" width="70px" height="30px" shape="block" />
				<BaseSkeleton width="90px" height="0.75em" />
			</template>
			<template v-else>
				<div class="stat-value">{{ stat.value }}</div>
				<p class="stat-sub">{{ stat.sub }}</p>
			</template>
		</BaseCard>
	</section>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import BaseCard from '../base/BaseCard.vue';
import BaseSkeleton from '../base/BaseSkeleton.vue';

// xp null: still loading the profile. loading: the round stats are still loading.
const props = defineProps<{
	xp: number | null;
	rounds: number;
	wins: number;
	losses: number;
	roundsSub: string;
	loading?: boolean;
}>();

const stats = computed(() => {
	const decided = props.wins + props.losses;
	const round = <T,>(value: T) => (props.loading ? null : value);
	return [
		{ label: 'Total XP', value: props.xp === null ? null : props.xp.toLocaleString(), sub: 'earned so far' },
		{ label: 'Rounds played', value: round(props.rounds), sub: props.roundsSub },
		{ label: 'Wins', value: round(props.wins), sub: `${props.losses} ${props.losses === 1 ? 'loss' : 'losses'}` },
		{
			label: 'Win rate',
			value: round(decided ? `${Math.round((props.wins / decided) * 100)}%` : '-'),
			sub: decided ? `of ${decided} finished ${decided === 1 ? 'round' : 'rounds'}` : 'no finished rounds yet',
		},
	];
});
</script>

<style scoped>
	.stat-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
		gap: 16px;
	}

	.stat-label {
		font-family: var(--font-mono);
		font-size: 11.5px;
		font-weight: 700;
		letter-spacing: 0.08em;
		text-transform: uppercase;
		color: var(--color-outline);
	}

	.stat-value {
		font-size: var(--text-2xl);
		font-weight: 800;
		margin: 8px 0 2px;
		letter-spacing: -0.02em;
	}

	.value-skeleton {
		margin: 10px 0 8px;
	}

	.stat-sub {
		margin: 0;
		font-size: 12.5px;
		color: var(--color-outline);
	}
</style>
