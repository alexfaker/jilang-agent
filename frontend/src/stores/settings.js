import { defineStore } from 'pinia';
import { settingsApi } from '../api';

/**
 * 用户设置存储
 */
export const useSettingsStore = defineStore('settings', {
  state: () => ({
    settings: {
      theme: {
        mode: 'light', // light, dark, system
        color: 'blue'  // primary color theme
      },
      language: 'zh-CN',
      dateFormat: 'yyyy-MM-dd',
      timeFormat: 'HH:mm:ss',
      timezone: 'Asia/Shanghai',
      layout: {
        sidebarPosition: 'left', // left, right
        contentWidth: 'contained', // contained, full
        compactMode: false
      },
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
    },
    loading: false,
    error: null
  }),
  
  getters: {
    isDarkMode: (state) => {
      if (state.settings.theme.mode === 'system') {
        return window.matchMedia('(prefers-color-scheme: dark)').matches;
      }
      return state.settings.theme.mode === 'dark';
    },
    
    currentLanguage: (state) => state.settings.language,
    
    primaryColor: (state) => {
      const colorMap = {
        blue: '#1976d2',
        green: '#4caf50',
        red: '#f44336',
        purple: '#9c27b0',
        orange: '#ff9800',
        cyan: '#00bcd4',
        pink: '#e91e63',
        brown: '#795548',
        gray: '#9e9e9e',
        black: '#212121'
      };
      
      return colorMap[state.settings.theme.color] || colorMap.blue;
    },
    
    getSettings: (state) => state.settings
  },
  
  actions: {
    /**
     * 初始化设置
     * 首先尝试从API获取设置，如果失败则从本地存储加载
     */
    async initSettings() {
      try {
        await this.fetchSettings();
      } catch (error) {
        console.error('初始化设置失败:', error);
        // 如果API调用失败，尝试从本地存储加载
        this.loadFromLocalStorage();
      }
      
      // 无论如何，确保主题被应用
      this.applyThemeSettings();
    },
    
    async fetchSettings() {
      this.loading = true;
      
      try {
        // 尝试从API获取设置
        const settings = await settingsApi.getSettings();
        this.updateSettings(settings);
      } catch (error) {
        console.error('获取设置失败:', error);
        this.error = error.message || '获取设置失败';
        
        // 如果API调用失败，尝试从本地存储加载
        this.loadFromLocalStorage();
      } finally {
        this.loading = false;
      }
    },
    
    async updateSettings(settings) {
      this.loading = true;
      
      try {
        // 更新本地设置
        if (settings.theme) {
          this.settings.theme = {
            ...this.settings.theme,
            ...settings.theme
          };
        }
        
        if (settings.language) {
          this.settings.language = settings.language;
        }
        
        if (settings.dateFormat) {
          this.settings.dateFormat = settings.dateFormat;
        }
        
        if (settings.timeFormat) {
          this.settings.timeFormat = settings.timeFormat;
        }
        
        if (settings.timezone) {
          this.settings.timezone = settings.timezone;
        }
        
        if (settings.layout) {
          this.settings.layout = {
            ...this.settings.layout,
            ...settings.layout
          };
        }
        
        if (settings.notifications) {
          this.settings.notifications = {
            ...this.settings.notifications,
            ...settings.notifications
          };
        }
        
        // 保存到API
        await settingsApi.updateSettings(this.settings);
        
        // 保存到本地存储
        this.saveToLocalStorage();
        
        // 应用设置
        this.applyThemeSettings();
        this.applyLayoutSettings();
        
        return this.settings;
      } catch (error) {
        console.error('更新设置失败:', error);
        this.error = error.message || '更新设置失败';
        throw error;
      } finally {
        this.loading = false;
      }
    },
    
    setThemeMode(mode) {
      this.settings.theme.mode = mode;
      this.applyThemeSettings();
      this.saveToLocalStorage();
    },
    
    setThemeColor(color) {
      this.settings.theme.color = color;
      this.applyThemeSettings();
      this.saveToLocalStorage();
    },
    
    setLanguage(language) {
      this.settings.language = language;
      this.saveToLocalStorage();
    },
    
    setLayoutOption(option, value) {
      if (this.settings.layout && option in this.settings.layout) {
        this.settings.layout[option] = value;
        this.applyLayoutSettings();
        this.saveToLocalStorage();
      }
    },
    
    /**
     * 应用主题设置
     * 包括暗黑模式和主题色
     */
    applyThemeSettings() {
      const isDark = this.isDarkMode;
      
      // 应用暗黑模式
      if (isDark) {
        document.documentElement.classList.add('dark');
      } else {
        document.documentElement.classList.remove('dark');
      }
      
      // 应用主题色
      document.documentElement.setAttribute('data-theme', this.settings.theme.color);
      
      // 设置CSS变量以便应用主题颜色
      document.documentElement.style.setProperty('--primary-color', this.primaryColor);
    },
    
    /**
     * 应用布局设置
     */
    applyLayoutSettings() {
      // 侧边栏位置
      if (this.settings.layout.sidebarPosition === 'right') {
        document.documentElement.classList.add('sidebar-right');
      } else {
        document.documentElement.classList.remove('sidebar-right');
      }
      
      // 内容宽度
      if (this.settings.layout.contentWidth === 'full') {
        document.documentElement.classList.add('content-full-width');
      } else {
        document.documentElement.classList.remove('content-full-width');
      }
      
      // 紧凑模式
      if (this.settings.layout.compactMode) {
        document.documentElement.classList.add('compact-mode');
      } else {
        document.documentElement.classList.remove('compact-mode');
      }
    },
    
    loadFromLocalStorage() {
      const savedSettings = localStorage.getItem('settings');
      
      if (savedSettings) {
        try {
          const settings = JSON.parse(savedSettings);
          
          // 更新设置，但不触发API调用
          if (settings.theme) {
            this.settings.theme = {
              ...this.settings.theme,
              ...settings.theme
            };
          }
          
          if (settings.language) {
            this.settings.language = settings.language;
          }
          
          if (settings.dateFormat) {
            this.settings.dateFormat = settings.dateFormat;
          }
          
          if (settings.timeFormat) {
            this.settings.timeFormat = settings.timeFormat;
          }
          
          if (settings.timezone) {
            this.settings.timezone = settings.timezone;
          }
          
          if (settings.layout) {
            this.settings.layout = {
              ...this.settings.layout,
              ...settings.layout
            };
          }
          
          if (settings.notifications) {
            this.settings.notifications = {
              ...this.settings.notifications,
              ...settings.notifications
            };
          }
          
          // 应用设置
          this.applyThemeSettings();
          this.applyLayoutSettings();
        } catch (e) {
          console.error('解析本地设置失败:', e);
        }
      } else {
        // 如果没有本地设置，应用默认设置并保存
        this.applyThemeSettings();
        this.applyLayoutSettings();
        this.saveToLocalStorage();
      }
    },
    
    saveToLocalStorage() {
      localStorage.setItem('settings', JSON.stringify(this.settings));
    },
    
    resetToDefaults() {
      this.settings = {
        theme: {
          mode: 'light',
          color: 'blue'
        },
        language: 'zh-CN',
        dateFormat: 'yyyy-MM-dd',
        timeFormat: 'HH:mm:ss',
        timezone: 'Asia/Shanghai',
        layout: {
          sidebarPosition: 'left',
          contentWidth: 'contained',
          compactMode: false
        },
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
      
      this.applyThemeSettings();
      this.applyLayoutSettings();
      this.saveToLocalStorage();
    }
  }
}); 