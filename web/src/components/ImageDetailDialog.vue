<template>
  <el-dialog
    v-model="visible"
    title="图片详细信息"
    :width="isMobile ? '95%' : '80%'"
    :close-on-click-modal="false"
    :append-to-body="true"
    :modal="true"
    @close="handleClose"
    class="detail-dialog"
  >
    <div v-loading="loading" class="detail-container">
      <div v-if="!loading && imageDetail" class="detail-content">
        <!-- 左侧原图展示 -->
        <div class="image-section">
          <el-image 
            :src="originImageUrl" 
            fit="contain"
            class="origin-image"
          >
            <template #error>
              <div class="image-error">
                <el-icon size="64"><Picture /></el-icon>
                <el-text type="info">原图加载失败</el-text>
              </div>
            </template>
          </el-image>
        </div>

        <!-- 右侧详细信息 -->
        <div class="info-section">
          <el-scrollbar :height="isMobile ? 'auto' : '600px'">
            <div class="info-content">
              <el-descriptions :column="isMobile ? 1 : 1" :border="!isMobile" :size="isMobile ? 'small' : 'default'">
                <el-descriptions-item label="图片ID">
                  {{ imageDetail.id }}
                </el-descriptions-item>
                
                <el-descriptions-item label="图片名称">
                  {{ imageDetail.name }}
                </el-descriptions-item>
                
                <el-descriptions-item label="拍摄时间">
                  {{ formatTimeDisplay(imageDetail.captured_time) }}
                </el-descriptions-item>
                
                <el-descriptions-item label="图片尺寸">
                  <span v-if="imageDetail.width_pixel && imageDetail.height_pixel">
                    {{ imageDetail.width_pixel }} × {{ imageDetail.height_pixel }} 像素
                  </span>
                  <el-text v-else type="info">无尺寸信息</el-text>
                </el-descriptions-item>
                
                <el-descriptions-item label="地理坐标">
                  <span v-if="imageDetail.latitude && imageDetail.longitude" class="coordinate-inline">
                    <span>纬度: {{ imageDetail.latitude.toFixed(6) }}</span>
                    <span class="coordinate-separator">|</span>
                    <span>经度: {{ imageDetail.longitude.toFixed(6) }}</span>
                  </span>
                  <el-text v-else type="info">无地理坐标</el-text>
                </el-descriptions-item>
                
                <el-descriptions-item label="详细地址">
                  {{ imageDetail.formatted_address || '无地址信息' }}
                </el-descriptions-item>
                
                <el-descriptions-item :label="isMobile ? '行政区域' : '省份'">
                  <template v-if="isMobile">
                    <span v-if="imageDetail.province || imageDetail.city || imageDetail.district">
                      {{ [imageDetail.province, imageDetail.city, imageDetail.district].filter(Boolean).join(' / ') || '未知' }}
                    </span>
                    <span v-else>未知</span>
                  </template>
                  <template v-else>
                    {{ imageDetail.province || '未知' }}
                  </template>
                </el-descriptions-item>
                
                <el-descriptions-item v-if="!isMobile" label="城市">
                  {{ imageDetail.city || '未知' }}
                </el-descriptions-item>
                
                <el-descriptions-item v-if="!isMobile" label="区县">
                  {{ imageDetail.district || '未知' }}
                </el-descriptions-item>
                
                <el-descriptions-item label="上传时间">
                  {{ formatTimeDisplay(imageDetail.created_time) }}
                </el-descriptions-item>
                
                <el-descriptions-item label="更新时间">
                  {{ formatTimeDisplay(imageDetail.updated_time) }}
                </el-descriptions-item>
              </el-descriptions>
            </div>
          </el-scrollbar>
        </div>
      </div>
    </div>

    <template #footer>
      <span class="dialog-footer">
        <el-button @click="handleClose">关闭</el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue';
import { Picture } from '@element-plus/icons-vue';
import { ImageAPI } from '@/api';
import { getOriginImageUrl } from '@/utils/url';
import { ElMessage } from 'element-plus';
import type { ImageDetail } from '@/type';
import { formatTime } from '@/utils/date';

// Props
const props = defineProps<{
  modelValue: boolean;
  imageId: number;
}>();

// Emits
const emit = defineEmits<{
  'update:modelValue': [value: boolean];
}>();

// 状态
const loading = ref(false);
const imageDetail = ref<ImageDetail | null>(null);

// 计算属性
const visible = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
});

const originImageUrl = computed(() => {
  if (imageDetail.value?.path) {
    return getOriginImageUrl(imageDetail.value.path);
  }
  return '';
});

// 移动端检测
const isMobile = computed(() => {
  if (typeof window === 'undefined') return false;
  return window.innerWidth <= 768;
});

// 监听对话框打开，加载详细信息
watch(visible, async (newVal) => {
  if (newVal && props.imageId) {
    await loadImageDetail();
  }
});

// 加载图片详细信息
const loadImageDetail = async () => {
  try {
    loading.value = true;
    const response = await ImageAPI.getDetail(props.imageId);
    imageDetail.value = response.data.data;
  } catch (error) {
    ElMessage.error('加载图片详细信息失败');
    console.error(error);
    handleClose();
  } finally {
    loading.value = false;
  }
};

// 格式化时间（使用工具函数）
const formatTimeDisplay = (time?: string) => {
  return formatTime(time);
};

// 关闭对话框
const handleClose = () => {
  visible.value = false;
  // 清空数据
  setTimeout(() => {
    imageDetail.value = null;
  }, 300);
};
</script>

<style scoped>
.detail-container {
  min-height: 600px;
}

.detail-content {
  display: flex;
  gap: 20px;
  height: 600px;
}

.image-section {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f5f7fa;
  border-radius: 4px;
  padding: 20px;
  overflow: hidden;
}

.origin-image {
  max-width: 100%;
  max-height: 100%;
}

.image-error {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: #909399;
  gap: 10px;
}

.info-section {
  flex: 0 0 400px;
  overflow: hidden;
}

.info-content {
  padding: 0 10px 10px 0;
}

/* 坐标显示优化 - 统一PC端和移动端 */
.coordinate-inline {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  white-space: nowrap !important;
  word-break: keep-all;
}

.coordinate-inline > span {
  white-space: nowrap;
  flex-shrink: 0;
}

.coordinate-separator {
  color: #dcdfe6;
  margin: 0 4px;
  flex-shrink: 0;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
}

/* PC端描述列表优化 */
:deep(.el-descriptions__label) {
  width: 100px;
  min-width: 100px;
}

/* 响应式设计 */
@media (max-width: 1200px) {
  .detail-content {
    flex-direction: column;
    height: auto;
  }
  
  .image-section {
    height: 400px;
  }
  
  .info-section {
    flex: 1;
  }
}

/* 移动端优化 */
@media (max-width: 768px) {
  .detail-container {
    min-height: auto;
  }

  .detail-content {
    flex-direction: column;
    gap: 16px;
    height: auto;
  }

  .image-section {
    height: 300px;
    padding: 12px;
    border-radius: 8px;
  }

  .origin-image {
    max-height: 280px;
  }

  .image-error {
    height: 280px;
  }

  .image-error .el-icon {
    font-size: 48px;
  }

  .info-section {
    flex: 1;
    width: 100%;
  }

  .info-content {
    padding: 0;
  }

  .dialog-footer {
    padding: 12px 0;
  }

  .dialog-footer .el-button {
    width: 100%;
    padding: 12px;
  }

  /* 移动端描述列表优化 */
  :deep(.el-descriptions) {
    font-size: 14px;
  }

  :deep(.el-descriptions__label) {
    font-size: 13px;
    font-weight: 600;
    width: 90px;
    padding: 8px 12px;
  }

  :deep(.el-descriptions__content) {
    font-size: 13px;
    padding: 8px 12px;
    word-break: break-word;
  }

  :deep(.el-descriptions__body) {
    background-color: #fff;
  }

  :deep(.el-descriptions__table) {
    width: 100%;
  }

  :deep(.el-descriptions__cell) {
    padding: 0;
  }
}

/* 小屏幕进一步优化 */
@media (max-width: 480px) {
  .detail-content {
    gap: 12px;
  }

  .image-section {
    height: 250px;
    padding: 10px;
  }

  .origin-image {
    max-height: 230px;
  }

  .image-error {
    height: 230px;
  }

  .image-error .el-icon {
    font-size: 40px;
  }

  :deep(.el-descriptions__label) {
    font-size: 12px;
    width: 80px;
    padding: 6px 10px;
  }

  :deep(.el-descriptions__content) {
    font-size: 12px;
    padding: 6px 10px;
  }

  .dialog-footer .el-button {
    padding: 10px;
    font-size: 14px;
  }
}

/* 对话框移动端优化 */
:deep(.detail-dialog .el-dialog) {
  margin: 5vh auto;
  border-radius: 12px;
  z-index: 3000 !important;
}

/* 确保对话框遮罩层也在最上层 */
:deep(.detail-dialog .el-overlay) {
  z-index: 2999 !important;
}

:deep(.detail-dialog .el-overlay-dialog) {
  z-index: 3000 !important;
}

:deep(.detail-dialog .el-dialog__header) {
  padding: 16px 20px;
  border-bottom: 1px solid #ebeef5;
}

:deep(.detail-dialog .el-dialog__title) {
  font-size: 18px;
  font-weight: 600;
}

:deep(.detail-dialog .el-dialog__body) {
  padding: 20px;
}

@media (max-width: 768px) {
  :deep(.detail-dialog .el-dialog) {
    width: 95% !important;
    margin: 2vh auto;
    border-radius: 10px;
    max-height: 90vh;
    display: flex;
    flex-direction: column;
  }

  :deep(.detail-dialog .el-dialog__header) {
    padding: 12px 16px;
    flex-shrink: 0;
  }

  :deep(.detail-dialog .el-dialog__title) {
    font-size: 16px;
  }

  :deep(.detail-dialog .el-dialog__body) {
    padding: 16px;
    flex: 1;
    overflow: hidden;
    display: flex;
    flex-direction: column;
  }

  :deep(.detail-dialog .el-dialog__footer) {
    padding: 12px 16px;
    flex-shrink: 0;
  }
}

@media (max-width: 480px) {
  :deep(.detail-dialog .el-dialog) {
    width: 98% !important;
    margin: 1vh auto;
    border-radius: 8px;
  }

  :deep(.detail-dialog .el-dialog__header) {
    padding: 10px 12px;
  }

  :deep(.detail-dialog .el-dialog__title) {
    font-size: 15px;
  }

  :deep(.detail-dialog .el-dialog__body) {
    padding: 12px;
  }

  :deep(.detail-dialog .el-dialog__footer) {
    padding: 10px 12px;
  }
}
</style>

