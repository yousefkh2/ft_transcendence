import { createRouter, createWebHistory } from 'vue-router';
import LandingPage from '../components/LandingPage.vue';

const router = createRouter({
	history: createWebHistory(),
	routes: [
		{ path: '/', name: 'landing', component: LandingPage },
		{ path: '/landingpage', redirect: '/' },
		{ path: '/dashboard', name: 'dashboard', component: () => import('../components/dashboard/Dashboard.vue') },
		{ path: '/dev', name: 'dev', component: () => import('../components/DevPrototype') },
		{ path: '/:pathMatch(.*)*', name: 'not-found', component: () => import('../components/NotFoundPage.vue') },
	],
});

export default router;
