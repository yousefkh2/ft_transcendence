import { t } from '.';
import type { MessageKey } from './types';

// The backend answers in English. Messages users can actually run into are mapped to a translated text.
const KNOWN: Record<string, MessageKey> = {
	'invalid credentials': 'apiErrors.invalidCredentials',
	'incorrect password': 'apiErrors.incorrectPassword',
	'username already taken': 'apiErrors.usernameTaken',
	'email already registered': 'apiErrors.emailTaken',
	'user already exists': 'apiErrors.userExists',
	'user not found': 'apiErrors.userNotFound',
	'friend request already sent': 'apiErrors.requestAlreadySent',
	'cannot send friend request to yourself': 'apiErrors.requestToYourself',
	'friend request not found': 'apiErrors.requestNotFound',
	'friend request already resolved': 'apiErrors.requestNotFound',
	'friend not found': 'apiErrors.friendNotFound',
	'Please choose a supported language.': 'apiErrors.unsupportedLanguage',
	'unsupported language': 'apiErrors.unsupportedLanguage',
	'Please upload a PNG or JPEG image.': 'apiErrors.avatarFormat',
	'The picture must be 2 MB or smaller.': 'apiErrors.avatarTooLarge',
	'The picture can be at most 1024 x 1024 pixels.': 'apiErrors.avatarDimensions',
	// Same rules as the frontend validation, so normally caught before the request is sent.
	'Please choose a username.': 'validation.usernameRequired',
	'Please enter your email.': 'validation.emailRequired',
	'Please enter a valid email address.': 'validation.emailInvalid',
	'Please choose a password.': 'validation.passwordRequired',
	'Password is too long.': 'validation.passwordTooLong',
	'Use at least one letter and one number.': 'validation.passwordLetterNumber',
};

// Unknown messages are technical (like "invalid json") and not helpful to show, so they get a generic text.
export function translateApiError(serverMessage: string | undefined): string {
	const key = serverMessage ? KNOWN[serverMessage] : undefined;
	return key ? t(key) : t('common.somethingWentWrong');
}
