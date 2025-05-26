import { reactive, watch } from 'vue';

// 默认设置
const defaultSettings = {
  // 系统设置
  system: {
    theme: 'system', // system, light, dark
    themeColor: 'blue',
    language: 'zh-CN',
    timezone: 'Asia/Shanghai',
    dateFormat: 'YYYY-MM-DD HH:mm:ss'
  },
  
  // 通知设置
  notification: {
    enabled: true,
    email: {
      workflowResult: true,
      workflowFailed: true,
      systemAnnouncements: true,
      digestFrequency: 'weekly'
    },
    browser: {
      enabled: false
    },
    desktop: {
      enabled: false
    },
    sound: true
  }
};

// 创建响应式状态
const state = reactive({
  // 从localStorage加载设置，如果没有则使用默认设置
  ...JSON.parse(JSON.stringify(defaultSettings)),
  // 是否已初始化
  initialized: false
});

// 加载设置
const loadSettings = () => {
  try {
    // 系统设置
    const systemSettings = localStorage.getItem('systemSettings');
    if (systemSettings) {
      Object.assign(state.system, JSON.parse(systemSettings));
    }
    
    // 通知设置
    const notificationSettings = localStorage.getItem('notificationSettings');
    if (notificationSettings) {
      Object.assign(state.notification, JSON.parse(notificationSettings));
    }
    
    // 应用主题
    applyTheme(state.system.theme);
    
    state.initialized = true;
  } catch (error) {
    console.error('加载设置失败:', error);
  }
};

// 保存设置
const saveSettings = () => {
  try {
    // 系统设置
    localStorage.setItem('systemSettings', JSON.stringify(state.system));
    
    // 通知设置
    localStorage.setItem('notificationSettings', JSON.stringify(state.notification));
  } catch (error) {
    console.error('保存设置失败:', error);
  }
};

// 重置设置
const resetSettings = () => {
  Object.assign(state.system, defaultSettings.system);
  Object.assign(state.notification, defaultSettings.notification);
  saveSettings();
  applyTheme(state.system.theme);
};

// 应用主题
const applyTheme = (theme) => {
  const html = document.documentElement;
  
  if (theme === 'system') {
    const isDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
    html.classList.toggle('dark', isDark);
  } else if (theme === 'dark') {
    html.classList.add('dark');
  } else {
    html.classList.remove('dark');
  }
};

// 监听系统主题变化
if (window.matchMedia) {
  const colorSchemeQuery = window.matchMedia('(prefers-color-scheme: dark)');
  colorSchemeQuery.addEventListener('change', (e) => {
    if (state.system.theme === 'system') {
      document.documentElement.classList.toggle('dark', e.matches);
    }
  });
}

// 监听设置变化，自动保存
watch(
  () => JSON.stringify(state),
  () => {
    if (state.initialized) {
      saveSettings();
    }
  },
  { deep: true }
);

// 导出设置状态和方法
export default {
  state,
  loadSettings,
  saveSettings,
  resetSettings,
  applyTheme
}; 