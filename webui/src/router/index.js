import { createRouter, createWebHashHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import LoginView from '../views/LoginView.vue'
import ChatView from '../views/ChatView.vue'
import ProfileView from '../views/ProfileView.vue'

const router = createRouter({
  history: createWebHashHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView
    },
    {
      path: '/login',
      name: 'login',
      component: LoginView
    },
	{
      path: '/chat/:id', 
      name: 'chat',
      component: ChatView
    },
	{
      path: '/profile',
      name: 'profile',
      component: ProfileView
    }
  ]
})

// Se un utente non ha fatto il login (non ha il token), non può vedere la Home, viene rimbalzato al Login
router.beforeEach((to, from, next) => {
	const isAuthenticated = !!localStorage.getItem('token');

	if (to.name !== 'login' && !isAuthenticated) {
		// Se prova ad andare ovunque tranne che sul login e NON è autenticato, rimandalo al login
		next({ name: 'login' });
	} else if (to.name === 'login' && isAuthenticated) {
		// Se prova ad andare sul login ma è autenticato, rimandalo alla home
		next({ name: 'home' });
	} else {
		next();
	}
});

export default router