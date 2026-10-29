<template>
	<section id="how-it-works" class="how-it-works">
		<div class="header">
			<BaseEyebrow>{{ t('landing.steps.eyebrow') }}</BaseEyebrow>
			<h2 class="title">{{ t('landing.steps.title') }}</h2>
		</div>

		<div class="grid">
			<BaseCard v-for="step in steps" :key="step.number" size="lg">
				<span class="step-number">{{ step.number }}</span>
				<h3 class="card-title">{{ t(step.title) }}</h3>
				<p class="card-desc">{{ t(step.description) }}</p>
			</BaseCard>
		</div>
	</section>
</template>

<script setup lang="ts">
import BaseEyebrow from '../base/BaseEyebrow.vue';
import BaseCard from '../base/BaseCard.vue';
import { t } from '../../i18n';
const steps = [
	{ number: '01', title: 'landing.steps.matchedTitle', description: 'landing.steps.matchedText' },
	{ number: '02', title: 'landing.steps.roleTitle', description: 'landing.steps.roleText' },
	{ number: '03', title: 'landing.steps.speakTitle', description: 'landing.steps.speakText' },
	{ number: '04', title: 'landing.steps.reviewTitle', description: 'landing.steps.reviewText' },
] as const;
</script>

<style scoped>
	.how-it-works {
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
		grid-template-columns: repeat(auto-fit, minmax(230px, 1fr));
		gap: 20px;
	}

	.step-number {
		display: block;
		font-family: var(--font-mono);
		font-weight: 800;
		font-size: 20px;
		color: var(--color-primary-container);
		margin-bottom: 12px;
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
		margin: 0;
	}

	@media (max-width: 900px) {
		.how-it-works {
			padding: 48px 16px;
		}
	}
</style>
