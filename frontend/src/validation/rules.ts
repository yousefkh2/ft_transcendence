// Each validator returns an error message, or an empty string when the value is valid.

const USERNAME_MIN = 3;
const USERNAME_MAX = 20;
const USERNAME_PATTERN = /^[A-Za-z0-9_-]+$/;

const EMAIL_MAX = 254;
const EMAIL_PATTERN = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

const PASSWORD_MIN = 8;
// bcrypt, which the backend uses for hashing, only accepts up to 72 bytes.
const PASSWORD_MAX_BYTES = 72;

export const PASSWORD_HINT = `At least ${PASSWORD_MIN} characters, with a letter and a number.`;

export function validateUsername(value: string): string {
	const username = value.trim();
	if (!username) return 'Please choose a username.';
	if (username.length < USERNAME_MIN) return `Username must be at least ${USERNAME_MIN} characters.`;
	if (username.length > USERNAME_MAX) return `Username can be at most ${USERNAME_MAX} characters.`;
	if (!USERNAME_PATTERN.test(username)) return 'Use only letters, numbers, "_" and "-".';
	return '';
}

export function validateEmail(value: string): string {
	const email = value.trim();
	if (!email) return 'Please enter your email.';
	if (email.length > EMAIL_MAX || !EMAIL_PATTERN.test(email)) return 'Please enter a valid email address.';
	return '';
}

export function validateNewPassword(value: string): string {
	if (!value) return 'Please choose a password.';
	if (value.length < PASSWORD_MIN) return `Password must be at least ${PASSWORD_MIN} characters.`;
	if (new TextEncoder().encode(value).length > PASSWORD_MAX_BYTES) return 'Password is too long.';
	if (!/[A-Za-z]/.test(value) || !/\d/.test(value)) return 'Use at least one letter and one number.';
	return '';
}

export function validateLoginPassword(value: string): string {
	return value ? '' : 'Please enter your password.';
}
