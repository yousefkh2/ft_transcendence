import { createRouter, createWebHistory, type RouteLocationNormalized } from 'vue-router';
import LandingPage from '../components/LandingPage.vue';
import { hasValidSession } from '../auth/session';

declare module 'vue-router' {
	interface RouteMeta {
		requiresAuth?: boolean;
		guestOnly?: boolean;
		// Shown in the browser tab as "<title> | LingoDash".
		title?: string | ((to: RouteLocationNormalized) => string);
	}
}

const router = createRouter({
	history: createWebHistory(),
	routes: [
		{ path: '/', name: 'landing', component: LandingPage, meta: { guestOnly: true, title: 'Learn languages through play' } },
		{ path: '/landingpage', redirect: '/' },
		{
			path: '/dashboard',
			component: () => import('../components/dashboard/Dashboard.vue'),
			meta: { requiresAuth: true },
			children: [
				{ path: '', name: 'dashboard', component: () => import('../components/dashboard/OverviewTab.vue'), meta: { title: 'Overview' } },
				{ path: 'friends', name: 'friends', component: () => import('../components/dashboard/FriendsTab.vue'), meta: { title: 'Friends' } },
				{ path: 'profile', name: 'profile', component: () => import('../components/dashboard/ProfileTab.vue'), meta: { title: 'Your profile' } },
				{
					path: 'profile/:username',
					name: 'user-profile',
					component: () => import('../components/dashboard/UserProfileTab.vue'),
					meta: { title: (to) => `${String(to.params.username)}'s profile` },
				},
			],
		},
		{
			path: '/play',
			name: 'play',
			component: () => import('../components/play/PlayPage.vue'),
			meta: { requiresAuth: true, title: 'Play' },
		},
		{ path: '/privacy', name: 'privacy', component: () => import('../components/legal/PrivacyPage.vue'), meta: { title: 'Privacy Policy' } },
		{ path: '/terms', name: 'terms', component: () => import('../components/legal/TermsPage.vue'), meta: { title: 'Terms of Service' } },
		{ path: '/dev', name: 'dev', component: () => import('../components/DevPrototype') },
		{ path: '/:pathMatch(.*)*', name: 'not-found', component: () => import('../components/NotFoundPage.vue'), meta: { title: 'Page not found' } },
	],
	scrollBehavior(_to, _from, savedPosition) {
		return savedPosition ?? { top: 0, behavior: 'instant' };
	},
});

router.beforeEach((to) => {
	if (to.meta.requiresAuth && !hasValidSession()) return { name: 'landing' };
	if (to.meta.guestOnly && hasValidSession()) return { name: 'dashboard' };
});

router.afterEach((to) => {
	const title = typeof to.meta.title === 'function' ? to.meta.title(to) : to.meta.title;
	document.title = title ? `${title} | LingoDash` : 'LingoDash';
});

export default router;
