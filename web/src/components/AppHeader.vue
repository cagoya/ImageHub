<template>
    <el-container class="header-container">
      <el-header height="60px" class="main-header">
        
        <div class="logo">
          <span class="logo-text">图片管理系统</span>
        </div>
  
        <el-menu
          :default-active="activeIndex"
          mode="horizontal"
          :router="true"            class="el-menu-demo"
        >
          <el-menu-item index="/dashboard">首页</el-menu-item>
          
          <el-menu-item index="/upload">上传图片</el-menu-item>
          
          <el-menu-item index="/gallery">查看图片</el-menu-item>
          
          <el-menu-item v-if="isAdmin" index="/manage">用户管理</el-menu-item>
        </el-menu>
  
        <div class="logout">
            <el-button type="text" @click="handleLogout">退出登录</el-button>
        </div>
  
      </el-header>
      
    </el-container>
  </template>
  
  <script lang="ts" setup>
import { ref, computed } from 'vue'
import { useRoute } from 'vue-router'
import { useRouter } from 'vue-router'
  const router = useRouter()
  
  // 获取当前路由对象
  const route = useRoute() 
  
  // 使用 computed 属性动态获取当前激活的菜单项
  // 确保 activeIndex 总是匹配当前路由的路径
  const activeIndex = computed(() => {
    return route.path
  })

  const userName = computed(() => {
    const userInfo = JSON.parse(localStorage.getItem('userInfo') || '{}')
    return userInfo.username
  })

  // 判断是否为管理员
  const isAdmin = computed(() => {
    const userInfo = JSON.parse(localStorage.getItem('userInfo') || '{}')
    return userInfo.is_admin === true
  })

  const handleLogout = () => {
    localStorage.removeItem('userToken')
    localStorage.removeItem('userInfo')
    localStorage.removeItem('userID')
    router.push('/login')
  }
  </script>
  
<style scoped>
.header-container {
  width: 100%;
}

.main-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 20px;
    box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
    background-color: #ffffff;
  }
  
  .logo-text {
    font-size: 20px;
    font-weight: bold;
    color: #403880;
  }
  
  .el-menu-demo {
    flex-grow: 1;
    border-bottom: none !important;
    justify-content: flex-start;
    margin-left: 50px;
  }
  
  .logout {
    display: flex;
    align-items: center;
    cursor: pointer;
  }
  </style>