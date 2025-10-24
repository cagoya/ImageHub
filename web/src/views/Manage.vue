<template>
  <AppHeader />
  <div class="manage-container">
    <el-card class="manage-card">
      <template #header>
        <div class="card-header">
          <span>用户管理</span>
          <el-button type="primary" @click="refreshUsers">刷新</el-button>
        </div>
      </template>

      <!-- 搜索和筛选区域 -->
      <div class="search-section">
        <el-input
          v-model="searchKeyword"
          placeholder="搜索用户名或邮箱"
          style="width: 300px; margin-right: 20px;"
          @input="handleSearch"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
        
        <el-select v-model="statusFilter" placeholder="用户状态" style="width: 150px; margin-right: 20px;" @change="handleFilter">
          <el-option label="全部" value="all" />
          <el-option label="正常" value="active" />
          <el-option label="封禁" value="suspended" />
        </el-select>

        <el-select v-model="roleFilter" placeholder="用户角色" style="width: 150px;" @change="handleFilter">
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
      >
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="username" label="用户名" width="150" />
        <el-table-column prop="email" label="邮箱" width="200" />
        <el-table-column prop="is_admin" label="角色" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.is_admin ? 'danger' : 'primary'">
              {{ scope.row.is_admin ? '管理员' : '普通用户' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="is_suspended" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.is_suspended ? 'warning' : 'success'">
              {{ scope.row.is_suspended ? '封禁' : '正常' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_time" label="注册时间" width="180">
          <template #default="scope">
            {{ formatDateDisplay(scope.row.created_time) }}
          </template>
        </el-table-column>
        <el-table-column prop="updated_time" label="更新时间" width="180">
          <template #default="scope">
            {{ formatDateDisplay(scope.row.updated_time) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="300">
          <template #default="scope">
            <el-button
              v-if="!scope.row.is_admin"
              type="success"
              size="small"
              @click="grantAdminPermission(scope.row)"
            >
              授予管理员权限
            </el-button>
            <el-button
              v-if="scope.row.is_admin"
              type="warning"
              size="small"
              @click="revokeAdminPermission(scope.row)"
            >
              撤销管理员权限
            </el-button>
            <el-button
              v-if="!scope.row.is_suspended"
              type="danger"
              size="small"
              @click="showSuspendDialog(scope.row)"
            >
              封禁用户
            </el-button>
            <el-button
              v-if="scope.row.is_suspended"
              type="success"
              size="small"
              @click="unsuspendUser(scope.row)"
            >
              解封用户
            </el-button>
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
        />
      </div>
    </el-card>

    <!-- 封禁用户对话框 -->
    <el-dialog
      v-model="suspendDialogVisible"
      title="封禁用户"
      width="400px"
    >
      <el-form :model="suspendForm" label-width="100px">
        <el-form-item label="用户名">
          <el-input v-model="suspendForm.username" disabled />
        </el-form-item>
        <el-form-item label="封禁时长">
          <el-select v-model="suspendForm.expiration" placeholder="选择封禁时长">
            <el-option label="1天" :value="1" />
            <el-option label="7天" :value="7" />
            <el-option label="30天" :value="30" />
            <el-option label="永久" :value="0" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="suspendDialogVisible = false">取消</el-button>
        <el-button type="danger" @click="confirmSuspend">确认封禁</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script lang="ts" setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search } from '@element-plus/icons-vue'
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
// 使用工具函数中的formatDate
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
  padding: 20px;
  background-color: #f0f2f5;
  min-height: 100vh;
}

.manage-card {
  max-width: 1200px;
  margin: 0 auto;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.search-section {
  display: flex;
  align-items: center;
  margin-bottom: 20px;
}

.pagination-container {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}
</style>