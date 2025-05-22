<template>
  <div class="min-h-screen flex flex-col items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full text-center">
      <div class="text-6xl font-bold text-brand mb-4">404</div>
      <h2 class="text-2xl font-semibold text-gray-800 mb-2">页面未找到</h2>
      <p class="text-gray-600 mb-8">抱歉，您请求的页面不存在或已被移动。</p>
      <div class="flex flex-col sm:flex-row justify-center gap-4">
        <router-link to="/" class="btn btn-primary flex items-center justify-center">
          <i class="fas fa-arrow-left mr-2"></i>
          返回主页
        </router-link>
        <button @click="goBack" class="btn btn-secondary flex items-center justify-center">
          <i class="fas fa-undo mr-2"></i>
          返回上一页
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'

const router = useRouter()

const goBack = () => {
  router.back()
}
</script> 