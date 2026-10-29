import { onMounted, onUnmounted } from 'vue';

// Runs fn once when the component mounts and then every intervalMs until it unmounts.
export function usePolling(fn: () => void, intervalMs: number) {
	let timer: ReturnType<typeof setInterval> | undefined;
	onMounted(() => {
		fn();
		timer = setInterval(fn, intervalMs);
	});
	onUnmounted(() => clearInterval(timer));
}
