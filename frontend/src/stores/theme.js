import { defineStore } from 'pinia';

export const useThemeStore = defineStore('theme', {
  state: () => ({
    mode: localStorage.getItem('themeMode') || 'light', // light, dark, system
    color: localStorage.getItem('themeColor') || 'blue', // blue, green, purple, orange, red
    availableColors: [
      { name: 'blue', value: '#1976d2', label: '蓝色' },
      { name: 'green', value: '#2e7d32', label: '绿色' },
      { name: 'purple', value: '#7b1fa2', label: '紫色' },
      { name: 'orange', value: '#e65100', label: '橙色' },
      { name: 'red', value: '#c62828', label: '红色' }
    ],
    systemDarkMode: window.matchMedia('(prefers-color-scheme: dark)').matches,
    loading: false
  }),
  
  getters: {
    // 当前主题色对象
    currentColorObject: (state) => {
      return state.availableColors.find(color => color.name === state.color) || state.availableColors[0];
    },
    
    // 当前主题色值
    currentColorValue: (state) => {
      const colorObj = state.availableColors.find(color => color.name === state.color);
      return colorObj ? colorObj.value : '#1976d2';
    },
    
    // 当前是否为暗黑模式
    isDarkMode: (state) => {
      if (state.mode === 'system') {
        return state.systemDarkMode;
      }
      return state.mode === 'dark';
    },
    
    // 主题模式选项
    modeOptions: () => [
      { value: 'light', label: '浅色模式' },
      { value: 'dark', label: '深色模式' },
      { value: 'system', label: '跟随系统' }
    ]
  },
  
  actions: {
    // 初始化主题
    initTheme() {
      // 监听系统主题变化
      const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)');
      mediaQuery.addEventListener('change', e => {
        this.systemDarkMode = e.matches;
        if (this.mode === 'system') {
          this.applyTheme();
        }
      });
      
      // 应用主题
      this.applyTheme();
    },
    
    // 设置主题模式
    setMode(mode) {
      if (['light', 'dark', 'system'].includes(mode)) {
        this.mode = mode;
        localStorage.setItem('themeMode', mode);
        this.applyTheme();
      }
    },
    
    // 设置主题颜色
    setColor(color) {
      if (this.availableColors.some(c => c.name === color)) {
        this.color = color;
        localStorage.setItem('themeColor', color);
        this.applyTheme();
      }
    },
    
    // 应用主题
    applyTheme() {
      const isDark = this.isDarkMode;
      const colorValue = this.currentColorValue;
      
      // 设置文档根元素的类
      document.documentElement.classList.toggle('dark-theme', isDark);
      
      // 设置CSS变量
      document.documentElement.style.setProperty('--primary-color', colorValue);
      
      // 设置meta theme-color
      const metaThemeColor = document.querySelector('meta[name="theme-color"]');
      if (metaThemeColor) {
        metaThemeColor.setAttribute('content', colorValue);
      } else {
        const meta = document.createElement('meta');
        meta.name = 'theme-color';
        meta.content = colorValue;
        document.head.appendChild(meta);
      }
      
      // 为暗黑模式设置不同的背景和文本颜色
      if (isDark) {
        document.documentElement.style.setProperty('--bg-color', '#121212');
        document.documentElement.style.setProperty('--text-color', '#ffffff');
        document.documentElement.style.setProperty('--card-bg', '#1e1e1e');
        document.documentElement.style.setProperty('--border-color', '#333333');
      } else {
        document.documentElement.style.setProperty('--bg-color', '#f5f5f5');
        document.documentElement.style.setProperty('--text-color', '#333333');
        document.documentElement.style.setProperty('--card-bg', '#ffffff');
        document.documentElement.style.setProperty('--border-color', '#e0e0e0');
      }
    },
    
    // 切换暗黑/浅色模式
    toggleDarkMode() {
      const newMode = this.mode === 'dark' ? 'light' : 'dark';
      this.setMode(newMode);
    },
    
    // 获取对比色（用于在主题色背景上显示文字）
    getContrastColor(hexColor) {
      // 将十六进制颜色转换为RGB
      const r = parseInt(hexColor.slice(1, 3), 16);
      const g = parseInt(hexColor.slice(3, 5), 16);
      const b = parseInt(hexColor.slice(5, 7), 16);
      
      // 计算亮度
      const brightness = (r * 299 + g * 587 + b * 114) / 1000;
      
      // 如果亮度大于128，返回黑色，否则返回白色
      return brightness > 128 ? '#000000' : '#ffffff';
    }
  }
}); 