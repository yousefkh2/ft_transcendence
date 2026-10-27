<template>
	<div class="profile">
		<ProfileHead :name="profile?.username" :src="apiAssetUrl(profile?.avatarUrl)">
			<p v-if="profile" class="email">
				{{ profile.email }}
				<span class="private">Only visible to you</span>
			</p>
			<p v-if="profile" class="member-since">
				Learning {{ languageName(profile.learningLanguage) }} · Member since {{ formatDate(profile.createdAt) }}
			</p>
			<template #actions>
				<AvatarEditor :profile="profile" @updated="emit('profile-updated', $event)" />
			</template>
		</ProfileHead>

		<ProfileStats
			v-if="!historyError"
			:xp="profile?.xp ?? null"
			:rounds="matches.length"
			:wins="wins"
			:losses="losses"
			:rounds-sub="`${countSince(matches, 7)} this week`"
		/>

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

		<EditProfileCard :profile="profile" @updated="emit('profile-updated', $event)" />

		<DeleteAccountCard />
	</div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { apiAssetUrl } from '../../api/client';
import { fetchFriends, type Friend } from '../../api/friends';
import { fetchMatchHistory, type MatchHistoryEntry, type Profile } from '../../api/users';
import { countSince, roundResult } from '../../utils/matches';
import { formatDate } from '../../utils/format';
import { languageName } from '../../constants/languages';
import BaseAlert from '../base/BaseAlert.vue';
import BaseButton from '../base/BaseButton.vue';
import BaseCard from '../base/BaseCard.vue';
import AvatarEditor from './AvatarEditor.vue';
import DeleteAccountCard from './DeleteAccountCard.vue';
import EditProfileCard from './EditProfileCard.vue';
import ProfileHead from './ProfileHead.vue';
import ProfileStats from './ProfileStats.vue';
import RoundList from './RoundList.vue';

defineProps<{ profile: Profile | null }>();
const emit = defineEmits<{ 'profile-updated': [profile: Profile] }>();

const matches = ref<MatchHistoryEntry[]>([]);
const friends = ref<Friend[]>([]);
const historyError = ref(false);
const friendsError = ref(false);

const wins = computed(() => matches.value.filter((m) => roundResult(m) === 'won').length);
const losses = computed(() => matches.value.filter((m) => roundResult(m) === 'lost').length);

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

	.member-since {
		margin: 6px 0 0;
		font-size: 12.5px;
		color: var(--color-outline);
	}

	.private {
		padding: 2px 8px;
		border-radius: var(--radius-pill);
		background: rgba(118, 117, 134, 0.1);
		font-size: 11.5px;
		font-weight: 600;
		color: var(--color-outline);
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
		.panels {
			grid-template-columns: 1fr;
		}
	}
</style>
