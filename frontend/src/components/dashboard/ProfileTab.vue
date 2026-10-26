<template>
	<div class="profile">
		<BaseCard class="profile-head">
			<BaseAvatar :name="profile?.username" size="lg" variant="gradient" />
			<div class="identity">
				<h1>{{ profile?.username ?? '…' }}</h1>
				<p v-if="profile" class="email">
					{{ profile.email }}
					<span class="private">Only visible to you</span>
				</p>
			</div>
		</BaseCard>

		<section v-if="!historyError" class="stat-grid" aria-label="Your stats">
			<BaseCard v-for="stat in stats" :key="stat.label" size="sm">
				<span class="stat-label">{{ stat.label }}</span>
				<div class="stat-value">{{ stat.value }}</div>
				<p class="stat-sub">{{ stat.sub }}</p>
			</BaseCard>
		</section>

		<section class="panels">
			<BaseCard>
				<div class="panel-header">
					<h2>Match history</h2>
					<span v-if="!historyError" class="count">
						{{ matches.length }} {{ matches.length === 1 ? 'round' : 'rounds' }}
					</span>
				</div>
				<BaseAlert v-if="historyError">Your match history could not be loaded right now.</BaseAlert>
				<RoundList v-else-if="matches.length" :matches="matches" show-role />
				<p v-else class="empty-hint">No rounds played yet - jump into a game to see them here.</p>
			</BaseCard>

			<BaseCard class="friends-card">
				<h2>Friends</h2>
				<BaseAlert v-if="friendsError">Your friends could not be loaded right now.</BaseAlert>
				<template v-else>
					<div class="friends-count">{{ friends.length }}</div>
					<p class="stat-sub">{{ friends.length === 1 ? 'friend' : 'friends' }}</p>
				</template>
				<BaseButton :to="{ name: 'friends' }" variant="outline" size="sm" pill>Manage friends</BaseButton>
			</BaseCard>
		</section>
	</div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { fetchFriends, type Friend } from '../../api/friends';
import { fetchMatchHistory, type MatchHistoryEntry, type Profile } from '../../api/users';
import { countSince, roundResult } from '../../utils/matches';
import BaseAlert from '../base/BaseAlert.vue';
import BaseAvatar from '../base/BaseAvatar.vue';
import BaseButton from '../base/BaseButton.vue';
import BaseCard from '../base/BaseCard.vue';
import RoundList from './RoundList.vue';

const props = defineProps<{ profile: Profile | null }>();

const matches = ref<MatchHistoryEntry[]>([]);
const friends = ref<Friend[]>([]);
const historyError = ref(false);
const friendsError = ref(false);

const xpDisplay = computed(() => (props.profile ? props.profile.xp.toLocaleString() : '…'));

const stats = computed(() => {
	const wins = matches.value.filter((m) => roundResult(m) === 'won').length;
	const losses = matches.value.filter((m) => roundResult(m) === 'lost').length;
	const decided = wins + losses;
	return [
		{ label: 'Total XP', value: xpDisplay.value, sub: 'earned so far' },
		{ label: 'Rounds played', value: matches.value.length, sub: `${countSince(matches.value, 7)} this week` },
		{ label: 'Wins', value: wins, sub: `${losses} ${losses === 1 ? 'loss' : 'losses'}` },
		{
			label: 'Win rate',
			value: decided ? `${Math.round((wins / decided) * 100)}%` : '-',
			sub: decided ? `of ${decided} finished ${decided === 1 ? 'round' : 'rounds'}` : 'no finished rounds yet',
		},
	];
});

// Loaded independently, so one failing request does not hide the rest of the profile.
onMounted(async () => {
	const [history, friendList] = await Promise.allSettled([fetchMatchHistory(), fetchFriends()]);
	if (history.status === 'fulfilled') matches.value = history.value;
	else historyError.value = true;
	if (friendList.status === 'fulfilled') friends.value = friendList.value;
	else friendsError.value = true;
});
</script>

<style scoped>
	.profile {
		display: grid;
		gap: 22px;
	}

	.profile-head {
		display: flex;
		align-items: center;
		gap: 20px;
	}

	.identity {
		flex: 1;
		min-width: 0;
	}

	.identity h1 {
		margin: 0 0 6px;
		font-size: 26px;
		font-weight: 800;
		line-height: 1.1;
		letter-spacing: -0.02em;
		overflow: hidden;
		text-overflow: ellipsis;
	}

	.email {
		display: flex;
		flex-wrap: wrap;
		align-items: center;
		gap: 8px;
		margin: 0;
		font-size: 14px;
		color: var(--color-on-surface-variant);
		overflow-wrap: anywhere;
	}

	.private {
		padding: 2px 8px;
		border-radius: var(--radius-pill);
		background: rgba(118, 117, 134, 0.1);
		font-size: 11.5px;
		font-weight: 600;
		color: var(--color-outline);
	}

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
		font-size: 30px;
		font-weight: 800;
		margin: 8px 0 2px;
		letter-spacing: -0.02em;
	}

	.stat-sub {
		margin: 0;
		font-size: 12.5px;
		color: var(--color-outline);
	}

	.panels {
		display: grid;
		grid-template-columns: minmax(300px, 1.6fr) minmax(220px, 1fr);
		gap: 22px;
		align-items: start;
	}

	.panel-header {
		display: flex;
		justify-content: space-between;
		align-items: baseline;
		margin-bottom: 18px;
	}

	h2 {
		margin: 0;
		font-size: 17px;
		font-weight: 800;
	}

	.count {
		font-size: 12.5px;
		color: var(--color-outline);
	}

	.friends-card {
		display: flex;
		flex-direction: column;
		align-items: flex-start;
	}

	.friends-count {
		margin: 12px 0 2px;
		font-size: 30px;
		font-weight: 800;
		letter-spacing: -0.02em;
	}

	.friends-card .stat-sub {
		margin-bottom: 16px;
	}

	.friends-card :deep(.alert) {
		margin: 12px 0 16px;
	}

	.empty-hint {
		margin: 0;
		font-size: 13px;
		color: var(--color-outline);
	}

	@media (max-width: 720px) {
		.profile-head {
			flex-wrap: wrap;
		}

		.identity h1 {
			font-size: 22px;
		}

		.panels {
			grid-template-columns: 1fr;
		}
	}
</style>
