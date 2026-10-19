import { createRouter, createWebHistory } from 'vue-router';
import LandingPage from '../components/LandingPage.vue';
import { hasValidSession } from '../auth/session';

declare module 'vue-router' {
	interface RouteMeta {
		requiresAuth?: boolean;
		guestOnly?: boolean;
	}
}

const router = createRouter({
	history: createWebHistory(),
	routes: [
		{ path: '/', name: 'landing', component: LandingPage, meta: { guestOnly: true } },
		{ path: '/landingpage', redirect: '/' },
		{
			path: '/dashboard',
			name: 'dashboard',
			component: () => import('../components/dashboard/Dashboard.vue'),
			meta: { requiresAuth: true },
		},
		{ path: '/privacy', name: 'privacy', component: () => import('../components/legal/PrivacyPage.vue') },
		{ path: '/terms', name: 'terms', component: () => import('../components/legal/TermsPage.vue') },
		{ path: '/dev', name: 'dev', component: () => import('../components/DevPrototype') },
		{ path: '/:pathMatch(.*)*', name: 'not-found', component: () => import('../components/NotFoundPage.vue') },
	],
	scrollBehavior(_to, _from, savedPosition) {
		return savedPosition ?? { top: 0, behavior: 'instant' };
	},
});

router.beforeEach((to) => {
	if (to.meta.requiresAuth && !hasValidSession()) return { name: 'landing' };
	if (to.meta.guestOnly && hasValidSession()) return { name: 'dashboard' };
});

export default router;
