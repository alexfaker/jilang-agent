<template>
  <div class="py-6">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex items-center justify-between">
        <h1 class="text-2xl font-semibold text-gray-900 dark:text-white">创建工作流</h1>
        <button
          @click="$router.back()"
          class="inline-flex items-center px-4 py-2 border border-gray-300 dark:border-gray-600 shadow-sm text-sm font-medium rounded-md text-gray-700 dark:text-gray-200 bg-white dark:bg-gray-800 hover:bg-gray-50 dark:hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
        >
          返回
        </button>
      </div>

      <div class="mt-6 bg-white dark:bg-gray-800 shadow overflow-hidden sm:rounded-lg">
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
                  v-model="workflow.name"
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
                  v-model="workflow.description"
                  rows="3"
                  :class="{'border-red-300 focus:ring-red-500 focus:border-red-500': validationErrors.description}"
                  class="shadow-sm focus:ring-primary-500 focus:border-primary-500 block w-full sm:text-sm border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white rounded-md"
                  placeholder="描述此工作流的用途和功能"
                ></textarea>
                <p v-if="validationErrors.description" class="mt-2 text-sm text-red-600">{{ validationErrors.description }}</p>
              </div>
              <p class="mt-2 text-sm text-gray-500 dark:text-gray-400">简要描述此工作流的目的和功能</p>
            </div>

            <!-- 触发器类型 -->
            <div class="sm:col-span-3">
              <label for="trigger_type" class="block text-sm font-medium text-gray-700 dark:text-gray-300">
                触发器类型
              </label>
              <div class="mt-1">
                <select
                  id="trigger_type"
                  v-model="workflow.trigger_type"
                  :class="{'border-red-300 focus:ring-red-500 focus:border-red-500': validationErrors.trigger_type}"
                  class="shadow-sm focus:ring-primary-500 focus:border-primary-500 block w-full sm:text-sm border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white rounded-md"
                  @change="onTriggerTypeChange"
                >
                  <option value="">选择触发器类型</option>
                  <option value="schedule">定时计划</option>
                  <option value="webhook">Webhook</option>
                  <option value="manual">手动触发</option>
                  <option value="event">事件触发</option>
                </select>
                <p v-if="validationErrors.trigger_type" class="mt-2 text-sm text-red-600">{{ validationErrors.trigger_type }}</p>
              </div>
            </div>

            <!-- 定时计划设置 -->
            <div v-if="workflow.trigger_type === 'schedule'" class="sm:col-span-3">
              <label for="schedule" class="block text-sm font-medium text-gray-700 dark:text-gray-300">
                计划表达式 (Cron)
              </label>
              <div class="mt-1">
                <input
                  type="text"
                  id="schedule"
                  v-model="workflow.trigger_config.cron"
                  :class="{'border-red-300 focus:ring-red-500 focus:border-red-500': validationErrors.trigger_config?.cron}"
                  class="shadow-sm focus:ring-primary-500 focus:border-primary-500 block w-full sm:text-sm border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white rounded-md"
                  placeholder="例如: 0 0 * * * (每天零点)"
                />
                <p v-if="validationErrors.trigger_config?.cron" class="mt-2 text-sm text-red-600">{{ validationErrors.trigger_config.cron }}</p>
              </div>
              <p class="mt-2 text-sm text-gray-500 dark:text-gray-400">使用 cron 表达式设置定时计划</p>
            </div>

            <!-- Webhook 设置 -->
            <div v-if="workflow.trigger_type === 'webhook'" class="sm:col-span-3">
              <label for="webhook_secret" class="block text-sm font-medium text-gray-700 dark:text-gray-300">
                Webhook 密钥 (可选)
              </label>
              <div class="mt-1">
                <input
                  type="text"
                  id="webhook_secret"
                  v-model="workflow.trigger_config.secret"
                  class="shadow-sm focus:ring-primary-500 focus:border-primary-500 block w-full sm:text-sm border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white rounded-md"
                  placeholder="webhook 安全密钥"
                />
              </div>
              <p class="mt-2 text-sm text-gray-500 dark:text-gray-400">用于验证 webhook 请求的密钥</p>
            </div>

            <!-- 事件触发设置 -->
            <div v-if="workflow.trigger_type === 'event'" class="sm:col-span-3">
              <label for="event_type" class="block text-sm font-medium text-gray-700 dark:text-gray-300">
                事件类型
              </label>
              <div class="mt-1">
                <select
                  id="event_type"
                  v-model="workflow.trigger_config.event_type"
                  :class="{'border-red-300 focus:ring-red-500 focus:border-red-500': validationErrors.trigger_config?.event_type}"
                  class="shadow-sm focus:ring-primary-500 focus:border-primary-500 block w-full sm:text-sm border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white rounded-md"
                >
                  <option value="">选择事件类型</option>
                  <option value="file_upload">文件上传</option>
                  <option value="database_change">数据库变更</option>
                  <option value="api_call">API 调用</option>
                  <option value="system_event">系统事件</option>
                </select>
                <p v-if="validationErrors.trigger_config?.event_type" class="mt-2 text-sm text-red-600">{{ validationErrors.trigger_config.event_type }}</p>
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
                  v-model="workflow.status"
                  class="shadow-sm focus:ring-primary-500 focus:border-primary-500 block w-full sm:text-sm border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white rounded-md"
                >
                  <option value="draft">草稿</option>
                  <option value="active">活跃</option>
                  <option value="inactive">暂停</option>
                </select>
              </div>
            </div>

            <!-- 超时设置 -->
            <div class="sm:col-span-3">
              <label for="timeout" class="block text-sm font-medium text-gray-700 dark:text-gray-300">
                执行超时（秒）
              </label>
              <div class="mt-1">
                <input
                  type="number"
                  id="timeout"
                  v-model="workflow.timeout"
                  min="0"
                  class="shadow-sm focus:ring-primary-500 focus:border-primary-500 block w-full sm:text-sm border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white rounded-md"
                  placeholder="工作流执行超时时间（秒）"
                />
              </div>
              <p class="mt-2 text-sm text-gray-500 dark:text-gray-400">0 表示无超时限制</p>
            </div>

            <!-- 重试设置 -->
            <div class="sm:col-span-3">
              <label for="max_retries" class="block text-sm font-medium text-gray-700 dark:text-gray-300">
                最大重试次数
              </label>
              <div class="mt-1">
                <input
                  type="number"
                  id="max_retries"
                  v-model="workflow.max_retries"
                  min="0"
                  class="shadow-sm focus:ring-primary-500 focus:border-primary-500 block w-full sm:text-sm border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white rounded-md"
                  placeholder="失败时重试次数"
                />
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
                    v-model="workflowDefinition"
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
                @click="$router.back()"
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
</template>

<script setup>
import { ref, reactive, computed, watch } from 'vue';
import { useRouter } from 'vue-router';
import { useToast } from 'vue-toastification';
import { useWorkflowStore } from '../../stores/workflow';

const router = useRouter();
const toast = useToast();
const workflowStore = useWorkflowStore();
const saving = ref(false);
const validationErrors = ref({});
const workflowDefinition = ref('{\n  "steps": []\n}');

// 工作流数据
const workflow = reactive({
  name: '',
  description: '',
  trigger_type: '',
  trigger_config: {},
  status: 'draft',
  timeout: 300,
  max_retries: 3,
  tasks: []
});

// 监听工作流定义变化，更新tasks
watch(workflowDefinition, (newValue) => {
  try {
    const parsed = JSON.parse(newValue);
    if (parsed && parsed.steps) {
      workflow.tasks = parsed.steps;
    }
  } catch (e) {
    // 解析错误时不更新tasks
  }
});

// 触发器类型变更
const onTriggerTypeChange = () => {
  workflow.trigger_config = {};
  
  if (workflow.trigger_type === 'schedule') {
    workflow.trigger_config.cron = '0 0 * * *'; // 默认每天零点
  } else if (workflow.trigger_type === 'webhook') {
    workflow.trigger_config.secret = '';
  } else if (workflow.trigger_type === 'event') {
    workflow.trigger_config.event_type = '';
  }
};

// 格式化JSON
const formatJSON = () => {
  try {
    const parsed = JSON.parse(workflowDefinition.value);
    workflowDefinition.value = JSON.stringify(parsed, null, 2);
    validationErrors.value.definition = '';
  } catch (error) {
    validationErrors.value.definition = `JSON格式错误: ${error.message}`;
  }
};

// 表单验证
const validateForm = () => {
  const errors = {};
  
  if (!workflow.name.trim()) {
    errors.name = '工作流名称不能为空';
  }
  
  if (!workflow.description.trim()) {
    errors.description = '工作流描述不能为空';
  }
  
  if (!workflow.trigger_type) {
    errors.trigger_type = '请选择触发器类型';
  } else {
    if (workflow.trigger_type === 'schedule' && !workflow.trigger_config.cron) {
      if (!errors.trigger_config) errors.trigger_config = {};
      errors.trigger_config.cron = '请设置计划表达式';
    }
    
    if (workflow.trigger_type === 'event' && !workflow.trigger_config.event_type) {
      if (!errors.trigger_config) errors.trigger_config = {};
      errors.trigger_config.event_type = '请选择事件类型';
    }
  }
  
  // 验证JSON
  try {
    JSON.parse(workflowDefinition.value);
  } catch (error) {
    errors.definition = `JSON格式错误: ${error.message}`;
  }
  
  validationErrors.value = errors;
  return Object.keys(errors).length === 0;
};

// 准备API提交数据
const prepareWorkflowData = () => {
  // 解析JSON定义
  let definition;
  try {
    definition = JSON.parse(workflowDefinition.value);
  } catch (error) {
    throw new Error(`JSON格式错误: ${error.message}`);
  }
  
  // 构建API数据结构
  return {
    name: workflow.name,
    description: workflow.description,
    status: workflow.status,
    trigger_type: workflow.trigger_type,
    trigger_config: workflow.trigger_config,
    timeout: parseInt(workflow.timeout) || 0,
    max_retries: parseInt(workflow.max_retries) || 0,
    definition: definition
  };
};

// 保存工作流
const saveWorkflow = async () => {
  if (!validateForm()) {
    toast.error('请修正表单中的错误');
    return;
  }
  
  saving.value = true;
  
  try {
    const workflowData = prepareWorkflowData();
    
    // 使用store创建工作流
    await workflowStore.createWorkflow(workflowData);
    
    toast.success('工作流创建成功');
    
    // 重定向到工作流列表页
    router.push('/workflows');
  } catch (error) {
    console.error('创建工作流失败', error);
    
    // 处理特定错误
    if (error.response && error.response.status === 400) {
      // 服务器返回的验证错误
      const serverErrors = error.response.data.errors || {};
      Object.keys(serverErrors).forEach(key => {
        validationErrors.value[key] = serverErrors[key];
      });
      toast.error('表单验证失败，请检查输入');
    } else {
      toast.error('创建工作流失败: ' + (error.message || '未知错误'));
    }
  } finally {
    saving.value = false;
  }
};
</script>

<style scoped>
/* 添加一些额外的样式 */
.font-mono {
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", "Courier New", monospace;
}
</style> 