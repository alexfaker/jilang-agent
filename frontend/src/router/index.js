import { createRouter, createWebHistory } from 'vue-router'
import Dashboard from '../views/DashboardView.vue'
import NotFound from '../views/NotFoundView.vue'
import Layout from '../components/layout/Layout.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      component: Layout,
      children: [
        {
          path: '',
          name: 'dashboard',
          component: Dashboard,
          meta: {
            title: '仪表盘',
            requiresAuth: true
          }
        },
        {
          path: 'workflows',
          name: 'workflows',
          component: () => import('../views/WorkflowView.vue'),
          meta: {
            title: '工作流管理',
            requiresAuth: true
          }
        },
        {
          path: 'executions',
          name: 'executions',
          component: () => import('../views/ExecutionView.vue'),
          meta: {
            title: '执行历史',
            requiresAuth: true
          }
        },
        {
          path: 'agents',
          name: 'agents',
          component: () => import('../views/AgentView.vue'),
          meta: {
            title: '代理管理',
            requiresAuth: true
          }
        },
        {
          path: 'settings',
          name: 'settings',
          component: () => import('../views/SettingsView.vue'),
          meta: {
            title: '系统设置',
            requiresAuth: true
          }
        },
        {
          path: 'profile',
          name: 'profile',
          component: () => import('../views/ProfileView.vue'),
          meta: {
            title: '个人资料',
            requiresAuth: true
          }
        }
      ]
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('../views/LoginView.vue'),
      meta: {
        title: '登录',
        requiresAuth: false
      }
    },
    {
      path: '/register',
      name: 'register',
      component: () => import('../views/RegisterView.vue'),
      meta: {
        title: '注册',
        requiresAuth: false
      }
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'NotFound',
      component: NotFound,
      meta: {
        title: '页面未找到',
        requiresAuth: false
      }
    }
  ]
})

// 路由前置守卫
router.beforeEach((to, from, next) => {
  // 设置页面标题
  document.title = to.meta.title ? `${to.meta.title} - JiLang Agent` : 'JiLang Agent';
  
  // 检查是否需要认证
  const isAuthenticated = localStorage.getItem('token');
  if (to.meta.requiresAuth && !isAuthenticated) {
    // 需要认证但未登录，重定向到登录页
    next({ name: 'login', query: { redirect: to.fullPath } });
  } else if ((to.name === 'login' || to.name === 'register') && isAuthenticated) {
    // 已登录用户访问登录/注册页，重定向到首页
    next({ name: 'dashboard' });
  } else {
    // 其他情况正常导航
    next();
  }
})

export default router 