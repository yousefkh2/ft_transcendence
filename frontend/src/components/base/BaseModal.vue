<template>
	<Transition name="modal">
		<div v-if="open" class="overlay" @click.self="close">
			<div ref="dialogEl" class="modal" role="dialog" aria-modal="true" :aria-labelledby="labelledby">
				<button type="button" class="close-btn" aria-label="Close" @click="close">
					<BaseIcon name="close" :size="16" />
				</button>
				<slot />
			</div>
		</div>
	</Transition>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useDialog } from '../../composables/useDialog';
import BaseIcon from './BaseIcon.vue';

// A centered dialog over a blurred page. Escape, the close button and a click outside close it.
// labelledby is the id of the dialog's heading, which screen readers read as its name.
defineProps<{ labelledby: string }>();
const open = defineModel<boolean>('open', { default: false });

const dialogEl = ref<HTMLElement | null>(null);

function close() {
	open.value = false;
}

useDialog(open, dialogEl, close);
</script>

<style scoped>
	.overlay {
		position: fixed;
		inset: 0;
		z-index: 200;
		display: flex;
		align-items: center;
		justify-content: center;
		padding: 16px;
		background: rgba(10, 12, 28, 0.6);
		backdrop-filter: blur(6px);
		-webkit-backdrop-filter: blur(6px);
	}

	.modal {
		position: relative;
		width: 100%;
		max-width: 400px;
		max-height: 100%;
		overflow-y: auto;
		padding: 40px 36px;
		border-radius: 28px;
		background: rgba(255, 255, 255, 0.86);
		backdrop-filter: blur(20px);
		-webkit-backdrop-filter: blur(20px);
		border: 1px solid rgba(255, 255, 255, 0.5);
		box-shadow: var(--shadow-lg);
	}

	.close-btn {
		position: absolute;
		top: 16px;
		right: 16px;
		z-index: 2;
		width: 32px;
		height: 32px;
		display: flex;
		align-items: center;
		justify-content: center;
		border: none;
		border-radius: 50%;
		background: rgba(17, 28, 45, 0.06);
		color: var(--color-on-surface-variant);
		cursor: pointer;
		transition: background 0.15s ease;
	}

	.close-btn:hover {
		background: rgba(17, 28, 45, 0.12);
	}

	.modal-enter-active {
		animation: overlay-in 0.22s ease;
	}

	.modal-enter-active .modal {
		animation: modal-in 0.32s cubic-bezier(0.16, 1, 0.3, 1);
	}

	@keyframes overlay-in {
		from {
			opacity: 0;
		}
	}

	@keyframes modal-in {
		from {
			opacity: 0;
			transform: translateY(24px) scale(0.96);
		}
	}

	@media (prefers-reduced-motion: reduce) {
		.modal-enter-active,
		.modal-enter-active .modal {
			animation: none;
		}
	}
</style>
