<template>
    <div>
      <div class="screen">
        <h2>用户注册</h2>
        
        <div class="name">
          <label for="name">用户名</label>
          <div class="item">
            <input type="text" name="name" placeholder="请输入用户名" v-model="registerForm.name" />
          </div>
        </div>
  
        <div class="email">
          <label for="email">邮箱</label>
          <div class="item">
            <input type="email" name="email" placeholder="请输入邮箱" v-model="registerForm.email" />
          </div>
        </div>
  
        <div class="password">
          <label for="password">密码</label>
          <div class="item">
            <input :type="passwordFieldType" name="password" placeholder="请输入密码" v-model="registerForm.password" />
            <ion-icon class="show-hide" @click="togglePasswordVisibility"></ion-icon>
          </div>
        </div>
  
        <button class="register-btn" @click="handleRegister">注册</button>
  
        <div class="footer">
          <span @click="handleToLogin">已有账号？去登录</span>
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
.screen {
  user-select: none;
  background: #f1f7fe;
  padding: 3em;
  display: flex;
  flex-direction: column;
  border-radius: 30px;
  box-shadow: 0 0 2em #e6e9f9;
  gap: 1.5em;
  min-width: 250px;
  max-width: 350px;
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
}

h2 {
  text-align: center;
  color: #3e4684;
  margin: 0 0 0.5em;
}

.name, .email, .password {
  background: white;
  box-shadow: 0 0 2em #e6e9f9;
  padding: 1em;
  display: flex;
  flex-direction: column;
  gap: 0.5em;
  border-radius: 20px;
  color: #4d4d4d;
}

.item {
  display: flex;
  align-items: center;
}

input {
  outline: none;
  border: none;
  flex-grow: 1;
  padding-left: 0.5em;
}

input::placeholder {
  color: black;
  font-size: 0.9em;
}

ion-icon {
  color: #4d4d4d;
  margin-bottom: -0.2em;
}

.show-hide {
  cursor: pointer;
}

.register-btn {
  padding: 1em;
  background: #3e4684;
  color: white;
  border: none;
  border-radius: 30px;
  font-weight: 600;
  cursor: pointer;
}

.footer {
  display: flex;
  justify-content: center;
  font-size: 0.7em;
  color: #5e5e5e;
}

.footer span {
  cursor: pointer;
}
</style>