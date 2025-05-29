<template>
  <div class="min-h-screen bg-gray-100 dark:bg-gray-900">
    <!-- 顶部导航栏 -->
    <nav class="bg-white dark:bg-gray-800 shadow-sm">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-16">
          <!-- 左侧Logo和菜单 -->
          <div class="flex">
            <!-- Logo -->
            <div class="flex-shrink-0 flex items-center">
              <router-link to="/">
                <BrandName size="xl" weight="bold" color="primary-600 dark:text-primary-400" />
              </router-link>
            </div>
            
            <!-- 主导航菜单 -->
            <div class="hidden sm:ml-6 sm:flex sm:space-x-8">
              <router-link 
                v-for="item in mainNavItems" 
                :key="item.path" 
                :to="item.path"
                :class="[
                  $route.path === item.path || $route.path.startsWith(item.activePattern || item.path)
                    ? 'border-primary-500 text-gray-900 dark:text-white' 
                    : 'border-transparent text-gray-500 dark:text-gray-300 hover:border-gray-300 dark:hover:border-gray-600 hover:text-gray-700 dark:hover:text-gray-200',
                  'inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium'
                ]"
              >
                {{ item.name }}
              </router-link>
            </div>
          </div>
          
          <!-- 右侧用户菜单 -->
          <div class="hidden sm:ml-6 sm:flex sm:items-center">
            <!-- 通知按钮 -->
            <button 
              type="button" 
              class="p-1 rounded-full text-gray-400 hover:text-gray-500 dark:hover:text-gray-300 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
            >
              <span class="sr-only">查看通知</span>
              <BellIcon class="h-6 w-6" />
            </button>
            
            <!-- 用户下拉菜单 -->
            <div class="ml-3 relative">
              <div>
                <button 
                  @click="isUserMenuOpen = !isUserMenuOpen" 
                  type="button" 
                  class="bg-white dark:bg-gray-800 rounded-full flex text-sm focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500" 
                  id="user-menu-button" 
                  aria-expanded="false" 
                  aria-haspopup="true"
                >
                  <span class="sr-only">打开用户菜单</span>
                  <div class="h-8 w-8 rounded-full bg-primary-600 flex items-center justify-center text-white">
                    {{ userInitials }}
                  </div>
                </button>
              </div>
              
              <!-- 下拉菜单 -->
              <div 
                v-if="isUserMenuOpen" 
                class="origin-top-right absolute right-0 mt-2 w-48 rounded-md shadow-lg py-1 bg-white dark:bg-gray-800 ring-1 ring-black ring-opacity-5 focus:outline-none z-10" 
                role="menu" 
                aria-orientation="vertical" 
                aria-labelledby="user-menu-button" 
                tabindex="-1"
              >
                <router-link 
                  v-for="item in userMenuItems" 
                  :key="item.path" 
                  :to="item.path" 
                  class="block px-4 py-2 text-sm text-gray-700 dark:text-gray-200 hover:bg-gray-100 dark:hover:bg-gray-700" 
                  role="menuitem" 
                  tabindex="-1" 
                  @click="isUserMenuOpen = false"
                >
                  {{ item.name }}
                </router-link>
                <button 
                  @click="handleLogout" 
                  class="block w-full text-left px-4 py-2 text-sm text-gray-700 dark:text-gray-200 hover:bg-gray-100 dark:hover:bg-gray-700" 
                  role="menuitem" 
                  tabindex="-1"
                >
                  退出登录
                </button>
              </div>
            </div>
          </div>
          
          <!-- 移动端菜单按钮 -->
          <div class="-mr-2 flex items-center sm:hidden">
            <button 
              @click="isMobileMenuOpen = !isMobileMenuOpen" 
              type="button" 
              class="inline-flex items-center justify-center p-2 rounded-md text-gray-400 hover:text-gray-500 hover:bg-gray-100 dark:hover:bg-gray-700 dark:hover:text-gray-300 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-primary-500" 
              aria-controls="mobile-menu" 
              aria-expanded="false"
            >
              <span class="sr-only">打开主菜单</span>
              <Bars3Icon v-if="!isMobileMenuOpen" class="block h-6 w-6" />
              <XMarkIcon v-else class="block h-6 w-6" />
            </button>
          </div>
        </div>
      </div>
      
      <!-- 移动端菜单 -->
      <div v-if="isMobileMenuOpen" class="sm:hidden" id="mobile-menu">
        <div class="pt-2 pb-3 space-y-1">
          <router-link 
            v-for="item in mainNavItems" 
            :key="item.path" 
            :to="item.path"
            :class="[
              $route.path === item.path || $route.path.startsWith(item.activePattern || item.path)
                ? 'bg-primary-50 dark:bg-primary-900 border-primary-500 text-primary-700 dark:text-primary-300' 
                : 'border-transparent text-gray-600 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-800 hover:border-gray-300 dark:hover:border-gray-600 hover:text-gray-800 dark:hover:text-white',
              'block pl-3 pr-4 py-2 border-l-4 text-base font-medium'
            ]"
            @click="isMobileMenuOpen = false"
          >
            {{ item.name }}
          </router-link>
        </div>
        <div class="pt-4 pb-3 border-t border-gray-200 dark:border-gray-700">
          <div class="flex items-center px-4">
            <div class="flex-shrink-0">
              <div class="h-10 w-10 rounded-full bg-primary-600 flex items-center justify-center text-white">
                {{ userInitials }}
              </div>
            </div>
            <div class="ml-3">
              <div class="text-base font-medium text-gray-800 dark:text-white">{{ userName }}</div>
              <div class="text-sm font-medium text-gray-500 dark:text-gray-400">{{ userEmail }}</div>
            </div>
            <button 
              type="button" 
              class="ml-auto flex-shrink-0 p-1 rounded-full text-gray-400 hover:text-gray-500 dark:hover:text-gray-300 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
            >
              <span class="sr-only">查看通知</span>
              <BellIcon class="h-6 w-6" />
            </button>
          </div>
          <div class="mt-3 space-y-1">
            <router-link 
              v-for="item in userMenuItems" 
              :key="item.path" 
              :to="item.path" 
              class="block px-4 py-2 text-base font-medium text-gray-500 dark:text-gray-400 hover:text-gray-800 dark:hover:text-white hover:bg-gray-100 dark:hover:bg-gray-700" 
              @click="isMobileMenuOpen = false"
            >
              {{ item.name }}
            </router-link>
            <button 
              @click="handleLogout" 
              class="block w-full text-left px-4 py-2 text-base font-medium text-gray-500 dark:text-gray-400 hover:text-gray-800 dark:hover:text-white hover:bg-gray-100 dark:hover:bg-gray-700"
            >
              退出登录
            </button>
          </div>
        </div>
      </div>
    </nav>
    
    <!-- 页面内容 -->
    <main class="py-6">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <router-view />
      </div>
    </main>
    
    <!-- 页脚 -->
    <footer class="bg-white dark:bg-gray-800 shadow-inner mt-auto">
      <div class="max-w-7xl mx-auto py-4 px-4 sm:px-6 lg:px-8">
        <p class="text-center text-sm text-gray-500 dark:text-gray-400">
          &copy; {{ new Date().getFullYear() }} <BrandName size="sm" weight="normal" color="gray-500 dark:text-gray-400" />. All rights reserved.
        </p>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount, watch } from 'vue';
import { useRouter } from 'vue-router';
import { BellIcon, Bars3Icon, XMarkIcon } from '@heroicons/vue/24/outline';
import { useUserStore } from '../../stores/user';
import notify from '../../utils/notification';
import BrandName from '../common/BrandName.vue';

// 路由
const router = useRouter();
const userStore = useUserStore();

// 菜单状态
const isUserMenuOpen = ref(false);
const isMobileMenuOpen = ref(false);

// 用户信息
const userName = computed(() => {
  return userStore.user?.nickname || userStore.user?.username || '用户';
});

const userEmail = computed(() => {
  return userStore.user?.email || 'user@example.com';
});

const userInitials = computed(() => {
  return userName.value.substring(0, 1).toUpperCase();
});

// 主导航菜单项
const mainNavItems = [
  { name: '仪表盘', path: '/dashboard', activePattern: '/dashboard' },
  { name: '工作流管理', path: '/workflows', activePattern: '/workflows' },
  { name: '执行历史', path: '/executions', activePattern: '/executions' },
  { name: '代理管理', path: '/agents', activePattern: '/agents' },
  { name: '积分管理', path: '/points', activePattern: '/points' },
  { name: '充值中心', path: '/recharge', activePattern: '/recharge' },
  { name: '服务购买', path: '/purchase', activePattern: '/purchase' }
];

// 用户菜单项
const userMenuItems = [
  { name: '个人资料', path: '/settings' },
  { name: '系统设置', path: '/settings' }
];

// 点击外部关闭用户菜单
const handleClickOutside = (event) => {
  const userMenu = document.getElementById('user-menu-button');
  if (userMenu && !userMenu.contains(event.target) && isUserMenuOpen.value) {
    isUserMenuOpen.value = false;
  }
};

// 退出登录
const handleLogout = async () => {
  try {
    await userStore.logout();
    router.push('/auth/login');
    notify.success('已成功退出登录');
  } catch (error) {
    notify.error('退出登录失败: ' + (error.message || '未知错误'));
  }
};

// 监听用户状态变化
watch(
  () => userStore.isAuthenticated,
  (isAuthenticated) => {
    if (!isAuthenticated && router.currentRoute.value.meta.requiresAuth) {
      router.push('/auth/login');
    }
  }
);

// 生命周期钩子
onMounted(() => {
  document.addEventListener('click', handleClickOutside);
  
  // 如果用户已登录但没有用户信息，获取用户资料
  if (userStore.isAuthenticated && !userStore.user) {
    userStore.fetchUserProfile();
  }
});

onBeforeUnmount(() => {
  document.removeEventListener('click', handleClickOutside);
});
</script> 