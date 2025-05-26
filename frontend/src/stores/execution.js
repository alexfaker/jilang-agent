import { defineStore } from 'pinia';
import { executionApi } from '../api';

export const useExecutionStore = defineStore('execution', {
  state: () => ({
    executions: [],
    currentExecution: null,
    loading: false,
    error: null,
    pagination: {
      page: 1,
      limit: 10,
      total: 0
    },
    filters: {
      workflowId: null,
      status: '',
      startDate: '',
      endDate: ''
    }
  }),
  
  getters: {
    getExecutionById: (state) => (id) => {
      return state.executions.find(execution => execution.id === id) || null;
    },
    
    filteredExecutions: (state) => {
      let result = [...state.executions];
      
      // 按工作流ID筛选
      if (state.filters.workflowId) {
        result = result.filter(execution => execution.workflowId === state.filters.workflowId);
      }
      
      // 按状态筛选
      if (state.filters.status) {
        result = result.filter(execution => execution.status === state.filters.status);
      }
      
      // 按日期范围筛选
      if (state.filters.startDate) {
        const startDate = new Date(state.filters.startDate);
        result = result.filter(execution => new Date(execution.startedAt) >= startDate);
      }
      
      if (state.filters.endDate) {
        const endDate = new Date(state.filters.endDate);
        endDate.setHours(23, 59, 59, 999); // 设置为当天结束时间
        result = result.filter(execution => new Date(execution.startedAt) <= endDate);
      }
      
      return result;
    },
    
    totalPages: (state) => {
      return Math.ceil(state.pagination.total / state.pagination.limit) || 1;
    },
    
    isFiltered: (state) => {
      return !!state.filters.workflowId || 
             !!state.filters.status || 
             !!state.filters.startDate || 
             !!state.filters.endDate;
    }
  },
  
  actions: {
    async fetchExecutions(params = {}) {
      this.loading = true;
      this.error = null;
      
      try {
        // 合并默认参数和传入的参数
        const queryParams = {
          limit: this.pagination.limit,
          offset: (this.pagination.page - 1) * this.pagination.limit,
          workflow_id: this.filters.workflowId,
          status: this.filters.status,
          start_date: this.filters.startDate,
          end_date: this.filters.endDate,
          ...params
        };
        
        const response = await executionApi.getExecutions(queryParams);
        
        // 假设API返回格式为 { data: [...], total: number }
        this.executions = response.data || [];
        this.pagination.total = response.total || 0;
        
        return response;
      } catch (error) {
        console.error('获取执行历史失败:', error);
        this.error = error.message || '获取执行历史失败';
        throw error;
      } finally {
        this.loading = false;
      }
    },
    
    async fetchExecutionById(id) {
      this.loading = true;
      this.error = null;
      
      try {
        const execution = await executionApi.getExecution(id);
        this.currentExecution = execution;
        
        // 更新列表中的执行记录
        const index = this.executions.findIndex(e => e.id === id);
        if (index !== -1) {
          this.executions[index] = execution;
        }
        
        return execution;
      } catch (error) {
        console.error(`获取执行记录 #${id} 失败:`, error);
        this.error = error.message || '获取执行记录详情失败';
        throw error;
      } finally {
        this.loading = false;
      }
    },
    
    async cancelExecution(id) {
      this.loading = true;
      this.error = null;
      
      try {
        await executionApi.cancelExecution(id);
        
        // 更新列表中的执行记录状态
        const index = this.executions.findIndex(e => e.id === id);
        if (index !== -1) {
          this.executions[index] = {
            ...this.executions[index],
            status: 'cancelled',
            statusText: '已取消',
            statusClass: 'warning'
          };
        }
        
        // 如果当前查看的是这个执行记录，也更新currentExecution
        if (this.currentExecution && this.currentExecution.id === id) {
          this.currentExecution = {
            ...this.currentExecution,
            status: 'cancelled',
            statusText: '已取消',
            statusClass: 'warning'
          };
        }
      } catch (error) {
        console.error(`取消执行 #${id} 失败:`, error);
        this.error = error.message || '取消执行失败';
        throw error;
      } finally {
        this.loading = false;
      }
    },
    
    setPage(page) {
      if (page >= 1 && page <= this.totalPages) {
        this.pagination.page = page;
        this.fetchExecutions();
      }
    },
    
    setLimit(limit) {
      this.pagination.limit = limit;
      this.pagination.page = 1; // 重置到第一页
      this.fetchExecutions();
    },
    
    setFilter(filters) {
      this.filters = { ...this.filters, ...filters };
      this.pagination.page = 1; // 重置到第一页
      this.fetchExecutions();
    },
    
    clearFilters() {
      this.filters = {
        workflowId: null,
        status: '',
        startDate: '',
        endDate: ''
      };
      this.pagination.page = 1;
      this.fetchExecutions();
    },
    
    // 格式化状态
    formatStatus(status) {
      const statusMap = {
        'pending': { text: '等待中', class: 'info' },
        'running': { text: '运行中', class: 'info' },
        'success': { text: '成功', class: 'success' },
        'failed': { text: '失败', class: 'error' },
        'cancelled': { text: '已取消', class: 'warning' }
      };
      
      return statusMap[status] || { text: status, class: 'default' };
    },
    
    // 格式化日期时间
    formatDateTime(dateString) {
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
    },
    
    // 格式化持续时间
    formatDuration(seconds) {
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
    }
  }
}); 