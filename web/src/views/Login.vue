<template>
  <div class="login-container">
    <div class="screen">
      <div class="header">
        <h2>用户登录</h2>
        <p>欢迎回来，请登录您的账户</p>
      </div>
      
      <div class="form-group">
        <div class="input-field">
          <div class="icon">
            <el-icon><User /></el-icon>
          </div>
          <input 
            type="text" 
            placeholder="请输入用户名" 
            v-model="loginForm.name" 
            class="input"
          />
        </div>
      </div>

      <div class="form-group">
        <div class="input-field">
          <div class="icon">
            <el-icon><Lock /></el-icon>
          </div>
          <input 
            :type="passwordFieldType" 
            placeholder="请输入密码" 
            v-model="loginForm.password" 
            class="input"
          />
          <div class="password-toggle" @click="togglePasswordVisibility">
            <el-icon>
              <View v-if="showPassword" />
              <Hide v-else />
            </el-icon>
          </div>
        </div>
      </div>

      <button 
        class="login-btn" 
        @click="handleLogin"
        :disabled="isLoading"
      >
        <span class="btn-text">{{ isLoading ? '登录中...' : '登录账户' }}</span>
        <el-icon class="btn-icon is-loading" v-if="isLoading">
          <Loading />
        </el-icon>
        <el-icon class="btn-icon" v-else>
          <ArrowRight />
        </el-icon>
      </button>

      <div class="footer">
        <span class="register-link" @click="handleRegister">注册账户</span>
        <span class="forgot-link" @click="handleForgotPassword">忘记密码</span>
      </div>
    </div>
    
    <el-dialog
      v-model="dialogVisible"
      title="修改密码"
      width="500"
      class="forgot-dialog"
    >
      <el-form>
        <el-form-item label="用户名">
          <el-input 
            v-model="forgetPasswordForm.username" 
            placeholder="请输入用户名"
            size="large"
          />
        </el-form-item>
        <el-form-item label="邮箱">
          <el-input 
            v-model="forgetPasswordForm.email" 
            placeholder="请输入邮箱"
            size="large"
          />
        </el-form-item>
        <el-form-item label="新密码">
          <el-input 
            v-model="forgetPasswordForm.newPassword" 
            type="password" 
            placeholder="请输入新密码"
            size="large"
            show-password
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false" size="large">取消</el-button>
          <el-button 
            type="primary" 
            @click="handleSubmitForgotPassword" 
            size="large"
            class="confirm-btn"
          >
            确定修改
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>
  
<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
// ⚠️ 假设这些是您项目中已有的导入和路径
import { UserAPI } from '../api';
import { ElDialog, ElForm, ElFormItem, ElInput, ElButton, ElMessage } from 'element-plus';
import type { LoginForm, ForgetPasswordForm } from '@/type';
import { isValidEmail, isValidPassword } from '@/utils/validation';
// Element Plus 图标导入
import { 
  User, 
  Lock, 
  View, 
  Hide, 
  ArrowRight,
  Loading // 引入 Loading 图标
} from '@element-plus/icons-vue';
  
const router = useRouter();

// 状态管理
const dialogVisible = ref(false);
const isLoading = ref(false); 
const showPassword = ref(false);
const passwordFieldType = ref<'password' | 'text'>('password');

// 表单数据
const forgetPasswordForm = ref<ForgetPasswordForm>({
  username: '',
  newPassword: '',
  email: ''
});
const loginForm = ref<LoginForm>({
  name: '',
  password: ''
});
  
// 切换密码的可见性
const togglePasswordVisibility = () => {
  showPassword.value = !showPassword.value;
  passwordFieldType.value = showPassword.value ? 'text' : 'password';
};
  
// 处理登录逻辑
const handleLogin = async () => {
  if (isLoading.value) return; // 防止重复点击

  try {
    if (!loginForm.value.name || !loginForm.value.password) {
      ElMessage.error('请输入用户名和密码');
      return;
    }
    
    isLoading.value = true; 
    const response = await UserAPI.login(loginForm.value.name, loginForm.value.password);
    
    if (response.status === 200) {
      localStorage.setItem('userToken', response.data.data.token);
      localStorage.setItem('userInfo', JSON.stringify(response.data.data.user));
      localStorage.setItem('userID', response.data.data.user.id);
      ElMessage.success('登录成功');
      router.push('/dashboard');
    }
  } catch (error: any) {
    // 提取后端返回的错误信息
    const errorMessage = error?.response?.data?.message || error?.message || '登录失败，请检查用户名和密码';
    ElMessage.error(errorMessage);
  } finally {
    isLoading.value = false; 
  }
};
  
// 处理注册链接点击
const handleRegister = () => {
  router.push('/register');
};
  
// 处理忘记密码链接点击
const handleForgotPassword = () => {
  dialogVisible.value = true;
};

// 处理提交修改密码表单
const handleSubmitForgotPassword = async () => {
  try {
    if (!forgetPasswordForm.value.username || !forgetPasswordForm.value.newPassword || !forgetPasswordForm.value.email) {
      ElMessage.error('请填写完整信息');
      return;
    }
    
    // 邮箱格式校验
    if (!isValidEmail(forgetPasswordForm.value.email)) {
      ElMessage.error('邮箱格式不正确，请检查');
      return;
    }
    
    // 密码长度校验
    if (!isValidPassword(forgetPasswordForm.value.newPassword)) {
      ElMessage.error('密码长度不能小于6位');
      return;
    }
    
    const response = await UserAPI.forgetPassword(forgetPasswordForm.value.username, forgetPasswordForm.value.email, forgetPasswordForm.value.newPassword);
    
    if (response.status === 200) {
      ElMessage.success('修改密码成功');
      dialogVisible.value = false;
      // 清空表单
      forgetPasswordForm.value = {
        username: '',
        newPassword: '',
        email: ''
      };
    }
  } catch (error: any) {
    // 提取后端返回的错误信息
    const errorMessage = error?.response?.data?.message || error?.message || '修改密码失败，请检查用户名和邮箱是否匹配';
    ElMessage.error(errorMessage);
  }
};
</script>
  
<style scoped>
/* 容器：保证全屏且居中，消除横向滚动 */
.login-container {
  min-height: 100vh;
  height: 100%; 
  background: linear-gradient(135deg, #f0f8ff 0%, #e6f3ff 50%, #dcedff 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
  box-sizing: border-box; 
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
  overflow-x: hidden; /* 消除横向滚动条 */
}

/* 屏幕/表单容器：控制尺寸和样式 */
.screen {
  background: rgba(255, 255, 255, 0.98); 
  backdrop-filter: blur(8px); 
  padding: 3em 2.5em;
  display: flex;
  flex-direction: column;
  border-radius: 24px;
  box-shadow: 
    0 8px 32px rgba(64, 158, 255, 0.15),
    0 2px 8px rgba(64, 158, 255, 0.1);
  gap: 1.5em; /* 略微减小整体间距，更紧凑 */
  min-width: 300px; 
  max-width: 400px; /* PC端保持在合理的最大宽度 */
  width: 100%;
  border: 1px solid rgba(64, 158, 255, 0.1);
}

.header {
  text-align: center;
  margin-bottom: 0.5em;
}

h2 {
  color: #2c3e50;
  margin: 0 0 0.3em;
  font-size: 26px; /* PC 端标题大小优化 */
  font-weight: 700;
  background: linear-gradient(135deg, #409EFF, #66B3FF);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.header p {
  color: #7f8c8d;
  margin: 0;
  font-size: 14px;
  font-weight: 400;
}

.form-group {
  margin-bottom: 0.5em;
}

.input-field {
  background: #f8fbff;
  border: 2px solid #e6f3ff;
  border-radius: 16px;
  padding: 0.8em 1.2em;
  display: flex;
  align-items: center;
  gap: 0.8em;
  transition: all 0.3s ease;
  position: relative;
}

.input-field:focus-within {
  border-color: #409EFF;
  background: white;
  box-shadow: 0 0 0 3px rgba(64, 158, 255, 0.1);
  transform: translateY(-1px);
}

.icon {
  display: flex;
  align-items: center;
  color: #409EFF;
  font-size: 18px;
}

.input {
  outline: none;
  border: none;
  flex-grow: 1;
  background: transparent;
  font-size: 15px;
  color: #2c3e50;
  font-weight: 500;
  padding: 0.2em 0;
}

.input::placeholder {
  color: #a0b4c0;
  font-weight: 400;
}

.password-toggle {
  display: flex;
  align-items: center;
  color: #7f8c8d;
  cursor: pointer;
  padding: 0.2em;
  border-radius: 6px;
  transition: all 0.2s ease;
}

.password-toggle:hover {
  color: #409EFF;
  background: rgba(64, 158, 255, 0.1);
}

.login-btn {
  padding: 1em 2em;
  background: linear-gradient(135deg, #409EFF, #66B3FF);
  color: white;
  border: none;
  border-radius: 16px;
  font-weight: 600;
  font-size: 16px;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.8em;
  margin-top: 1em; /* 增加与上方表单的间距 */
  box-shadow: 0 4px 15px rgba(64, 158, 255, 0.3);
}

.login-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(64, 158, 255, 0.4);
  background: linear-gradient(135deg, #66B3FF, #409EFF);
}

.login-btn:active {
  transform: translateY(0);
}

.btn-text {
  font-weight: 600;
}

.btn-icon {
  transition: transform 0.3s ease;
}

/* 非加载状态下，鼠标悬停时图标平移 */
.login-btn:hover:not(:disabled) .btn-icon:not(.is-loading) {
  transform: translateX(3px);
}

/* 底部链接：PC 端左右对齐，移动端堆叠对齐 */
.footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 14px;
  color: #7f8c8d;
  margin-top: 1em;
  padding: 0 0.5em;
}

.register-link, .forgot-link {
  color: #409EFF;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
  text-decoration: none;
  position: relative;
}

.register-link:hover, .forgot-link:hover {
  color: #66B3FF;
}

.register-link::after, .forgot-link::after {
  content: '';
  position: absolute;
  bottom: -2px;
  left: 0;
  width: 0;
  height: 2px;
  background: #409EFF;
  transition: width 0.3s ease;
}

.register-link:hover::after, .forgot-link:hover::after {
  width: 100%;
}

/* 响应式设计 - 针对手机/小屏幕优化 */
@media (max-width: 480px) {
  .login-container {
    padding: 15px; /* 移动端更小的内边距 */
  }

  .screen {
    width: 95%; 
    max-width: 100%; /* 确保在小屏幕上填满可用宽度 */
    padding: 2em 1.5em; /* 减小内边距 */
    gap: 1.2em; /* 进一步减小移动端组件间距 */
  }

  h2 {
    font-size: 22px; /* 移动端标题更小 */
  }

  .input-field {
    padding: 0.6em 1em; /* 减小输入框的内边距，更适合手指点击 */
  }

  .input {
    font-size: 14px; /* 移动端字体更小 */
  }

  .login-btn {
    padding: 0.8em 1.5em; /* 减小按钮的内边距 */
    font-size: 15px; /* 减小按钮字体 */
    margin-top: 0.8em;
  }
  
  .footer {
    /* 移动端改为垂直堆叠布局 */
    flex-direction: column;
    gap: 0.8em;
    text-align: center;
    font-size: 13px; /* 移动端页脚字体更小 */
  }
}

/* 加载状态 */
.login-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
  transform: none;
  box-shadow: 0 4px 15px rgba(64, 158, 255, 0.3);
}

.login-btn:disabled:hover {
  transform: none;
}
</style>

<style>
/* 确保弹窗在移动端有合适的宽度 */
@media (max-width: 480px) {
  .forgot-dialog {
    max-width: 90vw; 
  }
  .forgot-dialog .el-dialog {
    /* 强制覆盖 Element Plus 默认的 500px 宽度 */
    --el-dialog-width: 90% !important; 
    /* 强制居中 */
    margin-top: 5vh !important;
  }
}

/* 保持您对弹窗的其他美化样式 */
.forgot-dialog .el-dialog {
  border-radius: 16px;
  overflow: hidden;
}

.forgot-dialog .el-dialog__header {
  background: linear-gradient(135deg, #f8fbff, #f0f7ff);
  border-bottom: 1px solid rgba(64, 158, 255, 0.1);
  margin: 0;
  padding: 12px 20px;
}

.forgot-dialog .el-dialog__title {
  color: #2c3e50;
  font-weight: 600;
  font-size: 18px;
}

.forgot-dialog .el-dialog__body {
  padding-top: 24px;
}

.forgot-dialog .el-form-item__label {
  color: #409EFF;
  font-weight: 500;
}

.forgot-dialog .el-input__wrapper {
  border-radius: 12px;
  border: 2px solid #e6f3ff;
  background: #f8fbff;
  transition: all 0.3s ease;
}

.forgot-dialog .el-input__wrapper:hover,
.forgot-dialog .el-input__wrapper.is-focus {
  border-color: #409EFF;
  box-shadow: 0 0 0 3px rgba(64, 158, 255, 0.1);
}

.confirm-btn {
  background: linear-gradient(135deg, #409EFF, #66B3FF);
  border: none;
  border-radius: 12px;
  font-weight: 600;
}

.confirm-btn:hover {
  background: linear-gradient(135deg, #66B3FF, #409EFF);
}
</style>