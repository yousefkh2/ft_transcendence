import { createRouter, createWebHistory, type RouteLocationNormalized } from 'vue-router';
import LandingPage from '../components/LandingPage.vue';
import { hasValidSession } from '../auth/session';
import { t } from '../i18n';

declare module 'vue-router' {
	interface RouteMeta {
		requiresAuth?: boolean;
		guestOnly?: boolean;
		// Shown in the browser tab as "<title> | LingoDash". A function, so it follows the current language.
		title?: (to: RouteLocationNormalized) => string;
	}
}

const router = createRouter({
	history: createWebHistory(),
	routes: [
		{ path: '/', name: 'landing', component: LandingPage, meta: { guestOnly: true, title: () => t('titles.landing') } },
		{ path: '/landingpage', redirect: '/' },
		{
			path: '/dashboard',
			component: () => import('../components/dashboard/Dashboard.vue'),
			meta: { requiresAuth: true },
			children: [
				{ path: '', name: 'dashboard', component: () => import('../components/dashboard/OverviewTab.vue'), meta: { title: () => t('titles.overview') } },
				{ path: 'friends', name: 'friends', component: () => import('../components/dashboard/FriendsTab.vue'), meta: { title: () => t('titles.friends') } },
				{ path: 'profile', name: 'profile', component: () => import('../components/dashboard/ProfileTab.vue'), meta: { title: () => t('titles.profile') } },
				{
					path: 'profile/:username',
					name: 'user-profile',
					component: () => import('../components/dashboard/UserProfileTab.vue'),
					meta: { title: (to) => t('titles.userProfile', { username: String(to.params.username) }) },
				},
			],
		},
		{
			path: '/play',
			name: 'play',
			component: () => import('../components/play/PlayPage.vue'),
			meta: { requiresAuth: true, title: () => t('titles.play') },
		},
		{ path: '/privacy', name: 'privacy', component: () => import('../components/legal/PrivacyPage.vue'), meta: { title: () => t('titles.privacy') } },
		{ path: '/terms', name: 'terms', component: () => import('../components/legal/TermsPage.vue'), meta: { title: () => t('titles.terms') } },
		{ path: '/dev', name: 'dev', component: () => import('../components/DevPrototype') },
		{ path: '/:pathMatch(.*)*', name: 'not-found', component: () => import('../components/NotFoundPage.vue'), meta: { title: () => t('titles.notFound') } },
	],
	scrollBehavior(_to, _from, savedPosition) {
		return savedPosition ?? { top: 0, behavior: 'instant' };
	},
});

router.beforeEach((to) => {
	if (to.meta.requiresAuth && !hasValidSession()) return { name: 'landing' };
	if (to.meta.guestOnly && hasValidSession()) return { name: 'dashboard' };
});

export function updateTitle(to: RouteLocationNormalized) {
	const title = to.meta.title?.(to);
	document.title = title ? `${title} | LingoDash` : 'LingoDash';
}

router.afterEach(updateTitle);

export default router;
