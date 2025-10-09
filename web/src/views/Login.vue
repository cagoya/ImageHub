<template>
    <div class="body-wrapper">
      <div class="screen">
  
        <div class="name">
          <label for="name">用户名</label>
          <div class="item">
            <ion-icon name="name-outline"></ion-icon>
            <input type="name" name="name" placeholder="请输入用户名" v-model="loginForm.name" />
          </div>
        </div>
  
        <div class="password">
          <label for="password">密码</label>
          <div class="item">
            <ion-icon name="lock-closed-outline"></ion-icon>
            <input class="pas" :type="passwordFieldType" name="password" placeholder="请输入密码" v-model="loginForm.password" />
            <ion-icon class="show-hide" :name="showPassword ? 'eye-off-outline' : 'eye-outline'" @click="togglePasswordVisibility"></ion-icon>
          </div>
        </div>
  
        <button class="login" @click="handleLogin">登录</button>
  
        <div class="footer">
          <span @click="handleRegister">注册</span>
          <span @click="handleForgotPassword">忘记密码?</span>
        </div>
      </div>
    </div>
  </template>
  
  <script setup lang="ts">
  import { ref } from 'vue';
  import { UserAPI } from '../api';

  // 登录表单数据
  const loginForm = ref({
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
      alert('请输入用户名和密码');
      return;
    }
    const response = await UserAPI.login(loginForm.value.name, loginForm.value.password);
    if (response.status === 200) {
      localStorage.setItem('userToken', response.data.data.token);
      localStorage.setItem('userInfo', JSON.stringify(response.data.data.user));
      alert('登录成功');
    }
  } catch (error) {
    console.error('登录失败:', error);
    alert('登录失败: ' + (error instanceof Error ? error.message : '请检查用户名和密码'));
  }
};
  
  // 处理注册链接点击
  const handleRegister = () => {
    console.log('跳转到注册页面');
    alert('跳转到注册页面');
  };
  
  // 处理忘记密码链接点击
  const handleForgotPassword = () => {
    console.log('跳转到忘记密码页面');
    alert('跳转到忘记密码页面');
  };
  
  </script>
  
  <style scoped>
  
  .body-wrapper {
    user-select: none;
    overflow-y: hidden;
    display: flex;
    justify-content: center;
    align-items: center;
    height: 100vh;
    width: 100vw;
  }
  
  .screen {
    background: #f1f7fe;
    padding: 3em;
    display: flex;
    flex-direction: column;
    border-radius: 30px;
    box-shadow: 0 0 2em #e6e9f9;
    gap: 1.5em;
    min-width: 250px;
    max-width: 350px;
  }
  .screen .name {
    background: white;
    box-shadow: 0 0 2em #e6e9f9;
    padding: 1em;
    display: flex;
    flex-direction: column;
    gap: 0.5em;
    border-radius: 20px;
    color: #4d4d4d;
  }
  .screen .name .item {
      display: flex;
      align-items: center;
  }
  .screen .name input {
    outline: none;
    border: none;
    flex-grow: 1;
    padding-left: 0.5em;
  }
  .screen .name input::placeholder {
    color: black;
    font-size: 0.9em;
  }
  .screen .name ion-icon {
    color: #4d4d4d;
    margin-bottom: -0.2em;
  }
  .screen .password {
    background: white;
    box-shadow: 0 0 2em #e6e9f9;
    padding: 1em;
    display: flex;
    flex-direction: column;
    gap: 0.5em;
    border-radius: 20px;
    color: #4d4d4d;
  }
  .screen .password .item {
      display: flex;
      align-items: center;
  }
  .screen .password input {
    outline: none;
    border: none;
    flex-grow: 1;
    padding-left: 0.5em;
  }
  .screen .password input::placeholder {
    color: black;
    font-size: 0.9em;
  }
  .screen .password ion-icon {
    color: #4d4d4d;
    margin-bottom: -0.2em;
  }
  .screen .password .show-hide {
    margin-right: 0;
    cursor: pointer;
  }
  .screen .login {
    padding: 1em;
    background: #3e4684;
    color: white;
    border: none;
    border-radius: 30px;
    font-weight: 600;
  }
  .screen .footer {
    display: flex;
    justify-content: space-between;
    font-size: 0.7em;
    color: #5e5e5e;
    gap: 14em;
    padding-bottom: 10em;
  }
  .screen .footer span {
    cursor: pointer;
  }
  button {
    cursor: pointer;
  }
  </style>