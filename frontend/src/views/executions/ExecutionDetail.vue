<template>
  <div class="execution-detail">
    <!-- 标题和返回按钮 -->
    <div class="flex justify-between items-center mb-6">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">
          执行详情
        </h1>
        <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
          查看工作流执行的详细信息、日志和结果
        </p>
      </div>
      <router-link
        :to="{ name: 'Executions' }"
        class="px-4 py-2 text-sm font-medium text-gray-700 dark:text-gray-200 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm hover:bg-gray-50 dark:hover:bg-gray-700 focus:outline-none"
      >
        返回列表
      </router-link>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="flex justify-center items-center py-12">
      <div class="spinner"></div>
      <span class="ml-3">加载执行详情...</span>
    </div>

    <!-- 错误状态 -->
    <div v-else-if="error" class="bg-red-50 dark:bg-red-900/20 p-4 rounded-md">
      <div class="flex">
        <div class="flex-shrink-0">
          <svg class="h-5 w-5 text-red-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
            <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
          </svg>
        </div>
        <div class="ml-3">
          <h3 class="text-sm font-medium text-red-800 dark:text-red-200">加载失败</h3>
          <div class="mt-2 text-sm text-red-700 dark:text-red-300">
            <p>{{ error }}</p>
          </div>
        </div>
      </div>
    </div>

    <!-- 执行详情内容 -->
    <div v-else-if="execution" class="space-y-6">
      <!-- 基本信息卡片 -->
      <div class="bg-white dark:bg-gray-800 shadow rounded-lg overflow-hidden">
        <div class="px-6 py-5 border-b border-gray-200 dark:border-gray-700">
          <h2 class="text-lg font-medium text-gray-900 dark:text-white">基本信息</h2>
        </div>
        <div class="px-6 py-5">
          <dl class="grid grid-cols-1 md:grid-cols-2 gap-x-4 gap-y-6">
            <div class="col-span-1">
              <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">执行ID</dt>
              <dd class="mt-1 text-sm text-gray-900 dark:text-gray-200">{{ execution.id }}</dd>
            </div>
            <div class="col-span-1">
              <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">工作流</dt>
              <dd class="mt-1 text-sm text-gray-900 dark:text-gray-200">
                <router-link :to="`/workflows/${execution.workflowId}`" class="text-primary-600 dark:text-primary-400 hover:underline">
                  {{ execution.workflowName }}
                </router-link>
              </dd>
            </div>
            <div class="col-span-1">
              <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">触发方式</dt>
              <dd class="mt-1 text-sm text-gray-900 dark:text-gray-200">{{ triggerText(execution) }}</dd>
            </div>
            <div class="col-span-1">
              <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">状态</dt>
              <dd class="mt-1 text-sm">
                <span :class="statusBadgeClass(execution.status)" class="px-2 py-1 text-xs font-medium rounded-full">
                  {{ statusText(execution.status) }}
                </span>
              </dd>
            </div>
            <div class="col-span-1">
              <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">开始时间</dt>
              <dd class="mt-1 text-sm text-gray-900 dark:text-gray-200">{{ formatDate(execution.startedAt) }}</dd>
            </div>
            <div class="col-span-1">
              <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">结束时间</dt>
              <dd class="mt-1 text-sm text-gray-900 dark:text-gray-200">{{ execution.completedAt ? formatDate(execution.completedAt) : '尚未完成' }}</dd>
            </div>
            <div class="col-span-1">
              <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">执行时长</dt>
              <dd class="mt-1 text-sm text-gray-900 dark:text-gray-200">{{ formatDuration(execution.duration) }}</dd>
            </div>
            <div class="col-span-1" v-if="execution.agentId">
              <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">执行代理</dt>
              <dd class="mt-1 text-sm text-gray-900 dark:text-gray-200">
                <router-link :to="`/agents/${execution.agentId}`" class="text-primary-600 dark:text-primary-400 hover:underline">
                  {{ execution.agentName || execution.agentId }}
                </router-link>
              </dd>
            </div>
          </dl>
        </div>
      </div>

      <!-- 取消按钮 (仅当执行状态为运行中或等待中时显示) -->
      <div v-if="['running', 'pending'].includes(execution.status)" class="flex justify-end">
        <button
          @click="cancelExecution"
          class="px-4 py-2 bg-red-600 text-white rounded-md hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500"
          :disabled="cancelling"
        >
          <span v-if="cancelling">取消中...</span>
          <span v-else>取消执行</span>
        </button>
      </div>

      <!-- 选项卡导航 -->
      <div class="border-b border-gray-200 dark:border-gray-700">
        <nav class="-mb-px flex space-x-8">
          <button
            v-for="tab in tabs"
            :key="tab.id"
            @click="activeTab = tab.id"
            :class="[
              activeTab === tab.id
                ? 'border-primary-500 text-primary-600 dark:text-primary-400'
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300 dark:text-gray-400 dark:hover:text-gray-300 dark:hover:border-gray-600',
              'whitespace-nowrap py-4 px-1 border-b-2 font-medium text-sm'
            ]"
          >
            {{ tab.name }}
          </button>
        </nav>
      </div>

      <!-- 选项卡内容 -->
      <div class="mt-6">
        <!-- 日志选项卡 -->
        <div v-if="activeTab === 'logs'" class="bg-white dark:bg-gray-800 shadow rounded-lg overflow-hidden">
          <div class="px-6 py-5 border-b border-gray-200 dark:border-gray-700 flex justify-between items-center">
            <h3 class="text-lg font-medium text-gray-900 dark:text-white">执行日志</h3>
            <button
              @click="refreshLogs"
              class="text-sm text-primary-600 dark:text-primary-400 hover:text-primary-800 dark:hover:text-primary-300 focus:outline-none"
              :disabled="refreshingLogs"
            >
              <span v-if="refreshingLogs">刷新中...</span>
              <span v-else>刷新日志</span>
            </button>
          </div>
          <div class="p-6">
            <div v-if="execution.logs && execution.logs.length > 0" class="bg-gray-100 dark:bg-gray-900 p-4 rounded-md overflow-auto max-h-96 font-mono text-sm">
              <div v-for="(log, index) in execution.logs" :key="index" class="mb-1">
                <span class="text-gray-500 dark:text-gray-400">{{ formatLogTime(log.timestamp) }}</span>
                <span :class="getLogLevelClass(log.level)" class="ml-2">{{ log.message }}</span>
              </div>
            </div>
            <div v-else class="text-center py-8 text-gray-500 dark:text-gray-400">
              暂无日志记录
            </div>
          </div>
        </div>

        <!-- 输入数据选项卡 -->
        <div v-if="activeTab === 'input'" class="bg-white dark:bg-gray-800 shadow rounded-lg overflow-hidden">
          <div class="px-6 py-5 border-b border-gray-200 dark:border-gray-700">
            <h3 class="text-lg font-medium text-gray-900 dark:text-white">输入数据</h3>
          </div>
          <div class="p-6">
            <div v-if="execution.inputData" class="bg-gray-100 dark:bg-gray-900 p-4 rounded-md overflow-auto max-h-96">
              <pre class="text-sm font-mono">{{ formatJson(execution.inputData) }}</pre>
            </div>
            <div v-else class="text-center py-8 text-gray-500 dark:text-gray-400">
              无输入数据
            </div>
          </div>
        </div>

        <!-- 输出数据选项卡 -->
        <div v-if="activeTab === 'output'" class="bg-white dark:bg-gray-800 shadow rounded-lg overflow-hidden">
          <div class="px-6 py-5 border-b border-gray-200 dark:border-gray-700">
            <h3 class="text-lg font-medium text-gray-900 dark:text-white">输出结果</h3>
          </div>
          <div class="p-6">
            <div v-if="execution.outputData" class="bg-gray-100 dark:bg-gray-900 p-4 rounded-md overflow-auto max-h-96">
              <pre class="text-sm font-mono">{{ formatJson(execution.outputData) }}</pre>
            </div>
            <div v-else-if="execution.status === 'success'" class="text-center py-8 text-gray-500 dark:text-gray-400">
              执行成功但无返回数据
            </div>
            <div v-else-if="['running', 'pending'].includes(execution.status)" class="text-center py-8 text-gray-500 dark:text-gray-400">
              执行尚未完成，暂无输出数据
            </div>
            <div v-else class="text-center py-8 text-gray-500 dark:text-gray-400">
              无输出数据
            </div>
          </div>
        </div>

        <!-- 错误信息选项卡 -->
        <div v-if="activeTab === 'error'" class="bg-white dark:bg-gray-800 shadow rounded-lg overflow-hidden">
          <div class="px-6 py-5 border-b border-gray-200 dark:border-gray-700">
            <h3 class="text-lg font-medium text-gray-900 dark:text-white">错误信息</h3>
          </div>
          <div class="p-6">
            <div v-if="execution.errorMessage" class="bg-red-50 dark:bg-red-900/20 p-4 rounded-md text-red-700 dark:text-red-300">
              <pre class="text-sm font-mono whitespace-pre-wrap">{{ execution.errorMessage }}</pre>
            </div>
            <div v-else class="text-center py-8 text-gray-500 dark:text-gray-400">
              无错误信息
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 未找到执行记录 -->
    <div v-else class="bg-yellow-50 dark:bg-yellow-900/20 p-4 rounded-md">
      <div class="flex">
        <div class="flex-shrink-0">
          <svg class="h-5 w-5 text-yellow-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
            <path fill-rule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clip-rule="evenodd" />
          </svg>
        </div>
        <div class="ml-3">
          <h3 class="text-sm font-medium text-yellow-800 dark:text-yellow-200">未找到执行记录</h3>
          <div class="mt-2 text-sm text-yellow-700 dark:text-yellow-300">
            <p>无法找到ID为 "{{ $route.params.id }}" 的执行记录，请返回列表查看其他执行记录。</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useExecutionStore } from '../../stores/execution';
import { useToast } from 'vue-toastification';

const route = useRoute();
const router = useRouter();
const executionStore = useExecutionStore();
const toast = useToast();

// 状态变量
const loading = ref(true);
const error = ref(null);
const execution = ref(null);
const activeTab = ref('logs');
const cancelling = ref(false);
const refreshingLogs = ref(false);

// 选项卡定义
const tabs = [
  { id: 'logs', name: '执行日志' },
  { id: 'input', name: '输入数据' },
  { id: 'output', name: '输出结果' },
  { id: 'error', name: '错误信息' }
];

// 获取执行详情
const fetchExecutionDetails = async () => {
  const executionId = route.params.id;
  if (!executionId) {
    error.value = '无效的执行ID';
    loading.value = false;
    return;
  }

  loading.value = true;
  error.value = null;

  try {
    const result = await executionStore.fetchExecutionById(executionId);
    execution.value = result;
  } catch (err) {
    error.value = err.message || '获取执行详情失败';
    console.error('获取执行详情失败:', err);
  } finally {
    loading.value = false;
  }
};

// 刷新日志
const refreshLogs = async () => {
  if (!execution.value || refreshingLogs.value) return;
  
  refreshingLogs.value = true;
  try {
    await fetchExecutionDetails();
    toast.success('日志已刷新');
  } catch (err) {
    toast.error('刷新日志失败: ' + (err.message || '未知错误'));
  } finally {
    refreshingLogs.value = false;
  }
};

// 取消执行
const cancelExecution = async () => {
  if (!execution.value || cancelling.value) return;
  if (!['running', 'pending'].includes(execution.value.status)) return;
  
  cancelling.value = true;
  try {
    await executionStore.cancelExecution(execution.value.id);
    toast.success('执行已取消');
    // 刷新执行详情
    await fetchExecutionDetails();
  } catch (err) {
    toast.error('取消执行失败: ' + (err.message || '未知错误'));
  } finally {
    cancelling.value = false;
  }
};

// 格式化函数
const formatDate = (dateString) => {
  if (!dateString) return '';
  return executionStore.formatDateTime(dateString);
};

const formatDuration = (seconds) => {
  if (!seconds && seconds !== 0) return '计算中...';
  return executionStore.formatDuration(seconds);
};

const formatJson = (json) => {
  try {
    if (typeof json === 'string') {
      return JSON.stringify(JSON.parse(json), null, 2);
    }
    return JSON.stringify(json, null, 2);
  } catch (e) {
    return json;
  }
};

const formatLogTime = (timestamp) => {
  if (!timestamp) return '';
  const date = new Date(timestamp);
  return date.toLocaleTimeString('zh-CN');
};

// 状态文本和样式
const statusText = (status) => {
  const statusMap = {
    'pending': '等待中',
    'running': '运行中',
    'success': '成功',
    'failed': '失败',
    'cancelled': '已取消'
  };
  return statusMap[status] || status;
};

const statusBadgeClass = (status) => {
  const classMap = {
    'pending': 'bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-300',
    'running': 'bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-300',
    'success': 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-300',
    'failed': 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-300',
    'cancelled': 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900/30 dark:text-yellow-300'
  };
  return classMap[status] || 'bg-gray-100 text-gray-800 dark:bg-gray-900/30 dark:text-gray-300';
};

const triggerText = (execution) => {
  const triggerType = execution.triggerType || execution.trigger_type;
  const triggerMap = {
    'manual': '手动触发',
    'scheduled': '定时触发',
    'webhook': 'Webhook触发',
    'api': 'API触发'
  };
  return triggerMap[triggerType] || triggerType || '未知';
};

const getLogLevelClass = (level) => {
  const levelMap = {
    'error': 'text-red-600 dark:text-red-400',
    'warn': 'text-yellow-600 dark:text-yellow-400',
    'info': 'text-blue-600 dark:text-blue-400',
    'debug': 'text-gray-600 dark:text-gray-400'
  };
  return levelMap[level?.toLowerCase()] || 'text-gray-900 dark:text-gray-200';
};

// 组件挂载时获取执行详情
onMounted(() => {
  fetchExecutionDetails();
});
</script>

<style scoped>
.spinner {
  border: 2px solid rgba(0, 0, 0, 0.1);
  border-radius: 50%;
  border-top: 2px solid #3498db;
  width: 20px;
  height: 20px;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

@media (prefers-color-scheme: dark) {
  .spinner {
    border-color: rgba(255, 255, 255, 0.1);
    border-top-color: #3498db;
  }
}
</style> 