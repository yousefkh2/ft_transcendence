<template>
	<div class="toasts" aria-live="polite" aria-atomic="false">
		<TransitionGroup name="toast">
			<div v-for="a in visible" :key="a.id" class="toast">
				<span class="icon" aria-hidden="true"><BaseIcon :name="ACHIEVEMENT_ICONS[a.id]" :size="20" /></span>
				<div class="body">
					<p class="eyebrow">{{ t('achievements.toastTitle') }}</p>
					<p class="title">{{ t(`achievements.items.${a.id}.title`) }}</p>
					<p class="text">{{ t(`achievements.items.${a.id}.text`) }}</p>
				</div>
				<button type="button" class="close" :aria-label="t('common.close')" @click="dismiss(a.id)">
					<BaseIcon name="close" :size="16" />
				</button>
			</div>
		</TransitionGroup>
	</div>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue';
import { fetchMyAchievements, markAchievementsSeen, type Achievement } from '../../api/progress';
import { ACHIEVEMENT_ICONS } from '../../constants/achievements';
import { t } from '../../i18n';
import BaseIcon from '../base/BaseIcon.vue';

// Shows newly unlocked achievements once, for example after coming back from a round.
const AUTO_HIDE_MS = 8000;

const visible = ref<Achievement[]>([]);
const timers: ReturnType<typeof setTimeout>[] = [];

function dismiss(id: string) {
	visible.value = visible.value.filter((a) => a.id !== id);
}

onMounted(async () => {
	try {
		const fresh = (await fetchMyAchievements()).filter((a) => a.new);
		if (!fresh.length) return;
		visible.value = fresh;
		fresh.forEach((a, i) => timers.push(setTimeout(() => dismiss(a.id), AUTO_HIDE_MS + i * 1500)));
		await markAchievementsSeen();
	} catch {
		// Not important enough to show an error: they are still listed on the profile.
	}
});

onUnmounted(() => timers.forEach(clearTimeout));
</script>

<style scoped>
	.toasts {
		position: fixed;
		right: 20px;
		bottom: 20px;
		z-index: 60;
		display: grid;
		gap: 10px;
		width: min(360px, calc(100vw - 32px));
	}

	.toast {
		display: flex;
		align-items: flex-start;
		gap: 12px;
		padding: 14px;
		border-radius: var(--radius-lg);
		background: var(--color-surface-inverse);
		color: #fff;
		box-shadow: var(--shadow-md);
	}

	.icon {
		display: grid;
		place-items: center;
		width: 40px;
		height: 40px;
		flex-shrink: 0;
		border-radius: var(--radius-sm);
		background: linear-gradient(135deg, var(--color-primary), var(--color-primary-container));
	}

	.body {
		flex: 1;
		min-width: 0;
	}

	.body p {
		margin: 0;
	}

	.eyebrow {
		font-family: var(--font-mono);
		font-size: var(--text-2xs);
		font-weight: 700;
		letter-spacing: 0.08em;
		text-transform: uppercase;
		/* Light tint of the primary purple: readable on the dark toast. */
		color: #b9baff;
	}

	.title {
		font-weight: 800;
		font-size: 15px;
	}

	.text {
		font-size: 12.5px;
		color: rgba(255, 255, 255, 0.8);
	}

	.close {
		border: none;
		background: none;
		color: rgba(255, 255, 255, 0.8);
		cursor: pointer;
		padding: 2px;
		border-radius: 6px;
	}

	.close:focus-visible {
		outline: 2px solid #fff;
	}

	.toast-enter-active,
	.toast-leave-active {
		transition: opacity 0.25s ease, transform 0.25s ease;
	}

	.toast-enter-from,
	.toast-leave-to {
		opacity: 0;
		transform: translateY(12px);
	}

	@media (prefers-reduced-motion: reduce) {
		.toast-enter-active,
		.toast-leave-active {
			transition: none;
		}
	}
</style>
