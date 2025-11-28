<template>
  <AppHeader />
  <div class="manage-container">
    <el-card class="manage-card">
      <template #header>
        <div class="card-header">
          <div class="header-content">
            <h2 class="header-title">用户管理</h2>
            <p class="header-subtitle">管理系统所有用户账户和权限</p>
          </div>
          <el-button 
            type="primary" 
            @click="refreshUsers"
            class="refresh-btn"
          >
            <el-icon><Refresh /></el-icon>
            <span>刷新</span>
          </el-button>
        </div>
      </template>

      <!-- 搜索和筛选区域 -->
      <div class="search-section">
        <el-input
          v-model="searchKeyword"
          placeholder="搜索用户名或邮箱"
          class="search-input"
          @input="handleSearch"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
        
        <el-select 
          v-model="statusFilter" 
          placeholder="用户状态" 
          class="filter-select"
          @change="handleFilter"
        >
          <el-option label="全部" value="all" />
          <el-option label="正常" value="active" />
          <el-option label="封禁" value="suspended" />
        </el-select>

        <el-select 
          v-model="roleFilter" 
          placeholder="用户角色" 
          class="filter-select"
          @change="handleFilter"
        >
          <el-option label="全部" value="all" />
          <el-option label="普通用户" value="user" />
          <el-option label="管理员" value="admin" />
        </el-select>
      </div>

      <!-- 用户列表 -->
      <el-table
        :data="filteredUsers"
        v-loading="loading"
        style="width: 100%; margin-top: 20px;"
        stripe
        class="user-table"
      >
        <el-table-column prop="id" label="ID" width="80" align="center" />
        <el-table-column prop="username" label="用户名" width="150" />
        <el-table-column prop="email" label="邮箱" width="200" />
        <el-table-column prop="is_admin" label="角色" width="120" align="center">
          <template #default="scope">
            <el-tag 
              :type="scope.row.is_admin ? 'danger' : 'primary'" 
              class="role-tag"
              effect="light"
            >
              {{ scope.row.is_admin ? '管理员' : '普通用户' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="is_suspended" label="状态" width="120" align="center">
          <template #default="scope">
            <el-tag 
              :type="scope.row.is_suspended ? 'warning' : 'success'" 
              class="status-tag"
              effect="light"
            >
              {{ scope.row.is_suspended ? '封禁' : '正常' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_time" label="注册时间" width="180">
          <template #default="scope">
            <span class="time-text">{{ formatDateDisplay(scope.row.created_time) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="updated_time" label="更新时间" width="180">
          <template #default="scope">
            <span class="time-text">{{ formatDateDisplay(scope.row.updated_time) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="320" align="center">
          <template #default="scope">
            <div class="action-buttons">
              <el-button
                v-if="!scope.row.is_admin"
                type="success"
                size="small"
                @click="grantAdminPermission(scope.row)"
                class="action-btn"
              >
                <el-icon><Promotion /></el-icon>
                授权
              </el-button>
              <el-button
                v-if="scope.row.is_admin"
                type="warning"
                size="small"
                @click="revokeAdminPermission(scope.row)"
                class="action-btn"
              >
                <el-icon><Remove /></el-icon>
                撤销
              </el-button>
              <el-button
                v-if="!scope.row.is_suspended"
                type="danger"
                size="small"
                @click="showSuspendDialog(scope.row)"
                class="action-btn"
              >
                <el-icon><Lock /></el-icon>
                封禁
              </el-button>
              <el-button
                v-if="scope.row.is_suspended"
                type="success"
                size="small"
                @click="unsuspendUser(scope.row)"
                class="action-btn"
              >
                <el-icon><Unlock /></el-icon>
                解封
              </el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="totalUsers"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
          class="custom-pagination"
        />
      </div>
    </el-card>

    <!-- 封禁用户对话框 -->
    <el-dialog
      v-model="suspendDialogVisible"
      title="封禁用户"
      width="400px"
      class="suspend-dialog"
    >
      <div class="dialog-content">
        <el-form :model="suspendForm" label-width="100px">
          <el-form-item label="用户名">
            <el-input v-model="suspendForm.username" disabled class="disabled-input" />
          </el-form-item>
          <el-form-item label="封禁时长">
            <el-select v-model="suspendForm.expiration" placeholder="选择封禁时长" class="suspend-select">
              <el-option label="1天" :value="1" />
              <el-option label="7天" :value="7" />
              <el-option label="30天" :value="30" />
              <el-option label="永久" :value="0" />
            </el-select>
          </el-form-item>
        </el-form>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="suspendDialogVisible = false" class="cancel-btn">取消</el-button>
          <el-button type="danger" @click="confirmSuspend" class="confirm-suspend-btn">
            <el-icon><Lock /></el-icon>
            确认封禁
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script lang="ts" setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Refresh, Promotion, Remove, Lock, Unlock } from '@element-plus/icons-vue'
import { UserAPI } from '../api'
import AppHeader from '@/components/AppHeader.vue'
import type { User, UserListResponse } from '../type'
import { formatDate } from '@/utils/date'

// 响应式数据
const loading = ref(false)
const users = ref<User[]>([])
const searchKeyword = ref('')
const statusFilter = ref('all')
const roleFilter = ref('all')
const currentPage = ref(1)
const pageSize = ref(20)
const totalUsers = ref(0)

// 封禁对话框相关
const suspendDialogVisible = ref(false)
const suspendForm = ref({
  id: 0,
  username: '',
  expiration: 7
})

// 计算属性 - 过滤后的用户列表
const filteredUsers = computed(() => {
  let filtered = users.value

  // 搜索过滤
  if (searchKeyword.value) {
    const keyword = searchKeyword.value.toLowerCase()
    filtered = filtered.filter(user => 
      user.username.toLowerCase().includes(keyword) ||
      user.email.toLowerCase().includes(keyword)
    )
  }

  // 状态过滤
  if (statusFilter.value !== 'all') {
    const isSuspended = statusFilter.value === 'suspended'
    filtered = filtered.filter(user => user.is_suspended === isSuspended)
  }

  // 角色过滤
  if (roleFilter.value !== 'all') {
    const is_admin = roleFilter.value === 'admin'
    filtered = filtered.filter(user => user.is_admin === is_admin)
  }

  return filtered
})

// 获取用户列表
const fetchUsers = async () => {
  try {
    loading.value = true
    const response = await UserAPI.getAllUsers(currentPage.value, pageSize.value)
    if (response.data.code === 200) {
      users.value = response.data.data.users || []
      totalUsers.value = response.data.data.totalCount || 0
    } else {
      ElMessage.error('获取用户列表失败')
    }
  } catch (error) {
    console.error('获取用户列表失败:', error)
    ElMessage.error('获取用户列表失败')
  } finally {
    loading.value = false
  }
}

// 刷新用户列表
const refreshUsers = () => {
  fetchUsers()
}

// 搜索处理
const handleSearch = () => {
  // 搜索是实时的，通过计算属性处理
}

// 筛选处理
const handleFilter = () => {
  // 筛选是实时的，通过计算属性处理
}

// 分页处理
const handleSizeChange = (val: number) => {
  pageSize.value = val
  currentPage.value = 1
  fetchUsers()
}

const handleCurrentChange = (val: number) => {
  currentPage.value = val
  fetchUsers()
}

// 授予管理员权限
const grantAdminPermission = async (user: User) => {
  try {
    await ElMessageBox.confirm(
      `确定要授予用户 "${user.username}" 管理员权限吗？`,
      '确认操作',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )

    const response = await UserAPI.grantAdminPermission(user.id)
    if (response.data.code === 200) {
      ElMessage.success('授予管理员权限成功')
      fetchUsers()
    } else {
      ElMessage.error(response.data.message || '授予管理员权限失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('授予管理员权限失败:', error)
      ElMessage.error('授予管理员权限失败')
    }
  }
}

// 撤销管理员权限
const revokeAdminPermission = async (user: User) => {
  try {
    await ElMessageBox.confirm(
      `确定要撤销用户 "${user.username}" 的管理员权限吗？`,
      '确认操作',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )

    const response = await UserAPI.revokeAdminPermission(user.id)
    if (response.data.code === 200) {
      ElMessage.success('撤销管理员权限成功')
      fetchUsers()
    } else {
      ElMessage.error(response.data.message || '撤销管理员权限失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('撤销管理员权限失败:', error)
      ElMessage.error('撤销管理员权限失败')
    }
  }
}

// 显示封禁对话框
const showSuspendDialog = (user: User) => {
  suspendForm.value = {
    id: user.id,
    username: user.username,
    expiration: 7
  }
  suspendDialogVisible.value = true
}

// 确认封禁用户
const confirmSuspend = async () => {
  try {
    const response = await UserAPI.suspendUser(suspendForm.value.id, suspendForm.value.expiration)
    if (response.data.code === 200) {
      ElMessage.success('封禁用户成功')
      suspendDialogVisible.value = false
      fetchUsers()
    } else {
      ElMessage.error(response.data.message || '封禁用户失败')
    }
  } catch (error) {
    console.error('封禁用户失败:', error)
    ElMessage.error('封禁用户失败')
  }
}

// 解封用户
const unsuspendUser = async (user: User) => {
  try {
    await ElMessageBox.confirm(
      `确定要解封用户 "${user.username}" 吗？`,
      '确认操作',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )

    const response = await UserAPI.unsuspendUser(user.id)
    if (response.data.code === 200) {
      ElMessage.success('解封用户成功')
      fetchUsers()
    } else {
      ElMessage.error(response.data.message || '解封用户失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('解封用户失败:', error)
      ElMessage.error('解封用户失败')
    }
  }
}

// 格式化日期
const formatDateDisplay = (dateString: string) => {
  return formatDate(dateString);
};

// 组件挂载时获取用户列表
onMounted(() => {
  fetchUsers()
})
</script>

<style scoped>
.manage-container {
  min-height: 100vh;
  padding: 30px;
  background: linear-gradient(135deg, #f0f8ff 0%, #e6f3ff 50%, #dcedff 100%);
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
}

.manage-card {
  max-width: 1400px;
  margin: 0 auto;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 20px;
  box-shadow: 
    0 8px 32px rgba(64, 158, 255, 0.15),
    0 2px 8px rgba(64, 158, 255, 0.1);
  border: 1px solid rgba(64, 158, 255, 0.1);
  overflow: hidden;
}

:deep(.manage-card .el-card__header) {
  background: linear-gradient(135deg, #f8fbff, #f0f7ff);
  border-bottom: 1px solid rgba(64, 158, 255, 0.1);
  padding: 25px 30px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-content {
  flex: 1;
}

.header-title {
  color: #2c3e50;
  font-size: 28px;
  font-weight: 700;
  margin: 0 0 5px 0;
  background: linear-gradient(135deg, #409EFF, #66B3FF);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.header-subtitle {
  color: #7f8c8d;
  font-size: 14px;
  margin: 0;
  font-weight: 400;
}

.refresh-btn {
  background: linear-gradient(135deg, #409EFF, #66B3FF);
  border: none;
  border-radius: 12px;
  padding: 10px 20px;
  font-weight: 600;
  transition: all 0.3s ease;
  box-shadow: 0 4px 12px rgba(64, 158, 255, 0.3);
  display: flex;
  align-items: center;
  gap: 8px;
}

.refresh-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(64, 158, 255, 0.4);
  background: linear-gradient(135deg, #66B3FF, #409EFF);
}

/* 搜索区域样式 */
.search-section {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 20px;
  padding: 0 10px;
}

.search-input {
  width: 300px;
}

:deep(.search-input .el-input__wrapper) {
  background: #ffffff;
  border: 2px solid #e6f3ff;
  border-radius: 12px;
  transition: all 0.3s ease;
  box-shadow: none;
}

:deep(.search-input .el-input__wrapper:hover),
:deep(.search-input .el-input__wrapper.is-focus) {
  border-color: #409EFF;
  box-shadow: 0 0 0 3px rgba(64, 158, 255, 0.1);
}

.filter-select {
  width: 150px;
}

:deep(.filter-select .el-input__wrapper) {
  background: #ffffff;
  border: 2px solid #e6f3ff;
  border-radius: 12px;
  transition: all 0.3s ease;
  box-shadow: none;
}

:deep(.filter-select .el-input__wrapper:hover),
:deep(.filter-select .el-input__wrapper.is-focus) {
  border-color: #409EFF;
  box-shadow: 0 0 0 3px rgba(64, 158, 255, 0.1);
}

/* 表格样式 */
.user-table {
  margin-top: 20px;
}

:deep(.user-table .el-table) {
  border-radius: 12px;
  overflow: hidden;
}

:deep(.user-table .el-table__header) {
  background: linear-gradient(135deg, #f8fbff, #f0f7ff);
}

:deep(.user-table .el-table__header th) {
  background: transparent;
  color: #409EFF;
  font-weight: 600;
  border-bottom: 2px solid #e6f3ff;
}

:deep(.user-table .el-table__row) {
  transition: all 0.3s ease;
}

:deep(.user-table .el-table__row:hover) {
  background-color: #f8fbff !important;
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(64, 158, 255, 0.1);
}

.role-tag,
.status-tag {
  font-weight: 600;
  border: none;
  border-radius: 8px;
}

.time-text {
  color: #5D6D7E;
  font-size: 13px;
  font-weight: 500;
}

/* 操作按钮样式 */
.action-buttons {
  display: flex;
  gap: 8px;
  justify-content: center;
  flex-wrap: wrap;
}

.action-btn {
  border-radius: 8px;
  font-weight: 500;
  display: flex;
  align-items: center;
  gap: 4px;
  transition: all 0.3s ease;
}

.action-btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

/* 分页样式 */
.pagination-container {
  display: flex;
  justify-content: center;
  margin-top: 30px;
  padding: 20px 0;
}

:deep(.custom-pagination .el-pagination.is-background .btn-prev),
:deep(.custom-pagination .el-pagination.is-background .btn-next),
:deep(.custom-pagination .el-pagination.is-background .el-pager li) {
  background: #f8fbff;
  border: 1px solid #e6f3ff;
  border-radius: 8px;
  margin: 0 4px;
  transition: all 0.3s ease;
}

:deep(.custom-pagination .el-pagination.is-background .btn-prev:hover),
:deep(.custom-pagination .el-pagination.is-background .btn-next:hover),
:deep(.custom-pagination .el-pagination.is-background .el-pager li:hover) {
  border-color: #409EFF;
  background: rgba(64, 158, 255, 0.1);
}

:deep(.custom-pagination .el-pagination.is-background .el-pager li.is-active) {
  background: linear-gradient(135deg, #409EFF, #66B3FF);
  border-color: #409EFF;
  color: white;
}

/* 对话框样式 */
:deep(.suspend-dialog .el-dialog) {
  border-radius: 16px;
  overflow: hidden;
}

:deep(.suspend-dialog .el-dialog__header) {
  background: linear-gradient(135deg, #f8fbff, #f0f7ff);
  border-bottom: 1px solid rgba(64, 158, 255, 0.1);
  margin: 0;
  padding: 20px;
}

:deep(.suspend-dialog .el-dialog__title) {
  color: #2c3e50;
  font-weight: 600;
  font-size: 18px;
}

.dialog-content {
  padding: 20px 0;
}

:deep(.disabled-input .el-input__wrapper) {
  background-color: #f8fbff;
  border-color: #e6f3ff;
}

.suspend-select {
  width: 100%;
}

:deep(.suspend-select .el-input__wrapper) {
  background: #ffffff;
  border: 2px solid #e6f3ff;
  border-radius: 12px;
  transition: all 0.3s ease;
}

:deep(.suspend-select .el-input__wrapper:hover),
:deep(.suspend-select .el-input__wrapper.is-focus) {
  border-color: #409EFF;
  box-shadow: 0 0 0 3px rgba(64, 158, 255, 0.1);
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

.cancel-btn {
  border: 2px solid #e6f3ff;
  border-radius: 12px;
  padding: 10px 20px;
  color: #409EFF;
  font-weight: 500;
  transition: all 0.3s ease;
}

.cancel-btn:hover {
  border-color: #409EFF;
  background: rgba(64, 158, 255, 0.1);
}

.confirm-suspend-btn {
  background: linear-gradient(135deg, #F56C6C, #F78989);
  border: none;
  border-radius: 12px;
  padding: 10px 20px;
  font-weight: 600;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 6px;
}

.confirm-suspend-btn:hover {
  background: linear-gradient(135deg, #F78989, #F56C6C);
  transform: translateY(-1px);
}

/* 响应式设计 */
@media (max-width: 1200px) {
  .manage-card {
    max-width: 95%;
  }
}

@media (max-width: 768px) {
  .manage-container {
    padding: 20px 15px;
  }

  .search-section {
    flex-direction: column;
    align-items: stretch;
    gap: 12px;
  }

  .search-input,
  .filter-select {
    width: 100%;
  }

  .card-header {
    flex-direction: column;
    gap: 15px;
    align-items: flex-start;
  }

  .header-title {
    font-size: 24px;
  }

  .action-buttons {
    justify-content: flex-start;
  }

  .action-btn {
    flex: 1;
    min-width: 80px;
  }
}

@media (max-width: 480px) {
  .manage-container {
    padding: 15px 10px;
  }

  :deep(.manage-card .el-card__header) {
    padding: 20px;
  }

  .header-title {
    font-size: 22px;
  }

  .action-buttons {
    flex-direction: column;
    gap: 8px;
  }

  .action-btn {
    width: 100%;
    justify-content: center;
  }
}
</style>