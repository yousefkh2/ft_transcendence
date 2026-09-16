<template>
	<main class="dashboard-placeholder">
		<div class="card">
			<div class="bolt-mark"></div>
			<h1>You're in!</h1>
			<p>Logged in successfully. The dashboard itself is still being built.</p>
			<button type="button" class="logout-btn" @click="logout">Log out</button>
		</div>
	</main>
</template>

<script setup lang="ts">
function logout() {
	localStorage.removeItem('lingodash_token');
	window.location.href = '/';
}
</script>

<style scoped>
	.dashboard-placeholder {
		--color-primary: #4648d4;
		--color-surface: #f9f9ff;
		--color-on-surface: #111c2d;
		--color-on-surface-variant: #464554;

		--font-display: "Plus Jakarta Sans", sans-serif;

		display: flex;
		align-items: center;
		justify-content: center;
		min-height: 100vh;
		background: var(--color-surface);
		color: var(--color-on-surface);
		font-family: var(--font-display);
	}

	.card {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 12px;
		text-align: center;
		padding: 48px 40px;
	}

	.bolt-mark {
		width: 30px;
		height: 48px;
		margin-bottom: 8px;
		background: linear-gradient(160deg, #fea619 0%, #ff5d8f 45%, #6063ee 100%);
		clip-path: polygon(55% 0%, 25% 55%, 45% 55%, 35% 100%, 80% 40%, 55% 40%, 70% 0%);
	}

	h1 {
		margin: 0;
		font-size: 28px;
		font-weight: 800;
	}

	p {
		margin: 0;
		color: var(--color-on-surface-variant);
	}

	.logout-btn {
		margin-top: 16px;
		padding: 10px 20px;
		border: none;
		border-radius: 12px;
		background: var(--color-primary);
		color: #ffffff;
		font-family: var(--font-display);
		font-weight: 700;
		font-size: 14px;
		cursor: pointer;
	}
</style>
