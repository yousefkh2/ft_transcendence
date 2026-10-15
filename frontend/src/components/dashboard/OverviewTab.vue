<template>
	<div class="overview">
		<section class="hero">
			<div class="hero-glow"></div>
			<div class="hero-copy">
				<span class="hero-eyebrow">Next up · Apartment Setup</span>
				<h1>Ready for a round?</h1>
				<p>Jump into the live game and get matched with a partner for the next scenario.</p>
				<a href="/" class="play-btn">Play now</a>
			</div>
		</section>

		<p v-if="error" class="load-error">{{ error }}</p>

		<section v-else class="stat-grid">
			<div class="stat-card">
				<span class="stat-label">Total XP</span>
				<div class="stat-value">{{ xpDisplay }}</div>
			</div>
			<div class="stat-card">
				<span class="stat-label">Rounds played</span>
				<div class="stat-value">{{ matches.length }}</div>
				<p class="stat-sub">{{ roundsThisWeek }} this week</p>
			</div>
		</section>

		<section class="panels">
			<div class="panel">
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
			</div>

			<div class="panel">
				<div class="panel-header">
					<h2>Friends</h2>
					<button type="button" class="link-btn" @click="emit('switch-tab', 'friends')">See all</button>
				</div>
				<div v-if="friends.length" class="friend-list">
					<div v-for="friend in friendsPreview" :key="friend.friendId" class="friend-row">
						<div class="friend-avatar">{{ friend.friendUsername[0]?.toUpperCase() }}</div>
						<div class="friend-info">
							<div class="friend-name">{{ friend.friendUsername }}</div>
						</div>
					</div>
				</div>
				<p v-else class="empty-hint">No friends yet - add someone from the Friends tab.</p>
			</div>
		</section>
	</div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { fetchMatchHistory, fetchFriends, type MatchHistoryEntry, type Friend, type Profile } from '../../api/dashboard';

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
		background: linear-gradient(135deg, #111c2d 0%, #26295c 100%);
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

	.hero-eyebrow {
		display: inline-flex;
		padding: 5px 12px;
		border-radius: 9999px;
		background: rgba(255, 255, 255, 0.14);
		color: #fff;
		font-family: 'Space Grotesk', monospace;
		font-size: 11px;
		font-weight: 700;
		letter-spacing: 0.1em;
		text-transform: uppercase;
		margin-bottom: 14px;
	}

	.hero-copy h1 {
		color: #fff;
		font-size: 32px;
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

	.play-btn {
		display: inline-flex;
		height: 50px;
		align-items: center;
		padding: 0 28px;
		border-radius: 12px;
		background: linear-gradient(135deg, #4648d4, #6063ee);
		color: #fff;
		font-weight: 800;
		font-size: 15px;
		box-shadow: 0 12px 28px rgba(70, 72, 212, 0.4);
		text-decoration: none;
	}

	.stat-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
		gap: 16px;
	}

	.stat-card {
		background: #fff;
		border: 1px solid rgba(118, 117, 134, 0.18);
		border-radius: 16px;
		padding: 22px;
	}

	.stat-label {
		font-family: 'Space Grotesk', monospace;
		font-size: 11.5px;
		font-weight: 700;
		letter-spacing: 0.08em;
		text-transform: uppercase;
		color: #767586;
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
		color: #767586;
	}

	.panels {
		display: grid;
		grid-template-columns: minmax(300px, 1.6fr) minmax(260px, 1fr);
		gap: 22px;
	}

	.panel {
		background: #fff;
		border: 1px solid rgba(118, 117, 134, 0.18);
		border-radius: 18px;
		padding: 26px;
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
		color: #4648d4;
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
		background: #f9f9ff;
		border: 1px solid rgba(118, 117, 134, 0.12);
	}

	.round-badge {
		width: 36px;
		height: 36px;
		flex: none;
		border-radius: 10px;
		display: grid;
		place-items: center;
		font-family: 'Space Grotesk', monospace;
		font-weight: 700;
		font-size: 11px;
		color: #fff;
	}

	.round-badge.won {
		background: #22c55e;
	}

	.round-badge.lost {
		background: #ff5d8f;
	}

	.round-badge.pending {
		background: #767586;
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
		color: #767586;
	}

	.friend-avatar {
		width: 34px;
		height: 34px;
		flex: none;
		border-radius: 9999px;
		background: #4648d4;
		color: #fff;
		display: grid;
		place-items: center;
		font-weight: 800;
		font-size: 13px;
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
		color: #767586;
	}

	.load-error {
		padding: 16px 20px;
		border-radius: 12px;
		background: rgba(209, 53, 43, 0.08);
		color: #d1352b;
		font-weight: 600;
		font-size: 14px;
	}
</style>
