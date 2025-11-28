<template>
  <AppHeader />
  <div class="page-container">
    <div class="image-search-container">
      <!-- 搜索表单 -->
      <div class="search-header">
        <h1>图片库</h1>
        <p>查看和搜索您的图片</p>
      </div>

      <div class="search-card">
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
                <el-option v-if="isAdmin" label="所有图片" value="all-images" />
                <el-option v-if="isAdmin" label="用户图片" value="user-images" />
                <el-option label="标签查询" value="tag" />
                <el-option label="时间查询" value="time" />
                <el-option label="地区查询" value="area" />
              </el-select>
            </el-form-item>

            <!-- 用户选择器（仅管理员，查询类型为"用户图片"时显示） -->
            <el-form-item 
              v-if="isAdmin && searchType === 'user-images'" 
              label="选择用户" 
              class="control-group"
            >
              <el-select
                v-model="selectedUserID"
                class="search-select"
                placeholder="选择用户"
                filterable
                remote
                :remote-method="searchUsers"
                :loading="loadingUsers"
                clearable
              >
                <el-option
                  v-for="user in userOptions"
                  :key="user.id"
                  :label="`${user.username} (ID: ${user.id})`"
                  :value="user.id"
                />
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

            <!-- 区域查询：省份、城市、区县 -->
            <el-form-item 
              v-if="searchType === 'area'" 
              label="省份" 
              class="control-group"
            >
              <el-input
                v-model="areaProvince"
                class="search-input"
                placeholder="输入省份（可选）"
                clearable
              />
            </el-form-item>

            <el-form-item 
              v-if="searchType === 'area'" 
              label="城市" 
              class="control-group"
            >
              <el-input
                v-model="areaCity"
                class="search-input"
                placeholder="输入城市（可选）"
                clearable
              />
            </el-form-item>

            <el-form-item 
              v-if="searchType === 'area'" 
              label="区县" 
              class="control-group"
            >
              <el-input
                v-model="areaDistrict"
                class="search-input"
                placeholder="输入区县（可选）"
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
      </div>

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
                  @click="handleImageClick(image.image_id)"
                />
              </div>

              <div class="pagination-footer">
                <el-pagination
                  v-model:current-page="pagination.currentPage"
                  v-model:page-size="pagination.pageSize"
                  :total="pagination.total"
                  :page-sizes="pageSizeOptions.map(option => option.value)"
                  layout="total, sizes, prev, pager, next, jumper"
                  @size-change="handleSizeChange"
                  @current-change="handleCurrentChange"
                  :disabled="searching"
                  background
                />
              </div>
            </div>

            <div v-else class="welcome-message">
              <div class="welcome-content">
                <el-icon class="welcome-icon"><Search /></el-icon>
                <h3>欢迎来到图片库</h3>
                <p>选择查询类型并点击搜索按钮查看您的图片</p>
              </div>
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
    class="preview-dialog"
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
          class="confirm-btn"
        >
          确认提交
        </el-button>
      </span>
    </template>
  </el-dialog>

  <!-- 走马灯轮播对话框 -->
  <el-dialog
    v-model="carouselVisible"
    title="图片轮播"
    width="80%"
    :close-on-click-modal="true"
    class="carousel-dialog"
  >
    <div class="carousel-container">
      <el-carousel
        :initial-index="carouselInitialIndex"
        height="600px"
        indicator-position="outside"
        arrow="always"
      >
        <el-carousel-item v-for="image in images" :key="image.image_id">
          <div class="carousel-item-content">
            <el-image 
              :src="getOriginImageUrl(image.image_path)" 
              fit="contain"
              class="carousel-image"
            >
              <template #error>
                <div class="image-error">
                  <el-icon><Picture /></el-icon>
                  <span>图片加载失败</span>
                </div>
              </template>
            </el-image>
            <div class="carousel-info">
              <h3>{{ image.image_name }}</h3>
              <div class="carousel-tags">
                <el-tag
                  v-for="tag in image.tags"
                  :key="tag"
                  type="info"
                  effect="plain"
                  size="small"
                  round
                >
                  {{ tag }}
                </el-tag>
                <el-text v-if="image.tags.length === 0" size="small" type="info">
                  暂无标签
                </el-text>
              </div>
            </div>
          </div>
        </el-carousel-item>
      </el-carousel>
    </div>
  </el-dialog>
</template>

<script setup lang="ts">
import ImageCard from '@/components/ImageCard.vue';
import AppHeader from '@/components/AppHeader.vue';
import ImageEditor from '@/components/ImageEidtor.vue';
import { ref, reactive, onMounted, watch, computed } from 'vue'; 
import { ElMessage } from 'element-plus';
import { Search, Picture } from '@element-plus/icons-vue';
import { ImageAPI, UserAPI } from '@/api';
import { getOriginImageUrl } from '@/utils/url';
import type { ImageItem, PaginationState, User } from '@/type';
import { base64ToFile } from '@/utils/file';

const searchType = ref('all');
const searchValue = ref('');
const timeRange = ref<[string, string] | null>(null);
const searching = ref(false);
const hasSearched = ref(false);
const images = ref<ImageItem[]>([]);

// 区域查询相关状态
const areaProvince = ref('');
const areaCity = ref('');
const areaDistrict = ref('');

// 管理员相关状态
const isAdmin = computed(() => {
  const userInfo = JSON.parse(localStorage.getItem('userInfo') || '{}');
  return userInfo.is_admin === true;
});

// 用户选择相关状态
const selectedUserID = ref<number | null>(null);
const userOptions = ref<User[]>([]);
const loadingUsers = ref(false);

// 图片裁剪相关状态
const cropperVisible = ref(false);
const currentImageUrl = ref('');
const currentImageId = ref(0);

// 裁剪预览相关状态
const previewVisible = ref(false);
const croppedImageBase64 = ref('');
const isUploading = ref(false);

// 走马灯相关状态
const carouselVisible = ref(false);
const carouselInitialIndex = ref(0);

// 监听裁剪对话框关闭，释放 blob URL
watch(cropperVisible, (newVal, oldVal) => {
  if (oldVal && !newVal && currentImageUrl.value) {
    URL.revokeObjectURL(currentImageUrl.value);
    currentImageUrl.value = '';
  }
});

// 监听预览对话框关闭，清理所有缓存
watch(previewVisible, (newVal, oldVal) => {
  if (oldVal && !newVal) {
    clearImageCache();
  }
});

// 分页状态对象
const pagination = reactive<PaginationState>({
  currentPage: 1,
  pageSize: 12,
  total: 0,
});

// 页面大小选项
const pageSizeOptions = [
  { label: '12 张/页', value: 12 },
  { label: '24 张/页', value: 24 },
];

/**
 * 处理查询类型变化
 */
const handleTypeChange = () => {
  searchValue.value = '';
  timeRange.value = null;
  areaProvince.value = '';
  areaCity.value = '';
  areaDistrict.value = '';
  selectedUserID.value = null;
  
  // 如果切换到用户图片查询，加载用户列表
  if (isAdmin.value && searchType.value === 'user-images') {
    loadUsers();
  }
};

/**
 * 加载用户列表
 */
const loadUsers = async () => {
  if (loadingUsers.value) return;
  
  try {
    loadingUsers.value = true;
    const response = await UserAPI.getAllUsers(1, 100); // 获取前100个用户
    if (response.data && response.data.code === 200 && response.data.data) {
      userOptions.value = response.data.data.users || [];
    }
  } catch (error: any) {
    console.error('加载用户列表失败:', error);
    ElMessage.error('加载用户列表失败');
  } finally {
    loadingUsers.value = false;
  }
};

/**
 * 搜索用户（用于远程搜索）
 */
const searchUsers = async (query: string) => {
  if (!query) {
    loadUsers();
    return;
  }
  
  if (loadingUsers.value) return;
  
  try {
    loadingUsers.value = true;
    const response = await UserAPI.getAllUsers(1, 100);
    if (response.data && response.data.code === 200 && response.data.data) {
      const allUsers = response.data.data.users || [];
      // 根据查询字符串过滤用户
      userOptions.value = allUsers.filter((user: User) => 
        user.username.toLowerCase().includes(query.toLowerCase()) ||
        user.email.toLowerCase().includes(query.toLowerCase()) ||
        user.id.toString().includes(query)
      );
    }
  } catch (error: any) {
    console.error('搜索用户失败:', error);
  } finally {
    loadingUsers.value = false;
  }
};

/**
 * 核心搜索函数
 * @param resetPage 是否重置页码
 */
const handleSearch = async (resetPage: boolean = true) => {
  if (searchType.value === 'tag' && !searchValue.value.trim()) {
    ElMessage.warning('请输入标签');
    return;
  }
  if (searchType.value === 'time' && (!timeRange.value || timeRange.value.length !== 2)) {
    ElMessage.warning('请选择时间范围');
    return;
  }
  if (searchType.value === 'area') {
    if (!areaProvince.value.trim() && !areaCity.value.trim() && !areaDistrict.value.trim()) {
      ElMessage.warning('请至少输入一个区域条件（省份、城市或区县）');
      return;
    }
  }
  if (searchType.value === 'user-images' && !selectedUserID.value) {
    ElMessage.warning('请选择用户');
    return;
  }
  
  if (resetPage) {
    pagination.currentPage = 1;
  }

  searching.value = true;
  hasSearched.value = true;

  try {
    let response;
    
    switch (searchType.value) {
      case 'all':
        response = await ImageAPI.getPersonalAllImages(pagination.currentPage, pagination.pageSize);
        break;
      case 'all-images':
        // 管理员获取所有图片
        response = await ImageAPI.getAllImages(pagination.currentPage, pagination.pageSize);
        break;
      case 'user-images':
        // 管理员获取指定用户的图片
        if (selectedUserID.value) {
          response = await ImageAPI.getImagesByUserID(selectedUserID.value, pagination.currentPage, pagination.pageSize);
        }
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
      case 'area':
        response = await ImageAPI.getImagesByArea(
          areaProvince.value.trim() || undefined,
          areaCity.value.trim() || undefined,
          areaDistrict.value.trim() || undefined,
          pagination.currentPage,
          pagination.pageSize
        );
        break;
      default:
        ElMessage.error('未知的查询类型');
        return;
    }

    if (response && response.data) {
      const result = response.data;
      
      if (result.code === 200 && result.data) {
        images.value = result.data.images || [];
        pagination.total = result.data.totalCount || 0;

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
  areaProvince.value = '';
  areaCity.value = '';
  areaDistrict.value = '';
  selectedUserID.value = null;
  images.value = [];
  hasSearched.value = false;
  
  pagination.currentPage = 1;
  pagination.pageSize = 12;
  pagination.total = 0;
  
  ElMessage.info('已重置搜索条件');
};

/**
 * 处理操作成功事件
 */
const handleSuccess = () => {
  handleSearch(false);
};

/**
 * 处理编辑图片事件
 */
const handleEditImage = async (imageId: number, imagePath: string) => {
  try {
    currentImageId.value = imageId;
    
    const response = await ImageAPI.getFile(imageId);
    if (response.status !== 200) {
      ElMessage.error('获取图片失败');
      return;
    }
    
    const blob = response.data;
    currentImageUrl.value = URL.createObjectURL(blob);
    
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
  croppedImageBase64.value = base64;
  previewVisible.value = true;
  ElMessage.success('图片裁剪成功，请预览确认');
};

/**
 * 清理所有缓存的图片资源
 */
const clearImageCache = () => {
  if (currentImageUrl.value) {
    URL.revokeObjectURL(currentImageUrl.value);
    currentImageUrl.value = '';
  }
  
  croppedImageBase64.value = '';
  currentImageId.value = 0;
};

/**
 * 将 base64 转换为 File 对象
 */
const convertBase64ToFile = (base64: string, filename: string): File => {
  return base64ToFile(base64, filename);
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
    
    const file = convertBase64ToFile(croppedImageBase64.value, `image_${currentImageId.value}.png`);
    
    const response = await ImageAPI.update(currentImageId.value, file);
    
    if (response.data.code === 200) {
      ElMessage.success('图片更新成功');
      previewVisible.value = false;
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
  previewVisible.value = false;
  ElMessage.info('已取消提交');
};

/**
 * 处理图片卡片点击事件，打开走马灯
 */
const handleImageClick = (imageId: number) => {
  const index = images.value.findIndex(img => img.image_id === imageId);
  if (index !== -1) {
    carouselInitialIndex.value = index;
    carouselVisible.value = true;
  }
};

// 组件挂载时，如果是管理员且查询类型是用户图片，加载用户列表
onMounted(() => {
  // 如果是管理员且查询类型是用户图片，加载用户列表
  if (isAdmin.value && searchType.value === 'user-images') {
    loadUsers();
  }
  // 页面初始状态显示欢迎消息，不自动搜索
});
</script>

<style scoped>
.page-container {
  min-height: 100vh;
  padding: 30px;
  background: linear-gradient(135deg, #f0f8ff 0%, #e6f3ff 50%, #dcedff 100%);
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
}

.image-search-container {
  max-width: 1400px;
  margin: 0 auto;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 20px;
  box-shadow: 
    0 8px 32px rgba(64, 158, 255, 0.15),
    0 2px 8px rgba(64, 158, 255, 0.1);
  padding: 40px;
  border: 1px solid rgba(64, 158, 255, 0.1);
}

/* 头部标题样式 */
.search-header {
  text-align: center;
  margin-bottom: 30px;
}

.search-header h1 {
  font-size: 32px;
  font-weight: 700;
  color: #2c3e50;
  margin-bottom: 8px;
  background: linear-gradient(135deg, #409EFF, #66B3FF);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.search-header p {
  font-size: 16px;
  color: #7f8c8d;
  margin: 0;
}

/* 搜索卡片样式 */
.search-card {
  background: linear-gradient(135deg, #f8fbff, #f0f7ff);
  border-radius: 16px;
  padding: 30px;
  margin-bottom: 30px;
  border: 1px solid rgba(64, 158, 255, 0.1);
  box-shadow: 0 4px 15px rgba(64, 158, 255, 0.1);
}

.search-form {
  width: 100%;
}

/* 搜索控件布局 */
.search-controls {
  display: flex;
  flex-wrap: wrap;
  gap: 20px;
  align-items: flex-end;
}

.control-group {
  margin-bottom: 0 !important;
}

.control-group :deep(.el-form-item__label) {
  color: #409EFF;
  font-weight: 600;
  margin-bottom: 8px;
  display: block;
}

.search-select,
.search-input,
.time-picker {
  min-width: 200px;
}

:deep(.search-select .el-input__wrapper),
:deep(.search-input .el-input__wrapper),
:deep(.time-picker .el-input__wrapper) {
  background: #ffffff;
  border: 2px solid #e6f3ff;
  border-radius: 12px;
  transition: all 0.3s ease;
  box-shadow: none;
}

:deep(.search-select .el-input__wrapper:hover),
:deep(.search-input .el-input__wrapper:hover),
:deep(.time-picker .el-input__wrapper:hover),
:deep(.search-select .el-input__wrapper.is-focus),
:deep(.search-input .el-input__wrapper.is-focus),
:deep(.time-picker .el-input__wrapper.is-focus) {
  border-color: #409EFF;
  box-shadow: 0 0 0 3px rgba(64, 158, 255, 0.1);
}

/* 动作按钮 */
.action-buttons {
  margin-left: auto;
  display: flex;
  gap: 12px;
}

.search-btn {
  background: linear-gradient(135deg, #409EFF, #66B3FF);
  border: none;
  border-radius: 12px;
  padding: 12px 24px;
  font-weight: 600;
  transition: all 0.3s ease;
  box-shadow: 0 4px 12px rgba(64, 158, 255, 0.3);
}

.search-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(64, 158, 255, 0.4);
  background: linear-gradient(135deg, #66B3FF, #409EFF);
}

.reset-btn {
  border: 2px solid #e6f3ff;
  border-radius: 12px;
  padding: 12px 24px;
  color: #409EFF;
  font-weight: 500;
  transition: all 0.3s ease;
}

.reset-btn:hover {
  border-color: #409EFF;
  background: rgba(64, 158, 255, 0.1);
  transform: translateY(-1px);
}

/* 搜索结果区域 */
.search-results {
  margin-top: 20px;
}

.results-header {
  margin-bottom: 25px;
  padding-bottom: 15px;
  border-bottom: 2px solid #e6f3ff;
}

.results-header h3 {
  color: #2c3e50;
  font-size: 20px;
  font-weight: 600;
  margin: 0;
}

/* 图片网格布局 */
.images-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 24px;
  padding-bottom: 30px;
}

/* 分页样式 */
.pagination-footer {
  margin-top: 40px;
  display: flex;
  justify-content: center;
}

:deep(.el-pagination.is-background .btn-prev),
:deep(.el-pagination.is-background .btn-next),
:deep(.el-pagination.is-background .el-pager li) {
  background: #f8fbff;
  border: 1px solid #e6f3ff;
  border-radius: 8px;
  margin: 0 4px;
  transition: all 0.3s ease;
}

:deep(.el-pagination.is-background .btn-prev:hover),
:deep(.el-pagination.is-background .btn-next:hover),
:deep(.el-pagination.is-background .el-pager li:hover) {
  border-color: #409EFF;
  background: rgba(64, 158, 255, 0.1);
}

:deep(.el-pagination.is-background .el-pager li.is-active) {
  background: linear-gradient(135deg, #409EFF, #66B3FF);
  border-color: #409EFF;
  color: white;
}

/* 欢迎和无结果消息 */
.welcome-message,
.no-results {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 400px;
  text-align: center;
}

.welcome-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
}

.welcome-icon {
  font-size: 64px;
  color: #409EFF;
  margin-bottom: 8px;
}

.welcome-message h3 {
  color: #2c3e50;
  font-size: 24px;
  font-weight: 600;
  margin: 0;
}

.welcome-message p {
  color: #7f8c8d;
  font-size: 16px;
  margin: 0;
}

/* 预览对话框样式 */
.preview-container {
  padding: 10px 0;
}

.preview-hint {
  margin-bottom: 20px;
  padding: 16px;
  background: linear-gradient(135deg, #f8fbff, #f0f7ff);
  border-radius: 12px;
  text-align: center;
  border: 1px solid rgba(64, 158, 255, 0.1);
}

.preview-image-wrapper {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 400px;
  max-height: 500px;
  background: #f8fbff;
  border-radius: 12px;
  padding: 20px;
  overflow: hidden;
  border: 2px dashed #e6f3ff;
}

.preview-image {
  max-width: 100%;
  max-height: 460px;
  border-radius: 8px;
}

.image-error {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  color: #909399;
}

.image-error .el-icon {
  font-size: 48px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

.confirm-btn {
  background: linear-gradient(135deg, #409EFF, #66B3FF);
  border: none;
  border-radius: 12px;
  padding: 12px 24px;
  font-weight: 600;
}

.confirm-btn:hover {
  background: linear-gradient(135deg, #66B3FF, #409EFF);
}

/* 走马灯对话框样式 */
.carousel-container {
  width: 100%;
}

.carousel-item-content {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: #f8fbff;
  border-radius: 12px;
  overflow: hidden;
}

.carousel-image {
  flex: 1;
  width: 100%;
  height: 500px;
}

.carousel-info {
  padding: 20px;
  background: #ffffff;
  border-top: 1px solid #e6f3ff;
}

.carousel-info h3 {
  margin: 0 0 12px 0;
  font-size: 18px;
  font-weight: 600;
  color: #2c3e50;
}

.carousel-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  align-items: center;
}

/* 走马灯箭头样式优化 */
.carousel-container :deep(.el-carousel__arrow) {
  background-color: rgba(64, 158, 255, 0.8);
  border-radius: 50%;
}

.carousel-container :deep(.el-carousel__arrow:hover) {
  background-color: rgba(64, 158, 255, 0.95);
}

/* 对话框样式优化 */
:deep(.preview-dialog .el-dialog),
:deep(.carousel-dialog .el-dialog) {
  border-radius: 16px;
  overflow: hidden;
}

:deep(.preview-dialog .el-dialog__header),
:deep(.carousel-dialog .el-dialog__header) {
  background: linear-gradient(135deg, #f8fbff, #f0f7ff);
  border-bottom: 1px solid rgba(64, 158, 255, 0.1);
  margin: 0;
  padding: 20px;
}

:deep(.preview-dialog .el-dialog__title),
:deep(.carousel-dialog .el-dialog__title) {
  color: #2c3e50;
  font-weight: 600;
  font-size: 18px;
}

/* 响应式设计 */
@media (max-width: 1200px) {
  .image-search-container {
    max-width: 95%;
    padding: 30px;
  }
}

@media (max-width: 768px) {
  .page-container {
    padding: 20px 15px;
  }

  .image-search-container {
    padding: 25px 20px;
  }

  .search-card {
    padding: 20px;
  }

  .search-controls {
    flex-direction: column;
    align-items: stretch;
    gap: 16px;
  }

  .control-group {
    width: 100%;
  }

  .search-select,
  .search-input,
  .time-picker {
    min-width: 100%;
    width: 100%;
  }

  .action-buttons {
    margin-left: 0;
    justify-content: stretch;
  }

  .action-buttons .el-button {
    flex: 1;
  }

  .images-grid {
    grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
    gap: 16px;
  }
}

@media (max-width: 480px) {
  .image-search-container {
    padding: 20px 15px;
  }

  .search-header h1 {
    font-size: 28px;
  }

  .images-grid {
    grid-template-columns: 1fr;
    gap: 16px;
  }

  .welcome-icon {
    font-size: 56px;
  }

  .welcome-message h3 {
    font-size: 20px;
  }
}
</style>