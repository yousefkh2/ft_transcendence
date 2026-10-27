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
					<span v-if="profile" class="xp-value">{{ profile.xp.toLocaleString() }} XP</span>
					<BaseSkeleton v-else width="44px" height="0.8em" />
				</div>
				<RouterLink :to="{ name: 'profile' }" class="avatar-link" aria-label="Your profile">
					<BaseAvatar v-if="profile" :name="profile.username" :src="apiAssetUrl(profile.avatarUrl)" size="md" variant="gradient" />
					<BaseSkeleton v-else shape="circle" width="36px" height="36px" />
				</RouterLink>
				<button type="button" class="logout-btn" @click="logout"><BaseIcon name="logout" :size="16" />Log out</button>
			</div>
		</header>

		<main id="main" class="dashboard-main" tabindex="-1">
			<BaseAlert v-if="error" size="md" @retry="loadProfile">{{ error }}</BaseAlert>
			<RouterView v-else :profile="profile" @profile-updated="profile = $event" />
		</main>

		<footer class="dashboard-footer">
			<RouterLink to="/privacy">Privacy Policy</RouterLink>
			<RouterLink to="/terms">Terms of Service</RouterLink>
		</footer>
	</div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { apiAssetUrl } from '../../api/client';
import { fetchProfile, type Profile } from '../../api/users';
import { logout } from '../../auth/logout';
import { useHeartbeat } from '../../composables/useHeartbeat';
import BaseAlert from '../base/BaseAlert.vue';
import BaseAvatar from '../base/BaseAvatar.vue';
import BaseIcon from '../base/BaseIcon.vue';
import BaseSkeleton from '../base/BaseSkeleton.vue';

useHeartbeat();

const profile = ref<Profile | null>(null);
const error = ref('');

async function loadProfile() {
	error.value = '';
	try {
		profile.value = await fetchProfile();
	} catch {
		error.value = 'Your profile could not be loaded right now.';
	}
}

onMounted(loadProfile);
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
		border-bottom: 1px solid var(--color-border);
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
		font-size: var(--text-lg);
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
		/* Only these, so the keyboard focus ring appears instantly instead of being animated. */
		transition: color 0.15s ease, background 0.15s ease, box-shadow 0.15s ease;
	}

	.tab-btn.active {
		background: var(--color-surface-raised);
		color: var(--color-primary);
		font-weight: 800;
		box-shadow: var(--shadow-md);
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
		border: 1px solid var(--color-border-strong);
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
		display: inline-flex;
		align-items: center;
		gap: 6px;
		white-space: nowrap;
		border: none;
		background: none;
		color: var(--color-outline);
		font-family: inherit;
		font-weight: 700;
		font-size: var(--text-sm);
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
