import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { createRouter, createWebHistory } from 'vue-router'
import App from './App.vue'

import './style.css'
import '@fortawesome/fontawesome-free/css/all.min.css'

// 导入路由
import routes from './router'

// 创建路由实例
const router = createRouter({
  history: createWebHistory(),
  routes
})

// 创建Pinia状态管理实例
const pinia = createPinia()

// 创建应用实例
const app = createApp(App)

// 使用插件
app.use(pinia)
app.use(router)

// 挂载应用
app.mount('#app')
