import router from '../router';
import { clearToken } from './session';

export function logout() {
	clearToken();
	return router.replace({ name: 'landing' });
}
