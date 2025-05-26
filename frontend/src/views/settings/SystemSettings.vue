<template>
  <div>
    <h2 class="text-lg font-medium text-gray-900 dark:text-gray-100">系统设置</h2>
    <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
      自定义您的应用体验和偏好设置
    </p>

    <form @submit.prevent="saveSettings" class="mt-6 space-y-8">
      <!-- 主题设置 -->
      <div>
        <h3 class="text-md font-medium text-gray-900 dark:text-gray-100">主题设置</h3>
        <div class="mt-4 space-y-4">
          <div>
            <label for="theme-mode" class="block text-sm font-medium text-gray-700 dark:text-gray-300">
              显示模式
            </label>
            <select
              id="theme-mode"
              v-model="settings.theme.mode"
              class="mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white focus:outline-none focus:ring-primary-500 focus:border-primary-500 sm:text-sm rounded-md"
            >
              <option value="light">浅色模式</option>
              <option value="dark">深色模式</option>
              <option value="system">跟随系统</option>
            </select>
          </div>

          <div>
            <label for="theme-color" class="block text-sm font-medium text-gray-700 dark:text-gray-300">
              主题颜色
            </label>
            <div class="mt-1 grid grid-cols-5 gap-3">
              <div
                v-for="color in themeColors"
                :key="color.value"
                @click="settings.theme.color = color.value"
                class="relative flex items-center justify-center cursor-pointer rounded-md p-0.5"
                :class="{ 'ring-2 ring-offset-2 ring-primary-500': settings.theme.color === color.value }"
              >
                <span
                  class="h-8 w-8 rounded-full border border-gray-300 dark:border-gray-600"
                  :style="{ backgroundColor: color.hex }"
                  :title="color.name"
                ></span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 语言设置 -->
      <div>
        <h3 class="text-md font-medium text-gray-900 dark:text-gray-100">语言设置</h3>
        <div class="mt-4">
          <label for="language" class="block text-sm font-medium text-gray-700 dark:text-gray-300">
            界面语言
          </label>
          <select
            id="language"
            v-model="settings.language"
            class="mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white focus:outline-none focus:ring-primary-500 focus:border-primary-500 sm:text-sm rounded-md"
          >
            <option v-for="lang in languages" :key="lang.code" :value="lang.code">
              {{ lang.name }}
            </option>
          </select>
        </div>
      </div>

      <!-- 界面布局 -->
      <div>
        <h3 class="text-md font-medium text-gray-900 dark:text-gray-100">界面布局</h3>
        <div class="mt-4 space-y-4">
          <div>
            <label for="sidebar-position" class="block text-sm font-medium text-gray-700 dark:text-gray-300">
              侧边栏位置
            </label>
            <select
              id="sidebar-position"
              v-model="settings.layout.sidebarPosition"
              class="mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white focus:outline-none focus:ring-primary-500 focus:border-primary-500 sm:text-sm rounded-md"
            >
              <option value="left">左侧</option>
              <option value="right">右侧</option>
            </select>
          </div>

          <div>
            <label for="content-width" class="block text-sm font-medium text-gray-700 dark:text-gray-300">
              内容宽度
            </label>
            <select
              id="content-width"
              v-model="settings.layout.contentWidth"
              class="mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white focus:outline-none focus:ring-primary-500 focus:border-primary-500 sm:text-sm rounded-md"
            >
              <option value="full">全宽</option>
              <option value="contained">居中固定宽度</option>
            </select>
          </div>

          <div class="relative flex items-start">
            <div class="flex items-center h-5">
              <input
                id="compact-mode"
                v-model="settings.layout.compactMode"
                type="checkbox"
                class="focus:ring-primary-500 h-4 w-4 text-primary-600 border-gray-300 dark:border-gray-600 dark:bg-gray-700 rounded"
              />
            </div>
            <div class="ml-3 text-sm">
              <label for="compact-mode" class="font-medium text-gray-700 dark:text-gray-300">紧凑模式</label>
              <p class="text-gray-500 dark:text-gray-400">减小界面元素间距，显示更多内容</p>
            </div>
          </div>
        </div>
      </div>

      <!-- 时区设置 -->
      <div>
        <h3 class="text-md font-medium text-gray-900 dark:text-gray-100">时区设置</h3>
        <div class="mt-4">
          <label for="timezone" class="block text-sm font-medium text-gray-700 dark:text-gray-300">
            时区
          </label>
          <select
            id="timezone"
            v-model="settings.timezone"
            class="mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white focus:outline-none focus:ring-primary-500 focus:border-primary-500 sm:text-sm rounded-md"
          >
            <option v-for="tz in timezones" :key="tz.value" :value="tz.value">
              {{ tz.label }}
            </option>
          </select>
          <p class="mt-2 text-sm text-gray-500 dark:text-gray-400">
            当前时间: {{ currentTime }}
          </p>
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
import { ref, reactive, computed, onMounted, onUnmounted } from 'vue';
import { useToast } from 'vue-toastification';
import { useSettingsStore } from '../../stores/settings';
import { format } from 'date-fns';
import { zhCN } from 'date-fns/locale';

const toast = useToast();
const settingsStore = useSettingsStore();
const isSaving = ref(false);
const currentTime = ref('');
let timeInterval = null;

// 主题颜色选项
const themeColors = [
  { name: '蓝色', value: 'blue', hex: '#1976d2' },
  { name: '紫色', value: 'purple', hex: '#9c27b0' },
  { name: '绿色', value: 'green', hex: '#4caf50' },
  { name: '橙色', value: 'orange', hex: '#ff9800' },
  { name: '红色', value: 'red', hex: '#f44336' },
  { name: '青色', value: 'cyan', hex: '#00bcd4' },
  { name: '粉色', value: 'pink', hex: '#e91e63' },
  { name: '棕色', value: 'brown', hex: '#795548' },
  { name: '灰色', value: 'gray', hex: '#9e9e9e' },
  { name: '黑色', value: 'black', hex: '#212121' }
];

// 语言选项
const languages = [
  { code: 'zh-CN', name: '简体中文' },
  { code: 'en-US', name: 'English (US)' },
  { code: 'ja-JP', name: '日本語' },
  { code: 'ko-KR', name: '한국어' }
];

// 时区选项
const timezones = [
  { value: 'Asia/Shanghai', label: '(GMT+08:00) 北京, 上海, 香港, 台北' },
  { value: 'Asia/Tokyo', label: '(GMT+09:00) 东京, 大阪, 札幌' },
  { value: 'America/New_York', label: '(GMT-05:00) 纽约, 华盛顿特区' },
  { value: 'America/Los_Angeles', label: '(GMT-08:00) 洛杉矶, 旧金山' },
  { value: 'Europe/London', label: '(GMT+00:00) 伦敦, 爱丁堡' },
  { value: 'Europe/Paris', label: '(GMT+01:00) 巴黎, 柏林, 罗马, 马德里' },
  { value: 'Australia/Sydney', label: '(GMT+10:00) 悉尼, 墨尔本' }
];

// 默认设置
const defaultSettings = {
  theme: {
    mode: 'system',
    color: 'blue'
  },
  language: 'zh-CN',
  layout: {
    sidebarPosition: 'left',
    contentWidth: 'contained',
    compactMode: false
  },
  timezone: 'Asia/Shanghai'
};

// 初始设置
const initialSettings = {};

// 当前设置
const settings = reactive({
  theme: {
    mode: 'system',
    color: 'blue'
  },
  language: 'zh-CN',
  layout: {
    sidebarPosition: 'left',
    contentWidth: 'contained',
    compactMode: false
  },
  timezone: 'Asia/Shanghai'
});

// 检查设置是否已更改
const isSettingsChanged = computed(() => {
  return JSON.stringify(settings) !== JSON.stringify(initialSettings);
});

// 更新当前时间显示
const updateCurrentTime = () => {
  try {
    const now = new Date();
    currentTime.value = format(now, 'yyyy-MM-dd HH:mm:ss', { locale: zhCN });
  } catch (error) {
    console.error('更新时间显示失败:', error);
    currentTime.value = new Date().toLocaleString();
  }
};

// 加载设置
const loadSettings = () => {
  const userSettings = settingsStore.getSettings();
  
  // 合并默认设置和用户设置
  const mergedSettings = { ...defaultSettings };
  
  if (userSettings) {
    if (userSettings.theme) {
      mergedSettings.theme = {
        ...mergedSettings.theme,
        ...userSettings.theme
      };
    }
    
    if (userSettings.language) {
      mergedSettings.language = userSettings.language;
    }
    
    if (userSettings.layout) {
      mergedSettings.layout = {
        ...mergedSettings.layout,
        ...userSettings.layout
      };
    }
    
    if (userSettings.timezone) {
      mergedSettings.timezone = userSettings.timezone;
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
    
    // 更新系统设置
    const updatedSettings = {
      ...allSettings,
      theme: settings.theme,
      language: settings.language,
      layout: settings.layout,
      timezone: settings.timezone
    };
    
    // 保存更新后的设置
    await settingsStore.updateSettings(updatedSettings);
    
    // 应用主题设置
    settingsStore.setThemeMode(settings.theme.mode);
    settingsStore.setThemeColor(settings.theme.color);
    
    // 应用语言设置
    settingsStore.setLanguage(settings.language);
    
    // 更新初始设置以反映新的保存状态
    Object.assign(initialSettings, JSON.parse(JSON.stringify(settings)));
    
    toast.success('系统设置已保存');
  } catch (error) {
    toast.error('保存系统设置失败');
    console.error('保存系统设置失败:', error);
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
  updateCurrentTime();
  
  // 每秒更新一次当前时间
  timeInterval = setInterval(updateCurrentTime, 1000);
});

// 清理
onUnmounted(() => {
  if (timeInterval) {
    clearInterval(timeInterval);
  }
});
</script> 