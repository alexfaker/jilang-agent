<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-50 dark:bg-gray-900 py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full space-y-8">
      <div>
        <h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900 dark:text-white">
          登录到您的账户
        </h2>
        <p class="mt-2 text-center text-sm text-gray-600 dark:text-gray-400">
          还没有账户？
          <a href="#" class="font-medium text-primary-600 hover:text-primary-500">
            联系管理员
          </a>
        </p>
      </div>
      <form class="mt-8 space-y-6" @submit.prevent="handleLogin">
        <div v-if="userStore.error" class="rounded-md bg-red-50 dark:bg-red-900 p-4">
          <div class="flex">
            <div class="flex-shrink-0">
              <svg class="h-5 w-5 text-red-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
              </svg>
            </div>
            <div class="ml-3">
              <h3 class="text-sm font-medium text-red-800 dark:text-red-200">
                {{ userStore.error }}
              </h3>
            </div>
          </div>
        </div>
        
        <input type="hidden" name="remember" value="true" />
        <div class="rounded-md shadow-sm -space-y-px">
          <div>
            <label for="email-address" class="sr-only">Email 地址</label>
            <input
              id="email-address"
              name="email"
              type="email" 
              autocomplete="email"
              required
              v-model="credentials.email"
              :class="{ 'border-red-300': validationErrors.email }"
              class="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-300 dark:border-gray-700 dark:bg-gray-800 dark:text-white placeholder-gray-500 dark:placeholder-gray-400 text-gray-900 rounded-t-md focus:outline-none focus:ring-primary-500 focus:border-primary-500 focus:z-10 sm:text-sm"
              placeholder="邮箱地址"
            />
            <p v-if="validationErrors.email" class="mt-2 text-sm text-red-600 dark:text-red-400">{{ validationErrors.email }}</p>
          </div>
          <div>
            <label for="password" class="sr-only">密码</label>
            <input
              id="password"
              name="password"
              type="password"
              autocomplete="current-password"
              required
              v-model="credentials.password"
              :class="{ 'border-red-300': validationErrors.password }"
              class="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-300 dark:border-gray-700 dark:bg-gray-800 dark:text-white placeholder-gray-500 dark:placeholder-gray-400 text-gray-900 rounded-b-md focus:outline-none focus:ring-primary-500 focus:border-primary-500 focus:z-10 sm:text-sm"
              placeholder="密码"
            />
            <p v-if="validationErrors.password" class="mt-2 text-sm text-red-600 dark:text-red-400">{{ validationErrors.password }}</p>
          </div>
        </div>

        <div class="flex items-center justify-between">
          <div class="flex items-center">
            <input
              id="remember-me"
              name="remember-me"
              type="checkbox"
              v-model="credentials.remember"
              class="h-4 w-4 text-primary-600 focus:ring-primary-500 border-gray-300 dark:border-gray-700 dark:bg-gray-800 rounded"
            />
            <label for="remember-me" class="ml-2 block text-sm text-gray-900 dark:text-gray-300">
              记住我
            </label>
          </div>

          <div class="text-sm">
            <a href="#" class="font-medium text-primary-600 hover:text-primary-500">
              忘记密码?
            </a>
          </div>
        </div>

        <div>
          <button
            type="submit"
            :disabled="userStore.loading"
            class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
          >
            <span class="absolute left-0 inset-y-0 flex items-center pl-3">
              <svg
                v-if="!userStore.loading"
                class="h-5 w-5 text-primary-500 group-hover:text-primary-400"
                xmlns="http://www.w3.org/2000/svg"
                viewBox="0 0 20 20"
                fill="currentColor"
                aria-hidden="true"
              >
                <path
                  fill-rule="evenodd"
                  d="M5 9V7a5 5 0 0110 0v2a2 2 0 012 2v5a2 2 0 01-2 2H5a2 2 0 01-2-2v-5a2 2 0 012-2zm8-2v2H7V7a3 3 0 016 0z"
                  clip-rule="evenodd"
                />
              </svg>
              <svg 
                v-else 
                class="animate-spin h-5 w-5 text-white" 
                xmlns="http://www.w3.org/2000/svg" 
                fill="none" 
                viewBox="0 0 24 24"
              >
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
            </span>
            {{ userStore.loading ? '登录中...' : '登 录' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useUserStore } from '../../stores/user';
import { useToast } from 'vue-toastification';
import notify from '../../utils/notification';

const router = useRouter();
const route = useRoute();
const userStore = useUserStore();
const toast = useToast();

const credentials = reactive({
  email: '',
  password: '',
  remember: false
});

const validationErrors = ref({});

const validateForm = () => {
  const errors = {};
  if (!credentials.email) {
    errors.email = '请输入邮箱地址';
  } else if (!/^\S+@\S+\.\S+$/.test(credentials.email)) {
    errors.email = '请输入有效的邮箱地址';
  }
  
  if (!credentials.password) {
    errors.password = '请输入密码';
  } else if (credentials.password.length < 6) {
    errors.password = '密码长度不能小于6位';
  }
  
  validationErrors.value = errors;
  return Object.keys(errors).length === 0;
};

const handleLogin = async () => {
  // 重置错误
  userStore.error = null;
  validationErrors.value = {};
  
  // 表单验证
  if (!validateForm()) {
    return;
  }
  
  try {
    // 调用用户存储的登录方法
    await userStore.login({
      email: credentials.email,
      password: credentials.password,
      remember: credentials.remember
    });
    
    // 登录成功通知
    notify.success('登录成功，欢迎回来！');
    
    // 重定向到之前尝试访问的页面或默认到仪表盘
    const redirectPath = route.query.redirect || '/dashboard';
    router.push(redirectPath);
  } catch (error) {
    // 错误已在用户存储中处理，这里只需显示通知
    notify.error(userStore.error || '登录失败，请检查您的邮箱和密码');
  }
};
</script> 