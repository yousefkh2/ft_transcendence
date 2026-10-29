<template>
	<section id="gamification" class="gamification">
		<div class="header">
			<BaseEyebrow>{{ t('landing.gamification.eyebrow') }}</BaseEyebrow>
			<h2 class="title">{{ t('landing.gamification.title') }}</h2>
		</div>

		<div class="grid">
			<BaseCard size="lg">
				<h3 class="card-title">{{ t('landing.gamification.xpTitle') }}</h3>
				<p class="card-desc">{{ t('landing.gamification.xpText') }}</p>
				<div class="xp-widget">
					<div class="xp-row">
						<span class="xp-level">{{ t('landing.gamification.level', { level: 7 }) }}</span>
						<span class="xp-amount">{{ formatNumber(2140) }} / {{ t('common.xp', { xp: formatNumber(3000) }) }}</span>
					</div>
					<div class="xp-track">
						<div class="xp-fill" style="width: 71%"></div>
					</div>
				</div>
			</BaseCard>

			<BaseCard size="lg">
				<h3 class="card-title">{{ t('landing.gamification.flashcardsTitle') }}</h3>
				<p class="card-desc">{{ t('landing.gamification.flashcardsText') }}</p>
				<div class="flashcard">
					<span class="flashcard-wrong">Ich habe gegeht nach Hause.</span>
					<span class="flashcard-correct">Ich bin nach Hause gegangen.</span>
					<p class="flashcard-explanation">{{ t('landing.gamification.flashcardExplanation') }}</p>
				</div>
			</BaseCard>

			<BaseCard size="lg">
				<h3 class="card-title">{{ t('landing.gamification.friendsTitle') }}</h3>
				<p class="card-desc">{{ t('landing.gamification.friendsText') }}</p>
				<ul class="friends-list">
					<li v-for="(friend, index) in friends" :key="friend.xp" class="friend-row">
						<span class="friend-rank">{{ index + 1 }}</span>
						<span class="friend-name">{{ friend.name ?? t('landing.gamification.you') }}</span>
						<span class="friend-xp">{{ t('common.xp', { xp: formatNumber(friend.xp) }) }}</span>
					</li>
				</ul>
			</BaseCard>
		</div>
	</section>
</template>

<script setup lang="ts">
import BaseEyebrow from '../base/BaseEyebrow.vue';
import BaseCard from '../base/BaseCard.vue';
import { t } from '../../i18n';
import { formatNumber } from '../../utils/format';
const friends = [
	{ name: 'jonas.k', xp: 5820 },
	{ name: 'mei_learns', xp: 4390 },
	// The demo player shown as "you" in the current language.
	{ name: undefined, xp: 2140 },
];
</script>

<style scoped>
	.gamification {
		max-width: 1200px;
		margin: 0 auto;
		padding: 64px 40px;
	}

	.header {
		max-width: 620px;
		margin: 0 auto 48px;
		text-align: center;
	}

	.title {
		font-size: clamp(28px, 4vw, 40px);
		font-weight: 800;
		letter-spacing: -0.02em;
		line-height: 1.15;
		margin: 12px 0 0;
	}

	.grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
		gap: 20px;
	}

	.card-title {
		font-size: 18px;
		font-weight: 700;
		margin: 0 0 8px;
	}

	.card-desc {
		font-size: 14.5px;
		line-height: 1.5;
		color: var(--color-on-surface-variant);
		margin: 0 0 20px;
	}

	.xp-widget {
		display: flex;
		flex-direction: column;
		gap: 8px;
	}

	.xp-row {
		display: flex;
		align-items: baseline;
		justify-content: space-between;
		font-family: var(--font-mono);
	}

	.xp-level {
		font-weight: 700;
		font-size: var(--text-md);
		color: var(--color-on-surface);
	}

	.xp-amount {
		font-size: var(--text-xs);
		color: var(--color-on-surface-variant);
	}

	.xp-track {
		width: 100%;
		height: 10px;
		border-radius: var(--radius-pill);
		background: var(--color-border);
		overflow: hidden;
	}

	.xp-fill {
		height: 100%;
		border-radius: var(--radius-pill);
		background: linear-gradient(90deg, var(--color-primary-container), var(--color-primary));
	}

	.flashcard {
		display: flex;
		flex-direction: column;
		gap: 6px;
		padding: 16px;
		border-radius: 12px;
		background: var(--color-surface);
		border: 1px solid var(--color-border-strong);
	}

	.flashcard-wrong {
		font-size: var(--text-md);
		color: var(--color-on-surface-variant);
		text-decoration: line-through;
		text-decoration-color: var(--color-error);
	}

	.flashcard-correct {
		font-size: 15px;
		font-weight: 700;
		color: var(--color-on-surface);
	}

	.flashcard-explanation {
		font-size: var(--text-sm);
		line-height: 1.5;
		color: var(--color-on-surface-variant);
		margin: 4px 0 0;
	}

	.friends-list {
		list-style: none;
		display: flex;
		flex-direction: column;
		gap: 4px;
		margin: 0;
		padding: 0;
	}

	.friend-row {
		display: flex;
		align-items: center;
		gap: 12px;
		padding: 10px 12px;
		border-radius: 10px;
	}

	.friend-row:nth-child(odd) {
		background: var(--color-surface);
	}

	.friend-rank {
		font-family: var(--font-mono);
		font-weight: 700;
		font-size: var(--text-sm);
		color: var(--color-primary-container);
		width: 16px;
	}

	.friend-name {
		flex: 1;
		font-weight: 600;
		font-size: var(--text-md);
		color: var(--color-on-surface);
	}

	.friend-xp {
		font-family: var(--font-mono);
		font-size: 12.5px;
		color: var(--color-on-surface-variant);
	}

	@media (max-width: 900px) {
		.gamification {
			padding: 48px 16px;
		}
	}
</style>
