import { createRouter, createWebHistory } from 'vue-router'
import LayoutView from '../layout/LayoutView.vue'
import LoginView from '../views/LoginView.vue'

const router = createRouter({
  history: createWebHistory((import.meta as any).env.BASE_URL),
  routes: [
    { path: '/login', name: 'login', component: LoginView },
    {
      path: '/',
      component: LayoutView,
      redirect: '/dashboard',
      children: [
        {
          path: 'dashboard',
          name: 'dashboard',
          component: () => import('../views/DashboardView.vue'),
          meta: { title: '首页' }
        },
        {
          path: 'paths',
          name: 'paths',
          component: () => import('../views/PathView.vue'),
          meta: { title: 'API路径' }
        },
        {
          path: 'menus',
          name: 'menus',
          component: () => import('../views/MenuView.vue'),
          meta: { title: '菜单管理' }
        },
        {
          path: 'roles',
          name: 'roles',
          component: () => import('../views/RoleView.vue'),
          meta: { title: '角色管理' }
        },
        {
          path: 'permissions',
          name: 'permissions',
          component: () => import('../views/PermissionView.vue'),
          meta: { title: '权限管理' }
        },
        {
          path: 'users',
          name: 'users',
          component: () => import('../views/UserView.vue'),
          meta: { title: '用户管理' }
        },
        {
          path: 'operation-logs',
          name: 'operation-logs',
          component: () => import('../views/OperationLogView.vue'),
          meta: { title: '操作历史' }
        },
        {
          path: 'migrations',
          name: 'migrations',
          component: () => import('../views/MigrationView.vue'),
          meta: { title: '迁移记录' }
        }
      ]
    }
  ]
})

import type { NavigationGuardNext, RouteLocationNormalized } from 'vue-router'

router.beforeEach((to: RouteLocationNormalized, from: RouteLocationNormalized, next: NavigationGuardNext) => {
  const token = localStorage.getItem('token')
  if (to.name !== 'login' && !token) {
    next({ name: 'login' })
  } else {
    next()
  }
})

export default router
