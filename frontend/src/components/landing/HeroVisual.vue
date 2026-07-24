<template>
	<div class="hero-visual" aria-hidden="true">
		<div class="glow"></div>
		<div class="bolt"></div>
		<span
			v-for="(lang, index) in languages"
			:key="lang"
			class="chip"
			:class="`chip-${index}`"
			:style="{ animationDelay: `${index * 0.6}s` }"
		>
			{{ lang }}
		</span>
	</div>
</template>

<script setup lang="ts">
withDefaults(defineProps<{ languages?: string[] }>(), {
	languages: () => ['DE', 'ES', 'JA'],
});
</script>

<style scoped>
	.hero-visual {
		position: relative;
		aspect-ratio: 4 / 3;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.glow {
		position: absolute;
		width: 65%;
		height: 65%;
		border-radius: 50%;
		background: radial-gradient(circle, #6063ee 0%, transparent 70%);
		filter: blur(24px);
		animation: glow-pulse 5s ease-in-out infinite;
	}

	.bolt {
		position: relative;
		width: 170px;
		height: 270px;
		background: linear-gradient(160deg, #fea619 0%, #ff5d8f 45%, #6063ee 100%);
		clip-path: polygon(55% 0%, 25% 55%, 45% 55%, 35% 100%, 80% 40%, 55% 40%, 70% 0%);
		box-shadow: 0 30px 60px rgba(70, 72, 212, 0.35);
		animation: bolt-float 4s ease-in-out infinite;
	}

	.chip {
		position: absolute;
		display: inline-flex;
		align-items: center;
		justify-content: center;
		padding: 8px 16px;
		border-radius: var(--radius-pill);
		background: #ffffff;
		box-shadow: 0 8px 20px rgba(17, 28, 45, 0.14);
		font-family: var(--font-mono);
		font-weight: 700;
		font-size: 14px;
		color: #4648d4;
		animation: chip-bob 3.6s ease-in-out infinite;
	}

	.chip-0 {
		top: 8%;
		left: 4%;
	}

	.chip-1 {
		top: 42%;
		right: 2%;
	}

	.chip-2 {
		bottom: 10%;
		left: 10%;
	}

	@keyframes glow-pulse {
		0%, 100% {
			transform: scale(1);
			opacity: 0.5;
		}
		50% {
			transform: scale(1.08);
			opacity: 0.85;
		}
	}

	@keyframes bolt-float {
		0%, 100% {
			transform: translateY(0) rotate(-6deg);
		}
		50% {
			transform: translateY(-14px) rotate(-3deg);
		}
	}

	@keyframes chip-bob {
		0%, 100% {
			transform: translateY(0);
		}
		50% {
			transform: translateY(-10px);
		}
	}

	@media (prefers-reduced-motion: reduce) {
		.glow,
		.bolt,
		.chip {
			animation-play-state: paused;
		}
	}
</style>
