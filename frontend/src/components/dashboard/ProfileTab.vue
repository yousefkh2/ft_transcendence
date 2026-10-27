<template>
	<div class="profile">
		<ProfileHead :name="profile?.username" :src="apiAssetUrl(profile?.avatarUrl)">
			<p v-if="profile" class="email">
				{{ profile.email }}
				<BaseBadge icon="lock">{{ t('profile.onlyYou') }}</BaseBadge>
			</p>
			<p v-if="profile" class="member-since">
				{{ t('profile.learning', { language: languageName(profile.learningLanguage), date: formatDate(profile.createdAt) }) }}
			</p>
			<template #actions>
				<AvatarEditor :profile="profile" @updated="emit('profile-updated', $event)" />
			</template>
		</ProfileHead>

		<ProfileStats
			v-if="historyStatus !== 'error'"
			:xp="profile?.xp ?? null"
			:rounds="matches.length"
			:wins="wins"
			:losses="losses"
			:rounds-sub="t('overview.thisWeek', { n: countSince(matches, 7) })"
			:loading="historyStatus === 'loading'"
		/>

		<section class="panels">
			<BaseCard>
				<div class="panel-header">
					<h2>{{ t('profile.matchHistory') }}</h2>
					<span v-if="historyStatus === 'ready'" class="count">
						{{ t('profile.rounds', { n: matches.length }) }}
					</span>
				</div>
				<ListSkeleton v-if="historyStatus === 'loading'" :label="t('profile.loadingHistory')" />
				<BaseAlert v-else-if="historyStatus === 'error'" @retry="loadHistory">
					{{ t('profile.historyError') }}
				</BaseAlert>
				<RoundList v-else-if="matches.length" :matches="matches" show-role />
				<p v-else class="empty-hint">{{ t('overview.noRounds') }}</p>
			</BaseCard>

			<BaseCard class="friends-card">
				<h2>{{ t('profile.friends') }}</h2>
				<div v-if="friendsStatus === 'loading'" class="friends-skeleton" role="status">
					<span class="visually-hidden">{{ t('friends.loading') }}</span>
					<BaseSkeleton shape="block" width="40px" height="30px" />
					<BaseSkeleton width="50px" height="0.75em" />
				</div>
				<BaseAlert v-else-if="friendsStatus === 'error'" @retry="loadFriends">
					{{ t('profile.friendsError') }}
				</BaseAlert>
				<template v-else>
					<div class="friends-count">{{ friends.length }}</div>
					<p class="stat-sub">{{ t('profile.friendCount', { n: friends.length }) }}</p>
				</template>
				<BaseButton :to="{ name: 'friends' }" variant="outline" size="sm" pill>{{ t('profile.manageFriends') }}</BaseButton>
			</BaseCard>
		</section>

		<EditProfileCard v-if="profile" :profile="profile" @updated="emit('profile-updated', $event)" />

		<DeleteAccountCard />
	</div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue';
import { apiAssetUrl } from '../../api/client';
import { fetchFriends, type Friend } from '../../api/friends';
import { fetchMatchHistory, type MatchHistoryEntry, type Profile } from '../../api/users';
import { useLoader } from '../../composables/useLoader';
import { countSince, roundResult } from '../../utils/matches';
import { formatDate } from '../../utils/format';
import { languageName } from '../../constants/languages';
import BaseAlert from '../base/BaseAlert.vue';
import BaseBadge from '../base/BaseBadge.vue';
import { t } from '../../i18n';
import BaseButton from '../base/BaseButton.vue';
import BaseCard from '../base/BaseCard.vue';
import BaseSkeleton from '../base/BaseSkeleton.vue';
import AvatarEditor from './AvatarEditor.vue';
import DeleteAccountCard from './DeleteAccountCard.vue';
import EditProfileCard from './EditProfileCard.vue';
import ListSkeleton from './ListSkeleton.vue';
import ProfileHead from './ProfileHead.vue';
import ProfileStats from './ProfileStats.vue';
import RoundList from './RoundList.vue';

defineProps<{ profile: Profile | null }>();
const emit = defineEmits<{ 'profile-updated': [profile: Profile] }>();

// Loaded independently, so one failing request does not hide the rest of the profile.
const { data: matches, status: historyStatus, load: loadHistory } = useLoader(fetchMatchHistory, [] as MatchHistoryEntry[]);
const { data: friends, status: friendsStatus, load: loadFriends } = useLoader(fetchFriends, [] as Friend[]);

const wins = computed(() => matches.value.filter((m) => roundResult(m) === 'won').length);
const losses = computed(() => matches.value.filter((m) => roundResult(m) === 'lost').length);

onMounted(() => {
	loadHistory();
	loadFriends();
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
		font-size: var(--text-md);
		color: var(--color-on-surface-variant);
		overflow-wrap: anywhere;
	}

	.member-since {
		margin: 6px 0 0;
		font-size: 12.5px;
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
		font-size: var(--text-lg);
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
		font-size: var(--text-2xl);
		font-weight: 800;
		letter-spacing: -0.02em;
	}

	.friends-skeleton {
		display: grid;
		gap: 8px;
		margin: 14px 0 18px;
	}

	.friends-card .stat-sub {
		margin-bottom: 16px;
	}

	.friends-card :deep(.alert) {
		margin: 12px 0 16px;
	}

	.empty-hint {
		margin: 0;
		font-size: var(--text-sm);
		color: var(--color-outline);
	}

	@media (max-width: 720px) {
		.panels {
			grid-template-columns: 1fr;
		}
	}
</style>
