<template>
  <div>
    <h2 class="text-lg font-medium text-gray-900 dark:text-gray-100">安全设置</h2>
    <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
      管理API密钥、访问令牌和活动会话
    </p>

    <!-- API密钥管理 -->
    <div class="mt-8">
      <h3 class="text-md font-medium text-gray-900 dark:text-gray-100">API密钥</h3>
      <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
        使用API密钥可以从外部应用程序和服务访问<BrandName size="sm" weight="normal" color="gray-500 dark:text-gray-400" />。请妥善保管您的密钥，它们具有与您账户相同的权限。
      </p>

      <div class="mt-4 flex">
        <input
          v-model="apiKeyDescription"
          type="text"
          placeholder="密钥描述（例如：测试应用）"
          class="shadow-sm focus:ring-primary-500 focus:border-primary-500 block w-full sm:text-sm border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white rounded-md mr-2"
          :class="{ 'border-red-300 dark:border-red-700': descriptionError }"
          @input="descriptionError = false"
        />
        <button
          type="button"
          @click="generateApiKey"
          :disabled="isGenerating || !apiKeyDescription"
          class="inline-flex items-center px-4 py-2 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          <span v-if="isGenerating">生成中...</span>
          <span v-else>生成新密钥</span>
        </button>
      </div>
      <div v-if="descriptionError" class="mt-1 text-sm text-red-600 dark:text-red-400">
        请输入密钥描述
      </div>

      <!-- 新密钥显示 -->
      <div
        v-if="newApiKey"
        class="mt-4 p-4 bg-gray-50 dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-md"
      >
        <div class="flex justify-between items-center">
          <h4 class="text-sm font-medium text-gray-900 dark:text-gray-100">
            新密钥生成成功
          </h4>
          <button
            type="button"
            @click="newApiKey = ''"
            class="text-gray-500 hover:text-gray-600 dark:text-gray-400 dark:hover:text-gray-300"
          >
            <span class="sr-only">关闭</span>
            <svg class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
              <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd" />
            </svg>
          </button>
        </div>
        <p class="mt-2 text-sm text-gray-500 dark:text-gray-400">
          请立即保存此API密钥。为安全起见，我们不会再次显示它。
        </p>
        <div class="mt-2 flex rounded-md shadow-sm">
          <input
            type="text"
            readonly
            :value="newApiKey"
            class="flex-1 min-w-0 block w-full px-3 py-2 text-sm border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white rounded-l-md focus:outline-none"
          />
          <button
            type="button"
            @click="copyApiKey"
            class="inline-flex items-center px-3 py-2 border border-l-0 border-gray-300 dark:border-gray-600 bg-gray-100 dark:bg-gray-600 text-gray-700 dark:text-gray-300 rounded-r-md hover:bg-gray-200 dark:hover:bg-gray-500 focus:outline-none focus:ring-1 focus:ring-primary-500"
          >
            <span v-if="copied">已复制!</span>
            <span v-else>复制</span>
          </button>
        </div>
      </div>

      <!-- API密钥列表 -->
      <div class="mt-6 bg-white dark:bg-gray-800 shadow rounded-md overflow-hidden">
        <ul role="list" class="divide-y divide-gray-200 dark:divide-gray-700">
          <li v-if="apiKeys.length === 0 && !isLoading" class="px-4 py-4 sm:px-6">
            <div class="text-sm text-gray-500 dark:text-gray-400">
              您还没有创建API密钥
            </div>
          </li>
          <li v-if="isLoading" class="px-4 py-4 sm:px-6">
            <div class="text-sm text-gray-500 dark:text-gray-400">
              加载中...
            </div>
          </li>
          <li v-for="key in apiKeys" :key="key.id" class="px-4 py-4 sm:px-6">
            <div class="flex items-center justify-between">
              <div>
                <div class="text-sm font-medium text-gray-900 dark:text-gray-100">
                  {{ key.description }}
                </div>
                <div class="mt-1 flex items-center text-sm text-gray-500 dark:text-gray-400">
                  <span>ID: {{ key.id }}</span>
                  <span class="mx-2">&middot;</span>
                  <span>创建于: {{ formatDate(key.createdAt) }}</span>
                  <span v-if="key.lastUsed" class="mx-2">&middot;</span>
                  <span v-if="key.lastUsed">最后使用: {{ formatDate(key.lastUsed) }}</span>
                </div>
              </div>
              <button
                type="button"
                @click="confirmRevoke(key)"
                class="inline-flex items-center px-3 py-1 border border-transparent text-sm leading-5 font-medium rounded-md text-white bg-red-600 hover:bg-red-700 focus:outline-none focus:border-red-700 focus:shadow-outline-red active:bg-red-700 transition ease-in-out duration-150"
              >
                撤销
              </button>
            </div>
          </li>
        </ul>
      </div>
    </div>

    <!-- 活动会话 -->
    <div class="mt-8">
      <h3 class="text-md font-medium text-gray-900 dark:text-gray-100">活动会话</h3>
      <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
        以下是您当前的活动会话。如果您发现可疑活动，可以终止其他会话。
      </p>

      <!-- 会话列表 -->
      <div class="mt-4 bg-white dark:bg-gray-800 shadow rounded-md overflow-hidden">
        <ul role="list" class="divide-y divide-gray-200 dark:divide-gray-700">
          <li v-if="sessions.length === 0 && !isLoadingSessions" class="px-4 py-4 sm:px-6">
            <div class="text-sm text-gray-500 dark:text-gray-400">
              没有其他活动会话
            </div>
          </li>
          <li v-if="isLoadingSessions" class="px-4 py-4 sm:px-6">
            <div class="text-sm text-gray-500 dark:text-gray-400">
              加载中...
            </div>
          </li>
          <li v-for="session in sessions" :key="session.id" class="px-4 py-4 sm:px-6">
            <div class="flex items-center justify-between">
              <div>
                <div class="flex items-center">
                  <div class="text-sm font-medium text-gray-900 dark:text-gray-100">
                    {{ session.deviceName }}
                    <span
                      v-if="session.isCurrent"
                      class="ml-2 px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-green-100 dark:bg-green-800 text-green-800 dark:text-green-100"
                    >
                      当前会话
                    </span>
                  </div>
                </div>
                <div class="mt-1 flex items-center text-sm text-gray-500 dark:text-gray-400">
                  <span>IP: {{ session.ip }}</span>
                  <span class="mx-2">&middot;</span>
                  <span>登录时间: {{ formatDate(session.loginTime) }}</span>
                  <span class="mx-2">&middot;</span>
                  <span>最后活动: {{ formatFromNow(session.lastActivity) }}</span>
                </div>
                <div class="mt-1 text-xs text-gray-500 dark:text-gray-400">
                  <span>{{ session.browser }}</span>
                  <span class="mx-1">&middot;</span>
                  <span>{{ session.os }}</span>
                </div>
              </div>
              <button
                v-if="!session.isCurrent"
                type="button"
                @click="confirmTerminateSession(session)"
                class="inline-flex items-center px-3 py-1 border border-gray-300 dark:border-gray-600 text-sm leading-5 font-medium rounded-md text-gray-700 dark:text-gray-200 bg-white dark:bg-gray-700 hover:bg-gray-50 dark:hover:bg-gray-600 focus:outline-none focus:border-primary-300 focus:shadow-outline-primary"
              >
                终止
              </button>
            </div>
          </li>
        </ul>
      </div>
    </div>

    <!-- 确认撤销API密钥对话框 -->
    <div v-if="showRevokeConfirm" class="fixed z-10 inset-0 overflow-y-auto">
      <div class="flex items-center justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
        <div class="fixed inset-0 transition-opacity" aria-hidden="true">
          <div class="absolute inset-0 bg-gray-500 dark:bg-gray-900 opacity-75"></div>
        </div>
        <span class="hidden sm:inline-block sm:align-middle sm:h-screen" aria-hidden="true">&#8203;</span>
        <div
          class="inline-block align-bottom bg-white dark:bg-gray-800 rounded-lg px-4 pt-5 pb-4 text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full sm:p-6"
          role="dialog"
          aria-modal="true"
          aria-labelledby="modal-headline"
        >
          <div>
            <div class="mt-3 text-center sm:mt-0 sm:ml-4 sm:text-left">
              <h3
                class="text-lg leading-6 font-medium text-gray-900 dark:text-gray-100"
                id="modal-headline"
              >
                撤销API密钥
              </h3>
              <div class="mt-2">
                <p class="text-sm text-gray-500 dark:text-gray-400">
                  您确定要撤销此API密钥吗？此操作无法撤销，使用此密钥的应用程序将无法再访问<BrandName size="sm" weight="normal" color="gray-500 dark:text-gray-400" /> API。
                </p>
              </div>
            </div>
          </div>
          <div class="mt-5 sm:mt-4 sm:flex sm:flex-row-reverse">
            <button
              type="button"
              @click="revokeApiKey"
              :disabled="isRevoking"
              class="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-red-600 text-base font-medium text-white hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500 sm:ml-3 sm:w-auto sm:text-sm disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <span v-if="isRevoking">撤销中...</span>
              <span v-else>撤销</span>
            </button>
            <button
              type="button"
              @click="cancelRevoke"
              class="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 dark:border-gray-600 shadow-sm px-4 py-2 bg-white dark:bg-gray-700 text-base font-medium text-gray-700 dark:text-gray-200 hover:bg-gray-50 dark:hover:bg-gray-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 sm:mt-0 sm:w-auto sm:text-sm"
            >
              取消
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 确认终止会话对话框 -->
    <div v-if="showTerminateConfirm" class="fixed z-10 inset-0 overflow-y-auto">
      <div class="flex items-center justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
        <div class="fixed inset-0 transition-opacity" aria-hidden="true">
          <div class="absolute inset-0 bg-gray-500 dark:bg-gray-900 opacity-75"></div>
        </div>
        <span class="hidden sm:inline-block sm:align-middle sm:h-screen" aria-hidden="true">&#8203;</span>
        <div
          class="inline-block align-bottom bg-white dark:bg-gray-800 rounded-lg px-4 pt-5 pb-4 text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full sm:p-6"
          role="dialog"
          aria-modal="true"
          aria-labelledby="modal-headline"
        >
          <div>
            <div class="mt-3 text-center sm:mt-0 sm:ml-4 sm:text-left">
              <h3
                class="text-lg leading-6 font-medium text-gray-900 dark:text-gray-100"
                id="modal-headline"
              >
                终止会话
              </h3>
              <div class="mt-2">
                <p class="text-sm text-gray-500 dark:text-gray-400">
                  您确定要终止此会话吗？该设备上的用户将被登出并需要重新登录。
                </p>
              </div>
            </div>
          </div>
          <div class="mt-5 sm:mt-4 sm:flex sm:flex-row-reverse">
            <button
              type="button"
              @click="terminateSession"
              :disabled="isTerminating"
              class="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-red-600 text-base font-medium text-white hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500 sm:ml-3 sm:w-auto sm:text-sm disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <span v-if="isTerminating">终止中...</span>
              <span v-else>终止</span>
            </button>
            <button
              type="button"
              @click="cancelTerminate"
              class="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 dark:border-gray-600 shadow-sm px-4 py-2 bg-white dark:bg-gray-700 text-base font-medium text-gray-700 dark:text-gray-200 hover:bg-gray-50 dark:hover:bg-gray-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 sm:mt-0 sm:w-auto sm:text-sm"
            >
              取消
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useToast } from 'vue-toastification';
import { format, formatDistance } from 'date-fns';
import { zhCN } from 'date-fns/locale';
import { settingsApi } from '../../api';
import BrandName from '../../components/common/BrandName.vue';

const toast = useToast();

// API密钥管理
const apiKeys = ref([]);
const isLoading = ref(false);
const isGenerating = ref(false);
const isRevoking = ref(false);
const apiKeyDescription = ref('');
const newApiKey = ref('');
const descriptionError = ref(false);
const showRevokeConfirm = ref(false);
const keyToRevoke = ref(null);
const copied = ref(false);

// 会话管理
const sessions = ref([]);
const isLoadingSessions = ref(false);
const isTerminating = ref(false);
const showTerminateConfirm = ref(false);
const sessionToTerminate = ref(null);

// 加载API密钥
const loadApiKeys = async () => {
  isLoading.value = true;
  
  try {
    const response = await settingsApi.getApiKeys();
    apiKeys.value = response.data || [];
  } catch (error) {
    console.error('加载API密钥失败:', error);
    toast.error('无法加载API密钥列表');
  } finally {
    isLoading.value = false;
  }
};

// 加载活动会话
const loadSessions = async () => {
  isLoadingSessions.value = true;
  
  try {
    const response = await settingsApi.getSecuritySettings();
    sessions.value = response.sessions || [];
  } catch (error) {
    console.error('加载会话信息失败:', error);
    toast.error('无法加载活动会话列表');
  } finally {
    isLoadingSessions.value = false;
  }
};

// 生成新API密钥
const generateApiKey = async () => {
  if (!apiKeyDescription.value.trim()) {
    descriptionError.value = true;
    return;
  }
  
  isGenerating.value = true;
  
  try {
    const response = await settingsApi.generateApiKey(apiKeyDescription.value.trim());
    newApiKey.value = response.key || '';
    toast.success('API密钥生成成功');
    
    // 清空描述并刷新列表
    apiKeyDescription.value = '';
    await loadApiKeys();
  } catch (error) {
    console.error('生成API密钥失败:', error);
    toast.error('无法生成API密钥');
  } finally {
    isGenerating.value = false;
  }
};

// 复制API密钥到剪贴板
const copyApiKey = () => {
  if (!newApiKey.value) return;
  
  navigator.clipboard.writeText(newApiKey.value).then(() => {
    copied.value = true;
    setTimeout(() => {
      copied.value = false;
    }, 2000);
  });
};

// 确认撤销API密钥
const confirmRevoke = (key) => {
  keyToRevoke.value = key;
  showRevokeConfirm.value = true;
};

// 取消撤销
const cancelRevoke = () => {
  keyToRevoke.value = null;
  showRevokeConfirm.value = false;
};

// 撤销API密钥
const revokeApiKey = async () => {
  if (!keyToRevoke.value) return;
  
  isRevoking.value = true;
  
  try {
    await settingsApi.revokeApiKey(keyToRevoke.value.id);
    toast.success('API密钥已撤销');
    
    // 关闭对话框并刷新列表
    showRevokeConfirm.value = false;
    keyToRevoke.value = null;
    await loadApiKeys();
  } catch (error) {
    console.error('撤销API密钥失败:', error);
    toast.error('无法撤销API密钥');
  } finally {
    isRevoking.value = false;
  }
};

// 确认终止会话
const confirmTerminateSession = (session) => {
  sessionToTerminate.value = session;
  showTerminateConfirm.value = true;
};

// 取消终止会话
const cancelTerminate = () => {
  sessionToTerminate.value = null;
  showTerminateConfirm.value = false;
};

// 终止会话
const terminateSession = async () => {
  if (!sessionToTerminate.value) return;
  
  isTerminating.value = true;
  
  try {
    // 此处假设API有终止会话的接口
    // await settingsApi.terminateSession(sessionToTerminate.value.id);
    
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 500));
    
    toast.success('会话已终止');
    
    // 关闭对话框并刷新列表
    showTerminateConfirm.value = false;
    sessionToTerminate.value = null;
    await loadSessions();
  } catch (error) {
    console.error('终止会话失败:', error);
    toast.error('无法终止会话');
  } finally {
    isTerminating.value = false;
  }
};

// 格式化日期
const formatDate = (dateString) => {
  try {
    const date = new Date(dateString);
    return format(date, 'yyyy-MM-dd HH:mm:ss', { locale: zhCN });
  } catch (error) {
    return dateString || '未知';
  }
};

// 格式化相对时间
const formatFromNow = (dateString) => {
  try {
    const date = new Date(dateString);
    return formatDistance(date, new Date(), { addSuffix: true, locale: zhCN });
  } catch (error) {
    return dateString || '未知';
  }
};

// 初始化
onMounted(() => {
  loadApiKeys();
  loadSessions();
});
</script> 