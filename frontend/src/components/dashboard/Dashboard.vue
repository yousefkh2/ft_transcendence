<template>
	<div class="dashboard">
		<header class="dashboard-header">
			<div class="brand">
				<img src="../../images/logo-transparent.png" alt="LingoDash" class="brand-logo" />
				<span class="brand-name">LingoDash</span>
			</div>

			<nav class="tabs">
				<RouterLink :to="{ name: 'dashboard' }" class="tab-btn" exact-active-class="active">Overview</RouterLink>
				<RouterLink :to="{ name: 'friends' }" class="tab-btn" exact-active-class="active">Friends</RouterLink>
				<RouterLink :to="{ name: 'profile' }" class="tab-btn" exact-active-class="active">Profile</RouterLink>
			</nav>

			<div class="header-right">
				<div class="xp-badge">
					<span class="xp-dot"></span>
					<span class="xp-value">{{ xpLabel }}</span>
				</div>
				<RouterLink :to="{ name: 'profile' }" class="avatar-link" aria-label="Your profile">
					<BaseAvatar :name="profile?.username" size="md" variant="gradient" />
				</RouterLink>
				<button type="button" class="logout-btn" @click="logout">Log out</button>
			</div>
		</header>

		<main class="dashboard-main">
			<BaseAlert v-if="error" size="md">{{ error }}</BaseAlert>
			<RouterView v-else :profile="profile" @profile-updated="profile = $event" />
		</main>

		<footer class="dashboard-footer">
			<RouterLink to="/privacy">Privacy Policy</RouterLink>
			<RouterLink to="/terms">Terms of Service</RouterLink>
		</footer>
	</div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { fetchProfile, type Profile } from '../../api/users';
import { logout } from '../../auth/logout';
import BaseAlert from '../base/BaseAlert.vue';
import BaseAvatar from '../base/BaseAvatar.vue';

const profile = ref<Profile | null>(null);
const error = ref('');

const xpLabel = computed(() => (profile.value ? `${profile.value.xp.toLocaleString()} XP` : '… XP'));

onMounted(async () => {
	try {
		profile.value = await fetchProfile();
	} catch (err) {
		error.value = err instanceof Error ? err.message : 'Could not load your profile.';
	}
});
</script>

<style scoped>
	.dashboard {
		display: flex;
		flex-direction: column;
		min-height: 100vh;
		background: var(--color-surface);
		color: var(--color-on-surface);
		font-family: var(--font-display);
	}

	.dashboard-header {
		position: sticky;
		top: 0;
		z-index: 40;
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 24px;
		padding: 14px 32px;
		background: rgba(249, 249, 255, 0.85);
		backdrop-filter: blur(12px);
		-webkit-backdrop-filter: blur(12px);
		border-bottom: 1px solid rgba(118, 117, 134, 0.15);
	}

	.brand {
		display: flex;
		align-items: center;
		gap: 10px;
	}

	.brand-logo {
		width: 34px;
		height: 34px;
		object-fit: contain;
	}

	.brand-name {
		font-weight: 800;
		font-size: 17px;
		letter-spacing: -0.01em;
	}

	.tabs {
		display: flex;
		gap: 6px;
		background: rgba(118, 117, 134, 0.09);
		border-radius: 13px;
		padding: 4px;
	}

	.tab-btn {
		text-align: center;
		text-decoration: none;
		padding: 8px 18px;
		border: none;
		border-radius: 10px;
		background: transparent;
		color: var(--color-outline);
		font-family: inherit;
		font-weight: 700;
		font-size: 13.5px;
		cursor: pointer;
		transition: all 0.15s ease;
	}

	.tab-btn.active {
		background: var(--color-surface-raised);
		color: var(--color-primary);
		font-weight: 800;
		box-shadow: 0 3px 8px rgba(17, 28, 45, 0.09);
	}

	.header-right {
		display: flex;
		align-items: center;
		gap: 14px;
	}

	.xp-badge {
		display: flex;
		align-items: center;
		gap: 7px;
		padding: 7px 13px;
		border-radius: 9999px;
		background: var(--color-surface-raised);
		border: 1px solid rgba(118, 117, 134, 0.18);
	}

	.xp-dot {
		width: 7px;
		height: 7px;
		border-radius: 9999px;
		background: var(--color-secondary-container);
	}

	.xp-value {
		white-space: nowrap;
		font-family: var(--font-mono);
		font-size: 12.5px;
		font-weight: 700;
		color: var(--color-on-surface-variant);
	}

	.avatar-link {
		display: flex;
		text-decoration: none;
		border-radius: var(--radius-pill);
	}

	.avatar-link:focus-visible {
		outline: 2px solid var(--color-primary);
		outline-offset: 2px;
	}

	.logout-btn {
		white-space: nowrap;
		border: none;
		background: none;
		color: var(--color-outline);
		font-family: inherit;
		font-weight: 700;
		font-size: 13px;
		cursor: pointer;
		padding: 0;
	}

	.logout-btn:hover {
		color: var(--color-on-surface);
	}

	.dashboard-main {
		flex: 1;
		width: 100%;
		max-width: 1180px;
		margin: 0 auto;
		padding: 32px;
	}

	.dashboard-footer {
		display: flex;
		justify-content: center;
		gap: 20px;
		padding: 20px 32px 28px;
	}

	.dashboard-footer a {
		font-size: 12.5px;
		font-weight: 600;
		color: var(--color-outline);
		text-decoration: none;
	}

	.dashboard-footer a:hover {
		color: var(--color-on-surface);
	}

	@media (max-width: 720px) {
		.dashboard-header {
			flex-wrap: wrap;
			gap: 12px;
			padding: 12px 16px;
		}

		.tabs {
			order: 3;
			width: 100%;
		}

		.tab-btn {
			flex: 1;
		}

		.header-right {
			gap: 10px;
		}

		.dashboard-main {
			padding: 20px 16px;
		}

		.dashboard-footer {
			padding: 16px 16px 24px;
		}
	}
</style>
