<template>
  <el-dialog
    v-model="visible"
    title="图片详细信息"
    width="80%"
    :close-on-click-modal="false"
    @close="handleClose"
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
          <el-scrollbar height="600px">
            <div class="info-content">
              <el-descriptions :column="1" border>
                <el-descriptions-item label="图片ID">
                  {{ imageDetail.id }}
                </el-descriptions-item>
                
                <el-descriptions-item label="图片名称">
                  {{ imageDetail.name }}
                </el-descriptions-item>
                
                <el-descriptions-item label="文件路径">
                  {{ imageDetail.path }}
                </el-descriptions-item>
                
                <el-descriptions-item label="拍摄时间">
                  {{ formatTime(imageDetail.captured_time) }}
                </el-descriptions-item>
                
                <el-descriptions-item label="图片尺寸">
                  <span v-if="imageDetail.width_pixel && imageDetail.height_pixel">
                    {{ imageDetail.width_pixel }} × {{ imageDetail.height_pixel }} 像素
                  </span>
                  <el-text v-else type="info">无尺寸信息</el-text>
                </el-descriptions-item>
                
                <el-descriptions-item label="地理坐标">
                  <div v-if="imageDetail.latitude && imageDetail.longitude">
                    <div>纬度: {{ imageDetail.latitude.toFixed(6) }}</div>
                    <div>经度: {{ imageDetail.longitude.toFixed(6) }}</div>
                  </div>
                  <el-text v-else type="info">无地理坐标</el-text>
                </el-descriptions-item>
                
                <el-descriptions-item label="详细地址">
                  {{ imageDetail.formatted_address || '无地址信息' }}
                </el-descriptions-item>
                
                <el-descriptions-item label="省份">
                  {{ imageDetail.province || '未知' }}
                </el-descriptions-item>
                
                <el-descriptions-item label="城市">
                  {{ imageDetail.city || '未知' }}
                </el-descriptions-item>
                
                <el-descriptions-item label="区县">
                  {{ imageDetail.district || '未知' }}
                </el-descriptions-item>
                
                <el-descriptions-item label="上传时间">
                  {{ formatTime(imageDetail.created_time) }}
                </el-descriptions-item>
                
                <el-descriptions-item label="更新时间">
                  {{ formatTime(imageDetail.updated_time) }}
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

// 图片详细信息类型定义
interface ImageDetail {
  id: number;
  name: string;
  path: string;
  captured_time?: string;
  latitude?: number;
  longitude?: number;
  formatted_address?: string;
  province?: string;
  city?: string;
  district?: string;
  width_pixel?: number;
  height_pixel?: number;
  user_id: number;
  created_time: string;
  updated_time: string;
}

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

// 格式化时间
const formatTime = (time?: string) => {
  if (!time) return '无';
  
  try {
    const date = new Date(time);
    return date.toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit',
      hour12: false
    });
  } catch {
    return time;
  }
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

.dialog-footer {
  display: flex;
  justify-content: flex-end;
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
</style>

