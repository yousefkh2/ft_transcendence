<template>
	<section id="how-it-works" class="how-it-works">
		<div class="header">
			<span class="eyebrow">How It Works</span>
			<h2 class="title">Four steps to fluency under pressure</h2>
		</div>

		<div class="grid">
			<div class="card" v-for="step in steps" :key="step.number">
				<span class="step-number">{{ step.number }}</span>
				<h3 class="card-title">{{ step.title }}</h3>
				<p class="card-desc">{{ step.description }}</p>
			</div>
		</div>
	</section>
</template>

<script setup lang="ts">
const steps = [
	{
		number: '01',
		title: 'Get matched',
		description: 'Paired live with another learner, anywhere in the world.',
	},
	{
		number: '02',
		title: 'Get a secret role',
		description: 'Each player sees different information — nobody has the full picture alone.',
	},
	{
		number: '03',
		title: 'Speak it out',
		description: 'Describe, argue, negotiate — out loud, in the target language, before the clock runs out.',
	},
	{
		number: '04',
		title: 'Review & level up',
		description: 'Your mistakes become AI flashcards. Clear them to earn XP and climb levels.',
	},
];
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

	.eyebrow {
		display: inline-block;
		font-family: var(--font-mono);
		font-weight: 700;
		font-size: 13px;
		letter-spacing: 0.1em;
		text-transform: uppercase;
		color: var(--color-primary);
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

	.card {
		background: #ffffff;
		border: 1px solid rgba(118, 117, 134, 0.18);
		border-radius: 16px;
		padding: 28px;
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
