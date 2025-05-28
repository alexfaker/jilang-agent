<template>
  <div class="min-h-screen bg-gray-50">
    <!-- 顶部导航栏 -->
    <nav class="bg-white shadow-sm sticky top-0 z-50">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-16">
          <div class="flex items-center">
            <div class="flex-shrink-0 flex items-center">
              <CpuChipIcon class="h-8 w-8 text-indigo-600 mr-2" />
              <span class="text-xl font-bold text-gray-800">JiLang Agent</span>
            </div>
          </div>
          <div class="flex items-center space-x-4">
            <router-link 
              to="/auth/login" 
              class="text-gray-500 hover:text-gray-700 px-3 py-2 rounded-md text-sm font-medium transition-colors"
            >
              登录
            </router-link>
            <router-link 
              to="/auth/register" 
              class="px-4 py-2 rounded-md text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 transition-colors"
            >
              注册
            </router-link>
          </div>
        </div>
      </div>
    </nav>

    <!-- 欢迎横幅 -->
    <div class="bg-indigo-700 text-white">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12 md:py-20">
        <div class="md:flex md:items-center md:justify-between">
          <div class="md:w-3/5">
            <h1 class="text-3xl md:text-4xl font-bold mb-4">构建、部署和管理智能AI工作流</h1>
            <p class="text-indigo-100 text-lg mb-6">创建自定义工作流，轻松整合AI能力到您的业务中，提升效率，简化流程</p>
            <div class="flex flex-wrap gap-4">
              <router-link 
                to="/auth/register" 
                class="px-6 py-3 bg-white text-indigo-700 font-medium rounded-lg hover:bg-indigo-50 transition-colors"
              >
                开始使用
              </router-link>
              <a 
                href="#features" 
                class="px-6 py-3 bg-transparent border border-white text-white font-medium rounded-lg hover:bg-indigo-600 transition-colors"
                @click="scrollToFeatures"
              >
                了解更多
              </a>
            </div>
          </div>
          <div class="mt-8 md:mt-0 md:w-2/5">
            <img 
              src="https://images.unsplash.com/photo-1485827404703-89b55fcc595e?q=80&w=500&auto=format&fit=crop" 
              alt="AI机器人" 
              class="w-full rounded-lg shadow-xl"
            >
          </div>
        </div>
      </div>
    </div>

    <!-- 主要功能模块 -->
    <div id="features" class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <h2 class="text-2xl font-bold text-gray-900 text-center mb-12">主要功能</h2>
      
      <!-- 第一行：工作流和工作流市场 -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-8 mb-12">
        <!-- 我的工作流 -->
        <div class="module-card bg-white p-6 overflow-hidden rounded-2xl shadow-lg hover:shadow-xl transition-all duration-300 hover:-translate-y-1">
          <div class="flex items-center justify-between mb-4">
            <h3 class="text-xl font-semibold text-gray-800">我的工作流</h3>
            <div class="w-12 h-12 rounded-full bg-blue-100 flex items-center justify-center">
              <RectangleGroupIcon class="h-7 w-7 text-blue-600" />
            </div>
          </div>
          <p class="text-gray-600 mb-5">创建和管理您的自定义工作流，连接多个组件完成复杂任务</p>
          <img 
            src="https://images.unsplash.com/photo-1551288049-bebda4e38f71?q=80&w=500&auto=format&fit=crop" 
            alt="工作流示例" 
            class="w-full h-48 object-cover rounded-lg mb-5"
          >
          <router-link 
            to="/auth/login" 
            class="inline-block px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
          >
            查看我的工作流
          </router-link>
        </div>
        
        <!-- 工作流市场 -->
        <div class="module-card bg-white p-6 overflow-hidden rounded-2xl shadow-lg hover:shadow-xl transition-all duration-300 hover:-translate-y-1">
          <div class="flex items-center justify-between mb-4">
            <h3 class="text-xl font-semibold text-gray-800">工作流市场</h3>
            <div class="w-12 h-12 rounded-full bg-purple-100 flex items-center justify-center">
              <BuildingStorefrontIcon class="h-7 w-7 text-purple-600" />
            </div>
          </div>
          <p class="text-gray-600 mb-5">浏览和使用预构建的工作流，立即添加智能功能到您的应用中</p>
          <img 
            src="https://images.unsplash.com/photo-1655720031554-a929595ffad7?q=80&w=500&auto=format&fit=crop" 
            alt="工作流市场" 
            class="w-full h-48 object-cover rounded-lg mb-5"
          >
          <router-link 
            to="/auth/login" 
            class="inline-block px-4 py-2 bg-purple-600 text-white rounded-lg hover:bg-purple-700 transition-colors"
          >
            浏览工作流
          </router-link>
        </div>
      </div>
      
      <!-- 第二行：执行历史和统计分析 -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-8 mb-12">
        <!-- 执行历史 -->
        <div class="module-card bg-white p-6 overflow-hidden rounded-2xl shadow-lg hover:shadow-xl transition-all duration-300 hover:-translate-y-1">
          <div class="flex items-center justify-between mb-4">
            <h3 class="text-xl font-semibold text-gray-800">执行历史</h3>
            <div class="w-12 h-12 rounded-full bg-green-100 flex items-center justify-center">
              <ClockIcon class="h-7 w-7 text-green-600" />
            </div>
          </div>
          <p class="text-gray-600 mb-5">查看您的工作流执行历史记录，了解处理结果和性能</p>
          <img 
            src="https://images.unsplash.com/photo-1551288049-bebda4e38f71?q=80&w=500&auto=format&fit=crop" 
            alt="执行历史" 
            class="w-full h-48 object-cover rounded-lg mb-5"
          >
          <router-link 
            to="/auth/login" 
            class="inline-block px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors"
          >
            查看执行历史
          </router-link>
        </div>
        
        <!-- 统计分析 -->
        <div class="module-card bg-white p-6 overflow-hidden rounded-2xl shadow-lg hover:shadow-xl transition-all duration-300 hover:-translate-y-1">
          <div class="flex items-center justify-between mb-4">
            <h3 class="text-xl font-semibold text-gray-800">统计分析</h3>
            <div class="w-12 h-12 rounded-full bg-amber-100 flex items-center justify-center">
              <ChartBarIcon class="h-7 w-7 text-amber-600" />
            </div>
          </div>
          <p class="text-gray-600 mb-5">获取工作流执行的性能指标和统计数据，优化您的工作流</p>
          <img 
            src="https://images.unsplash.com/photo-1551288049-bebda4e38f71?q=80&w=500&auto=format&fit=crop" 
            alt="统计分析" 
            class="w-full h-48 object-cover rounded-lg mb-5"
          >
          <router-link 
            to="/auth/login" 
            class="inline-block px-4 py-2 bg-amber-600 text-white rounded-lg hover:bg-amber-700 transition-colors"
          >
            查看统计数据
          </router-link>
        </div>
      </div>
      
      <!-- 第三行：个人中心 -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
        <div class="module-card bg-white p-6 overflow-hidden rounded-2xl shadow-lg hover:shadow-xl transition-all duration-300 hover:-translate-y-1 md:col-span-2 md:max-w-lg md:mx-auto">
          <div class="flex items-center justify-between mb-4">
            <h3 class="text-xl font-semibold text-gray-800">个人中心</h3>
            <div class="w-12 h-12 rounded-full bg-indigo-100 flex items-center justify-center">
              <UserIcon class="h-7 w-7 text-indigo-600" />
            </div>
          </div>
          <p class="text-gray-600 mb-5">管理您的账户设置、API密钥和个人偏好</p>
          <img 
            src="https://images.unsplash.com/photo-1551288049-bebda4e38f71?q=80&w=500&auto=format&fit=crop" 
            alt="个人中心" 
            class="w-full h-48 object-cover rounded-lg mb-5"
          >
          <router-link 
            to="/auth/login" 
            class="inline-block px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-colors"
          >
            查看个人设置
          </router-link>
        </div>
      </div>
    </div>

    <!-- 底部信息 -->
    <footer class="bg-gray-800 text-white py-12">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="grid grid-cols-1 md:grid-cols-4 gap-8">
          <div>
            <h3 class="text-lg font-semibold mb-4">JiLang Agent</h3>
            <p class="text-gray-400">智能工作流平台，让您轻松创建和管理工作流程。</p>
          </div>
          <div>
            <h3 class="text-lg font-semibold mb-4">功能</h3>
            <ul class="space-y-2 text-gray-400">
              <li><router-link to="/auth/login" class="hover:text-white transition-colors">工作流</router-link></li>
              <li><router-link to="/auth/login" class="hover:text-white transition-colors">工作流市场</router-link></li>
              <li><router-link to="/auth/login" class="hover:text-white transition-colors">统计分析</router-link></li>
            </ul>
          </div>
          <div>
            <h3 class="text-lg font-semibold mb-4">资源</h3>
            <ul class="space-y-2 text-gray-400">
              <li><a href="#" class="hover:text-white transition-colors">文档</a></li>
              <li><a href="#" class="hover:text-white transition-colors">API参考</a></li>
              <li><a href="#" class="hover:text-white transition-colors">示例</a></li>
              <li><a href="#" class="hover:text-white transition-colors">教程</a></li>
            </ul>
          </div>
          <div>
            <h3 class="text-lg font-semibold mb-4">联系我们</h3>
            <ul class="space-y-2 text-gray-400">
              <li><a href="#" class="hover:text-white transition-colors">关于我们</a></li>
              <li><a href="#" class="hover:text-white transition-colors">联系方式</a></li>
              <li><a href="#" class="hover:text-white transition-colors">帮助中心</a></li>
            </ul>
          </div>
        </div>
        <div class="border-t border-gray-700 mt-8 pt-8 flex flex-col md:flex-row justify-between items-center">
          <p class="text-gray-400">© {{ new Date().getFullYear() }} JiLang Agent. 保留所有权利。</p>
          <div class="flex space-x-6 mt-4 md:mt-0">
            <a href="#" class="text-gray-400 hover:text-white transition-colors">
              <LinkIcon class="h-5 w-5" />
            </a>
            <a href="#" class="text-gray-400 hover:text-white transition-colors">
              <LinkIcon class="h-5 w-5" />
            </a>
            <a href="#" class="text-gray-400 hover:text-white transition-colors">
              <LinkIcon class="h-5 w-5" />
            </a>
          </div>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { 
  CpuChipIcon, 
  RectangleGroupIcon, 
  BuildingStorefrontIcon, 
  ClockIcon, 
  ChartBarIcon, 
  UserIcon, 
  LinkIcon 
} from '@heroicons/vue/24/outline'

// 平滑滚动到功能区域
const scrollToFeatures = (event) => {
  event.preventDefault()
  const featuresElement = document.getElementById('features')
  if (featuresElement) {
    featuresElement.scrollIntoView({
      behavior: 'smooth'
    })
  }
}
</script>

<style scoped>
.module-card {
  transition: all 0.3s ease;
}

.module-card:hover {
  transform: translateY(-5px);
}
</style> 