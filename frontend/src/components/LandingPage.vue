<template>
	<main class="landing">
		<section class="hero">
			<div class="hero-copy">
				<span class="badge">Beta · Two-Player Mode</span>
				<h1 class="headline">
					Talk Fast, Build Together
					<span class="accent">In German.</span>
				</h1>
				<p class="subtext">
					One player sees the room. The other holds the plan. Describe it, argue
          			about it, and race the clock to furnish the apartment before time runs out.
				</p>
				<div class="cta-row">
					<button class="btn btn-primary">Play Now</button>
					<button class="btn btn-ghost">How It Works</button>
				</div>
			</div>
			<div class="hero-visual" aria-hidden="true"></div>
		</section>
	</main>
</template>

<script setup lang="ts">
</script>

<style scoped>
	.landing {
		--color-primary: #4648d4;
		--color-primary-container: #6063ee;
		--color-on-primary: #ffffff;
		--color-secondary-container: #fea619;
		--color-on-secondary-container: #684000;
		--color-surface: #f9f9ff;
		--color-on-surface: #111c2d;
		--color-on-surface-variant: #464554;
		--color-outline: #767586;

		--font-display: "Plus Jakarta Sans", sans-serif;
		--font-mono: "Space Grotesk", monospace;

		--radius-lg: 1rem;
		--radius-pill: 9999px;

		background: var(--color-surface);
		color: var(--color-on-surface);
		font-family: var(--font-display);
		min-height: 100vh;
	}

	.hero {
		display: grid;
		grid-template-columns: 1fr 1fr;
		align-items: center;
		gap: 48px;
		max-width: 1200px;
		margin: 0 auto;
		padding: 96px 40px;
	}

	.badge {
		display: inline-flex;
		padding: 6px 14px;
		border-radius: var(--radius-pill);
		background: var(--color-secondary-container);
		color: var(--color-on-secondary-container);
		font-family: (--font-mono);
		font-size: 12px;
		font-weight: 700;
		letter-spacing: 0.1em;
		text-transform: uppercase;
		margin-bottom: 20px;
	}

	.headling {
		font-size: clamp(32px, 5vw, 48px);
		font-weight: 800;
		line-height: 1.1;
		letter-spacing: -0.2em;
		margin: 0 0 20px;
	}

	.accent {
		color: var(--color-primary);
	}

	.subtext {
		max-width: 480px;
		color: var(--color-on-surface-variant);
		font-size: 16px;
		line-height: 1.6;
		margin: 0 0 32px;
	}

	.cta-row {
		display: flex;
		gap: 12px;
		flex-wrap: wrap;
	}

	.btn {
		font-family: var(--font-display);
		font-weight: 700;
		font-size: 15px;
		padding: 14px 24px;
		border-radius: var(--radius-pill);
		border: none;
		cursor: pointer;
	}

	.btn-primary {
		background: var(--color-primary);
		color: var(--color-on-primary);
		border-bottom: 4px solid #2f2ebe;
	}

	.btn-primary:active {
		border-bottom-width: 0;
		transform: translateY(4px);
	}

	.btn-ghost {
		background: transparent;
		color: var(--color-on-surface);
		border: 2px solid var(--color-outline);
	}

	.hero-visual {
		aspect-ratio: 4 / 3;
		border-radius: var(--radius-lg);
		background: linear-gradient(135deg, var(--color-primary-container), var(--color-primary));
	}

	@media (max-width: 900px) {
		.hero {
			grid-template-columns: 1fr;
			text-align: center;
			padding: 48px 16px;
		}

		.subtext {
			margin-left: auto;
			margin-right: auto;
		}

		.hero-visual {
			order: -1;
		}
	}


</style>