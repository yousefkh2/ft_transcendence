// The languages the game supports. Must match SupportedLobbyLangs in backend/internal/db/lobby.go.
export const LEARNING_LANGUAGES = [
	{ value: 'de', label: 'German' },
	{ value: 'en', label: 'English' },
	{ value: 'pl', label: 'Polish' },
	{ value: 'tr', label: 'Turkish' },
];

export const DEFAULT_LEARNING_LANGUAGE = 'de';

export function languageName(code: string) {
	return LEARNING_LANGUAGES.find((l) => l.value === code)?.label ?? code;
}
