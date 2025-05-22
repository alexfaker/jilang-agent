<template>
  <div class="dashboard-page">
    <!-- 欢迎区域 -->
    <div class="bg-white rounded-lg p-6 shadow-card mb-6">
      <div class="flex items-center justify-between">
        <div>
          <h2 class="text-xl font-semibold text-gray-800">欢迎回来，{{ userInfo.fullName || '用户' }}</h2>
          <p class="text-gray-600 mt-1">{{ currentTimeGreeting }}</p>
        </div>
        <div class="hidden md:block">
          <router-link to="/workflows/create" class="btn btn-primary">
            <i class="fas fa-plus mr-2"></i>创建工作流
          </router-link>
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
          <p v-if="stats.workflowGrowth > 0" class="text-sm text-green-600 ml-2 mb-1 flex items-center">
            <i class="fas fa-arrow-up mr-1"></i> {{ stats.workflowGrowth }}%
          </p>
          <p v-else-if="stats.workflowGrowth < 0" class="text-sm text-red-600 ml-2 mb-1 flex items-center">
            <i class="fas fa-arrow-down mr-1"></i> {{ Math.abs(stats.workflowGrowth) }}%
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
          <p v-if="stats.executionGrowth > 0" class="text-sm text-green-600 ml-2 mb-1 flex items-center">
            <i class="fas fa-arrow-up mr-1"></i> {{ stats.executionGrowth }}%
          </p>
          <p v-else-if="stats.executionGrowth < 0" class="text-sm text-red-600 ml-2 mb-1 flex items-center">
            <i class="fas fa-arrow-down mr-1"></i> {{ Math.abs(stats.executionGrowth) }}%
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
          <p v-if="stats.successRateDiff > 0" class="text-sm text-green-600 ml-2 mb-1 flex items-center">
            <i class="fas fa-arrow-up mr-1"></i> {{ stats.successRateDiff }}%
          </p>
          <p v-else-if="stats.successRateDiff < 0" class="text-sm text-red-600 ml-2 mb-1 flex items-center">
            <i class="fas fa-arrow-down mr-1"></i> {{ Math.abs(stats.successRateDiff) }}%
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
        
        <div v-if="loading" class="p-6 flex justify-center">
          <div class="spinner"></div>
        </div>
        
        <div v-else-if="error" class="p-6 text-center">
          <p class="text-red-500">{{ error }}</p>
          <button @click="fetchDashboardData" class="mt-2 text-sm text-brand hover:underline">
            重试
          </button>
        </div>
        
        <div v-else-if="recentExecutions.length === 0" class="p-6 text-center text-gray-500">
          <i class="fas fa-history text-2xl mb-2"></i>
          <p>暂无执行记录</p>
        </div>
        
        <div v-else class="divide-y divide-gray-100">
          <div v-for="(execution, index) in recentExecutions" :key="index" class="p-4 hover:bg-gray-50">
            <div class="flex items-center justify-between mb-2">
              <div class="flex items-center">
                <div class="w-8 h-8 rounded-full bg-gray-100 flex items-center justify-center mr-3">
                  <i class="fas fa-robot text-gray-600"></i>
                </div>
                <router-link :to="`/workflows/${execution.workflowId}`" class="font-medium text-gray-800 hover:text-brand">
                  {{ execution.workflowName }}
                </router-link>
              </div>
              <span :class="`badge badge-${execution.statusClass}`">{{ execution.statusText }}</span>
            </div>
            <div class="flex items-center text-sm text-gray-500">
              <i class="fas fa-calendar-alt mr-2"></i>
              <span>{{ formatDate(execution.startedAt) }}</span>
              <span class="mx-2">•</span>
              <i class="fas fa-clock mr-2"></i>
              <span>{{ execution.duration }}秒</span>
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
        
        <div v-if="loading" class="p-6 flex justify-center">
          <div class="spinner"></div>
        </div>
        
        <div v-else-if="error" class="p-6 text-center">
          <p class="text-red-500">{{ error }}</p>
          <button @click="fetchDashboardData" class="mt-2 text-sm text-brand hover:underline">
            重试
          </button>
        </div>
        
        <div v-else class="p-6">
          <div ref="chartContainer" class="w-full h-64"></div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, computed, onBeforeUnmount } from 'vue';
import * as echarts from 'echarts';
import { statsApi } from '../api';

export default {
  name: 'DashboardView',
  setup() {
    const loading = ref(true);
    const error = ref(null);
    const userInfo = ref({
      fullName: localStorage.getItem('user') ? JSON.parse(localStorage.getItem('user')).fullName : '用户'
    });
    
    // 统计数据
    const stats = ref({
      workflowCount: 0,
      workflowGrowth: 0,
      executionCount: 0,
      executionGrowth: 0,
      successRate: 0,
      successRateDiff: 0,
      usageHours: 0,
      remainingHours: 0
    });
    
    // 最近的执行记录
    const recentExecutions = ref([]);
    
    // 图表容器ref
    const chartContainer = ref(null);
    // 图表实例
    let chartInstance = null;
    
    // 时间段选择
    const periods = [
      { label: '7天', value: 7 },
      { label: '30天', value: 30 },
      { label: '90天', value: 90 }
    ];
    const selectedPeriod = ref(7);
    
    // 当前时间问候语
    const currentTimeGreeting = computed(() => {
      const hour = new Date().getHours();
      if (hour < 6) return '夜深了，注意休息';
      if (hour < 12) return '早上好，开始新的一天';
      if (hour < 18) return '下午好，今天过得如何';
      return '晚上好，工作有条不紊吗';
    });
    
    // 获取仪表盘数据
    const fetchDashboardData = async () => {
      loading.value = true;
      error.value = null;
      
      try {
        // 实际项目中应该调用API获取数据
        // const response = await statsApi.getDashboardStats();
        
        // 模拟API响应数据
        setTimeout(() => {
          stats.value = {
            workflowCount: 12,
            workflowGrowth: 8.5,
            executionCount: 125,
            executionGrowth: 12.3,
            successRate: 95.7,
            successRateDiff: -1.2,
            usageHours: 42.5,
            remainingHours: 57.5
          };
          
          recentExecutions.value = [
            {
              id: 1,
              workflowId: 1,
              workflowName: '每日数据采集',
              status: 'success',
              statusText: '成功',
              statusClass: 'success',
              startedAt: new Date(Date.now() - 3600000).toISOString(),
              duration: 120
            },
            {
              id: 2,
              workflowId: 2,
              workflowName: '客户数据分析',
              status: 'failed',
              statusText: '失败',
              statusClass: 'error',
              startedAt: new Date(Date.now() - 86400000).toISOString(),
              duration: 300
            },
            {
              id: 3,
              workflowId: 3,
              workflowName: '周报自动生成',
              status: 'success',
              statusText: '成功',
              statusClass: 'success',
              startedAt: new Date(Date.now() - 172800000).toISOString(),
              duration: 180
            },
            {
              id: 4,
              workflowId: 4,
              workflowName: '社交媒体内容创建',
              status: 'partial',
              statusText: '部分成功',
              statusClass: 'warning',
              startedAt: new Date(Date.now() - 259200000).toISOString(),
              duration: 420
            }
          ];
          
          loading.value = false;
          
          // 初始化图表
          initChart();
        }, 1000);
      } catch (err) {
        error.value = '获取仪表盘数据失败: ' + (err.message || '未知错误');
        loading.value = false;
      }
    };
    
    // 初始化图表
    const initChart = () => {
      if (!chartContainer.value) return;
      
      if (chartInstance) {
        chartInstance.dispose();
      }
      
      chartInstance = echarts.init(chartContainer.value);
      updateChart();
    };
    
    // 更新图表数据
    const updateChart = () => {
      // 这里应该根据selectedPeriod从API获取数据
      // 暂时使用模拟数据
      const days = selectedPeriod.value;
      const dates = Array.from({ length: days }, (_, i) => {
        const date = new Date();
        date.setDate(date.getDate() - (days - i - 1));
        return `${date.getMonth() + 1}/${date.getDate()}`;
      });
      
      const successData = Array.from({ length: days }, () => Math.floor(Math.random() * 10) + 5);
      const failData = Array.from({ length: days }, () => Math.floor(Math.random() * 3));
      
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
      };
      
      chartInstance?.setOption(option);
    };
    
    // 格式化日期
    const formatDate = (dateString) => {
      const date = new Date(dateString);
      const now = new Date();
      const diffTime = now - date;
      const diffDays = Math.floor(diffTime / (1000 * 60 * 60 * 24));
      
      if (diffDays === 0) {
        return '今天 ' + date.getHours().toString().padStart(2, '0') + ':' + 
               date.getMinutes().toString().padStart(2, '0');
      } else if (diffDays === 1) {
        return '昨天 ' + date.getHours().toString().padStart(2, '0') + ':' + 
               date.getMinutes().toString().padStart(2, '0');
      } else if (diffDays < 7) {
        return diffDays + '天前';
      } else {
        return date.getFullYear() + '-' + 
               (date.getMonth() + 1).toString().padStart(2, '0') + '-' + 
               date.getDate().toString().padStart(2, '0');
      }
    };
    
    // 监听窗口大小变化，调整图表大小
    const handleResize = () => {
      chartInstance?.resize();
    };
    
    onMounted(() => {
      fetchDashboardData();
      window.addEventListener('resize', handleResize);
    });
    
    onBeforeUnmount(() => {
      window.removeEventListener('resize', handleResize);
      if (chartInstance) {
        chartInstance.dispose();
        chartInstance = null;
      }
    });
    
    return {
      loading,
      error,
      userInfo,
      stats,
      recentExecutions,
      chartContainer,
      periods,
      selectedPeriod,
      currentTimeGreeting,
      fetchDashboardData,
      formatDate
    };
  }
};
</script>

<style scoped>
.dashboard-page {
  padding: 1.5rem;
}

.spinner {
  border: 3px solid #f3f3f3;
  border-top: 3px solid #3498db;
  border-radius: 50%;
  width: 30px;
  height: 30px;
  animation: spin 1s linear infinite;
  margin: 1rem auto;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.badge {
  @apply px-2 py-0.5 rounded-full text-xs font-semibold;
}

.badge-success {
  @apply bg-green-100 text-green-800;
}

.badge-error {
  @apply bg-red-100 text-red-800;
}

.badge-warning {
  @apply bg-amber-100 text-amber-800;
}

.badge-info {
  @apply bg-blue-100 text-blue-800;
}
</style> 