<template>
  <AppHeader />
    <div class="page-container">
      <div class="image-search-container">
        <!-- 搜索表单 -->
        <div class="search-header">
          <h1>图片库</h1>
          <p>查看和搜索您的图片</p>
        </div>

        <el-form class="search-form" inline>
          <div class="search-controls">
            <el-form-item label="查询类型" class="control-group">
              <el-select 
                v-model="searchType" 
                class="search-select" 
                placeholder="选择查询类型"
                @change="handleTypeChange"
              >
                <el-option label="个人所有" value="all" />
                <el-option label="标签查询" value="tag" />
                <el-option label="时间查询" value="time" />
                <el-option label="省份查询" value="province" />
                <el-option label="城市查询" value="city" />
                <el-option label="区县查询" value="district" />
              </el-select>
            </el-form-item>

            <!-- 标签查询 -->
            <el-form-item 
              v-if="searchType === 'tag'" 
              label="标签" 
              class="control-group"
            >
              <el-input
                v-model="searchValue"
                class="search-input"
                placeholder="输入标签"
                clearable
              />
            </el-form-item>

            <!-- 时间查询 -->
            <el-form-item 
              v-if="searchType === 'time'" 
              label="时间范围" 
              class="control-group"
            >
              <el-date-picker
                v-model="timeRange"
                type="datetimerange"
                range-separator="至"
                start-placeholder="开始时间"
                end-placeholder="结束时间"
                value-format="YYYY-MM-DD HH:mm:ss"
                class="time-picker"
              />
            </el-form-item>

            <!-- 省份查询 -->
            <el-form-item 
              v-if="searchType === 'province'" 
              label="省份" 
              class="control-group"
            >
              <el-input
                v-model="searchValue"
                class="search-input"
                placeholder="输入省份"
                clearable
              />
            </el-form-item>

            <!-- 城市查询 -->
            <el-form-item 
              v-if="searchType === 'city'" 
              label="城市" 
              class="control-group"
            >
              <el-input
                v-model="searchValue"
                class="search-input"
                placeholder="输入城市"
                clearable
              />
            </el-form-item>

            <!-- 区县查询 -->
            <el-form-item 
              v-if="searchType === 'district'" 
              label="区县" 
              class="control-group"
            >
              <el-input
                v-model="searchValue"
                class="search-input"
                placeholder="输入区县"
                clearable
              />
            </el-form-item>

            <div class="action-buttons">
              <el-button 
                type="primary" 
                class="search-btn"
                @click="handleSearch(true)"
                :loading="searching"
              >
                <el-icon><Search /></el-icon>
                <span>搜索</span>
              </el-button>
              <el-button 
                class="reset-btn"
                @click="handleReset"
                :disabled="searching"
              >
                重置
              </el-button>
            </div>
          </div>
        </el-form>

        <!-- 搜索结果 -->
        <div class="search-results">
          <div v-loading="searching" class="loading-container" element-loading-text="加载中...">
            <template v-if="!searching">
              <div v-if="images.length === 0 && hasSearched" class="no-results">
                <el-empty description="未找到相关图片" image-size="100" />
              </div>
  
              <div v-else-if="images.length > 0" class="results-grid">
                <div class="results-header">
                  <h3>图片列表 (共 {{ pagination.total }} 张图片)</h3>
                </div>
                <div class="images-grid">
                  <imageCard
                    v-for="image in images"
                    :key="image.image_id"
                    :image="image"
                    @delete-success="handleSuccess"
                    @edit-image="handleEditImage"
                  />
                </div>
  
                <div class="pagination-footer">
                  <el-pagination
                    v-model:current-page="pagination.currentPage"
                    v-model:page-size="pagination.pageSize"
                    :total="pagination.total"
                    layout="prev, pager, next, jumper"
                    @size-change="handleSizeChange"
                    @current-change="handleCurrentChange"
                    :disabled="searching"
                    background
                  />
                </div>
              </div>
  
              <div v-else class="welcome-message">
                <el-icon class="welcome-icon"><Search /></el-icon>
                <h3>欢迎来到图片库</h3>
                <p>点击搜索按钮查看您的图片</p>
              </div>
            </template>
          </div>
        </div>
      </div>
    </div>
  
    <!-- 图片裁剪对话框 -->
    <ImageEditor
      v-model:visible="cropperVisible"
      :image-url="currentImageUrl"
      :container-height="500"
      @cropped="handleCropped"
    />
    
    <!-- 裁剪预览对话框 -->
    <el-dialog
      v-model="previewVisible"
      title="预览裁剪后的图片"
      width="700px"
      :close-on-click-modal="false"
    >
      <div class="preview-container">
        <div class="preview-hint">
          <el-text type="info">请确认裁剪效果，点击"确认提交"将更新图片</el-text>
        </div>
        
        <div class="preview-image-wrapper">
          <el-image 
            :src="croppedImageBase64" 
            fit="contain"
            class="preview-image"
          >
            <template #error>
              <div class="image-error">
                <el-icon><Picture /></el-icon>
                <span>图片加载失败</span>
              </div>
            </template>
          </el-image>
        </div>
      </div>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="handleCancelPreview" :disabled="isUploading">
            取消
          </el-button>
          <el-button 
            type="primary" 
            @click="handleConfirmCrop"
            :loading="isUploading"
          >
            确认提交
          </el-button>
        </span>
      </template>
    </el-dialog>
  </template>
  
<script setup lang="ts">
import ImageCard from '@/components/ImageCard.vue';
import AppHeader from '@/components/AppHeader.vue';
import ImageEditor from '@/components/ImageEidtor.vue';
import { ref, reactive, onMounted, watch } from 'vue'; 
import { ElMessage } from 'element-plus';
import { Search, Picture } from '@element-plus/icons-vue';
import { ImageAPI } from '@/api';

// 图片接口定义，完整的定义，适配各种接口
interface ImageItem {
  image_id: number;
  image_name: string;
  image_path: string;
  tags: string[];
}

// 分页状态接口
interface PaginationState {
  currentPage: number;
  pageSize: number;
  total: number;
}

const searchType = ref('all');
const searchValue = ref('');
const timeRange = ref<[string, string] | null>(null);
const searching = ref(false);
const hasSearched = ref(false);
const images = ref<ImageItem[]>([]);

// 图片裁剪相关状态
const cropperVisible = ref(false);
const currentImageUrl = ref('');
const currentImageId = ref(0);

// 裁剪预览相关状态
const previewVisible = ref(false);
const croppedImageBase64 = ref('');
const isUploading = ref(false);

// 监听裁剪对话框关闭，释放 blob URL
watch(cropperVisible, (newVal, oldVal) => {
  if (oldVal && !newVal && currentImageUrl.value) {
    // 对话框从打开变为关闭时，释放 blob URL
    URL.revokeObjectURL(currentImageUrl.value);
    currentImageUrl.value = '';
  }
});

// 监听预览对话框关闭，清理所有缓存
watch(previewVisible, (newVal, oldVal) => {
  if (oldVal && !newVal) {
    // 对话框从打开变为关闭时，清理所有缓存资源
    clearImageCache();
  }
});

// 分页状态对象
const pagination = reactive<PaginationState>({
  currentPage: 1,
  pageSize: 9,
  total: 0,
});

/**
 * 处理查询类型变化
 */
const handleTypeChange = () => {
  searchValue.value = '';
  timeRange.value = null;
};

/**
 * 核心搜索函数
 * @param resetPage 是否重置页码
 */
const handleSearch = async (resetPage: boolean = true) => {
  // 验证输入
  if (searchType.value === 'tag' && !searchValue.value.trim()) {
    ElMessage.warning('请输入标签');
    return;
  }
  if (searchType.value === 'time' && (!timeRange.value || timeRange.value.length !== 2)) {
    ElMessage.warning('请选择时间范围');
    return;
  }
  if (['province', 'city', 'district'].includes(searchType.value) && !searchValue.value.trim()) {
    ElMessage.warning(`请输入${searchType.value === 'province' ? '省份' : searchType.value === 'city' ? '城市' : '区县'}`);
    return;
  }
  
  // 重置页码
  if (resetPage) {
    pagination.currentPage = 1;
  }

  searching.value = true;
  hasSearched.value = true;

  try {
    let response;
    
    // 根据查询类型调用不同的API
    switch (searchType.value) {
      case 'all':
        response = await ImageAPI.getPersonalAllImages(pagination.currentPage, pagination.pageSize);
        break;
      case 'tag':
        response = await ImageAPI.getImagesByTag(searchValue.value.trim(), pagination.currentPage, pagination.pageSize);
        break;
      case 'time':
        if (timeRange.value && timeRange.value.length === 2) {
          response = await ImageAPI.getImagesByTimeRange(
            timeRange.value[0], 
            timeRange.value[1], 
            pagination.currentPage, 
            pagination.pageSize
          );
        }
        break;
      case 'province':
        response = await ImageAPI.getImagesByProvince(searchValue.value.trim(), pagination.currentPage, pagination.pageSize);
        break;
      case 'city':
        response = await ImageAPI.getImagesByCity(searchValue.value.trim(), pagination.currentPage, pagination.pageSize);
        break;
      case 'district':
        response = await ImageAPI.getImagesByDistrict(searchValue.value.trim(), pagination.currentPage, pagination.pageSize);
        break;
      default:
        ElMessage.error('未知的查询类型');
        return;
    }

    if (response && response.data) {
      const result = response.data;
      
      if (result.code === 200 && result.data) {
        images.value = result.data.images || [];
        pagination.total = result.data.pageTotal || 0;

        if (images.value.length === 0) {
          ElMessage.info('未找到相关图片');
        } else {
          ElMessage.success(`找到 ${pagination.total} 张图片`);
        }
      } else {
        ElMessage.error(result.message || '查询失败');
        images.value = [];
        pagination.total = 0;
      }
    }
  } catch (error: any) {
    console.error('查询失败:', error);
    ElMessage.error(error.message || '查询失败，请重试');
    images.value = [];
    pagination.total = 0;
  } finally {
    searching.value = false;
  }
};

/**
 * 处理每页大小变化
 */
const handleSizeChange = (val: number) => {
  handleSearch(true);
};

/**
 * 处理当前页码变化
 */
const handleCurrentChange = (val: number) => {
  handleSearch(false);
};

/**
 * 重置搜索
 */
const handleReset = () => {
  searchType.value = 'all';
  searchValue.value = '';
  timeRange.value = null;
  images.value = [];
  hasSearched.value = false;
  
  pagination.currentPage = 1;
  pagination.pageSize = 16;
  pagination.total = 0;
  
  ElMessage.info('已重置搜索条件');
};

/**
 * 处理操作成功事件
 */
const handleSuccess = () => {
  // 重新加载当前页数据
  handleSearch(false);
};

/**
 * 处理编辑图片事件
 */
const handleEditImage = async (imageId: number, imagePath: string) => {
  try {
    currentImageId.value = imageId;
    
    // 获取图片文件（blob格式）
    const response = await ImageAPI.getFile(imageId);
    if (response.status !== 200) {
      ElMessage.error('获取图片失败');
      return;
    }
    
    // 将blob转换为可用的URL
    const blob = response.data;
    currentImageUrl.value = URL.createObjectURL(blob);
    
    // 打开裁剪对话框
    cropperVisible.value = true;
  } catch (error: any) {
    console.error('获取图片失败:', error);
    ElMessage.error(error.message || '获取图片失败');
  }
};

/**
 * 处理裁剪完成事件
 */
const handleCropped = (base64: string) => {
  // 保存裁剪后的图片数据
  croppedImageBase64.value = base64;
  
  // 打开预览对话框
  previewVisible.value = true;
  
  ElMessage.success('图片裁剪成功，请预览确认');
};

/**
 * 清理所有缓存的图片资源
 */
const clearImageCache = () => {
  // 释放 blob URL
  if (currentImageUrl.value) {
    URL.revokeObjectURL(currentImageUrl.value);
    currentImageUrl.value = '';
  }
  
  // 清空裁剪后的图片数据
  croppedImageBase64.value = '';
  
  // 重置图片 ID
  currentImageId.value = 0;
};

/**
 * 将 base64 转换为 File 对象
 */
const base64ToFile = (base64: string, filename: string): File => {
  const arr = base64.split(',');
  const match = arr[0]?.match(/:(.*?);/);
  const mime = match?.[1] || 'image/png';
  const bstr = atob(arr[1] || '');
  let n = bstr.length;
  const u8arr = new Uint8Array(n);
  while (n--) {
    u8arr[n] = bstr.charCodeAt(n);
  }
  return new File([u8arr], filename, { type: mime });
};

/**
 * 确认提交裁剪后的图片
 */
const handleConfirmCrop = async () => {
  if (!croppedImageBase64.value) {
    ElMessage.error('没有可提交的图片');
    return;
  }
  
  try {
    isUploading.value = true;
    
    // 将 base64 转换为 File 对象
    const file = base64ToFile(croppedImageBase64.value, `image_${currentImageId.value}.png`);
    
    // 调用更新 API
    const response = await ImageAPI.update(currentImageId.value, file);
    
    if (response.data.code === 200) {
      ElMessage.success('图片更新成功');
      
      // 关闭预览对话框（watch 会自动清理缓存）
      previewVisible.value = false;
      
      // 刷新图片列表
      handleSearch(false);
    } else {
      ElMessage.error(response.data.message || '图片更新失败');
    }
  } catch (error: any) {
    console.error('图片更新失败:', error);
    ElMessage.error(error.message || '图片更新失败');
  } finally {
    isUploading.value = false;
  }
};

/**
 * 取消预览
 */
const handleCancelPreview = () => {
  // 关闭预览对话框（watch 会自动清理缓存）
  previewVisible.value = false;
  
  ElMessage.info('已取消提交');
};

// 组件挂载时自动加载个人所有图片
onMounted(() => {
  handleSearch(true);
});
</script>

<style scoped>
/* ---------------------------------- */
/* 整体和容器美化 (CSS) */
/* ---------------------------------- */

.page-container {
    min-height: 100vh;
    padding: 30px;
    background-color: #f5f7fa; /* 浅灰色背景 */
    font-family: "Helvetica Neue", Helvetica, Arial, "PingFang SC", "Hiragino Sans GB", "Microsoft YaHei", sans-serif;
}

.image-search-container {
    max-width: 1400px;
    margin: 0 auto;
    background: #ffffff;
    border-radius: 12px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05); /* 提升感 */
    padding: 30px;
}

/* 头部标题样式 */
.search-header h1 {
    font-size: 28px;
    font-weight: bold;
    color: #303133;
    margin-bottom: 5px;
}

.search-header p {
    font-size: 16px;
    color: #606266;
}

/* 搜索表单区域美化 */
.search-form {
    padding: 20px;
    margin-bottom: 25px;
    border: 1px solid #ebeef5;
    border-radius: 8px;
    background-color: #fcfcfc;
}

/* 使用 Flexbox 布局搜索控件，使其在同一行排列 */
.search-controls {
    display: flex;
    flex-wrap: wrap; /* 允许控件换行 */
    gap: 20px; /* 控件间间距 */
    align-items: center;
}

/* 重置 Element Plus 默认的表单项间距 */
.search-controls .el-form-item {
    margin-right: 0 !important;
    margin-bottom: 0 !important;
}

/* 动作按钮推到最右边 */
.search-controls .action-buttons {
    margin-left: auto;
    display: flex;
    gap: 10px;
}

/* ---------------------------------- */
/* 搜索结果和图片网格布局 (CSS Grid) */
/* ---------------------------------- */

.search-results {
    margin-top: 30px;
}

.results-header {
    display: flex;
    align-items: center;
    margin-bottom: 15px;
    padding-bottom: 10px;
    border-bottom: 1px solid #e4e7ed; /* 分割线 */
}

/* 核心：使用 CSS Grid 实现多列图片排列 */
.images-grid {
    display: grid;
    /* 自动填充列，每列最小宽度为 300px，最大占满可用空间 (1fr) */
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr)); 
    gap: 20px; /* 图片卡片之间的间距 */
    padding-bottom: 20px;
}

/* 确保 ImageCard.vue 中的 el-card 样式美观 */
/* 注意: 如果 ImageCard 是一个单独的组件，你可能需要在其内部定义样式 */
.images-grid :deep(.image-card) {
    border-radius: 8px;
    overflow: hidden;
    transition: box-shadow 0.3s;
}

.images-grid :deep(.image-card:hover) {
    box-shadow: 0 6px 16px rgba(0, 0, 0, 0.1);
}

/* 分页居中 */
.pagination-footer {
    margin-top: 30px;
    display: flex;
    justify-content: center;
}

/* 欢迎和无结果消息的样式 */
.welcome-message, .no-results {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    min-height: 300px;
    color: #909399;
}

.welcome-message .welcome-icon {
    font-size: 48px;
    margin-bottom: 10px;
}

/* ---------------------------------- */
/* 响应式调整 (CSS Media Queries) */
/* ---------------------------------- */

@media (max-width: 992px) {
    /* 在中等屏幕上，时间选择器的宽度调整 */
    .search-controls .control-group .time-picker {
        width: 100%;
        max-width: 350px;
    }
}

@media (max-width: 768px) {
    .page-container {
        padding: 15px;
    }

    .image-search-container {
        padding: 20px;
    }

    /* 在小屏幕上，搜索控件垂直排列 */
    .search-controls {
        flex-direction: column;
        align-items: stretch;
    }
    
    .search-controls .control-group {
        width: 100%; /* 表单项占满宽度 */
    }

    .search-controls .action-buttons {
        margin-left: 0;
        justify-content: flex-end; /* 按钮靠右对齐 */
    }
    
    /* 调整图片网格的最小列宽，避免在小屏上太拥挤 */
    .images-grid {
        grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
    }
}

/* ---------------------------------- */
/* 预览对话框样式 */
/* ---------------------------------- */

.preview-container {
    padding: 10px 0;
}

.preview-hint {
    margin-bottom: 20px;
    padding: 12px;
    background-color: #f4f4f5;
    border-radius: 4px;
    text-align: center;
}

.preview-image-wrapper {
    display: flex;
    justify-content: center;
    align-items: center;
    min-height: 400px;
    max-height: 500px;
    background-color: #f5f7fa;
    border-radius: 8px;
    padding: 20px;
    overflow: hidden;
}

.preview-image {
    max-width: 100%;
    max-height: 460px;
    border-radius: 4px;
}

.image-error {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 10px;
    color: #909399;
}

.image-error .el-icon {
    font-size: 48px;
}

.dialog-footer {
    display: flex;
    justify-content: flex-end;
    gap: 10px;
}
</style>