<template>
	<main class="not-found">
		<div class="bolt-mark"></div>
		<span class="code">404</span>
		<h1>Page not found</h1>
		<p>The page you're looking for doesn't exist or has moved.</p>
		<RouterLink to="/" class="home-btn">Back to home</RouterLink>
	</main>
</template>

<style scoped>
	.not-found {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		gap: 12px;
		min-height: 100vh;
		padding: 32px 16px;
		text-align: center;
		background: var(--color-surface);
		color: var(--color-on-surface);
		font-family: var(--font-display);
	}

	.bolt-mark {
		width: 30px;
		height: 48px;
		margin-bottom: 8px;
		background: linear-gradient(160deg, var(--color-secondary-container) 0%, var(--color-tertiary) 45%, var(--color-primary-container) 100%);
		clip-path: polygon(55% 0%, 25% 55%, 45% 55%, 35% 100%, 80% 40%, 55% 40%, 70% 0%);
	}

	.code {
		font-family: var(--font-mono);
		font-size: 13px;
		font-weight: 700;
		letter-spacing: 0.1em;
		color: var(--color-primary);
	}

	h1 {
		margin: 0;
		font-size: 28px;
		line-height: 0.95;
		font-weight: 800;
	}

	p {
		margin: 0;
		color: var(--color-on-surface-variant);
	}

	.home-btn {
		margin-top: 16px;
		padding: 10px 20px;
		border-radius: 12px;
		background: var(--color-primary);
		color: var(--color-on-primary);
		font-weight: 700;
		font-size: 14px;
		text-decoration: none;
	}
</style>
