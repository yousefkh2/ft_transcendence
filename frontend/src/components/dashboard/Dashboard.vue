<template>
	<div class="dashboard">
		<header class="dashboard-header">
			<div class="brand">
				<img src="../../images/logo-transparent.png" alt="LingoDash" class="brand-logo" />
				<span class="brand-name">LingoDash</span>
			</div>

			<nav class="tabs">
				<button
					type="button"
					class="tab-btn"
					:class="{ active: tab === 'overview' }"
					@click="tab = 'overview'"
				>
					Overview
				</button>
				<button
					type="button"
					class="tab-btn"
					:class="{ active: tab === 'friends' }"
					@click="tab = 'friends'"
				>
					Friends
				</button>
			</nav>

			<div class="header-right">
				<div class="xp-badge">
					<span class="xp-dot"></span>
					<span class="xp-value">{{ xpLabel }}</span>
				</div>
				<div class="avatar" :title="profile?.username">{{ avatarInitial }}</div>
				<button type="button" class="logout-btn" @click="logout">Log out</button>
			</div>
		</header>

		<main class="dashboard-main">
			<p v-if="error" class="load-error">{{ error }}</p>
			<OverviewTab v-else-if="tab === 'overview'" :profile="profile" @switch-tab="tab = $event" />
			<FriendsTab v-else-if="tab === 'friends'" />
		</main>
	</div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { fetchProfile, type Profile } from '../../api/dashboard';
import OverviewTab from './OverviewTab.vue';
import FriendsTab from './FriendsTab.vue';

const tab = ref<'overview' | 'friends'>('overview');
const profile = ref<Profile | null>(null);
const error = ref('');

const xpLabel = computed(() => (profile.value ? `${profile.value.xp.toLocaleString()} XP` : '… XP'));
const avatarInitial = computed(() => profile.value?.username?.[0]?.toUpperCase() ?? '?');

function logout() {
	localStorage.removeItem('lingodash_token');
	window.location.href = '/';
}

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
		--color-primary: #4648d4;
		--color-surface: #f9f9ff;
		--color-on-surface: #111c2d;
		--color-on-surface-variant: #464554;
		--color-muted: #767586;

		min-height: 100vh;
		background: var(--color-surface);
		color: var(--color-on-surface);
		font-family: 'Plus Jakarta Sans', sans-serif;
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
		padding: 8px 18px;
		border: none;
		border-radius: 10px;
		background: transparent;
		color: var(--color-muted);
		font-family: inherit;
		font-weight: 700;
		font-size: 13.5px;
		cursor: pointer;
		transition: all 0.15s ease;
	}

	.tab-btn.active {
		background: #fff;
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
		background: #fff;
		border: 1px solid rgba(118, 117, 134, 0.18);
	}

	.xp-dot {
		width: 7px;
		height: 7px;
		border-radius: 9999px;
		background: #fea619;
	}

	.xp-value {
		font-family: 'Space Grotesk', monospace;
		font-size: 12.5px;
		font-weight: 700;
		color: var(--color-on-surface-variant);
	}

	.avatar {
		width: 36px;
		height: 36px;
		border-radius: 9999px;
		background: linear-gradient(135deg, #4648d4, #6063ee);
		color: #fff;
		display: grid;
		place-items: center;
		font-weight: 800;
		font-size: 14px;
	}

	.logout-btn {
		border: none;
		background: none;
		color: var(--color-muted);
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
		max-width: 1180px;
		margin: 0 auto;
		padding: 32px;
	}

	.load-error {
		padding: 16px 20px;
		border-radius: 12px;
		background: rgba(209, 53, 43, 0.08);
		color: #d1352b;
		font-weight: 600;
		font-size: 14px;
	}
</style>
