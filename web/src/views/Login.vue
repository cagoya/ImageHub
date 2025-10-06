<template>
    <div class="body-wrapper">
      <div class="screen-1">
        <svg class="logo" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" version="1.1" width="300" height="300" viewbox="0 0 640 480" xml:space="preserve">
          <g transform="matrix(3.31 0 0 3.31 320.4 240.4)">
            <circle style="stroke: rgb(0,0,0); stroke-width: 0; stroke-dasharray: none; stroke-linecap: butt; stroke-dashoffset: 0; stroke-linejoin: miter; stroke-miterlimit: 4; fill: rgb(61,71,133); fill-rule: nonzero; opacity: 1;" cx="0" cy="0" r="40"></circle>
          </g>
          <g transform="matrix(0.98 0 0 0.98 268.7 213.7)">
            <circle style="stroke: rgb(0,0,0); stroke-width: 0; stroke-dasharray: none; stroke-linecap: butt; stroke-dashoffset: 0; stroke-linejoin: miter; stroke-miterlimit: 4; fill: rgb(255,255,255); fill-rule: nonzero; opacity: 1;" cx="0" cy="0" r="40"></circle>
          </g>
          <g transform="matrix(1.01 0 0 1.01 362.9 210.9)">
            <circle style="stroke: rgb(0,0,0); stroke-width: 0; stroke-dasharray: none; stroke-linecap: butt; stroke-dashoffset: 0; stroke-linejoin: miter; stroke-miterlimit: 4; fill: rgb(255,255,255); fill-rule: nonzero; opacity: 1;" cx="0" cy="0" r="40"></circle>
          </g>
          <g transform="matrix(0.92 0 0 0.92 318.5 286.5)">
            <circle style="stroke: rgb(0,0,0); stroke-width: 0; stroke-dasharray: none; stroke-linecap: butt; stroke-dashoffset: 0; stroke-linejoin: miter; stroke-miterlimit: 4; fill: rgb(255,255,255); fill-rule: nonzero; opacity: 1;" cx="0" cy="0" r="40"></circle>
          </g>
          <g transform="matrix(0.16 -0.12 0.49 0.66 290.57 243.57)">
            <polygon style="stroke: rgb(0,0,0); stroke-width: 0; stroke-dasharray: none; stroke-linecap: butt; stroke-dashoffset: 0; stroke-linejoin: miter; stroke-miterlimit: 4; fill: rgb(255,255,255); fill-rule: nonzero; opacity: 1;" points="-50,-50 -50,50 50,50 50,-50 "></polygon>
          </g>
          <g transform="matrix(0.16 0.1 -0.44 0.69 342.03 248.34)">
            <polygon style="stroke: rgb(0,0,0); stroke-width: 0; stroke-dasharray: none; stroke-linecap: butt; stroke-dashoffset: 0; stroke-linejoin: miter; stroke-miterlimit: 4; fill: rgb(255,255,255); fill-rule: nonzero; opacity: 1;" vector-effect="non-scaling-stroke" points="-50,-50 -50,50 50,50 50,-50 "></polygon>
          </g>
        </svg>
  
        <div class="email">
          <label for="email">邮箱</label>
          <div class="sec-2">
            <ion-icon name="mail-outline"></ion-icon>
            <input type="email" name="email" placeholder="输入邮箱" v-model="loginForm.email" />
          </div>
        </div>
  
        <div class="password">
          <label for="password">密码</label>
          <div class="sec-2">
            <ion-icon name="lock-closed-outline"></ion-icon>
            <input class="pas" :type="passwordFieldType" name="password" placeholder="输入密码" v-model="loginForm.password" />
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
  // 导入 ElementPlus 的相关组件和图标（如果决定使用）
  // import { ElButton, ElInput } from 'element-plus';
  // import { Message, Lock, View, Hide } from '@element-plus/icons-vue';
  
  // 登录表单数据
  const loginForm = ref({
    email: '',
    password: ''
  });
  
  // 密码显示/隐藏状态
  const showPassword = ref(false);
  const passwordFieldType = ref<'password' | 'text'>('password');
  
  /**
   * 切换密码的可见性
   */
  const togglePasswordVisibility = () => {
    showPassword.value = !showPassword.value;
    passwordFieldType.value = showPassword.value ? 'text' : 'password';
  };
  
  /**
   * 处理登录逻辑
   */
  const handleLogin = () => {
    console.log('尝试登录:', loginForm.value);
    // 可以在这里添加表单验证和 API 调用逻辑
    alert(`尝试登录：\n邮箱: ${loginForm.value.email}\n密码: ${loginForm.value.password}`);
  };
  
  /**
   * 处理注册链接点击
   */
  const handleRegister = () => {
    console.log('跳转到注册页面');
    alert('跳转到注册页面');
  };
  
  /**
   * 处理忘记密码链接点击
   */
  const handleForgotPassword = () => {
    console.log('跳转到忘记密码页面');
    alert('跳转到忘记密码页面');
  };
  
  </script>
  
  <style scoped>
  
  .body-wrapper {
    /* 模拟原 body 样式，让组件居中 */
    user-select: none;
    overflow-y: hidden;
    display: flex;
    justify-content: center;
    align-items: center;
    background: #dde5f4;
    height: 100vh;
    width: 100vw; /* 确保占据整个视口 */
  }
  
  /* 重置样式 (原 CSS 中的 *) */
  * {
    font-family: "Poppins", sans-serif; /* 添加 fallback 字体 */
  }
  
  /* 从这里开始是原 screen-1 及其子元素的样式 */
  
  .screen-1 {
    background: #f1f7fe;
    padding: 2em;
    display: flex;
    flex-direction: column;
    border-radius: 30px;
    box-shadow: 0 0 2em #e6e9f9;
    gap: 2em;
  }
  .screen-1 .logo {
    margin-top: -3em;
  }
  .screen-1 .email {
    background: white;
    box-shadow: 0 0 2em #e6e9f9;
    padding: 1em;
    display: flex;
    flex-direction: column;
    gap: 0.5em;
    border-radius: 20px;
    color: #4d4d4d;
    margin-top: -3em;
  }
  .screen-1 .email .sec-2 {
      display: flex; /* 确保图标和输入框在同一行 */
      align-items: center; /* 垂直居中 */
  }
  .screen-1 .email input {
    outline: none;
    border: none;
    flex-grow: 1; /* 让输入框占据剩余空间 */
    padding-left: 0.5em; /* 稍微增加间隔 */
  }
  .screen-1 .email input::placeholder {
    color: black;
    font-size: 0.9em;
  }
  .screen-1 .email ion-icon {
    color: #4d4d4d;
    margin-bottom: -0.2em;
  }
  .screen-1 .password {
    background: white;
    box-shadow: 0 0 2em #e6e9f9;
    padding: 1em;
    display: flex;
    flex-direction: column;
    gap: 0.5em;
    border-radius: 20px;
    color: #4d4d4d;
  }
  .screen-1 .password .sec-2 {
      display: flex; /* 确保图标和输入框在同一行 */
      align-items: center; /* 垂直居中 */
  }
  .screen-1 .password input {
    outline: none;
    border: none;
    flex-grow: 1;
    padding-left: 0.5em;
  }
  .screen-1 .password input::placeholder {
    color: black;
    font-size: 0.9em;
  }
  .screen-1 .password ion-icon {
    color: #4d4d4d;
    margin-bottom: -0.2em;
  }
  .screen-1 .password .show-hide {
    /* 调整这个负边距以适应 flex 布局 */
    margin-right: 0;
    cursor: pointer;
    /* 原本的 -5em 边距在 flex 布局下通常不适用，这里删除了它 */
  }
  .screen-1 .login {
    padding: 1em;
    background: #3e4684;
    color: white;
    border: none;
    border-radius: 30px;
    font-weight: 600;
  }
  .screen-1 .footer {
    display: flex;
    justify-content: space-between; /* 使用 space-between 替代 space-around 并调整 gap */
    font-size: 0.7em;
    color: #5e5e5e;
    /* gap: 14em; 极大的 gap 在 flex 下可能导致溢出，使用 space-between 更好 */
    padding-bottom: 10em;
  }
  .screen-1 .footer span {
    cursor: pointer;
  }
  button {
    cursor: pointer;
  }
  </style>