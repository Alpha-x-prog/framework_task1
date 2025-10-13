import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('../pages/Login.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('../pages/Register.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/forbidden',
    name: 'Forbidden',
    component: () => import('../pages/Forbidden.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/',
    redirect: '/projects'
  },
  {
    path: '/projects',
    name: 'Projects',
    component: () => import('../pages/Dashboard.vue'),
    meta: { requiresAuth: true, roles: ['manager', 'engineer', 'lead', 'viewer'] }
  },
  {
    path: '/defects',
    name: 'Defects',
    component: () => import('../pages/Defects.vue'),
    meta: { requiresAuth: true, roles: ['manager', 'engineer', 'lead', 'viewer'] }
  },
  {
    path: '/defects/:id',
    name: 'DefectDetail',
    component: () => import('../pages/DefectDetail.vue'),
    meta: { requiresAuth: true, roles: ['manager', 'engineer', 'lead', 'viewer'] }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore()
  
  // Check if route requires authentication
  if (to.meta.requiresAuth !== false) {
    // Check if user is authenticated
    if (!authStore.isAuthed) {
      // Try to restore auth from localStorage
      const isAuthenticated = await authStore.checkAuth()
      if (!isAuthenticated) {
        next('/login')
        return
      }
    }
    
    // Check role-based access
    if (to.meta.roles && to.meta.roles.length > 0) {
      const userRole = authStore.role
      if (!to.meta.roles.includes(userRole)) {
        next('/forbidden')
        return
      }
    }
  } else {
    // If user is already authenticated and trying to access login/register
    // Only redirect if they're manually navigating to these pages
    if (authStore.isAuthed && (to.name === 'Login' || to.name === 'Register') && from.name !== 'Register' && from.name !== 'Login') {
      next('/projects')
      return
    }
  }
  
  next()
})

export default router
