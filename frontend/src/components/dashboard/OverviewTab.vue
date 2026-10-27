<template>
	<div class="overview">
		<section class="hero">
			<div class="hero-glow"></div>
			<div class="hero-copy">
				<BaseEyebrow variant="chip">Next up · Apartment Setup</BaseEyebrow>
				<h1>Ready for a round?</h1>
				<p>Jump into the live game and get matched with a partner for the next scenario.</p>
				<BaseButton to="/play" variant="gradient" size="lg">Play now</BaseButton>
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
				<RoundList v-if="recentRounds.length" :matches="recentRounds" />
				<p v-else class="empty-hint">No rounds played yet - jump into a game to see them here.</p>
			</BaseCard>

			<BaseCard>
				<div class="panel-header">
					<h2>Friends</h2>
					<RouterLink :to="{ name: 'friends' }" class="link-btn">See all</RouterLink>
				</div>
				<div v-if="friends.length" class="friend-list">
					<div v-for="friend in friendsPreview" :key="friend.friendId" class="friend-row">
						<BaseAvatar :name="friend.friendUsername" :status="friend.online ? 'online' : 'offline'" />
						<div class="friend-info">
							<RouterLink :to="{ name: 'user-profile', params: { username: friend.friendUsername } }" class="friend-name">{{ friend.friendUsername }}</RouterLink>
							<div class="friend-meta" :class="{ online: friend.online }">{{ formatLastSeen(friend.online, friend.lastSeenAt) }}</div>
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
import RoundList from './RoundList.vue';
import { countSince } from '../../utils/matches';
import { formatLastSeen } from '../../utils/format';
import { sortByPresence } from '../../utils/friends';
import { PRESENCE_INTERVAL_MS } from '../../composables/useHeartbeat';
import { usePolling } from '../../composables/usePolling';

const props = defineProps<{ profile: Profile | null }>();

const matches = ref<MatchHistoryEntry[]>([]);
const friends = ref<Friend[]>([]);
const error = ref('');

const xpDisplay = computed(() => (props.profile ? props.profile.xp.toLocaleString() : '…'));

const roundsThisWeek = computed(() => countSince(matches.value, 7));

const recentRounds = computed(() => matches.value.slice(0, 4));

const friendsPreview = computed(() => sortByPresence(friends.value).slice(0, 3));

onMounted(async () => {
	try {
		matches.value = await fetchMatchHistory();
	} catch (err) {
		error.value = err instanceof Error ? err.message : 'Could not load your stats.';
	}
});

usePolling(async () => {
	try {
		friends.value = await fetchFriends();
	} catch (err) {
		error.value = err instanceof Error ? err.message : 'Could not load your friends.';
	}
}, PRESENCE_INTERVAL_MS);
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
		text-decoration: none;
		padding: 1px 6px;
		background: none;
		border: none;
		color: var(--color-primary);
		font-weight: 700;
		font-size: 13px;
		cursor: pointer;
		font-family: inherit;
	}

	.friend-list {
		display: grid;
		gap: 10px;
	}

	.friend-row {
		display: flex;
		align-items: center;
		gap: 14px;
		padding: 13px 15px;
		border-radius: 13px;
		background: var(--color-surface);
		border: 1px solid rgba(118, 117, 134, 0.12);
	}

	.friend-info {
		flex: 1;
		min-width: 0;
	}

	.friend-name {
		display: block;
		color: inherit;
		text-decoration: none;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
		font-size: 13.5px;
		font-weight: 700;
	}

	.friend-name:hover {
		text-decoration: underline;
	}

	.friend-meta {
		font-size: 11.5px;
		color: var(--color-outline);
	}

	.friend-meta.online {
		color: var(--color-on-success-container);
		font-weight: 700;
	}

	.empty-hint {
		margin: 0;
		font-size: 13px;
		color: var(--color-outline);
	}

	@media (max-width: 720px) {
		.hero {
			padding: 28px 22px;
		}

		.hero-copy h1 {
			font-size: 26px;
		}

		.panels {
			grid-template-columns: 1fr;
		}
	}
</style>
