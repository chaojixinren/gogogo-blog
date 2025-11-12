import { createRouter, createWebHistory } from 'vue-router'
import HomePage from '@/pages/HomePage.vue'
import PostDetailPage from '@/pages/PostDetailPage.vue'
import LoginPage from '@/pages/LoginPage.vue'
import RegisterPage from '@/pages/RegisterPage.vue'
import DashboardLayout from '@/pages/dashboard/DashboardLayout.vue'
import DashboardPostsPage from '@/pages/dashboard/DashboardPostsPage.vue'
import DashboardCategoriesPage from '@/pages/dashboard/DashboardCategoriesPage.vue'
import DashboardTagsPage from '@/pages/dashboard/DashboardTagsPage.vue'
import { useAuthStore } from '@/store/auth'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomePage,
    },
    {
      path: '/posts/:slug',
      name: 'post-detail',
      component: PostDetailPage,
      props: true,
    },
    {
      path: '/login',
      name: 'login',
      component: LoginPage,
      meta: { guestOnly: true },
    },
    {
      path: '/register',
      name: 'register',
      component: RegisterPage,
      meta: { guestOnly: true },
    },
    {
      path: '/dashboard',
      component: DashboardLayout,
      meta: { requiresAuth: true },
      children: [
        {
          path: '',
          redirect: { name: 'dashboard-posts' },
        },
        {
          path: 'posts',
          name: 'dashboard-posts',
          component: DashboardPostsPage,
        },
        {
          path: 'categories',
          name: 'dashboard-categories',
          component: DashboardCategoriesPage,
        },
        {
          path: 'tags',
          name: 'dashboard-tags',
          component: DashboardTagsPage,
        },
      ],
    },
  ],
})

router.beforeEach(async (to) => {
  const auth = useAuthStore()
  await auth.initialize()

  if (to.meta.requiresAuth && !auth.isAuthenticated) {
    return {
      name: 'login',
      query: { redirect: to.fullPath },
    }
  }

  if (to.meta.guestOnly && auth.isAuthenticated) {
    return { name: 'home' }
  }

  return true
})

export default router
