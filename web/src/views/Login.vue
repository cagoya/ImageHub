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

      <button class="login-btn" @click="handleLogin">
        <span class="btn-text">登录账户</span>
        <el-icon class="btn-icon"><ArrowRight /></el-icon>
      </button>

      <div class="footer">
        <span class="register-link" @click="handleRegister">注册新账户</span>
        <span class="forgot-link" @click="handleForgotPassword">忘记密码?</span>
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
import { UserAPI } from '../api';
import { ElDialog, ElForm, ElFormItem, ElInput, ElButton, ElMessage } from 'element-plus';
import type { LoginForm, ForgetPasswordForm } from '@/type';
import { isValidEmail, isValidPassword } from '@/utils/validation';
import { User, Lock, View, Hide, ArrowRight } from '@element-plus/icons-vue';
  
const router = useRouter();

const dialogVisible = ref(false);
const forgetPasswordForm = ref<ForgetPasswordForm>({
  username: '',
  newPassword: '',
  email: ''
});

// 登录表单数据
const loginForm = ref<LoginForm>({
  name: '',
  password: ''
});
  
// 密码显示/隐藏状态
const showPassword = ref(false);
const passwordFieldType = ref<'password' | 'text'>('password');
  
// 切换密码的可见性
const togglePasswordVisibility = () => {
  showPassword.value = !showPassword.value;
  passwordFieldType.value = showPassword.value ? 'text' : 'password';
};
  
// 处理登录逻辑
const handleLogin = async () => {
  try {
    if (!loginForm.value.name || !loginForm.value.password) {
      ElMessage.error('请输入用户名和密码');
      return;
    }
    const response = await UserAPI.login(loginForm.value.name, loginForm.value.password);
    if (response.status === 200) {
      localStorage.setItem('userToken', response.data.data.token);
      localStorage.setItem('userInfo', JSON.stringify(response.data.data.user));
      localStorage.setItem('userID', response.data.data.user.id);
      ElMessage.success('登录成功');
      router.push('/dashboard');
    }
  } catch (error) {
    ElMessage.error('登录失败: ' + (error instanceof Error ? error.message : '请检查用户名和密码'));
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
    if (forgetPasswordForm.value.newPassword.length < 6) {
      ElMessage.error('密码长度不能小于6位');
      return;
    }
    const response = await UserAPI.forgetPassword(forgetPasswordForm.value.username, forgetPasswordForm.value.email, forgetPasswordForm.value.newPassword);
    if (response.status === 200) {
      ElMessage.success('修改密码成功');
    }
    dialogVisible.value = false;
    // 清空表单
    forgetPasswordForm.value = {
      username: '',
      newPassword: '',
      email: ''
    };
  } catch (error) {
    ElMessage.error('修改密码失败: ' + (error instanceof Error ? error.message : '请检查用户名和邮箱是否匹配'));
  }
};
</script>
  
<style scoped>
.login-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #f0f8ff 0%, #e6f3ff 50%, #dcedff 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
}

.screen {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  padding: 3em 2.5em;
  display: flex;
  flex-direction: column;
  border-radius: 24px;
  box-shadow: 
    0 8px 32px rgba(64, 158, 255, 0.15),
    0 2px 8px rgba(64, 158, 255, 0.1);
  gap: 1.8em;
  min-width: 320px;
  max-width: 420px;
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
  font-size: 28px;
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
  margin-top: 0.5em;
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

.login-btn:hover .btn-icon {
  transform: translateX(3px);
}

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

/* 响应式设计 */
@media (max-width: 480px) {
  .screen {
    padding: 2em 1.5em;
    margin: 1em;
  }
  
  h2 {
    font-size: 24px;
  }
  
  .input-field {
    padding: 0.7em 1em;
  }
  
  .login-btn {
    padding: 0.9em 1.5em;
  }
  
  .footer {
    flex-direction: column;
    gap: 0.8em;
    text-align: center;
  }
}

/* 加载状态 */
.login-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
  transform: none;
}

.login-btn:disabled:hover {
  transform: none;
  box-shadow: 0 4px 15px rgba(64, 158, 255, 0.3);
}
</style>

<style>
.forgot-dialog .el-dialog {
  border-radius: 16px;
  overflow: hidden;
}

.forgot-dialog .el-dialog__header {
  background: linear-gradient(135deg, #f8fbff, #f0f7ff);
  border-bottom: 1px solid rgba(64, 158, 255, 0.1);
  margin: 0;
  padding: 20px;
}

.forgot-dialog .el-dialog__title {
  color: #2c3e50;
  font-weight: 600;
  font-size: 18px;
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