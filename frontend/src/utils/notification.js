import { useToast } from 'vue-toastification';

/**
 * 通知辅助函数
 * 提供一种更简单的方式在组件外部使用通知功能
 * 基于vue-toastification库
 */
export const notify = {
  /**
   * 显示成功通知
   * @param {string} message 通知消息
   * @param {Object} options 额外选项
   * @returns {number} 通知ID
   */
  success(message, options = {}) {
    const toast = useToast();
    return toast.success(message, {
      timeout: 4000,
      ...options
    });
  },
  
  /**
   * 显示错误通知
   * @param {string} message 通知消息
   * @param {Object} options 额外选项
   * @returns {number} 通知ID
   */
  error(message, options = {}) {
    const toast = useToast();
    return toast.error(message, {
      timeout: 6000,
      ...options
    });
  },
  
  /**
   * 显示警告通知
   * @param {string} message 通知消息
   * @param {Object} options 额外选项
   * @returns {number} 通知ID
   */
  warning(message, options = {}) {
    const toast = useToast();
    return toast.warning(message, {
      timeout: 5000,
      ...options
    });
  },
  
  /**
   * 显示信息通知
   * @param {string} message 通知消息
   * @param {Object} options 额外选项
   * @returns {number} 通知ID
   */
  info(message, options = {}) {
    const toast = useToast();
    return toast.info(message, {
      timeout: 4000,
      ...options
    });
  },
  
  /**
   * 显示API错误通知
   * 自动处理常见的API错误情况
   * @param {Error} error 错误对象
   * @param {string} fallbackMessage 当无法解析错误时的后备消息
   * @param {Object} options 额外选项
   * @returns {number} 通知ID
   */
  apiError(error, fallbackMessage = '操作失败，请稍后重试', options = {}) {
    // 尝试获取错误详情
    let message = fallbackMessage;
    
    if (error.response && error.response.data) {
      // 优先使用服务器返回的错误消息
      message = error.response.data.message || error.response.data.error || fallbackMessage;
    } else if (error.message) {
      // 使用错误对象的消息
      message = error.message;
    }
    
    return this.error(message, {
      timeout: 6000,
      ...options
    });
  },
  
  /**
   * 清除所有通知
   */
  clearAll() {
    const toast = useToast();
    toast.clear();
  },
  
  /**
   * 清除特定通知
   * @param {number} id 通知ID
   */
  clear(id) {
    const toast = useToast();
    toast.dismiss(id);
  }
};

/**
 * 创建一个通知的Vue插件
 * 使用方法：app.use(notificationPlugin)
 * 然后可以在组件中使用 this.$notify
 */
export const notificationPlugin = {
  install(app) {
    app.config.globalProperties.$notify = notify;
  }
};

export default notify; 