<template>
  <div class="space-y-6">
    <form @submit.prevent="saveSystemSettings" class="space-y-6">
      <!-- 主题设置 -->
      <div class="bg-white rounded-lg border border-gray-200 p-6">
        <h3 class="text-lg font-medium mb-4">主题设置</h3>
        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">主题模式</label>
            <div class="grid grid-cols-3 gap-4">
              <div 
                v-for="theme in themes" 
                :key="theme.value" 
                class="relative"
              >
                <input 
                  type="radio" 
                  :id="theme.value" 
                  name="theme" 
                  :value="theme.value" 
                  v-model="systemForm.theme"
                  class="sr-only"
                />
                <label 
                  :for="theme.value" 
                  class="cursor-pointer block p-3 border-2 rounded-lg transition-all"
                  :class="systemForm.theme === theme.value ? 'border-primary-500 ring-2 ring-primary-200' : 'border-gray-200 hover:border-gray-300'"
                >
                  <div class="h-20 rounded flex items-center justify-center" :class="theme.previewClass">
                    <component :is="theme.icon" class="w-8 h-8" :class="theme.iconClass" />
                  </div>
                  <div class="mt-2 text-center text-sm font-medium">{{ theme.label }}</div>
                </label>
              </div>
            </div>
          </div>
          
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">主题色</label>
            <div class="grid grid-cols-6 gap-2">
              <div 
                v-for="color in themeColors" 
                :key="color.value" 
                class="relative"
              >
                <input 
                  type="radio" 
                  :id="color.value" 
                  name="themeColor" 
                  :value="color.value" 
                  v-model="systemForm.themeColor"
                  class="sr-only"
                />
                <label 
                  :for="color.value" 
                  class="cursor-pointer block p-1 border-2 rounded-lg transition-all"
                  :class="systemForm.themeColor === color.value ? 'border-primary-500 ring-2 ring-primary-200' : 'border-transparent hover:border-gray-300'"
                >
                  <div class="h-8 w-full rounded-md" :class="color.previewClass"></div>
                </label>
              </div>
            </div>
          </div>
        </div>
      </div>
      
      <!-- 语言设置 -->
      <div class="bg-white rounded-lg border border-gray-200 p-6">
        <h3 class="text-lg font-medium mb-4">语言设置</h3>
        <div class="space-y-4">
          <div>
            <label for="language" class="block text-sm font-medium text-gray-700 mb-1">界面语言</label>
            <select
              id="language"
              v-model="systemForm.language"
              class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-primary-500 focus:border-primary-500"
            >
              <option v-for="lang in languages" :key="lang.value" :value="lang.value">
                {{ lang.label }} ({{ lang.nativeName }})
              </option>
            </select>
          </div>
        </div>
      </div>
      
      <!-- 时区设置 -->
      <div class="bg-white rounded-lg border border-gray-200 p-6">
        <h3 class="text-lg font-medium mb-4">时区设置</h3>
        <div class="space-y-4">
          <div>
            <label for="timezone" class="block text-sm font-medium text-gray-700 mb-1">时区</label>
            <select
              id="timezone"
              v-model="systemForm.timezone"
              class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-primary-500 focus:border-primary-500"
            >
              <option v-for="tz in timezones" :key="tz.value" :value="tz.value">
                {{ tz.label }}
              </option>
            </select>
          </div>
          
          <div>
            <label for="dateFormat" class="block text-sm font-medium text-gray-700 mb-1">日期格式</label>
            <select
              id="dateFormat"
              v-model="systemForm.dateFormat"
              class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-primary-500 focus:border-primary-500"
            >
              <option v-for="format in dateFormats" :key="format.value" :value="format.value">
                {{ format.label }}
              </option>
            </select>
            <p class="text-xs text-gray-500 mt-1">示例: {{ formatDateExample(systemForm.dateFormat) }}</p>
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
import { ref, reactive, onMounted } from 'vue';
import { useToast } from 'vue-toastification';
import { 
  SunIcon, 
  MoonIcon, 
  ComputerDesktopIcon 
} from '@heroicons/vue/24/outline';
import settingsStore from '../../store/settings';

// Toast通知
const toast = useToast();

// 表单数据 - 使用状态管理中的数据
const systemForm = reactive({...settingsStore.state.system});

// 提交状态
const isSubmitting = ref(false);

// 主题选项
const themes = [
  { 
    label: '跟随系统', 
    value: 'system', 
    icon: ComputerDesktopIcon, 
    iconClass: 'text-gray-700',
    previewClass: 'bg-gradient-to-r from-gray-100 to-gray-300' 
  },
  { 
    label: '浅色模式', 
    value: 'light', 
    icon: SunIcon, 
    iconClass: 'text-amber-500',
    previewClass: 'bg-white border border-gray-200' 
  },
  { 
    label: '深色模式', 
    value: 'dark', 
    icon: MoonIcon, 
    iconClass: 'text-indigo-200',
    previewClass: 'bg-gray-800' 
  }
];

// 主题色选项
const themeColors = [
  { label: '蓝色', value: 'blue', previewClass: 'bg-blue-600' },
  { label: '绿色', value: 'green', previewClass: 'bg-green-600' },
  { label: '红色', value: 'red', previewClass: 'bg-red-600' },
  { label: '紫色', value: 'purple', previewClass: 'bg-purple-600' },
  { label: '橙色', value: 'orange', previewClass: 'bg-orange-600' },
  { label: '青色', value: 'cyan', previewClass: 'bg-cyan-600' }
];

// 语言选项
const languages = [
  { label: '简体中文', value: 'zh-CN', nativeName: '简体中文' },
  { label: 'English', value: 'en-US', nativeName: 'English' },
  { label: '日本語', value: 'ja-JP', nativeName: '日本語' },
  { label: '한국어', value: 'ko-KR', nativeName: '한국어' }
];

// 时区选项
const timezones = [
  { label: '(GMT+08:00) 北京, 香港, 上海', value: 'Asia/Shanghai' },
  { label: '(GMT+08:00) 台北', value: 'Asia/Taipei' },
  { label: '(GMT+09:00) 东京', value: 'Asia/Tokyo' },
  { label: '(GMT+09:00) 首尔', value: 'Asia/Seoul' },
  { label: '(GMT+00:00) 伦敦', value: 'Europe/London' },
  { label: '(GMT-05:00) 纽约', value: 'America/New_York' },
  { label: '(GMT-08:00) 洛杉矶', value: 'America/Los_Angeles' }
];

// 日期格式选项
const dateFormats = [
  { label: 'YYYY-MM-DD HH:mm:ss', value: 'YYYY-MM-DD HH:mm:ss' },
  { label: 'YYYY/MM/DD HH:mm:ss', value: 'YYYY/MM/DD HH:mm:ss' },
  { label: 'DD/MM/YYYY HH:mm:ss', value: 'DD/MM/YYYY HH:mm:ss' },
  { label: 'MM/DD/YYYY HH:mm:ss', value: 'MM/DD/YYYY HH:mm:ss' },
  { label: 'YYYY年MM月DD日 HH:mm:ss', value: 'YYYY年MM月DD日 HH:mm:ss' }
];

// 日期格式化示例
const formatDateExample = (format) => {
  const now = new Date();
  const year = now.getFullYear();
  const month = String(now.getMonth() + 1).padStart(2, '0');
  const day = String(now.getDate()).padStart(2, '0');
  const hours = String(now.getHours()).padStart(2, '0');
  const minutes = String(now.getMinutes()).padStart(2, '0');
  const seconds = String(now.getSeconds()).padStart(2, '0');
  
  return format
    .replace('YYYY', year)
    .replace('MM', month)
    .replace('DD', day)
    .replace('HH', hours)
    .replace('mm', minutes)
    .replace('ss', seconds);
};

// 获取系统设置
onMounted(() => {
  // 从localStorage加载设置
  const savedSettings = localStorage.getItem('systemSettings');
  if (savedSettings) {
    try {
      const settings = JSON.parse(savedSettings);
      Object.assign(systemForm, settings);
    } catch (error) {
      console.error('加载系统设置失败:', error);
    }
  }
  
  // 应用主题
  settingsStore.applyTheme(systemForm.theme);
});

// 保存系统设置
const saveSystemSettings = async () => {
  isSubmitting.value = true;
  
  try {
    // 更新状态管理中的数据
    Object.assign(settingsStore.state.system, systemForm);
    
    // 应用主题
    settingsStore.applyTheme(systemForm.theme);
    
    toast.success('系统设置已保存');
  } catch (error) {
    toast.error('保存设置失败: ' + error.message);
  } finally {
    isSubmitting.value = false;
  }
};
</script> 