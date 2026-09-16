<template>
	<header class="navbar">
		<div class="brand">
			<img :src="logo" alt="LingoDash Logo" class="brand-logo" />
			<span class="brand-name">LingoDash</span>
		</div>

		<nav class="nav-links">
			<a href="#how-it-works">How it Works</a>
			<a href="#game-modes">Games</a>
		</nav>

		<div class="auth-buttons">
			<button class="btn-ghost" type="button" @click="emit('open-auth', 'login')">Login</button>
			<button class="btn-login" type="button" @click="emit('open-auth', 'register')">Sign Up</button>
		</div>
	</header>
</template>

<script setup lang="ts">
import logo from "../../images/logo-transparent.png"

const emit = defineEmits<{ 'open-auth': [tab: 'login' | 'register'] }>();
</script>

<style scoped>
	.navbar {
		position: sticky;
		top: 0;
		z-index: 50;
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 24px;
		padding: 16px 40px;
		background: rgba(249, 249, 255, 0.72);
		backdrop-filter: blur(12px);
		-webkit-backdrop-filter: blur(12px);
		border-bottom: 1px solid rgba(118, 117, 134, 0.15);
		width: 80vw;
		margin: auto;
	}

	.brand {
		display: flex;
		align-items: center;
		gap: 10px;
	}

	.brand-logo {
		width: 45px;
		height: 42px;
		border-radius: var(--radius-md);
	}

	.brand-name {
		font-family: var(--font-display);
		font-weight: 800;
		font-size: 18px;
		letter-spacing: -0.01em;
		color: var(--color-on-surface);
	}

	.nav-links {
		display: flex;
		align-items: center;
		gap: 32px;
	}

	.nav-links a {
		position: relative;
		color: var(--color-on-surface-variant);
		text-decoration: none;
		font-weight: 600;
		font-size: 14px;
		padding: 4px 0;
		transition: color 0.2s ease;
	}

	.nav-links a::after {
		content: "";
		position: absolute;
		left: 0;
		bottom: 0;
		width: 0%;
		height: 2px;
		background: var(--color-primary);
		transition: width 0.25s ease;
	}

	.nav-links a:hover {
		color: var(--color-on-surface);
	}

	.nav-links a:hover::after {
		width: 100%;
	}

	.auth-buttons {
		display: flex;
		align-items: center;
		gap: 8px;
	}

	.btn-ghost {
		font-family: var(--font-display);
		font-weight: 700;
		font-size: 14px;
		padding: 10px 18px;
		border-radius: var(--radius-pill);
		border: none;
		background: transparent;
		color: var(--color-on-surface);
		cursor: pointer;
		transition: background 0.15s ease;
	}

	.btn-ghost:hover {
		background: rgba(17, 28, 45, 0.06);
	}

	.btn-login {
		font-family: var(--font-display);
		font-weight: 700;
		font-size: 14px;
		padding: 10px 22px;
		border-radius: var(--radius-pill);
		border: none;
		background: var(--color-primary);
		color: var(--color-on-primary);
		cursor: pointer;
		transition: tansform 0.15s ease, box-shadow 0.15s ease;
	}

	.btn-login:hover {
		transform: translateY(-2px);
		box-shadow: 0 8px 20px rgba(70, 72, 212, 0.35);
	}

	.btn-login:active {
		transform: translateY(0);
		box-shadow: none;
	}

	@media (max-width: 720px) {
		.nav-links {
			display: none;;
		}
	}

</style>