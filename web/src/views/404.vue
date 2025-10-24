<template>
    <div class="error-page">
      <section class="error-container">
        <span class="four"><span class="screen-reader-text">4</span></span>
        <span class="zero"><span class="screen-reader-text">0</span></span>
        <span class="four"><span class="screen-reader-text">4</span></span>
      </section>
      <div class="link-container">
        <router-link to="/login" class="more-link">返回首页</router-link>
      </div>
    </div>
  </template>
  
  <script setup lang="ts">

  import { defineComponent } from 'vue';
  
  defineComponent({
    name: 'NotFoundBlue',
  });
  </script>
  
  <style scoped>
  /* 统一配色与 Login 页面一致 */
  

  .error-page {
    color: #4d4d4d;

    min-height: 100vh;
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    text-align: center;
  }
  
  h1 {
    text-align: center;
    margin: 30px 15px;
  }
  
  /* 404 容器样式 */
  .error-container {
    text-align: center;
    font-size: 106px;
    font-weight: 800;
    margin: 70px 15px;
  }
  
  .error-container > span {
    display: inline-block;
    position: relative;
  }
  
  /* 第一个和第三个 "4" 的样式，与Login按钮一致 */
  .error-container > span.four {
    width: 136px;
    height: 43px;
    border-radius: 999px;
    background: linear-gradient(140deg, rgba(0, 0, 0, 0.1) 0%, rgba(0, 0, 0, 0.07) 43%, transparent 44%, transparent 100%),
      linear-gradient(105deg, transparent 0%, transparent 40%, rgba(0, 0, 0, 0.06) 41%, rgba(0, 0, 0, 0.07) 76%, transparent 77%, transparent 100%),
      linear-gradient(to right, #5a5fb5, #3e4684); /* 与Login按钮一致的渐变 */
  }
  
  .error-container > span.four:before,
  .error-container > span.four:after {
    content: '';
    display: block;
    position: absolute;
    border-radius: 999px;
  }
  
  .error-container > span.four:before {
    width: 43px;
    height: 156px;
    left: 60px;
    bottom: -43px;
    background: linear-gradient(128deg, rgba(0, 0, 0, 0.1) 0%, rgba(0, 0, 0, 0.07) 40%, transparent 41%, transparent 100%),
      linear-gradient(116deg, rgba(0, 0, 0, 0.1) 0%, rgba(0, 0, 0, 0.07) 50%, transparent 51%, transparent 100%),
      linear-gradient(to top, #3e4684, #4a5491, #5a5fb5, #6b6fd1);
  }
  
  .error-container > span.four:after {
    width: 137px;
    height: 43px;
    transform: rotate(-49.5deg);
    left: -18px;
    bottom: 36px;
    background: linear-gradient(to right, #3e4684, #4a5491, #5a5fb5, #6b6fd1);
  }
  
  /* "0" 的样式，与Login按钮一致 */
  .error-container > span.zero {
    vertical-align: text-top;
    width: 156px;
    height: 156px;
    border-radius: 999px;
    background: linear-gradient(-45deg, transparent 0%, rgba(0, 0, 0, 0.06) 50%, transparent 51%, transparent 100%),
      linear-gradient(to top right, #3e4684, #3e4684, #4a5491, #5a5fb5, #6b6fd1, #6b6fd1);
    overflow: hidden;
    /* 移除动画：animation: bgshadow 5s infinite; */
  }
  
  .error-container > span.zero:before {
    content: '';
    display: block;
    position: absolute;
    transform: rotate(45deg);
    width: 90px;
    height: 90px;
    background-color: transparent;
    left: 0px;
    bottom: 0px;
    background: linear-gradient(95deg, transparent 0%, transparent 8%, rgba(0, 0, 0, 0.07) 9%, transparent 50%, transparent 100%),
      linear-gradient(85deg, transparent 0%, transparent 19%, rgba(0, 0, 0, 0.05) 20%, rgba(0, 0, 0, 0.07) 91%, transparent 92%, transparent 100%);
  }
  
  .error-container > span.zero:after {
    content: '';
    display: block;
    position: absolute;
    border-radius: 999px;
    width: 70px;
    height: 70px;
    left: 43px;
    bottom: 43px;
    background: #dde5f4; /* 与页面背景保持一致 */
    box-shadow: -2px 2px 2px 0px rgba(0, 0, 0, 0.1);
  }
  
  .screen-reader-text {
    position: absolute;
    top: -9999em;
    left: -9999em;
  }
  
  /* demo stuff */
  .error-page * {
    box-sizing: border-box;
  }
  
  .link-container {
    text-align: center;
  }
  
  .more-link {
    text-transform: uppercase;
    font-size: 13px;
    background-color: #3e4684; /* 与Login按钮一致 */
    padding: 10px 20px;
    border-radius: 30px; /* 圆角与 Login 按钮保持一致 */
    color: #fff;
    display: inline-block;
    margin-right: 5px;
    margin-bottom: 5px;
    line-height: 1.5;
    text-decoration: none;
    margin-top: 50px;
    letter-spacing: 1px;
  }
  </style>