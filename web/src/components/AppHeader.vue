<template>
  <el-container class="header-container">
    <el-header height="60px" class="main-header">
      <el-menu
        :default-active="activeIndex"
        mode="horizontal"
        :router="true"
        class="el-menu-desktop"
        :class="{ 'hidden-sm': isMobile }"
      >
        <el-menu-item index="/dashboard">首页</el-menu-item>
        <el-menu-item index="/upload">上传图片</el-menu-item>
        <el-menu-item index="/gallery">查看图片</el-menu-item>
        <el-menu-item index="/chat">AI搜图</el-menu-item>
        <el-menu-item v-if="isAdmin" index="/manage">用户管理</el-menu-item>
        <el-menu-item index="/config">设置</el-menu-item>
      </el-menu>

      <div class="mobile-menu-toggle" v-if="isMobile">
        <el-icon :size="24" @click="drawerVisible = true">
          <Menu />
        </el-icon>
      </div>

      <div class="logout">
        <el-button type="text" @click="handleLogout">退出登录</el-button>
      </div>
    </el-header>

    <el-drawer
      v-model="drawerVisible"
      title="导航"
      direction="ltr"
      :size="isMobile ? '70%' : '60%'"
      :with-header="false"
      class="mobile-drawer"
      :z-index="3000"
    >
      <div class="mobile-menu-container">
        <el-menu
          :default-active="activeIndex"
          :router="true"
          class="el-menu-mobile"
          @select="handleMenuSelect"
        >
          <el-menu-item index="/dashboard">
            <span>首页</span>
          </el-menu-item>
          <el-menu-item index="/upload">
            <span>上传图片</span>
          </el-menu-item>
          <el-menu-item index="/gallery">
            <span>查看图片</span>
          </el-menu-item>
          <el-menu-item index="/chat">
            <span>AI搜图</span>
          </el-menu-item>
          <el-menu-item v-if="isAdmin" index="/manage">
            <span>用户管理</span>
          </el-menu-item>
          <el-menu-item index="/config">
            <span>设置</span>
          </el-menu-item>
        </el-menu>
      </div>
    </el-drawer>
  </el-container>
</template>

<script lang="ts" setup>
import { ref, computed, onMounted, onBeforeUnmount } from 'vue';
import { useRoute, useRouter } from 'vue-router';
// 导入 Element Plus 汉堡包图标
import { Menu } from '@element-plus/icons-vue';

const router = useRouter();
const route = useRoute();

// --- 响应式和菜单状态 ---
const activeIndex = computed(() => route.path);
const drawerVisible = ref(false); // 控制抽屉菜单的显示/隐藏
const isMobile = ref(false); // 判断是否为手机屏幕

// --- 用户信息状态 ---
const userInfo = ref(JSON.parse(localStorage.getItem('userInfo') || '{}'));

// --- 用户信息计算属性 ---
const isAdmin = computed(() => {
  return userInfo.value.is_admin === true;
});

// 监听用户信息更新事件
const handleUserInfoUpdate = (event: CustomEvent) => {
  console.log('收到用户信息更新事件，刷新界面');
  userInfo.value = event.detail;
  // 触发Vue的响应式更新
  userInfo.value = { ...event.detail };
};

// --- 方法 ---

// 监听页面宽度变化
const checkMobile = () => {
  // 假设宽度小于 768px 为手机端，您可以根据实际设计调整
  isMobile.value = window.innerWidth < 768;
};

// 菜单项被选中时关闭抽屉
const handleMenuSelect = (index: string) => {
  console.log('移动端菜单选择:', index);
  // 延迟关闭抽屉，确保路由跳转先执行
  setTimeout(() => {
    drawerVisible.value = false;
  }, 100);
  
  // 如果路由没有自动跳转，手动跳转
  if (route.path !== index) {
    router.push(index);
  }
};

const handleLogout = () => {
  localStorage.removeItem('userToken');
  localStorage.removeItem('userInfo');
  localStorage.removeItem('userID');
  router.push('/login');
};

// --- 生命周期钩子 ---
onMounted(() => {
  checkMobile(); // 初始化检查一次
  window.addEventListener('resize', checkMobile); // 监听窗口大小变化
  // 监听用户信息更新事件
  window.addEventListener('userInfoUpdated', handleUserInfoUpdate as EventListener);
  // 初始化用户信息
  userInfo.value = JSON.parse(localStorage.getItem('userInfo') || '{}');
});

onBeforeUnmount(() => {
  window.removeEventListener('resize', checkMobile); // 组件销毁时移除监听
  window.removeEventListener('userInfoUpdated', handleUserInfoUpdate as EventListener);
});
</script>

<style scoped>
/* 容器和头部样式保持不变 */
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

/* 桌面端菜单样式 */
.el-menu-desktop {
  flex-grow: 1;
  border-bottom: none !important;
  justify-content: flex-start;
  margin-left: 50px;
}

/* 隐藏桌面端菜单的样式类 */
.hidden-sm {
  display: none !important;
}

/* 手机端汉堡包图标容器 */
.mobile-menu-toggle {
  display: none; /* 默认隐藏 */
  margin-right: auto; /* 推到左边 */
  cursor: pointer;
}

.logout {
  display: flex;
  align-items: center;
  /* 确保退出按钮始终在右边 */
}

/* 媒体查询：应用于屏幕宽度小于 768px 的设备 */
@media (max-width: 767.98px) {
  /* 隐藏桌面端菜单（已通过 JS/v-if 实现，这里作为保险） */
  .el-menu-desktop {
    display: none;
  }

  /* 显示汉堡包图标 */
  .mobile-menu-toggle {
    display: flex;
    align-items: center;
    /* 调整菜单和退出按钮之间的空间 */
    margin-left: 20px;
  }
}

/* 手机端抽屉菜单容器 */
.mobile-menu-container {
  width: 100%;
  height: 100%;
  overflow-y: auto;
}

/* 手机端抽屉菜单的样式 */
.el-menu-mobile {
  border-right: none; /* 移除垂直菜单的右边框 */
  padding-top: 20px;
  width: 100%;
  background-color: transparent;
}

/* 确保移动端菜单项可点击 */
.el-menu-mobile :deep(.el-menu-item) {
  height: 48px;
  line-height: 48px;
  padding: 0 20px;
  margin: 4px 12px;
  border-radius: 8px;
  cursor: pointer;
  -webkit-tap-highlight-color: rgba(0, 0, 0, 0.1);
  user-select: none;
  transition: all 0.3s ease;
  font-size: 15px;
}

.el-menu-mobile :deep(.el-menu-item:hover) {
  background-color: #f5f7fa;
}

.el-menu-mobile :deep(.el-menu-item.is-active) {
  color: #409EFF;
  background-color: #ecf5ff;
  font-weight: 500;
}

.el-menu-mobile :deep(.el-menu-item span) {
  display: block;
  width: 100%;
}

/* 确保 drawer 在移动端正确显示 */
:deep(.mobile-drawer .el-drawer) {
  z-index: 3000;
}

:deep(.mobile-drawer .el-drawer__body) {
  padding: 0;
  overflow: hidden;
}
</style>