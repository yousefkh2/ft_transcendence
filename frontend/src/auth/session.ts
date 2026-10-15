const TOKEN_KEY = 'lingodash_token';

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
	const token = localStorage.getItem(TOKEN_KEY);
	if (!token) return false;

	const expiresAt = tokenExpiry(token);
	if (expiresAt === null || expiresAt <= Date.now()) {
		localStorage.removeItem(TOKEN_KEY);
		return false;
	}
	return true;
}
