<template>
  <div class="flex-1 overflow-y-auto p-6">
    <!-- 顶部操作栏 -->
    <div class="flex items-center justify-between mb-6">
      <div>
        <h1 class="text-xl font-semibold text-gray-800">执行历史</h1>
      </div>
      
      <!-- 右侧操作区 -->
      <div class="flex items-center space-x-4">
        <button 
          @click="fetchExecutions" 
          :disabled="loading"
          class="text-gray-600 hover:text-gray-800 px-3 py-2 flex items-center transition-colors"
        >
          <ArrowPathIcon :class="loading ? 'animate-spin' : ''" class="w-5 h-5 mr-1" />
          <span class="hidden md:inline">刷新</span>
        </button>
        
        <button 
          @click="exportExecutions"
          class="text-gray-600 hover:text-gray-800 px-3 py-2 flex items-center transition-colors"
        >
          <ArrowDownTrayIcon class="w-5 h-5 mr-1" />
          <span class="hidden md:inline">导出</span>
        </button>
      </div>
    </div>
    
    <!-- 筛选器和搜索 -->
    <div class="bg-white rounded-lg p-4 mb-6 shadow-sm">
      <div class="flex flex-col md:flex-row md:items-center space-y-4 md:space-y-0">
        <!-- 搜索框 -->
        <div class="flex-1 md:mr-4">
          <div class="relative">
            <input 
              type="text" 
              placeholder="搜索执行记录..." 
              v-model="searchQuery"
              @input="debounceSearch"
              class="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
            >
            <MagnifyingGlassIcon class="w-5 h-5 absolute left-3 top-3 text-gray-400" />
          </div>
        </div>
        
        <!-- 过滤器 -->
        <div class="flex flex-wrap gap-2">
          <select 
            v-model="workflowFilter"
            @change="fetchExecutions"
            class="border border-gray-300 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 bg-white"
          >
            <option value="">所有工作流</option>
            <option value="wf1">数据处理自动化</option>
            <option value="wf2">客户反馈分析</option>
            <option value="wf3">内容生成与分发</option>
          </select>
          
          <select 
            v-model="statusFilter"
            @change="fetchExecutions"
            class="border border-gray-300 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 bg-white"
          >
            <option value="">所有状态</option>
            <option value="success">成功</option>
            <option value="failed">失败</option>
            <option value="running">运行中</option>
            <option value="cancelled">已取消</option>
          </select>
          
          <select 
            v-model="timeFilter"
            @change="fetchExecutions"
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
    <div v-if="loading" class="flex justify-center items-center py-12">
      <ArrowPathIcon class="w-8 h-8 animate-spin text-indigo-600 mr-3" />
      <span class="text-gray-600">加载执行历史...</span>
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

    <!-- 执行记录为空 -->
    <div v-else-if="executions.length === 0" class="text-center py-12">
      <RectangleStackIcon class="mx-auto h-12 w-12 text-gray-400" />
      <h3 class="mt-2 text-sm font-medium text-gray-900">暂无执行记录</h3>
      <p class="mt-1 text-sm text-gray-500">
        {{ searchQuery || statusFilter || workflowFilter || timeFilter !== 'all' ? '没有符合条件的执行记录，请尝试修改筛选条件。' : '暂无任何工作流执行记录，请先执行一个工作流。' }}
      </p>
    </div>
    
    <!-- 执行历史列表 -->
    <div v-else class="bg-white rounded-lg overflow-hidden shadow-sm">
      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">执行ID</th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">工作流</th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">开始时间</th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">持续时间</th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">状态</th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">触发方式</th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">操作</th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-for="execution in executions" :key="execution.id" class="hover:bg-gray-50">
              <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                #EXE-{{ execution.id.substring(0, 5) }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                <div class="flex items-center">
                  <RectangleGroupIcon class="w-5 h-5 text-indigo-500 mr-2" />
                  <router-link 
                    :to="`/workflows/${execution.workflow_id}`" 
                    class="hover:text-indigo-600 transition-colors"
                  >
                    {{ execution.workflow_name || execution.workflowName }}
                  </router-link>
                </div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                {{ formatDate(execution.start_time || execution.startedAt) }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                {{ formatDuration(execution.duration) }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <span 
                  :class="statusBadgeClass(execution.status)" 
                  class="px-2 py-1 inline-flex text-xs leading-5 font-semibold rounded-full"
                >
                  {{ statusText(execution.status) }}
                </span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                {{ triggerText(execution.trigger_type || execution.triggerType) }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
                <router-link
                  :to="`/executions/${execution.id}`"
                  class="text-indigo-600 hover:text-indigo-900 mr-2 transition-colors"
                >
                  查看日志
                </router-link>
                <button
                  @click="rerunExecution(execution)"
                  class="text-indigo-600 hover:text-indigo-900 transition-colors"
                >
                  重新执行
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
    
    <!-- 分页 -->
    <div v-if="executions.length > 0" class="py-4 flex items-center justify-between border-t border-gray-200 mt-4">
      <div class="flex-1 flex justify-between sm:hidden">
        <button
          @click="goToPage(currentPage - 1)"
          :disabled="currentPage === 1"
          class="relative inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 transition-colors"
          :class="currentPage === 1 ? 'opacity-50 cursor-not-allowed' : ''"
        >
          上一页
        </button>
        <button
          @click="goToPage(currentPage + 1)"
          :disabled="currentPage === totalPages"
          class="ml-3 relative inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 transition-colors"
          :class="currentPage === totalPages ? 'opacity-50 cursor-not-allowed' : ''"
        >
          下一页
        </button>
      </div>
      <div class="hidden sm:flex-1 sm:flex sm:items-center sm:justify-between">
        <div>
          <p class="text-sm text-gray-700">
            显示第 <span class="font-medium">{{ startItem }}</span> 至 <span class="font-medium">{{ endItem }}</span> 项，共 <span class="font-medium">{{ totalItems }}</span> 项
          </p>
        </div>
        <div>
          <nav class="relative z-0 inline-flex rounded-md shadow-sm -space-x-px" aria-label="Pagination">
            <button
              @click="goToPage(currentPage - 1)"
              :disabled="currentPage === 1"
              class="relative inline-flex items-center px-2 py-2 rounded-l-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 transition-colors"
              :class="currentPage === 1 ? 'opacity-50 cursor-not-allowed' : ''"
            >
              <span class="sr-only">上一页</span>
              <ChevronLeftIcon class="h-5 w-5" />
            </button>
            
            <template v-for="page in displayedPages" :key="page">
              <button
                v-if="page !== '...'"
                @click="goToPage(page)"
                :class="page === currentPage ? 'z-10 bg-indigo-50 border-indigo-500 text-indigo-600' : 'bg-white border-gray-300 text-gray-500 hover:bg-gray-50'"
                class="relative inline-flex items-center px-4 py-2 border text-sm font-medium transition-colors"
              >
                {{ page }}
              </button>
              <span
                v-else
                class="relative inline-flex items-center px-4 py-2 border border-gray-300 bg-white text-sm font-medium text-gray-700"
              >
                ...
              </span>
            </template>
            
            <button
              @click="goToPage(currentPage + 1)"
              :disabled="currentPage === totalPages"
              class="relative inline-flex items-center px-2 py-2 rounded-r-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 transition-colors"
              :class="currentPage === totalPages ? 'opacity-50 cursor-not-allowed' : ''"
            >
              <span class="sr-only">下一页</span>
              <ChevronRightIcon class="h-5 w-5" />
            </button>
          </nav>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { format, formatDistanceStrict } from 'date-fns';
import { zhCN } from 'date-fns/locale';
import { useToast } from 'vue-toastification';
import { useExecutionStore } from '../../stores/execution';
import { useRouter } from 'vue-router';
import {
  ArrowPathIcon,
  ArrowDownTrayIcon,
  MagnifyingGlassIcon,
  ExclamationTriangleIcon,
  RectangleStackIcon,
  RectangleGroupIcon,
  ChevronLeftIcon,
  ChevronRightIcon
} from '@heroicons/vue/24/outline';

const toast = useToast();
const executionStore = useExecutionStore();
const router = useRouter();

const loading = ref(false);
const executions = ref([]);
const searchQuery = ref('');
const statusFilter = ref('');
const workflowFilter = ref('');
const timeFilter = ref('7d');
const totalItems = ref(0);
const currentPage = ref(1);
const pageSize = ref(10);
const error = ref(null);

// 计算属性
const totalPages = computed(() => Math.ceil(totalItems.value / pageSize.value) || 1);
const startItem = computed(() => ((currentPage.value - 1) * pageSize.value) + 1);
const endItem = computed(() => Math.min(currentPage.value * pageSize.value, totalItems.value));

const displayedPages = computed(() => {
  const pages = [];
  if (totalPages.value <= 7) {
    for (let i = 1; i <= totalPages.value; i++) {
      pages.push(i);
    }
  } else {
    pages.push(1);
    if (currentPage.value > 3) {
      pages.push('...');
    }
    
    let startPage = Math.max(2, currentPage.value - 1);
    let endPage = Math.min(totalPages.value - 1, currentPage.value + 1);
    
    if (currentPage.value <= 3) {
      endPage = Math.min(5, totalPages.value - 1);
    }
    
    if (currentPage.value >= totalPages.value - 2) {
      startPage = Math.max(2, totalPages.value - 4);
    }
    
    for (let i = startPage; i <= endPage; i++) {
      pages.push(i);
    }
    
    if (currentPage.value < totalPages.value - 2) {
      pages.push('...');
    }
    
    pages.push(totalPages.value);
  }
  return pages;
});

// 方法
const fetchExecutions = async () => {
  loading.value = true;
  error.value = null;
  
  try {
    // 模拟数据，实际应该调用API
    const mockData = [
      {
        id: 'exe12345',
        workflow_id: 'wf1',
        workflow_name: '数据处理自动化',
        start_time: new Date(Date.now() - 10 * 60 * 1000),
        duration: 200,
        status: 'success',
        trigger_type: 'manual'
      },
      {
        id: 'exe12344',
        workflow_id: 'wf2',
        workflow_name: '客户反馈分析',
        start_time: new Date(Date.now() - 2 * 60 * 60 * 1000),
        duration: 432,
        status: 'warning',
        trigger_type: 'api'
      },
      {
        id: 'exe12343',
        workflow_id: 'wf1',
        workflow_name: '数据处理自动化',
        start_time: new Date(Date.now() - 24 * 60 * 60 * 1000),
        duration: 222,
        status: 'success',
        trigger_type: 'schedule'
      },
      {
        id: 'exe12342',
        workflow_id: 'wf3',
        workflow_name: '内容生成与分发',
        start_time: new Date(Date.now() - 2 * 24 * 60 * 60 * 1000),
        duration: 725,
        status: 'failed',
        trigger_type: 'manual'
      },
      {
        id: 'exe12341',
        workflow_id: 'wf1',
        workflow_name: '数据处理自动化',
        start_time: new Date(Date.now() - 3 * 24 * 60 * 60 * 1000),
        duration: 235,
        status: 'failed',
        trigger_type: 'schedule'
      }
    ];
    
    executions.value = mockData;
    totalItems.value = mockData.length;
    
    await new Promise(resolve => setTimeout(resolve, 500));
  } catch (err) {
    console.error('获取执行历史失败:', err);
    error.value = err.message || '获取执行历史失败';
    toast.error(`获取执行历史失败: ${err.message || '未知错误'}`);
  } finally {
    loading.value = false;
  }
};

const debounceSearch = (() => {
  let timeout;
  return () => {
    clearTimeout(timeout);
    timeout = setTimeout(() => {
      currentPage.value = 1;
      fetchExecutions();
    }, 300);
  };
})();

const formatDate = (date) => {
  if (!date) return '未知';
  return format(new Date(date), 'yyyy-MM-dd HH:mm', { locale: zhCN });
};

const formatDuration = (seconds) => {
  if (seconds === null || seconds === undefined) return '进行中';
  
  if (seconds < 60) {
    return `${seconds}秒`;
  }
  
  const minutes = Math.floor(seconds / 60);
  const remainingSeconds = seconds % 60;
  return `${minutes}分${remainingSeconds}秒`;
};

const goToPage = (page) => {
  if (page < 1 || page > totalPages.value) return;
  currentPage.value = page;
  fetchExecutions();
};

const statusText = (status) => {
  const statusMap = {
    success: '成功',
    failed: '失败',
    warning: '部分成功',
    running: '运行中',
    pending: '等待中',
    cancelled: '已取消'
  };
  return statusMap[status] || status;
};

const statusBadgeClass = (status) => {
  const badgeMap = {
    success: 'bg-green-100 text-green-800',
    failed: 'bg-red-100 text-red-800',
    warning: 'bg-yellow-100 text-yellow-800',
    running: 'bg-blue-100 text-blue-800',
    pending: 'bg-gray-100 text-gray-800',
    cancelled: 'bg-gray-100 text-gray-800'
  };
  return badgeMap[status] || 'bg-gray-100 text-gray-800';
};

const triggerText = (triggerType) => {
  const triggerMap = {
    schedule: '定时',
    manual: '手动',
    api: 'API触发',
    webhook: 'Webhook',
    event: '事件触发'
  };
  return triggerMap[triggerType] || triggerType || '未知';
};

const exportExecutions = () => {
  toast.info('导出功能开发中...');
};

const rerunExecution = (execution) => {
  if (confirm(`确定要重新执行工作流"${execution.workflow_name}"吗？`)) {
    toast.success('工作流已开始重新执行');
  }
};

onMounted(() => {
  fetchExecutions();
});
</script> 