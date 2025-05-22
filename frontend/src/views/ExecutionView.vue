<template>
  <div class="execution-page">
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-semibold text-gray-800">执行历史</h1>
      <div class="flex items-center space-x-4">
        <button @click="fetchExecutions" class="btn btn-secondary">
          <i class="fas fa-sync-alt mr-2"></i> 刷新
        </button>
      </div>
    </div>

    <!-- 筛选器 -->
    <div class="bg-white rounded-lg p-4 shadow-card mb-6">
      <div class="grid grid-cols-1 md:grid-cols-3 lg:grid-cols-4 gap-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">工作流</label>
          <select v-model="filters.workflowId" class="form-select w-full">
            <option value="">所有工作流</option>
            <option v-for="workflow in workflows" :key="workflow.id" :value="workflow.id">
              {{ workflow.name }}
            </option>
          </select>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">状态</label>
          <select v-model="filters.status" class="form-select w-full">
            <option value="">所有状态</option>
            <option value="success">成功</option>
            <option value="failed">失败</option>
            <option value="running">运行中</option>
            <option value="pending">等待中</option>
            <option value="cancelled">已取消</option>
          </select>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">日期范围</label>
          <div class="flex space-x-2">
            <input 
              type="date" 
              v-model="filters.startDate" 
              class="form-control w-full"
              :max="today"
            />
            <input 
              type="date" 
              v-model="filters.endDate" 
              class="form-control w-full"
              :min="filters.startDate"
              :max="today"
            />
          </div>
        </div>
        <div class="flex items-end">
          <button @click="applyFilters" class="btn btn-primary w-full">
            <i class="fas fa-filter mr-2"></i> 应用筛选
          </button>
        </div>
      </div>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="bg-white rounded-lg p-10 shadow-card flex flex-col items-center justify-center">
      <div class="spinner mb-4"></div>
      <p class="text-gray-600">加载执行历史...</p>
    </div>

    <!-- 错误状态 -->
    <div v-else-if="error" class="bg-white rounded-lg p-10 shadow-card flex flex-col items-center justify-center">
      <i class="fas fa-exclamation-circle text-red-500 text-4xl mb-4"></i>
      <p class="text-red-500 mb-4">{{ error }}</p>
      <button @click="fetchExecutions" class="btn btn-primary">
        <i class="fas fa-sync-alt mr-2"></i> 重试
      </button>
    </div>

    <!-- 空状态 -->
    <div v-else-if="executions.length === 0" class="bg-white rounded-lg p-10 shadow-card flex flex-col items-center justify-center">
      <i class="fas fa-history text-gray-400 text-4xl mb-4"></i>
      <p class="text-gray-600 mb-2">暂无执行历史记录</p>
      <p v-if="isFiltered" class="text-gray-500 text-sm mb-4">尝试调整筛选条件</p>
      <router-link to="/workflows" class="btn btn-primary">
        <i class="fas fa-project-diagram mr-2"></i> 查看工作流
      </router-link>
    </div>

    <!-- 执行历史列表 -->
    <div v-else>
      <div class="bg-white rounded-lg shadow-card overflow-hidden mb-6">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                ID
              </th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                工作流
              </th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                状态
              </th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                开始时间
              </th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                执行时长
              </th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                操作
              </th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-for="execution in executions" :key="execution.id" class="hover:bg-gray-50">
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                {{ execution.id }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="flex items-center">
                  <div class="flex-shrink-0 h-8 w-8 rounded-full bg-gray-100 flex items-center justify-center">
                    <i class="fas fa-robot text-gray-600"></i>
                  </div>
                  <div class="ml-4">
                    <div class="text-sm font-medium text-gray-900">
                      {{ execution.workflowName }}
                    </div>
                    <div class="text-sm text-gray-500">
                      ID: {{ execution.workflowId }}
                    </div>
                  </div>
                </div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <span :class="`badge badge-${execution.statusClass}`">
                  {{ execution.statusText }}
                </span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                {{ formatDateTime(execution.startedAt) }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                {{ formatDuration(execution.duration) }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm">
                <div class="flex space-x-2">
                  <button 
                    @click="viewExecution(execution)" 
                    class="text-indigo-600 hover:text-indigo-900"
                    title="查看详情"
                  >
                    <i class="fas fa-eye"></i>
                  </button>
                  <button 
                    v-if="execution.status === 'success'"
                    @click="rerunExecution(execution)" 
                    class="text-green-600 hover:text-green-900"
                    title="重新运行"
                  >
                    <i class="fas fa-redo"></i>
                  </button>
                  <button 
                    v-if="execution.status === 'running' || execution.status === 'pending'"
                    @click="cancelExecution(execution)" 
                    class="text-red-600 hover:text-red-900"
                    title="取消执行"
                  >
                    <i class="fas fa-stop"></i>
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- 分页控件 -->
      <div class="flex justify-between items-center">
        <div class="text-sm text-gray-700">
          显示 <span class="font-medium">{{ startItem }}</span> 到 
          <span class="font-medium">{{ endItem }}</span> 条，共 
          <span class="font-medium">{{ totalItems }}</span> 条记录
        </div>
        <div class="flex space-x-2">
          <button 
            @click="changePage(currentPage - 1)" 
            :disabled="currentPage === 1"
            class="btn btn-sm btn-secondary"
            :class="{ 'opacity-50 cursor-not-allowed': currentPage === 1 }"
          >
            <i class="fas fa-chevron-left mr-1"></i> 上一页
          </button>
          <button 
            @click="changePage(currentPage + 1)" 
            :disabled="currentPage === totalPages"
            class="btn btn-sm btn-secondary"
            :class="{ 'opacity-50 cursor-not-allowed': currentPage === totalPages }"
          >
            下一页 <i class="fas fa-chevron-right ml-1"></i>
          </button>
        </div>
      </div>
    </div>

    <!-- 执行详情模态框 -->
    <div v-if="selectedExecution" class="fixed inset-0 bg-gray-600 bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg shadow-xl w-full max-w-4xl max-h-[90vh] overflow-hidden">
        <div class="px-6 py-4 border-b border-gray-200 flex justify-between items-center">
          <h3 class="text-lg font-medium text-gray-900">
            执行详情 #{{ selectedExecution.id }}
          </h3>
          <button @click="selectedExecution = null" class="text-gray-400 hover:text-gray-500">
            <i class="fas fa-times"></i>
          </button>
        </div>
        <div class="p-6 overflow-y-auto" style="max-height: calc(90vh - 120px);">
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
            <div>
              <h4 class="text-sm font-medium text-gray-500 mb-2">工作流信息</h4>
              <div class="bg-gray-50 p-4 rounded-lg">
                <div class="mb-2">
                  <span class="text-sm text-gray-500">名称：</span>
                  <span class="font-medium">{{ selectedExecution.workflowName }}</span>
                </div>
                <div>
                  <span class="text-sm text-gray-500">ID：</span>
                  <span class="font-medium">{{ selectedExecution.workflowId }}</span>
                </div>
              </div>
            </div>
            <div>
              <h4 class="text-sm font-medium text-gray-500 mb-2">执行信息</h4>
              <div class="bg-gray-50 p-4 rounded-lg">
                <div class="mb-2">
                  <span class="text-sm text-gray-500">状态：</span>
                  <span :class="`badge badge-${selectedExecution.statusClass}`">
                    {{ selectedExecution.statusText }}
                  </span>
                </div>
                <div class="mb-2">
                  <span class="text-sm text-gray-500">开始时间：</span>
                  <span class="font-medium">{{ formatDateTime(selectedExecution.startedAt) }}</span>
                </div>
                <div class="mb-2">
                  <span class="text-sm text-gray-500">完成时间：</span>
                  <span class="font-medium">
                    {{ selectedExecution.completedAt ? formatDateTime(selectedExecution.completedAt) : '未完成' }}
                  </span>
                </div>
                <div>
                  <span class="text-sm text-gray-500">执行时长：</span>
                  <span class="font-medium">{{ formatDuration(selectedExecution.duration) }}</span>
                </div>
              </div>
            </div>
          </div>

          <div class="mb-6">
            <h4 class="text-sm font-medium text-gray-500 mb-2">输入数据</h4>
            <pre class="bg-gray-50 p-4 rounded-lg overflow-x-auto text-sm">{{ formatJson(selectedExecution.inputData) }}</pre>
          </div>

          <div class="mb-6">
            <h4 class="text-sm font-medium text-gray-500 mb-2">输出数据</h4>
            <pre v-if="selectedExecution.outputData" class="bg-gray-50 p-4 rounded-lg overflow-x-auto text-sm">{{ formatJson(selectedExecution.outputData) }}</pre>
            <div v-else class="bg-gray-50 p-4 rounded-lg text-gray-500 text-sm">暂无输出数据</div>
          </div>

          <div v-if="selectedExecution.errorMessage">
            <h4 class="text-sm font-medium text-gray-500 mb-2">错误信息</h4>
            <div class="bg-red-50 text-red-700 p-4 rounded-lg text-sm">
              {{ selectedExecution.errorMessage }}
            </div>
          </div>

          <div>
            <h4 class="text-sm font-medium text-gray-500 mb-2">执行日志</h4>
            <div class="bg-gray-50 p-4 rounded-lg overflow-x-auto h-64 font-mono text-sm">
              <div v-if="selectedExecution.logs" v-html="formatLogs(selectedExecution.logs)"></div>
              <div v-else class="text-gray-500">暂无日志</div>
            </div>
          </div>
        </div>
        <div class="px-6 py-4 border-t border-gray-200 flex justify-end">
          <button @click="selectedExecution = null" class="btn btn-secondary mr-2">
            关闭
          </button>
          <button 
            v-if="selectedExecution.status === 'success'"
            @click="rerunExecution(selectedExecution)" 
            class="btn btn-primary"
          >
            <i class="fas fa-redo mr-2"></i> 重新运行
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue';
import { executionApi, workflowApi } from '../api';

export default {
  name: 'ExecutionView',
  setup() {
    const executions = ref([]);
    const workflows = ref([]);
    const loading = ref(true);
    const error = ref(null);
    const currentPage = ref(1);
    const itemsPerPage = ref(10);
    const totalItems = ref(0);
    const selectedExecution = ref(null);
    
    // 今天日期，用于日期选择器的最大值
    const today = new Date().toISOString().split('T')[0];
    
    // 筛选条件
    const filters = ref({
      workflowId: '',
      status: '',
      startDate: '',
      endDate: ''
    });
    
    // 计算属性
    const totalPages = computed(() => Math.ceil(totalItems.value / itemsPerPage.value) || 1);
    const startItem = computed(() => (currentPage.value - 1) * itemsPerPage.value + 1);
    const endItem = computed(() => Math.min(currentPage.value * itemsPerPage.value, totalItems.value));
    const isFiltered = computed(() => {
      return filters.value.workflowId !== '' || 
             filters.value.status !== '' || 
             filters.value.startDate !== '' || 
             filters.value.endDate !== '';
    });
    
    // 获取执行历史
    const fetchExecutions = async () => {
      loading.value = true;
      error.value = null;
      
      try {
        // 实际项目中应该调用API获取数据
        // const params = {
        //   workflow_id: filters.value.workflowId,
        //   status: filters.value.status,
        //   start_date: filters.value.startDate,
        //   end_date: filters.value.endDate,
        //   limit: itemsPerPage.value,
        //   offset: (currentPage.value - 1) * itemsPerPage.value
        // };
        // const response = await executionApi.getExecutions(params);
        
        // 模拟API响应数据
        setTimeout(() => {
          const mockExecutions = [
            {
              id: 1001,
              workflowId: 1,
              workflowName: '每日数据采集',
              status: 'success',
              statusText: '成功',
              statusClass: 'success',
              startedAt: '2023-10-15T08:30:00Z',
              completedAt: '2023-10-15T08:32:00Z',
              duration: 120,
              inputData: { source: 'api', date: '2023-10-15' },
              outputData: { records: 150, processed: 150, skipped: 0 },
              logs: 'INFO: 开始执行工作流\nINFO: 正在连接数据源\nINFO: 成功获取150条记录\nINFO: 数据处理完成\nINFO: 工作流执行成功'
            },
            {
              id: 1002,
              workflowId: 2,
              workflowName: '客户数据分析',
              status: 'failed',
              statusText: '失败',
              statusClass: 'error',
              startedAt: '2023-10-14T15:45:00Z',
              completedAt: '2023-10-14T15:48:00Z',
              duration: 180,
              inputData: { customer_segment: 'premium', period: 'last_month' },
              outputData: null,
              errorMessage: '数据库连接超时',
              logs: 'INFO: 开始执行工作流\nINFO: 正在连接数据库\nERROR: 数据库连接超时\nERROR: 工作流执行失败'
            },
            {
              id: 1003,
              workflowId: 3,
              workflowName: '周报自动生成',
              status: 'success',
              statusText: '成功',
              statusClass: 'success',
              startedAt: '2023-10-13T08:00:00Z',
              completedAt: '2023-10-13T08:03:00Z',
              duration: 180,
              inputData: { week: '2023-W41' },
              outputData: { report_url: 'https://example.com/reports/2023-W41.pdf' },
              logs: 'INFO: 开始执行工作流\nINFO: 正在收集数据\nINFO: 生成报告\nINFO: 报告已保存\nINFO: 工作流执行成功'
            },
            {
              id: 1004,
              workflowId: 4,
              workflowName: '社交媒体内容创建',
              status: 'running',
              statusText: '运行中',
              statusClass: 'info',
              startedAt: '2023-10-15T10:20:00Z',
              completedAt: null,
              duration: 300,
              inputData: { platform: 'twitter', topic: 'AI news' },
              outputData: null,
              logs: 'INFO: 开始执行工作流\nINFO: 正在收集话题数据\nINFO: 生成内容草稿\nINFO: 正在优化内容...'
            }
          ];
          
          executions.value = mockExecutions;
          totalItems.value = 24; // 模拟总记录数
          loading.value = false;
        }, 800);
      } catch (err) {
        error.value = '获取执行历史失败: ' + (err.message || '未知错误');
        loading.value = false;
      }
    };
    
    // 获取工作流列表（用于筛选）
    const fetchWorkflows = async () => {
      try {
        // 实际项目中应该调用API获取数据
        // const response = await workflowApi.getWorkflows();
        
        // 模拟API响应数据
        workflows.value = [
          { id: 1, name: '每日数据采集' },
          { id: 2, name: '客户数据分析' },
          { id: 3, name: '周报自动生成' },
          { id: 4, name: '社交媒体内容创建' }
        ];
      } catch (err) {
        console.error('获取工作流列表失败:', err);
      }
    };
    
    // 应用筛选
    const applyFilters = () => {
      currentPage.value = 1;
      fetchExecutions();
    };
    
    // 切换页面
    const changePage = (page) => {
      if (page >= 1 && page <= totalPages.value) {
        currentPage.value = page;
        fetchExecutions();
      }
    };
    
    // 查看执行详情
    const viewExecution = (execution) => {
      selectedExecution.value = { ...execution };
    };
    
    // 重新运行工作流
    const rerunExecution = (execution) => {
      if (confirm(`确定要重新运行工作流 "${execution.workflowName}" 吗？`)) {
        alert(`已开始重新运行工作流: ${execution.workflowName}`);
        // 实际项目中应该调用API重新运行工作流
        // workflowApi.executeWorkflow(execution.workflowId, execution.inputData);
      }
    };
    
    // 取消执行
    const cancelExecution = (execution) => {
      if (confirm(`确定要取消工作流 "${execution.workflowName}" 的执行吗？`)) {
        alert(`已取消工作流执行: ${execution.workflowName}`);
        // 实际项目中应该调用API取消工作流执行
        // executionApi.cancelExecution(execution.id);
      }
    };
    
    // 格式化日期时间
    const formatDateTime = (dateString) => {
      if (!dateString) return '';
      const date = new Date(dateString);
      return date.toLocaleString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit'
      });
    };
    
    // 格式化持续时间
    const formatDuration = (seconds) => {
      if (!seconds) return '0秒';
      
      if (seconds < 60) {
        return `${seconds}秒`;
      } else if (seconds < 3600) {
        const minutes = Math.floor(seconds / 60);
        const remainingSeconds = seconds % 60;
        return `${minutes}分${remainingSeconds > 0 ? remainingSeconds + '秒' : ''}`;
      } else {
        const hours = Math.floor(seconds / 3600);
        const minutes = Math.floor((seconds % 3600) / 60);
        return `${hours}小时${minutes > 0 ? minutes + '分' : ''}`;
      }
    };
    
    // 格式化JSON
    const formatJson = (jsonData) => {
      if (!jsonData) return '';
      try {
        if (typeof jsonData === 'string') {
          return JSON.stringify(JSON.parse(jsonData), null, 2);
        }
        return JSON.stringify(jsonData, null, 2);
      } catch (err) {
        return String(jsonData);
      }
    };
    
    // 格式化日志
    const formatLogs = (logs) => {
      if (!logs) return '';
      
      return logs.split('\n').map(line => {
        if (line.includes('ERROR:')) {
          return `<div class="text-red-600">${line}</div>`;
        } else if (line.includes('WARNING:')) {
          return `<div class="text-amber-600">${line}</div>`;
        } else if (line.includes('INFO:')) {
          return `<div class="text-blue-600">${line}</div>`;
        }
        return `<div>${line}</div>`;
      }).join('');
    };
    
    onMounted(() => {
      fetchExecutions();
      fetchWorkflows();
    });
    
    return {
      executions,
      workflows,
      loading,
      error,
      currentPage,
      itemsPerPage,
      totalItems,
      filters,
      selectedExecution,
      today,
      totalPages,
      startItem,
      endItem,
      isFiltered,
      fetchExecutions,
      applyFilters,
      changePage,
      viewExecution,
      rerunExecution,
      cancelExecution,
      formatDateTime,
      formatDuration,
      formatJson,
      formatLogs
    };
  }
};
</script>

<style scoped>
.execution-page {
  padding: 1.5rem;
}

.spinner {
  border: 3px solid #f3f3f3;
  border-top: 3px solid #3498db;
  border-radius: 50%;
  width: 30px;
  height: 30px;
  animation: spin 1s linear infinite;
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

.form-control,
.form-select {
  @apply block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm;
}

.btn {
  @apply inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-offset-2;
}

.btn-primary {
  @apply text-white bg-indigo-600 hover:bg-indigo-700 focus:ring-indigo-500;
}

.btn-secondary {
  @apply text-gray-700 bg-white border-gray-300 hover:bg-gray-50 focus:ring-indigo-500;
}

.btn-sm {
  @apply px-3 py-1 text-xs;
}
</style> 