<template>
	<Transition name="overlay">
		<div v-if="open" class="overlay" @click.self="close">
			<div class="modal">
				<div class="glow glow-top"></div>
				<div class="glow glow-bottom"></div>

				<button type="button" class="close-btn" aria-label="Close" @click="close">×</button>

				<div class="modal-content">
					<div class="header">
						<div class="bolt-mark"></div>
						<h2 class="title">{{ tab === 'login' ? 'Welcome back' : 'Create your account' }}</h2>
						<p class="subtitle">
							{{
								tab === 'login'
									? 'Log in to keep your streak going.'
									: 'Start learning through play - free forever.'
							}}
						</p>
					</div>

					<div class="tab-switcher">
						<div class="tab-indicator" :class="{ 'is-register': tab === 'register' }"></div>
						<button
							type="button"
							class="tab-btn"
							:class="{ active: tab === 'login' }"
							@click="switchTab('login')"
						>
							Log In
						</button>
						<button
							type="button"
							class="tab-btn"
							:class="{ active: tab === 'register' }"
							@click="switchTab('register')"
						>
							Sign Up
						</button>
					</div>

					<div class="form-area">
						<Transition :name="direction === 'forward' ? 'slide-fwd' : 'slide-back'">
							<form v-if="tab === 'login'" key="login" class="auth-form" @submit.prevent>
								<label class="field">
									<span class="field-label">Email</span>
									<input type="email" required placeholder="you@example.com" />
								</label>
								<label class="field">
									<span class="field-label">Password</span>
									<input type="password" required placeholder="••••••••" />
								</label>
								<a href="#" class="forgot-link">Forgot password?</a>
								<button type="submit" class="submit-btn">Log In</button>
								<p class="switch-line">
									Don't have an account?
									<a href="#" @click.prevent="switchTab('register')">Sign Up</a>
								</p>
							</form>
							<form v-else key="register" class="auth-form" @submit.prevent>
								<label class="field">
									<span class="field-label">Username</span>
									<input type="text" required placeholder="yourname" />
								</label>
								<label class="field">
									<span class="field-label">Email</span>
									<input type="email" required placeholder="you@example.com" />
								</label>
								<label class="field">
									<span class="field-label">Password</span>
									<input type="password" required placeholder="••••••••" />
								</label>
								<label class="field">
									<span class="field-label">Learning language</span>
									<select v-model="learningLanguage">
										<option value="de">German</option>
										<option value="es" disabled>Spanish — coming soon</option>
										<option value="fr" disabled>French — coming soon</option>
										<option value="ja" disabled>Japanese — coming soon</option>
									</select>
								</label>
								<button type="submit" class="submit-btn">Create Account</button>
								<p class="switch-line">
									Already have an account?
									<a href="#" @click.prevent="switchTab('login')">Log In</a>
								</p>
							</form>
						</Transition>
					</div>
				</div>
			</div>
		</div>
	</Transition>
</template>

<script setup lang="ts">
import { ref } from 'vue';

const open = defineModel<boolean>('open', { default: false });
const tab = defineModel<'login' | 'register'>('tab', { default: 'login' });

const direction = ref<'forward' | 'backward'>('forward');
const learningLanguage = ref('de');

function switchTab(next: 'login' | 'register') {
	if (next === tab.value) return;
	direction.value = next === 'register' ? 'forward' : 'backward';
	tab.value = next;
}

function close() {
	open.value = false;
}
</script>

<style scoped>
	.overlay {
		position: fixed;
		inset: 0;
		z-index: 200;
		display: flex;
		align-items: center;
		justify-content: center;
		padding: 16px;
		background: rgba(10, 12, 28, 0.6);
		backdrop-filter: blur(6px);
		-webkit-backdrop-filter: blur(6px);
	}

	.overlay-enter-active {
		animation: overlay-in 0.22s ease;
	}

	.modal {
		position: relative;
		width: 100%;
		max-width: 400px;
		padding: 40px 36px;
		border-radius: 28px;
		overflow: hidden;
		background: rgba(255, 255, 255, 0.86);
		backdrop-filter: blur(20px);
		-webkit-backdrop-filter: blur(20px);
		border: 1px solid rgba(255, 255, 255, 0.5);
		box-shadow: 0 50px 100px rgba(17, 28, 45, 0.4);
		animation: modal-in 0.32s cubic-bezier(0.16, 1, 0.3, 1);
	}

	.glow {
		position: absolute;
		width: 55%;
		aspect-ratio: 1;
		border-radius: 50%;
		filter: blur(30px);
		pointer-events: none;
		z-index: 0;
		animation: glow-pulse 6s ease-in-out infinite;
	}

	.glow-top {
		top: -15%;
		right: -15%;
		background: radial-gradient(circle, rgba(96, 99, 238, 0.18) 0%, transparent 70%);
	}

	.glow-bottom {
		bottom: -15%;
		left: -15%;
		background: radial-gradient(circle, rgba(254, 166, 25, 0.14) 0%, transparent 70%);
		animation-delay: 3s;
	}

	.close-btn {
		position: absolute;
		top: 16px;
		right: 16px;
		z-index: 2;
		width: 32px;
		height: 32px;
		display: flex;
		align-items: center;
		justify-content: center;
		border: none;
		border-radius: 50%;
		background: rgba(17, 28, 45, 0.06);
		color: var(--color-on-surface-variant);
		font-size: 18px;
		line-height: 1;
		cursor: pointer;
		transition: background 0.15s ease;
	}

	.close-btn:hover {
		background: rgba(17, 28, 45, 0.12);
	}

	.modal-content {
		position: relative;
		z-index: 1;
	}

	.header {
		text-align: center;
		margin-bottom: 28px;
	}

	.bolt-mark {
		width: 30px;
		height: 48px;
		margin: 0 auto 16px;
		background: linear-gradient(160deg, #fea619 0%, #ff5d8f 45%, #6063ee 100%);
		clip-path: polygon(55% 0%, 25% 55%, 45% 55%, 35% 100%, 80% 40%, 55% 40%, 70% 0%);
		filter: drop-shadow(0 8px 14px rgba(70, 72, 212, 0.35));
	}

	.title {
		font-size: 24px;
		font-weight: 800;
		color: var(--color-on-surface);
		margin: 0 0 6px;
	}

	.subtitle {
		font-size: 14px;
		color: var(--color-on-surface-variant);
		margin: 0;
	}

	.tab-switcher {
		position: relative;
		display: flex;
		padding: 4px;
		margin-bottom: 24px;
		border-radius: 13px;
		background: rgba(118, 117, 134, 0.1);
	}

	.tab-indicator {
		position: absolute;
		top: 4px;
		left: 4px;
		width: calc(50% - 4px);
		height: calc(100% - 8px);
		border-radius: 10px;
		background: #ffffff;
		box-shadow: 0 1px 3px rgba(17, 28, 45, 0.12);
		transition: transform 0.28s cubic-bezier(0.2, 0.8, 0.3, 1);
	}

	.tab-indicator.is-register {
		transform: translateX(100%);
	}

	.tab-btn {
		position: relative;
		z-index: 1;
		flex: 1;
		padding: 10px 0;
		border: none;
		border-radius: 10px;
		background: transparent;
		font-family: var(--font-display);
		font-weight: 700;
		font-size: 14px;
		color: var(--color-on-surface-variant);
		cursor: pointer;
		transition: color 0.2s ease;
	}

	.tab-btn.active {
		color: var(--color-on-surface);
	}

	.form-area {
		position: relative;
		display: grid;
		overflow: hidden;
	}

	.auth-form {
		grid-area: 1 / 1;
		display: flex;
		flex-direction: column;
		gap: 14px;
	}

	.field {
		display: flex;
		flex-direction: column;
		gap: 6px;
	}

	.field-label {
		font-size: 13px;
		font-weight: 600;
		color: var(--color-on-surface-variant);
	}

	.field input,
	.field select {
		height: 47px;
		padding: 0 14px;
		border-radius: 12px;
		border: 1px solid rgba(118, 117, 134, 0.18);
		background: rgba(249, 249, 255, 0.8);
		font-family: var(--font-display);
		font-size: 14px;
		color: var(--color-on-surface);
		transition: border-color 0.2s ease;
	}

	.field input:focus,
	.field select:focus {
		outline: none;
		border-color: var(--color-primary);
	}

	.forgot-link {
		align-self: flex-end;
		margin-top: -4px;
		font-size: 13px;
		font-weight: 600;
		color: var(--color-primary);
		text-decoration: none;
	}

	.forgot-link:hover {
		text-decoration: underline;
	}

	.submit-btn {
		height: 47px;
		margin-top: 4px;
		border: none;
		border-radius: 12px;
		background: linear-gradient(135deg, #3739b8 0%, #4648d4 55%, #6063ee 100%);
		color: #ffffff;
		font-family: var(--font-display);
		font-weight: 700;
		font-size: 15px;
		cursor: pointer;
		box-shadow: 0 10px 20px -8px rgba(53, 55, 170, 0.55), 0 20px 40px -16px rgba(53, 55, 170, 0.4);
		transition: transform 0.15s ease, box-shadow 0.15s ease, filter 0.15s ease;
	}

	.submit-btn:hover {
		transform: translateY(-2px);
		filter: brightness(1.08);
		box-shadow: 0 14px 26px -8px rgba(53, 55, 170, 0.6), 0 26px 48px -16px rgba(53, 55, 170, 0.48);
	}

	.switch-line {
		margin: 4px 0 0;
		text-align: center;
		font-size: 13.5px;
		color: var(--color-on-surface-variant);
	}

	.switch-line a {
		margin-left: 4px;
		color: var(--color-primary);
		font-weight: 700;
		text-decoration: none;
	}

	.switch-line a:hover {
		text-decoration: underline;
	}

	.slide-fwd-enter-active,
	.slide-fwd-leave-active,
	.slide-back-enter-active,
	.slide-back-leave-active {
		transition: transform 0.32s cubic-bezier(0.2, 0.8, 0.3, 1), opacity 0.32s ease;
	}

	.slide-fwd-enter-from,
	.slide-back-leave-to {
		transform: translateX(28px);
		opacity: 0;
	}

	.slide-fwd-leave-to,
	.slide-back-enter-from {
		transform: translateX(-28px);
		opacity: 0;
	}

	@keyframes overlay-in {
		from {
			opacity: 0;
		}
		to {
			opacity: 1;
		}
	}

	@keyframes modal-in {
		from {
			opacity: 0;
			transform: translateY(24px) scale(0.96);
		}
		to {
			opacity: 1;
			transform: translateY(0) scale(1);
		}
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

	@media (prefers-reduced-motion: reduce) {
		.overlay-enter-active,
		.modal {
			animation: none;
		}

		.glow {
			animation-play-state: paused;
		}
	}
</style>
