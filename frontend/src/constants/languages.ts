import { locale } from '../i18n';

// The languages the game supports. Must match SupportedLobbyLangs in backend/internal/db/lobby.go.
const LEARNING_LANGUAGE_CODES = ['de', 'en', 'pl', 'tr'];

export const DEFAULT_LEARNING_LANGUAGE = 'de';

// The name of a language in the current interface language, e.g. "German", "Deutsch", "Almanca".
export function languageName(code: string) {
	const name = new Intl.DisplayNames(locale.value, { type: 'language' }).of(code) ?? code;
	return name.charAt(0).toLocaleUpperCase(locale.value) + name.slice(1);
}

// Options for a select; call it inside a computed so the names follow a language switch.
export function learningLanguageOptions() {
	return LEARNING_LANGUAGE_CODES.map((value) => ({ value, label: languageName(value) }));
}
