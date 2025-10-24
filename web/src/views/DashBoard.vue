<template>
  <AppHeader />
  <div class="dashboard-container">
    <!-- 欢迎区域 -->
    <div class="welcome-section">
      <div class="welcome-content">
        <h1 class="welcome-title">欢迎回来，{{ userInfo.username }}！</h1>
        <p class="welcome-subtitle">祝您使用愉快</p>
      </div>
      <div class="welcome-decoration">
        <el-icon :size="50" color="#409EFF">
          <Picture />
        </el-icon>
      </div>
    </div>

    <!-- 主要内容区域 -->
    <div class="main-content">
      <!-- 第一排：系统信息 -->
      <div class="first-row">
        <el-card class="system-info-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <el-icon :size="20"><InfoFilled /></el-icon>
              <span>系统信息</span>
            </div>
          </template>
          
          <div class="system-info">
            <div class="info-row">
              <el-icon :size="18" color="#67C23A"><Calendar /></el-icon>
              <span class="info-text">今天是 {{ currentDate }}</span>
            </div>
            <div class="info-row">
              <el-icon :size="18" color="#409EFF"><Clock /></el-icon>
              <span class="info-text">当前时间 {{ currentTime }}</span>
            </div>
          </div>
        </el-card>
      </div>

      <!-- 第二排：用户信息 -->
      <div class="second-row">
        <el-card class="user-info-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <el-icon :size="20"><User /></el-icon>
              <span>个人信息</span>
            </div>
          </template>
          
          <div class="user-details">
            <div class="info-item">
              <label class="info-label">用户名</label>
              <div class="info-value">{{ userInfo.username }}</div>
            </div>
            
            <div class="info-item">
              <label class="info-label">邮箱</label>
              <div class="info-value">{{ userInfo.email }}</div>
            </div>
            
            <div class="info-item">
              <label class="info-label">用户ID</label>
              <div class="info-value">{{ userInfo.id }}</div>
            </div>
            
            <div class="info-item">
              <label class="info-label">账户类型</label>
              <div class="info-value">
                <el-tag :type="userInfo.is_admin ? 'danger' : 'success'" size="small">
                  {{ userInfo.is_admin ? '管理员' : '普通用户' }}
                </el-tag>
              </div>
            </div>
            
            <div class="info-item">
              <label class="info-label">注册时间</label>
              <div class="info-value">{{ formatDateDisplay(userInfo.created_time) }}</div>
            </div>
            
            <div class="info-item">
              <label class="info-label">最后更新</label>
              <div class="info-value">{{ formatDateDisplay(userInfo.updated_time) }}</div>
            </div>
          </div>
        </el-card>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import AppHeader from '@/components/AppHeader.vue'
import { formatDate } from '@/utils/date'
import { 
  User, 
  Picture, 
  InfoFilled,
  Calendar,
  Clock
} from '@element-plus/icons-vue'

// 获取用户信息
const userInfo = computed(() => {
  const info = localStorage.getItem('userInfo')
  if (info) {
    return JSON.parse(info)
  }
  return {
    id: '',
    username: '未登录',
    email: '',
    is_admin: false,
    created_time: '',
    updated_time: ''
  }
})

// 当前日期
const currentDate = ref('')
const currentTime = ref('')

// 更新时间
const updateDateTime = () => {
  const now = new Date()
  currentDate.value = now.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    weekday: 'long'
  })
  currentTime.value = now.toLocaleTimeString('zh-CN', {
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  })
}

// 定时器ID
let timer: number | null = null

// 格式化日期
// 使用工具函数中的formatDate
const formatDateDisplay = (dateString: string) => {
  return formatDate(dateString);
};

onMounted(() => {
  updateDateTime()
  timer = window.setInterval(updateDateTime, 1000)
})

onUnmounted(() => {
  if (timer) {
    clearInterval(timer)
  }
})
</script>

<style scoped>
.dashboard-container {
  min-height: calc(100vh - 60px);
  background: linear-gradient(135deg, #f5f7fa 0%, #e4e7ed 100%);
  padding: 30px;
}

/* 欢迎区域 */
.welcome-section {
  max-width: 1130px;
  margin: 0 auto 25px auto;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 12px;
  padding: 25px 35px;
  color: white;
  display: flex;
  justify-content: space-between;
  align-items: center;
  box-shadow: 0 4px 16px rgba(102, 126, 234, 0.25);
}

.welcome-content {
  flex: 1;
}

.welcome-title {
  font-size: 24px;
  font-weight: 600;
  margin: 0 0 6px 0;
}

.welcome-subtitle {
  font-size: 14px;
  opacity: 0.9;
  margin: 0;
}

.welcome-decoration {
  opacity: 0.2;
}

/* 主要内容区域 */
.main-content {
  max-width: 1200px;
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  gap: 30px;
}

/* 卡片头部样式 */
.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  font-size: 16px;
  color: #303133;
}

/* 第一排：系统信息 */
.first-row {
  width: 100%;
}

/* 系统信息卡片 */
.system-info-card {
  border-radius: 12px;
  overflow: hidden;
}

.system-info {
  display: flex;
  justify-content: space-around;
  gap: 30px;
}

.info-row {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  padding: 20px;
  background-color: #f5f7fa;
  border-radius: 8px;
}

.info-text {
  font-size: 15px;
  color: #606266;
  font-weight: 500;
}

/* 第二排：用户信息 */
.second-row {
  width: 100%;
}

/* 用户信息卡片 */
.user-info-card {
  border-radius: 12px;
  overflow: hidden;
}

.user-details {
  padding: 10px 0;
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 20px;
}

.info-item {
  padding: 15px;
  background-color: #fafbfc;
  border-radius: 8px;
  transition: all 0.3s ease;
}

.info-item:hover {
  background-color: #f5f7fa;
  transform: translateY(-2px);
}

.info-label {
  display: block;
  font-size: 13px;
  color: #909399;
  margin-bottom: 8px;
}

.info-value {
  font-size: 15px;
  color: #303133;
  font-weight: 500;
}

/* 响应式设计 */
@media (max-width: 1024px) {
  .system-info {
    flex-direction: column;
  }
  
  .user-details {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .dashboard-container {
    padding: 20px;
  }
  
  .welcome-section {
    flex-direction: column;
    text-align: center;
    padding: 20px 15px;
  }
  
  .welcome-decoration {
    margin-top: 15px;
  }
  
  .welcome-title {
    font-size: 20px;
  }
  
  .welcome-subtitle {
    font-size: 13px;
  }
}
</style>