<template>
	<div class="overview">
		<section class="hero">
			<div class="hero-glow"></div>
			<div class="hero-copy">
				<BaseEyebrow variant="chip">Next up · Apartment Setup</BaseEyebrow>
				<h1>Ready for a round?</h1>
				<p>Jump into the live game and get matched with a partner for the next scenario.</p>
				<BaseButton to="/dev" variant="gradient" size="lg">Play now</BaseButton>
			</div>
		</section>

		<BaseAlert v-if="error" size="md">{{ error }}</BaseAlert>

		<section v-else class="stat-grid">
			<BaseCard size="sm">
				<span class="stat-label">Total XP</span>
				<div class="stat-value">{{ xpDisplay }}</div>
			</BaseCard>
			<BaseCard size="sm">
				<span class="stat-label">Rounds played</span>
				<div class="stat-value">{{ matches.length }}</div>
				<p class="stat-sub">{{ roundsThisWeek }} this week</p>
			</BaseCard>
		</section>

		<section class="panels">
			<BaseCard>
				<div class="panel-header">
					<h2>Recent rounds</h2>
				</div>
				<div v-if="recentRounds.length" class="round-list">
					<div v-for="round in recentRounds" :key="round.sessionId" class="round-row">
						<div class="round-badge" :class="round.badgeClass">{{ round.badgeLabel }}</div>
						<div class="round-info">
							<div class="round-mode">{{ round.modeLabel }}</div>
							<div class="round-when">{{ round.whenLabel }}</div>
						</div>
					</div>
				</div>
				<p v-else class="empty-hint">No rounds played yet - jump into a game to see them here.</p>
			</BaseCard>

			<BaseCard>
				<div class="panel-header">
					<h2>Friends</h2>
					<button type="button" class="link-btn" @click="emit('switch-tab', 'friends')">See all</button>
				</div>
				<div v-if="friends.length" class="friend-list">
					<div v-for="friend in friendsPreview" :key="friend.friendId" class="friend-row">
						<BaseAvatar :name="friend.friendUsername" />
						<div class="friend-info">
							<div class="friend-name">{{ friend.friendUsername }}</div>
						</div>
					</div>
				</div>
				<p v-else class="empty-hint">No friends yet - add someone from the Friends tab.</p>
			</BaseCard>
		</section>
	</div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { fetchFriends, type Friend } from '../../api/friends';
import { fetchMatchHistory, type MatchHistoryEntry, type Profile } from '../../api/users';
import BaseAlert from '../base/BaseAlert.vue';
import BaseAvatar from '../base/BaseAvatar.vue';
import BaseCard from '../base/BaseCard.vue';
import BaseButton from '../base/BaseButton.vue';
import BaseEyebrow from '../base/BaseEyebrow.vue';

const props = defineProps<{ profile: Profile | null }>();
const emit = defineEmits<{ (e: 'switch-tab', tab: 'friends'): void }>();

const matches = ref<MatchHistoryEntry[]>([]);
const friends = ref<Friend[]>([]);
const error = ref('');

const xpDisplay = computed(() => (props.profile ? props.profile.xp.toLocaleString() : '…'));

const roundsThisWeek = computed(() => {
	const weekAgo = Date.now() - 7 * 24 * 60 * 60 * 1000;
	return matches.value.filter((m) => new Date(m.startedAt).getTime() >= weekAgo).length;
});

function formatMode(mode: string) {
	return mode.replace(/_/g, ' ').replace(/\b\w/g, (c) => c.toUpperCase());
}

function formatWhen(iso: string) {
	const diffDays = Math.floor((Date.now() - new Date(iso).getTime()) / (24 * 60 * 60 * 1000));
	if (diffDays <= 0) return 'Today';
	if (diffDays === 1) return 'Yesterday';
	return `${diffDays} days ago`;
}

const recentRounds = computed(() =>
	matches.value.slice(0, 4).map((m) => ({
		sessionId: m.sessionId,
		modeLabel: formatMode(m.gameMode),
		whenLabel: formatWhen(m.startedAt),
		badgeLabel: m.isWinner === true ? 'Won' : m.isWinner === false ? 'Lost' : '—',
		badgeClass: m.isWinner === true ? 'won' : m.isWinner === false ? 'lost' : 'pending',
	})),
);

const friendsPreview = computed(() => friends.value.slice(0, 3));

onMounted(async () => {
	try {
		[matches.value, friends.value] = await Promise.all([fetchMatchHistory(), fetchFriends()]);
	} catch (err) {
		error.value = err instanceof Error ? err.message : 'Could not load your stats.';
	}
});
</script>

<style scoped>
	.overview {
		display: grid;
		gap: 22px;
	}

	.hero {
		position: relative;
		overflow: hidden;
		background: linear-gradient(135deg, var(--color-surface-inverse) 0%, #26295c 100%);
		border-radius: 24px;
		padding: 40px;
	}

	.hero-glow {
		position: absolute;
		top: -45%;
		right: 8%;
		width: 40%;
		aspect-ratio: 1;
		border-radius: 9999px;
		background: radial-gradient(circle, rgba(96, 99, 238, 0.45) 0%, transparent 70%);
		pointer-events: none;
	}

	.hero-copy {
		position: relative;
		max-width: 460px;
	}

	.hero-copy h1 {
		color: var(--color-on-primary);
		font-size: 32px;
		line-height: 0.95;
		font-weight: 800;
		margin: 0 0 10px;
		letter-spacing: -0.02em;
	}

	.hero-copy p {
		color: #c7c9e0;
		font-size: 15px;
		line-height: 1.6;
		margin: 0 0 24px;
	}

	.stat-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
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
		grid-template-columns: minmax(300px, 1.6fr) minmax(260px, 1fr);
		gap: 22px;
	}

	.panel-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 18px;
	}

	.panel-header h2 {
		font-size: 17px;
		font-weight: 800;
		margin: 0;
	}

	.link-btn {
		background: none;
		border: none;
		color: var(--color-primary);
		font-weight: 700;
		font-size: 13px;
		cursor: pointer;
		font-family: inherit;
	}

	.round-list,
	.friend-list {
		display: grid;
		gap: 10px;
	}

	.round-row,
	.friend-row {
		display: flex;
		align-items: center;
		gap: 14px;
		padding: 13px 15px;
		border-radius: 13px;
		background: var(--color-surface);
		border: 1px solid rgba(118, 117, 134, 0.12);
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
		font-size: 11px;
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
		font-size: 14px;
		font-weight: 700;
	}

	.round-when {
		font-size: 12.5px;
		color: var(--color-outline);
	}

	.friend-info {
		flex: 1;
		min-width: 0;
	}

	.friend-name {
		font-size: 13.5px;
		font-weight: 700;
	}

	.empty-hint {
		margin: 0;
		font-size: 13px;
		color: var(--color-outline);
	}
</style>
