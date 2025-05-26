<template>
  <div class="space-y-8">
    <!-- 个人资料表单 -->
    <div class="bg-white rounded-lg border border-gray-200 p-6">
      <h3 class="text-lg font-medium mb-4">个人资料</h3>
      <form @submit.prevent="updateProfile" class="space-y-4">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <!-- 用户名 -->
          <div>
            <label for="username" class="block text-sm font-medium text-gray-700 mb-1">用户名</label>
            <input
              id="username"
              v-model="profileForm.username"
              type="text"
              class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-primary-500 focus:border-primary-500"
              :disabled="true"
            />
            <p class="text-xs text-gray-500 mt-1">用户名不可修改</p>
          </div>
          
          <!-- 昵称 -->
          <div>
            <label for="nickname" class="block text-sm font-medium text-gray-700 mb-1">昵称</label>
            <input
              id="nickname"
              v-model="profileForm.nickname"
              type="text"
              class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-primary-500 focus:border-primary-500"
            />
          </div>
          
          <!-- 邮箱 -->
          <div>
            <label for="email" class="block text-sm font-medium text-gray-700 mb-1">邮箱</label>
            <input
              id="email"
              v-model="profileForm.email"
              type="email"
              class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-primary-500 focus:border-primary-500"
            />
            <p v-if="errors.email" class="text-xs text-red-500 mt-1">{{ errors.email }}</p>
          </div>
          
          <!-- 手机号 -->
          <div>
            <label for="phone" class="block text-sm font-medium text-gray-700 mb-1">手机号</label>
            <input
              id="phone"
              v-model="profileForm.phone"
              type="text"
              class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-primary-500 focus:border-primary-500"
            />
            <p v-if="errors.phone" class="text-xs text-red-500 mt-1">{{ errors.phone }}</p>
          </div>
        </div>
        
        <!-- 提交按钮 -->
        <div class="flex justify-end">
          <button
            type="submit"
            class="px-4 py-2 bg-primary-600 text-white rounded-md hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2"
            :disabled="isProfileSubmitting"
          >
            <span v-if="isProfileSubmitting" class="flex items-center">
              <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              保存中...
            </span>
            <span v-else>保存修改</span>
          </button>
        </div>
      </form>
    </div>
    
    <!-- 修改密码表单 -->
    <div class="bg-white rounded-lg border border-gray-200 p-6">
      <h3 class="text-lg font-medium mb-4">修改密码</h3>
      <form @submit.prevent="changePassword" class="space-y-4">
        <div class="grid grid-cols-1 gap-4">
          <!-- 当前密码 -->
          <div>
            <label for="currentPassword" class="block text-sm font-medium text-gray-700 mb-1">当前密码</label>
            <input
              id="currentPassword"
              v-model="passwordForm.currentPassword"
              type="password"
              class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-primary-500 focus:border-primary-500"
            />
            <p v-if="errors.currentPassword" class="text-xs text-red-500 mt-1">{{ errors.currentPassword }}</p>
          </div>
          
          <!-- 新密码 -->
          <div>
            <label for="newPassword" class="block text-sm font-medium text-gray-700 mb-1">新密码</label>
            <input
              id="newPassword"
              v-model="passwordForm.newPassword"
              type="password"
              class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-primary-500 focus:border-primary-500"
            />
            <p v-if="errors.newPassword" class="text-xs text-red-500 mt-1">{{ errors.newPassword }}</p>
            <p class="text-xs text-gray-500 mt-1">密码至少包含8个字符，且必须包含字母和数字</p>
          </div>
          
          <!-- 确认新密码 -->
          <div>
            <label for="confirmPassword" class="block text-sm font-medium text-gray-700 mb-1">确认新密码</label>
            <input
              id="confirmPassword"
              v-model="passwordForm.confirmPassword"
              type="password"
              class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-primary-500 focus:border-primary-500"
            />
            <p v-if="errors.confirmPassword" class="text-xs text-red-500 mt-1">{{ errors.confirmPassword }}</p>
          </div>
        </div>
        
        <!-- 提交按钮 -->
        <div class="flex justify-end">
          <button
            type="submit"
            class="px-4 py-2 bg-primary-600 text-white rounded-md hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2"
            :disabled="isPasswordSubmitting"
          >
            <span v-if="isPasswordSubmitting" class="flex items-center">
              <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              修改中...
            </span>
            <span v-else>修改密码</span>
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue';
import { userApi } from '../../api/index';
import { useToast } from 'vue-toastification';

// Toast通知
const toast = useToast();

// 表单数据
const profileForm = reactive({
  username: '',
  nickname: '',
  email: '',
  phone: ''
});

const passwordForm = reactive({
  currentPassword: '',
  newPassword: '',
  confirmPassword: ''
});

// 表单错误
const errors = reactive({});

// 提交状态
const isProfileSubmitting = ref(false);
const isPasswordSubmitting = ref(false);

// 获取用户信息
onMounted(async () => {
  try {
    const userData = await userApi.getCurrentUser();
    profileForm.username = userData.username;
    profileForm.nickname = userData.nickname || '';
    profileForm.email = userData.email || '';
    profileForm.phone = userData.phone || '';
  } catch (error) {
    toast.error('获取用户信息失败：' + error.message);
  }
});

// 更新个人资料
const updateProfile = async () => {
  // 清除错误
  errors.email = null;
  errors.phone = null;
  
  // 表单验证
  let isValid = true;
  
  // 验证邮箱
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  if (profileForm.email && !emailRegex.test(profileForm.email)) {
    errors.email = '请输入有效的邮箱地址';
    isValid = false;
  }
  
  // 验证手机号
  const phoneRegex = /^1[3-9]\d{9}$/;
  if (profileForm.phone && !phoneRegex.test(profileForm.phone)) {
    errors.phone = '请输入有效的手机号码';
    isValid = false;
  }
  
  if (!isValid) return;
  
  isProfileSubmitting.value = true;
  
  try {
    await userApi.updateProfile({
      nickname: profileForm.nickname,
      email: profileForm.email,
      phone: profileForm.phone
    });
    
    toast.success('个人资料更新成功');
  } catch (error) {
    toast.error('更新失败：' + error.message);
  } finally {
    isProfileSubmitting.value = false;
  }
};

// 修改密码
const changePassword = async () => {
  // 清除错误
  errors.currentPassword = null;
  errors.newPassword = null;
  errors.confirmPassword = null;
  
  // 表单验证
  let isValid = true;
  
  if (!passwordForm.currentPassword) {
    errors.currentPassword = '请输入当前密码';
    isValid = false;
  }
  
  // 验证新密码
  const passwordRegex = /^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d]{8,}$/;
  if (!passwordRegex.test(passwordForm.newPassword)) {
    errors.newPassword = '密码至少包含8个字符，且必须包含字母和数字';
    isValid = false;
  }
  
  // 验证确认密码
  if (passwordForm.newPassword !== passwordForm.confirmPassword) {
    errors.confirmPassword = '两次输入的密码不一致';
    isValid = false;
  }
  
  if (!isValid) return;
  
  isPasswordSubmitting.value = true;
  
  try {
    await userApi.changePassword({
      currentPassword: passwordForm.currentPassword,
      newPassword: passwordForm.newPassword
    });
    
    toast.success('密码修改成功');
    
    // 清空表单
    passwordForm.currentPassword = '';
    passwordForm.newPassword = '';
    passwordForm.confirmPassword = '';
  } catch (error) {
    if (error.status === 401) {
      errors.currentPassword = '当前密码不正确';
    } else {
      toast.error('修改失败：' + error.message);
    }
  } finally {
    isPasswordSubmitting.value = false;
  }
};
</script> 