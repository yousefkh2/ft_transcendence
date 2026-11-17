<template>
	<div class="chat" :class="{ 'has-selection': !!username }">
		<BaseCard tag="section" size="sm" class="conversations" :aria-label="t('chat.conversations')">
			<h1 class="panel-title">{{ t('chat.title') }}</h1>
			<ListSkeleton v-if="status === 'loading'" :label="t('chat.loadingFriends')" avatar />
			<BaseAlert v-else-if="status === 'error'" @retry="loadFriends">{{ t('chat.friendsError') }}</BaseAlert>
			<ul v-else-if="friends.length" class="friend-list">
				<li v-for="friend in friends" :key="friend.friendId">
					<RouterLink
						:to="{ name: 'chat', params: { username: friend.friendUsername } }"
						class="friend"
						:class="{ active: friend.friendId === selected?.friendId }"
						:aria-current="friend.friendId === selected?.friendId ? 'page' : undefined"
					>
						<BaseAvatar :name="friend.friendUsername" :src="apiAssetUrl(friend.avatarUrl)" :status="friend.online ? 'online' : 'offline'" />
						<span class="friend-info">
							<span class="friend-name">{{ friend.friendUsername }}</span>
							<span class="friend-meta" :class="{ online: friend.online }">{{ formatLastSeen(friend.online, friend.lastSeenAt) }}</span>
						</span>
						<span v-if="unreadFor(friend.friendId)" class="unread" :aria-label="t('chat.unread', { n: unreadFor(friend.friendId) })">
							{{ unreadFor(friend.friendId) }}
						</span>
					</RouterLink>
				</li>
			</ul>
			<div v-else class="empty-hint">
				<p>{{ t('chat.noFriends') }}</p>
				<BaseButton :to="{ name: 'friends' }" size="sm" pill><BaseIcon name="user-plus" :size="16" />{{ t('chat.findFriends') }}</BaseButton>
			</div>
		</BaseCard>

		<BaseCard tag="section" size="sm" class="conversation">
			<template v-if="selected">
				<header class="conversation-head">
					<RouterLink :to="{ name: 'chat' }" class="back" :aria-label="t('chat.back')">
						<BaseIcon name="arrow-left" :size="18" />
					</RouterLink>
					<BaseAvatar :name="selected.friendUsername" :src="apiAssetUrl(selected.avatarUrl)" :status="selected.online ? 'online' : 'offline'" size="md" />
					<div class="head-info">
						<h2 class="head-name">{{ selected.friendUsername }}</h2>
						<span class="friend-meta" :class="{ online: selected.online }">{{ formatLastSeen(selected.online, selected.lastSeenAt) }}</span>
					</div>
					<BaseButton :to="{ name: 'user-profile', params: { username: selected.friendUsername } }" variant="ghost" size="sm" pill>
						{{ t('chat.viewProfile') }}
					</BaseButton>
				</header>

				<p v-if="!connected" class="connection" role="status">{{ t('chat.reconnecting') }}</p>

				<div ref="scroller" class="messages" role="log" :aria-label="t('chat.historyLabel', { name: selected.friendUsername })" tabindex="0">
					<ListSkeleton v-if="messagesStatus === 'loading'" :label="t('chat.loading')" />
					<BaseAlert v-else-if="messagesStatus === 'error'" @retry="openConversation">{{ t('chat.loadError') }}</BaseAlert>
					<template v-else>
						<div v-if="hasMore" class="older">
							<BaseButton variant="outline" size="sm" pill :disabled="loadingOlder" @click="loadOlder">{{ t('chat.loadOlder') }}</BaseButton>
						</div>
						<p v-if="!messages.length" class="empty-hint">{{ t('chat.empty', { name: selected.friendUsername }) }}</p>
						<div v-for="message in messages" :key="message.id" class="bubble-row" :class="{ mine: message.senderId === myId }">
							<div class="bubble">
								<span class="sr-only">{{ message.senderId === myId ? t('chat.you') : selected.friendUsername }}:</span>
								<p class="bubble-text">{{ message.body }}</p>
								<time class="bubble-time" :datetime="message.createdAt">{{ formatMessageTime(message.createdAt) }}</time>
							</div>
						</div>
					</template>
				</div>

				<form class="composer" @submit.prevent="send">
					<label :for="inputId" class="sr-only">{{ t('chat.inputLabel', { name: selected.friendUsername }) }}</label>
					<textarea
						:id="inputId"
						ref="input"
						v-model="draft"
						class="composer-input"
						:class="{ invalid: !!sendError }"
						rows="1"
						:maxlength="CHAT_MESSAGE_MAX"
						:placeholder="t('chat.placeholder', { name: selected.friendUsername })"
						:aria-invalid="sendError ? 'true' : undefined"
						:aria-describedby="sendError ? errorId : counterId"
						@keydown.enter.exact.prevent="send"
					></textarea>
					<BaseButton type="submit" :disabled="sending || !draft.trim()" :aria-label="t('chat.send')">
						<BaseIcon name="send" :size="16" /><span class="send-text">{{ t('chat.send') }}</span>
					</BaseButton>
				</form>
				<p v-if="sendError" :id="errorId" class="composer-error" role="alert">{{ sendError }}</p>
				<p v-else :id="counterId" class="composer-hint" :class="{ warn: draftLength > CHAT_MESSAGE_MAX - 100 }">
					{{ t('chat.counter', { n: draftLength, max: CHAT_MESSAGE_MAX }) }}
				</p>
			</template>

			<BaseAlert v-else-if="username && status === 'ready'">{{ t('chat.notFriend', { name: username }) }}</BaseAlert>
			<div v-else class="placeholder">
				<BaseIcon name="message" :size="32" />
				<p>{{ t('chat.pickFriend') }}</p>
			</div>
		</BaseCard>
	</div>
</template>

<script setup lang="ts">
import { computed, nextTick, ref, useId, watch } from 'vue';
import { useRoute } from 'vue-router';
import { apiAssetUrl } from '../../api/client';
import { CHAT_MESSAGE_MAX, fetchChatMessages, fetchUnreadCounts, markChatRead, sendChatMessage, type ChatMessage } from '../../api/chat';
import { fetchFriends, type Friend } from '../../api/friends';
import type { Profile } from '../../api/users';
import { useChatSocket } from '../../composables/useChatSocket';
import { PRESENCE_INTERVAL_MS } from '../../composables/useHeartbeat';
import { useLoader, type LoadStatus } from '../../composables/useLoader';
import { usePolling } from '../../composables/usePolling';
import { locale, t } from '../../i18n';
import { formatLastSeen } from '../../utils/format';
import { sortByPresence } from '../../utils/friends';
import BaseAlert from '../base/BaseAlert.vue';
import BaseAvatar from '../base/BaseAvatar.vue';
import BaseButton from '../base/BaseButton.vue';
import BaseCard from '../base/BaseCard.vue';
import BaseIcon from '../base/BaseIcon.vue';
import ListSkeleton from './ListSkeleton.vue';

const PAGE_SIZE = 50;

const props = defineProps<{ profile: Profile | null }>();
const myId = computed(() => props.profile?.id ?? '');

const route = useRoute();
const username = computed(() => (typeof route.params.username === 'string' ? route.params.username : ''));

// Friends and unread counters, polled like the friends tab so the online status stays fresh.
const unread = ref<Record<string, number>>({});
const {
	data: friendsData,
	status,
	load: loadFriends,
} = useLoader(async () => {
	const [list, counts] = await Promise.all([fetchFriends(), fetchUnreadCounts()]);
	unread.value = counts;
	return sortByPresence(list);
}, [] as Friend[]);
const friends = computed(() => friendsData.value);
usePolling(loadFriends, PRESENCE_INTERVAL_MS);

const selected = computed(() => friends.value.find((f) => f.friendUsername === username.value) ?? null);

// The open conversation counts as read, even if a poll returns before the read was saved.
function unreadFor(friendId: string) {
	return friendId === selected.value?.friendId ? 0 : (unread.value[friendId] ?? 0);
}

// Messages of the open conversation.
const messages = ref<ChatMessage[]>([]);
const messagesStatus = ref<LoadStatus>('loading');
const hasMore = ref(false);
const loadingOlder = ref(false);
const scroller = ref<HTMLElement | null>(null);
const input = ref<HTMLTextAreaElement | null>(null);
// The message being written; cleared when another conversation is opened.
const draft = ref('');
const sendError = ref('');

function addMessage(message: ChatMessage) {
	if (messages.value.some((m) => m.id === message.id)) return;
	messages.value = [...messages.value, message];
}

function isNearBottom() {
	const el = scroller.value;
	return !el || el.scrollHeight - el.scrollTop - el.clientHeight < 80;
}

async function scrollToBottom() {
	await nextTick();
	if (scroller.value) scroller.value.scrollTop = scroller.value.scrollHeight;
}

function markRead(friendId: string) {
	markChatRead(friendId).catch(() => {});
	unread.value = { ...unread.value, [friendId]: 0 };
}

async function openConversation() {
	const friend = selected.value;
	if (!friend) return;
	messagesStatus.value = 'loading';
	try {
		const page = await fetchChatMessages(friend.friendId);
		// The user may have switched to another friend while this was loading.
		if (selected.value?.friendId !== friend.friendId) return;
		messages.value = page;
		hasMore.value = page.length === PAGE_SIZE;
		messagesStatus.value = 'ready';
		markRead(friend.friendId);
		await scrollToBottom();
	} catch {
		if (selected.value?.friendId === friend.friendId) messagesStatus.value = 'error';
	}
}

async function loadOlder() {
	const friend = selected.value;
	const oldest = messages.value[0];
	if (!friend || !oldest) return;
	loadingOlder.value = true;
	const el = scroller.value;
	const distanceFromBottom = el ? el.scrollHeight - el.scrollTop : 0;
	try {
		const page = await fetchChatMessages(friend.friendId, oldest.createdAt);
		if (selected.value?.friendId !== friend.friendId) return;
		messages.value = [...page, ...messages.value];
		hasMore.value = page.length === PAGE_SIZE;
		// Keep the messages the user was reading in the same place.
		await nextTick();
		if (el) el.scrollTop = el.scrollHeight - distanceFromBottom;
	} catch {
		sendError.value = t('chat.loadError');
	} finally {
		loadingOlder.value = false;
	}
}

watch(
	() => selected.value?.friendId,
	(id, previous) => {
		if (id === previous) return;
		messages.value = [];
		draft.value = '';
		sendError.value = '';
		if (id) openConversation();
	},
	{ immediate: true },
);

// Live messages: add them to the open conversation, or count them as unread for the other friends.
const { connected } = useChatSocket(
	(message) => {
		const friend = selected.value;
		const partner = message.senderId === myId.value ? message.recipientId : message.senderId;
		if (friend && partner === friend.friendId) {
			const stick = message.senderId === myId.value || isNearBottom();
			addMessage(message);
			if (message.senderId !== myId.value) markRead(friend.friendId);
			if (stick) scrollToBottom();
		} else if (message.recipientId === myId.value) {
			unread.value = { ...unread.value, [message.senderId]: (unread.value[message.senderId] ?? 0) + 1 };
		}
	},
	// After a reconnect, reload what may have arrived while the connection was down.
	() => {
		loadFriends();
		if (selected.value) openConversation();
	},
);

// Writing a message.
const inputId = useId();
const errorId = useId();
const counterId = useId();
const sending = ref(false);
const draftLength = computed(() => [...draft.value].length);

watch(draft, () => {
	sendError.value = '';
});

function validateDraft(text: string): string {
	if (!text) return t('chat.emptyMessage');
	if ([...text].length > CHAT_MESSAGE_MAX) return t('chat.tooLong', { max: CHAT_MESSAGE_MAX });
	return '';
}

async function send() {
	const friend = selected.value;
	const text = draft.value.trim();
	if (!friend || sending.value) return;
	sendError.value = validateDraft(text);
	if (sendError.value) return;

	sending.value = true;
	try {
		const message = await sendChatMessage(friend.friendId, text);
		addMessage(message);
		draft.value = '';
		scrollToBottom();
	} catch (err) {
		sendError.value = err instanceof Error ? err.message : t('chat.sendError');
	} finally {
		sending.value = false;
		input.value?.focus();
	}
}

// "14:05" for today, otherwise also the date.
function formatMessageTime(iso: string) {
	const date = new Date(iso);
	const sameDay = date.toDateString() === new Date().toDateString();
	return date.toLocaleString(locale.value, sameDay ? { hour: '2-digit', minute: '2-digit' } : { day: 'numeric', month: 'short', hour: '2-digit', minute: '2-digit' });
}
</script>

<style scoped>
	.chat {
		display: grid;
		grid-template-columns: minmax(220px, 300px) minmax(0, 1fr);
		gap: 18px;
		height: calc(100vh - 220px);
		min-height: 420px;
	}

	.conversations,
	.conversation {
		display: flex;
		flex-direction: column;
		min-height: 0;
	}

	.panel-title {
		font-size: var(--text-lg);
		font-weight: 800;
		margin: 0 0 12px;
	}

	.friend-list {
		list-style: none;
		margin: 0;
		padding: 0;
		display: grid;
		gap: 6px;
		overflow-y: auto;
	}

	.friend {
		display: flex;
		align-items: center;
		gap: 10px;
		padding: 9px 10px;
		border-radius: var(--radius-md);
		color: inherit;
		text-decoration: none;
		border: 1px solid transparent;
	}

	.friend:hover {
		background: var(--color-fill-muted);
	}

	.friend.active {
		background: var(--color-surface);
		border-color: var(--color-border-strong);
	}

	.friend:focus-visible,
	.back:focus-visible,
	.messages:focus-visible {
		outline: 2px solid var(--color-primary);
		outline-offset: 2px;
	}

	.friend-info {
		flex: 1;
		min-width: 0;
		display: flex;
		flex-direction: column;
	}

	.friend-name {
		font-size: 13.5px;
		font-weight: 700;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}

	.friend-meta {
		font-size: 11.5px;
		color: var(--color-outline);
	}

	.friend-meta.online {
		color: var(--color-on-success-container);
		font-weight: 700;
	}

	.unread {
		min-width: 20px;
		padding: 2px 6px;
		border-radius: var(--radius-pill);
		background: var(--color-primary);
		color: var(--color-on-primary);
		font-family: var(--font-mono);
		font-size: var(--text-2xs);
		font-weight: 700;
		text-align: center;
	}

	.conversation-head {
		display: flex;
		align-items: center;
		gap: 12px;
		padding-bottom: 12px;
		border-bottom: 1px solid var(--color-border-subtle);
	}

	.back {
		display: none;
		color: var(--color-outline);
		border-radius: var(--radius-sm);
	}

	.head-info {
		flex: 1;
		min-width: 0;
		display: flex;
		flex-direction: column;
	}

	.head-name {
		margin: 0;
		font-size: 15px;
		font-weight: 800;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}

	.connection {
		margin: 8px 0 0;
		font-size: var(--text-xs);
		color: var(--color-outline);
	}

	.messages {
		flex: 1;
		min-height: 0;
		overflow-y: auto;
		display: flex;
		flex-direction: column;
		gap: 6px;
		padding: 14px 2px;
	}

	.older {
		display: flex;
		justify-content: center;
		margin-bottom: 8px;
	}

	.bubble-row {
		display: flex;
	}

	.bubble-row.mine {
		justify-content: flex-end;
	}

	.bubble {
		max-width: min(75%, 520px);
		padding: 8px 12px;
		border-radius: 16px 16px 16px 4px;
		background: var(--color-surface);
		border: 1px solid var(--color-border-subtle);
	}

	.mine .bubble {
		border-radius: 16px 16px 4px 16px;
		background: var(--color-primary);
		border-color: var(--color-primary);
		color: var(--color-on-primary);
	}

	.bubble-text {
		margin: 0;
		font-size: var(--text-md);
		line-height: 1.4;
		white-space: pre-wrap;
		overflow-wrap: anywhere;
	}

	.bubble-time {
		display: block;
		margin-top: 2px;
		font-size: var(--text-2xs);
		color: var(--color-outline);
		text-align: right;
	}

	.mine .bubble-time {
		color: rgba(255, 255, 255, 0.85);
	}

	.composer {
		display: flex;
		align-items: flex-end;
		gap: 10px;
		padding-top: 12px;
		border-top: 1px solid var(--color-border-subtle);
	}

	.composer-input {
		flex: 1;
		min-height: 42px;
		max-height: 140px;
		resize: vertical;
		padding: 10px 14px;
		border-radius: var(--radius-md);
		border: 1px solid var(--color-border-strong);
		background: var(--color-surface-raised);
		color: var(--color-on-surface);
		font: inherit;
		font-size: var(--text-md);
	}

	.composer-input:focus {
		outline: none;
		border-color: var(--color-primary);
		box-shadow: 0 0 0 3px var(--color-focus-ring);
	}

	.composer-input.invalid {
		border-color: var(--color-error);
	}

	.composer-hint,
	.composer-error {
		margin: 6px 0 0;
		font-size: var(--text-2xs);
		color: var(--color-outline);
	}

	.composer-hint.warn,
	.composer-error {
		color: var(--color-error);
	}

	.empty-hint {
		margin: 0;
		font-size: var(--text-sm);
		color: var(--color-outline);
	}

	.empty-hint p {
		margin: 0 0 12px;
	}

	.messages .empty-hint {
		margin: auto;
		text-align: center;
	}

	.placeholder {
		flex: 1;
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		gap: 8px;
		color: var(--color-outline);
		font-size: var(--text-sm);
	}

	.sr-only {
		position: absolute;
		width: 1px;
		height: 1px;
		overflow: hidden;
		clip: rect(0 0 0 0);
		white-space: nowrap;
	}

	/* Phones: one column. The list shows until a friend is picked, then only the conversation with a back arrow. */
	@media (max-width: 720px) {
		.chat {
			grid-template-columns: minmax(0, 1fr);
			height: calc(100vh - 230px);
		}

		.chat.has-selection .conversations,
		.chat:not(.has-selection) .conversation {
			display: none;
		}

		.back {
			display: flex;
		}

		.send-text {
			position: absolute;
			width: 1px;
			height: 1px;
			overflow: hidden;
			clip: rect(0 0 0 0);
		}
	}
</style>
