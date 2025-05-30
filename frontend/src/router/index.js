import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '../stores/user'
import SidebarLayout from '../components/layout/SidebarLayout.vue'

// 路由配置
const routes = [
  {
    path: '/',
    redirect: '/home'
  },
  {
    path: '/home',
    name: 'Home',
    component: () => import('../views/HomePage.vue'),
    meta: { guest: true }
  },
  {
    path: '/',
    component: SidebarLayout,
    meta: { requiresAuth: true },
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('../views/Dashboard.vue'),
        meta: { requiresAuth: true }
      },
      {
        path: 'settings',
        name: 'Settings',
        component: () => import('../views/settings/SettingsView.vue'),
        meta: { requiresAuth: true }
      },
      {
        path: 'workflows',
        name: 'Workflows',
        component: () => import('../views/workflows/WorkflowList.vue'),
        meta: { requiresAuth: true }
      },
      {
        path: 'workflows/create',
        name: 'CreateWorkflow',
        component: () => import('../views/workflows/WorkflowCreate.vue'),
        meta: { requiresAuth: true }
      },
      {
        path: 'workflows/:id',
        name: 'WorkflowDetail',
        component: () => import('../views/workflows/WorkflowDetail.vue'),
        meta: { requiresAuth: true }
      },
      {
        path: 'agents',
        name: 'Agents',
        component: () => import('../views/agents/AgentMarket.vue'),
        meta: { requiresAuth: true }
      },
      {
        path: 'executions',
        name: 'Executions',
        component: () => import('../views/executions/ExecutionList.vue'),
        meta: { requiresAuth: true }
      },
      {
        path: 'executions/:id',
        name: 'ExecutionDetail',
        component: () => import('../views/executions/ExecutionDetail.vue'),
        meta: { requiresAuth: true }
      },
      {
        path: 'stats',
        name: 'Stats',
        component: () => import('../views/stats/StatsView.vue'),
        meta: { requiresAuth: true }
      },
      {
        path: 'help',
        name: 'Help',
        component: () => import('../views/help/HelpCenter.vue'),
        meta: { requiresAuth: true }
      },
      {
        path: 'profile',
        name: 'Profile',
        component: () => import('../views/profile/ProfileView.vue'),
        meta: { requiresAuth: true }
      },
      {
        path: 'recharge',
        name: 'Recharge',
        component: () => import('../views/recharge/RechargeView.vue'),
        meta: { requiresAuth: true }
      },
      {
        path: 'purchase',
        name: 'Purchase',
        component: () => import('../views/purchase/PurchaseView.vue'),
        meta: { requiresAuth: true }
      }
    ]
  },
  {
    path: '/auth/login',
    name: 'Login',
    component: () => import('../views/auth/Login.vue'),
    meta: { guest: true }
  },
  {
    path: '/auth/register',
    name: 'Register',
    component: () => import('../views/auth/Register.vue'),
    meta: { guest: true }
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('../views/NotFound.vue')
  }
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes
})

// 导航守卫
router.beforeEach((to, from, next) => {
  const userStore = useUserStore()
  
  // 如果用户存储未初始化，从本地存储加载用户状态
  if (!userStore.isAuthenticated && localStorage.getItem('token')) {
    userStore.initializeFromLocalStorage()
  }

  // 检查是否需要登录
  if (to.matched.some(record => record.meta.requiresAuth) && !userStore.isAuthenticated) {
    next({ 
      name: 'Login', 
      query: { redirect: to.fullPath } 
    })
  } 
  // 检查是否为游客页面(如登录页)且用户已登录
  else if (to.matched.some(record => record.meta.guest) && userStore.isAuthenticated) {
    next({ name: 'Dashboard' })
  } 
  else {
    // 正常导航
    next()
  }
})

export default router 