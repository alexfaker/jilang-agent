<template>
  <div class="bg-white dark:bg-gray-800 shadow rounded-lg">
    <div class="border-b border-gray-200 dark:border-gray-700">
      <nav class="flex flex-wrap -mb-px">
        <button
          v-for="tab in tabs"
          :key="tab.id"
          @click="activeTab = tab.id"
          :class="[
            activeTab === tab.id
              ? 'border-primary-500 text-primary-600 dark:text-primary-400'
              : 'border-transparent text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300 hover:border-gray-300 dark:hover:border-gray-600',
            'whitespace-nowrap py-4 px-6 border-b-2 font-medium text-sm'
          ]"
        >
          {{ tab.name }}
        </button>
      </nav>
    </div>

    <div class="p-6">
      <keep-alive>
        <component :is="currentTabComponent" />
      </keep-alive>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue';
import ProfileSettings from './ProfileSettings.vue';
import SystemSettings from './SystemSettings.vue';
import NotificationSettings from './NotificationSettings.vue';
import SecuritySettings from './SecuritySettings.vue';

const activeTab = ref('profile');

const tabs = [
  { id: 'profile', name: '个人资料', component: ProfileSettings },
  { id: 'system', name: '系统设置', component: SystemSettings },
  { id: 'notification', name: '通知设置', component: NotificationSettings },
  { id: 'security', name: '安全设置', component: SecuritySettings }
];

const currentTabComponent = computed(() => {
  const tab = tabs.find(t => t.id === activeTab.value);
  return tab ? tab.component : null;
});
</script> 