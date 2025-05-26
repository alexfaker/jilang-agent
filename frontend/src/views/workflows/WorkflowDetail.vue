<template>
  <div class="py-6">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-2xl font-semibold text-gray-900 dark:text-white">
            {{ loading ? '加载中...' : workflow.name }}
          </h1>
          <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
            {{ loading ? '' : workflow.description }}
          </p>
        </div>
        <div class="flex space-x-3">
          <button
            @click="$router.push('/workflows')"
            class="inline-flex items-center px-4 py-2 border border-gray-300 dark:border-gray-600 shadow-sm text-sm font-medium rounded-md text-gray-700 dark:text-gray-200 bg-white dark:bg-gray-800 hover:bg-gray-50 dark:hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
          >
            返回
          </button>
          <button
            v-if="!loading && workflow.id"
            @click="executeWorkflow"
            :disabled="executing"
            class="inline-flex items-center px-4 py-2 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
          >
            <svg v-if="executing" class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            {{ executing ? '执行中...' : '执行工作流' }}
          </button>
          <button
            v-if="!loading && workflow.id"
            @click="editMode = !editMode"
            class="inline-flex items-center px-4 py-2 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
          >
            {{ editMode ? '取消编辑' : '编辑工作流' }}
          </button>
        </div>
      </div>

      <!-- 加载中 -->
      <div v-if="loading" class="mt-6">
        <div class="animate-pulse flex space-x-4">
          <div class="flex-1 space-y-4 py-1">
            <div class="h-4 bg-gray-200 dark:bg-gray-700 rounded w-3/4"></div>
            <div class="space-y-2">
              <div class="h-4 bg-gray-200 dark:bg-gray-700 rounded"></div>
              <div class="h-4 bg-gray-200 dark:bg-gray-700 rounded w-5/6"></div>
            </div>
          </div>
        </div>
      </div>

      <!-- 错误提示 -->
      <div v-else-if="error" class="mt-6 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-md p-4">
        <div class="flex">
          <div class="flex-shrink-0">
            <svg class="h-5 w-5 text-red-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
              <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
            </svg>
          </div>
          <div class="ml-3">
            <h3 class="text-sm font-medium text-red-800 dark:text-red-200">加载工作流失败</h3>
            <div class="mt-2 text-sm text-red-700 dark:text-red-300">
              <p>{{ error }}</p>
            </div>
            <div class="mt-4">
              <button
                @click="fetchWorkflow"
                class="inline-flex items-center px-3 py-2 border border-transparent text-sm leading-4 font-medium rounded-md text-red-700 dark:text-red-300 bg-red-50 dark:bg-red-900/30 hover:bg-red-100 dark:hover:bg-red-900/50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500"
              >
                重试
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- 工作流详情 -->
      <div v-else class="mt-6">
        <!-- 查看模式 -->
        <div v-if="!editMode" class="bg-white dark:bg-gray-800 shadow overflow-hidden sm:rounded-lg">
          <div class="px-4 py-5 sm:px-6">
            <h3 class="text-lg leading-6 font-medium text-gray-900 dark:text-white">工作流详情</h3>
            <p class="mt-1 max-w-2xl text-sm text-gray-500 dark:text-gray-400">工作流配置和定义</p>
          </div>
          <div class="border-t border-gray-200 dark:border-gray-700 px-4 py-5 sm:p-0">
            <dl class="sm:divide-y sm:divide-gray-200 dark:sm:divide-gray-700">
              <div class="py-4 sm:py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6">
                <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">ID</dt>
                <dd class="mt-1 text-sm text-gray-900 dark:text-white sm:mt-0 sm:col-span-2">{{ workflow.id }}</dd>
              </div>
              <div class="py-4 sm:py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6">
                <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">名称</dt>
                <dd class="mt-1 text-sm text-gray-900 dark:text-white sm:mt-0 sm:col-span-2">{{ workflow.name }}</dd>
              </div>
              <div class="py-4 sm:py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6">
                <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">描述</dt>
                <dd class="mt-1 text-sm text-gray-900 dark:text-white sm:mt-0 sm:col-span-2">{{ workflow.description || '无描述' }}</dd>
              </div>
              <div class="py-4 sm:py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6">
                <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">状态</dt>
                <dd class="mt-1 text-sm sm:mt-0 sm:col-span-2">
                  <span :class="statusBadgeClass">{{ statusText }}</span>
                </dd>
              </div>
              <div class="py-4 sm:py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6">
                <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">创建时间</dt>
                <dd class="mt-1 text-sm text-gray-900 dark:text-white sm:mt-0 sm:col-span-2">{{ formatDate(workflow.createdAt) }}</dd>
              </div>
              <div class="py-4 sm:py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6">
                <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">更新时间</dt>
                <dd class="mt-1 text-sm text-gray-900 dark:text-white sm:mt-0 sm:col-span-2">{{ formatDate(workflow.updatedAt) }}</dd>
              </div>
              <div class="py-4 sm:py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6">
                <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">最后执行时间</dt>
                <dd class="mt-1 text-sm text-gray-900 dark:text-white sm:mt-0 sm:col-span-2">{{ workflow.lastRunAt ? formatDate(workflow.lastRunAt) : '从未执行' }}</dd>
              </div>
              <div class="py-4 sm:py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6">
                <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">执行次数</dt>
                <dd class="mt-1 text-sm text-gray-900 dark:text-white sm:mt-0 sm:col-span-2">{{ workflow.runCount }}</dd>
              </div>
              <div class="py-4 sm:py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6">
                <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">触发方式</dt>
                <dd class="mt-1 text-sm text-gray-900 dark:text-white sm:mt-0 sm:col-span-2">
                  {{ triggerTypeText }}
                  <div v-if="workflow.trigger_type === 'schedule'" class="mt-1 text-xs text-gray-500 dark:text-gray-400">
                    Cron表达式: {{ workflow.trigger_config?.cron || 'N/A' }}
                  </div>
                </dd>
              </div>
              <div class="py-4 sm:py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6">
                <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">工作流定义</dt>
                <dd class="mt-1 text-sm text-gray-900 dark:text-white sm:mt-0 sm:col-span-2">
                  <div class="bg-gray-50 dark:bg-gray-900 p-4 rounded-md overflow-auto max-h-96">
                    <pre class="text-xs font-mono">{{ formattedDefinition }}</pre>
                  </div>
                </dd>
              </div>
            </dl>
          </div>
        </div>

        <!-- 编辑模式 -->
        <div v-else class="bg-white dark:bg-gray-800 shadow overflow-hidden sm:rounded-lg">
          <form @submit.prevent="saveWorkflow" class="p-6">
            <div class="grid grid-cols-1 gap-y-6 gap-x-4 sm:grid-cols-6">
              <!-- 名称 -->
              <div class="sm:col-span-4">
                <label for="name" class="block text-sm font-medium text-gray-700 dark:text-gray-300">
                  工作流名称
                </label>
                <div class="mt-1">
                  <input
                    type="text"
                    id="name"
                    v-model="editedWorkflow.name"
                    :class="{'border-red-300 focus:ring-red-500 focus:border-red-500': validationErrors.name}"
                    class="shadow-sm focus:ring-primary-500 focus:border-primary-500 block w-full sm:text-sm border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white rounded-md"
                    placeholder="输入工作流名称"
                  />
                  <p v-if="validationErrors.name" class="mt-2 text-sm text-red-600">{{ validationErrors.name }}</p>
                </div>
              </div>

              <!-- 描述 -->
              <div class="sm:col-span-6">
                <label for="description" class="block text-sm font-medium text-gray-700 dark:text-gray-300">
                  描述
                </label>
                <div class="mt-1">
                  <textarea
                    id="description"
                    v-model="editedWorkflow.description"
                    rows="3"
                    :class="{'border-red-300 focus:ring-red-500 focus:border-red-500': validationErrors.description}"
                    class="shadow-sm focus:ring-primary-500 focus:border-primary-500 block w-full sm:text-sm border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white rounded-md"
                    placeholder="描述此工作流的用途和功能"
                  ></textarea>
                  <p v-if="validationErrors.description" class="mt-2 text-sm text-red-600">{{ validationErrors.description }}</p>
                </div>
              </div>

              <!-- 工作流状态 -->
              <div class="sm:col-span-3">
                <label for="status" class="block text-sm font-medium text-gray-700 dark:text-gray-300">
                  状态
                </label>
                <div class="mt-1">
                  <select
                    id="status"
                    v-model="editedWorkflow.status"
                    class="shadow-sm focus:ring-primary-500 focus:border-primary-500 block w-full sm:text-sm border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white rounded-md"
                  >
                    <option value="draft">草稿</option>
                    <option value="active">活跃</option>
                    <option value="inactive">暂停</option>
                    <option value="archived">归档</option>
                  </select>
                </div>
              </div>

              <!-- 工作流定义 - JSON编辑器 -->
              <div class="sm:col-span-6 border-t border-gray-200 dark:border-gray-700 pt-5 mt-4">
                <h3 class="text-lg leading-6 font-medium text-gray-900 dark:text-white">工作流定义</h3>
                <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">使用JSON格式定义工作流的步骤和逻辑</p>
                
                <div class="mt-4">
                  <label for="definition" class="block text-sm font-medium text-gray-700 dark:text-gray-300">
                    JSON 定义
                  </label>
                  <div class="mt-1">
                    <textarea
                      id="definition"
                      v-model="editedDefinition"
                      rows="10"
                      :class="{'border-red-300 focus:ring-red-500 focus:border-red-500': validationErrors.definition}"
                      class="font-mono shadow-sm focus:ring-primary-500 focus:border-primary-500 block w-full sm:text-sm border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white rounded-md"
                      placeholder='{"steps": []}'
                    ></textarea>
                    <p v-if="validationErrors.definition" class="mt-2 text-sm text-red-600">{{ validationErrors.definition }}</p>
                  </div>
                </div>
                
                <div class="mt-4 flex justify-end">
                  <button
                    type="button"
                    @click="formatJSON"
                    class="inline-flex items-center px-3 py-1.5 border border-gray-300 dark:border-gray-600 shadow-sm text-sm font-medium rounded-md text-gray-700 dark:text-gray-300 bg-white dark:bg-gray-800 hover:bg-gray-50 dark:hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
                  >
                    格式化 JSON
                  </button>
                </div>
              </div>
            </div>

            <!-- 表单操作 -->
            <div class="pt-5 mt-6 border-t border-gray-200 dark:border-gray-700">
              <div class="flex justify-end">
                <button
                  type="button"
                  @click="editMode = false"
                  class="bg-white dark:bg-gray-800 py-2 px-4 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm text-sm font-medium text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
                >
                  取消
                </button>
                <button
                  type="submit"
                  :disabled="saving"
                  class="ml-3 inline-flex justify-center py-2 px-4 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
                >
                  <svg v-if="saving" class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                  {{ saving ? '保存中...' : '保存工作流' }}
                </button>
              </div>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, reactive } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useToast } from 'vue-toastification';
import { useWorkflowStore } from '../../stores/workflow';

const route = useRoute();
const router = useRouter();
const toast = useToast();
const workflowStore = useWorkflowStore();

const workflow = ref({});
const editedWorkflow = ref({});
const editedDefinition = ref('');
const loading = ref(true);
const error = ref(null);
const editMode = ref(false);
const saving = ref(false);
const executing = ref(false);
const validationErrors = ref({});

// 获取工作流详情
const fetchWorkflow = async () => {
  const workflowId = parseInt(route.params.id);
  if (!workflowId) {
    error.value = '无效的工作流ID';
    loading.value = false;
    return;
  }

  loading.value = true;
  error.value = null;

  try {
    const result = await workflowStore.fetchWorkflowById(workflowId);
    workflow.value = result;
    
    // 初始化编辑数据
    editedWorkflow.value = { ...result };
    editedDefinition.value = JSON.stringify(result.definition, null, 2);
  } catch (err) {
    error.value = err.message || '加载工作流失败';
    console.error('获取工作流详情失败:', err);
  } finally {
    loading.value = false;
  }
};

// 格式化JSON
const formatJSON = () => {
  try {
    const parsed = JSON.parse(editedDefinition.value);
    editedDefinition.value = JSON.stringify(parsed, null, 2);
    validationErrors.value.definition = '';
  } catch (error) {
    validationErrors.value.definition = `JSON格式错误: ${error.message}`;
  }
};

// 表单验证
const validateForm = () => {
  const errors = {};
  
  if (!editedWorkflow.value.name?.trim()) {
    errors.name = '工作流名称不能为空';
  }
  
  // 验证JSON
  try {
    JSON.parse(editedDefinition.value);
  } catch (error) {
    errors.definition = `JSON格式错误: ${error.message}`;
  }
  
  validationErrors.value = errors;
  return Object.keys(errors).length === 0;
};

// 保存工作流
const saveWorkflow = async () => {
  if (!validateForm()) {
    toast.error('请修正表单中的错误');
    return;
  }
  
  saving.value = true;
  
  try {
    // 准备更新数据
    const updateData = {
      name: editedWorkflow.value.name,
      description: editedWorkflow.value.description,
      status: editedWorkflow.value.status,
      definition: JSON.parse(editedDefinition.value)
    };
    
    // 更新工作流
    const updatedWorkflow = await workflowStore.updateWorkflow(workflow.value.id, updateData);
    
    // 更新本地数据
    workflow.value = updatedWorkflow;
    
    toast.success('工作流更新成功');
    editMode.value = false;
  } catch (error) {
    console.error('更新工作流失败', error);
    
    // 处理特定错误
    if (error.response && error.response.status === 400) {
      // 服务器返回的验证错误
      const serverErrors = error.response.data.errors || {};
      Object.keys(serverErrors).forEach(key => {
        validationErrors.value[key] = serverErrors[key];
      });
      toast.error('表单验证失败，请检查输入');
    } else {
      toast.error('更新工作流失败: ' + (error.message || '未知错误'));
    }
  } finally {
    saving.value = false;
  }
};

// 执行工作流
const executeWorkflow = async () => {
  executing.value = true;
  
  try {
    const result = await workflowStore.executeWorkflow(workflow.value.id);
    toast.success('工作流执行已启动');
    
    // 跳转到执行详情页
    router.push(`/executions/${result.id}`);
  } catch (error) {
    toast.error('执行工作流失败: ' + (error.message || '未知错误'));
  } finally {
    executing.value = false;
  }
};

// 格式化日期
const formatDate = (dateString) => {
  if (!dateString) return '';
  const date = new Date(dateString);
  return new Intl.DateTimeFormat('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    hour12: false
  }).format(date);
};

// 工作流状态文本
const statusText = computed(() => {
  const statusMap = {
    'draft': '草稿',
    'active': '活跃',
    'inactive': '暂停',
    'archived': '归档'
  };
  return statusMap[workflow.value.status] || workflow.value.status;
});

// 工作流状态徽章样式
const statusBadgeClass = computed(() => {
  const baseClass = 'px-2 inline-flex text-xs leading-5 font-semibold rounded-full';
  
  switch (workflow.value.status) {
    case 'active':
      return `${baseClass} bg-green-100 dark:bg-green-900/30 text-green-800 dark:text-green-300`;
    case 'inactive':
      return `${baseClass} bg-yellow-100 dark:bg-yellow-900/30 text-yellow-800 dark:text-yellow-300`;
    case 'draft':
      return `${baseClass} bg-gray-100 dark:bg-gray-700 text-gray-800 dark:text-gray-300`;
    case 'archived':
      return `${baseClass} bg-red-100 dark:bg-red-900/30 text-red-800 dark:text-red-300`;
    default:
      return `${baseClass} bg-gray-100 dark:bg-gray-700 text-gray-800 dark:text-gray-300`;
  }
});

// 触发器类型文本
const triggerTypeText = computed(() => {
  const triggerMap = {
    'manual': '手动触发',
    'schedule': '定时计划',
    'webhook': 'Webhook',
    'event': '事件触发'
  };
  return triggerMap[workflow.value.trigger_type] || '未知';
});

// 格式化定义
const formattedDefinition = computed(() => {
  if (!workflow.value.definition) return '';
  
  try {
    if (typeof workflow.value.definition === 'string') {
      return JSON.stringify(JSON.parse(workflow.value.definition), null, 2);
    } else {
      return JSON.stringify(workflow.value.definition, null, 2);
    }
  } catch (e) {
    return String(workflow.value.definition);
  }
});

// 初始化
onMounted(() => {
  fetchWorkflow();
});
</script>

<style scoped>
.font-mono {
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", "Courier New", monospace;
}
</style> 