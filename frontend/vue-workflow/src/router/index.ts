import { RouteRecordRaw } from 'vue-router'

// 导入页面组件
const Dashboard = () => import('../views/Dashboard.vue')
const Workflows = () => import('../views/Workflows.vue')
const WorkflowDetail = () => import('../views/WorkflowDetail.vue')
const WorkflowCreate = () => import('../views/WorkflowCreate.vue')
const WorkflowExecution = () => import('../views/WorkflowExecution.vue')
const ExecutionHistory = () => import('../views/ExecutionHistory.vue')
const AgentMarket = () => import('../views/AgentMarket.vue')
const AgentDetail = () => import('../views/AgentDetail.vue')
const Statistics = () => import('../views/Statistics.vue')
const Profile = () => import('../views/Profile.vue')
const Settings = () => import('../views/Settings.vue')
const Help = () => import('../views/Help.vue')
const Login = () => import('../views/Login.vue')
const Register = () => import('../views/Register.vue')
const NotFound = () => import('../views/NotFound.vue')

// 路由配置
const routes: Array<RouteRecordRaw> = [
  {
    path: '/',
    name: 'Dashboard',
    component: Dashboard,
    meta: { requiresAuth: true }
  },
  {
    path: '/workflows',
    name: 'Workflows',
    component: Workflows,
    meta: { requiresAuth: true }
  },
  {
    path: '/workflows/create',
    name: 'WorkflowCreate',
    component: WorkflowCreate,
    meta: { requiresAuth: true }
  },
  {
    path: '/workflows/:id',
    name: 'WorkflowDetail',
    component: WorkflowDetail,
    meta: { requiresAuth: true }
  },
  {
    path: '/workflows/:id/execute',
    name: 'WorkflowExecution',
    component: WorkflowExecution,
    meta: { requiresAuth: true }
  },
  {
    path: '/executions',
    name: 'ExecutionHistory',
    component: ExecutionHistory,
    meta: { requiresAuth: true }
  },
  {
    path: '/market',
    name: 'AgentMarket',
    component: AgentMarket,
    meta: { requiresAuth: true }
  },
  {
    path: '/market/:id',
    name: 'AgentDetail',
    component: AgentDetail,
    meta: { requiresAuth: true }
  },
  {
    path: '/stats',
    name: 'Statistics',
    component: Statistics,
    meta: { requiresAuth: true }
  },
  {
    path: '/profile',
    name: 'Profile',
    component: Profile,
    meta: { requiresAuth: true }
  },
  {
    path: '/settings',
    name: 'Settings',
    component: Settings,
    meta: { requiresAuth: true }
  },
  {
    path: '/help',
    name: 'Help',
    component: Help,
    meta: { requiresAuth: true }
  },
  {
    path: '/login',
    name: 'Login',
    component: Login,
    meta: { requiresAuth: false }
  },
  {
    path: '/register',
    name: 'Register',
    component: Register,
    meta: { requiresAuth: false }
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: NotFound
  }
]

export default routes 