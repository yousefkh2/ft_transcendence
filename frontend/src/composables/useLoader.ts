import { ref, shallowRef, type Ref } from 'vue';

export type LoadStatus = 'loading' | 'error' | 'ready';

// Loads data and tracks whether it is still loading, failed or ready, so a page can show exactly one
// of skeleton, error or content. After the first success, a failed reload (for example from polling)
// keeps the old data on screen instead of replacing it with an error.
export function useLoader<T>(fetcher: () => Promise<T>, initial: T) {
	const data = shallowRef(initial) as Ref<T>;
	const status = ref<LoadStatus>('loading');

	async function load() {
		if (status.value !== 'ready') status.value = 'loading';
		try {
			data.value = await fetcher();
			status.value = 'ready';
		} catch {
			if (status.value !== 'ready') status.value = 'error';
		}
	}

	return { data, status, load };
}
