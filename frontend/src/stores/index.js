// 导出所有状态存储
export * from './user';
export * from './settings';
export * from './workflow';
export * from './execution';
export * from './agent';
export * from './stats';
export * from './notification';
export * from './theme';
export * from './loading';

// 初始化函数
import { useThemeStore } from './theme';

/**
 * 初始化所有需要在应用启动时立即初始化的存储
 */
export function initializeStores() {
  // 初始化主题
  const themeStore = useThemeStore();
  themeStore.initTheme();
  
  // 未来可以在这里添加其他需要初始化的存储
} 