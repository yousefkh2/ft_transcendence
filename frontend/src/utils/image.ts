// Size of the uploaded profile picture. Large enough for the biggest avatar (72 px) on sharp screens.
const AVATAR_SIZE = 256;

// Files above this are refused before decoding, so a huge photo cannot freeze the tab.
export const MAX_PICTURE_BYTES = 15 * 1024 * 1024;

// Crops the middle square of an image and scales it to a small JPEG, so the upload stays around 30 KB.
// Any format the browser can show works as input. Throws when the file is not a readable image.
export async function toSquareJpeg(file: Blob): Promise<Blob> {
	const bitmap = await createImageBitmap(file);
	const side = Math.min(bitmap.width, bitmap.height);
	const canvas = document.createElement('canvas');
	canvas.width = canvas.height = AVATAR_SIZE;
	const ctx = canvas.getContext('2d');
	if (!ctx) throw new Error('Canvas is not available.');

	// JPEG has no transparency: transparent parts become white instead of black.
	ctx.fillStyle = '#fff';
	ctx.fillRect(0, 0, AVATAR_SIZE, AVATAR_SIZE);
	ctx.drawImage(bitmap, (bitmap.width - side) / 2, (bitmap.height - side) / 2, side, side, 0, 0, AVATAR_SIZE, AVATAR_SIZE);
	bitmap.close();

	return new Promise((resolve, reject) => {
		canvas.toBlob((blob) => (blob ? resolve(blob) : reject(new Error('Could not process the picture.'))), 'image/jpeg', 0.9);
	});
}
