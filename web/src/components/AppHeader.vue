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
        </el-menu>
  
        <div class="user-info">
            <el-avatar :size="30" :icon="UserFilled" />
            <span class="username">{{ userName }}</span>
            <el-button type="text" @click="handleLogout">退出登录</el-button>
        </div>
  
      </el-header>
      
    </el-container>
  </template>
  
  <script lang="ts" setup>
  import { ref, computed } from 'vue'
  import { useRoute } from 'vue-router'
  import { UserFilled } from '@element-plus/icons-vue'
  import { useRouter } from 'vue-router'
  const router = useRouter()
  
  // 获取当前路由对象
  const route = useRoute() 
  
  // 使用 computed 属性动态获取当前激活的菜单项
  // 确保 activeIndex 总是匹配当前路由的路径
  const activeIndex = computed(() => {
    return route.path
  })
  
  // 在 router 模式下，不需要手动处理 handleSelect
  // 因为 ElMenu 会自动使用 index 作为路径进行路由跳转

  const userName = computed(() => {
    const userInfo = JSON.parse(localStorage.getItem('userInfo') || '{}')
    return userInfo.username
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
  
  .user-info {
    display: flex;
    align-items: center;
    cursor: pointer;
  }
  
  .username {
    margin-left: 8px;
    color: #606266;
    font-size: 14px;
  }
  
  .el-main {
    padding: 20px;
    background-color: #f0f2f5;
  }
  </style>