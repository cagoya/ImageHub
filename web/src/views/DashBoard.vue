<template>
  <AppHeader />
  <div class="dashboard-container">
    <div class="welcome-section">
      <div class="welcome-content">
        <h1 class="welcome-title">欢迎回来，{{ userInfo.username }}！</h1>
        <p class="welcome-subtitle">祝您使用愉快</p>
      </div>
      <div class="welcome-decoration">
        <div class="decoration-circle circle-1"></div>
        <div class="decoration-circle circle-2"></div>
        <el-icon :size="50" color="rgba(255,255,255,0.9)">
          <Picture />
        </el-icon>
      </div>
    </div>

    <div class="main-content">
      <div class="first-row">
        <el-card class="system-info-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <div class="card-icon">
                <el-icon :size="20"><InfoFilled /></el-icon>
              </div>
              <span>系统信息</span>
            </div>
          </template>
          
          <div class="system-info">
            <div class="info-row">
              <div class="info-icon">
                <el-icon :size="20" color="#409EFF"><Calendar /></el-icon>
              </div>
              <div class="info-content">
                <span class="info-label">当前日期</span>
                <span class="info-text">{{ currentDate }}</span>
              </div>
            </div>
            <div class="info-row">
              <div class="info-icon">
                <el-icon :size="20" color="#409EFF"><Clock /></el-icon>
              </div>
              <div class="info-content">
                <span class="info-label">当前时间</span>
                <span class="info-text">{{ currentTime }}</span>
              </div>
            </div>
          </div>
        </el-card>
      </div>

      <div class="second-row">
        <el-card class="user-info-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <div class="card-icon">
                <el-icon :size="20"><User /></el-icon>
              </div>
              <span>个人信息</span>
            </div>
          </template>
          
          <div class="user-details">
            <div class="info-item">
              <div class="info-icon-item">
                <el-icon :size="16" color="#409EFF"><User /></el-icon>
              </div>
              <div class="info-content-item">
                <label class="info-label">用户名</label>
                <div class="info-value">{{ userInfo.username }}</div>
              </div>
            </div>
            
            <div class="info-item">
              <div class="info-icon-item">
                <el-icon :size="16" color="#409EFF"><Message /></el-icon>
              </div>
              <div class="info-content-item">
                <label class="info-label">邮箱</label>
                <div class="info-value">{{ userInfo.email }}</div>
              </div>
            </div>
            
            <div class="info-item">
              <div class="info-icon-item">
                <el-icon :size="16" color="#409EFF"><Key /></el-icon>
              </div>
              <div class="info-content-item">
                <label class="info-label">用户ID</label>
                <div class="info-value">{{ userInfo.id }}</div>
              </div>
            </div>
            
            <div class="info-item">
              <div class="info-icon-item">
                <el-icon :size="16" color="#409EFF"><Star /></el-icon>
              </div>
              <div class="info-content-item">
                <label class="info-label">账户类型</label>
                <div class="info-value">
                  <el-tag 
                    :type="userInfo.is_admin ? 'danger' : 'success'" 
                    size="small"
                    class="role-tag"
                  >
                    {{ userInfo.is_admin ? '管理员' : '普通用户' }}
                  </el-tag>
                </div>
              </div>
            </div>
            
            <div class="info-item">
              <div class="info-icon-item">
                <el-icon :size="16" color="#409EFF"><Calendar /></el-icon>
              </div>
              <div class="info-content-item">
                <label class="info-label">注册时间</label>
                <div class="info-value">{{ formatDateDisplay(userInfo.created_time) }}</div>
              </div>
            </div>
            
            <div class="info-item">
              <div class="info-icon-item">
                <el-icon :size="16" color="#409EFF"><Refresh /></el-icon>
              </div>
              <div class="info-content-item">
                <label class="info-label">最后更新</label>
                <div class="info-value">{{ formatDateDisplay(userInfo.updated_time) }}</div>
              </div>
            </div>
          </div>
        </el-card>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
// ⚠️ 确保 AppHeader 组件路径和 formatDate 工具函数存在
import AppHeader from '@/components/AppHeader.vue'
import { formatDate } from '@/utils/date'
import { 
  User, 
  Picture, 
  InfoFilled,
  Calendar,
  Clock,
  Message,
  Key,
  Star,
  Refresh
} from '@element-plus/icons-vue'

// 获取用户信息
const userInfo = computed(() => {
  const info = localStorage.getItem('userInfo')
  if (info) {
    try {
      // 假设 userInfo 包含 username, email, id, is_admin, created_time, updated_time
      return JSON.parse(info)
    } catch (e) {
      console.error("解析用户信息失败:", e);
    }
  }
  return {
    id: '',
    username: '游客',
    email: 'N/A',
    is_admin: false,
    created_time: new Date().toISOString(),
    updated_time: new Date().toISOString()
  }
})

// 当前日期和时间
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
const formatDateDisplay = (dateString: string) => {
  // 假设 formatDate 函数能处理您的 ISO 时间字符串
  if (!dateString) return 'N/A';
  // return formatDate(dateString); 
  // 临时使用内置函数，直到您提供 formatDate 的实现
  const date = new Date(dateString);
  return date.toLocaleDateString('zh-CN') + ' ' + date.toLocaleTimeString('zh-CN');
};

onMounted(() => {
  updateDateTime()
  // 1000 毫秒更新一次时间
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
  background: linear-gradient(135deg, #f0f8ff 0%, #e6f3ff 50%, #dcedff 100%);
  padding: 30px;
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
  overflow-x: hidden; /* 确保无横向滚动 */
}

/* 欢迎区域 */
.welcome-section {
  max-width: 1200px; /* 统一最大宽度 */
  margin: 0 auto 30px auto;
  background: linear-gradient(135deg, #409EFF 0%, #66B3FF 50%, #87CEEB 100%);
  border-radius: 16px;
  padding: 30px 40px;
  color: white;
  display: flex;
  justify-content: space-between;
  align-items: center;
  box-shadow: 0 8px 32px rgba(64, 158, 255, 0.3);
  position: relative;
  overflow: hidden;
}

.welcome-section::before {
  content: '';
  position: absolute;
  top: -50%;
  right: -50%;
  width: 100%;
  height: 100%;
  background: radial-gradient(circle, rgba(255,255,255,0.1) 1%, transparent 1%);
  background-size: 20px 20px;
  opacity: 0.3;
}

.welcome-content {
  flex: 1;
  z-index: 2;
}

.welcome-title {
  font-size: 28px;
  font-weight: 700;
  margin: 0 0 8px 0;
  text-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.welcome-subtitle {
  font-size: 16px;
  opacity: 0.9;
  margin: 0;
  font-weight: 400;
}

.welcome-decoration {
  position: relative;
  z-index: 2;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100px;
  height: 100px;
}

.decoration-circle {
  position: absolute;
  border-radius: 50%;
  background: rgba(255,255,255,0.15);
}

.circle-1 {
  width: 60px;
  height: 60px;
  top: -10px;
  right: -10px;
}

.circle-2 {
  width: 40px;
  height: 40px;
  bottom: -5px;
  left: -5px;
}

/* 主要内容区域 - 采用 Grid 布局控制 PC 端分栏 */
.main-content {
  max-width: 1200px; /* 统一最大宽度 */
  margin: 0 auto;
  display: grid; 
  grid-template-columns: 1fr; /* 默认单列堆叠 */
  gap: 30px;
}

/* 卡片头部样式 */
.card-header {
  display: flex;
  align-items: center;
  gap: 10px;
  font-weight: 600;
  font-size: 17px;
  color: #2c3e50;
}

.card-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  background: linear-gradient(135deg, #409EFF, #66B3FF);
  border-radius: 10px;
  color: white;
}

/* 第一排：系统信息卡片 */
.first-row {
  width: 100%;
  grid-column: 1 / span 1; 
}

/* 系统信息内部布局 */
.system-info-card {
  border-radius: 16px;
  overflow: hidden;
  border: none;
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(10px);
}

.system-info-card :deep(.el-card__header) {
  background: linear-gradient(135deg, #f8fbff, #f0f7ff);
  border-bottom: 1px solid rgba(64, 158, 255, 0.1);
  padding: 18px 24px;
}

.system-info {
  display: flex;
  justify-content: space-around;
  gap: 30px;
  padding: 10px 0;
}

.info-row {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 24px;
  background: linear-gradient(135deg, #f8fbff, #f0f7ff);
  border-radius: 12px;
  border: 1px solid rgba(64, 158, 255, 0.1);
  transition: all 0.3s ease;
  cursor: pointer;
}

.info-row:hover {
  transform: translateY(-3px);
  box-shadow: 0 6px 20px rgba(64, 158, 255, 0.15);
  border-color: rgba(64, 158, 255, 0.3);
}

.info-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 48px;
  height: 48px;
  background: rgba(64, 158, 255, 0.1);
  border-radius: 12px;
  flex-shrink: 0;
}

.info-content {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.info-label {
  font-size: 13px;
  color: #409EFF;
  font-weight: 500;
}

.info-text {
  font-size: 16px;
  color: #2c3e50;
  font-weight: 600;
}

/* 第二排：用户信息卡片 */
.second-row {
  width: 100%;
  grid-column: 1 / span 1;
}

.user-info-card {
  border-radius: 16px;
  overflow: hidden;
  border: none;
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(10px);
}

.user-info-card :deep(.el-card__header) {
  background: linear-gradient(135deg, #f8fbff, #f0f7ff);
  border-bottom: 1px solid rgba(64, 158, 255, 0.1);
  padding: 18px 24px;
}

/* 用户详情布局 - Grid 2列 */
.user-details {
  padding: 10px 0;
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 20px;
}

.info-item {
  padding: 20px;
  background: linear-gradient(135deg, #fafcff, #f5f9ff);
  border-radius: 12px;
  border: 1px solid rgba(64, 158, 255, 0.1);
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 15px;
  cursor: pointer;
}

.info-item:hover {
  background: linear-gradient(135deg, #f5f9ff, #f0f7ff);
  transform: translateY(-2px);
  box-shadow: 0 4px 15px rgba(64, 158, 255, 0.15);
  border-color: rgba(64, 158, 255, 0.2);
}

.info-icon-item {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  background: rgba(64, 158, 255, 0.1);
  border-radius: 10px;
  flex-shrink: 0;
}

.info-content-item {
  flex: 1;
}

.info-label {
  display: block;
  font-size: 13px;
  color: #409EFF;
  margin-bottom: 6px;
  font-weight: 500;
}

.info-value {
  font-size: 15px;
  color: #2c3e50;
  font-weight: 600;
  /* 确保长文本不溢出 */
  word-break: break-all; 
}

.role-tag {
  font-weight: 600;
  border: none;
}


/* 响应式设计 - 重点优化 PC/移动布局切换 */

/* **PC/大屏幕布局 (宽度 > 1024px) - 分栏显示** */
@media (min-width: 1025px) {
    .main-content {
        /* 左侧固定 350px (系统信息)，右侧自适应 (个人信息) */
        grid-template-columns: 350px 1fr; 
        gap: 30px;
    }
    
    .first-row {
        grid-column: 1 / 2;
        order: 1; /* 确保在左侧 */
    }

    .second-row {
        grid-column: 2 / 3; 
        order: 2; /* 确保在右侧 */
    }
    
    /* 系统信息内部：垂直堆叠，在侧边栏宽度下更美观 */
    .system-info {
        flex-direction: column; 
        gap: 20px;
    }
    
    /* 用户详情保持 2 列 */
    .user-details {
        grid-template-columns: repeat(2, 1fr); 
    }
}

/* 平板/中等屏幕布局 (769px - 1024px) - 全宽堆叠，系统信息并排 */
@media (min-width: 769px) and (max-width: 1024px) {
    .main-content {
        grid-template-columns: 1fr; /* 恢复单列堆叠 */
    }

    .system-info {
        flex-direction: row; /* 日期和时间保持并排 */
        justify-content: space-around;
        gap: 20px;
    }
    
    .user-details {
        grid-template-columns: repeat(2, 1fr); /* 个人信息保持 2 列 */
    }
}


/* 移动端布局 (max-width: 768px 及更小) - 全宽堆叠，信息单列 */
@media (max-width: 768px) {
  .dashboard-container {
    padding: 20px 15px;
  }
  
  .main-content {
      gap: 20px; /* 移动端间距略微减小 */
  }

  .welcome-section {
    flex-direction: column;
    text-align: center;
    padding: 25px 20px;
    gap: 20px;
  }
  
  .welcome-title {
    font-size: 24px;
  }
  
  .welcome-subtitle {
    font-size: 14px;
  }
  
  .welcome-decoration {
    width: 80px;
    height: 80px;
  }
  
  /* 系统信息在移动端垂直堆叠 */
  .system-info {
    flex-direction: column;
    gap: 15px;
    padding: 0;
  }
  
  .info-row {
    padding: 16px;
  }
  
  /* 用户详情在移动端改为单列堆叠 */
  .user-details {
    grid-template-columns: 1fr;
    gap: 15px;
  }
  
  .info-item {
    padding: 14px;
  }
}

/* 极小屏幕 (max-width: 480px) */
@media (max-width: 480px) {
  .welcome-title {
    font-size: 20px;
  }
  
  .card-header {
    font-size: 15px;
  }
  
  .info-text, .info-value {
    font-size: 14px;
  }
}
</style>