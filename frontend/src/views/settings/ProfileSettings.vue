<template>
  <div>
    <h2 class="text-lg font-medium text-gray-900 dark:text-gray-100">个人资料</h2>
    <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
      更新您的个人信息和账户设置
    </p>

    <form @submit.prevent="updateProfile" class="mt-6 space-y-6">
      <div class="grid grid-cols-1 gap-y-6 gap-x-4 sm:grid-cols-6">
        <!-- 头像 -->
        <div class="sm:col-span-6">
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300">头像</label>
          <div class="mt-2 flex items-center">
            <div v-if="profileImage" class="h-12 w-12 rounded-full overflow-hidden bg-gray-100 dark:bg-gray-700">
              <img :src="profileImage" alt="Profile" class="h-full w-full object-cover" />
            </div>
            <div v-else class="h-12 w-12 rounded-full bg-primary-600 flex items-center justify-center text-white">
              {{ userInitials }}
            </div>
            <button
              type="button"
              class="ml-5 bg-white dark:bg-gray-700 py-2 px-3 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm text-sm leading-4 font-medium text-gray-700 dark:text-gray-200 hover:bg-gray-50 dark:hover:bg-gray-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
              @click="$refs.fileInput.click()"
            >
              更改
            </button>
            <input
              ref="fileInput"
              type="file"
              class="hidden"
              accept="image/*"
              @change="handleFileChange"
            />
          </div>
          <p v-if="fileError" class="mt-2 text-sm text-red-600 dark:text-red-400">
            {{ fileError }}
          </p>
          <p v-else class="mt-2 text-xs text-gray-500 dark:text-gray-400">
            支持JPG、PNG格式，文件大小不超过2MB
          </p>
        </div>

        <!-- 用户名 -->
        <div class="sm:col-span-3">
          <label for="username" class="block text-sm font-medium text-gray-700 dark:text-gray-300">用户名</label>
          <div class="mt-1">
            <input
              id="username"
              v-model="form.username"
              type="text"
              autocomplete="username"
              class="shadow-sm focus:ring-primary-500 focus:border-primary-500 block w-full sm:text-sm border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white rounded-md"
              :disabled="true"
            />
          </div>
          <p class="mt-2 text-sm text-gray-500 dark:text-gray-400">用户名不可更改</p>
        </div>

        <!-- 昵称 -->
        <div class="sm:col-span-3">
          <label for="nickname" class="block text-sm font-medium text-gray-700 dark:text-gray-300">昵称</label>
          <div class="mt-1">
            <input
              id="nickname"
              v-model="form.nickname"
              type="text"
              autocomplete="nickname"
              class="shadow-sm focus:ring-primary-500 focus:border-primary-500 block w-full sm:text-sm border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white rounded-md"
              :class="{ 'border-red-300 dark:border-red-700': v$.form.nickname.$error }"
            />
          </div>
          <p v-if="v$.form.nickname.$error" class="mt-2 text-sm text-red-600 dark:text-red-400">
            {{ v$.form.nickname.$errors[0].$message }}
          </p>
        </div>

        <!-- 邮箱 -->
        <div class="sm:col-span-4">
          <label for="email" class="block text-sm font-medium text-gray-700 dark:text-gray-300">邮箱地址</label>
          <div class="mt-1">
            <input
              id="email"
              v-model="form.email"
              type="email"
              autocomplete="email"
              class="shadow-sm focus:ring-primary-500 focus:border-primary-500 block w-full sm:text-sm border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white rounded-md"
              :class="{ 'border-red-300 dark:border-red-700': v$.form.email.$error }"
            />
          </div>
          <p v-if="v$.form.email.$error" class="mt-2 text-sm text-red-600 dark:text-red-400">
            {{ v$.form.email.$errors[0].$message }}
          </p>
        </div>

        <!-- 个人简介 -->
        <div class="sm:col-span-6">
          <label for="bio" class="block text-sm font-medium text-gray-700 dark:text-gray-300">个人简介</label>
          <div class="mt-1">
            <textarea
              id="bio"
              v-model="form.bio"
              rows="3"
              class="shadow-sm focus:ring-primary-500 focus:border-primary-500 block w-full sm:text-sm border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white rounded-md"
            />
          </div>
          <p class="mt-2 text-sm text-gray-500 dark:text-gray-400">简短介绍一下您自己</p>
        </div>
      </div>

      <div class="pt-5 flex justify-end">
        <button
          type="submit"
          :disabled="isSubmitting || !v$.form.$dirty || v$.form.$invalid"
          class="ml-3 inline-flex justify-center py-2 px-4 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          <span v-if="isSubmitting">保存中...</span>
          <span v-else>保存</span>
        </button>
      </div>
    </form>

    <div class="mt-10 pt-10 border-t border-gray-200 dark:border-gray-700">
      <h2 class="text-lg font-medium text-gray-900 dark:text-gray-100">修改密码</h2>
      <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
        确保您的账户使用强密码以保障安全
      </p>

      <form @submit.prevent="changePassword" class="mt-6 space-y-6">
        <div class="grid grid-cols-1 gap-y-6 gap-x-4 sm:grid-cols-6">
          <!-- 当前密码 -->
          <div class="sm:col-span-4">
            <label for="current-password" class="block text-sm font-medium text-gray-700 dark:text-gray-300">当前密码</label>
            <div class="mt-1 relative">
              <input
                id="current-password"
                v-model="passwordForm.currentPassword"
                :type="showCurrentPassword ? 'text' : 'password'"
                autocomplete="current-password"
                class="shadow-sm focus:ring-primary-500 focus:border-primary-500 block w-full sm:text-sm border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white rounded-md pr-10"
                :class="{ 'border-red-300 dark:border-red-700': v$.passwordForm.currentPassword.$error }"
              />
              <button
                type="button"
                @click="showCurrentPassword = !showCurrentPassword"
                class="absolute inset-y-0 right-0 pr-3 flex items-center text-gray-400 hover:text-gray-500"
              >
                <span v-if="showCurrentPassword" class="text-xs">隐藏</span>
                <span v-else class="text-xs">显示</span>
              </button>
            </div>
            <p v-if="v$.passwordForm.currentPassword.$error" class="mt-2 text-sm text-red-600 dark:text-red-400">
              {{ v$.passwordForm.currentPassword.$errors[0].$message }}
            </p>
          </div>

          <!-- 新密码 -->
          <div class="sm:col-span-4">
            <label for="new-password" class="block text-sm font-medium text-gray-700 dark:text-gray-300">新密码</label>
            <div class="mt-1 relative">
              <input
                id="new-password"
                v-model="passwordForm.newPassword"
                :type="showNewPassword ? 'text' : 'password'"
                autocomplete="new-password"
                class="shadow-sm focus:ring-primary-500 focus:border-primary-500 block w-full sm:text-sm border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white rounded-md pr-10"
                :class="{ 'border-red-300 dark:border-red-700': v$.passwordForm.newPassword.$error }"
              />
              <button
                type="button"
                @click="showNewPassword = !showNewPassword"
                class="absolute inset-y-0 right-0 pr-3 flex items-center text-gray-400 hover:text-gray-500"
              >
                <span v-if="showNewPassword" class="text-xs">隐藏</span>
                <span v-else class="text-xs">显示</span>
              </button>
            </div>
            <p v-if="v$.passwordForm.newPassword.$error" class="mt-2 text-sm text-red-600 dark:text-red-400">
              {{ v$.passwordForm.newPassword.$errors[0].$message }}
            </p>
          </div>

          <!-- 确认新密码 -->
          <div class="sm:col-span-4">
            <label for="confirm-password" class="block text-sm font-medium text-gray-700 dark:text-gray-300">确认新密码</label>
            <div class="mt-1 relative">
              <input
                id="confirm-password"
                v-model="passwordForm.confirmPassword"
                :type="showConfirmPassword ? 'text' : 'password'"
                autocomplete="new-password"
                class="shadow-sm focus:ring-primary-500 focus:border-primary-500 block w-full sm:text-sm border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white rounded-md pr-10"
                :class="{ 'border-red-300 dark:border-red-700': v$.passwordForm.confirmPassword.$error }"
              />
              <button
                type="button"
                @click="showConfirmPassword = !showConfirmPassword"
                class="absolute inset-y-0 right-0 pr-3 flex items-center text-gray-400 hover:text-gray-500"
              >
                <span v-if="showConfirmPassword" class="text-xs">隐藏</span>
                <span v-else class="text-xs">显示</span>
              </button>
            </div>
            <p v-if="v$.passwordForm.confirmPassword.$error" class="mt-2 text-sm text-red-600 dark:text-red-400">
              {{ v$.passwordForm.confirmPassword.$errors[0].$message }}
            </p>
          </div>
        </div>

        <div class="pt-5 flex justify-end">
          <button
            type="submit"
            :disabled="isChangingPassword || !v$.passwordForm.$dirty || v$.passwordForm.$invalid"
            class="ml-3 inline-flex justify-center py-2 px-4 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <span v-if="isChangingPassword">更新中...</span>
            <span v-else>更新密码</span>
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue';
import { useToast } from 'vue-toastification';
import { useVuelidate } from '@vuelidate/core';
import { required, email, minLength, sameAs, helpers } from '@vuelidate/validators';
import { authApi } from '../../api';
import { useUserStore } from '../../stores/user';
import { createImagePreview, revokeImagePreview, validateFileSize, isImageFile, compressImage } from '../../utils/fileUpload';

const toast = useToast();
const userStore = useUserStore();

// 表单状态
const form = reactive({
  username: '',
  nickname: '',
  email: '',
  bio: ''
});

const passwordForm = reactive({
  currentPassword: '',
  newPassword: '',
  confirmPassword: ''
});

// 表单验证规则
const rules = {
  form: {
    nickname: { required: helpers.withMessage('昵称不能为空', required) },
    email: { 
      required: helpers.withMessage('邮箱不能为空', required),
      email: helpers.withMessage('请输入有效的邮箱地址', email)
    }
  },
  passwordForm: {
    currentPassword: { 
      required: helpers.withMessage('当前密码不能为空', required) 
    },
    newPassword: { 
      required: helpers.withMessage('新密码不能为空', required),
      minLength: helpers.withMessage('密码长度至少为8个字符', minLength(8))
    },
    confirmPassword: { 
      required: helpers.withMessage('请确认新密码', required),
      sameAsPassword: helpers.withMessage('两次输入的密码不一致', sameAs(computed(() => passwordForm.newPassword)))
    }
  }
};

const v$ = useVuelidate(rules, { form, passwordForm });

// 提交状态
const isSubmitting = ref(false);
const isChangingPassword = ref(false);
const fileError = ref('');

// 密码显示状态
const showCurrentPassword = ref(false);
const showNewPassword = ref(false);
const showConfirmPassword = ref(false);

// 头像相关
const profileImage = ref('');
const avatarFile = ref(null);

// 用户首字母缩写
const userInitials = computed(() => {
  return form.nickname 
    ? form.nickname.charAt(0).toUpperCase() 
    : form.username 
      ? form.username.charAt(0).toUpperCase() 
      : 'U';
});

// 初始化表单数据
onMounted(async () => {
  // 如果用户未登录，则获取用户信息
  if (!userStore.user) {
    try {
      await userStore.fetchUserProfile();
    } catch (error) {
      toast.error('获取用户信息失败');
      return;
    }
  }

  // 填充表单数据
  if (userStore.user) {
    form.username = userStore.user.username || '';
    form.nickname = userStore.user.fullName || '';
    form.email = userStore.user.email || '';
    form.bio = userStore.user.bio || '';
    
    // 设置头像
    if (userStore.user.avatar) {
      profileImage.value = userStore.user.avatar;
    }
  }
});

// 更新个人资料
const updateProfile = async () => {
  // 验证表单
  const isFormValid = await v$.value.form.$validate();
  if (!isFormValid) return;
  
  isSubmitting.value = true;
  
  try {
    // 准备要提交的数据
    const profileData = new FormData();
    profileData.append('nickname', form.nickname);
    profileData.append('email', form.email);
    profileData.append('bio', form.bio || '');
    
    // 如果有新头像，添加到表单数据中
    if (avatarFile.value) {
      profileData.append('avatar', avatarFile.value);
    }
    
    // 调用API更新个人资料
    await userStore.updateProfile(profileData);
    
    toast.success('个人资料已更新');
    
    // 清除头像文件引用
    avatarFile.value = null;
  } catch (error) {
    toast.error(error.message || '更新个人资料失败');
  } finally {
    isSubmitting.value = false;
  }
};

// 修改密码
const changePassword = async () => {
  // 验证表单
  const isPasswordFormValid = await v$.value.passwordForm.$validate();
  if (!isPasswordFormValid) return;
  
  isChangingPassword.value = true;
  
  try {
    // 调用API修改密码
    await userStore.changePassword({
      currentPassword: passwordForm.currentPassword,
      newPassword: passwordForm.newPassword
    });
    
    toast.success('密码已更新');
    
    // 重置表单
    passwordForm.currentPassword = '';
    passwordForm.newPassword = '';
    passwordForm.confirmPassword = '';
    v$.value.passwordForm.$reset();
  } catch (error) {
    toast.error(error.message || '修改密码失败');
  } finally {
    isChangingPassword.value = false;
  }
};

// 处理头像上传
const handleFileChange = async (event) => {
  const file = event.target.files[0];
  if (!file) return;
  
  // 清除之前的错误
  fileError.value = '';
  
  // 验证文件类型
  if (!isImageFile(file)) {
    fileError.value = '请上传图片文件（JPG、PNG、GIF等格式）';
    return;
  }
  
  // 验证文件大小
  if (!validateFileSize(file, 2)) {
    fileError.value = '图片大小不能超过2MB';
    return;
  }
  
  try {
    // 压缩图片
    const compressedImage = await compressImage(file);
    
    // 清除之前的预览
    if (profileImage.value && !profileImage.value.startsWith('http')) {
      revokeImagePreview(profileImage.value);
    }
    
    // 创建新的预览
    profileImage.value = createImagePreview(compressedImage);
    
    // 保存文件引用，用于后续上传
    avatarFile.value = new File([compressedImage], file.name, { type: file.type });
  } catch (error) {
    fileError.value = '处理图片时出错';
    console.error('图片处理错误:', error);
  }
};
</script> 