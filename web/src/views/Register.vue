<template>
  <div class="register-container">
    <div class="screen">
      <div class="header">
        <h2>用户注册</h2>
        <p>创建您的新账户</p>
      </div>
      
      <div class="form-group">
        <div class="input-field">
          <div class="icon">
            <el-icon><User /></el-icon>
          </div>
          <input 
            type="text" 
            placeholder="请输入用户名" 
            v-model="registerForm.name" 
            class="input"
          />
        </div>
      </div>

      <div class="form-group">
        <div class="input-field">
          <div class="icon">
            <el-icon><Message /></el-icon>
          </div>
          <input 
            type="email" 
            placeholder="请输入邮箱" 
            v-model="registerForm.email" 
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
            v-model="registerForm.password" 
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
        class="register-btn" 
        @click="handleRegister"
        :disabled="isLoading" 
      >
        <span class="btn-text">{{ isLoading ? '注册中...' : '注册账户' }}</span>
        <el-icon class="btn-icon is-loading" v-if="isLoading">
          <Loading />
        </el-icon>
        <el-icon class="btn-icon" v-else>
          <ArrowRight />
        </el-icon>
      </button>

      <div class="footer">
        <span>已有账号？</span>
        <span class="login-link" @click="handleToLogin">立即登录</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { UserAPI } from '../api';
import { ElMessage } from 'element-plus';
import type { RegisterForm } from '@/type'; 
import { isValidEmail, isValidPassword, isValidUsername } from '@/utils/validation'; 
// Element Plus 图标导入
import { 
  User, 
  Message, 
  Lock, 
  View, 
  Hide, 
  ArrowRight,
  Loading // 引入 Loading 图标
} from '@element-plus/icons-vue';

const router = useRouter();

// 注册表单数据
const registerForm = ref<RegisterForm>({
  name: '',
  email: '',
  password: ''
});

// 密码显示/隐藏状态
const showPassword = ref(false);
const passwordFieldType = ref<'password' | 'text'>('password');

// 加载状态
const isLoading = ref(false); 

// 切换密码的可见性
const togglePasswordVisibility = () => {
  showPassword.value = !showPassword.value;
  passwordFieldType.value = showPassword.value ? 'text' : 'password';
};

// 邮箱格式校验（使用您提供的函数）
const isValidEmailFormat = (email: string): boolean => {
  return isValidEmail(email);
};

// 处理注册逻辑
const handleRegister = async () => {
  if (isLoading.value) return; // 防止重复点击

  const { name, email, password } = registerForm.value;

  // 1. 基础非空校验
  if (!name || !email || !password) {
    ElMessage.error('所有字段都是必填项');
    return;
  }

  // 2. 邮箱格式校验
  if (!isValidEmailFormat(email)) {
    ElMessage.error('邮箱格式不正确，请检查');
    return;
  }

  // 3. 密码长度校验
  if (password.length < 6) {
    ElMessage.error('密码长度不能小于6位');
    return;
  }

  // 开启加载状态
  isLoading.value = true; 

  try {
    // 假设 UserAPI.register 成功返回 status: 200
    const response = await UserAPI.register(name, email, password);
    
    // 您原代码中的注册成功判断，此处保持一致
    if (response.status === 200) {
      ElMessage.success('注册成功，即将跳转到登录页');
    }
    
    // 注册成功后，清除表单并延迟跳转
    registerForm.value.name = '';
    registerForm.value.email = '';
    registerForm.value.password = '';
    // 延迟 1 秒跳转，让用户看到成功提示
    await new Promise(resolve => setTimeout(resolve, 1000)); 
    router.push('/login');
  } catch (error) {
    ElMessage.error('注册失败: ' + (error instanceof Error ? error.message : '请稍后重试'));
  } finally {
    // 无论成功失败，关闭加载状态
    isLoading.value = false; 
  }
};

// 处理跳转到登录链接点击
const handleToLogin = () => {
  router.push('/login');
};
</script>

<style scoped>
/* 容器：保证全屏且居中，消除横向滚动 */
.register-container {
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
  gap: 1.5em; 
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

.register-btn {
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

.register-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(64, 158, 255, 0.4);
  background: linear-gradient(135deg, #66B3FF, #409EFF);
}

.register-btn:active {
  transform: translateY(0);
}

.btn-text {
  font-weight: 600;
}

.btn-icon {
  transition: transform 0.3s ease;
}

.register-btn:hover:not(:disabled) .btn-icon:not(.is-loading) {
  transform: translateX(3px);
}

.footer {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 0.5em;
  font-size: 14px;
  color: #7f8c8d;
  margin-top: 1em;
}

.login-link {
  color: #409EFF;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
  text-decoration: none;
  position: relative;
}

.login-link:hover {
  color: #66B3FF;
}

.login-link::after {
  content: '';
  position: absolute;
  bottom: -2px;
  left: 0;
  width: 0;
  height: 2px;
  background: #409EFF;
  transition: width 0.3s ease;
}

.login-link:hover::after {
  width: 100%;
}

/* 响应式设计 - 针对手机/小屏幕优化 */
@media (max-width: 480px) {
  .register-container {
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

  .register-btn {
    padding: 0.8em 1.5em; /* 减小按钮的内边距 */
    font-size: 15px; /* 减小按钮字体 */
    margin-top: 0.8em;
  }

  .footer {
    font-size: 13px; /* 移动端页脚字体更小 */
  }
}

/* 加载状态 */
.register-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
  transform: none;
  box-shadow: 0 4px 15px rgba(64, 158, 255, 0.3);
}

.register-btn:disabled:hover {
  transform: none;
}
</style>