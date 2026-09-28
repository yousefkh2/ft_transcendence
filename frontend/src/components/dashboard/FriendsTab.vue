<template>
	<div class="friends">
		<div class="friends-toolbar">
			<div>
				<h1>Friends</h1>
				<p class="subtitle">Play with people you know.</p>
			</div>
			<form class="add-form" @submit.prevent="addFriend">
				<input v-model="newUsername" type="text" placeholder="Add by username" :disabled="sending" />
				<button type="submit" :disabled="sending || !newUsername.trim()">Add</button>
			</form>
		</div>

		<p v-if="addError" class="form-error">{{ addError }}</p>
		<p v-if="addNotice" class="form-notice">{{ addNotice }}</p>
		<p v-if="loadError" class="load-error">{{ loadError }}</p>

		<section v-if="requests.length" class="panel">
			<h2>Pending requests</h2>
			<div class="request-list">
				<div v-for="request in requests" :key="request.id" class="row">
					<div class="avatar">{{ request.requesterUsername[0]?.toUpperCase() }}</div>
					<div class="row-info">
						<div class="row-name">{{ request.requesterUsername }}</div>
						<div class="row-meta">Sent {{ formatDate(request.createdAt) }}</div>
					</div>
					<div class="row-actions">
						<button type="button" class="accept-btn" @click="respond(request.id, 'accept')">Accept</button>
						<button type="button" class="decline-btn" @click="respond(request.id, 'decline')">Ignore</button>
					</div>
				</div>
			</div>
		</section>

		<section class="panel">
			<h2>Your friends</h2>
			<div v-if="friends.length" class="friend-list">
				<div v-for="friend in friends" :key="friend.friendId" class="row">
					<div class="avatar">{{ friend.friendUsername[0]?.toUpperCase() }}</div>
					<div class="row-info">
						<div class="row-name">{{ friend.friendUsername }}</div>
						<div class="row-meta">Friends since {{ formatDate(friend.createdAt) }}</div>
					</div>
				</div>
			</div>
			<p v-else class="empty-hint">No friends yet — add someone above.</p>
		</section>
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
} from '../../api/dashboard';

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
		font-weight: 800;
		margin: 0 0 6px;
		letter-spacing: -0.02em;
	}

	.subtitle {
		margin: 0;
		color: #767586;
		font-size: 14.5px;
	}

	.add-form {
		display: flex;
		gap: 10px;
	}

	.add-form input {
		height: 44px;
		border-radius: 11px;
		border: 1.5px solid rgba(118, 117, 134, 0.2);
		background: #fff;
		padding: 0 15px;
		font-size: 14px;
		font-family: inherit;
		color: #111c2d;
		outline: none;
		min-width: 200px;
	}

	.add-form button {
		height: 44px;
		padding: 0 20px;
		border: none;
		border-radius: 11px;
		background: #4648d4;
		color: #fff;
		font-weight: 800;
		font-size: 14px;
		cursor: pointer;
		font-family: inherit;
	}

	.add-form button:disabled {
		opacity: 0.6;
		cursor: not-allowed;
	}

	.form-error,
	.load-error {
		margin: 0;
		padding: 10px 14px;
		border-radius: 10px;
		background: rgba(209, 53, 43, 0.08);
		color: #d1352b;
		font-size: 13px;
		font-weight: 600;
	}

	.form-notice {
		margin: 0;
		padding: 10px 14px;
		border-radius: 10px;
		background: rgba(34, 197, 94, 0.1);
		color: #15803d;
		font-size: 13px;
		font-weight: 600;
	}

	.panel {
		background: #fff;
		border: 1px solid rgba(118, 117, 134, 0.18);
		border-radius: 18px;
		padding: 26px;
	}

	.panel h2 {
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
		background: #f9f9ff;
		border: 1px solid rgba(118, 117, 134, 0.12);
	}

	.avatar {
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

	.row-info {
		flex: 1;
		min-width: 0;
	}

	.row-name {
		font-size: 13.5px;
		font-weight: 700;
	}

	.row-meta {
		font-size: 11.5px;
		color: #767586;
	}

	.row-actions {
		display: flex;
		gap: 8px;
	}

	.accept-btn {
		border: none;
		background: #4648d4;
		color: #fff;
		font-weight: 700;
		font-size: 12px;
		padding: 7px 15px;
		border-radius: 9999px;
		cursor: pointer;
		font-family: inherit;
	}

	.decline-btn {
		border: 1.5px solid rgba(118, 117, 134, 0.3);
		background: none;
		color: #767586;
		font-weight: 700;
		font-size: 12px;
		padding: 7px 14px;
		border-radius: 9999px;
		cursor: pointer;
		font-family: inherit;
	}

	.empty-hint {
		margin: 0;
		font-size: 13px;
		color: #767586;
	}
</style>
