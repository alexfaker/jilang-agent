import { defineStore } from 'pinia';
import { v4 as uuidv4 } from 'uuid';

export const useNotificationStore = defineStore('notification', {
  state: () => ({
    notifications: [],
    maxNotifications: 5, // 最多同时显示的通知数量
  }),
  
  actions: {
    /**
     * 添加一条通知
     * @param {Object} notification 通知对象
     * @param {string} notification.type 通知类型: 'success', 'error', 'warning', 'info'
     * @param {string} notification.title 通知标题
     * @param {string} notification.message 通知内容
     * @param {number} notification.duration 显示时长(ms)，默认4000ms，设为0则不自动关闭
     * @returns {string} 返回通知ID
     */
    add({ type = 'info', title = '', message = '', duration = 4000 }) {
      const id = uuidv4();
      const notification = {
        id,
        type,
        title,
        message,
        duration,
        timestamp: Date.now(),
        visible: true
      };
      
      // 添加到通知列表开头
      this.notifications.unshift(notification);
      
      // 如果超过最大数量，移除最旧的通知
      if (this.notifications.length > this.maxNotifications) {
        const oldNotifications = this.notifications.splice(this.maxNotifications);
        // 确保被移除的通知的DOM元素也会被移除
        oldNotifications.forEach(n => {
          n.visible = false;
        });
      }
      
      // 如果设置了持续时间，则在指定时间后自动移除
      if (duration > 0) {
        setTimeout(() => {
          this.remove(id);
        }, duration);
      }
      
      return id;
    },
    
    /**
     * 移除一条通知
     * @param {string} id 通知ID
     */
    remove(id) {
      const index = this.notifications.findIndex(n => n.id === id);
      if (index !== -1) {
        // 先将visible设为false以触发动画
        this.notifications[index].visible = false;
        
        // 然后在动画结束后移除
        setTimeout(() => {
          const currentIndex = this.notifications.findIndex(n => n.id === id);
          if (currentIndex !== -1) {
            this.notifications.splice(currentIndex, 1);
          }
        }, 300); // 假设动画持续300ms
      }
    },
    
    /**
     * 清除所有通知
     */
    clearAll() {
      // 先将所有通知的visible设为false以触发动画
      this.notifications.forEach(notification => {
        notification.visible = false;
      });
      
      // 然后在动画结束后清空
      setTimeout(() => {
        this.notifications = [];
      }, 300); // 假设动画持续300ms
    },
    
    /**
     * 添加成功通知
     */
    success(message, title = '成功', duration = 4000) {
      return this.add({ type: 'success', title, message, duration });
    },
    
    /**
     * 添加错误通知
     */
    error(message, title = '错误', duration = 6000) {
      return this.add({ type: 'error', title, message, duration });
    },
    
    /**
     * 添加警告通知
     */
    warning(message, title = '警告', duration = 5000) {
      return this.add({ type: 'warning', title, message, duration });
    },
    
    /**
     * 添加信息通知
     */
    info(message, title = '提示', duration = 4000) {
      return this.add({ type: 'info', title, message, duration });
    }
  }
}); 