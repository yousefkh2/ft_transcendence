<template>
	<section id="gamification" class="gamification">
		<div class="header">
			<span class="eyebrow">Gamification</span>
			<h2 class="title">Progress you can feel every round</h2>
		</div>

		<div class="grid">
			<div class="card">
				<h3 class="card-title">XP &amp; levels</h3>
				<p class="card-desc">
					Every round you finish earns XP. Win, lose, or run out the clock -
					you always walk away a little further along.
				</p>
				<div class="xp-widget">
					<div class="xp-row">
						<span class="xp-level">Level 7</span>
						<span class="xp-amount">2,140 / 3,000 XP</span>
					</div>
					<div class="xp-track">
						<div class="xp-fill" style="width: 71%"></div>
					</div>
				</div>
			</div>

			<div class="card">
				<h3 class="card-title">AI-reviewed flashcards</h3>
				<p class="card-desc">
					Slip-ups from the round get turned into flashcards, each with a
					plain-language explanation of what to say instead.
				</p>
				<div class="flashcard">
					<span class="flashcard-wrong">Ich habe gegeht nach Hause.</span>
					<span class="flashcard-correct">Ich bin nach Hause gegangen.</span>
					<p class="flashcard-explanation">
						Movement verbs like "gehen" use "sein", not "haben", as the
						auxiliary in the perfect tense.
					</p>
				</div>
			</div>

			<div class="card">
				<h3 class="card-title">Friends</h3>
				<p class="card-desc">
					Add friends and see who's climbing fastest. Send a request, get
					accepted, and play together whenever you're both online.
				</p>
				<ul class="friends-list">
					<li v-for="(friend, index) in friends" :key="friend.name" class="friend-row">
						<span class="friend-rank">{{ index + 1 }}</span>
						<span class="friend-name">{{ friend.name }}</span>
						<span class="friend-xp">{{ friend.xp.toLocaleString() }} XP</span>
					</li>
				</ul>
			</div>
		</div>
	</section>
</template>

<script setup lang="ts">
const friends = [
	{ name: 'jonas.k', xp: 5820 },
	{ name: 'mei_learns', xp: 4390 },
	{ name: 'you', xp: 2140 },
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
		grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
		gap: 20px;
	}

	.card {
		background: #ffffff;
		border: 1px solid rgba(118, 117, 134, 0.18);
		border-radius: 16px;
		padding: 28px;
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
		font-size: 14px;
		color: var(--color-on-surface);
	}

	.xp-amount {
		font-size: 12px;
		color: var(--color-on-surface-variant);
	}

	.xp-track {
		width: 100%;
		height: 10px;
		border-radius: var(--radius-pill);
		background: rgba(118, 117, 134, 0.15);
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
		border: 1px solid rgba(118, 117, 134, 0.18);
	}

	.flashcard-wrong {
		font-size: 14px;
		color: var(--color-on-surface-variant);
		text-decoration: line-through;
		text-decoration-color: #d1352b;
	}

	.flashcard-correct {
		font-size: 15px;
		font-weight: 700;
		color: var(--color-on-surface);
	}

	.flashcard-explanation {
		font-size: 13px;
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
		font-size: 13px;
		color: var(--color-primary-container);
		width: 16px;
	}

	.friend-name {
		flex: 1;
		font-weight: 600;
		font-size: 14px;
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
