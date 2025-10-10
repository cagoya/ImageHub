<template>
  <el-card class="image-card" shadow="hover" :body-style="{ padding: '0px' }">
    
    <template #header>
      <div class="card-header">
        <el-text size="small" type="info">ID: {{ image.image_id }}</el-text>
        <span class="card-title">{{ image.image_name }}</span>
        <div class="tags-footer">
      <el-text class="tags-label" size="small">标签:</el-text>
      <div class="tags-container">
        <el-tag
          v-for="tag in image.tags"
          :key="tag"
          class="tag-item"
          type="info"
          effect="plain"
          round
        >
          {{ tag }}
        </el-tag>
        <el-text v-if="image.tags.length === 0" size="small" type="danger">
          暂无标签
        </el-text>
      </div>
    </div>
      </div>
    </template>

    <div class="image-body-container">
      <el-image 
        class="image-display"
        :src="imageUrl" 
        fit="cover" 
        loading="lazy"
      >
        <template #error>
          <div class="image-slot">
            <el-icon><Picture /></el-icon>
            <el-text size="small" type="info">图片加载失败</el-text>
          </div>
        </template>
      </el-image>

      <div class="toolbar">
        <el-button 
          v-for="btn in toolbarButtons"
          :key="btn.name"
          :icon="btn.icon"
          :title="btn.name"
          size="small"
          circle
          @click.stop="handleAction(btn.action)"
        />
      </div>
    </div>


  </el-card>
</template>

<script setup lang="ts">
import { defineProps, defineEmits, type PropType, computed } from 'vue';
import { Picture, InfoFilled, PriceTag, Edit, Delete } from '@element-plus/icons-vue';
import { getImageUrl } from '@/utils/url';
import { ImageAPI } from '@/api';
import { ElMessage } from 'element-plus';

interface Image {
    image_id: number;
    image_name: string;
    image_path: string;
    tags: string[];
}

// 使用 defineProps 接收 image 数据
const props = defineProps({
  image: {
    type: Object as PropType<Image>,
    required: true,
  },
});

// 定义 emits
const emit = defineEmits<{
  'delete-success': [];
}>();

// 获取图片 URL
const imageUrl = computed(() => getImageUrl(props.image.image_path));

console.log(props.image);

// 定义工具栏按钮
const toolbarButtons = [
  { name: '查看详细信息', icon: InfoFilled, action: 'view' },
  { name: '添加标签', icon: PriceTag, action: 'addTag' },
  { name: '编辑图片', icon: Edit, action: 'edit' },
  { name: '删除图片', icon: Delete, action: 'delete' },
];

// 处理按钮点击事件（这里仅为演示，实际应用中你可能需要使用 emit 向上级组件报告）
const handleAction = (action: string) => {

  // 示例:
  switch (action) {
    case 'view':
      handleView();
      break;
    case 'addTag':
      handleAddTag();
      break;
    case 'edit':
      handleEdit();
      break;
    case 'delete':
      handleDelete();
      break;
  }
};

const handleView = () => {
  // TODO: 实现查看详细信息功能
  console.log('查看详细信息');
};

const handleAddTag = () => {
  // TODO: 实现添加标签功能
  console.log('添加标签');
};

const handleEdit = () => {
  // TODO: 实现编辑图片功能
  console.log('编辑图片');
};

const handleDelete = async () => {
  try {
    const userID = localStorage.getItem('userID');
    await ImageAPI.delete(props.image.image_id, parseInt(userID || '0'));
    ElMessage.success('删除图片成功');
    // 通知父组件刷新列表
    emit('delete-success');
  } catch (error) {
    ElMessage.error('删除图片失败');
  }
};
</script>

<style scoped>
/* 限制卡片宽度和高度 */
.image-card {
  width: 300px;
  height: 400px;
  margin: 20px;
}

.card-header {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.card-title {
  font-size: 16px;
  font-weight: bold;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.image-body-container {
  position: relative; /* 为绝对定位的工具栏提供定位参照 */
  height: 250px;
  overflow: hidden;
}

.image-display {
  width: 100%;
  height: 100%;
  display: block;
}

.image-slot {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  width: 100%;
  height: 100%;
  background: var(--el-fill-color-light);
  color: var(--el-text-color-secondary);
}

.tags-footer {
  padding: 0px;
}

.tags-label {
  margin-bottom: 6px;
  display: block;
  font-weight: 600;
}

.tags-container {
  display: flex;
  flex-wrap: wrap; 
  gap: 6px;
  max-height: 40px;
  overflow: hidden;
}

.tag-item {
  margin: 0; 
}

.toolbar {
  /* 🌟 重点：设置绝对定位 */
  position: absolute;
  /* 🌟 重点：定位到右上角 */
  top: 8px; /* 距离顶部的间距 */
  right: 8px; /* 距离右侧的间距 */
  
  /* 布局和外观优化 */
  display: flex; /* 让按钮水平排列 */
  gap: 4px; /* 按钮之间的间距 */
  padding: 4px;
  background-color: rgba(0, 0, 0, 0.4); /* 半透明背景，更易于查看 */
  border-radius: 4px;
  z-index: 10; /* 确保工具栏在图片上方显示 */

  /* 悬停效果 (可选) */
  opacity: 0; /* 默认隐藏 */
  transition: opacity 0.3s ease-in-out;
}

/* 鼠标悬停在卡片上时显示工具栏 */
.image-card:hover .toolbar {
  opacity: 1; /* 悬停时显示 */
}
</style>