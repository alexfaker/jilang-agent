import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '../stores/user'
import Layout from '../components/layout/Layout.vue'

// 路由配置
const routes = [
  {
    path: '/',
    redirect: '/dashboard'
  },
  {
    path: '/',
    component: Layout,
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
        component: () => import('../views/agents/AgentList.vue'),
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
      }
    ]
  },
  {
    path: '/auth/login',
    name: 'Login',
    component: () => import('../views/auth/Login.vue'),
    meta: { guest: true }
  },
  // 暂时注释掉不存在的注册页面
  // {
  //   path: '/auth/register',
  //   name: 'Register',
  //   component: () => import('../views/auth/Register.vue'),
  //   meta: { guest: true }
  // },
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