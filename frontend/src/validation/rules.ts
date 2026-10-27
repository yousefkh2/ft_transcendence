import { t } from '../i18n';

// Each validator returns an error message in the current language, or an empty string when the value is valid.

const USERNAME_MIN = 3;
const USERNAME_MAX = 20;
const USERNAME_PATTERN = /^[A-Za-z0-9_-]+$/;

const EMAIL_MAX = 254;
const EMAIL_PATTERN = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

const PASSWORD_MIN = 8;
// bcrypt, which the backend uses for hashing, only accepts up to 72 bytes.
const PASSWORD_MAX_BYTES = 72;

export function passwordHint() {
	return t('validation.passwordHint', { min: PASSWORD_MIN });
}

export function validateUsername(value: string): string {
	const username = value.trim();
	if (!username) return t('validation.usernameRequired');
	if (username.length < USERNAME_MIN) return t('validation.usernameTooShort', { min: USERNAME_MIN });
	if (username.length > USERNAME_MAX) return t('validation.usernameTooLong', { max: USERNAME_MAX });
	if (!USERNAME_PATTERN.test(username)) return t('validation.usernameChars');
	return '';
}

export function validateEmail(value: string): string {
	const email = value.trim();
	if (!email) return t('validation.emailRequired');
	if (email.length > EMAIL_MAX || !EMAIL_PATTERN.test(email)) return t('validation.emailInvalid');
	return '';
}

export function validateNewPassword(value: string): string {
	if (!value) return t('validation.passwordRequired');
	if (value.length < PASSWORD_MIN) return t('validation.passwordTooShort', { min: PASSWORD_MIN });
	if (new TextEncoder().encode(value).length > PASSWORD_MAX_BYTES) return t('validation.passwordTooLong');
	if (!/[A-Za-z]/.test(value) || !/\d/.test(value)) return t('validation.passwordLetterNumber');
	return '';
}

export function validateLoginPassword(value: string): string {
	return value ? '' : t('validation.enterPassword');
}
