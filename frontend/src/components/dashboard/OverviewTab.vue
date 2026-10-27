<template>
	<div class="overview">
		<section class="hero">
			<div class="hero-glow"></div>
			<div class="hero-copy">
				<BaseEyebrow variant="chip">Next up · Apartment Setup</BaseEyebrow>
				<h1>Ready for a round?</h1>
				<p>Jump into the live game and get matched with a partner for the next scenario.</p>
				<BaseButton to="/play" variant="gradient" size="lg"><BaseIcon name="play" />Play now</BaseButton>
			</div>
		</section>

		<section class="stat-grid" aria-label="Stats">
			<BaseCard size="sm">
				<span class="stat-label">Total XP</span>
				<div v-if="profile" class="stat-value">{{ profile.xp.toLocaleString() }}</div>
				<BaseSkeleton v-else class="value-skeleton" shape="block" width="70px" height="30px" />
			</BaseCard>
			<BaseCard size="sm">
				<span class="stat-label">Rounds played</span>
				<template v-if="matchesStatus === 'ready'">
					<div class="stat-value">{{ matches.length }}</div>
					<p class="stat-sub">{{ roundsThisWeek }} this week</p>
				</template>
				<template v-else-if="matchesStatus === 'error'">
					<div class="stat-value">-</div>
					<p class="stat-sub">not available right now</p>
				</template>
				<template v-else>
					<BaseSkeleton class="value-skeleton" shape="block" width="70px" height="30px" />
					<BaseSkeleton width="90px" height="0.75em" />
				</template>
			</BaseCard>
		</section>

		<section class="panels">
			<BaseCard>
				<div class="panel-header">
					<h2>Recent rounds</h2>
				</div>
				<ListSkeleton v-if="matchesStatus === 'loading'" label="Loading recent rounds" />
				<BaseAlert v-else-if="matchesStatus === 'error'" @retry="loadMatches">
					Your rounds could not be loaded right now.
				</BaseAlert>
				<RoundList v-else-if="recentRounds.length" :matches="recentRounds" />
				<p v-else class="empty-hint">No rounds played yet - jump into a game to see them here.</p>
			</BaseCard>

			<BaseCard>
				<div class="panel-header">
					<h2>Friends</h2>
					<RouterLink :to="{ name: 'friends' }" class="link-btn">See all</RouterLink>
				</div>
				<ListSkeleton v-if="friendsStatus === 'loading'" label="Loading friends" avatar />
				<BaseAlert v-else-if="friendsStatus === 'error'" @retry="loadFriends">
					Your friends could not be loaded right now.
				</BaseAlert>
				<div v-else-if="friends.length" class="friend-list">
					<div v-for="friend in friendsPreview" :key="friend.friendId" class="friend-row">
						<BaseAvatar :name="friend.friendUsername" :src="apiAssetUrl(friend.avatarUrl)" :status="friend.online ? 'online' : 'offline'" />
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
import { computed, onMounted } from 'vue';
import { apiAssetUrl } from '../../api/client';
import { fetchFriends, type Friend } from '../../api/friends';
import { fetchMatchHistory, type MatchHistoryEntry, type Profile } from '../../api/users';
import BaseAlert from '../base/BaseAlert.vue';
import BaseAvatar from '../base/BaseAvatar.vue';
import BaseCard from '../base/BaseCard.vue';
import BaseButton from '../base/BaseButton.vue';
import BaseEyebrow from '../base/BaseEyebrow.vue';
import BaseIcon from '../base/BaseIcon.vue';
import BaseSkeleton from '../base/BaseSkeleton.vue';
import ListSkeleton from './ListSkeleton.vue';
import RoundList from './RoundList.vue';
import { countSince } from '../../utils/matches';
import { formatLastSeen } from '../../utils/format';
import { sortByPresence } from '../../utils/friends';
import { PRESENCE_INTERVAL_MS } from '../../composables/useHeartbeat';
import { usePolling } from '../../composables/usePolling';
import { useLoader } from '../../composables/useLoader';

defineProps<{ profile: Profile | null }>();

// Rounds and friends load independently, so one failing request does not hide the other card.
const { data: matches, status: matchesStatus, load: loadMatches } = useLoader(fetchMatchHistory, [] as MatchHistoryEntry[]);
const { data: friends, status: friendsStatus, load: loadFriends } = useLoader(fetchFriends, [] as Friend[]);

const roundsThisWeek = computed(() => countSince(matches.value, 7));

const recentRounds = computed(() => matches.value.slice(0, 4));

const friendsPreview = computed(() => sortByPresence(friends.value).slice(0, 3));

onMounted(loadMatches);
usePolling(loadFriends, PRESENCE_INTERVAL_MS);
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
		font-size: var(--text-lg);
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
		font-size: var(--text-sm);
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
		border: 1px solid var(--color-border-subtle);
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
		font-size: var(--text-sm);
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
