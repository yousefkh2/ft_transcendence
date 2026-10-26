<template>
	<div class="friends">
		<div class="friends-toolbar">
			<div>
				<h1>Friends</h1>
				<p class="subtitle">Play with people you know.</p>
			</div>
			<form class="add-form" @submit.prevent="addFriend">
				<BaseInput v-model="newUsername" size="sm" class="add-input" type="text" placeholder="Add by username" aria-label="Username to add as friend" :disabled="sending" />
				<BaseButton type="submit" :disabled="sending || !newUsername.trim()">Add</BaseButton>
			</form>
		</div>

		<BaseAlert v-if="addError">{{ addError }}</BaseAlert>
		<BaseAlert v-if="addNotice" variant="success">{{ addNotice }}</BaseAlert>
		<BaseAlert v-if="loadError">{{ loadError }}</BaseAlert>

		<BaseCard v-if="requests.length" tag="section">
			<h2 class="panel-title">Pending requests</h2>
			<div class="request-list">
				<div v-for="request in requests" :key="request.id" class="row">
					<BaseAvatar :name="request.requesterUsername" />
					<div class="row-info">
						<div class="row-name">{{ request.requesterUsername }}</div>
						<div class="row-meta">Sent {{ formatDate(request.createdAt) }}</div>
					</div>
					<div class="row-actions">
						<BaseButton size="sm" pill @click="respond(request.id, 'accept')">Accept</BaseButton>
						<BaseButton variant="outline" size="sm" pill @click="respond(request.id, 'decline')">Ignore</BaseButton>
					</div>
				</div>
			</div>
		</BaseCard>

		<BaseCard tag="section">
			<h2 class="panel-title">Your friends</h2>
			<div v-if="friends.length" class="friend-list">
				<div v-for="friend in friends" :key="friend.friendId" class="row">
					<BaseAvatar :name="friend.friendUsername" />
					<div class="row-info">
						<div class="row-name">{{ friend.friendUsername }}</div>
						<div class="row-meta">Friends since {{ formatDate(friend.createdAt) }}</div>
					</div>
				</div>
			</div>
			<p v-else class="empty-hint">No friends yet - add someone above.</p>
		</BaseCard>
	</div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue';
import {
	fetchFriends,
	fetchFriendRequests,
	sendFriendRequest,
	acceptFriendRequest,
	declineFriendRequest,
	type Friend,
	type FriendRequest,
} from '../../api/friends';
import BaseAlert from '../base/BaseAlert.vue';
import BaseAvatar from '../base/BaseAvatar.vue';
import BaseCard from '../base/BaseCard.vue';
import BaseInput from '../base/BaseInput.vue';
import BaseButton from '../base/BaseButton.vue';

const friends = ref<Friend[]>([]);
const requests = ref<FriendRequest[]>([]);
const loadError = ref('');

const newUsername = ref('');
const sending = ref(false);
const addError = ref('');
const addNotice = ref('');

function formatDate(iso: string) {
	return new Date(iso).toLocaleDateString(undefined, { month: 'short', day: 'numeric', year: 'numeric' });
}

async function loadFriendsAndRequests() {
	try {
		[friends.value, requests.value] = await Promise.all([fetchFriends(), fetchFriendRequests()]);
	} catch (err) {
		loadError.value = err instanceof Error ? err.message : 'Could not load your friends.';
	}
}

async function addFriend() {
	const username = newUsername.value.trim();
	if (!username) return;

	addError.value = '';
	addNotice.value = '';
	sending.value = true;
	try {
		await sendFriendRequest(username);
		addNotice.value = `Friend request sent to ${username}.`;
		newUsername.value = '';
	} catch (err) {
		addError.value = err instanceof Error ? err.message : 'Could not send that friend request.';
	} finally {
		sending.value = false;
	}
}

async function respond(id: string, action: 'accept' | 'decline') {
	try {
		if (action === 'accept') {
			await acceptFriendRequest(id);
		} else {
			await declineFriendRequest(id);
		}
		await loadFriendsAndRequests();
	} catch (err) {
		loadError.value = err instanceof Error ? err.message : 'Could not update that request.';
	}
}

onMounted(loadFriendsAndRequests);
</script>

<style scoped>
	.friends {
		display: grid;
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
		font-size: 17px;
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
		border: 1px solid rgba(118, 117, 134, 0.12);
	}

	.row-info {
		flex: 1;
		min-width: 0;
	}

	.row-name {
		font-size: 13.5px;
		font-weight: 700;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}

	.row-meta {
		font-size: 11.5px;
		color: var(--color-outline);
	}

	.row-actions {
		display: flex;
		gap: 8px;
	}

	.empty-hint {
		margin: 0;
		font-size: 13px;
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
