import { defineStore } from 'pinia';
import { workflowApi } from '../api';

export const useWorkflowStore = defineStore('workflow', {
  state: () => ({
    workflows: [],
    currentWorkflow: null,
    loading: false,
    error: null,
    pagination: {
      page: 1,
      limit: 10,
      total: 0
    },
    filters: {
      status: '',
      search: ''
    }
  }),
  
  getters: {
    getWorkflowById: (state) => (id) => {
      return state.workflows.find(workflow => workflow.id === id) || null;
    },
    
    filteredWorkflows: (state) => {
      let result = [...state.workflows];
      
      // 按状态筛选
      if (state.filters.status) {
        result = result.filter(workflow => workflow.status === state.filters.status);
      }
      
      // 按搜索词筛选
      if (state.filters.search) {
        const searchLower = state.filters.search.toLowerCase();
        result = result.filter(workflow => 
          workflow.name.toLowerCase().includes(searchLower) || 
          (workflow.description && workflow.description.toLowerCase().includes(searchLower))
        );
      }
      
      return result;
    },
    
    totalPages: (state) => {
      return Math.ceil(state.pagination.total / state.pagination.limit) || 1;
    }
  },
  
  actions: {
    async fetchWorkflows(params = {}) {
      this.loading = true;
      this.error = null;
      
      try {
        // 合并默认参数和传入的参数
        const queryParams = {
          limit: this.pagination.limit,
          offset: (this.pagination.page - 1) * this.pagination.limit,
          status: this.filters.status,
          ...params
        };
        
        const response = await workflowApi.getWorkflows(queryParams);
        
        // 假设API返回格式为 { data: [...], total: number }
        this.workflows = response.data || [];
        this.pagination.total = response.total || 0;
        
        return response;
      } catch (error) {
        console.error('获取工作流列表失败:', error);
        this.error = error.message || '获取工作流列表失败';
        throw error;
      } finally {
        this.loading = false;
      }
    },
    
    async fetchWorkflowById(id) {
      this.loading = true;
      this.error = null;
      
      try {
        const workflow = await workflowApi.getWorkflow(id);
        this.currentWorkflow = workflow;
        
        // 更新列表中的工作流
        const index = this.workflows.findIndex(w => w.id === id);
        if (index !== -1) {
          this.workflows[index] = workflow;
        }
        
        return workflow;
      } catch (error) {
        console.error(`获取工作流 #${id} 失败:`, error);
        this.error = error.message || '获取工作流详情失败';
        throw error;
      } finally {
        this.loading = false;
      }
    },
    
    async createWorkflow(workflowData) {
      this.loading = true;
      this.error = null;
      
      try {
        const newWorkflow = await workflowApi.createWorkflow(workflowData);
        
        // 添加到列表开头
        this.workflows.unshift(newWorkflow);
        this.pagination.total += 1;
        
        return newWorkflow;
      } catch (error) {
        console.error('创建工作流失败:', error);
        this.error = error.message || '创建工作流失败';
        throw error;
      } finally {
        this.loading = false;
      }
    },
    
    async updateWorkflow(id, workflowData) {
      this.loading = true;
      this.error = null;
      
      try {
        const updatedWorkflow = await workflowApi.updateWorkflow(id, workflowData);
        
        // 更新列表中的工作流
        const index = this.workflows.findIndex(w => w.id === id);
        if (index !== -1) {
          this.workflows[index] = updatedWorkflow;
        }
        
        // 如果当前查看的是这个工作流，也更新currentWorkflow
        if (this.currentWorkflow && this.currentWorkflow.id === id) {
          this.currentWorkflow = updatedWorkflow;
        }
        
        return updatedWorkflow;
      } catch (error) {
        console.error(`更新工作流 #${id} 失败:`, error);
        this.error = error.message || '更新工作流失败';
        throw error;
      } finally {
        this.loading = false;
      }
    },
    
    async deleteWorkflow(id) {
      this.loading = true;
      this.error = null;
      
      try {
        await workflowApi.deleteWorkflow(id);
        
        // 从列表中移除
        this.workflows = this.workflows.filter(w => w.id !== id);
        this.pagination.total -= 1;
        
        // 如果当前查看的是这个工作流，清除currentWorkflow
        if (this.currentWorkflow && this.currentWorkflow.id === id) {
          this.currentWorkflow = null;
        }
      } catch (error) {
        console.error(`删除工作流 #${id} 失败:`, error);
        this.error = error.message || '删除工作流失败';
        throw error;
      } finally {
        this.loading = false;
      }
    },
    
    async executeWorkflow(id, params = {}) {
      this.loading = true;
      this.error = null;
      
      try {
        const execution = await workflowApi.executeWorkflow(id, params);
        return execution;
      } catch (error) {
        console.error(`执行工作流 #${id} 失败:`, error);
        this.error = error.message || '执行工作流失败';
        throw error;
      } finally {
        this.loading = false;
      }
    },
    
    setPage(page) {
      if (page >= 1 && page <= this.totalPages) {
        this.pagination.page = page;
        this.fetchWorkflows();
      }
    },
    
    setLimit(limit) {
      this.pagination.limit = limit;
      this.pagination.page = 1; // 重置到第一页
      this.fetchWorkflows();
    },
    
    setFilter(filters) {
      this.filters = { ...this.filters, ...filters };
      this.pagination.page = 1; // 重置到第一页
      this.fetchWorkflows();
    },
    
    clearFilters() {
      this.filters = {
        status: '',
        search: ''
      };
      this.fetchWorkflows();
    }
  }
}); 