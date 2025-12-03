<template>
  <AppHeader />
  <div class="config-container">
    <div class="config-content">
      <el-card class="config-card" shadow="hover">
        <template #header>
          <div class="card-header">
            <div class="card-icon">
              <el-icon :size="20"><Setting /></el-icon>
            </div>
            <span>个人设置</span>
          </div>
        </template>

        <el-form :model="userForm" label-width="120px" class="config-form">
          <!-- 用户名 -->
          <el-form-item label="用户名">
            <el-input 
              v-model="userForm.username" 
              placeholder="请输入用户名"
              :disabled="isLoading"
              clearable
            />
          </el-form-item>

          <!-- 邮箱 -->
          <el-form-item label="邮箱">
            <el-input 
              v-model="userForm.email" 
              placeholder="请输入邮箱"
              :disabled="isLoading"
              clearable
            />
          </el-form-item>

          <!-- 密码 -->
          <el-form-item label="新密码">
            <el-input 
              v-model="userForm.password" 
              type="password" 
              placeholder="留空则不修改密码"
              :disabled="isLoading"
              show-password
              clearable
            />
          </el-form-item>

          <!-- EXIF 设置 -->
          <el-form-item label="解析 EXIF">
            <div class="switch-wrapper">
              <el-switch
                v-model="userForm.parse_exif"
                :disabled="isLoading"
                active-text="开启"
                inactive-text="关闭"
              />
            </div>
            <div class="form-tip">开启后，上传图片时将自动解析 EXIF 信息（拍摄时间、位置等）</div>
          </el-form-item>

          <!-- 提交按钮 -->
          <el-form-item>
            <el-button 
              type="primary" 
              @click="handleSave"
              :loading="isLoading"
              size="large"
            >
              保存
            </el-button>
            <el-button 
              @click="handleReset"
              :disabled="isLoading"
              size="large"
            >
              重置
            </el-button>
          </el-form-item>
        </el-form>
      </el-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Setting } from '@element-plus/icons-vue'
import AppHeader from '@/components/AppHeader.vue'
import { UserAPI } from '@/api'
import { isValidEmail, isValidPassword, isValidUsername } from '@/utils/validation'
import type { User } from '@/type'

// 加载状态
const isLoading = ref(false)

// 原始用户信息
const originalUserInfo = ref<User | null>(null)

// 表单数据
const userForm = ref({
  username: '',
  email: '',
  password: '',
  parse_exif: true
})

// 从 localStorage 加载用户信息
const loadUserInfo = () => {
  const userInfoStr = localStorage.getItem('userInfo')
  if (userInfoStr) {
    try {
      const userInfo: User = JSON.parse(userInfoStr)
      originalUserInfo.value = { ...userInfo }
      userForm.value = {
        username: userInfo.username || '',
        email: userInfo.email || '',
        password: '',
        parse_exif: userInfo.parse_exif !== undefined ? userInfo.parse_exif : true
      }
    } catch (e) {
      console.error('解析用户信息失败:', e)
      ElMessage.error('加载用户信息失败')
    }
  } else {
    ElMessage.warning('未找到用户信息，请重新登录')
  }
}

// 验证表单
const validateForm = (): boolean => {
  if (!userForm.value.username || !isValidUsername(userForm.value.username)) {
    ElMessage.error('用户名长度应在2-20个字符之间')
    return false
  }

  if (!userForm.value.email || !isValidEmail(userForm.value.email)) {
    ElMessage.error('请输入有效的邮箱地址')
    return false
  }

  if (userForm.value.password && !isValidPassword(userForm.value.password)) {
    ElMessage.error('密码长度不能小于6位')
    return false
  }

  return true
}

// 保存设置
const handleSave = async () => {
  if (!validateForm()) {
    return
  }

  if (!originalUserInfo.value) {
    ElMessage.error('用户信息不存在')
    return
  }

  try {
    isLoading.value = true

    // 更新用户基本信息
    await UserAPI.updateUser(
      originalUserInfo.value.id,
      userForm.value.username,
      userForm.value.email,
      userForm.value.password || undefined
    )

    // 更新 EXIF 设置
    if (userForm.value.parse_exif !== originalUserInfo.value.parse_exif) {
      await UserAPI.setParseExif(userForm.value.parse_exif)
    }

    // 更新 localStorage 中的用户信息
    const updatedUserInfo: User = {
      ...originalUserInfo.value,
      username: userForm.value.username,
      email: userForm.value.email,
      parse_exif: userForm.value.parse_exif
    }
    localStorage.setItem('userInfo', JSON.stringify(updatedUserInfo))
    originalUserInfo.value = updatedUserInfo

    // 清空密码字段
    userForm.value.password = ''

    ElMessage.success('设置保存成功')
  } catch (error: any) {
    console.error('保存设置失败:', error)
    ElMessage.error('保存设置失败: ' + (error?.message || '未知错误'))
  } finally {
    isLoading.value = false
  }
}

// 重置表单
const handleReset = () => {
  if (originalUserInfo.value) {
    userForm.value = {
      username: originalUserInfo.value.username || '',
      email: originalUserInfo.value.email || '',
      password: '',
      parse_exif: originalUserInfo.value.parse_exif !== undefined ? originalUserInfo.value.parse_exif : true
    }
    ElMessage.info('已重置为原始值')
  }
}

// 监听用户信息更新事件
const handleUserInfoUpdate = (event: CustomEvent) => {
  console.log('Config: 收到用户信息更新事件，重新加载用户信息');
  loadUserInfo();
};

// 组件挂载时加载用户信息
onMounted(() => {
  loadUserInfo();
  // 监听用户信息更新事件
  window.addEventListener('userInfoUpdated', handleUserInfoUpdate as EventListener);
});

// 组件卸载时移除事件监听
onBeforeUnmount(() => {
  window.removeEventListener('userInfoUpdated', handleUserInfoUpdate as EventListener);
});
</script>

<style scoped>
.config-container {
  min-height: calc(100vh - 60px);
  background: linear-gradient(135deg, #f0f8ff 0%, #e6f3ff 50%, #dcedff 100%);
  padding: 30px;
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
  overflow-x: hidden;
}

.config-content {
  max-width: 800px;
  margin: 0 auto;
}

.config-card {
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 8px 32px rgba(64, 158, 255, 0.15);
}

.card-header {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 18px;
  font-weight: 600;
  color: #2c3e50;
}

.card-icon {
  display: flex;
  align-items: center;
  color: #409EFF;
}

.config-form {
  padding: 20px 0;
}

.form-tip {
  font-size: 12px;
  color: #909399;
  margin-top: 5px;
  line-height: 1.5;
}

:deep(.el-form-item__label) {
  font-weight: 500;
  color: #606266;
}

:deep(.el-input__wrapper) {
  border-radius: 8px;
  transition: all 0.3s ease;
}

:deep(.el-input__wrapper:hover) {
  box-shadow: 0 0 0 1px #409EFF inset;
}

.switch-wrapper {
  display: flex;
  align-items: center;
  margin-bottom: 5px;
}

:deep(.el-switch) {
  margin-right: 10px;
}

:deep(.el-button) {
  border-radius: 8px;
  font-weight: 500;
  padding: 12px 24px;
}

:deep(.el-button--primary) {
  background: linear-gradient(135deg, #409EFF, #66B3FF);
  border: none;
}

:deep(.el-button--primary:hover) {
  background: linear-gradient(135deg, #66B3FF, #409EFF);
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(64, 158, 255, 0.3);
}

/* 响应式设计 */

/* 平板/中等屏幕 (769px - 1024px) */
@media (min-width: 769px) and (max-width: 1024px) {
  .config-container {
    padding: 25px 20px;
  }

  .config-content {
    max-width: 100%;
  }
}

/* 移动端布局 (max-width: 768px) */
@media (max-width: 768px) {
  .config-container {
    padding: 20px 15px;
    min-height: calc(100vh - 60px);
  }

  .config-content {
    max-width: 100%;
  }

  .config-card {
    border-radius: 12px;
  }

  .card-header {
    font-size: 16px;
    gap: 8px;
  }

  .card-icon {
    font-size: 18px;
  }

  .config-form {
    padding: 15px 0;
  }

  :deep(.el-form-item) {
    margin-bottom: 20px;
  }

  :deep(.el-form-item__label) {
    width: 90px !important;
    font-size: 14px;
    padding-right: 10px;
  }

  :deep(.el-input__wrapper) {
    font-size: 14px;
  }

  :deep(.el-input__inner) {
    font-size: 14px;
  }

  .form-tip {
    font-size: 11px;
    margin-top: 8px;
    padding-left: 0;
  }

  .switch-wrapper {
    margin-bottom: 8px;
  }

  :deep(.el-switch) {
    margin-right: 8px;
  }

  /* 按钮组在移动端水平排列并居中 */
  :deep(.el-form-item:last-child) {
    margin-bottom: 0;
  }

  :deep(.el-form-item:last-child .el-form-item__label) {
    width: 0 !important;
    padding: 0 !important;
  }

  :deep(.el-form-item:last-child .el-form-item__content) {
    display: flex !important;
    flex-direction: row !important;
    justify-content: center !important;
    align-items: center !important;
    flex-wrap: nowrap !important;
    gap: 12px;
    margin-left: 0 !important;
    width: 100%;
  }

  :deep(.el-form-item:last-child .el-button) {
    width: auto !important;
    min-width: 100px;
    flex: 0 0 auto;
    margin-bottom: 0;
    padding: 10px 20px;
    font-size: 14px;
  }
}

/* 极小屏幕 (max-width: 480px) */
@media (max-width: 480px) {
  .config-container {
    padding: 15px 10px;
  }

  .config-card {
    border-radius: 10px;
  }

  .card-header {
    font-size: 15px;
    gap: 6px;
  }

  .card-icon {
    font-size: 16px;
  }

  .config-form {
    padding: 10px 0;
  }

  :deep(.el-form-item) {
    margin-bottom: 18px;
  }

  :deep(.el-form-item__label) {
    width: 80px !important;
    font-size: 13px;
    padding-right: 8px;
    line-height: 1.4;
  }

  :deep(.el-input__wrapper) {
    font-size: 13px;
    padding: 8px 12px;
  }

  :deep(.el-input__inner) {
    font-size: 13px;
  }

  .form-tip {
    font-size: 10px;
    margin-top: 6px;
    line-height: 1.4;
  }

  .switch-wrapper {
    margin-bottom: 6px;
  }

  :deep(.el-switch) {
    margin-right: 6px;
  }

  /* 按钮组在极小屏幕也水平排列 */
  :deep(.el-form-item:last-child .el-form-item__label) {
    width: 0 !important;
    padding: 0 !important;
  }

  :deep(.el-form-item:last-child .el-form-item__content) {
    flex-wrap: nowrap !important;
    gap: 10px;
    margin-left: 0 !important;
    width: 100%;
  }

  :deep(.el-form-item:last-child .el-button) {
    min-width: 90px;
    flex: 0 0 auto;
    padding: 9px 16px;
    font-size: 13px;
    border-radius: 6px;
  }

  :deep(.el-switch__label) {
    font-size: 12px;
  }
}
</style>
