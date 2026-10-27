<template>
	<div class="profile">
		<template v-if="notFound">
			<BaseAlert size="md">There is no player called "{{ username }}".</BaseAlert>
			<BaseButton :to="{ name: 'friends' }" variant="outline" size="sm" pill class="back">Back to friends</BaseButton>
		</template>

		<template v-else>
			<ProfileHead :name="user?.username ?? username" :status="user ? (user.online ? 'online' : 'offline') : undefined">
				<p v-if="user" class="presence" :class="{ online: user.online }">{{ formatLastSeen(user.online, user.lastSeenAt) }}</p>
				<p v-if="user" class="member-since">
					Learning {{ languageName(user.learningLanguage) }} · Member since {{ formatDate(user.createdAt) }}
				</p>
				<template #actions>
					<span v-if="relation === 'friend'" class="badge">Friends</span>
					<span v-else-if="relation === 'sent'" class="badge">Request sent</span>
					<BaseButton v-else-if="relation === 'incoming'" :to="{ name: 'friends' }" size="sm" pill>
						Answer friend request
					</BaseButton>
					<BaseButton v-else-if="relation === 'none'" size="sm" pill :disabled="sending" @click="addFriend">
						Add friend
					</BaseButton>
				</template>
			</ProfileHead>

			<BaseAlert v-if="loadError" size="md">This profile could not be loaded right now.</BaseAlert>
			<BaseAlert v-if="actionError">{{ actionError }}</BaseAlert>

			<ProfileStats
				v-if="user"
				:xp="user.xp"
				:rounds="user.roundsPlayed"
				:wins="user.wins"
				:losses="user.losses"
				rounds-sub="completed rounds"
			/>
		</template>
	</div>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { ApiError } from '../../api/client';
import { fetchFriendRequests, fetchFriends, sendFriendRequest, type Friend, type FriendRequest } from '../../api/friends';
import { fetchPublicProfile, type Profile, type PublicProfile } from '../../api/users';
import { formatDate, formatLastSeen } from '../../utils/format';
import { languageName } from '../../constants/languages';
import BaseAlert from '../base/BaseAlert.vue';
import BaseButton from '../base/BaseButton.vue';
import ProfileHead from './ProfileHead.vue';
import ProfileStats from './ProfileStats.vue';

const props = defineProps<{ profile: Profile | null }>();

const route = useRoute();
const router = useRouter();

const username = computed(() => String(route.params.username ?? ''));
const user = ref<PublicProfile | null>(null);
const friends = ref<Friend[]>([]);
const requests = ref<FriendRequest[]>([]);
const notFound = ref(false);
const loadError = ref(false);
const requestSent = ref(false);
const sending = ref(false);
const actionError = ref('');

const relation = computed(() => {
	if (!user.value) return 'unknown';
	if (friends.value.some((f) => f.friendUsername === user.value?.username)) return 'friend';
	if (requests.value.some((r) => r.requesterUsername === user.value?.username)) return 'incoming';
	if (requestSent.value) return 'sent';
	return 'none';
});

async function load(name: string) {
	user.value = null;
	notFound.value = false;
	loadError.value = false;
	requestSent.value = false;
	actionError.value = '';
	try {
		[user.value, friends.value, requests.value] = await Promise.all([
			fetchPublicProfile(name),
			fetchFriends(),
			fetchFriendRequests(),
		]);
	} catch (err) {
		if (err instanceof ApiError && err.status === 404) notFound.value = true;
		else loadError.value = true;
	}
}

async function addFriend() {
	if (!user.value) return;
	sending.value = true;
	actionError.value = '';
	try {
		await sendFriendRequest(user.value.username);
		requestSent.value = true;
	} catch (err) {
		actionError.value = err instanceof Error ? err.message : 'Could not send that friend request.';
	} finally {
		sending.value = false;
	}
}

// Your own name leads to your own profile, which also shows private data like your email.
watch(
	[username, () => props.profile?.username],
	([name, own]) => {
		if (own && name === own) router.replace({ name: 'profile' });
		else if (name && name !== user.value?.username) load(name);
	},
	{ immediate: true },
);
</script>

<style scoped>
	.profile {
		display: grid;
		gap: 22px;
	}

	.presence {
		margin: 0 0 2px;
		font-size: 13px;
		font-weight: 700;
		color: var(--color-outline);
	}

	.presence.online {
		color: var(--color-on-success-container);
	}

	.member-since {
		margin: 0;
		font-size: 13px;
		color: var(--color-outline);
	}

	.badge {
		display: inline-flex;
		padding: 6px 12px;
		border-radius: var(--radius-pill);
		background: rgba(34, 197, 94, 0.12);
		color: var(--color-on-success-container);
		font-size: 12.5px;
		font-weight: 700;
	}

	.back {
		justify-self: start;
	}
</style>
