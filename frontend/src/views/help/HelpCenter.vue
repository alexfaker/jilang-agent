<template>
  <div class="flex-1 overflow-y-auto p-6">
    <!-- 页面标题 -->
    <div class="mb-6">
      <h1 class="text-2xl font-bold text-gray-900">帮助中心</h1>
      <p class="text-gray-600 mt-1">在这里找到您需要的帮助和支持</p>
    </div>

    <!-- 搜索框 -->
    <div class="mb-8">
      <div class="relative max-w-md">
        <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
          <i class="fas fa-search text-gray-400"></i>
        </div>
        <input
          type="text"
          v-model="searchQuery"
          placeholder="搜索帮助文档..."
          class="block w-full pl-10 pr-3 py-2 border border-gray-300 rounded-lg focus:ring-indigo-500 focus:border-indigo-500"
        >
      </div>
    </div>

    <!-- 快速帮助卡片 -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-8">
      <div
        v-for="category in helpCategories"
        :key="category.id"
        @click="activeCategory = category.id"
        class="bg-white rounded-lg shadow-sm border border-gray-200 p-6 cursor-pointer hover:shadow-md transition-shadow"
      >
        <div class="flex items-center mb-3">
          <i :class="category.icon" class="text-2xl text-indigo-600 mr-3"></i>
          <h3 class="text-lg font-semibold text-gray-900">{{ category.title }}</h3>
        </div>
        <p class="text-gray-600 text-sm">{{ category.description }}</p>
        <div class="mt-4">
          <span class="text-indigo-600 text-sm font-medium">查看详情 →</span>
        </div>
      </div>
    </div>

    <!-- 详细内容区域 -->
    <div class="bg-white rounded-lg shadow-sm border border-gray-200">
      <!-- 标签导航 -->
      <div class="border-b">
        <nav class="flex space-x-8 px-6">
          <button
            v-for="category in helpCategories"
            :key="category.id"
            @click="activeCategory = category.id"
            :class="[
              activeCategory === category.id
                ? 'py-4 px-1 border-b-2 border-indigo-600 text-indigo-600 font-medium text-sm'
                : 'py-4 px-1 text-gray-500 hover:text-indigo-600 font-medium text-sm'
            ]"
          >
            {{ category.title }}
          </button>
        </nav>
      </div>

      <!-- 内容区域 -->
      <div class="p-6">
        <!-- 快速开始 -->
        <div v-if="activeCategory === 'quickstart'">
          <h2 class="text-xl font-semibold text-gray-900 mb-4">快速开始</h2>
          <div class="space-y-4">
            <div v-for="item in quickstartItems" :key="item.id" class="border-l-4 border-indigo-500 pl-4">
              <h3 class="font-medium text-gray-900 mb-2">{{ item.title }}</h3>
              <p class="text-gray-600 text-sm">{{ item.description }}</p>
            </div>
          </div>
        </div>

        <!-- 常见问题 -->
        <div v-if="activeCategory === 'faq'">
          <h2 class="text-xl font-semibold text-gray-900 mb-4">常见问题</h2>
          <div class="space-y-4">
            <div
              v-for="item in filteredFaqItems"
              :key="item.id"
              class="border border-gray-200 rounded-lg"
            >
              <button
                @click="toggleFaq(item.id)"
                class="w-full px-4 py-3 text-left flex items-center justify-between hover:bg-gray-50"
              >
                <span class="font-medium text-gray-900">{{ item.question }}</span>
                <i :class="expandedFaq === item.id ? 'fas fa-chevron-up' : 'fas fa-chevron-down'" class="text-gray-400"></i>
              </button>
              <div v-if="expandedFaq === item.id" class="px-4 pb-3">
                <p class="text-gray-600">{{ item.answer }}</p>
              </div>
            </div>
          </div>
        </div>

        <!-- 使用指南 -->
        <div v-if="activeCategory === 'guides'">
          <h2 class="text-xl font-semibold text-gray-900 mb-4">使用指南</h2>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div v-for="guide in guideItems" :key="guide.id" class="border border-gray-200 rounded-lg p-4">
              <h3 class="font-medium text-gray-900 mb-2">{{ guide.title }}</h3>
              <p class="text-gray-600 text-sm mb-3">{{ guide.description }}</p>
              <button class="text-indigo-600 text-sm font-medium hover:text-indigo-500">
                阅读指南 →
              </button>
            </div>
          </div>
        </div>

        <!-- API文档 -->
        <div v-if="activeCategory === 'api'">
          <h2 class="text-xl font-semibold text-gray-900 mb-4">API文档</h2>
          <div class="space-y-6">
            <div class="bg-gray-50 rounded-lg p-4">
              <h3 class="font-medium text-gray-900 mb-2">API概述</h3>
              <p class="text-gray-600 text-sm mb-3">
                我们的API允许您以编程方式管理工作流、执行任务和获取数据。
              </p>
              <div class="flex space-x-4">
                <button class="text-indigo-600 text-sm font-medium hover:text-indigo-500">
                  查看文档 →
                </button>
                <button class="text-indigo-600 text-sm font-medium hover:text-indigo-500">
                  API参考 →
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- 联系支持 -->
        <div v-if="activeCategory === 'contact'">
          <h2 class="text-xl font-semibold text-gray-900 mb-4">联系支持</h2>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div class="space-y-4">
              <div class="border border-gray-200 rounded-lg p-4">
                <h3 class="font-medium text-gray-900 mb-2">在线客服</h3>
                <p class="text-gray-600 text-sm mb-3">工作日 9:00-18:00 为您提供实时帮助</p>
                <button class="bg-indigo-600 text-white px-4 py-2 rounded-lg text-sm hover:bg-indigo-700">
                  开始对话
                </button>
              </div>
              <div class="border border-gray-200 rounded-lg p-4">
                <h3 class="font-medium text-gray-900 mb-2">邮件支持</h3>
                <p class="text-gray-600 text-sm mb-3">发送邮件至 support@example.com</p>
                <p class="text-gray-500 text-xs">通常在24小时内回复</p>
              </div>
            </div>
            <div class="border border-gray-200 rounded-lg p-4">
              <h3 class="font-medium text-gray-900 mb-4">提交工单</h3>
              <form @submit.prevent="submitTicket" class="space-y-3">
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">问题类型</label>
                  <select v-model="ticketForm.type" class="w-full border border-gray-300 rounded px-3 py-2 text-sm">
                    <option value="">请选择</option>
                    <option value="technical">技术问题</option>
                    <option value="billing">账单问题</option>
                    <option value="feature">功能建议</option>
                    <option value="other">其他</option>
                  </select>
                </div>
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">问题描述</label>
                  <textarea
                    v-model="ticketForm.description"
                    rows="4"
                    class="w-full border border-gray-300 rounded px-3 py-2 text-sm"
                    placeholder="请详细描述您遇到的问题..."
                  ></textarea>
                </div>
                <button
                  type="submit"
                  class="w-full bg-indigo-600 text-white px-4 py-2 rounded-lg text-sm hover:bg-indigo-700"
                >
                  提交工单
                </button>
              </form>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'

const searchQuery = ref('')
const activeCategory = ref('quickstart')
const expandedFaq = ref(null)

const ticketForm = ref({
  type: '',
  description: ''
})

const helpCategories = [
  {
    id: 'quickstart',
    title: '快速开始',
    description: '了解如何开始使用平台的基本功能',
    icon: 'fas fa-rocket'
  },
  {
    id: 'faq',
    title: '常见问题',
    description: '查找用户最常遇到的问题和解答',
    icon: 'fas fa-question-circle'
  },
  {
    id: 'guides',
    title: '使用指南',
    description: '详细的功能使用教程和最佳实践',
    icon: 'fas fa-book'
  },
  {
    id: 'api',
    title: 'API文档',
    description: '开发者集成指南和API参考文档',
    icon: 'fas fa-code'
  },
  {
    id: 'contact',
    title: '联系支持',
    description: '获取个性化帮助和技术支持',
    icon: 'fas fa-headset'
  }
]

const quickstartItems = [
  {
    id: 1,
    title: '创建您的第一个工作流',
    description: '学习如何从零开始创建和配置一个工作流'
  },
  {
    id: 2,
    title: '从市场获取模板',
    description: '浏览工作流市场，找到适合您需求的模板'
  },
  {
    id: 3,
    title: '执行和监控',
    description: '了解如何执行工作流并监控执行状态'
  },
  {
    id: 4,
    title: '查看统计数据',
    description: '使用统计分析页面了解工作流性能'
  }
]

const faqItems = [
  {
    id: 1,
    question: '如何创建一个新的工作流？',
    answer: '点击"我的工作流"页面右上角的"创建工作流"按钮，然后按照向导步骤配置您的工作流。您也可以从工作流市场选择一个模板开始。'
  },
  {
    id: 2,
    question: '工作流执行失败了怎么办？',
    answer: '请检查执行历史页面的错误信息。常见原因包括配置错误、输入数据格式不正确或API密钥过期。您可以在设置页面更新相关配置。'
  },
  {
    id: 3,
    question: '如何查看工作流的执行历史？',
    answer: '访问"执行历史"页面，您可以看到所有工作流的执行记录，包括状态、开始时间、执行时长等信息。'
  },
  {
    id: 4,
    question: '如何修改账户设置？',
    answer: '点击左侧菜单的"设置"，您可以修改个人信息、密码、主题偏好等设置。'
  },
  {
    id: 5,
    question: '如何从工作流市场获取模板？',
    answer: '访问"工作流市场"页面，浏览可用的模板。点击模板可以查看详情，然后选择"免费获取"或"购买"按钮将模板添加到您的账户。'
  }
]

const guideItems = [
  {
    id: 1,
    title: '工作流设计最佳实践',
    description: '学习如何设计高效、可维护的工作流'
  },
  {
    id: 2,
    title: '数据格式和类型',
    description: '了解支持的数据格式和如何正确配置输入输出'
  },
  {
    id: 3,
    title: '错误处理和调试',
    description: '掌握工作流调试技巧和错误处理策略'
  },
  {
    id: 4,
    title: '性能优化指南',
    description: '提升工作流执行效率的技巧和建议'
  }
]

const filteredFaqItems = computed(() => {
  if (!searchQuery.value) return faqItems
  
  return faqItems.filter(item =>
    item.question.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
    item.answer.toLowerCase().includes(searchQuery.value.toLowerCase())
  )
})

const toggleFaq = (id) => {
  expandedFaq.value = expandedFaq.value === id ? null : id
}

const submitTicket = async () => {
  if (!ticketForm.value.type || !ticketForm.value.description) {
    alert('请填写完整信息')
    return
  }
  
  // 这里应该调用API提交工单
  alert('工单提交成功！我们会尽快回复您。')
  
  // 重置表单
  ticketForm.value = {
    type: '',
    description: ''
  }
}
</script>

<style scoped>
/* 这里可以添加组件特定的样式 */
</style> 