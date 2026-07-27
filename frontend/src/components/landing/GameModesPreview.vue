<template>
	<section id="game-modes" class="game-modes">
		<div class="header">
			<span class="eyebrow">Game Modes</span>
			<h2 class="title">One scenario live. Many more landing.</h2>
		</div>

		<div class="card">
			<div class="copy">
				<h3 class="mode-title">Apartment Setup</h3>
				<p class="mode-desc">
					One player sees the room, the other reads the whiteboard. Move the
					furniture exactly where the clues say — before the 60-second clock
					runs out.
				</p>
				<div class="roles">
					<div class="role-chip" v-for="role in roles" :key="role.name">
						<span class="role-name">{{ role.name }}</span>
						<span class="role-desc">{{ role.description }}</span>
					</div>
				</div>
			</div>

			<div class="split-screen">
				<div class="pane pane-onsite">
					<img :src="livingroomBg" alt="" class="pane-bg" />
					<img :src="lamp" alt="" class="furniture furniture-lamp" />
					<img :src="sofa" alt="" class="furniture furniture-sofa" />
					<img :src="plant" alt="" class="furniture furniture-plant" />
					<span class="pane-label">On-Site</span>
				</div>
				<div class="divider"></div>
				<div class="pane pane-mc">
					<img :src="missionControlBg" alt="" class="pane-bg" />
					<span class="clue clue-1">plant → right of sofa</span>
					<span class="clue clue-2">lamp → left of sofa</span>
					<span class="pane-label">Mission Control</span>
				</div>
			</div>
		</div>

		<div class="future-modes">
			<span class="soon-chip" v-for="mode in futureModes" :key="mode">
				{{ mode }} <em>Soon</em>
			</span>
		</div>
	</section>
</template>

<script setup lang="ts">
import livingroomBg from '../../images/livingroom/livingroom_background.png';
import missionControlBg from '../../images/livingroom/mission_control_background.png';
import sofa from '../../images/livingroom/sofa.png';
import plant from '../../images/livingroom/plant.png';
import lamp from '../../images/livingroom/lamp.png';

const roles = [
	{ name: 'On-Site', description: 'Sees the room, moves the furniture.' },
	{ name: 'Mission Control', description: 'Reads the whiteboard, calls out directions.' },
];

const futureModes = ['Market Haggling', "Doctor's Visit", 'Job Interview'];
</script>

<style scoped>
	.game-modes {
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

	.card {
		display: grid;
		grid-template-columns: 1fr 1fr;
		align-items: center;
		gap: 40px;
		background: #111c2d;
		border-radius: 24px;
		padding: 48px;
	}

	.dot {
		width: 6px;
		height: 6px;
		border-radius: 50%;
		background: #4ade80;
	}

	.mode-title {
		font-size: clamp(24px, 3vw, 32px);
		font-weight: 800;
		color: #ffffff;
		margin: 0 0 12px;
	}

	.mode-desc {
		max-width: 420px;
		color: rgba(255, 255, 255, 0.65);
		font-size: 15px;
		line-height: 1.6;
		margin: 0 0 24px;
	}

	.roles {
		display: flex;
		flex-direction: column;
		gap: 10px;
	}

	.role-chip {
		display: flex;
		flex-wrap: wrap;
		align-items: baseline;
		gap: 10px;
		padding: 10px 16px;
		border-radius: 12px;
		background: rgba(255, 255, 255, 0.06);
		border: 1px solid rgba(255, 255, 255, 0.1);
	}

	.role-name {
		font-family: var(--font-mono);
		font-weight: 700;
		font-size: 13px;
		letter-spacing: 0.04em;
		text-transform: uppercase;
		color: var(--color-secondary-container);
	}

	.role-desc {
		font-size: 13.5px;
		color: rgba(255, 255, 255, 0.6);
	}

	.split-screen {
		position: relative;
		display: flex;
		aspect-ratio: 8 / 5;
		border-radius: 16px;
		overflow: hidden;
		box-shadow: 0 30px 60px rgba(0, 0, 0, 0.35);
	}

	.pane {
		position: relative;
		flex: 1;
		overflow: hidden;
	}

	.pane-bg {
		position: absolute;
		inset: 0;
		width: 100%;
		height: 100%;
		object-fit: cover;
	}

	.pane-onsite .pane-bg {
		object-position: 35% 65%;
	}

	.pane-mc .pane-bg {
		object-position: center 30%;
	}

	.divider {
		width: 3px;
		background: rgba(255, 255, 255, 0.25);
		z-index: 2;
	}

	.pane-label {
		position: absolute;
		bottom: 8px;
		left: 50%;
		transform: translateX(-50%);
		z-index: 3;
		padding: 3px 10px;
		border-radius: var(--radius-pill);
		background: rgba(0, 0, 0, 0.45);
		font-family: var(--font-mono);
		font-size: 10px;
		font-weight: 700;
		letter-spacing: 0.08em;
		text-transform: uppercase;
		color: #ffffff;
	}

	.furniture {
		position: absolute;
		image-rendering: pixelated;
		filter: drop-shadow(0 8px 10px rgba(0, 0, 0, 0.35));
	}

	.furniture-sofa {
		width: 46%;
		left: 28%;
		bottom: 12%;
	}

	.furniture-plant {
		width: 18%;
		left: 76%;
		bottom: 15%;
	}

	.furniture-lamp {
		width: 16%;
		left: 10%;
		bottom: 16%;
	}

	.clue {
		position: absolute;
		font-family: var(--font-mono);
		font-weight: 700;
		font-style: italic;
		font-size: clamp(11px, 1.6vw, 14px);
	}

	.clue-1 {
		top: 32%;
		left: 10%;
		color: #d1352b;
		transform: rotate(-3deg);
	}

	.clue-2 {
		top: 50%;
		left: 14%;
		color: #2854c4;
		transform: rotate(2deg);
	}

	.future-modes {
		display: flex;
		flex-wrap: wrap;
		justify-content: center;
		gap: 12px;
		margin-top: 32px;
	}

	.soon-chip {
		display: inline-flex;
		align-items: center;
		gap: 8px;
		padding: 10px 20px;
		border-radius: var(--radius-pill);
		border: 1.5px dashed var(--color-outline);
		color: var(--color-on-surface-variant);
		font-size: 14px;
		font-weight: 600;
	}

	.soon-chip em {
		font-style: normal;
		font-family: var(--font-mono);
		font-size: 10px;
		font-weight: 700;
		letter-spacing: 0.08em;
		text-transform: uppercase;
		color: var(--color-primary);
	}

	@media (max-width: 900px) {
		.game-modes {
			padding: 48px 16px;
		}

		.card {
			grid-template-columns: 1fr;
			padding: 32px 24px;
		}
	}
</style>
