import { defineStore } from 'pinia';

export const useLoadingStore = defineStore('loading', {
  state: () => ({
    isLoading: false,
    message: '',
    progress: 0,
    loadingStates: {}, // 用于存储各个模块的加载状态
    loadingCount: 0, // 记录当前正在加载的请求数量
  }),
  
  getters: {
    // 判断指定模块是否正在加载
    isModuleLoading: (state) => (module) => {
      return !!state.loadingStates[module];
    },
    
    // 获取所有正在加载的模块
    activeLoadingModules: (state) => {
      return Object.keys(state.loadingStates).filter(key => state.loadingStates[key]);
    }
  },
  
  actions: {
    /**
     * 开始全局加载
     * @param {string} message 加载提示消息
     */
    startLoading(message = '加载中...') {
      this.isLoading = true;
      this.message = message;
      this.progress = 0;
    },
    
    /**
     * 更新加载进度
     * @param {number} progress 进度值(0-100)
     * @param {string} message 可选的新消息
     */
    updateProgress(progress, message = null) {
      this.progress = Math.min(Math.max(progress, 0), 100);
      if (message !== null) {
        this.message = message;
      }
    },
    
    /**
     * 结束全局加载
     */
    endLoading() {
      // 先将进度设为100%，然后再结束加载
      this.progress = 100;
      
      // 延迟一小段时间再隐藏加载器，让用户看到100%的状态
      setTimeout(() => {
        this.isLoading = false;
        this.message = '';
        this.progress = 0;
      }, 300);
    },
    
    /**
     * 开始模块加载
     * @param {string} module 模块名称
     * @param {string} message 加载提示消息
     */
    startModuleLoading(module, message = '') {
      this.loadingStates[module] = { active: true, message };
      this.loadingCount++;
    },
    
    /**
     * 结束模块加载
     * @param {string} module 模块名称
     */
    endModuleLoading(module) {
      if (this.loadingStates[module]?.active) {
        this.loadingStates[module].active = false;
        this.loadingCount = Math.max(0, this.loadingCount - 1);
      }
    },
    
    /**
     * 清除所有模块的加载状态
     */
    clearAllModuleLoading() {
      this.loadingStates = {};
      this.loadingCount = 0;
    },
    
    /**
     * 为异步操作添加加载状态
     * @param {string} module 模块名称
     * @param {Promise} promise 要执行的Promise
     * @param {string} message 加载提示消息
     * @returns {Promise} 原始Promise的结果
     */
    async withLoading(module, promise, message = '') {
      this.startModuleLoading(module, message);
      
      try {
        const result = await promise;
        return result;
      } finally {
        this.endModuleLoading(module);
      }
    },
    
    /**
     * 为异步操作添加全局加载状态
     * @param {Promise} promise 要执行的Promise
     * @param {string} message 加载提示消息
     * @returns {Promise} 原始Promise的结果
     */
    async withGlobalLoading(promise, message = '加载中...') {
      this.startLoading(message);
      
      try {
        const result = await promise;
        return result;
      } finally {
        this.endLoading();
      }
    }
  }
}); 