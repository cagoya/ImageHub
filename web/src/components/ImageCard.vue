<template>
  <el-card class="image-card clickable-card" shadow="hover" :body-style="{ padding: '0px' }" @click="handleCardClick">
    
    <template #header>
      <div class="card-header">
        <span class="card-title">{{ image.image_name }}</span>
        <div class="tags-footer">
      <div class="tags-container">
        <!-- AI标签 -->
        <el-tag
          v-for="tag in (image.ai_tags || [])"
          :key="`ai-${tag}`"
          class="tag-item"
          type="primary"
          effect="plain"
          round
        >
          {{ tag }}
        </el-tag>
        <!-- 手动标签 -->
        <el-tag
          v-for="tag in (image.manual_tags || [])"
          :key="`manual-${tag}`"
          class="tag-item"
          type="success"
          effect="plain"
          round
        >
          {{ tag }}
        </el-tag>
        <el-text 
          v-if="(!image.ai_tags || image.ai_tags.length === 0) && (!image.manual_tags || image.manual_tags.length === 0)" 
          size="small" 
          type="info"
          class="no-tags-text"
        >
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

  <!-- 图片详情对话框 -->
  <ImageDetailDialog v-model="detailDialogVisible" :image-id="props.image.image_id" />

  <!-- 添加标签对话框 -->
  <el-dialog
    v-model="tagDialogVisible"
    title="添加标签"
    :width="isMobile ? '90%' : '500px'"
    :close-on-click-modal="false"
    :append-to-body="true"
    class="tag-dialog"
  >
    <div class="tag-dialog-content">
      <div class="hint-section">
        <el-text class="hint-text" size="small" type="info">
          每个标签不超过5个字，按回车键添加标签
        </el-text>
        <el-button
          type="primary"
          :icon="MagicStick"
          size="small"
          :loading="generatingTags"
          @click="handleGenerateAITags"
          class="generate-ai-tags-btn"
        >
          AI生成标签
        </el-button>
      </div>
      
      <div class="tags-input-container">
        <!-- AI生成的标签（蓝色） -->
        <el-tag
          v-for="tag in currentAITags"
          :key="`ai-${tag}`"
          closable
          :disable-transitions="false"
          @close="handleRemoveAITag(tag)"
          class="tag-item-large"
          type="primary"
        >
          <el-icon class="tag-icon"><MagicStick /></el-icon>
          {{ tag }}
        </el-tag>
        <!-- 手动添加的标签（绿色） -->
        <el-tag
          v-for="tag in currentManualTags"
          :key="`manual-${tag}`"
          closable
          :disable-transitions="false"
          @close="handleRemoveManualTag(tag)"
          class="tag-item-large"
          type="success"
        >
          {{ tag }}
        </el-tag>
        
        <el-input
          v-if="inputVisible"
          ref="inputRef"
          v-model="inputValue"
          class="tag-input"
          size="small"
          maxlength="5"
          @keyup.enter="handleInputConfirm"
          @blur="handleInputConfirm"
          placeholder="输入标签"
        />
        
        <el-button
          v-else
          class="button-new-tag"
          size="small"
          @click="showInput"
        >
          + 新标签
        </el-button>
      </div>
    </div>

    <template #footer>
      <span class="dialog-footer">
        <el-button @click="tagDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSaveTags" :loading="saving">
          保存
        </el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { defineProps, defineEmits, withDefaults, type PropType, computed, ref, nextTick } from 'vue';
import { Picture, InfoFilled, PriceTag, Edit, Delete, Download, MagicStick } from '@element-plus/icons-vue';
import { getImageUrl } from '@/utils/url';
import { ImageAPI } from '@/api';
import { ElMessage, ElMessageBox } from 'element-plus';
import ImageDetailDialog from './ImageDetailDialog.vue';
import type { Image } from '@/type';

// 使用 defineProps 接收 image 数据
const props = withDefaults(defineProps<{
  image: Image;
  userId?: number;
  allowEdit?: boolean;
}>(), {
  userId: undefined,
  allowEdit: true,
});

// 定义 emits
const emit = defineEmits<{
  'delete-success': [imageId: number];
  'edit-image': [imageId: number, imagePath: string];
  'click': [imageId: number];
}>();

// 获取图片 URL，如果提供了 userId 则使用它，否则使用当前登录用户的ID
const imageUrl = computed(() => getImageUrl(props.image.image_path, props.userId));

// 移动端检测
const isMobile = computed(() => {
  if (typeof window === 'undefined') return false;
  return window.innerWidth <= 768;
});

// 定义工具栏按钮
const toolbarButtons = computed(() => {
  const buttons = [
    { name: '查看详细信息', icon: InfoFilled, action: 'view' },
    { name: '添加标签', icon: PriceTag, action: 'addTag' },
    { name: '编辑图片', icon: Edit, action: 'edit' },
    { name: '下载图片', icon: Download, action: 'download' },
    { name: '删除图片', icon: Delete, action: 'delete' },
  ];
  
  // 如果不允许编辑，则过滤掉编辑按钮
  if (!props.allowEdit) {
    return buttons.filter(btn => btn.action !== 'edit');
  }
  
  return buttons;
});

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
    case 'download':
      handleDownload();
      break;
    case 'delete':
      handleDelete();
      break;
  }
};

// 图片详情对话框相关状态
const detailDialogVisible = ref(false);

// 标签对话框相关状态
const tagDialogVisible = ref(false);
const currentAITags = ref<string[]>([]);
const currentManualTags = ref<string[]>([]);
const inputVisible = ref(false);
const inputValue = ref('');
const inputRef = ref();
const saving = ref(false);
const generatingTags = ref(false);

const handleCardClick = () => {
  // 触发点击事件，传递图片 ID 给父组件
  emit('click', props.image.image_id);
};

const handleView = () => {
  detailDialogVisible.value = true;
};

const handleAddTag = () => {
  // 初始化当前标签：AI标签和手动标签分开存储
  currentAITags.value = [...(props.image.ai_tags || [])];
  currentManualTags.value = [...(props.image.manual_tags || [])];
  tagDialogVisible.value = true;
};

// AI生成标签
const handleGenerateAITags = async () => {
  try {
    generatingTags.value = true;
    const response = await ImageAPI.generateTags(props.image.image_id);
    
    if (response.data && response.data.code === 200 && response.data.data) {
      const aiTags = response.data.data;
      
      if (Array.isArray(aiTags) && aiTags.length > 0) {
        // 将AI生成的标签添加到AI标签列表（去重，同时从手动标签中移除）
        const existingAITags = new Set(currentAITags.value);
        const existingManualTags = new Set(currentManualTags.value);
        
        aiTags.forEach((tag: string) => {
          if (tag && !existingAITags.has(tag)) {
            currentAITags.value.push(tag);
            existingAITags.add(tag);
            // 如果该标签在手动标签中，则移除（AI标签优先级更高）
            if (existingManualTags.has(tag)) {
              currentManualTags.value = currentManualTags.value.filter(t => t !== tag);
            }
          }
        });
        ElMessage.success(`已生成 ${aiTags.length} 个标签`);
      } else {
        ElMessage.warning('AI未生成标签');
      }
    } else {
      ElMessage.error(response.data?.message || '生成标签失败');
    }
  } catch (error: any) {
    console.error('生成标签失败:', error);
    ElMessage.error(error.message || '生成标签失败，请重试');
  } finally {
    generatingTags.value = false;
  }
};

const handleEdit = () => {
  // 触发编辑事件，将图片信息传递给父组件
  emit('edit-image', props.image.image_id, props.image.image_path);
};

const handleDownload = async () => {
  try {
    // 调用 API 获取图片文件
    const response = await ImageAPI.getFile(props.image.image_id);
    
    // 创建 Blob 对象
    const blob = new Blob([response.data], { type: 'image/jpeg' });
    
    // 创建下载链接
    const url = window.URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    
    // 设置下载文件名
    const fileName = props.image.image_name || `image_${props.image.image_id}.jpg`;
    link.download = fileName;
    
    // 触发下载
    document.body.appendChild(link);
    link.click();
    
    // 清理
    document.body.removeChild(link);
    window.URL.revokeObjectURL(url);
    
    ElMessage.success('图片下载成功');
  } catch (error) {
    ElMessage.error('图片下载失败');
    console.error(error);
  }
};

// 显示输入框
const showInput = () => {
  inputVisible.value = true;
  nextTick(() => {
    inputRef.value?.focus();
  });
};

// 移除AI标签
const handleRemoveAITag = (tag: string) => {
  currentAITags.value = currentAITags.value.filter(t => t !== tag);
};

// 移除手动标签
const handleRemoveManualTag = (tag: string) => {
  currentManualTags.value = currentManualTags.value.filter(t => t !== tag);
};

// 确认输入
const handleInputConfirm = () => {
  const value = inputValue.value.trim();
  
  if (value) {
    // 验证标签长度
    if (value.length > 5) {
      ElMessage.warning('标签不能超过5个字');
      return;
    }
    
    // 检查是否已存在（在AI标签或手动标签中）
    if (currentAITags.value.includes(value)) {
      ElMessage.warning('该标签已作为AI标签存在');
    } else if (currentManualTags.value.includes(value)) {
      ElMessage.warning('标签已存在');
    } else {
      // 添加到手动标签
      currentManualTags.value.push(value);
    }
  }
  
  inputVisible.value = false;
  inputValue.value = '';
};

// 保存标签
const handleSaveTags = async () => {
  try {
    saving.value = true;
    // 使用对话框中的AI标签和手动标签
    const aiTags = currentAITags.value;
    const manualTags = currentManualTags.value;
    
    await ImageAPI.setTags(props.image.image_id, aiTags, manualTags);
    
    // 更新本地数据
    props.image.ai_tags = [...aiTags];
    props.image.manual_tags = [...manualTags];
    
    ElMessage.success('标签保存成功');
    tagDialogVisible.value = false;
  } catch (error) {
    ElMessage.error('标签保存失败');
    console.error(error);
  } finally {
    saving.value = false;
  }
};

const handleDelete = () => {
  // 使用 Element Plus 的确认对话框
  ElMessageBox.confirm(
    `确定要删除图片 "${props.image.image_name}" 吗？此操作不可撤销。`,
    '删除确认',
    {
      confirmButtonText: '确定删除',
      cancelButtonText: '取消',
      type: 'warning',
      confirmButtonClass: 'el-button--danger',
    }
  ).then(async () => {
    try {
      const userID = localStorage.getItem('userID');
      await ImageAPI.delete(props.image.image_id, parseInt(userID || '0'));
      ElMessage.success('删除图片成功');
      // 通知父组件删除成功，传递图片ID
      emit('delete-success', props.image.image_id);
    } catch (error) {
      ElMessage.error('删除图片失败');
    }
  }).catch(() => {
    // 用户取消删除
    ElMessage.info('已取消删除');
  });
};
</script>

<style scoped>
/* 限制卡片宽度和高度 */
.image-card {
  width: 250px;
  height: 300px;
  margin: 20px;
}

/* 点击效果 */
.clickable-card {
  cursor: pointer;
  transition: transform 0.2s, box-shadow 0.2s;
}

.clickable-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.15) !important;
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

.tags-container {
  display: flex;
  flex-wrap: nowrap; 
  gap: 6px;
  height: 34px;
  min-height: 34px;
  overflow-x: auto;
  overflow-y: hidden;
  align-items: center;
  white-space: nowrap;
  width: 100%;
  -webkit-overflow-scrolling: touch;
  scrollbar-width: thin;
  scrollbar-color: rgba(64, 158, 255, 0.5) transparent;
}

.tags-container::-webkit-scrollbar {
  height: 4px;
}

.tags-container::-webkit-scrollbar-track {
  background: transparent;
}

.tags-container::-webkit-scrollbar-thumb {
  background-color: rgba(64, 158, 255, 0.5);
  border-radius: 2px;
}

.tags-container::-webkit-scrollbar-thumb:hover {
  background-color: rgba(64, 158, 255, 0.7);
}

.tag-item {
  margin: 0;
  flex-shrink: 0;
  white-space: nowrap;
  overflow: visible;
}

.no-tags-text {
  flex-shrink: 0;
  line-height: 34px;
  height: 34px;
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

/* 标签对话框样式 */
.tag-dialog-content {
  padding: 10px 0;
}

.hint-section {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  gap: 12px;
}

.hint-text {
  flex: 1;
  color: #909399;
}

.generate-ai-tags-btn {
  flex-shrink: 0;
  white-space: nowrap;
}

.tags-input-container {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  align-items: center;
  min-height: 80px;
  padding: 12px;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  background-color: #f5f7fa;
}

.tag-item-large {
  font-size: 14px;
  height: 28px;
  line-height: 26px;
}

.tag-item-large .tag-icon {
  margin-right: 4px;
  font-size: 12px;
}

.tag-input {
  width: 120px;
}

.button-new-tag {
  height: 28px;
  padding: 0 12px;
  border-style: dashed;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

/* 标签对话框移动端优化 */
@media (max-width: 768px) {
  .tag-dialog-content {
    padding: 8px 0;
  }

  .hint-section {
    flex-direction: column;
    align-items: flex-start;
    gap: 10px;
    margin-bottom: 12px;
  }

  .hint-text {
    font-size: 13px;
    width: 100%;
  }

  .generate-ai-tags-btn {
    width: 100%;
  }

  .tags-input-container {
    min-height: 100px;
    padding: 10px;
    gap: 8px;
  }

  .tag-item-large {
    font-size: 13px;
    height: 32px;
    line-height: 30px;
    padding: 0 10px;
  }

  .tag-input {
    width: 100%;
    max-width: 150px;
  }

  .button-new-tag {
    height: 32px;
    padding: 0 16px;
    font-size: 14px;
  }

  .dialog-footer {
    flex-direction: row;
    justify-content: center;
    gap: 12px;
  }

  .dialog-footer .el-button {
    flex: 0 0 auto;
    min-width: 100px;
    max-width: 140px;
    padding: 10px 20px;
    font-size: 14px;
  }
}

@media (max-width: 480px) {
  .tag-dialog-content {
    padding: 6px 0;
  }

  .hint-section {
    gap: 8px;
    margin-bottom: 10px;
  }

  .hint-text {
    font-size: 12px;
  }

  .generate-ai-tags-btn {
    font-size: 12px;
    padding: 8px 12px;
  }

  .tags-input-container {
    min-height: 90px;
    padding: 8px;
    gap: 6px;
  }

  .tag-item-large {
    font-size: 12px;
    height: 28px;
    line-height: 26px;
    padding: 0 8px;
  }

  .tag-input {
    width: 100%;
    max-width: 120px;
  }

  .button-new-tag {
    height: 28px;
    padding: 0 12px;
    font-size: 13px;
  }

  .dialog-footer {
    gap: 10px;
  }

  .dialog-footer .el-button {
    min-width: 90px;
    max-width: 120px;
    padding: 10px 16px;
    font-size: 13px;
  }
}

/* 对话框移动端优化 */
:deep(.tag-dialog .el-dialog) {
  border-radius: 12px;
}

:deep(.tag-dialog .el-dialog__header) {
  padding: 16px 20px;
  border-bottom: 1px solid #ebeef5;
}

:deep(.tag-dialog .el-dialog__title) {
  font-size: 18px;
  font-weight: 600;
}

:deep(.tag-dialog .el-dialog__body) {
  padding: 20px;
}

@media (max-width: 768px) {
  :deep(.tag-dialog .el-dialog) {
    width: 90% !important;
    margin: 5vh auto;
    border-radius: 10px;
  }

  :deep(.tag-dialog .el-dialog__header) {
    padding: 12px 16px;
  }

  :deep(.tag-dialog .el-dialog__title) {
    font-size: 16px;
  }

  :deep(.tag-dialog .el-dialog__body) {
    padding: 16px;
  }

  :deep(.tag-dialog .el-dialog__footer) {
    padding: 12px 16px;
  }
}

@media (max-width: 480px) {
  :deep(.tag-dialog .el-dialog) {
    width: 95% !important;
    margin: 3vh auto;
    border-radius: 8px;
  }

  :deep(.tag-dialog .el-dialog__header) {
    padding: 10px 12px;
  }

  :deep(.tag-dialog .el-dialog__title) {
    font-size: 15px;
  }

  :deep(.tag-dialog .el-dialog__body) {
    padding: 12px;
  }

  :deep(.tag-dialog .el-dialog__footer) {
    padding: 10px 12px;
  }
}
</style>