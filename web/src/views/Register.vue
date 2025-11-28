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

      <button class="register-btn" @click="handleRegister">
        <span class="btn-text">注册账户</span>
        <el-icon class="btn-icon"><ArrowRight /></el-icon>
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
import { User, Message, Lock, View, Hide, ArrowRight } from '@element-plus/icons-vue';

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

// 切换密码的可见性
const togglePasswordVisibility = () => {
  showPassword.value = !showPassword.value;
  passwordFieldType.value = showPassword.value ? 'text' : 'password';
};

// 邮箱格式校验
const isValidEmailFormat = (email: string): boolean => {
return isValidEmail(email);
};

// 处理注册逻辑
const handleRegister = async () => {
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

  try {
    const response = await UserAPI.register(name, email, password);
    if (response.status === 200) {
      ElMessage.success('注册成功');
    }
    
    // 注册成功后，清除表单并跳转到登录页
    registerForm.value.name = '';
    registerForm.value.email = '';
    registerForm.value.password = '';
    router.push('/login');
  } catch (error) {
    ElMessage.error('注册失败: ' + (error instanceof Error ? error.message : '请稍后重试'));
  }
};

// 处理跳转到登录链接点击
const handleToLogin = () => {
  router.push('/login');
};

</script>

<style scoped>
.register-container {
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
margin-top: 0.5em;
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

.register-btn:hover .btn-icon {
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

.register-btn {
  padding: 0.9em 1.5em;
}
}

/* 加载状态 */
.register-btn:disabled {
opacity: 0.7;
cursor: not-allowed;
transform: none;
}

.register-btn:disabled:hover {
transform: none;
box-shadow: 0 4px 15px rgba(64, 158, 255, 0.3);
}
</style>