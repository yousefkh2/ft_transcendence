<template>
	<div class="friends">
		<div class="friends-toolbar">
			<div>
				<h1>{{ t('friends.title') }}</h1>
				<p class="subtitle">{{ t('friends.subtitle') }}</p>
			</div>
			<form class="add-form" @submit.prevent="addFriend">
				<BaseInput v-model="newUsername" size="sm" class="add-input" :error="addError" type="text" :placeholder="t('friends.addPlaceholder')" :aria-label="t('friends.addLabel')" :disabled="sending" />
				<BaseButton type="submit" :disabled="sending || !newUsername.trim()"><BaseIcon name="user-plus" :size="16" />{{ t('friends.add') }}</BaseButton>
			</form>
		</div>

		<BaseAlert v-if="addError">{{ addError }}</BaseAlert>
		<BaseAlert v-if="addNotice" variant="success">{{ addNotice }}</BaseAlert>
		<BaseAlert v-if="actionError">{{ actionError }}</BaseAlert>

		<BaseCard v-if="requests.length" tag="section">
			<h2 class="panel-title">{{ t('friends.pending') }}</h2>
			<div class="request-list">
				<div v-for="request in requests" :key="request.id" class="row">
					<BaseAvatar :name="request.requesterUsername" :src="apiAssetUrl(request.requesterAvatarUrl)" />
					<div class="row-info">
						<RouterLink :to="{ name: 'user-profile', params: { username: request.requesterUsername } }" class="row-name">{{ request.requesterUsername }}</RouterLink>
						<div class="row-meta">{{ t('friends.sent', { date: formatDate(request.createdAt) }) }}</div>
					</div>
					<div class="row-actions">
						<BaseButton size="sm" pill @click="respond(request.id, 'accept')">{{ t('friends.accept') }}</BaseButton>
						<BaseButton variant="outline" size="sm" pill @click="respond(request.id, 'decline')">{{ t('friends.ignore') }}</BaseButton>
					</div>
				</div>
			</div>
		</BaseCard>

		<BaseCard tag="section">
			<h2 class="panel-title">{{ t('friends.yours') }}</h2>
			<ListSkeleton v-if="status === 'loading'" :label="t('friends.loading')" avatar />
			<BaseAlert v-else-if="status === 'error'" @retry="loadFriendsAndRequests">
				{{ t('friends.loadError') }}
			</BaseAlert>
			<div v-else-if="friends.length" class="friend-list">
				<div v-for="friend in friends" :key="friend.friendId" class="row">
					<BaseAvatar :name="friend.friendUsername" :src="apiAssetUrl(friend.avatarUrl)" :status="friend.online ? 'online' : 'offline'" />
					<div class="row-info">
						<RouterLink :to="{ name: 'user-profile', params: { username: friend.friendUsername } }" class="row-name">{{ friend.friendUsername }}</RouterLink>
						<div class="row-meta">
							<span :class="{ online: friend.online }">{{ formatLastSeen(friend.online, friend.lastSeenAt) }}</span>
							· {{ t('friends.since', { date: formatDate(friend.createdAt) }) }}
						</div>
					</div>
					<div class="row-actions">
						<template v-if="confirmingId === friend.friendId">
							<BaseButton variant="danger" size="sm" pill :disabled="removing" @click="remove(friend)">
								{{ t('friends.removeConfirm') }}
							</BaseButton>
							<BaseButton variant="outline" size="sm" pill :disabled="removing" @click="confirmingId = ''">
								{{ t('common.cancel') }}
							</BaseButton>
						</template>
						<BaseButton
							v-else
							variant="ghost"
							size="sm"
							pill
							:aria-label="t('friends.removeLabel', { name: friend.friendUsername })"
							@click="confirmingId = friend.friendId"
						>
							{{ t('friends.remove') }}
						</BaseButton>
					</div>
				</div>
			</div>
			<p v-else class="empty-hint">{{ t('friends.none') }}</p>
		</BaseCard>
	</div>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { apiAssetUrl } from '../../api/client';
import {
	fetchFriends,
	fetchFriendRequests,
	sendFriendRequest,
	acceptFriendRequest,
	declineFriendRequest,
	removeFriend,
	type Friend,
	type FriendRequest,
} from '../../api/friends';
import type { Profile } from '../../api/users';
import { PRESENCE_INTERVAL_MS } from '../../composables/useHeartbeat';
import { usePolling } from '../../composables/usePolling';
import { useLoader } from '../../composables/useLoader';
import { formatDate, formatLastSeen } from '../../utils/format';
import { sortByPresence } from '../../utils/friends';
import BaseAlert from '../base/BaseAlert.vue';
import BaseAvatar from '../base/BaseAvatar.vue';
import BaseCard from '../base/BaseCard.vue';
import BaseInput from '../base/BaseInput.vue';
import BaseButton from '../base/BaseButton.vue';
import BaseIcon from '../base/BaseIcon.vue';
import { t } from '../../i18n';
import ListSkeleton from './ListSkeleton.vue';

const props = defineProps<{ profile: Profile | null }>();

// Friends and requests load together: both lists are needed to check a new friend request.
const {
	data,
	status,
	load: loadFriendsAndRequests,
} = useLoader(
	async () => {
		const [friendList, requests] = await Promise.all([fetchFriends(), fetchFriendRequests()]);
		return { friends: sortByPresence(friendList), requests };
	},
	{ friends: [] as Friend[], requests: [] as FriendRequest[] },
);
const friends = computed(() => data.value.friends);
const requests = computed(() => data.value.requests);
const actionError = ref('');

const newUsername = ref('');
const sending = ref(false);
const addError = ref('');
const addNotice = ref('');

watch(newUsername, () => {
	addError.value = '';
});

// Catches requests the backend would reject anyway, so they never leave the browser.
function checkNewFriend(username: string): string {
	if (username === props.profile?.username) return t('friends.self');
	if (friends.value.some((f) => f.friendUsername === username)) return t('friends.already', { name: username });
	if (requests.value.some((r) => r.requesterUsername === username)) {
		return t('friends.incoming', { name: username });
	}
	return '';
}

async function addFriend() {
	const username = newUsername.value.trim();
	if (!username) return;

	addError.value = checkNewFriend(username);
	addNotice.value = '';
	if (addError.value) return;
	sending.value = true;
	try {
		await sendFriendRequest(username);
		addNotice.value = t('friends.requestSent', { name: username });
		newUsername.value = '';
	} catch (err) {
		addError.value = err instanceof Error ? err.message : t('friends.sendError');
	} finally {
		sending.value = false;
	}
}

const confirmingId = ref('');
const removing = ref(false);

async function remove(friend: Friend) {
	removing.value = true;
	addError.value = '';
	actionError.value = '';
	try {
		await removeFriend(friend.friendId);
		addNotice.value = t('friends.removed', { name: friend.friendUsername });
		confirmingId.value = '';
		await loadFriendsAndRequests();
	} catch (err) {
		actionError.value = err instanceof Error ? err.message : t('friends.removeError');
	} finally {
		removing.value = false;
	}
}

async function respond(id: string, action: 'accept' | 'decline') {
	actionError.value = '';
	try {
		if (action === 'accept') {
			await acceptFriendRequest(id);
		} else {
			await declineFriendRequest(id);
		}
		await loadFriendsAndRequests();
	} catch (err) {
		actionError.value = err instanceof Error ? err.message : t('friends.respondError');
	}
}

// Polling keeps the online status and new requests up to date.
usePolling(loadFriendsAndRequests, PRESENCE_INTERVAL_MS);
</script>

<style scoped>
	.friends {
		display: grid;
		/* minmax(0, ...) lets the column shrink below its content, so long text wraps instead of overflowing. */
		grid-template-columns: minmax(0, 1fr);
		gap: 22px;
	}

	.friends-toolbar {
		display: flex;
		justify-content: space-between;
		align-items: flex-end;
		gap: 20px;
		flex-wrap: wrap;
	}

	.friends-toolbar h1 {
		font-size: 26px;
		line-height: 0.95;
		font-weight: 800;
		margin: 0 0 6px;
		letter-spacing: -0.02em;
	}

	.subtitle {
		margin: 0;
		color: var(--color-outline);
		font-size: 14.5px;
	}

	.add-form {
		display: flex;
		gap: 10px;
	}

	.add-input {
		min-width: 200px;
	}

	.panel-title {
		font-size: var(--text-lg);
		font-weight: 800;
		margin: 0 0 16px;
	}

	.request-list,
	.friend-list {
		display: grid;
		gap: 10px;
	}

	.row {
		display: flex;
		align-items: center;
		gap: 12px;
		padding: 12px 14px;
		border-radius: 13px;
		background: var(--color-surface);
		border: 1px solid var(--color-border-subtle);
	}

	.row-info {
		flex: 1;
		min-width: 0;
	}

	.row-name {
		display: block;
		color: inherit;
		text-decoration: none;
		font-size: 13.5px;
		font-weight: 700;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}

	.row-name:hover {
		text-decoration: underline;
	}

	.row-meta {
		font-size: 11.5px;
		color: var(--color-outline);
	}

	.row-meta .online {
		color: var(--color-on-success-container);
		font-weight: 700;
	}

	.row-actions {
		display: flex;
		gap: 8px;
	}

	.empty-hint {
		margin: 0;
		font-size: var(--text-sm);
		color: var(--color-outline);
	}

	@media (max-width: 720px) {
		.add-form {
			width: 100%;
		}

		.add-input {
			flex: 1;
			min-width: 0;
		}
	}

	@media (max-width: 480px) {
		.row {
			flex-wrap: wrap;
		}

		.row-actions {
			width: 100%;
			justify-content: flex-end;
		}
	}
</style>
