import { defineStore } from 'pinia';
import { authApi } from '../api';

export const useUserStore = defineStore('user', {
  state: () => ({
    user: null,
    isAuthenticated: false,
    loading: false,
    error: null
  }),
  
  getters: {
    userProfile: (state) => state.user,
    isLoggedIn: (state) => state.isAuthenticated,
    isAdmin: (state) => state.user?.role === 'admin'
  },
  
  actions: {
    async login(credentials) {
      this.loading = true;
      this.error = null;
      
      try {
        const response = await authApi.login(credentials);
        this.user = response.data.user;
        this.isAuthenticated = true;
        localStorage.setItem('token', response.data.token);
        localStorage.setItem('user', JSON.stringify(response.data.user));
        
        // 登录成功后同步设置
        this.syncSettingsAfterLogin();
        
        return response;
      } catch (error) {
        this.error = error.message || '登录失败';
        throw error;
      } finally {
        this.loading = false;
      }
    },
    
    async register(userData) {
      this.loading = true;
      this.error = null;
      
      try {
        const response = await authApi.register(userData);
        this.user = response.data.user;
        this.isAuthenticated = true;
        localStorage.setItem('token', response.data.token);
        localStorage.setItem('user', JSON.stringify(response.data.user));
        
        // 注册成功后同步设置
        this.syncSettingsAfterLogin();
        
        return response;
      } catch (error) {
        this.error = error.message || '注册失败';
        throw error;
      } finally {
        this.loading = false;
      }
    },
    
    async logout() {
      this.loading = true;
      
      try {
        await authApi.logout();
      } catch (error) {
        console.error('注销时发生错误:', error);
      } finally {
        this.user = null;
        this.isAuthenticated = false;
        localStorage.removeItem('token');
        localStorage.removeItem('user');
        this.loading = false;
      }
    },
    
    async fetchUserProfile() {
      if (!localStorage.getItem('token')) {
        return;
      }
      
      this.loading = true;
      
      try {
        const user = await authApi.getProfile();
        this.user = user;
        this.isAuthenticated = true;
        localStorage.setItem('user', JSON.stringify(user));
      } catch (error) {
        this.error = error.message || '获取用户信息失败';
        this.user = null;
        this.isAuthenticated = false;
        localStorage.removeItem('token');
        localStorage.removeItem('user');
      } finally {
        this.loading = false;
      }
    },
    
    async updateProfile(profileData) {
      this.loading = true;
      
      try {
        let updatedUser;
        if (profileData instanceof FormData) {
          updatedUser = await authApi.updateProfileWithAvatar(profileData);
        } else {
          updatedUser = await authApi.updateProfile(profileData);
        }
        
        this.user = { ...this.user, ...updatedUser };
        localStorage.setItem('user', JSON.stringify(this.user));
        return updatedUser;
      } catch (error) {
        this.error = error.message || '更新用户信息失败';
        throw error;
      } finally {
        this.loading = false;
      }
    },
    
    async changePassword(passwordData) {
      this.loading = true;
      
      try {
        await authApi.changePassword(passwordData);
        return { success: true };
      } catch (error) {
        this.error = error.message || '修改密码失败';
        throw error;
      } finally {
        this.loading = false;
      }
    },
    
    initializeFromLocalStorage() {
      const token = localStorage.getItem('token');
      const user = localStorage.getItem('user');
      
      if (token && user) {
        try {
          this.user = JSON.parse(user);
          this.isAuthenticated = true;
        } catch (e) {
          localStorage.removeItem('user');
          localStorage.removeItem('token');
        }
      }
    },

    // 登录后同步设置的辅助方法
    syncSettingsAfterLogin() {
      try {
        // 动态导入settings store以避免循环依赖
        import('./settings').then(({ useSettingsStore }) => {
          const settingsStore = useSettingsStore();
          settingsStore.fetchSettings().catch(error => {
            console.warn('登录后同步设置失败:', error.message);
          });
        });
      } catch (error) {
        console.warn('无法同步设置:', error.message);
      }
    }
  }
}); 