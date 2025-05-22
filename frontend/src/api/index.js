import axios from 'axios';
import { useRouter } from 'vue-router';

// 创建axios实例
const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
    'Accept': 'application/json'
  }
});

// 请求拦截器
api.interceptors.request.use(
  config => {
    // 添加认证令牌到请求头
    const token = localStorage.getItem('token');
    if (token) {
      config.headers['Authorization'] = `Bearer ${token}`;
    }
    return config;
  },
  error => {
    return Promise.reject(error);
  }
);

// 响应拦截器
api.interceptors.response.use(
  response => {
    return response.data;
  },
  error => {
    const router = useRouter();
    
    // 处理401未授权错误
    if (error.response && error.response.status === 401) {
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      
      // 如果当前页面不是登录页，重定向到登录页
      if (router && router.currentRoute.value.name !== 'login') {
        router.push({
          name: 'login',
          query: { redirect: router.currentRoute.value.fullPath }
        });
      }
    }
    
    // 构建错误消息
    let errorMessage = '请求失败';
    if (error.response && error.response.data) {
      errorMessage = error.response.data.message || errorMessage;
    } else if (error.message) {
      errorMessage = error.message;
    }
    
    return Promise.reject({
      message: errorMessage,
      status: error.response ? error.response.status : null,
      data: error.response ? error.response.data : null
    });
  }
);

// 认证相关API
export const authApi = {
  // 用户登录
  login(credentials) {
    return api.post('/auth/login', credentials);
  },
  
  // 用户注册
  register(userData) {
    return api.post('/auth/register', userData);
  },
  
  // 刷新令牌
  refreshToken() {
    return api.post('/auth/refresh');
  },
  
  // 退出登录
  logout() {
    return api.post('/auth/logout');
  }
};

// 用户相关API
export const userApi = {
  // 获取当前用户信息
  getCurrentUser() {
    return api.get('/users/profile');
  },
  
  // 更新用户信息
  updateProfile(data) {
    return api.put('/users/profile', data);
  },
  
  // 修改密码
  changePassword(data) {
    return api.post('/users/change-password', data);
  }
};

// 工作流相关API
export const workflowApi = {
  // 获取工作流列表
  getWorkflows(params) {
    return api.get('/workflows', { params });
  },
  
  // 获取单个工作流
  getWorkflow(id) {
    return api.get(`/workflows/${id}`);
  },
  
  // 创建工作流
  createWorkflow(data) {
    return api.post('/workflows', data);
  },
  
  // 更新工作流
  updateWorkflow(id, data) {
    return api.put(`/workflows/${id}`, data);
  },
  
  // 删除工作流
  deleteWorkflow(id) {
    return api.delete(`/workflows/${id}`);
  },
  
  // 执行工作流
  executeWorkflow(id, data) {
    return api.post(`/workflows/${id}/execute`, data);
  }
};

// 执行历史相关API
export const executionApi = {
  // 获取执行历史列表
  getExecutions(params) {
    return api.get('/executions', { params });
  },
  
  // 获取单个执行记录
  getExecution(id) {
    return api.get(`/executions/${id}`);
  }
};

// 代理相关API
export const agentApi = {
  // 获取代理列表
  getAgents(params) {
    return api.get('/agents', { params });
  },
  
  // 获取单个代理
  getAgent(id) {
    return api.get(`/agents/${id}`);
  }
};

// 统计数据相关API
export const statsApi = {
  // 获取仪表盘统计数据
  getDashboardStats() {
    return api.get('/stats/dashboard');
  }
};

export default api; 