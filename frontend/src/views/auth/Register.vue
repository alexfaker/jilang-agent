<template>
  <div class="min-h-screen bg-gray-50">
    <!-- 顶部导航栏 -->
    <nav class="bg-white shadow-sm">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-16">
          <div class="flex items-center">
            <router-link to="/" class="flex-shrink-0">
              <BrandLogo size="base" />
            </router-link>
          </div>
          <div class="flex items-center">
            <router-link 
              to="/auth/login" 
              class="text-gray-500 hover:text-gray-700 px-3 py-2 rounded-md text-sm font-medium transition-colors"
            >
              登录
            </router-link>
          </div>
        </div>
      </div>
    </nav>

    <!-- 注册表单区域 -->
    <div class="flex-grow flex items-center justify-center px-4 sm:px-6 lg:px-8 py-12">
      <div class="max-w-md w-full">
        <div class="auth-card bg-white p-8 rounded-2xl shadow-lg">
          <!-- Logo和注册标题 -->
          <div class="flex items-center justify-center mb-8">
            <BrandLogo size="lg" spacing="loose" />
          </div>

          <!-- 错误提示 -->
          <div v-if="notification.show" :class="[
            'rounded-md p-4 mb-6',
            notification.type === 'success' ? 'bg-green-50' : 'bg-red-50'
          ]">
            <div class="flex">
              <div class="flex-shrink-0">
                <CheckCircleIcon v-if="notification.type === 'success'" class="h-5 w-5 text-green-400" />
                <ExclamationTriangleIcon v-else class="h-5 w-5 text-red-400" />
              </div>
              <div class="ml-3">
                <h3 :class="[
                  'text-sm font-medium',
                  notification.type === 'success' ? 'text-green-800' : 'text-red-800'
                ]">
                  {{ notification.message }}
                </h3>
              </div>
              <div class="ml-auto pl-3">
                <button @click="notification.show = false" :class="[
                  'inline-flex text-sm',
                  notification.type === 'success' ? 'text-green-500 hover:text-green-400' : 'text-red-500 hover:text-red-400'
                ]">
                  <XMarkIcon class="h-4 w-4" />
                </button>
              </div>
            </div>
          </div>
          
          <!-- 注册表单 -->
          <form class="space-y-6" @submit.prevent="handleRegister">
            <!-- 用户名输入 -->
            <div>
              <label for="username" class="block text-sm font-medium text-gray-700">用户名</label>
              <div class="mt-1 relative">
                <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                  <UserIcon class="h-5 w-5 text-gray-400" />
                </div>
                <input 
                  id="username" 
                  name="username" 
                  type="text" 
                  autocomplete="username" 
                  required
                  v-model="form.username"
                  :class="{ 'border-red-300': errors.username }"
                  class="auth-input appearance-none block w-full pl-10 pr-3 py-3 border border-gray-300 rounded-lg placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all"
                  placeholder="请输入用户名"
                />
              </div>
              <p v-if="errors.username" class="mt-2 text-sm text-red-600">{{ errors.username }}</p>
            </div>

            <!-- 邮箱输入 -->
            <div>
              <label for="email" class="block text-sm font-medium text-gray-700">邮箱地址</label>
              <div class="mt-1 relative">
                <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                  <EnvelopeIcon class="h-5 w-5 text-gray-400" />
                </div>
                <input 
                  id="email" 
                  name="email" 
                  type="email" 
                  autocomplete="email" 
                  required
                  v-model="form.email"
                  :class="{ 'border-red-300': errors.email }"
                  class="auth-input appearance-none block w-full pl-10 pr-3 py-3 border border-gray-300 rounded-lg placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all"
                  placeholder="请输入邮箱地址"
                />
              </div>
              <p v-if="errors.email" class="mt-2 text-sm text-red-600">{{ errors.email }}</p>
            </div>

            <!-- 密码输入 -->
            <div>
              <label for="password" class="block text-sm font-medium text-gray-700">密码</label>
              <div class="mt-1 relative">
                <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                  <LockClosedIcon class="h-5 w-5 text-gray-400" />
                </div>
                <input 
                  id="password" 
                  name="password" 
                  :type="showPassword ? 'text' : 'password'" 
                  autocomplete="new-password" 
                  required
                  v-model="form.password"
                  :class="{ 'border-red-300': errors.password }"
                  class="auth-input appearance-none block w-full pl-10 pr-10 py-3 border border-gray-300 rounded-lg placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all"
                  placeholder="请输入密码（至少8位）"
                />
                <button 
                  type="button" 
                  @click="togglePassword"
                  class="absolute inset-y-0 right-0 pr-3 flex items-center"
                >
                  <EyeIcon v-if="!showPassword" class="h-5 w-5 text-gray-400 hover:text-gray-600" />
                  <EyeSlashIcon v-else class="h-5 w-5 text-gray-400 hover:text-gray-600" />
                </button>
              </div>
              <p class="mt-1 text-xs text-gray-500">密码长度至少8位，包含字母和数字</p>
              <p v-if="errors.password" class="mt-2 text-sm text-red-600">{{ errors.password }}</p>
            </div>

            <!-- 确认密码 -->
            <div>
              <label for="confirmPassword" class="block text-sm font-medium text-gray-700">确认密码</label>
              <div class="mt-1 relative">
                <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                  <LockClosedIcon class="h-5 w-5 text-gray-400" />
                </div>
                <input 
                  id="confirmPassword" 
                  name="confirmPassword" 
                  type="password" 
                  autocomplete="new-password" 
                  required
                  v-model="form.confirmPassword"
                  :class="{ 'border-red-300': errors.confirmPassword }"
                  class="auth-input appearance-none block w-full pl-10 pr-3 py-3 border border-gray-300 rounded-lg placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all"
                  placeholder="请再次输入密码"
                />
              </div>
              <p v-if="errors.confirmPassword" class="mt-2 text-sm text-red-600">{{ errors.confirmPassword }}</p>
            </div>

            <!-- 条款同意 -->
            <div class="flex items-start">
              <div class="flex items-center h-5">
                <input 
                  id="terms" 
                  name="terms" 
                  type="checkbox" 
                  required
                  v-model="form.agreeToTerms"
                  class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                />
              </div>
              <div class="ml-3 text-sm">
                <label for="terms" class="font-medium text-gray-700">我同意
                  <a href="#" class="text-indigo-600 hover:text-indigo-500">服务条款</a>
                  和
                  <a href="#" class="text-indigo-600 hover:text-indigo-500">隐私政策</a>
                </label>
              </div>
            </div>

            <!-- 注册按钮 -->
            <div>
              <button 
                type="submit"
                :disabled="loading"
                class="auth-btn w-full flex justify-center py-3 px-4 border border-transparent text-sm font-medium rounded-lg text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 transition-all hover:transform hover:-translate-y-0.5"
              >
                <span v-if="loading" class="flex items-center">
                  <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                  注册中...
                </span>
                <span v-else>注册</span>
              </button>
            </div>
          </form>

          <!-- 社交注册 -->
          <div class="mt-6">
            <div class="relative">
              <div class="absolute inset-0 flex items-center">
                <div class="w-full border-t border-gray-300"></div>
              </div>
              <div class="relative flex justify-center text-sm">
                <span class="px-2 bg-white text-gray-500">
                  或通过以下方式注册
                </span>
              </div>
            </div>

            <div class="mt-6 grid grid-cols-3 gap-3">
              <button
                type="button"
                class="w-full inline-flex justify-center py-2 px-4 border border-gray-300 rounded-md shadow-sm bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 transition-colors"
              >
                <svg class="h-5 w-5" fill="currentColor" viewBox="0 0 24 24">
                  <path d="M12 0C5.37 0 0 5.37 0 12 0 17.31 3.435 21.795 8.205 23.385 8.805 23.49 9.03 23.13 9.03 22.815L9.015 20.91C5.672 21.675 4.968 19.35 4.968 19.35 4.422 17.94 3.633 17.583 3.633 17.583 2.546 16.791 3.717 16.806 3.717 16.806 4.922 16.896 5.555 18.084 5.555 18.084 6.625 19.932 8.364 19.401 9.048 19.088 9.155 18.311 9.466 17.78 9.81 17.475 7.145 17.167 4.344 16.125 4.344 11.475 4.344 10.147 4.809 9.057 5.579 8.207 5.444 7.896 5.039 6.651 5.684 4.992 5.684 4.992 6.689 4.651 8.984 6.207 10.189 5.962 11.496 5.962 12.803 6.207 15.097 4.651 16.103 4.992 16.103 4.992 16.748 6.651 16.343 7.896 16.208 8.207 16.978 9.057 17.443 10.147 17.443 11.475 17.443 16.134 14.641 17.167 11.976 17.475 12.32 17.78 12.631 18.311 12.738 19.088 13.422 19.401 15.161 19.932 16.231 18.084 16.864 16.896 18.069 16.806 18.069 16.806 19.24 16.791 18.153 17.583 17.364 17.94 16.818 19.35 16.818 19.35 16.114 21.675 12.97 20.91L12.955 22.815C12.955 23.13 13.18 23.49 13.78 23.385 18.55 21.795 21.985 17.31 21.985 12 22 5.37 16.63.015 12 0Z"/>
                </svg>
              </button>
              <button
                type="button"
                class="w-full inline-flex justify-center py-2 px-4 border border-gray-300 rounded-md shadow-sm bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 transition-colors"
              >
                <svg class="h-5 w-5" fill="currentColor" viewBox="0 0 24 24">
                  <path d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/>
                  <path d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
                  <path d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/>
                  <path d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
                </svg>
              </button>
              <button
                type="button"
                class="w-full inline-flex justify-center py-2 px-4 border border-gray-300 rounded-md shadow-sm bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 transition-colors"
              >
                <svg class="h-5 w-5" fill="currentColor" viewBox="0 0 24 24">
                  <path d="M18.73 5.73h-5.06V3.59c0-1.35.87-1.66 1.48-1.66h3.45V.11S17.56 0 16.53 0c-2.68 0-4.37 1.04-4.37 3.14v2.59H9.33v1.94h2.83v7.75c.85.13 1.72.2 2.6.2.88 0 1.75-.07 2.6-.2V7.67h2.25l.37-1.94z"/>
                </svg>
              </button>
            </div>
          </div>
        </div>

        <!-- 登录提示 -->
        <p class="mt-6 text-center text-sm text-gray-600">
          已有账户?
          <router-link to="/auth/login" class="font-medium text-indigo-600 hover:text-indigo-500">
            登录
          </router-link>
        </p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref, computed, watch } from 'vue';
import { useRouter } from 'vue-router';
import { authApi } from '../../api/index';
import { 
  CpuChipIcon, 
  EnvelopeIcon, 
  LockClosedIcon, 
  UserIcon,
  EyeIcon, 
  EyeSlashIcon,
  ExclamationTriangleIcon,
  CheckCircleIcon,
  XMarkIcon 
} from '@heroicons/vue/24/outline';
import BrandLogo from '../../components/common/BrandLogo.vue';

const router = useRouter();

// 表单数据
const form = reactive({
  username: '',
  email: '',
  password: '',
  confirmPassword: '',
  agreeToTerms: false
});

// 状态管理
const loading = ref(false);
const showPassword = ref(false);
const errors = reactive({});

// 通知
const notification = reactive({
  show: false,
  type: 'error',
  message: ''
});

// 切换密码显示
const togglePassword = () => {
  showPassword.value = !showPassword.value;
};

// 显示通知
const showNotification = (type, message) => {
  notification.show = true;
  notification.type = type;
  notification.message = message;
  
  setTimeout(() => {
    notification.show = false;
  }, 5000);
};

// 验证函数
const validateForm = () => {
  const newErrors = {};

  // 用户名验证
  if (!form.username.trim()) {
    newErrors.username = '用户名不能为空';
  } else if (form.username.length < 3) {
    newErrors.username = '用户名至少3个字符';
  } else if (!/^[a-zA-Z0-9_]+$/.test(form.username)) {
    newErrors.username = '用户名只能包含字母、数字和下划线';
  }

  // 邮箱验证
  if (!form.email.trim()) {
    newErrors.email = '邮箱不能为空';
  } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.email)) {
    newErrors.email = '请输入有效的邮箱地址';
  }

  // 密码验证
  if (!form.password) {
    newErrors.password = '密码不能为空';
  } else if (form.password.length < 8) {
    newErrors.password = '密码长度至少8位';
  }

  // 确认密码验证
  if (!form.confirmPassword) {
    newErrors.confirmPassword = '请确认密码';
  } else if (form.password !== form.confirmPassword) {
    newErrors.confirmPassword = '两次输入的密码不一致';
  }

  // 更新错误状态
  Object.keys(errors).forEach(key => delete errors[key]);
  Object.assign(errors, newErrors);

  return Object.keys(newErrors).length === 0;
};

// 注册处理
const handleRegister = async () => {
  if (!validateForm()) {
    showNotification('error', '请检查表单输入');
    return;
  }

  if (!form.agreeToTerms) {
    showNotification('error', '请同意服务条款和隐私政策');
    return;
  }

  loading.value = true;
  
  try {
    const registerData = {
      username: form.username.trim(),
      email: form.email.trim(),
      password: form.password
    };

    const response = await authApi.register(registerData);
    
    showNotification('success', '注册成功！正在跳转到登录页面...');
    
    // 延迟跳转，让用户看到成功消息
    setTimeout(() => {
      router.push({
        name: 'Login',
        query: { message: 'registration_success' }
      });
    }, 2000);
    
  } catch (error) {
    console.error('注册失败:', error);
    
    let errorMessage = '注册失败，请稍后重试';
    
    if (error.response && error.response.data && error.response.data.message) {
      errorMessage = error.response.data.message;
    } else if (error.message) {
      errorMessage = error.message;
    }
    
    showNotification('error', errorMessage);
  } finally {
    loading.value = false;
  }
};
</script>

<style scoped>
.animate-spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

.auth-card {
  transition: all 0.3s ease;
}

.auth-input {
  transition: all 0.2s ease;
  height: 3rem;
}

.auth-btn {
  transition: all 0.2s ease;
  height: 3rem;
}

.auth-btn:hover:not(:disabled) {
  transform: translateY(-2px);
}
</style> 