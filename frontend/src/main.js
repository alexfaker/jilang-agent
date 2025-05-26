import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'
import router from './router'
import Toast from "vue-toastification"
import "vue-toastification/dist/index.css"
import App from './App.vue'
import { useSettingsStore } from './stores/settings'
import { useUserStore } from './stores/user'
import notify from './utils/notification'

// 创建应用
const app = createApp(App)

// 注册Pinia
const pinia = createPinia()
app.use(pinia)

// 初始化设置
const settingsStore = useSettingsStore(pinia)
await settingsStore.initSettings()

// 初始化用户状态
const userStore = useUserStore(pinia)
userStore.initializeFromLocalStorage()

// 注册路由
app.use(router)

// 注册Toast通知
app.use(Toast, {
  transition: "Vue-Toastification__fade",
  maxToasts: 5,
  newestOnTop: true,
  position: "top-right",
  timeout: 5000,
  closeOnClick: true,
  pauseOnFocusLoss: true,
  pauseOnHover: true,
  draggable: true,
  draggablePercent: 0.6,
  showCloseButtonOnHover: false,
  hideProgressBar: false,
  closeButton: "button",
  icon: true,
  rtl: false
})

// 将设置存储作为全局属性
app.config.globalProperties.$settings = settingsStore

// 将用户存储作为全局属性
app.config.globalProperties.$user = userStore

// 将通知工具作为全局属性
app.config.globalProperties.$notify = notify

// 挂载应用
app.mount('#app') 