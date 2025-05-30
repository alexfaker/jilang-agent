<template>
  <div class="h-screen flex overflow-hidden bg-gray-50">
    <!-- 侧边栏 -->
    <aside class="w-20 md:w-64 bg-white shadow-md flex flex-col">
      <!-- 侧边栏头部 - Logo -->
      <div class="p-4 border-b flex items-center justify-center md:justify-start">
        <i class="fas fa-robot text-indigo-600 text-2xl mr-2"></i>
        <span class="hidden md:block text-xl font-bold text-gray-800">AI 工作流平台</span>
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
            <router-link 
              to="/dashboard" 
              :class="[
                'group relative flex items-center px-4 py-3 rounded-lg mx-2 transition-all duration-200',
                $route.path === '/dashboard' 
                  ? 'text-indigo-600 bg-indigo-50' 
                  : 'text-gray-700 hover:bg-indigo-50 hover:text-indigo-600'
              ]"
            >
              <i class="fas fa-tachometer-alt text-lg w-6 text-center"></i>
              <span class="hidden md:block ml-3">仪表盘</span>
              <!-- 工具提示 -->
              <div class="md:hidden group-hover:block absolute left-full ml-2 px-2 py-1 bg-gray-800 text-white text-xs rounded whitespace-nowrap z-50">查看平台概览</div>
            </router-link>
          </li>
          
          <!-- 工作流 -->
          <li>
            <router-link 
              to="/workflows" 
              :class="[
                'group relative flex items-center px-4 py-3 rounded-lg mx-2 transition-all duration-200',
                $route.path.startsWith('/workflows') 
                  ? 'text-indigo-600 bg-indigo-50' 
                  : 'text-gray-700 hover:bg-indigo-50 hover:text-indigo-600'
              ]"
            >
              <i class="fas fa-project-diagram text-lg w-6 text-center"></i>
              <span class="hidden md:block ml-3">我的工作流</span>
              <span class="hidden md:flex ml-auto bg-indigo-100 text-indigo-800 text-xs font-semibold px-2 py-1 rounded-full">{{ workflowCount }}</span>
              <!-- 工具提示 -->
              <div class="md:hidden group-hover:block absolute left-full ml-2 px-2 py-1 bg-gray-800 text-white text-xs rounded whitespace-nowrap z-50">管理您的工作流</div>
            </router-link>
          </li>
          
          <!-- 工作流市场/代理市场 -->
          <li>
            <router-link 
              to="/agents" 
              :class="[
                'group relative flex items-center px-4 py-3 rounded-lg mx-2 transition-all duration-200',
                $route.path.startsWith('/agents') 
                  ? 'text-indigo-600 bg-indigo-50' 
                  : 'text-gray-700 hover:bg-indigo-50 hover:text-indigo-600'
              ]"
            >
              <i class="fas fa-store text-lg w-6 text-center"></i>
              <span class="hidden md:block ml-3">工作流市场</span>
              <span class="hidden md:flex ml-auto bg-green-100 text-green-800 text-xs font-semibold px-2 py-1 rounded-full">新</span>
              <!-- 工具提示 -->
              <div class="md:hidden group-hover:block absolute left-full ml-2 px-2 py-1 bg-gray-800 text-white text-xs rounded whitespace-nowrap z-50">探索和使用工作流模板</div>
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
            <router-link 
              to="/executions" 
              :class="[
                'group relative flex items-center px-4 py-3 rounded-lg mx-2 transition-all duration-200',
                $route.path.startsWith('/executions') 
                  ? 'text-indigo-600 bg-indigo-50' 
                  : 'text-gray-700 hover:bg-indigo-50 hover:text-indigo-600'
              ]"
            >
              <i class="fas fa-history text-lg w-6 text-center"></i>
              <span class="hidden md:block ml-3">执行历史</span>
              <!-- 工具提示 -->
              <div class="md:hidden group-hover:block absolute left-full ml-2 px-2 py-1 bg-gray-800 text-white text-xs rounded whitespace-nowrap z-50">查看历史执行记录</div>
            </router-link>
          </li>
          
          <!-- 统计分析 -->
          <li>
            <router-link 
              to="/stats" 
              :class="[
                'group relative flex items-center px-4 py-3 rounded-lg mx-2 transition-all duration-200',
                $route.path.startsWith('/stats') 
                  ? 'text-indigo-600 bg-indigo-50' 
                  : 'text-gray-700 hover:bg-indigo-50 hover:text-indigo-600'
              ]"
            >
              <i class="fas fa-chart-line text-lg w-6 text-center"></i>
              <span class="hidden md:block ml-3">统计分析</span>
              <!-- 工具提示 -->
              <div class="md:hidden group-hover:block absolute left-full ml-2 px-2 py-1 bg-gray-800 text-white text-xs rounded whitespace-nowrap z-50">查看工作流执行统计和数据分析</div>
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
            <router-link 
              to="/settings" 
              :class="[
                'group relative flex items-center px-4 py-3 rounded-lg mx-2 transition-all duration-200',
                $route.path.startsWith('/settings') 
                  ? 'text-indigo-600 bg-indigo-50' 
                  : 'text-gray-700 hover:bg-indigo-50 hover:text-indigo-600'
              ]"
            >
              <i class="fas fa-cog text-lg w-6 text-center"></i>
              <span class="hidden md:block ml-3">设置</span>
              <!-- 工具提示 -->
              <div class="md:hidden group-hover:block absolute left-full ml-2 px-2 py-1 bg-gray-800 text-white text-xs rounded whitespace-nowrap z-50">系统设置和偏好</div>
            </router-link>
          </li>
          
          <!-- 帮助中心 -->
          <li>
            <router-link 
              to="/help" 
              :class="[
                'group relative flex items-center px-4 py-3 rounded-lg mx-2 transition-all duration-200',
                $route.path.startsWith('/help') 
                  ? 'text-indigo-600 bg-indigo-50' 
                  : 'text-gray-700 hover:bg-indigo-50 hover:text-indigo-600'
              ]"
            >
              <i class="fas fa-question-circle text-lg w-6 text-center"></i>
              <span class="hidden md:block ml-3">帮助中心</span>
              <!-- 工具提示 -->
              <div class="md:hidden group-hover:block absolute left-full ml-2 px-2 py-1 bg-gray-800 text-white text-xs rounded whitespace-nowrap z-50">帮助和支持</div>
            </router-link>
          </li>
        </ul>
      </nav>
      
      <!-- 侧边栏底部 - 用户信息 -->
      <div class="border-t p-4">
        <div class="user-menu-container relative">
          <div class="group flex items-center text-gray-700 hover:text-indigo-600 cursor-pointer" @click="toggleUserMenu">
            <div class="relative">
              <div class="w-8 h-8 rounded-full bg-indigo-600 flex items-center justify-center text-white text-sm font-medium">
                {{ userInitials }}
              </div>
              <span class="absolute -top-1 -right-1 w-3.5 h-3.5 bg-green-500 rounded-full border-2 border-white"></span>
            </div>
            <div class="hidden md:block ml-3 flex-1">
              <p class="text-sm font-medium">{{ userName }}</p>
              <p class="text-xs text-gray-500">查看个人资料</p>
            </div>
            <!-- 工具提示 -->
            <div class="md:hidden group-hover:block absolute left-full ml-2 px-2 py-1 bg-gray-800 text-white text-xs rounded whitespace-nowrap z-50">管理您的个人资料和设置</div>
          </div>
          
          <!-- 用户菜单下拉 -->
          <div v-if="isUserMenuOpen" class="absolute bottom-16 left-4 right-4 bg-white rounded-lg shadow-lg py-2 z-50 md:left-4 md:right-4">
            <router-link to="/profile" class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100" @click="handleProfileClick">
              个人资料
            </router-link>
            <button @click="handleLogout" class="block w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100">
              退出登录
            </button>
          </div>
        </div>
      </div>
      
      <!-- 主题开关 -->
      <div class="border-t p-4 flex items-center justify-center md:justify-between text-gray-700">
        <span class="hidden md:block text-xs text-gray-500">暗黑模式</span>
        <label class="relative inline-flex items-center cursor-pointer">
          <input 
            type="checkbox" 
            v-model="isDarkMode" 
            @change="toggleTheme"
            class="sr-only peer"
          >
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
              <input 
                type="text" 
                placeholder="搜索..." 
                v-model="searchQuery"
                class="pl-8 pr-4 py-2 rounded-lg border border-gray-300 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
              >
              <i class="fas fa-search absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400"></i>
            </div>
            
            <!-- 通知 -->
            <button 
              type="button" 
              class="relative p-2 text-gray-600 hover:text-indigo-600 focus:outline-none"
              @click="showNotifications"
            >
              <i class="fas fa-bell text-xl"></i>
              <span v-if="notificationCount > 0" class="absolute top-0 right-0 block h-4 w-4 rounded-full bg-red-500 text-white text-xs flex items-center justify-center">{{ notificationCount }}</span>
            </button>
            
            <!-- 帮助 -->
            <button 
              type="button" 
              class="p-2 text-gray-600 hover:text-indigo-600 focus:outline-none"
              @click="showHelp"
            >
              <i class="fas fa-question-circle text-xl"></i>
            </button>
          </div>
        </div>
      </header>

      <!-- 内容区域 -->
      <div class="flex-1 overflow-y-auto">
        <router-view />
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useUserStore } from '../../stores/user';
import notify from '../../utils/notification';

const router = useRouter();
const route = useRoute();
const userStore = useUserStore();

// 状态管理
const isUserMenuOpen = ref(false);
const isDarkMode = ref(false);
const searchQuery = ref('');
const notificationCount = ref(3);
const workflowCount = ref(12);

// 计算属性
const userName = computed(() => userStore.user?.username || '用户');
const userInitials = computed(() => {
  const name = userName.value;
  if (name.length <= 2) return name.toUpperCase();
  return name.charAt(0).toUpperCase();
});

const pageTitle = computed(() => {
  const titles = {
    '/dashboard': '仪表盘',
    '/workflows': '我的工作流',
    '/agents': '工作流市场',
    '/executions': '执行历史',
    '/stats': '统计分析',
    '/settings': '设置',
    '/help': '帮助中心',
    '/profile': '个人资料',
    '/purchase': '购买服务'
  };
  return titles[route.path] || '仪表盘';
});

// 方法
const toggleUserMenu = () => {
  console.log('toggleUserMenu 被调用，当前状态:', isUserMenuOpen.value)
  isUserMenuOpen.value = !isUserMenuOpen.value;
  console.log('切换后状态:', isUserMenuOpen.value)
};

const handleProfileClick = () => {
  console.log('点击个人资料，准备跳转到 /profile')
  isUserMenuOpen.value = false;
  router.push('/profile').then(() => {
    console.log('成功跳转到个人资料页面')
  }).catch((error) => {
    console.error('跳转失败:', error)
  })
};

const handleLogout = async () => {
  try {
    await userStore.logout();
    notify.success('已成功退出登录');
    router.push('/auth/login');
  } catch (error) {
    notify.error('退出登录失败');
  }
  isUserMenuOpen.value = false;
};

const toggleTheme = () => {
  // 主题切换逻辑
  document.documentElement.classList.toggle('dark', isDarkMode.value);
  localStorage.setItem('theme', isDarkMode.value ? 'dark' : 'light');
};

const showNotifications = () => {
  notify.info('通知功能开发中...');
};

const showHelp = () => {
  notify.info('帮助中心开发中...');
};

// 生命周期
onMounted(() => {
  // 初始化主题
  const savedTheme = localStorage.getItem('theme');
  if (savedTheme === 'dark') {
    isDarkMode.value = true;
    document.documentElement.classList.add('dark');
  }
  
  // 点击外部关闭用户菜单
  document.addEventListener('click', (e) => {
    if (!e.target.closest('.user-menu-container')) {
      isUserMenuOpen.value = false;
    }
  });
});
</script>

<style scoped>
.sidebar-icon {
  transition: all 0.2s ease;
}

.sidebar-icon:hover {
  transform: translateY(-2px);
}

.card {
  border-radius: 1rem;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
}

/* FontAwesome样式确保 */
.fas {
  font-family: "Font Awesome 5 Free";
  font-weight: 900;
}
</style> 