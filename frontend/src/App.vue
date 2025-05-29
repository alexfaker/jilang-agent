<template>
  <div :class="{ 'dark': isDarkMode }" class="min-h-screen">
    <router-view />
  </div>
</template>

<script setup>
import { computed, onMounted, onUnmounted, watch } from 'vue';
import { useSettingsStore } from './stores/settings';
import { useToast } from 'vue-toastification';
import notify from './utils/notification';
import { useBrandStore } from './stores/brand';

const settingsStore = useSettingsStore();
const brandStore = useBrandStore();
const toast = useToast();

// 是否为深色模式
const isDarkMode = computed(() => settingsStore.isDarkMode);

// 监听系统颜色方案变化
const handleColorSchemeChange = (e) => {
  if (settingsStore.settings.theme.mode === 'system') {
    settingsStore.applyThemeSettings();
  }
};

// 监听主题设置变化
watch(
  () => settingsStore.settings.theme,
  () => {
    settingsStore.applyThemeSettings();
  },
  { deep: true }
);

onMounted(() => {
  // 监听系统颜色方案变化
  const colorSchemeMediaQuery = window.matchMedia('(prefers-color-scheme: dark)');
  
  // 添加事件监听器
  if (colorSchemeMediaQuery.addEventListener) {
    colorSchemeMediaQuery.addEventListener('change', handleColorSchemeChange);
  } else if (colorSchemeMediaQuery.addListener) {
    // 旧版浏览器兼容
    colorSchemeMediaQuery.addListener(handleColorSchemeChange);
  }
  
  // 应用主题设置
  settingsStore.applyThemeSettings();
  
  // 初始化通知，显示欢迎信息
  setTimeout(() => {
    notify.info(brandStore.welcomeMessage);
  }, 1000);
});

onUnmounted(() => {
  // 移除事件监听器
  const colorSchemeMediaQuery = window.matchMedia('(prefers-color-scheme: dark)');
  
  if (colorSchemeMediaQuery.removeEventListener) {
    colorSchemeMediaQuery.removeEventListener('change', handleColorSchemeChange);
  } else if (colorSchemeMediaQuery.removeListener) {
    // 旧版浏览器兼容
    colorSchemeMediaQuery.removeListener(handleColorSchemeChange);
  }
});
</script>

<style>
/* CSS变量 - 主题色 */
:root {
  /* 默认蓝色主题 */
  --color-primary-50: 239 246 255;  /* blue-50 */
  --color-primary-100: 219 234 254; /* blue-100 */
  --color-primary-200: 191 219 254; /* blue-200 */
  --color-primary-300: 147 197 253; /* blue-300 */
  --color-primary-400: 96 165 250;  /* blue-400 */
  --color-primary-500: 59 130 246;  /* blue-500 */
  --color-primary-600: 37 99 235;   /* blue-600 */
  --color-primary-700: 29 78 216;   /* blue-700 */
  --color-primary-800: 30 64 175;   /* blue-800 */
  --color-primary-900: 30 58 138;   /* blue-900 */
  --color-primary-950: 23 37 84;    /* blue-950 */
}

/* 基础样式 */
body {
  font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  @apply bg-gray-50 text-gray-900;
}

.dark body {
  @apply bg-gray-900 text-gray-100;
}

/* 主题颜色定义 */
[data-theme="blue"] {
  --color-primary-50: 239 246 255;  /* blue-50 */
  --color-primary-100: 219 234 254; /* blue-100 */
  --color-primary-200: 191 219 254; /* blue-200 */
  --color-primary-300: 147 197 253; /* blue-300 */
  --color-primary-400: 96 165 250;  /* blue-400 */
  --color-primary-500: 59 130 246;  /* blue-500 */
  --color-primary-600: 37 99 235;   /* blue-600 */
  --color-primary-700: 29 78 216;   /* blue-700 */
  --color-primary-800: 30 64 175;   /* blue-800 */
  --color-primary-900: 30 58 138;   /* blue-900 */
  --color-primary-950: 23 37 84;    /* blue-950 */
}

[data-theme="green"] {
  --color-primary-50: 240 253 244;  /* green-50 */
  --color-primary-100: 220 252 231; /* green-100 */
  --color-primary-200: 187 247 208; /* green-200 */
  --color-primary-300: 134 239 172; /* green-300 */
  --color-primary-400: 74 222 128;  /* green-400 */
  --color-primary-500: 34 197 94;   /* green-500 */
  --color-primary-600: 22 163 74;   /* green-600 */
  --color-primary-700: 21 128 61;   /* green-700 */
  --color-primary-800: 22 101 52;   /* green-800 */
  --color-primary-900: 20 83 45;    /* green-900 */
  --color-primary-950: 5 46 22;     /* green-950 */
}

[data-theme="red"] {
  --color-primary-50: 254 242 242;  /* red-50 */
  --color-primary-100: 254 226 226; /* red-100 */
  --color-primary-200: 254 202 202; /* red-200 */
  --color-primary-300: 252 165 165; /* red-300 */
  --color-primary-400: 248 113 113; /* red-400 */
  --color-primary-500: 239 68 68;   /* red-500 */
  --color-primary-600: 220 38 38;   /* red-600 */
  --color-primary-700: 185 28 28;   /* red-700 */
  --color-primary-800: 153 27 27;   /* red-800 */
  --color-primary-900: 127 29 29;   /* red-900 */
  --color-primary-950: 69 10 10;    /* red-950 */
}

[data-theme="purple"] {
  --color-primary-50: 250 245 255;  /* purple-50 */
  --color-primary-100: 243 232 255; /* purple-100 */
  --color-primary-200: 233 213 255; /* purple-200 */
  --color-primary-300: 216 180 254; /* purple-300 */
  --color-primary-400: 192 132 252; /* purple-400 */
  --color-primary-500: 168 85 247;  /* purple-500 */
  --color-primary-600: 147 51 234;  /* purple-600 */
  --color-primary-700: 126 34 206;  /* purple-700 */
  --color-primary-800: 107 33 168;  /* purple-800 */
  --color-primary-900: 88 28 135;   /* purple-900 */
  --color-primary-950: 59 7 100;    /* purple-950 */
}

[data-theme="orange"] {
  --color-primary-50: 255 247 237;  /* orange-50 */
  --color-primary-100: 255 237 213; /* orange-100 */
  --color-primary-200: 254 215 170; /* orange-200 */
  --color-primary-300: 253 186 116; /* orange-300 */
  --color-primary-400: 251 146 60;  /* orange-400 */
  --color-primary-500: 249 115 22;  /* orange-500 */
  --color-primary-600: 234 88 12;   /* orange-600 */
  --color-primary-700: 194 65 12;   /* orange-700 */
  --color-primary-800: 154 52 18;   /* orange-800 */
  --color-primary-900: 124 45 18;   /* orange-900 */
  --color-primary-950: 67 20 7;     /* orange-950 */
}

[data-theme="cyan"] {
  --color-primary-50: 236 254 255;  /* cyan-50 */
  --color-primary-100: 207 250 254; /* cyan-100 */
  --color-primary-200: 165 243 252; /* cyan-200 */
  --color-primary-300: 103 232 249; /* cyan-300 */
  --color-primary-400: 34 211 238;  /* cyan-400 */
  --color-primary-500: 6 182 212;   /* cyan-500 */
  --color-primary-600: 8 145 178;   /* cyan-600 */
  --color-primary-700: 14 116 144;  /* cyan-700 */
  --color-primary-800: 21 94 117;   /* cyan-800 */
  --color-primary-900: 22 78 99;    /* cyan-900 */
  --color-primary-950: 8 51 68;     /* cyan-950 */
}
</style> 