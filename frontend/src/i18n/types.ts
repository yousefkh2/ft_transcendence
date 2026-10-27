import type en from './locales/en';

// Plural forms as defined by Intl.PluralRules. English only needs one and other, Polish also few and many.
export interface Plural {
	one?: string;
	few?: string;
	many?: string;
	other: string;
}

// The shape of the English messages with every text widened to string, so the other languages
// must have exactly the same keys. A missing or extra key is a type error.
type Widen<T> = T extends string ? string : T extends { other: string } ? Plural : { [K in keyof T]: Widen<T[K]> };
export type Messages = Widen<typeof en>;

// Every valid key as a dotted path, like 'friends.title'.
type Paths<T, Prefix extends string = ''> = {
	[K in keyof T & string]: T[K] extends string | { other: string } ? `${Prefix}${K}` : Paths<T[K], `${Prefix}${K}.`>;
}[keyof T & string];
export type MessageKey = Paths<typeof en>;
