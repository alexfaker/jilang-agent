<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import AppLayout from './components/layout/AppLayout.vue'

const route = useRoute()

// 判断是否是需要认证的页面
const isAuthPage = computed(() => {
  return route.name === 'Login' || route.name === 'Register' || route.meta.requiresAuth === false
})
</script>

<template>
  <router-view v-slot="{ Component }">
    <template v-if="isAuthPage">
      <component :is="Component" />
    </template>
    <template v-else>
      <AppLayout>
        <component :is="Component" />
      </AppLayout>
    </template>
  </router-view>
</template>

<style scoped>
.logo {
  height: 6em;
  padding: 1.5em;
  will-change: filter;
  transition: filter 300ms;
}
.logo:hover {
  filter: drop-shadow(0 0 2em #646cffaa);
}
.logo.vue:hover {
  filter: drop-shadow(0 0 2em #42b883aa);
}
</style>
