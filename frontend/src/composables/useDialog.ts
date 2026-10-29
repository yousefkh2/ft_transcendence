import { nextTick, onUnmounted, watch, type Ref } from 'vue';

const FOCUSABLE = 'a[href], button:not([disabled]), input:not([disabled]), select:not([disabled]), textarea:not([disabled]), [tabindex]:not([tabindex="-1"])';

// Keyboard behavior every modal dialog needs:
// - focus moves into the dialog when it opens and back to where it was when it closes,
// - Tab and Shift+Tab stay inside the dialog,
// - Escape closes it,
// - the page behind does not scroll.
export function useDialog(open: Ref<boolean>, dialog: Ref<HTMLElement | null>, close: () => void) {
	let returnFocusTo: HTMLElement | null = null;

	function focusable(): HTMLElement[] {
		return Array.from(dialog.value?.querySelectorAll<HTMLElement>(FOCUSABLE) ?? []).filter((el) => el.offsetParent !== null);
	}

	function onKeydown(event: KeyboardEvent) {
		if (event.key === 'Escape') {
			event.preventDefault();
			close();
			return;
		}
		if (event.key !== 'Tab') return;

		const items = focusable();
		if (!items.length) return;
		const first = items[0];
		const last = items[items.length - 1];
		const active = document.activeElement;
		if (event.shiftKey && (active === first || !dialog.value?.contains(active))) {
			event.preventDefault();
			last.focus();
		} else if (!event.shiftKey && (active === last || !dialog.value?.contains(active))) {
			event.preventDefault();
			first.focus();
		}
	}

	function activate() {
		returnFocusTo = document.activeElement instanceof HTMLElement ? document.activeElement : null;
		document.addEventListener('keydown', onKeydown);
		document.body.style.overflow = 'hidden';
		// The first form field, so the user can start typing right away.
		nextTick(() => (dialog.value?.querySelector<HTMLElement>('input') ?? focusable()[0])?.focus());
	}

	function deactivate() {
		document.removeEventListener('keydown', onKeydown);
		document.body.style.overflow = '';
		returnFocusTo?.focus();
		returnFocusTo = null;
	}

	watch(open, (isOpen) => (isOpen ? activate() : deactivate()), { immediate: true });
	onUnmounted(() => open.value && deactivate());
}
