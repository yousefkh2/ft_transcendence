import { request } from './client';

// Same limit as the backend (validate.ChatMessageMax) and the database constraint.
export const CHAT_MESSAGE_MAX = 1000;

export interface ChatMessage {
	id: string;
	senderId: string;
	recipientId: string;
	body: string;
	createdAt: string;
	readAt?: string;
}

// Newest messages with a friend, oldest first. With before, only messages older than that time.
export function fetchChatMessages(friendId: string, before?: string): Promise<ChatMessage[]> {
	const query = before ? `?before=${encodeURIComponent(before)}` : '';
	return request<ChatMessage[]>(`/api/chat/${encodeURIComponent(friendId)}/messages${query}`);
}

export function sendChatMessage(friendId: string, body: string): Promise<ChatMessage> {
	return request<ChatMessage>(`/api/chat/${encodeURIComponent(friendId)}/messages`, { method: 'POST', body: { body } });
}

export function markChatRead(friendId: string): Promise<void> {
	return request<void>(`/api/chat/${encodeURIComponent(friendId)}/read`, { method: 'POST' });
}

// Unread messages per friend id.
export function fetchUnreadCounts(): Promise<Record<string, number>> {
	return request<Record<string, number>>('/api/chat/unread');
}
