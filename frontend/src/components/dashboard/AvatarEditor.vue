<template>
	<div class="avatar-editor">
		<input ref="fileInput" type="file" accept="image/*" class="file-input" tabindex="-1" aria-hidden="true" @change="onFileChosen" />
		<div class="buttons">
			<BaseButton size="sm" pill :disabled="busy" @click="fileInput?.click()">
				<BaseIcon name="upload" :size="14" />
				{{ profile?.avatarUrl ? t('avatar.change') : t('avatar.upload') }}
			</BaseButton>
			<BaseButton v-if="profile?.avatarUrl" variant="ghost" size="sm" pill :disabled="busy" :aria-label="t('avatar.removeLabel')" @click="remove">
				{{ t('avatar.remove') }}
			</BaseButton>
		</div>
		<p class="status" role="status" aria-live="polite">{{ busy ? t('avatar.saving') : '' }}</p>
		<BaseAlert v-if="error">{{ error }}</BaseAlert>
	</div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { deleteAvatar, uploadAvatar, type Profile } from '../../api/users';
import { MAX_PICTURE_BYTES, toSquareJpeg } from '../../utils/image';
import BaseAlert from '../base/BaseAlert.vue';
import BaseButton from '../base/BaseButton.vue';
import { t } from '../../i18n';
import BaseIcon from '../base/BaseIcon.vue';

defineProps<{ profile: Profile | null }>();
const emit = defineEmits<{ updated: [profile: Profile] }>();

const fileInput = ref<HTMLInputElement | null>(null);
const busy = ref(false);
const error = ref('');

async function save(action: () => Promise<Profile>, fallback: string) {
	busy.value = true;
	error.value = '';
	try {
		emit('updated', await action());
	} catch (err) {
		error.value = err instanceof Error ? err.message : fallback;
	} finally {
		busy.value = false;
	}
}

async function onFileChosen() {
	const file = fileInput.value?.files?.[0];
	// Reset, so choosing the same file again still triggers a change.
	if (fileInput.value) fileInput.value.value = '';
	if (!file) return;

	if (file.size > MAX_PICTURE_BYTES) {
		error.value = t('avatar.tooBig');
		return;
	}
	let picture: Blob;
	try {
		picture = await toSquareJpeg(file);
	} catch {
		error.value = t('avatar.unreadable');
		return;
	}
	await save(() => uploadAvatar(picture), t('avatar.uploadError'));
}

function remove() {
	save(deleteAvatar, t('avatar.removeError'));
}
</script>

<style scoped>
	.avatar-editor {
		display: grid;
		gap: 8px;
		justify-items: end;
	}

	.file-input {
		display: none;
	}

	.buttons {
		display: flex;
		gap: 8px;
	}

	.status {
		margin: 0;
		font-size: var(--text-xs);
		color: var(--color-outline);
	}

	.status:empty {
		display: none;
	}

	@media (max-width: 720px) {
		.avatar-editor {
			justify-items: start;
		}
	}
</style>
