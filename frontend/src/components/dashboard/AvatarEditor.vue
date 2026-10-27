<template>
	<div class="avatar-editor">
		<input ref="fileInput" type="file" accept="image/*" class="file-input" tabindex="-1" aria-hidden="true" @change="onFileChosen" />
		<div class="buttons">
			<BaseButton size="sm" pill :disabled="busy" @click="fileInput?.click()">
				{{ profile?.avatarUrl ? 'Change picture' : 'Upload picture' }}
			</BaseButton>
			<BaseButton v-if="profile?.avatarUrl" variant="ghost" size="sm" pill :disabled="busy" aria-label="Remove profile picture" @click="remove">
				Remove
			</BaseButton>
		</div>
		<p class="status" role="status" aria-live="polite">{{ busy ? 'Saving picture...' : '' }}</p>
		<BaseAlert v-if="error">{{ error }}</BaseAlert>
	</div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { deleteAvatar, uploadAvatar, type Profile } from '../../api/users';
import { MAX_PICTURE_BYTES, toSquareJpeg } from '../../utils/image';
import BaseAlert from '../base/BaseAlert.vue';
import BaseButton from '../base/BaseButton.vue';

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
		error.value = 'This file is too big. Please choose a picture under 15 MB.';
		return;
	}
	let picture: Blob;
	try {
		picture = await toSquareJpeg(file);
	} catch {
		error.value = 'This file could not be read as a picture. Try a JPEG or PNG.';
		return;
	}
	await save(() => uploadAvatar(picture), 'Could not upload your picture.');
}

function remove() {
	save(deleteAvatar, 'Could not remove your picture.');
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
		font-size: 12px;
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
