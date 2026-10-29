<template>
	<BaseCard tag="section" aria-labelledby="edit-profile-title">
		<h2 id="edit-profile-title">{{ t('editProfile.title') }}</h2>

		<div class="forms">
			<form class="form" novalidate @submit.prevent="saveDetails">
				<h3>{{ t('editProfile.details') }}</h3>
				<BaseAlert v-if="detailsNotice" variant="success">{{ detailsNotice }}</BaseAlert>
				<BaseAlert v-if="detailsError">{{ detailsError }}</BaseAlert>
				<BaseInput v-model="details.username" :label="t('common.username')" autocomplete="username" :error="detailsErrors.username" />
				<BaseInput v-model="details.email" :label="t('common.email')" type="email" autocomplete="email" :error="detailsErrors.email" />
				<BaseSelect v-model="details.learningLanguage" :label="t('common.learningLanguage')" :options="languageOptions" />
				<BaseInput
					v-if="emailChanged"
					v-model="details.currentPassword"
					:label="t('common.currentPassword')"
					type="password"
					autocomplete="current-password"
					:hint="t('editProfile.emailPasswordHint')"
					:error="detailsErrors.currentPassword"
				/>
				<BaseButton type="submit" size="sm" pill class="submit" :disabled="!detailsChanged || savingDetails">
					{{ savingDetails ? t('common.saving') : t('editProfile.save') }}
				</BaseButton>
			</form>

			<form class="form" novalidate @submit.prevent="savePassword">
				<h3>{{ t('editProfile.changePassword') }}</h3>
				<BaseAlert v-if="passwordNotice" variant="success">{{ passwordNotice }}</BaseAlert>
				<BaseAlert v-if="passwordError">{{ passwordError }}</BaseAlert>
				<BaseInput
					v-model="passwords.current"
					:label="t('common.currentPassword')"
					type="password"
					autocomplete="current-password"
					:error="passwordErrors.current"
				/>
				<BaseInput
					v-model="passwords.next"
					:label="t('editProfile.newPassword')"
					type="password"
					autocomplete="new-password"
					:hint="passwordHint()"
					:error="passwordErrors.next"
				/>
				<BaseInput
					v-model="passwords.repeat"
					:label="t('editProfile.repeatPassword')"
					type="password"
					autocomplete="new-password"
					:error="passwordErrors.repeat"
				/>
				<BaseButton type="submit" size="sm" pill class="submit" :disabled="savingPassword">
					{{ savingPassword ? t('common.saving') : t('editProfile.changePassword') }}
				</BaseButton>
			</form>
		</div>
	</BaseCard>
</template>

<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue';
import { ApiError } from '../../api/client';
import { changePassword, updateProfile, type Profile, type ProfileChanges } from '../../api/users';
import { passwordHint, validateEmail, validateNewPassword, validateUsername } from '../../validation/rules';
import BaseAlert from '../base/BaseAlert.vue';
import BaseButton from '../base/BaseButton.vue';
import BaseCard from '../base/BaseCard.vue';
import BaseInput from '../base/BaseInput.vue';
import BaseSelect from '../base/BaseSelect.vue';
import { learningLanguageOptions } from '../../constants/languages';
import { t } from '../../i18n';

const props = defineProps<{ profile: Profile | null }>();
const emit = defineEmits<{ updated: [profile: Profile] }>();

const languageOptions = computed(learningLanguageOptions);

const details = reactive({ username: '', email: '', learningLanguage: '', currentPassword: '' });
const detailsErrors = reactive({ username: '', email: '', currentPassword: '' });
const detailsNotice = ref('');
const detailsError = ref('');
const savingDetails = ref(false);

const passwords = reactive({ current: '', next: '', repeat: '' });
const passwordErrors = reactive({ current: '', next: '', repeat: '' });
const passwordNotice = ref('');
const passwordError = ref('');
const savingPassword = ref(false);

// Fill the form with the current values once the profile has loaded.
watch(
	() => props.profile,
	(p) => {
		if (!p) return;
		details.username = p.username;
		details.email = p.email;
		details.learningLanguage = p.learningLanguage;
	},
	{ immediate: true },
);

const usernameChanged = computed(() => !!props.profile && details.username.trim() !== props.profile.username);
const emailChanged = computed(() => !!props.profile && details.email.trim() !== props.profile.email);
const languageChanged = computed(() => !!props.profile && details.learningLanguage !== props.profile.learningLanguage);
const detailsChanged = computed(() => usernameChanged.value || emailChanged.value || languageChanged.value);

async function saveDetails() {
	detailsNotice.value = '';
	detailsError.value = '';
	detailsErrors.username = usernameChanged.value ? validateUsername(details.username) : '';
	detailsErrors.email = emailChanged.value ? validateEmail(details.email) : '';
	detailsErrors.currentPassword = emailChanged.value && !details.currentPassword ? t('validation.enterPassword') : '';
	if (detailsErrors.username || detailsErrors.email || detailsErrors.currentPassword) return;

	const changes: ProfileChanges = {};
	if (usernameChanged.value) changes.username = details.username.trim();
	if (languageChanged.value) changes.learningLanguage = details.learningLanguage;
	if (emailChanged.value) {
		changes.email = details.email.trim();
		changes.currentPassword = details.currentPassword;
	}

	savingDetails.value = true;
	try {
		const updated = await updateProfile(changes);
		details.currentPassword = '';
		detailsNotice.value = t('editProfile.updated');
		emit('updated', updated);
	} catch (err) {
		const serverMessage = err instanceof ApiError ? err.serverMessage : '';
		if (err instanceof ApiError && err.status === 403) detailsErrors.currentPassword = t('apiErrors.incorrectPassword');
		else if (serverMessage === 'username already taken') detailsErrors.username = t('apiErrors.usernameTaken');
		else if (serverMessage === 'email already registered') detailsErrors.email = t('apiErrors.emailTaken');
		else detailsError.value = err instanceof Error ? err.message : t('editProfile.updateError');
	} finally {
		savingDetails.value = false;
	}
}

async function savePassword() {
	passwordNotice.value = '';
	passwordError.value = '';
	passwordErrors.current = passwords.current ? '' : t('validation.enterCurrentPassword');
	passwordErrors.next = validateNewPassword(passwords.next);
	passwordErrors.repeat = passwords.repeat === passwords.next ? '' : t('validation.passwordsDontMatch');
	if (passwordErrors.current || passwordErrors.next || passwordErrors.repeat) return;

	savingPassword.value = true;
	try {
		await changePassword(passwords.current, passwords.next);
		Object.assign(passwords, { current: '', next: '', repeat: '' });
		passwordNotice.value = t('editProfile.passwordChanged');
	} catch (err) {
		if (err instanceof ApiError && err.status === 403) passwordErrors.current = t('apiErrors.incorrectPassword');
		else passwordError.value = err instanceof Error ? err.message : t('editProfile.passwordError');
	} finally {
		savingPassword.value = false;
	}
}
</script>

<style scoped>
	h2 {
		margin: 0 0 16px;
		font-size: var(--text-lg);
		font-weight: 800;
	}

	h3 {
		margin: 0;
		font-size: var(--text-md);
		font-weight: 700;
	}

	.forms {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
		gap: 28px;
	}

	.form {
		display: grid;
		align-content: start;
		gap: 14px;
	}

	.submit {
		justify-self: start;
	}
</style>
