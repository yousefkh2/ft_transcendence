import { setToken } from '../auth/session';
import { request } from './client';

export async function login(email: string, password: string) {
	const { token } = await request<{ token: string }>('/api/auth/login', {
		method: 'POST',
		body: { email, password },
		auth: false,
	});
	setToken(token);
}

export async function register(username: string, email: string, password: string) {
	await request<{ id: string }>('/api/auth/register', {
		method: 'POST',
		body: { username, email, password },
		auth: false,
	});
}
