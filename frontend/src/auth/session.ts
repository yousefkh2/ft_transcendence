const TOKEN_KEY = 'lingodash_token';

export function getToken(): string | null {
	return localStorage.getItem(TOKEN_KEY);
}

export function setToken(token: string) {
	localStorage.setItem(TOKEN_KEY, token);
}

export function clearToken() {
	localStorage.removeItem(TOKEN_KEY);
}

function tokenExpiry(token: string): number | null {
	try {
		const payload = token.split('.')[1].replace(/-/g, '+').replace(/_/g, '/');
		const { exp } = JSON.parse(atob(payload));
		return typeof exp === 'number' ? exp * 1000 : null;
	} catch {
		return null;
	}
}

export function hasValidSession(): boolean {
	const token = getToken();
	if (!token) return false;

	const expiresAt = tokenExpiry(token);
	if (expiresAt === null || expiresAt <= Date.now()) {
		clearToken();
		return false;
	}
	return true;
}
