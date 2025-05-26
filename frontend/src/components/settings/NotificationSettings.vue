<template>
  <div class="space-y-6">
    <form @submit.prevent="saveNotificationSettings" class="space-y-6">
      <!-- 通知开关 -->
      <div class="bg-white rounded-lg border border-gray-200 p-6">
        <h3 class="text-lg font-medium mb-4">通知设置</h3>
        <div class="space-y-4">
          <div class="flex items-center justify-between">
            <div>
              <h4 class="text-base font-medium text-gray-900">启用通知</h4>
              <p class="text-sm text-gray-500">控制是否接收任何通知</p>
            </div>
            <label class="relative inline-flex items-center cursor-pointer">
              <input 
                type="checkbox" 
                v-model="notificationForm.enabled" 
                class="sr-only peer"
              >
              <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-primary-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-primary-600"></div>
            </label>
          </div>
        </div>
      </div>
      
      <!-- 邮件通知 -->
      <div class="bg-white rounded-lg border border-gray-200 p-6" :class="{ 'opacity-50': !notificationForm.enabled }">
        <h3 class="text-lg font-medium mb-4">邮件通知</h3>
        <div class="space-y-4">
          <div class="flex items-center justify-between">
            <div>
              <h4 class="text-base font-medium text-gray-900">工作流执行结果</h4>
              <p class="text-sm text-gray-500">当工作流执行完成时发送邮件通知</p>
            </div>
            <label class="relative inline-flex items-center cursor-pointer" :class="{ 'cursor-not-allowed': !notificationForm.enabled }">
              <input 
                type="checkbox" 
                v-model="notificationForm.email.workflowResult" 
                class="sr-only peer"
                :disabled="!notificationForm.enabled"
              >
              <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-primary-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-primary-600"></div>
            </label>
          </div>
          
          <div class="flex items-center justify-between">
            <div>
              <h4 class="text-base font-medium text-gray-900">工作流执行失败</h4>
              <p class="text-sm text-gray-500">仅当工作流执行失败时发送邮件通知</p>
            </div>
            <label class="relative inline-flex items-center cursor-pointer" :class="{ 'cursor-not-allowed': !notificationForm.enabled }">
              <input 
                type="checkbox" 
                v-model="notificationForm.email.workflowFailed" 
                class="sr-only peer"
                :disabled="!notificationForm.enabled"
              >
              <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-primary-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-primary-600"></div>
            </label>
          </div>
          
          <div class="flex items-center justify-between">
            <div>
              <h4 class="text-base font-medium text-gray-900">系统公告</h4>
              <p class="text-sm text-gray-500">接收系统更新和重要公告</p>
            </div>
            <label class="relative inline-flex items-center cursor-pointer" :class="{ 'cursor-not-allowed': !notificationForm.enabled }">
              <input 
                type="checkbox" 
                v-model="notificationForm.email.systemAnnouncements" 
                class="sr-only peer"
                :disabled="!notificationForm.enabled"
              >
              <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-primary-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-primary-600"></div>
            </label>
          </div>
          
          <div class="pt-2">
            <label for="emailDigest" class="block text-sm font-medium text-gray-700 mb-1">邮件摘要频率</label>
            <select
              id="emailDigest"
              v-model="notificationForm.email.digestFrequency"
              class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-primary-500 focus:border-primary-500"
              :disabled="!notificationForm.enabled"
            >
              <option value="never">从不发送摘要</option>
              <option value="daily">每日摘要</option>
              <option value="weekly">每周摘要</option>
              <option value="monthly">每月摘要</option>
            </select>
            <p class="text-xs text-gray-500 mt-1">定期接收工作流执行摘要</p>
          </div>
        </div>
      </div>
      
      <!-- 系统通知 -->
      <div class="bg-white rounded-lg border border-gray-200 p-6" :class="{ 'opacity-50': !notificationForm.enabled }">
        <h3 class="text-lg font-medium mb-4">系统通知</h3>
        <div class="space-y-4">
          <div class="flex items-center justify-between">
            <div>
              <h4 class="text-base font-medium text-gray-900">浏览器通知</h4>
              <p class="text-sm text-gray-500">在浏览器中显示通知</p>
            </div>
            <label class="relative inline-flex items-center cursor-pointer" :class="{ 'cursor-not-allowed': !notificationForm.enabled }">
              <input 
                type="checkbox" 
                v-model="notificationForm.browser.enabled" 
                class="sr-only peer"
                :disabled="!notificationForm.enabled"
                @change="requestNotificationPermission"
              >
              <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-primary-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-primary-600"></div>
            </label>
          </div>
          
          <div class="flex items-center justify-between">
            <div>
              <h4 class="text-base font-medium text-gray-900">桌面通知</h4>
              <p class="text-sm text-gray-500">在桌面显示通知</p>
            </div>
            <label class="relative inline-flex items-center cursor-pointer" :class="{ 'cursor-not-allowed': !notificationForm.enabled || !notificationForm.browser.enabled }">
              <input 
                type="checkbox" 
                v-model="notificationForm.desktop.enabled" 
                class="sr-only peer"
                :disabled="!notificationForm.enabled || !notificationForm.browser.enabled"
              >
              <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-primary-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-primary-600"></div>
            </label>
          </div>
          
          <div class="flex items-center justify-between">
            <div>
              <h4 class="text-base font-medium text-gray-900">声音提醒</h4>
              <p class="text-sm text-gray-500">收到通知时播放声音</p>
            </div>
            <label class="relative inline-flex items-center cursor-pointer" :class="{ 'cursor-not-allowed': !notificationForm.enabled }">
              <input 
                type="checkbox" 
                v-model="notificationForm.sound" 
                class="sr-only peer"
                :disabled="!notificationForm.enabled"
              >
              <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-primary-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-primary-600"></div>
            </label>
          </div>
        </div>
      </div>
      
      <!-- 提交按钮 -->
      <div class="flex justify-end">
        <button
          type="submit"
          class="px-4 py-2 bg-primary-600 text-white rounded-md hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2"
          :disabled="isSubmitting"
        >
          <span v-if="isSubmitting" class="flex items-center">
            <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            保存中...
          </span>
          <span v-else>保存设置</span>
        </button>
      </div>
    </form>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, watch } from 'vue';
import { useToast } from 'vue-toastification';
import settingsStore from '../../store/settings';

// Toast通知
const toast = useToast();

// 表单数据 - 使用状态管理中的数据
const notificationForm = reactive({...settingsStore.state.notification});

// 提交状态
const isSubmitting = ref(false);

// 监听浏览器通知状态变化
watch(() => notificationForm.browser.enabled, (newValue) => {
  if (newValue) {
    requestNotificationPermission();
  }
});

// 请求浏览器通知权限
const requestNotificationPermission = async () => {
  if (!notificationForm.browser.enabled) return;
  
  if (!('Notification' in window)) {
    notificationForm.browser.enabled = false;
    toast.error('你的浏览器不支持通知功能');
    return;
  }
  
  if (Notification.permission === 'granted') {
    return;
  }
  
  if (Notification.permission !== 'denied') {
    const permission = await Notification.requestPermission();
    
    if (permission !== 'granted') {
      notificationForm.browser.enabled = false;
      toast.error('通知权限被拒绝，请在浏览器设置中启用通知');
    }
  } else {
    notificationForm.browser.enabled = false;
    toast.error('通知权限被拒绝，请在浏览器设置中启用通知');
  }
};

// 检查浏览器通知权限
onMounted(() => {
  // 检查浏览器通知权限
  if ('Notification' in window) {
    if (Notification.permission !== 'granted') {
      notificationForm.browser.enabled = false;
      notificationForm.desktop.enabled = false;
    }
  } else {
    notificationForm.browser.enabled = false;
    notificationForm.desktop.enabled = false;
  }
});

// 保存通知设置
const saveNotificationSettings = async () => {
  isSubmitting.value = true;
  
  try {
    // 更新状态管理中的数据
    Object.assign(settingsStore.state.notification, notificationForm);
    
    // 如果启用了浏览器通知，确保有权限
    if (notificationForm.browser.enabled) {
      await requestNotificationPermission();
    }
    
    toast.success('通知设置已保存');
  } catch (error) {
    toast.error('保存设置失败: ' + error.message);
  } finally {
    isSubmitting.value = false;
  }
};

// 测试通知
const testNotification = () => {
  if (!notificationForm.browser.enabled) return;
  
  if ('Notification' in window && Notification.permission === 'granted') {
    new Notification('测试通知', {
      body: '这是一条测试通知，表示通知功能正常工作',
      icon: '/favicon.ico'
    });
  }
};
</script> 