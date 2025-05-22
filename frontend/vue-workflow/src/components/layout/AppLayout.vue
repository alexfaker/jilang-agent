<template>
  <div class="h-screen flex overflow-hidden">
    <!-- 侧边栏 -->
    <aside class="w-20 md:w-64 bg-white shadow-md flex flex-col">
      <!-- 侧边栏头部 - Logo -->
      <div class="p-4 border-b flex items-center justify-center md:justify-start">
        <i class="fas fa-robot text-brand text-2xl mr-2"></i>
        <span class="hidden md:block text-xl font-bold text-gray-800">AI Agent 平台</span>
      </div>
      
      <!-- 侧边栏导航 -->
      <nav class="flex-1 overflow-y-auto py-4">
        <!-- 分类标签：工作台 -->
        <div class="px-4 py-2">
          <span class="text-xs font-semibold text-gray-500 uppercase tracking-wider">工作台</span>
        </div>
        <ul class="mb-4">
          <!-- 仪表盘 -->
          <li>
            <router-link to="/" class="group relative flex items-center px-4 py-3 text-gray-700 hover:bg-indigo-50 hover:text-brand rounded-lg mx-2 transition-all duration-200" :class="{ 'sidebar-active': $route.name === 'Dashboard' }">
              <i class="fas fa-tachometer-alt text-lg w-6 text-center"></i>
              <span class="hidden md:block ml-3">仪表盘</span>
              <!-- 工具提示 -->
              <div class="hidden group-hover:block absolute left-full ml-2 px-2 py-1 bg-gray-800 text-white text-xs rounded whitespace-nowrap z-50">查看平台概览</div>
            </router-link>
          </li>
          
          <!-- 工作流 -->
          <li>
            <router-link to="/workflows" class="group relative flex items-center px-4 py-3 text-gray-700 hover:bg-indigo-50 hover:text-brand rounded-lg mx-2 transition-all duration-200" :class="{ 'sidebar-active': $route.name?.startsWith('Workflow') }">
              <i class="fas fa-project-diagram text-lg w-6 text-center"></i>
              <span class="hidden md:block ml-3">我的工作流</span>
              <span class="hidden md:flex ml-auto bg-indigo-100 text-indigo-800 text-xs font-semibold px-2 py-1 rounded-full">{{ workflowsCount }}</span>
              <!-- 工具提示 -->
              <div class="hidden group-hover:block absolute left-full ml-2 px-2 py-1 bg-gray-800 text-white text-xs rounded whitespace-nowrap z-50">管理您的AI工作流</div>
            </router-link>
          </li>
          
          <!-- 工作流市场 -->
          <li>
            <router-link to="/market" class="group relative flex items-center px-4 py-3 text-gray-700 hover:bg-indigo-50 hover:text-brand rounded-lg mx-2 transition-all duration-200" :class="{ 'sidebar-active': $route.name?.startsWith('Agent') }">
              <i class="fas fa-store text-lg w-6 text-center"></i>
              <span class="hidden md:block ml-3">工作流市场</span>
              <span class="hidden md:flex ml-auto bg-green-100 text-green-800 text-xs font-semibold px-2 py-1 rounded-full">新</span>
              <!-- 工具提示 -->
              <div class="hidden group-hover:block absolute left-full ml-2 px-2 py-1 bg-gray-800 text-white text-xs rounded whitespace-nowrap z-50">探索和使用AI代理</div>
            </router-link>
          </li>
        </ul>
        
        <!-- 分类标签：分析 -->
        <div class="px-4 py-2">
          <span class="text-xs font-semibold text-gray-500 uppercase tracking-wider">分析</span>
        </div>
        <ul class="mb-4">
          <!-- 执行历史 -->
          <li>
            <router-link to="/executions" class="group relative flex items-center px-4 py-3 text-gray-700 hover:bg-indigo-50 hover:text-brand rounded-lg mx-2 transition-all duration-200" :class="{ 'sidebar-active': $route.name === 'ExecutionHistory' }">
              <i class="fas fa-history text-lg w-6 text-center"></i>
              <span class="hidden md:block ml-3">执行历史</span>
              <!-- 工具提示 -->
              <div class="hidden group-hover:block absolute left-full ml-2 px-2 py-1 bg-gray-800 text-white text-xs rounded whitespace-nowrap z-50">查看历史执行记录</div>
            </router-link>
          </li>
          
          <!-- 统计分析 -->
          <li>
            <router-link to="/stats" class="group relative flex items-center px-4 py-3 text-gray-700 hover:bg-indigo-50 hover:text-brand rounded-lg mx-2 transition-all duration-200" :class="{ 'sidebar-active': $route.name === 'Statistics' }">
              <i class="fas fa-chart-bar text-lg w-6 text-center"></i>
              <span class="hidden md:block ml-3">统计分析</span>
              <!-- 工具提示 -->
              <div class="hidden group-hover:block absolute left-full ml-2 px-2 py-1 bg-gray-800 text-white text-xs rounded whitespace-nowrap z-50">查看使用数据和分析</div>
            </router-link>
          </li>
        </ul>
        
        <!-- 分类标签：系统 -->
        <div class="px-4 py-2">
          <span class="text-xs font-semibold text-gray-500 uppercase tracking-wider">系统</span>
        </div>
        <ul>
          <!-- 设置 -->
          <li>
            <router-link to="/settings" class="group relative flex items-center px-4 py-3 text-gray-700 hover:bg-indigo-50 hover:text-brand rounded-lg mx-2 transition-all duration-200" :class="{ 'sidebar-active': $route.name === 'Settings' }">
              <i class="fas fa-cog text-lg w-6 text-center"></i>
              <span class="hidden md:block ml-3">设置</span>
              <!-- 工具提示 -->
              <div class="hidden group-hover:block absolute left-full ml-2 px-2 py-1 bg-gray-800 text-white text-xs rounded whitespace-nowrap z-50">系统设置和偏好</div>
            </router-link>
          </li>
          
          <!-- 帮助中心 -->
          <li>
            <router-link to="/help" class="group relative flex items-center px-4 py-3 text-gray-700 hover:bg-indigo-50 hover:text-brand rounded-lg mx-2 transition-all duration-200" :class="{ 'sidebar-active': $route.name === 'Help' }">
              <i class="fas fa-question-circle text-lg w-6 text-center"></i>
              <span class="hidden md:block ml-3">帮助中心</span>
              <!-- 工具提示 -->
              <div class="hidden group-hover:block absolute left-full ml-2 px-2 py-1 bg-gray-800 text-white text-xs rounded whitespace-nowrap z-50">获取帮助和教程</div>
            </router-link>
          </li>
        </ul>
      </nav>
      
      <!-- 侧边栏底部 - 用户信息 -->
      <div class="border-t p-4">
        <router-link to="/profile" class="group relative flex items-center text-gray-700 hover:text-brand">
          <div class="relative">
            <img :src="userData.avatar || 'https://images.unsplash.com/photo-1535713875002-d1d0cf377fde?q=80&w=40&auto=format&fit=crop'" class="w-8 h-8 rounded-full border-2 border-white" alt="用户头像">
            <span class="absolute -top-1 -right-1 w-3.5 h-3.5 bg-green-500 rounded-full border-2 border-white"></span>
          </div>
          <div class="hidden md:block ml-3">
            <p class="text-sm font-medium">{{ userData.name || '用户' }}</p>
            <p class="text-xs text-gray-500">查看个人资料</p>
          </div>
          <!-- 工具提示 -->
          <div class="hidden group-hover:block absolute left-full ml-2 px-2 py-1 bg-gray-800 text-white text-xs rounded whitespace-nowrap z-50">管理您的个人资料和设置</div>
        </router-link>
      </div>
      
      <!-- 主题开关 -->
      <div class="border-t p-4 flex items-center justify-center md:justify-between text-gray-700">
        <span class="hidden md:block text-xs text-gray-500">暗黑模式</span>
        <label class="relative inline-flex items-center cursor-pointer">
          <input type="checkbox" v-model="darkMode" class="sr-only peer">
          <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-2 peer-focus:ring-indigo-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-indigo-600"></div>
        </label>
      </div>
    </aside>

    <!-- 主要内容区 -->
    <main class="flex-1 flex flex-col overflow-hidden">
      <!-- 顶部导航栏 -->
      <header class="bg-white shadow-sm z-10">
        <div class="flex items-center justify-between p-4">
          <!-- 面包屑导航 -->
          <div class="flex items-center">
            <h1 class="text-xl font-semibold text-gray-800">{{ pageTitle }}</h1>
          </div>
          
          <!-- 右侧操作区 -->
          <div class="flex items-center space-x-4">
            <!-- 搜索框 -->
            <div class="hidden md:block relative">
              <input type="text" placeholder="搜索..." class="pl-8 pr-4 py-2 rounded-lg border border-gray-300 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent">
              <i class="fas fa-search absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400"></i>
            </div>
            
            <!-- 通知图标 -->
            <button class="relative p-2 text-gray-500 hover:text-brand hover:bg-gray-100 rounded-full focus:outline-none">
              <i class="fas fa-bell"></i>
              <span v-if="unreadNotifications > 0" class="absolute top-0 right-0 bg-red-500 text-white text-xs w-4 h-4 flex items-center justify-center rounded-full">{{ unreadNotifications > 9 ? '9+' : unreadNotifications }}</span>
            </button>
            
            <slot name="header-actions"></slot>
          </div>
        </div>
      </header>

      <!-- 内容区域 -->
      <div class="flex-1 overflow-y-auto p-6 bg-gray-50">
        <slot></slot>
      </div>
    </main>
  </div>
</template>

<script lang="ts" setup>
import { ref, computed, watch } from 'vue'
import { useRoute } from 'vue-router'

// 用户数据（实际项目中应该从状态管理中获取）
const userData = ref({
  name: '张三',
  avatar: ''
})

// 工作流数量（实际项目中应该从API获取）
const workflowsCount = ref(12)

// 未读通知数量（实际项目中应该从API获取）
const unreadNotifications = ref(3)

// 暗黑模式切换
const darkMode = ref(false)

// 监听暗黑模式变化，添加/移除暗黑模式类
watch(darkMode, (newValue) => {
  if (newValue) {
    document.documentElement.classList.add('dark')
  } else {
    document.documentElement.classList.remove('dark')
  }
}, { immediate: true })

// 获取当前路由
const route = useRoute()

// 根据路由计算页面标题
const pageTitle = computed(() => {
  switch(route.name) {
    case 'Dashboard': return '仪表盘'
    case 'Workflows': return '我的工作流'
    case 'WorkflowCreate': return '创建工作流'
    case 'WorkflowDetail': return '工作流详情'
    case 'WorkflowExecution': return '执行工作流'
    case 'ExecutionHistory': return '执行历史'
    case 'AgentMarket': return '工作流市场'
    case 'AgentDetail': return '代理详情'
    case 'Statistics': return '统计分析'
    case 'Profile': return '个人资料'
    case 'Settings': return '系统设置'
    case 'Help': return '帮助中心'
    default: return '工作流平台'
  }
})
</script>

<style scoped>
/* 侧边栏和内容的自定义样式 */
.sidebar-active {
  @apply border-l-3 border-brand bg-indigo-50 text-brand;
}
</style>