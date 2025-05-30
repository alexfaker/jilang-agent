<template>
  <div class="flex-1 overflow-y-auto p-6">
    <!-- 页面标题和操作 -->
    <div class="flex items-center justify-between mb-6">
      <div>
        <h1 class="text-xl font-semibold text-gray-800">统计分析</h1>
        <p class="text-gray-600 mt-1">查看工作流执行统计和数据分析</p>
      </div>
      
      <!-- 右侧操作区 -->
      <div class="flex items-center space-x-4">
        <button 
          @click="refreshStats" 
          :disabled="loading"
          class="text-gray-600 hover:text-gray-800 px-3 py-2 flex items-center transition-colors"
        >
          <ArrowPathIcon :class="loading ? 'animate-spin' : ''" class="w-5 h-5 mr-1" />
          <span class="hidden md:inline">刷新</span>
        </button>
        
        <button 
          @click="exportReport"
          class="text-gray-600 hover:text-gray-800 px-3 py-2 flex items-center transition-colors"
        >
          <ArrowDownTrayIcon class="w-5 h-5 mr-1" />
          <span class="hidden md:inline">导出报告</span>
        </button>
      </div>
    </div>

    <!-- 时间范围筛选器 -->
    <div class="bg-white rounded-lg p-4 mb-6 shadow-sm">
      <div class="flex flex-col md:flex-row md:items-center space-y-4 md:space-y-0">
        <div class="flex flex-wrap gap-2">
          <select 
            v-model="timeRange"
            @change="fetchExecutionStats"
            class="border border-gray-300 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 bg-white"
          >
            <option value="7d">最近7天</option>
            <option value="30d">最近30天</option>
            <option value="90d">最近90天</option>
            <option value="all">所有时间</option>
          </select>
        </div>
      </div>
    </div>
    
    <!-- 加载状态 -->
    <div v-if="loading && !stats.total_workflows" class="flex justify-center items-center py-12">
      <ArrowPathIcon class="w-8 h-8 animate-spin text-indigo-600 mr-3" />
      <span class="text-gray-600">正在加载统计数据...</span>
    </div>

    <!-- 错误状态 -->
    <div v-else-if="error" class="bg-red-50 p-4 rounded-lg">
      <div class="flex">
        <div class="flex-shrink-0">
          <ExclamationTriangleIcon class="h-5 w-5 text-red-400" />
        </div>
        <div class="ml-3">
          <h3 class="text-sm font-medium text-red-800">加载失败</h3>
          <div class="mt-2 text-sm text-red-700">
            <p>{{ error }}</p>
          </div>
        </div>
      </div>
    </div>

    <!-- 统计卡片 -->
    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
      <!-- 工作流执行次数 -->
      <div class="bg-white rounded-lg p-6 shadow-sm">
        <div class="flex items-center">
          <div class="p-2 bg-blue-100 rounded-lg">
            <PlayCircleIcon class="w-6 h-6 text-blue-600" />
          </div>
          <div class="ml-4">
            <h3 class="text-sm font-medium text-gray-500">工作流执行次数</h3>
            <div class="flex items-baseline">
              <p class="text-2xl font-semibold text-gray-900">{{ formatNumber(stats.total_executions || 0) }}</p>
              <p class="ml-2 text-sm font-medium text-green-600 flex items-center">
                <ArrowTrendingUpIcon class="w-4 h-4 mr-1" />
                {{ getChangePercent('executions') }}%
              </p>
            </div>
          </div>
        </div>
      </div>

      <!-- 资源点数消耗 -->
      <div class="bg-white rounded-lg p-6 shadow-sm">
        <div class="flex items-center">
          <div class="p-2 bg-yellow-100 rounded-lg">
            <CurrencyDollarIcon class="w-6 h-6 text-yellow-600" />
          </div>
          <div class="ml-4">
            <h3 class="text-sm font-medium text-gray-500">资源点数消耗</h3>
            <div class="flex items-baseline">
              <p class="text-2xl font-semibold text-gray-900">{{ formatNumber(estimatedPointsUsed) }}</p>
              <p class="ml-2 text-sm font-medium text-green-600 flex items-center">
                <ArrowTrendingUpIcon class="w-4 h-4 mr-1" />
                {{ getChangePercent('points') }}%
              </p>
            </div>
          </div>
        </div>
      </div>

      <!-- 工作流成功率 -->
      <div class="bg-white rounded-lg p-6 shadow-sm">
        <div class="flex items-center">
          <div class="p-2 bg-green-100 rounded-lg">
            <CheckCircleIcon class="w-6 h-6 text-green-600" />
          </div>
          <div class="ml-4">
            <h3 class="text-sm font-medium text-gray-500">工作流成功率</h3>
            <div class="flex items-baseline">
              <p class="text-2xl font-semibold text-gray-900">{{ formatPercent(stats.success_rate || 0) }}%</p>
              <p class="ml-2 text-sm font-medium text-green-600 flex items-center">
                <ArrowTrendingUpIcon class="w-4 h-4 mr-1" />
                {{ getChangePercent('success_rate') }}%
              </p>
            </div>
          </div>
        </div>
      </div>

      <!-- 平均执行时间 -->
      <div class="bg-white rounded-lg p-6 shadow-sm">
        <div class="flex items-center">
          <div class="p-2 bg-indigo-100 rounded-lg">
            <ClockIcon class="w-6 h-6 text-indigo-600" />
          </div>
          <div class="ml-4">
            <h3 class="text-sm font-medium text-gray-500">平均执行时间</h3>
            <div class="flex items-baseline">
              <p class="text-2xl font-semibold text-gray-900">{{ formatDuration(executionStats.avg_duration || 0) }}</p>
              <p class="ml-2 text-sm font-medium" :class="getDurationChangeClass()">
                <ArrowTrendingUpIcon v-if="getDurationChangeClass().includes('green')" class="w-4 h-4 mr-1" />
                <ArrowTrendingDownIcon v-else class="w-4 h-4 mr-1" />
                {{ getChangePercent('duration') }}%
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 图表区域 -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
      <!-- 资源点数使用趋势 -->
      <div class="bg-white rounded-lg p-6 shadow-sm">
        <div class="flex items-center justify-between mb-4">
          <h3 class="text-lg font-medium text-gray-900">资源点数使用趋势</h3>
          <div class="flex items-center space-x-2">
            <button 
              @click="chartType = 'points'"
              :class="chartType === 'points' ? 'bg-indigo-100 text-indigo-700' : 'text-gray-500 hover:text-gray-700'"
              class="px-3 py-1 rounded-lg text-sm font-medium transition-colors"
            >
              点数消耗
            </button>
            <button 
              @click="chartType = 'executions'"
              :class="chartType === 'executions' ? 'bg-indigo-100 text-indigo-700' : 'text-gray-500 hover:text-gray-700'"
              class="px-3 py-1 rounded-lg text-sm font-medium transition-colors"
            >
              执行次数
            </button>
          </div>
        </div>
        
        <!-- 这里可以放置图表 -->
        <div class="h-80 flex items-center justify-center bg-gray-50 rounded-lg">
          <div class="text-center">
            <ChartBarIcon class="w-12 h-12 mx-auto text-gray-400 mb-2" />
            <p class="text-gray-500">{{ chartType === 'points' ? '点数消耗' : '执行次数' }}趋势图</p>
            <p class="text-xs text-gray-400 mt-1">
              过去{{ timeRange === '7d' ? '7天' : timeRange === '30d' ? '30天' : timeRange === '90d' ? '90天' : '所有时间' }}
            </p>
          </div>
        </div>
      </div>

      <!-- 执行状态分布 -->
      <div class="bg-white rounded-lg p-6 shadow-sm">
        <h3 class="text-lg font-medium text-gray-900 mb-4">执行状态分布</h3>
        
        <!-- 状态统计 -->
        <div class="space-y-4">
          <div v-for="(count, status) in executionStats.status_counts" :key="status" class="flex items-center justify-between">
            <div class="flex items-center">
              <div 
                :class="getStatusColor(status)"
                class="w-4 h-4 rounded-full mr-3"
              ></div>
              <span class="text-sm font-medium text-gray-700">{{ getStatusText(status) }}</span>
            </div>
            <div class="flex items-center">
              <span class="text-sm text-gray-900 font-medium">{{ count }}</span>
              <span class="text-xs text-gray-500 ml-2">
                ({{ getStatusPercent(status, count) }}%)
              </span>
            </div>
          </div>
        </div>
        
        <!-- 状态条 -->
        <div class="mt-6">
          <div class="flex rounded-lg overflow-hidden h-3">
            <div 
              v-for="(count, status) in executionStats.status_counts" 
              :key="status"
              :style="{ width: getStatusPercent(status, count) + '%' }"
              :class="getStatusBgColor(status)"
            ></div>
          </div>
        </div>
      </div>
    </div>

    <!-- 工作流统计表格 -->
    <div v-if="workflowStats.length > 0" class="mt-8 bg-white rounded-lg shadow-sm">
      <div class="px-6 py-4 border-b border-gray-200">
        <h3 class="text-lg font-medium text-gray-900">工作流性能统计</h3>
        <p class="text-sm text-gray-500 mt-1">查看各个工作流的执行情况和成功率</p>
      </div>
      
      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">工作流名称</th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">总执行次数</th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">成功次数</th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">失败次数</th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">成功率</th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">平均耗时</th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-for="workflow in workflowStats" :key="workflow.workflow_id" class="hover:bg-gray-50">
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="flex items-center">
                  <RectangleGroupIcon class="w-5 h-5 text-indigo-500 mr-2" />
                  <div class="text-sm font-medium text-gray-900">{{ workflow.workflow_name }}</div>
                </div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                {{ workflow.total_runs }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-green-600">
                {{ workflow.success_runs }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-red-600">
                {{ workflow.failure_runs }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="flex items-center">
                  <div class="w-16 bg-gray-200 rounded-full h-2 mr-2">
                    <div 
                      class="bg-green-500 h-2 rounded-full" 
                      :style="{ width: workflow.success_rate + '%' }"
                    ></div>
                  </div>
                  <span class="text-sm text-gray-900">{{ formatPercent(workflow.success_rate) }}%</span>
                </div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                {{ formatDuration(workflow.avg_duration_ms) }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { statsApi } from '../../api';
import { useToast } from 'vue-toastification';
import {
  ArrowPathIcon,
  ArrowDownTrayIcon,
  PlayCircleIcon,
  CurrencyDollarIcon,
  CheckCircleIcon,
  ClockIcon,
  ChartBarIcon,
  RectangleGroupIcon,
  ArrowTrendingUpIcon,
  ArrowTrendingDownIcon,
  ExclamationTriangleIcon
} from '@heroicons/vue/24/outline';

const toast = useToast();

// 响应式数据
const loading = ref(false);
const error = ref(null);
const timeRange = ref('30d');
const chartType = ref('points');

// 统计数据
const stats = ref({});
const executionStats = ref({});
const workflowStats = ref([]);

// 上期数据用于计算变化百分比（模拟数据）
const previousStats = ref({
  total_executions: 0,
  success_rate: 0,
  avg_duration: 0
});

// 计算属性
const estimatedPointsUsed = computed(() => {
  // 假设每次执行平均消耗25个积分
  return (stats.value.total_executions || 0) * 25;
});

// 方法
const fetchDashboardStats = async () => {
  try {
    const response = await statsApi.getDashboardStats();
    if (response.status === 'success') {
      stats.value = response.data;
      return response.data;
    } else {
      throw new Error(response.message || '获取仪表盘统计失败');
    }
  } catch (err) {
    console.error('获取仪表盘统计失败:', err);
    throw err;
  }
};

const fetchExecutionStats = async () => {
  try {
    const params = {};
    
    // 根据时间范围设置参数
    if (timeRange.value !== 'all') {
      const now = new Date();
      let startDate;
      
      switch (timeRange.value) {
        case '7d':
          startDate = new Date(now.getTime() - 7 * 24 * 60 * 60 * 1000);
          break;
        case '30d':
          startDate = new Date(now.getTime() - 30 * 24 * 60 * 60 * 1000);
          break;
        case '90d':
          startDate = new Date(now.getTime() - 90 * 24 * 60 * 60 * 1000);
          break;
      }
      
      if (startDate) {
        params.start_date = startDate.toISOString().split('T')[0];
        params.end_date = now.toISOString().split('T')[0];
      }
    }
    
    const response = await statsApi.getExecutionStats(params);
    if (response.status === 'success') {
      executionStats.value = response.data;
      return response.data;
    } else {
      throw new Error(response.message || '获取执行统计失败');
    }
  } catch (err) {
    console.error('获取执行统计失败:', err);
    throw err;
  }
};

const fetchWorkflowStats = async () => {
  try {
    const response = await statsApi.getWorkflowStats();
    if (response.status === 'success') {
      workflowStats.value = response.data || [];
      return response.data;
    } else {
      throw new Error(response.message || '获取工作流统计失败');
    }
  } catch (err) {
    console.error('获取工作流统计失败:', err);
    throw err;
  }
};

const refreshStats = async () => {
  loading.value = true;
  error.value = null;
  
  try {
    await Promise.all([
      fetchDashboardStats(),
      fetchExecutionStats(),
      fetchWorkflowStats()
    ]);
    toast.success('统计数据已刷新');
  } catch (err) {
    console.error('刷新统计数据失败:', err);
    error.value = err.message || '刷新统计数据失败';
    toast.error('刷新统计数据失败：' + (err.message || '未知错误'));
  } finally {
    loading.value = false;
  }
};

const exportReport = () => {
  toast.info('导出功能开发中...');
};

// 格式化函数
const formatNumber = (num) => {
  if (num >= 1000000) {
    return (num / 1000000).toFixed(1) + 'M';
  } else if (num >= 1000) {
    return (num / 1000).toFixed(1) + 'K';
  }
  return num.toString();
};

const formatPercent = (percent) => {
  return Math.round(percent * 10) / 10;
};

const formatDuration = (milliseconds) => {
  if (!milliseconds) return '0秒';
  
  const seconds = Math.floor(milliseconds / 1000);
  if (seconds < 60) {
    return `${seconds}秒`;
  }
  
  const minutes = Math.floor(seconds / 60);
  const remainingSeconds = seconds % 60;
  return `${minutes}分${remainingSeconds}秒`;
};

// 状态相关函数
const getStatusText = (status) => {
  const statusMap = {
    succeeded: '成功',
    failed: '失败',
    running: '运行中',
    pending: '等待中',
    cancelled: '已取消'
  };
  return statusMap[status] || status;
};

const getStatusColor = (status) => {
  const colorMap = {
    succeeded: 'bg-green-500',
    failed: 'bg-red-500',
    running: 'bg-blue-500',
    pending: 'bg-yellow-500',
    cancelled: 'bg-gray-500'
  };
  return colorMap[status] || 'bg-gray-500';
};

const getStatusBgColor = (status) => {
  const colorMap = {
    succeeded: 'bg-green-500',
    failed: 'bg-red-500',
    running: 'bg-blue-500',
    pending: 'bg-yellow-500',
    cancelled: 'bg-gray-500'
  };
  return colorMap[status] || 'bg-gray-500';
};

const getStatusPercent = (status, count) => {
  const total = Object.values(executionStats.value.status_counts || {}).reduce((sum, c) => sum + c, 0);
  return total > 0 ? Math.round((count / total) * 100) : 0;
};

// 变化百分比计算（模拟数据）
const getChangePercent = (type) => {
  const changes = {
    executions: 12.5,
    points: 8.2,
    success_rate: 1.2,
    duration: 5.3
  };
  return changes[type] || 0;
};

const getDurationChangeClass = () => {
  // 执行时间增加是负面的，所以用红色表示
  return 'text-red-600';
};

onMounted(() => {
  refreshStats();
});
</script> 