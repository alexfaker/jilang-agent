<template>
  <div>
    <h2 class="text-lg font-medium text-gray-900 dark:text-gray-100">通知设置</h2>
    <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
      管理您接收通知的方式和频率
    </p>

    <form @submit.prevent="saveSettings" class="mt-6 space-y-8">
      <!-- 电子邮件通知 -->
      <div>
        <h3 class="text-md font-medium text-gray-900 dark:text-gray-100">电子邮件通知</h3>
        <div class="mt-4 space-y-4">
          <div class="relative flex items-start">
            <div class="flex items-center h-5">
              <input
                id="email-enabled"
                v-model="settings.notifications.email.enabled"
                type="checkbox"
                class="focus:ring-primary-500 h-4 w-4 text-primary-600 border-gray-300 dark:border-gray-600 dark:bg-gray-700 rounded"
              />
            </div>
            <div class="ml-3 text-sm">
              <label for="email-enabled" class="font-medium text-gray-700 dark:text-gray-300">启用电子邮件通知</label>
              <p class="text-gray-500 dark:text-gray-400">接收重要通知和更新的电子邮件</p>
            </div>
          </div>

          <div class="ml-7 space-y-4" v-if="settings.notifications.email.enabled">
            <div v-for="(value, key) in settings.notifications.email.types" :key="key">
              <div class="relative flex items-start">
                <div class="flex items-center h-5">
                  <input
                    :id="`email-${key}`"
                    v-model="settings.notifications.email.types[key]"
                    type="checkbox"
                    class="focus:ring-primary-500 h-4 w-4 text-primary-600 border-gray-300 dark:border-gray-600 dark:bg-gray-700 rounded"
                  />
                </div>
                <div class="ml-3 text-sm">
                  <label :for="`email-${key}`" class="font-medium text-gray-700 dark:text-gray-300">
                    {{ getNotificationTypeLabel(key) }}
                  </label>
                  <p class="text-gray-500 dark:text-gray-400">{{ getNotificationTypeDescription(key) }}</p>
                </div>
              </div>
            </div>

            <div>
              <label for="email-frequency" class="block text-sm font-medium text-gray-700 dark:text-gray-300">发送频率</label>
              <select
                id="email-frequency"
                v-model="settings.notifications.email.frequency"
                class="mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white focus:outline-none focus:ring-primary-500 focus:border-primary-500 sm:text-sm rounded-md"
              >
                <option value="immediate">立即发送</option>
                <option value="daily">每日摘要</option>
                <option value="weekly">每周摘要</option>
              </select>
            </div>
          </div>
        </div>
      </div>

      <!-- 浏览器通知 -->
      <div>
        <h3 class="text-md font-medium text-gray-900 dark:text-gray-100">浏览器通知</h3>
        <div class="mt-4 space-y-4">
          <div class="relative flex items-start">
            <div class="flex items-center h-5">
              <input
                id="browser-enabled"
                v-model="settings.notifications.browser.enabled"
                type="checkbox"
                class="focus:ring-primary-500 h-4 w-4 text-primary-600 border-gray-300 dark:border-gray-600 dark:bg-gray-700 rounded"
                @change="requestBrowserPermissions"
              />
            </div>
            <div class="ml-3 text-sm">
              <label for="browser-enabled" class="font-medium text-gray-700 dark:text-gray-300">启用浏览器通知</label>
              <p class="text-gray-500 dark:text-gray-400">在浏览器中接收实时通知提醒</p>
            </div>
          </div>

          <div class="ml-7 space-y-4" v-if="settings.notifications.browser.enabled">
            <div v-for="(value, key) in settings.notifications.browser.types" :key="key">
              <div class="relative flex items-start">
                <div class="flex items-center h-5">
                  <input
                    :id="`browser-${key}`"
                    v-model="settings.notifications.browser.types[key]"
                    type="checkbox"
                    class="focus:ring-primary-500 h-4 w-4 text-primary-600 border-gray-300 dark:border-gray-600 dark:bg-gray-700 rounded"
                  />
                </div>
                <div class="ml-3 text-sm">
                  <label :for="`browser-${key}`" class="font-medium text-gray-700 dark:text-gray-300">
                    {{ getNotificationTypeLabel(key) }}
                  </label>
                  <p class="text-gray-500 dark:text-gray-400">{{ getNotificationTypeDescription(key) }}</p>
                </div>
              </div>
            </div>

            <div>
              <label for="browser-sound" class="block text-sm font-medium text-gray-700 dark:text-gray-300">通知声音</label>
              <select
                id="browser-sound"
                v-model="settings.notifications.browser.sound"
                class="mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white focus:outline-none focus:ring-primary-500 focus:border-primary-500 sm:text-sm rounded-md"
              >
                <option value="none">无声音</option>
                <option value="default">默认声音</option>
                <option value="bell">铃声</option>
                <option value="chime">风铃声</option>
              </select>
            </div>
          </div>
        </div>
      </div>

      <!-- 应用内通知 -->
      <div>
        <h3 class="text-md font-medium text-gray-900 dark:text-gray-100">应用内通知</h3>
        <div class="mt-4 space-y-4">
          <div class="relative flex items-start">
            <div class="flex items-center h-5">
              <input
                id="inapp-enabled"
                v-model="settings.notifications.inApp.enabled"
                type="checkbox"
                class="focus:ring-primary-500 h-4 w-4 text-primary-600 border-gray-300 dark:border-gray-600 dark:bg-gray-700 rounded"
              />
            </div>
            <div class="ml-3 text-sm">
              <label for="inapp-enabled" class="font-medium text-gray-700 dark:text-gray-300">启用应用内通知</label>
              <p class="text-gray-500 dark:text-gray-400">在应用中显示通知横幅和提醒</p>
            </div>
          </div>

          <div class="ml-7 space-y-4" v-if="settings.notifications.inApp.enabled">
            <div v-for="(value, key) in settings.notifications.inApp.types" :key="key">
              <div class="relative flex items-start">
                <div class="flex items-center h-5">
                  <input
                    :id="`inapp-${key}`"
                    v-model="settings.notifications.inApp.types[key]"
                    type="checkbox"
                    class="focus:ring-primary-500 h-4 w-4 text-primary-600 border-gray-300 dark:border-gray-600 dark:bg-gray-700 rounded"
                  />
                </div>
                <div class="ml-3 text-sm">
                  <label :for="`inapp-${key}`" class="font-medium text-gray-700 dark:text-gray-300">
                    {{ getNotificationTypeLabel(key) }}
                  </label>
                  <p class="text-gray-500 dark:text-gray-400">{{ getNotificationTypeDescription(key) }}</p>
                </div>
              </div>
            </div>

            <div>
              <label for="inapp-position" class="block text-sm font-medium text-gray-700 dark:text-gray-300">通知位置</label>
              <select
                id="inapp-position"
                v-model="settings.notifications.inApp.position"
                class="mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white focus:outline-none focus:ring-primary-500 focus:border-primary-500 sm:text-sm rounded-md"
              >
                <option value="top-right">右上角</option>
                <option value="top-left">左上角</option>
                <option value="bottom-right">右下角</option>
                <option value="bottom-left">左下角</option>
              </select>
            </div>
          </div>
        </div>
      </div>

      <!-- 提交按钮 -->
      <div class="pt-5 flex justify-end">
        <button
          type="button"
          @click="resetSettings"
          class="py-2 px-4 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm text-sm font-medium text-gray-700 dark:text-gray-200 bg-white dark:bg-gray-700 hover:bg-gray-50 dark:hover:bg-gray-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
        >
          重置为默认
        </button>
        <button
          type="submit"
          :disabled="isSaving || !isSettingsChanged"
          class="ml-3 inline-flex justify-center py-2 px-4 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          <span v-if="isSaving">保存中...</span>
          <span v-else>保存设置</span>
        </button>
      </div>
    </form>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue';
import { useToast } from 'vue-toastification';
import { useSettingsStore } from '../../stores/settings';

const toast = useToast();
const settingsStore = useSettingsStore();
const isSaving = ref(false);

// 默认设置
const defaultSettings = {
  notifications: {
    email: {
      enabled: true,
      frequency: 'immediate',
      types: {
        workflow_completed: true,
        workflow_failed: true,
        agent_error: true,
        new_updates: true,
        security_alerts: true
      }
    },
    browser: {
      enabled: false,
      sound: 'default',
      types: {
        workflow_completed: true,
        workflow_failed: true,
        agent_error: true,
        security_alerts: true
      }
    },
    inApp: {
      enabled: true,
      position: 'top-right',
      types: {
        workflow_completed: true,
        workflow_failed: true,
        agent_error: true,
        security_alerts: true,
        new_updates: true
      }
    }
  }
};

// 初始设置
const initialSettings = {};

// 当前设置
const settings = reactive({
  notifications: {
    email: {
      enabled: true,
      frequency: 'immediate',
      types: {
        workflow_completed: true,
        workflow_failed: true,
        agent_error: true,
        new_updates: true,
        security_alerts: true
      }
    },
    browser: {
      enabled: false,
      sound: 'default',
      types: {
        workflow_completed: true,
        workflow_failed: true,
        agent_error: true,
        security_alerts: true
      }
    },
    inApp: {
      enabled: true,
      position: 'top-right',
      types: {
        workflow_completed: true,
        workflow_failed: true,
        agent_error: true,
        security_alerts: true,
        new_updates: true
      }
    }
  }
});

// 通知类型标签和描述
const notificationTypes = {
  workflow_completed: {
    label: '工作流完成',
    description: '工作流执行成功完成时通知您'
  },
  workflow_failed: {
    label: '工作流失败',
    description: '工作流执行失败时通知您'
  },
  agent_error: {
    label: '代理错误',
    description: '代理发生错误时通知您'
  },
  security_alerts: {
    label: '安全警报',
    description: '检测到安全相关事件时通知您'
  },
  new_updates: {
    label: '新更新',
    description: '系统有新更新可用时通知您'
  }
};

// 获取通知类型标签
const getNotificationTypeLabel = (type) => {
  return notificationTypes[type]?.label || type;
};

// 获取通知类型描述
const getNotificationTypeDescription = (type) => {
  return notificationTypes[type]?.description || '';
};

// 检查设置是否已更改
const isSettingsChanged = computed(() => {
  return JSON.stringify(settings) !== JSON.stringify(initialSettings);
});

// 请求浏览器通知权限
const requestBrowserPermissions = async () => {
  if (!settings.notifications.browser.enabled) return;
  
  try {
    if (Notification && Notification.permission !== 'granted') {
      const permission = await Notification.requestPermission();
      
      if (permission !== 'granted') {
        toast.warning('浏览器通知权限被拒绝，您将无法接收浏览器通知');
        settings.notifications.browser.enabled = false;
      }
    }
  } catch (error) {
    console.error('请求通知权限失败:', error);
    toast.error('浏览器不支持通知功能');
    settings.notifications.browser.enabled = false;
  }
};

// 加载设置
const loadSettings = () => {
  const userSettings = settingsStore.getSettings();
  
  // 合并默认设置和用户设置
  const mergedSettings = { ...defaultSettings };
  
  if (userSettings && userSettings.notifications) {
    if (userSettings.notifications.email) {
      mergedSettings.notifications.email = {
        ...mergedSettings.notifications.email,
        ...userSettings.notifications.email
      };
      
      if (userSettings.notifications.email.types) {
        mergedSettings.notifications.email.types = {
          ...mergedSettings.notifications.email.types,
          ...userSettings.notifications.email.types
        };
      }
    }
    
    if (userSettings.notifications.browser) {
      mergedSettings.notifications.browser = {
        ...mergedSettings.notifications.browser,
        ...userSettings.notifications.browser
      };
      
      if (userSettings.notifications.browser.types) {
        mergedSettings.notifications.browser.types = {
          ...mergedSettings.notifications.browser.types,
          ...userSettings.notifications.browser.types
        };
      }
    }
    
    if (userSettings.notifications.inApp) {
      mergedSettings.notifications.inApp = {
        ...mergedSettings.notifications.inApp,
        ...userSettings.notifications.inApp
      };
      
      if (userSettings.notifications.inApp.types) {
        mergedSettings.notifications.inApp.types = {
          ...mergedSettings.notifications.inApp.types,
          ...userSettings.notifications.inApp.types
        };
      }
    }
  }
  
  // 更新当前设置
  Object.assign(settings, mergedSettings);
  
  // 更新初始设置以便比较变更
  Object.assign(initialSettings, JSON.parse(JSON.stringify(settings)));
};

// 保存设置
const saveSettings = async () => {
  isSaving.value = true;
  
  try {
    // 获取当前完整的用户设置
    const allSettings = settingsStore.getSettings();
    
    // 更新通知设置
    const updatedSettings = {
      ...allSettings,
      notifications: settings.notifications
    };
    
    // 保存更新后的设置
    await settingsStore.updateSettings(updatedSettings);
    
    // 更新初始设置以反映新的保存状态
    Object.assign(initialSettings, JSON.parse(JSON.stringify(settings)));
    
    toast.success('通知设置已保存');
  } catch (error) {
    toast.error('保存通知设置失败');
    console.error('保存通知设置失败:', error);
  } finally {
    isSaving.value = false;
  }
};

// 重置设置
const resetSettings = () => {
  // 重置为默认值
  Object.assign(settings, JSON.parse(JSON.stringify(defaultSettings)));
};

// 初始化
onMounted(() => {
  loadSettings();
});
</script> 