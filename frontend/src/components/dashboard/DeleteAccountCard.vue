<template>
	<BaseCard tag="section" class="delete-card" aria-labelledby="delete-account-title">
		<h2 id="delete-account-title">{{ t('deleteAccount.title') }}</h2>
		<p class="text">{{ t('deleteAccount.text') }}</p>

		<BaseButton v-if="!open" variant="outline" size="sm" pill class="start" @click="open = true">
			{{ t('deleteAccount.start') }}
		</BaseButton>

		<form v-else class="confirm" novalidate @submit.prevent="submit">
			<BaseInput
				v-model="password"
				:label="t('deleteAccount.confirmLabel')"
				type="password"
				autocomplete="current-password"
				:error="error"
			/>
			<div class="buttons">
				<BaseButton type="submit" variant="danger" size="sm" pill :disabled="deleting">
					{{ deleting ? t('deleteAccount.deleting') : t('deleteAccount.confirm') }}
				</BaseButton>
				<BaseButton variant="outline" size="sm" pill :disabled="deleting" @click="cancel">{{ t('common.cancel') }}</BaseButton>
			</div>
		</form>
	</BaseCard>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { ApiError } from '../../api/client';
import { deleteAccount } from '../../api/users';
import { logout } from '../../auth/logout';
import BaseButton from '../base/BaseButton.vue';
import BaseCard from '../base/BaseCard.vue';
import BaseInput from '../base/BaseInput.vue';
import { t } from '../../i18n';

const open = ref(false);
const password = ref('');
const error = ref('');
const deleting = ref(false);

function cancel() {
	open.value = false;
	password.value = '';
	error.value = '';
}

async function submit() {
	if (!password.value) {
		error.value = t('validation.enterPassword');
		return;
	}
	deleting.value = true;
	error.value = '';
	try {
		await deleteAccount(password.value);
		await logout();
	} catch (err) {
		error.value =
			err instanceof ApiError && err.status === 403 ? t('apiErrors.incorrectPassword') : t('deleteAccount.error');
	} finally {
		deleting.value = false;
	}
}
</script>

<style scoped>
	h2 {
		margin: 0 0 6px;
		font-size: var(--text-lg);
		font-weight: 800;
	}

	.text {
		margin: 0 0 16px;
		max-width: 560px;
		font-size: 13.5px;
		line-height: 1.5;
		color: var(--color-on-surface-variant);
	}

	.start {
		color: var(--color-error);
		border-color: rgba(209, 53, 43, 0.35);
	}

	.confirm {
		display: grid;
		gap: 14px;
		max-width: 360px;
	}

	.buttons {
		display: flex;
		flex-wrap: wrap;
		gap: 10px;
	}
</style>
