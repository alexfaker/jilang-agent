<template>
  <div class="notifications-container">
    <TransitionGroup name="notification">
      <div
        v-for="notification in notifications"
        :key="notification.id"
        :class="['notification', `notification-${notification.type}`]"
        v-show="notification.visible"
      >
        <div class="notification-icon">
          <i v-if="notification.type === 'success'" class="fas fa-check-circle"></i>
          <i v-else-if="notification.type === 'error'" class="fas fa-exclamation-circle"></i>
          <i v-else-if="notification.type === 'warning'" class="fas fa-exclamation-triangle"></i>
          <i v-else class="fas fa-info-circle"></i>
        </div>
        
        <div class="notification-content">
          <div v-if="notification.title" class="notification-title">{{ notification.title }}</div>
          <div class="notification-message">{{ notification.message }}</div>
        </div>
        
        <button class="notification-close" @click="removeNotification(notification.id)">
          <i class="fas fa-times"></i>
        </button>
      </div>
    </TransitionGroup>
  </div>
</template>

<script>
import { defineComponent, computed } from 'vue';
import { storeToRefs } from 'pinia';
import { useNotificationStore } from '@/stores';

export default defineComponent({
  name: 'NotificationComponent',
  
  setup() {
    const notificationStore = useNotificationStore();
    const { notifications } = storeToRefs(notificationStore);
    
    const removeNotification = (id) => {
      notificationStore.remove(id);
    };
    
    return {
      notifications,
      removeNotification
    };
  }
});
</script>

<style scoped>
.notifications-container {
  position: fixed;
  top: 20px;
  right: 20px;
  z-index: 9999;
  width: 350px;
  max-width: calc(100vw - 40px);
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.notification {
  display: flex;
  align-items: flex-start;
  padding: 12px 16px;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  background-color: white;
  margin-bottom: 10px;
  position: relative;
  overflow: hidden;
}

.notification-icon {
  margin-right: 12px;
  font-size: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.notification-content {
  flex: 1;
  min-width: 0;
}

.notification-title {
  font-weight: 600;
  margin-bottom: 4px;
}

.notification-message {
  font-size: 14px;
  word-break: break-word;
}

.notification-close {
  background: transparent;
  border: none;
  cursor: pointer;
  font-size: 16px;
  padding: 4px;
  margin-left: 8px;
  color: #999;
  transition: color 0.2s;
}

.notification-close:hover {
  color: #333;
}

/* 通知类型样式 */
.notification-success {
  border-left: 4px solid #4caf50;
}

.notification-success .notification-icon {
  color: #4caf50;
}

.notification-error {
  border-left: 4px solid #f44336;
}

.notification-error .notification-icon {
  color: #f44336;
}

.notification-warning {
  border-left: 4px solid #ff9800;
}

.notification-warning .notification-icon {
  color: #ff9800;
}

.notification-info {
  border-left: 4px solid #2196f3;
}

.notification-info .notification-icon {
  color: #2196f3;
}

/* 动画效果 */
.notification-enter-active,
.notification-leave-active {
  transition: all 0.3s ease;
}

.notification-enter-from {
  opacity: 0;
  transform: translateX(30px);
}

.notification-leave-to {
  opacity: 0;
  transform: translateY(-30px);
}

/* 暗黑模式样式 */
:deep(.dark-theme) .notification {
  background-color: #333;
  color: #fff;
}

:deep(.dark-theme) .notification-close {
  color: #ccc;
}

:deep(.dark-theme) .notification-close:hover {
  color: #fff;
}
</style> 