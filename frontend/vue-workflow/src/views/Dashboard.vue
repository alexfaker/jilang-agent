<template>
  <div>
    <!-- 欢迎区域 -->
    <div class="bg-white rounded-lg p-6 shadow-card mb-6">
      <div class="flex items-center justify-between">
        <div>
          <h2 class="text-xl font-semibold text-gray-800">欢迎回来，{{ userData.name }}</h2>
          <p class="text-gray-600 mt-1">{{ currentTimeGreeting }}</p>
        </div>
        <div class="hidden md:block">
          <button class="btn btn-primary">
            <i class="fas fa-plus mr-2"></i>创建工作流
          </button>
        </div>
      </div>
    </div>

    <!-- 数据概览 -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-6">
      <!-- 工作流数量卡片 -->
      <div class="bg-white rounded-lg p-6 shadow-card">
        <div class="flex items-center justify-between mb-4">
          <h3 class="text-gray-500 text-sm font-medium">我的工作流</h3>
          <span class="w-10 h-10 rounded-full bg-blue-100 flex items-center justify-center">
            <i class="fas fa-project-diagram text-blue-600"></i>
          </span>
        </div>
        <div class="flex items-end">
          <p class="text-3xl font-semibold text-gray-800">{{ stats.workflowCount }}</p>
          <p class="text-sm text-green-600 ml-2 mb-1 flex items-center">
            <i class="fas fa-arrow-up mr-1"></i> {{ stats.workflowGrowth }}%
          </p>
        </div>
        <p class="text-gray-600 text-sm mt-1">相比上月</p>
      </div>

      <!-- 执行次数卡片 -->
      <div class="bg-white rounded-lg p-6 shadow-card">
        <div class="flex items-center justify-between mb-4">
          <h3 class="text-gray-500 text-sm font-medium">累计执行次数</h3>
          <span class="w-10 h-10 rounded-full bg-green-100 flex items-center justify-center">
            <i class="fas fa-play text-green-600"></i>
          </span>
        </div>
        <div class="flex items-end">
          <p class="text-3xl font-semibold text-gray-800">{{ stats.executionCount }}</p>
          <p class="text-sm text-green-600 ml-2 mb-1 flex items-center">
            <i class="fas fa-arrow-up mr-1"></i> {{ stats.executionGrowth }}%
          </p>
        </div>
        <p class="text-gray-600 text-sm mt-1">相比上月</p>
      </div>

      <!-- 成功率卡片 -->
      <div class="bg-white rounded-lg p-6 shadow-card">
        <div class="flex items-center justify-between mb-4">
          <h3 class="text-gray-500 text-sm font-medium">执行成功率</h3>
          <span class="w-10 h-10 rounded-full bg-amber-100 flex items-center justify-center">
            <i class="fas fa-check-circle text-amber-600"></i>
          </span>
        </div>
        <div class="flex items-end">
          <p class="text-3xl font-semibold text-gray-800">{{ stats.successRate }}%</p>
          <p class="text-sm text-red-600 ml-2 mb-1 flex items-center">
            <i class="fas fa-arrow-down mr-1"></i> {{ stats.successRateDiff }}%
          </p>
        </div>
        <p class="text-gray-600 text-sm mt-1">相比上月</p>
      </div>

      <!-- 使用时长卡片 -->
      <div class="bg-white rounded-lg p-6 shadow-card">
        <div class="flex items-center justify-between mb-4">
          <h3 class="text-gray-500 text-sm font-medium">本月使用时长</h3>
          <span class="w-10 h-10 rounded-full bg-indigo-100 flex items-center justify-center">
            <i class="fas fa-clock text-indigo-600"></i>
          </span>
        </div>
        <div class="flex items-end">
          <p class="text-3xl font-semibold text-gray-800">{{ stats.usageHours }}</p>
          <span class="text-sm text-gray-600 ml-1 mb-1">小时</span>
        </div>
        <p class="text-gray-600 text-sm mt-1">剩余 {{ stats.remainingHours }} 小时</p>
      </div>
    </div>

    <!-- 最近工作流和图表区域 -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- 最近工作流 -->
      <div class="lg:col-span-1 bg-white rounded-lg shadow-card overflow-hidden">
        <div class="p-6 border-b border-gray-100 flex justify-between items-center">
          <h3 class="font-semibold text-gray-800">最近执行的工作流</h3>
          <router-link to="/executions" class="text-sm text-brand hover:underline">
            查看全部
          </router-link>
        </div>
        <div class="divide-y divide-gray-100">
          <div v-for="(workflow, index) in recentWorkflows" :key="index" class="p-4 hover:bg-gray-50">
            <div class="flex items-center justify-between mb-2">
              <div class="flex items-center">
                <div class="w-8 h-8 rounded-full bg-gray-100 flex items-center justify-center mr-3">
                  <i class="fas fa-robot text-gray-600"></i>
                </div>
                <router-link :to="`/workflows/${workflow.id}`" class="font-medium text-gray-800 hover:text-brand">
                  {{ workflow.name }}
                </router-link>
              </div>
              <span :class="`badge badge-${workflow.status}`">{{ workflow.statusText }}</span>
            </div>
            <div class="flex items-center text-sm text-gray-500">
              <i class="fas fa-calendar-alt mr-2"></i>
              <span>{{ workflow.executedAt }}</span>
              <span class="mx-2">•</span>
              <i class="fas fa-clock mr-2"></i>
              <span>{{ workflow.duration }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 执行统计图表 -->
      <div class="lg:col-span-2 bg-white rounded-lg shadow-card overflow-hidden">
        <div class="p-6 border-b border-gray-100 flex justify-between items-center">
          <h3 class="font-semibold text-gray-800">工作流执行统计</h3>
          <div class="flex space-x-2">
            <button 
              v-for="(period, index) in periods" 
              :key="index" 
              @click="selectedPeriod = period.value" 
              class="px-3 py-1 text-sm rounded-md" 
              :class="selectedPeriod === period.value ? 'bg-indigo-100 text-brand' : 'text-gray-600 hover:bg-gray-100'"
            >
              {{ period.label }}
            </button>
          </div>
        </div>
        <div class="p-6">
          <div ref="chartContainer" class="w-full h-64"></div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import * as echarts from 'echarts'

// 用户数据（实际项目中应该从状态管理中获取）
const userData = ref({
  name: '张三',
  avatar: ''
})

// 统计数据
const stats = ref({
  workflowCount: 12,
  workflowGrowth: 8.5,
  executionCount: 125,
  executionGrowth: 12.3,
  successRate: 95.7,
  successRateDiff: 1.2,
  usageHours: 42.5,
  remainingHours: 57.5
})

// 最近的工作流
const recentWorkflows = ref([
  {
    id: 1,
    name: '每日数据采集',
    status: 'success',
    statusText: '成功',
    executedAt: '今天 09:30',
    duration: '2分钟'
  },
  {
    id: 2,
    name: '客户数据分析',
    status: 'error',
    statusText: '失败',
    executedAt: '昨天 15:45',
    duration: '5分钟'
  },
  {
    id: 3,
    name: '周报自动生成',
    status: 'success',
    statusText: '成功',
    executedAt: '2天前 08:00',
    duration: '3分钟'
  },
  {
    id: 4,
    name: '社交媒体内容创建',
    status: 'warning',
    statusText: '部分成功',
    executedAt: '3天前 14:20',
    duration: '7分钟'
  }
])

// 图表容器ref
const chartContainer = ref<HTMLElement | null>(null)
// 图表实例
let chartInstance: echarts.ECharts | null = null

// 时间段选择
const periods = [
  { label: '7天', value: 7 },
  { label: '30天', value: 30 },
  { label: '90天', value: 90 }
]
const selectedPeriod = ref(7)

// 当前时间问候语
const currentTimeGreeting = computed(() => {
  const hour = new Date().getHours()
  if (hour < 6) return '夜深了，注意休息'
  if (hour < 12) return '早上好，开始新的一天'
  if (hour < 18) return '下午好，今天过得如何'
  return '晚上好，工作有条不紊吗'
})

// 初始化图表
const initChart = () => {
  if (!chartContainer.value) return
  
  chartInstance = echarts.init(chartContainer.value)
  updateChart()
  
  window.addEventListener('resize', () => {
    chartInstance?.resize()
  })
}

// 更新图表数据
const updateChart = () => {
  // 这里应该根据selectedPeriod从API获取数据
  // 暂时使用模拟数据
  const days = selectedPeriod.value
  const dates = Array.from({ length: days }, (_, i) => {
    const date = new Date()
    date.setDate(date.getDate() - (days - i - 1))
    return `${date.getMonth() + 1}/${date.getDate()}`
  })
  
  const successData = Array.from({ length: days }, () => Math.floor(Math.random() * 10) + 5)
  const failData = Array.from({ length: days }, () => Math.floor(Math.random() * 3))
  
  const option = {
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow'
      }
    },
    legend: {
      data: ['成功', '失败'],
      right: 10,
      top: 0
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      data: dates,
      axisLine: {
        lineStyle: {
          color: '#E5E7EB'
        }
      },
      axisLabel: {
        color: '#6B7280'
      }
    },
    yAxis: {
      type: 'value',
      splitLine: {
        lineStyle: {
          color: '#E5E7EB'
        }
      },
      axisLabel: {
        color: '#6B7280'
      }
    },
    series: [
      {
        name: '成功',
        type: 'bar',
        stack: 'total',
        itemStyle: {
          color: '#10B981'
        },
        data: successData
      },
      {
        name: '失败',
        type: 'bar',
        stack: 'total',
        itemStyle: {
          color: '#EF4444'
        },
        data: failData
      }
    ]
  }
  
  chartInstance?.setOption(option)
}

onMounted(() => {
  initChart()
})
</script> 