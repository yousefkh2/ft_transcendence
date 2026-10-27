<template>
	<BaseModal v-model:open="open" :labelledby="titleId">
		<div class="glow glow-top"></div>
		<div class="glow glow-bottom"></div>

		<div ref="modalEl" class="modal-content">
			<div class="header">
				<BrandBolt class="bolt-mark" />
				<h2 :id="titleId" class="title">{{ tab === 'login' ? 'Welcome back' : 'Create your account' }}</h2>
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
					:aria-pressed="tab === 'login'"
					@click="switchTab('login')"
				>
					Log In
				</button>
				<button
					type="button"
					class="tab-btn"
					:class="{ active: tab === 'register' }"
					:aria-pressed="tab === 'register'"
					@click="switchTab('register')"
				>
					Sign Up
				</button>
			</div>

			<div class="form-area">
				<Transition :name="direction === 'forward' ? 'slide-fwd' : 'slide-back'">
					<form v-if="tab === 'login'" key="login" class="auth-form" novalidate @submit.prevent="handleLogin">
						<BaseAlert v-if="loginError">{{ loginError }}</BaseAlert>
						<BaseInput
							v-model="loginForm.email"
							label="Email"
							type="email"
							autocomplete="email"
							required
							placeholder="you@example.com"
							:error="loginErrors.email"
						/>
						<BaseInput
							v-model="loginForm.password"
							label="Password"
							type="password"
							autocomplete="current-password"
							required
							placeholder="••••••••"
							:error="loginErrors.password"
						/>
						<a href="#" class="forgot-link" @click.prevent="showForgotNotice = true">Forgot password?</a>
						<BaseAlert v-if="showForgotNotice" variant="info" class="forgot-notice">
							Password reset isn't available yet - hang tight, it's coming soon.
						</BaseAlert>
						<BaseButton type="submit" variant="gradient" size="lg" block class="submit-btn" :disabled="loginLoading">
							{{ loginLoading ? 'Logging in…' : 'Log In' }}
						</BaseButton>
						<p class="switch-line">
							Don't have an account?
							<a href="#" @click.prevent="switchTab('register')">Sign Up</a>
						</p>
					</form>
					<form v-else key="register" class="auth-form" novalidate @submit.prevent="handleRegister">
						<BaseAlert v-if="registerError">{{ registerError }}</BaseAlert>
						<BaseInput
							v-model="registerForm.username"
							label="Username"
							type="text"
							autocomplete="username"
							required
							placeholder="yourname"
							:error="registerErrors.username"
						/>
						<BaseInput
							v-model="registerForm.email"
							label="Email"
							type="email"
							autocomplete="email"
							required
							placeholder="you@example.com"
							:error="registerErrors.email"
						/>
						<BaseInput
							v-model="registerForm.password"
							label="Password"
							type="password"
							autocomplete="new-password"
							required
							placeholder="••••••••"
							:hint="PASSWORD_HINT"
							:error="registerErrors.password"
						/>
						<BaseSelect v-model="registerForm.language" label="Learning language" :options="LEARNING_LANGUAGES" />
						<BaseButton type="submit" variant="gradient" size="lg" block class="submit-btn" :disabled="registerLoading">
							{{ registerLoading ? 'Creating account…' : 'Create Account' }}
						</BaseButton>
						<p class="legal-note">
							By creating an account, you agree to our
							<RouterLink to="/terms">Terms of Service</RouterLink> and
							<RouterLink to="/privacy">Privacy Policy</RouterLink>.
						</p>
						<p class="switch-line">
							Already have an account?
							<a href="#" @click.prevent="switchTab('login')">Log In</a>
						</p>
					</form>
				</Transition>
			</div>
		</div>
	</BaseModal>
</template>

<script setup lang="ts">
import { nextTick, reactive, ref, useId, watch } from 'vue';
import { useRouter } from 'vue-router';
import { login, register } from '../../api/auth';
import {
	PASSWORD_HINT,
	validateEmail,
	validateLoginPassword,
	validateNewPassword,
	validateUsername,
} from '../../validation/rules';
import BaseAlert from '../base/BaseAlert.vue';
import BaseInput from '../base/BaseInput.vue';
import BaseSelect from '../base/BaseSelect.vue';
import { DEFAULT_LEARNING_LANGUAGE, LEARNING_LANGUAGES } from '../../constants/languages';
import BaseButton from '../base/BaseButton.vue';
import BaseModal from '../base/BaseModal.vue';
import BrandBolt from '../base/BrandBolt.vue';

const router = useRouter();

const open = defineModel<boolean>('open', { default: false });
const tab = defineModel<'login' | 'register'>('tab', { default: 'login' });

const direction = ref<'forward' | 'backward'>('forward');
const showForgotNotice = ref(false);

const modalEl = ref<HTMLElement | null>(null);
const titleId = useId();

const loginForm = reactive({ email: '', password: '' });
const loginErrors = reactive({ email: '', password: '' });
const loginSubmitted = ref(false);
const loginError = ref('');
const loginLoading = ref(false);

const registerForm = reactive({ username: '', email: '', password: '', language: DEFAULT_LEARNING_LANGUAGE });
const registerErrors = reactive({ username: '', email: '', password: '' });
const registerSubmitted = ref(false);
const registerError = ref('');
const registerLoading = ref(false);

function validateLoginForm() {
	loginErrors.email = validateEmail(loginForm.email);
	loginErrors.password = validateLoginPassword(loginForm.password);
	return !loginErrors.email && !loginErrors.password;
}

function validateRegisterForm() {
	registerErrors.username = validateUsername(registerForm.username);
	registerErrors.email = validateEmail(registerForm.email);
	registerErrors.password = validateNewPassword(registerForm.password);
	return !registerErrors.username && !registerErrors.email && !registerErrors.password;
}

// After a first submit attempt, errors update live while the user fixes them.
watch(loginForm, () => loginSubmitted.value && validateLoginForm());
watch(registerForm, () => registerSubmitted.value && validateRegisterForm());

function focusFirstInvalidField() {
	nextTick(() => modalEl.value?.querySelector<HTMLElement>('[aria-invalid="true"]')?.focus());
}

function resetValidation() {
	loginSubmitted.value = false;
	registerSubmitted.value = false;
	Object.assign(loginErrors, { email: '', password: '' });
	Object.assign(registerErrors, { username: '', email: '', password: '' });
}

function switchTab(next: 'login' | 'register') {
	if (next === tab.value) return;
	direction.value = next === 'register' ? 'forward' : 'backward';
	tab.value = next;
	loginError.value = '';
	registerError.value = '';
	showForgotNotice.value = false;
	resetValidation();
}

function close() {
	open.value = false;
}

// However the dialog was closed (button, Escape, click outside, success), it opens clean next time.
watch(open, (isOpen) => {
	if (isOpen) return;
	loginError.value = '';
	registerError.value = '';
	showForgotNotice.value = false;
	resetValidation();
});

async function handleLogin() {
	loginError.value = '';
	loginSubmitted.value = true;
	if (!validateLoginForm()) {
		focusFirstInvalidField();
		return;
	}
	loginLoading.value = true;
	try {
		await login(loginForm.email.trim(), loginForm.password);
		close();
		router.push('/dashboard');
	} catch (err) {
		loginError.value = err instanceof Error ? err.message : 'Something went wrong. Please try again.';
	} finally {
		loginLoading.value = false;
	}
}

async function handleRegister() {
	registerError.value = '';
	registerSubmitted.value = true;
	if (!validateRegisterForm()) {
		focusFirstInvalidField();
		return;
	}
	registerLoading.value = true;
	const username = registerForm.username.trim();
	const email = registerForm.email.trim();
	try {
		await register(username, email, registerForm.password, registerForm.language);
		await login(email, registerForm.password);
		close();
		router.push('/dashboard');
	} catch (err) {
		registerError.value = err instanceof Error ? err.message : 'Something went wrong. Please try again.';
	} finally {
		registerLoading.value = false;
	}
}
</script>

<style scoped>
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

	.modal-content {
		position: relative;
		z-index: 1;
	}

	.header {
		text-align: center;
		margin-bottom: 28px;
	}

	.bolt-mark {
		margin: 0 auto 16px;
		filter: drop-shadow(0 8px 14px rgba(70, 72, 212, 0.35));
	}

	.title {
		font-size: var(--text-xl);
		font-weight: 800;
		color: var(--color-on-surface);
		margin: 0 0 6px;
	}

	.subtitle {
		font-size: var(--text-md);
		color: var(--color-on-surface-variant);
		margin: 0;
	}

	.tab-switcher {
		position: relative;
		display: flex;
		padding: 4px;
		margin-bottom: 24px;
		border-radius: 13px;
		background: var(--color-fill-muted);
	}

	.tab-indicator {
		position: absolute;
		top: 4px;
		left: 4px;
		width: calc(50% - 4px);
		height: calc(100% - 8px);
		border-radius: 10px;
		background: var(--color-surface-raised);
		box-shadow: var(--shadow-sm);
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
		font-size: var(--text-md);
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

	.forgot-notice {
		margin-top: -6px;
	}

	.forgot-link {
		align-self: flex-end;
		margin-top: -4px;
		font-size: var(--text-sm);
		font-weight: 600;
		color: var(--color-primary);
		text-decoration: none;
	}

	.forgot-link:hover {
		text-decoration: underline;
	}

	.submit-btn {
		margin-top: 4px;
	}

	.legal-note {
		margin: 0;
		text-align: center;
		font-size: var(--text-xs);
		line-height: 1.5;
		color: var(--color-outline);
	}

	.legal-note a {
		color: var(--color-on-surface-variant);
		font-weight: 600;
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
		.glow {
			animation-play-state: paused;
		}
	}
</style>
