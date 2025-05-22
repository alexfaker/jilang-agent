<template>
  <div class="workflow-page">
    <h1 class="page-title">工作流管理</h1>
    
    <div class="actions-bar">
      <button class="btn btn-primary" @click="showCreateModal = true">
        <i class="fas fa-plus"></i> 创建工作流
      </button>
      <div class="filters">
        <select v-model="statusFilter" class="form-select">
          <option value="">所有状态</option>
          <option value="active">激活</option>
          <option value="inactive">未激活</option>
          <option value="draft">草稿</option>
        </select>
        <input 
          type="text" 
          v-model="searchQuery" 
          placeholder="搜索工作流..." 
          class="form-control search-input"
        />
      </div>
    </div>
    
    <div v-if="loading" class="loading-state">
      <div class="spinner"></div>
      <p>加载工作流...</p>
    </div>
    
    <div v-else-if="error" class="error-state">
      <i class="fas fa-exclamation-circle"></i>
      <p>{{ error }}</p>
      <button @click="fetchWorkflows" class="btn btn-secondary">重试</button>
    </div>
    
    <div v-else-if="filteredWorkflows.length === 0" class="empty-state">
      <i class="fas fa-file-alt"></i>
      <p>没有找到工作流</p>
      <p v-if="searchQuery || statusFilter" class="empty-hint">
        尝试清除过滤条件或创建新工作流
      </p>
      <button v-else @click="showCreateModal = true" class="btn btn-primary">
        创建第一个工作流
      </button>
    </div>
    
    <div v-else class="workflow-grid">
      <div 
        v-for="workflow in filteredWorkflows" 
        :key="workflow.id"
        class="workflow-card"
        :class="{'active': workflow.status === 'active'}"
      >
        <div class="workflow-header">
          <span 
            class="status-badge"
            :class="workflow.status"
          >
            {{ statusLabels[workflow.status] || workflow.status }}
          </span>
          <div class="workflow-actions">
            <button class="action-btn" @click="editWorkflow(workflow)">
              <i class="fas fa-edit"></i>
            </button>
            <button class="action-btn" @click="executeWorkflow(workflow)">
              <i class="fas fa-play"></i>
            </button>
            <button class="action-btn" @click="confirmDelete(workflow)">
              <i class="fas fa-trash"></i>
            </button>
          </div>
        </div>
        
        <h3 class="workflow-title">{{ workflow.name }}</h3>
        <p class="workflow-description">{{ workflow.description || '无描述' }}</p>
        
        <div class="workflow-footer">
          <span class="created-time">
            创建于 {{ formatDate(workflow.created_at) }}
          </span>
          <span class="execution-count">
            <i class="fas fa-history"></i> {{ workflow.execution_count || 0 }}次执行
          </span>
        </div>
      </div>
    </div>
    
    <!-- 分页控件 -->
    <div class="pagination-controls" v-if="totalPages > 1">
      <button 
        :disabled="currentPage === 1" 
        @click="changePage(currentPage - 1)"
        class="btn btn-sm btn-secondary"
      >
        <i class="fas fa-chevron-left"></i> 上一页
      </button>
      
      <span class="page-info">{{ currentPage }} / {{ totalPages }}</span>
      
      <button 
        :disabled="currentPage === totalPages" 
        @click="changePage(currentPage + 1)"
        class="btn btn-sm btn-secondary"
      >
        下一页 <i class="fas fa-chevron-right"></i>
      </button>
    </div>
    
    <!-- 创建/编辑工作流模态框 -->
    <div v-if="showCreateModal || editingWorkflow" class="modal-backdrop">
      <div class="modal-container">
        <div class="modal-header">
          <h2>{{ editingWorkflow ? '编辑工作流' : '创建工作流' }}</h2>
          <button class="close-btn" @click="closeModal">
            <i class="fas fa-times"></i>
          </button>
        </div>
        
        <div class="modal-body">
          <div class="form-group">
            <label for="workflow-name">名称 *</label>
            <input 
              id="workflow-name"
              v-model="workflowForm.name" 
              type="text" 
              class="form-control"
              :class="{'is-invalid': validationErrors.name}"
            />
            <div v-if="validationErrors.name" class="invalid-feedback">
              {{ validationErrors.name }}
            </div>
          </div>
          
          <div class="form-group">
            <label for="workflow-description">描述</label>
            <textarea 
              id="workflow-description"
              v-model="workflowForm.description" 
              class="form-control"
              rows="3"
            ></textarea>
          </div>
          
          <div class="form-group">
            <label for="workflow-status">状态</label>
            <select 
              id="workflow-status"
              v-model="workflowForm.status" 
              class="form-select"
            >
              <option value="draft">草稿</option>
              <option value="active">激活</option>
              <option value="inactive">停用</option>
            </select>
          </div>
          
          <div class="form-group">
            <label for="workflow-definition">工作流定义 (JSON) *</label>
            <textarea 
              id="workflow-definition"
              v-model="workflowForm.definition" 
              class="form-control code-editor"
              rows="10"
              :class="{'is-invalid': validationErrors.definition}"
            ></textarea>
            <div v-if="validationErrors.definition" class="invalid-feedback">
              {{ validationErrors.definition }}
            </div>
          </div>
        </div>
        
        <div class="modal-footer">
          <button class="btn btn-secondary" @click="closeModal">
            取消
          </button>
          <button 
            class="btn btn-primary" 
            @click="saveWorkflow"
            :disabled="saveInProgress"
          >
            <span v-if="saveInProgress" class="spinner-sm"></span>
            {{ editingWorkflow ? '保存更改' : '创建工作流' }}
          </button>
        </div>
      </div>
    </div>
    
    <!-- 删除确认模态框 -->
    <div v-if="showDeleteConfirm" class="modal-backdrop">
      <div class="modal-container confirm-modal">
        <div class="modal-header">
          <h2>确认删除</h2>
          <button class="close-btn" @click="showDeleteConfirm = false">
            <i class="fas fa-times"></i>
          </button>
        </div>
        
        <div class="modal-body">
          <p>
            您确定要删除工作流 <strong>{{ workflowToDelete?.name }}</strong> 吗？
          </p>
          <p class="text-danger">
            此操作无法撤销，所有相关的执行历史也将被删除。
          </p>
        </div>
        
        <div class="modal-footer">
          <button class="btn btn-secondary" @click="showDeleteConfirm = false">
            取消
          </button>
          <button 
            class="btn btn-danger" 
            @click="deleteWorkflow"
            :disabled="deleteInProgress"
          >
            <span v-if="deleteInProgress" class="spinner-sm"></span>
            确认删除
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'WorkflowView',
  data() {
    return {
      workflows: [],
      loading: true,
      error: null,
      currentPage: 1,
      totalItems: 0,
      itemsPerPage: 12,
      searchQuery: '',
      statusFilter: '',
      showCreateModal: false,
      editingWorkflow: null,
      workflowForm: {
        name: '',
        description: '',
        status: 'draft',
        definition: '{}'
      },
      validationErrors: {},
      saveInProgress: false,
      showDeleteConfirm: false,
      workflowToDelete: null,
      deleteInProgress: false,
      statusLabels: {
        active: '激活',
        inactive: '未激活',
        draft: '草稿'
      }
    };
  },
  computed: {
    totalPages() {
      return Math.ceil(this.totalItems / this.itemsPerPage);
    },
    filteredWorkflows() {
      let result = [...this.workflows];
      
      if (this.searchQuery) {
        const query = this.searchQuery.toLowerCase();
        result = result.filter(workflow => 
          workflow.name.toLowerCase().includes(query) || 
          (workflow.description && workflow.description.toLowerCase().includes(query))
        );
      }
      
      if (this.statusFilter) {
        result = result.filter(workflow => workflow.status === this.statusFilter);
      }
      
      return result;
    }
  },
  mounted() {
    this.fetchWorkflows();
  },
  methods: {
    fetchWorkflows() {
      this.loading = true;
      this.error = null;
      
      // 模拟API调用
      setTimeout(() => {
        try {
          // 这里替换为实际的API调用
          this.workflows = [
            {
              id: 1,
              name: '文档分析工作流',
              description: '使用AI代理分析文档内容并提取关键信息',
              status: 'active',
              created_at: '2023-10-15T10:20:30Z',
              execution_count: 24,
              definition: JSON.stringify({
                steps: [
                  { type: 'document_input', id: 'input' },
                  { type: 'text_analysis', id: 'analysis', input: 'input' },
                  { type: 'summary_generation', id: 'summary', input: 'analysis' }
                ]
              }, null, 2)
            },
            {
              id: 2,
              name: '客户支持代理',
              description: '处理客户查询并提供自动响应',
              status: 'inactive',
              created_at: '2023-09-28T14:15:20Z',
              execution_count: 124,
              definition: JSON.stringify({
                steps: [
                  { type: 'query_input', id: 'query' },
                  { type: 'intent_classification', id: 'intent', input: 'query' },
                  { type: 'response_generation', id: 'response', input: 'intent' }
                ]
              }, null, 2)
            },
            {
              id: 3,
              name: '代码审查助手',
              description: '分析代码并提供改进建议',
              status: 'draft',
              created_at: '2023-10-10T16:40:15Z',
              execution_count: 5,
              definition: JSON.stringify({
                steps: [
                  { type: 'code_input', id: 'code' },
                  { type: 'static_analysis', id: 'analysis', input: 'code' },
                  { type: 'improvement_suggestions', id: 'suggestions', input: 'analysis' }
                ]
              }, null, 2)
            }
          ];
          
          this.totalItems = this.workflows.length;
          this.loading = false;
        } catch (err) {
          this.error = '加载工作流失败：' + (err.message || '未知错误');
          this.loading = false;
        }
      }, 1000);
    },
    changePage(page) {
      if (page >= 1 && page <= this.totalPages) {
        this.currentPage = page;
        this.fetchWorkflows();
      }
    },
    editWorkflow(workflow) {
      this.editingWorkflow = workflow;
      this.workflowForm = {
        name: workflow.name,
        description: workflow.description || '',
        status: workflow.status,
        definition: workflow.definition
      };
      this.validationErrors = {};
    },
    closeModal() {
      this.showCreateModal = false;
      this.editingWorkflow = null;
      this.workflowForm = {
        name: '',
        description: '',
        status: 'draft',
        definition: '{}'
      };
      this.validationErrors = {};
    },
    saveWorkflow() {
      // 验证表单
      this.validationErrors = {};
      
      if (!this.workflowForm.name.trim()) {
        this.validationErrors.name = '工作流名称不能为空';
      }
      
      try {
        JSON.parse(this.workflowForm.definition);
      } catch (err) {
        this.validationErrors.definition = 'JSON格式无效：' + err.message;
      }
      
      if (Object.keys(this.validationErrors).length > 0) {
        return;
      }
      
      this.saveInProgress = true;
      
      // 模拟API调用
      setTimeout(() => {
        try {
          // 这里替换为实际的API调用
          if (this.editingWorkflow) {
            // 更新现有工作流
            const index = this.workflows.findIndex(w => w.id === this.editingWorkflow.id);
            if (index !== -1) {
              this.workflows[index] = {
                ...this.editingWorkflow,
                name: this.workflowForm.name,
                description: this.workflowForm.description,
                status: this.workflowForm.status,
                definition: this.workflowForm.definition
              };
            }
          } else {
            // 创建新工作流
            const newWorkflow = {
              id: Date.now(), // 模拟ID生成
              name: this.workflowForm.name,
              description: this.workflowForm.description,
              status: this.workflowForm.status,
              created_at: new Date().toISOString(),
              execution_count: 0,
              definition: this.workflowForm.definition
            };
            
            this.workflows.unshift(newWorkflow);
            this.totalItems++;
          }
          
          this.closeModal();
          this.saveInProgress = false;
        } catch (err) {
          alert('保存工作流失败：' + (err.message || '未知错误'));
          this.saveInProgress = false;
        }
      }, 1000);
    },
    executeWorkflow(workflow) {
      // 跳转到执行页面或启动执行
      alert(`开始执行工作流：${workflow.name}`);
      // this.$router.push(`/workflow/execute/${workflow.id}`);
    },
    confirmDelete(workflow) {
      this.workflowToDelete = workflow;
      this.showDeleteConfirm = true;
    },
    deleteWorkflow() {
      if (!this.workflowToDelete) return;
      
      this.deleteInProgress = true;
      
      // 模拟API调用
      setTimeout(() => {
        try {
          // 这里替换为实际的API调用
          const index = this.workflows.findIndex(w => w.id === this.workflowToDelete.id);
          if (index !== -1) {
            this.workflows.splice(index, 1);
            this.totalItems--;
          }
          
          this.showDeleteConfirm = false;
          this.workflowToDelete = null;
          this.deleteInProgress = false;
        } catch (err) {
          alert('删除工作流失败：' + (err.message || '未知错误'));
          this.deleteInProgress = false;
        }
      }, 1000);
    },
    formatDate(dateString) {
      const date = new Date(dateString);
      return date.toLocaleDateString('zh-CN', {
        year: 'numeric',
        month: 'short',
        day: 'numeric'
      });
    }
  }
};
</script>

<style scoped>
.workflow-page {
  padding: 1.5rem;
}

.page-title {
  margin-bottom: 1.5rem;
  font-size: 1.8rem;
  font-weight: 600;
}

.actions-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}

.filters {
  display: flex;
  gap: 0.5rem;
}

.search-input {
  width: 250px;
}

.workflow-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.workflow-card {
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  padding: 1.25rem;
  transition: transform 0.2s, box-shadow 0.2s;
  display: flex;
  flex-direction: column;
  border-left: 4px solid #ddd;
}

.workflow-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.workflow-card.active {
  border-left-color: #4CAF50;
}

.workflow-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.75rem;
}

.status-badge {
  font-size: 0.75rem;
  padding: 0.25rem 0.5rem;
  border-radius: 12px;
  background-color: #eee;
}

.status-badge.active {
  background-color: #e8f5e9;
  color: #2e7d32;
}

.status-badge.inactive {
  background-color: #f5f5f5;
  color: #757575;
}

.status-badge.draft {
  background-color: #e3f2fd;
  color: #1565c0;
}

.workflow-actions {
  display: flex;
  gap: 0.25rem;
}

.action-btn {
  width: 28px;
  height: 28px;
  border-radius: 4px;
  border: none;
  background-color: #f5f5f5;
  color: #555;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background-color 0.2s;
}

.action-btn:hover {
  background-color: #e0e0e0;
}

.workflow-title {
  font-size: 1.1rem;
  font-weight: 600;
  margin-bottom: 0.5rem;
}

.workflow-description {
  font-size: 0.9rem;
  color: #666;
  flex-grow: 1;
  margin-bottom: 1rem;
}

.workflow-footer {
  display: flex;
  justify-content: space-between;
  font-size: 0.8rem;
  color: #888;
}

.execution-count {
  display: flex;
  align-items: center;
  gap: 0.25rem;
}

.loading-state,
.error-state,
.empty-state {
  min-height: 200px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  text-align: center;
  padding: 2rem;
}

.spinner {
  border: 3px solid #f3f3f3;
  border-top: 3px solid #3498db;
  border-radius: 50%;
  width: 30px;
  height: 30px;
  animation: spin 1s linear infinite;
  margin-bottom: 1rem;
}

.spinner-sm {
  border: 2px solid #f3f3f3;
  border-top: 2px solid #fff;
  border-radius: 50%;
  width: 16px;
  height: 16px;
  animation: spin 1s linear infinite;
  display: inline-block;
  margin-right: 0.5rem;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.error-state i,
.empty-state i {
  font-size: 2.5rem;
  color: #ccc;
  margin-bottom: 1rem;
}

.error-state p {
  color: #e53935;
  margin-bottom: 1rem;
}

.empty-state p {
  margin-bottom: 1rem;
}

.empty-hint {
  font-size: 0.9rem;
  color: #777;
}

.pagination-controls {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 1rem;
  margin-top: 1rem;
}

.page-info {
  font-size: 0.9rem;
  color: #666;
}

.modal-backdrop {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal-container {
  background-color: #fff;
  border-radius: 8px;
  width: 90%;
  max-width: 700px;
  max-height: 90vh;
  display: flex;
  flex-direction: column;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.2);
}

.confirm-modal {
  max-width: 500px;
}

.modal-header {
  padding: 1rem 1.5rem;
  border-bottom: 1px solid #eee;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.modal-header h2 {
  margin: 0;
  font-size: 1.5rem;
  font-weight: 600;
}

.close-btn {
  background: none;
  border: none;
  font-size: 1.2rem;
  color: #666;
  cursor: pointer;
}

.modal-body {
  padding: 1.5rem;
  overflow-y: auto;
}

.modal-footer {
  padding: 1rem 1.5rem;
  border-top: 1px solid #eee;
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
}

.form-group {
  margin-bottom: 1.25rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 500;
}

.form-control,
.form-select {
  width: 100%;
  padding: 0.5rem 0.75rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 1rem;
}

.code-editor {
  font-family: monospace;
  background-color: #f8f9fa;
}

.text-danger {
  color: #dc3545;
}

.is-invalid {
  border-color: #dc3545;
}

.invalid-feedback {
  color: #dc3545;
  font-size: 0.875rem;
  margin-top: 0.25rem;
}

.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 0.375rem 0.75rem;
  border-radius: 4px;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.2s, border-color 0.2s;
  border: 1px solid transparent;
}

.btn-sm {
  padding: 0.25rem 0.5rem;
  font-size: 0.875rem;
}

.btn-primary {
  background-color: #1976d2;
  color: #fff;
  border-color: #1976d2;
}

.btn-primary:hover {
  background-color: #1565c0;
  border-color: #1565c0;
}

.btn-secondary {
  background-color: #f5f5f5;
  color: #333;
  border-color: #ddd;
}

.btn-secondary:hover {
  background-color: #e0e0e0;
  border-color: #ccc;
}

.btn-danger {
  background-color: #d32f2f;
  color: #fff;
  border-color: #d32f2f;
}

.btn-danger:hover {
  background-color: #c62828;
  border-color: #c62828;
}

.btn[disabled] {
  opacity: 0.65;
  cursor: not-allowed;
}
</style> 