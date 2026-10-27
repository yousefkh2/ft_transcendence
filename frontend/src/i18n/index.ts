import { ref } from 'vue';
import en from './locales/en';
import de from './locales/de';
import pl from './locales/pl';
import tr from './locales/tr';
import type { MessageKey, Messages, Plural } from './types';

// The interface languages. They match the languages the game teaches (constants/languages.ts).
export const LOCALES = [
	{ code: 'en', name: 'English' },
	{ code: 'de', name: 'Deutsch' },
	{ code: 'pl', name: 'Polski' },
	{ code: 'tr', name: 'Türkçe' },
] as const;

export type Locale = (typeof LOCALES)[number]['code'];

const MESSAGES: Record<Locale, Messages> = { en, de, pl, tr };
const STORAGE_KEY = 'lingodash_locale';

function isLocale(value: unknown): value is Locale {
	return LOCALES.some((l) => l.code === value);
}

// A saved choice wins, then the first browser language we support, then English.
function initialLocale(): Locale {
	try {
		const saved = localStorage.getItem(STORAGE_KEY);
		if (isLocale(saved)) return saved;
	} catch {
		// Storage can be blocked (private mode); fall back to the browser language.
	}
	for (const tag of navigator.languages ?? [navigator.language]) {
		const code = tag.slice(0, 2).toLowerCase();
		if (isLocale(code)) return code;
	}
	return 'en';
}

// Reading locale inside t() makes every template and computed that calls t() update on a switch.
export const locale = ref<Locale>(initialLocale());
document.documentElement.lang = locale.value;

export function setLocale(next: Locale) {
	locale.value = next;
	document.documentElement.lang = next;
	try {
		localStorage.setItem(STORAGE_KEY, next);
	} catch {
		// Not saved, but the switch still works for this visit.
	}
}

function lookup(messages: Messages, key: string): string | Plural | undefined {
	let node: unknown = messages;
	for (const part of key.split('.')) node = (node as Record<string, unknown> | undefined)?.[part];
	return typeof node === 'string' || (typeof node === 'object' && node !== null && 'other' in node)
		? (node as string | Plural)
		: undefined;
}

export type Params = Record<string, string | number>;

// Translates a key like 'friends.title'. {name} placeholders are filled from params; for plural
// messages params.n picks the form (one, few, many, other) by the rules of the current language.
export function t(key: MessageKey, params: Params = {}): string {
	const message = lookup(MESSAGES[locale.value], key) ?? lookup(en, key) ?? key;
	let text: string;
	if (typeof message === 'string') {
		text = message;
	} else {
		const form = new Intl.PluralRules(locale.value).select(Number(params.n ?? 0)) as keyof Plural;
		text = message[form] ?? message.other;
	}
	return text.replace(/\{(\w+)\}/g, (match, name: string) => (name in params ? String(params[name]) : match));
}
