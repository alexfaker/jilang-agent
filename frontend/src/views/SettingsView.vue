<template>
  <div class="container mx-auto px-4 py-6">
    <h1 class="text-2xl font-bold mb-6">系统设置</h1>
    
    <!-- 标签页导航 -->
    <div class="mb-6 border-b border-gray-200">
      <ul class="flex flex-wrap -mb-px text-sm font-medium text-center">
        <li class="mr-2" v-for="tab in tabs" :key="tab.id">
          <button 
            @click="activeTab = tab.id" 
            :class="[
              'inline-block p-4 rounded-t-lg',
              activeTab === tab.id 
                ? 'text-primary-600 border-b-2 border-primary-600 active' 
                : 'text-gray-500 hover:text-gray-600 hover:border-gray-300 border-b-2 border-transparent'
            ]"
          >
            <span class="flex items-center">
              <component :is="tab.icon" class="w-5 h-5 mr-2" />
              {{ tab.name }}
            </span>
          </button>
        </li>
      </ul>
    </div>
    
    <!-- 标签页内容 -->
    <div class="bg-white rounded-lg shadow p-6">
      <!-- 个人设置 -->
      <div v-if="activeTab === 'profile'" class="space-y-6">
        <h2 class="text-xl font-semibold mb-4">个人设置</h2>
        <ProfileSettings />
      </div>
      
      <!-- 系统设置 -->
      <div v-if="activeTab === 'system'" class="space-y-6">
        <h2 class="text-xl font-semibold mb-4">系统设置</h2>
        <SystemSettings />
      </div>
      
      <!-- 通知设置 -->
      <div v-if="activeTab === 'notifications'" class="space-y-6">
        <h2 class="text-xl font-semibold mb-4">通知设置</h2>
        <NotificationSettings />
      </div>
      
      <!-- 安全设置 -->
      <div v-if="activeTab === 'security'" class="space-y-6">
        <h2 class="text-xl font-semibold mb-4">安全设置</h2>
        <SecuritySettings />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { 
  UserIcon, 
  ComputerDesktopIcon, 
  BellIcon, 
  ShieldCheckIcon 
} from '@heroicons/vue/24/outline';

// 引入设置组件
import ProfileSettings from '../components/settings/ProfileSettings.vue';
import SystemSettings from '../components/settings/SystemSettings.vue';
import NotificationSettings from '../components/settings/NotificationSettings.vue';
import SecuritySettings from '../components/settings/SecuritySettings.vue';

// 标签页配置
const tabs = [
  { id: 'profile', name: '个人设置', icon: UserIcon },
  { id: 'system', name: '系统设置', icon: ComputerDesktopIcon },
  { id: 'notifications', name: '通知设置', icon: BellIcon },
  { id: 'security', name: '安全设置', icon: ShieldCheckIcon }
];

// 当前激活的标签页
const activeTab = ref('profile');
</script>

<style scoped>
/* 可以添加自定义样式 */
</style> 