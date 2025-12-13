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
            <template v-if="searchType === 'time'">
              <el-form-item 
                label="开始日期" 
                class="control-group"
              >
                <el-date-picker
                  v-model="startTime"
                  type="date"
                  :placeholder="isMobile ? '开始日期' : '选择开始日期'"
                  value-format="YYYY-MM-DD"
                  class="time-picker"
                  :teleported="false"
                  style="width: 100%"
                  size="default"
                  clearable
                />
              </el-form-item>
              <el-form-item 
                label="结束日期" 
                class="control-group"
              >
                <el-date-picker
                  v-model="endTime"
                  type="date"
                  :placeholder="isMobile ? '结束日期' : '选择结束日期'"
                  value-format="YYYY-MM-DD"
                  class="time-picker"
                  :teleported="false"
                  style="width: 100%"
                  size="default"
                  clearable
                />
              </el-form-item>
            </template>

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
                  :user-id="currentSearchUserID !== null ? currentSearchUserID : undefined"
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
                  :layout="isMobile ? 'prev, pager, next' : 'total, sizes, prev, pager, next, jumper'"
                  @size-change="handleSizeChange"
                  @current-change="handleCurrentChange"
                  :disabled="searching"
                  background
                  small
                  class="pagination-component"
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
    :width="isMobile ? '100%' : '700px'"
    :close-on-click-modal="false"
    :fullscreen="isMobile"
    :show-close="!isMobile"
    class="preview-dialog"
    :class="{ 'mobile-dialog': isMobile }"
  >
    <div class="preview-container" :class="{ 'mobile-layout': isMobile }">
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
      <div class="dialog-footer" :class="{ 'mobile-footer': isMobile }">
        <el-button 
          @click="handleCancelPreview" 
          :disabled="isUploading"
          :size="isMobile ? 'large' : 'default'"
        >
          取消
        </el-button>
        <el-button 
          type="primary" 
          @click="handleConfirmCrop"
          :loading="isUploading"
          class="confirm-btn"
          :size="isMobile ? 'large' : 'default'"
        >
          确认提交
        </el-button>
      </div>
    </template>
  </el-dialog>

  <!-- 走马灯轮播对话框 -->
  <el-dialog
    v-model="carouselVisible"
    title="图片轮播"
    width="90%"
    :close-on-click-modal="true"
    :close-on-press-escape="true"
    :show-close="true"
    class="carousel-dialog"
    @closed="handleCarouselClosed"
  >
    <div class="carousel-container" @keydown="handleKeydown" tabindex="0">
      <div class="carousel-header-info">
        <span class="carousel-counter">
          {{ carouselCurrentIndex + 1 }} / {{ images.length }}
        </span>
      </div>
      <el-carousel
        ref="carouselRef"
        :key="`carousel-${carouselKey}`"
        :initial-index="carouselInitialIndex"
        height="70vh"
        indicator-position="outside"
        arrow="always"
        :interval="0"
        :loop="true"
        @change="handleCarouselChange"
      >
        <el-carousel-item v-for="image in images" :key="image.image_id">
          <div class="carousel-item-content">
            <div class="carousel-image-wrapper" :class="{ 'image-zoomed': imageZoomed }">
            <el-image 
              :src="getOriginImageUrl(image.image_path, currentSearchUserID !== null ? currentSearchUserID : undefined)" 
              fit="contain"
              class="carousel-image"
                :class="{ 'zoom-in': imageZoomed }"
                :preview-src-list="[getOriginImageUrl(image.image_path, currentSearchUserID !== null ? currentSearchUserID : undefined)]"
                :initial-index="0"
                :preview-teleported="true"
                loading="lazy"
                @load="handleImageLoad"
                @error="handleImageError"
                @dblclick="toggleZoom"
            >
              <template #error>
                <div class="image-error">
                  <el-icon><Picture /></el-icon>
                  <span>图片加载失败</span>
                </div>
              </template>
                <template #placeholder>
                  <div class="image-loading">
                    <el-icon class="is-loading"><Loading /></el-icon>
                    <span>加载中...</span>
                  </div>
                </template>
            </el-image>
            </div>
            <div class="carousel-info">
              <div class="carousel-info-header">
              <h3>{{ image.image_name }}</h3>
                <div class="carousel-actions">
                  <el-button
                    circle
                    size="small"
                    :icon="imageZoomed ? ZoomOut : ZoomIn"
                    @click="toggleZoom"
                    title="缩放图片"
                  />
                </div>
              </div>
              <div class="carousel-tags">
                <!-- AI标签 -->
                <el-tag
                  v-for="tag in (image.ai_tags || [])"
                  :key="`ai-${tag}`"
                  type="primary"
                  effect="plain"
                  size="small"
                  round
                >
                  <el-icon class="tag-icon"><MagicStick /></el-icon>
                  {{ tag }}
                </el-tag>
                <!-- 手动标签 -->
                <el-tag
                  v-for="tag in (image.manual_tags || [])"
                  :key="`manual-${tag}`"
                  type="success"
                  effect="plain"
                  size="small"
                  round
                >
                  <el-icon class="tag-icon"><PriceTag /></el-icon>
                  {{ tag }}
                </el-tag>
                <el-text v-if="(!image.ai_tags || image.ai_tags.length === 0) && (!image.manual_tags || image.manual_tags.length === 0)" size="small" type="info">
                  暂无标签
                </el-text>
              </div>
            </div>
          </div>
        </el-carousel-item>
      </el-carousel>
      <div class="carousel-hint" v-if="!isMobile">
        <el-text type="info" size="small">
          <el-icon><InfoFilled /></el-icon>
          使用键盘 ← → 切换图片，ESC 关闭，双击图片放大
        </el-text>
      </div>
      <div class="carousel-hint-mobile" v-else>
        <el-text type="info" size="small">
          <el-icon><InfoFilled /></el-icon>
          左右滑动切换图片，双击放大
        </el-text>
      </div>
    </div>
  </el-dialog>
</template>

<script setup lang="ts">
import ImageCard from '@/components/ImageCard.vue';
import AppHeader from '@/components/AppHeader.vue';
import ImageEditor from '@/components/ImageEidtor.vue';
import { ref, reactive, onMounted, watch, computed, nextTick, onUnmounted } from 'vue'; 
import { ElMessage } from 'element-plus';
import { Search, Picture, ZoomIn, ZoomOut, Loading, InfoFilled, MagicStick, PriceTag } from '@element-plus/icons-vue';
import { ImageAPI, UserAPI } from '@/api';
import { getOriginImageUrl } from '@/utils/url';
import type { ImageItem, PaginationState, User } from '@/type';
import { base64ToFile } from '@/utils/file';
import type { ElCarousel } from 'element-plus';

const searchType = ref('all');
const searchValue = ref('');
const startTime = ref<string | null>(null);
const endTime = ref<string | null>(null);
const searching = ref(false);
const hasSearched = ref(false);
const images = ref<ImageItem[]>([]);

// 区域查询相关状态
const areaProvince = ref('');
const areaCity = ref('');
const areaDistrict = ref('');

// 管理员相关状态
const userInfo = ref(JSON.parse(localStorage.getItem('userInfo') || '{}'));
const isAdmin = computed(() => {
  return userInfo.value.is_admin === true;
});

// 监听用户信息更新事件
const handleUserInfoUpdate = (event: CustomEvent) => {
  console.log('Gallery: 收到用户信息更新事件，刷新界面');
  userInfo.value = event.detail;
  userInfo.value = { ...event.detail };
};

// 移动端检测
const isMobile = computed(() => {
  if (typeof window === 'undefined') return false;
  return window.innerWidth <= 768;
});

// 用户选择相关状态
const selectedUserID = ref<number | null>(null);
const userOptions = ref<User[]>([]);
const loadingUsers = ref(false);
// 实际搜索时使用的用户ID（用于图片URL生成）
const currentSearchUserID = ref<number | null>(null);

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
const carouselCurrentIndex = ref(0);
const carouselRef = ref<InstanceType<typeof ElCarousel> | null>(null);
const imageZoomed = ref(false);
const carouselKey = ref(0); // 用于强制重新渲染轮播组件

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

// 监听走马灯对话框打开，设置焦点
watch(carouselVisible, (newVal) => {
  if (newVal) {
    nextTick(() => {
      const container = document.querySelector('.carousel-container') as HTMLElement;
      if (container) {
        container.focus();
      }
    });
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
  startTime.value = null;
  endTime.value = null;
  areaProvince.value = '';
  areaCity.value = '';
  areaDistrict.value = '';
  selectedUserID.value = null;
  currentSearchUserID.value = null; // 清空搜索时使用的用户ID
  images.value = []; // 清空图片列表
  hasSearched.value = false; // 重置搜索状态
  
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
      const allUsers = response.data.data.users || [];
      // 获取当前用户ID，排除自己
      const currentUserId = userInfo.value.id || parseInt(localStorage.getItem('userID') || '0');
      // 过滤掉当前用户
      userOptions.value = allUsers.filter((user: User) => user.id !== currentUserId);
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
      // 获取当前用户ID，排除自己
      const currentUserId = userInfo.value.id || parseInt(localStorage.getItem('userID') || '0');
      // 根据查询字符串过滤用户，同时排除当前用户
      userOptions.value = allUsers.filter((user: User) => 
        user.id !== currentUserId &&
        (user.username.toLowerCase().includes(query.toLowerCase()) ||
         user.email.toLowerCase().includes(query.toLowerCase()) ||
         user.id.toString().includes(query))
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
  if (searchType.value === 'time') {
    if (!startTime.value || !endTime.value) {
      ElMessage.warning('请选择开始日期和结束日期');
      return;
    }
    if (startTime.value >= endTime.value) {
      ElMessage.warning('结束日期必须晚于开始日期');
      return;
    }
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

  // 保存实际搜索时使用的用户ID
  if (searchType.value === 'user-images' && selectedUserID.value) {
    currentSearchUserID.value = selectedUserID.value;
  } else {
    currentSearchUserID.value = null;
  }

  try {
    let response;
    
    switch (searchType.value) {
      case 'all':
        response = await ImageAPI.getPersonalAllImages(pagination.currentPage, pagination.pageSize);
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
        if (startTime.value && endTime.value) {
          // 将日期转换为时间范围：开始日期 00:00:00，结束日期 23:59:59
          const startDateTime = `${startTime.value} 00:00:00`;
          const endDateTime = `${endTime.value} 23:59:59`;
          response = await ImageAPI.getImagesByTimeRange(
            startDateTime, 
            endDateTime, 
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
        // 标准化图片数据，确保 ai_tags 和 manual_tags 始终是数组
        const normalizedImages = (result.data.images || []).map((img: any) => ({
          ...img,
          ai_tags: Array.isArray(img.ai_tags) ? img.ai_tags : [],
          manual_tags: Array.isArray(img.manual_tags) ? img.manual_tags : [],
        }));
        images.value = normalizedImages;
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
  startTime.value = null;
  endTime.value = null;
  areaProvince.value = '';
  areaCity.value = '';
  areaDistrict.value = '';
  selectedUserID.value = null;
  currentSearchUserID.value = null; // 清空搜索时使用的用户ID
  images.value = [];
  hasSearched.value = false;
  
  pagination.currentPage = 1;
  pagination.pageSize = 12;
  pagination.total = 0;
  
  ElMessage.info('已重置搜索条件');
};

/**
 * 处理删除成功事件
 * @param imageId 被删除的图片ID
 */
const handleSuccess = (imageId: number) => {
  // 从前端数组中移除被删除的图片
  const index = images.value.findIndex(img => img.image_id === imageId);
  if (index !== -1) {
    // 如果删除的图片正在走马灯中显示，需要处理走马灯状态
    const isViewingDeletedImage = carouselVisible.value && carouselCurrentIndex.value === index;
    
    // 移除图片
    images.value.splice(index, 1);
    // 更新总数
    pagination.total = Math.max(0, pagination.total - 1);
    
    // 处理走马灯状态
    if (isViewingDeletedImage) {
      if (images.value.length === 0) {
        // 删除后没有图片了，关闭走马灯
        carouselVisible.value = false;
      } else if (index >= images.value.length) {
        // 删除的是最后一张，显示新的最后一张（索引减1）
        carouselCurrentIndex.value = images.value.length - 1;
        carouselKey.value++; // 强制重新渲染轮播组件
      } else {
        // 删除的不是最后一张，当前索引不变（因为后面的元素前移了，正好显示下一张）
        carouselKey.value++; // 强制重新渲染轮播组件
      }
    } else if (carouselVisible.value && carouselCurrentIndex.value > index) {
      // 如果走马灯显示的是被删除图片之后的图片，索引需要减1
      carouselCurrentIndex.value--;
    }
    
    // 如果当前页没有图片了，且不是第一页，则回到上一页
    if (images.value.length === 0 && pagination.currentPage > 1) {
      pagination.currentPage--;
      // 需要重新查询上一页的数据
      handleSearch(false);
    }
  }
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
      
      // 直接更新图片路径，而不是重新查询
      const newStoreName = response.data.data?.storeName;
      if (newStoreName) {
        const imageIndex = images.value.findIndex(img => img.image_id === currentImageId.value);
        if (imageIndex !== -1) {
          images.value[imageIndex].image_path = newStoreName;
          // 如果图片正在走马灯中显示，需要强制更新
          if (carouselVisible.value && carouselCurrentIndex.value === imageIndex) {
            carouselKey.value++; // 强制重新渲染轮播组件
          }
        }
      } else {
        // 如果没有返回新图片名，则重新查询
        handleSearch(false);
      }
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
    carouselCurrentIndex.value = index;
    imageZoomed.value = false;
    carouselKey.value++; // 递增 key 确保轮播组件重新渲染
    carouselVisible.value = true;
  }
};

/**
 * 处理走马灯切换事件
 */
const handleCarouselChange = (index: number) => {
  carouselCurrentIndex.value = index;
  imageZoomed.value = false; // 切换图片时重置缩放
};

/**
 * 处理键盘事件
 */
const handleKeydown = (event: KeyboardEvent) => {
  if (!carouselVisible.value || !carouselRef.value) return;
  
  switch (event.key) {
    case 'ArrowLeft':
      event.preventDefault();
      carouselRef.value.prev();
      break;
    case 'ArrowRight':
      event.preventDefault();
      carouselRef.value.next();
      break;
    case 'Escape':
      event.preventDefault();
      carouselVisible.value = false;
      break;
  }
};

/**
 * 切换图片缩放
 */
const toggleZoom = () => {
  imageZoomed.value = !imageZoomed.value;
};

/**
 * 处理图片加载成功
 */
const handleImageLoad = () => {
  // 图片加载成功后的处理
};

/**
 * 处理图片加载失败
 */
const handleImageError = () => {
  ElMessage.error('图片加载失败');
};

/**
 * 处理走马灯关闭
 */
const handleCarouselClosed = () => {
  imageZoomed.value = false;
  carouselCurrentIndex.value = 0;
  carouselInitialIndex.value = 0;
};

// 窗口大小变化监听
const handleResize = () => {
  // 触发响应式更新
};

// 组件挂载时，如果是管理员且查询类型是用户图片，加载用户列表
onMounted(() => {
  // 监听用户信息更新事件
  window.addEventListener('userInfoUpdated', handleUserInfoUpdate as EventListener);
  // 如果是管理员且查询类型是用户图片，加载用户列表
  if (isAdmin.value && searchType.value === 'user-images') {
    loadUsers();
  }
  // 页面初始状态显示欢迎消息，不自动搜索
  window.addEventListener('resize', handleResize);
});

// 组件卸载时清理
onUnmounted(() => {
  window.removeEventListener('resize', handleResize);
  window.removeEventListener('userInfoUpdated', handleUserInfoUpdate as EventListener);
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
  position: relative;
  outline: none;
}

.carousel-header-info {
  position: absolute;
  top: -50px;
  right: 60px; /* 为关闭按钮留出空间 */
  z-index: 1;
  pointer-events: none;
}

.carousel-counter {
  display: inline-flex;
  align-items: center;
  padding: 6px 16px;
  background: rgba(64, 158, 255, 0.9);
  color: white;
  border-radius: 20px;
  font-size: 14px;
  font-weight: 600;
  box-shadow: 0 2px 8px rgba(64, 158, 255, 0.3);
  pointer-events: auto;
}

.carousel-item-content {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: linear-gradient(135deg, #f8fbff, #f0f7ff);
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
  animation: fadeIn 0.4s ease-in-out;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.carousel-image-wrapper {
  flex: 1;
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
  overflow: hidden;
  position: relative;
  transition: all 0.3s ease;
}

.carousel-image-wrapper.image-zoomed {
  padding: 10px;
}

.carousel-image {
  width: 100%;
  height: 100%;
  max-width: 100%;
  max-height: 100%;
  object-fit: contain;
  border-radius: 12px;
  transition: transform 0.4s cubic-bezier(0.4, 0, 0.2, 1), opacity 0.3s ease;
  cursor: zoom-in;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
  user-select: none;
}

.carousel-image.zoom-in {
  cursor: zoom-out;
  transform: scale(1.8);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.25);
}

.carousel-image:hover {
  opacity: 0.95;
}

.carousel-info {
  padding: 24px;
  background: linear-gradient(135deg, #ffffff, #f8fbff);
  border-top: 2px solid rgba(64, 158, 255, 0.1);
  backdrop-filter: blur(10px);
}

.carousel-info-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.carousel-info h3 {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
  color: #2c3e50;
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.carousel-actions {
  display: flex;
  gap: 8px;
  margin-left: 16px;
}

.carousel-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  align-items: center;
}

.carousel-tags .tag-icon {
  margin-right: 4px;
  font-size: 12px;
}

.carousel-hint {
  margin-top: 16px;
  text-align: center;
  padding: 12px;
  background: linear-gradient(135deg, #f8fbff, #f0f7ff);
  border-radius: 8px;
  border: 1px solid rgba(64, 158, 255, 0.1);
}

.carousel-hint .el-icon {
  margin-right: 6px;
  vertical-align: middle;
}

/* 图片加载状态 */
.image-loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 16px;
  color: #409EFF;
  font-size: 16px;
  min-height: 400px;
}

.image-loading .el-icon {
  font-size: 48px;
}

.image-loading .is-loading {
  animation: rotating 2s linear infinite;
}

@keyframes rotating {
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
}

/* 走马灯箭头样式优化 */
.carousel-container :deep(.el-carousel__arrow) {
  background-color: rgba(64, 158, 255, 0.85);
  border-radius: 50%;
  width: 48px;
  height: 48px;
  border: 2px solid rgba(255, 255, 255, 0.8);
  box-shadow: 0 4px 12px rgba(64, 158, 255, 0.3);
  transition: all 0.3s ease;
}

.carousel-container :deep(.el-carousel__arrow:hover) {
  background-color: rgba(64, 158, 255, 1);
  transform: scale(1.1);
  box-shadow: 0 6px 16px rgba(64, 158, 255, 0.4);
}

.carousel-container :deep(.el-carousel__arrow--left) {
  left: 20px;
}

.carousel-container :deep(.el-carousel__arrow--right) {
  right: 20px;
}

.carousel-container :deep(.el-carousel__arrow i) {
  font-size: 20px;
  color: white;
}

/* 走马灯指示器样式优化 */
.carousel-container :deep(.el-carousel__indicators) {
  margin-top: 20px;
}

.carousel-container :deep(.el-carousel__indicator) {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background-color: rgba(64, 158, 255, 0.4);
  margin: 0 6px;
  transition: all 0.3s ease;
}

.carousel-container :deep(.el-carousel__indicator:hover) {
  background-color: rgba(64, 158, 255, 0.7);
  transform: scale(1.2);
}

.carousel-container :deep(.el-carousel__indicator.is-active) {
  background-color: #409EFF;
  width: 24px;
  border-radius: 5px;
}

.carousel-container :deep(.el-carousel__button) {
  width: 100%;
  height: 100%;
  background-color: transparent;
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

/* 走马灯对话框关闭按钮样式 */
:deep(.carousel-dialog .el-dialog__header) {
  position: relative;
  z-index: 1;
}

:deep(.carousel-dialog .el-dialog__headerbtn) {
  top: 20px !important;
  right: 20px !important;
  width: 40px !important;
  height: 40px !important;
  z-index: 9999 !important;
  background-color: rgba(255, 255, 255, 0.9) !important;
  border-radius: 50% !important;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15) !important;
  transition: all 0.3s ease !important;
  pointer-events: auto !important;
  cursor: pointer !important;
}

:deep(.carousel-dialog .el-dialog__headerbtn::before) {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: -1;
}

:deep(.carousel-dialog .el-dialog__headerbtn:hover) {
  background-color: rgba(255, 255, 255, 1);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
  transform: scale(1.1);
}

:deep(.carousel-dialog .el-dialog__close) {
  font-size: 20px;
  color: #2c3e50;
  transition: all 0.3s ease;
  font-weight: 600;
}

:deep(.carousel-dialog .el-dialog__close:hover) {
  color: #409EFF;
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
    padding: 15px 10px;
  }

  .image-search-container {
    padding: 20px 15px;
    border-radius: 16px;
  }

  .search-header {
    margin-bottom: 20px;
  }

  .search-header h1 {
    font-size: 24px;
  }

  .search-header p {
    font-size: 14px;
  }

  .search-card {
    padding: 16px;
    border-radius: 12px;
  }

  .search-controls {
    flex-direction: column;
    align-items: stretch;
    gap: 12px;
  }

  .control-group {
    width: 100%;
    margin-bottom: 0;
  }

  .control-group :deep(.el-form-item__label) {
    font-size: 13px;
    margin-bottom: 5px;
  }

  .search-select,
  .search-input,
  .time-picker {
    min-width: 100%;
    width: 100%;
  }

  .time-picker {
    font-size: 12px;
  }

  .time-picker :deep(.el-input__wrapper) {
    width: 100%;
    padding: 3px 5px;
    min-height: 32px;
  }

  .time-picker :deep(.el-input__inner) {
    font-size: 12px;
    line-height: 1.4;
  }

  .time-picker :deep(.el-input__suffix) {
    font-size: 14px;
  }

  .time-picker :deep(.el-input__suffix .el-icon) {
    font-size: 14px;
  }

  .action-buttons {
    margin-left: 0;
    justify-content: stretch;
    width: 100%;
    margin-top: 8px;
    gap: 8px;
  }

  .action-buttons .el-button {
    flex: 1;
    padding: 10px 12px;
    font-size: 14px;
    min-width: 0;
  }

  .search-btn {
    flex: 1.2;
  }

  .reset-btn {
    flex: 0.8;
  }

  .results-header {
    margin-bottom: 16px;
    padding-bottom: 12px;
  }

  .results-header h3 {
    font-size: 16px;
  }

  .images-grid {
    grid-template-columns: repeat(2, 1fr);
    gap: 8px;
  }

  /* 图片卡片在移动端的优化 */
  .images-grid :deep(.image-card) {
    border-radius: 8px;
    width: 100% !important;
    height: auto !important;
    margin: 0 !important;
    min-height: 200px;
  }

  .images-grid :deep(.image-card .image-body-container) {
    height: 160px !important;
  }

  .images-grid :deep(.image-card .card-header) {
    padding: 8px 12px;
  }

  .images-grid :deep(.image-card .card-title) {
    font-size: 13px;
  }

  .images-grid :deep(.image-card .tag-item) {
    font-size: 10px;
    padding: 2px 6px;
    margin: 0;
  }

  .images-grid :deep(.image-card .tags-container) {
    max-height: 40px;
    gap: 4px;
  }

  .images-grid :deep(.image-card .toolbar) {
    top: 3px;
    right: 3px;
    gap: 1px;
    padding: 1px;
    opacity: 0.8; /* 移动端稍微显示工具栏 */
    border-radius: 3px;
    background-color: rgba(0, 0, 0, 0.5);
  }

  .images-grid :deep(.image-card:hover .toolbar),
  .images-grid :deep(.image-card:active .toolbar) {
    opacity: 1;
  }

  .images-grid :deep(.image-card .toolbar .el-button) {
    width: 20px;
    height: 20px;
    padding: 0;
    min-width: 20px;
  }

  .images-grid :deep(.image-card .toolbar .el-button.is-circle) {
    border-radius: 3px;
  }

  .images-grid :deep(.image-card .toolbar .el-icon) {
    font-size: 10px;
  }

  .pagination-footer {
    margin-top: 24px;
    padding: 0 8px;
  }

  .pagination-component {
    justify-content: center;
  }

  .pagination-footer :deep(.el-pagination) {
    justify-content: center;
    flex-wrap: wrap;
  }

  .pagination-footer :deep(.el-pagination .el-pagination__sizes),
  .pagination-footer :deep(.el-pagination .el-pagination__jump) {
    margin-top: 8px;
    width: 100%;
    justify-content: center;
  }

  .pagination-footer :deep(.el-pagination.is-background .btn-prev),
  .pagination-footer :deep(.el-pagination.is-background .btn-next),
  .pagination-footer :deep(.el-pagination.is-background .el-pager li) {
    margin: 0 2px;
  }

  /* 对话框响应式 */
  :deep(.preview-dialog.mobile-dialog .el-dialog) {
    width: 100% !important;
    height: 100vh !important;
    max-height: 100vh !important;
    margin: 0 !important;
    border-radius: 0 !important;
    display: flex;
    flex-direction: column;
  }

  :deep(.preview-dialog.mobile-dialog .el-dialog__header) {
    padding: 10px 12px;
    border-bottom: 1px solid #ebeef5;
    flex-shrink: 0;
  }

  :deep(.preview-dialog.mobile-dialog .el-dialog__title) {
    font-size: 17px;
    font-weight: 600;
  }

  :deep(.preview-dialog.mobile-dialog .el-dialog__headerbtn) {
    top: 10px;
    right: 12px;
    width: 32px;
    height: 32px;
  }

  :deep(.preview-dialog.mobile-dialog .el-dialog__body) {
    padding: 12px;
    flex: 1;
    overflow-y: auto;
    -webkit-overflow-scrolling: touch;
    overscroll-behavior: contain;
  }

  :deep(.preview-dialog.mobile-dialog .el-dialog__footer) {
    padding: 10px 12px;
    flex-shrink: 0;
    border-top: 1px solid #ebeef5;
    background-color: #fff;
  }

  :deep(.preview-dialog .el-dialog),
  :deep(.carousel-dialog .el-dialog) {
    width: 95% !important;
    margin: 2vh auto;
    border-radius: 12px;
  }

  :deep(.preview-dialog .el-dialog__header),
  :deep(.carousel-dialog .el-dialog__header) {
    padding: 16px;
  }

  :deep(.preview-dialog .el-dialog__title),
  :deep(.carousel-dialog .el-dialog__title) {
    font-size: 16px;
  }

  /* 移动端关闭按钮优化 */
  :deep(.carousel-dialog .el-dialog__headerbtn) {
    top: 12px;
    right: 12px;
    width: 36px;
    height: 36px;
    background-color: rgba(255, 255, 255, 0.95);
    border-radius: 50%;
    box-shadow: 0 2px 6px rgba(0, 0, 0, 0.2);
  }

  :deep(.carousel-dialog .el-dialog__headerbtn:active) {
    transform: scale(0.95);
  }

  :deep(.carousel-dialog .el-dialog__close) {
    font-size: 18px;
    font-weight: 600;
  }

  /* 确保计数器不遮挡关闭按钮 */
  .carousel-header-info {
    top: -45px;
    right: 50px;
  }

  .preview-container {
    padding: 8px 0;
  }

  .preview-container.mobile-layout {
    padding: 0;
    height: 100%;
    display: flex;
    flex-direction: column;
  }

  .preview-hint {
    margin-bottom: 12px;
    padding: 12px;
    font-size: 13px;
  }

  .mobile-layout .preview-hint {
    margin-bottom: 10px;
    padding: 10px;
    font-size: 12px;
    border-radius: 8px;
  }

  .preview-image-wrapper {
    min-height: 300px;
    max-height: 400px;
    padding: 12px;
  }

  .mobile-layout .preview-image-wrapper {
    flex: 1;
    min-height: 0;
    max-height: none;
    padding: 10px;
    display: flex;
    align-items: center;
    justify-content: center;
    overflow: hidden;
  }

  .preview-image {
    max-height: 360px;
  }

  .mobile-layout .preview-image {
    max-width: 100%;
    max-height: 100%;
    width: auto;
    height: auto;
    object-fit: contain;
  }

  .dialog-footer {
    padding: 12px 16px;
  }

  .dialog-footer.mobile-footer {
    flex-direction: row;
    justify-content: space-between;
    gap: 10px;
    padding: 10px 12px;
  }

  .dialog-footer.mobile-footer .el-button {
    flex: 1;
    max-width: 45%;
    margin: 0;
    height: 44px;
    font-size: 16px;
  }

  .dialog-footer .el-button {
    padding: 10px 20px;
    font-size: 14px;
  }

  /* 走马灯响应式 */
  .carousel-container :deep(.el-carousel) {
    height: 55vh !important;
  }

  .carousel-header-info {
    top: -45px;
    right: 0;
  }

  .carousel-counter {
    padding: 4px 12px;
    font-size: 12px;
  }

  .carousel-image-wrapper {
    padding: 12px;
  }

  .carousel-info {
    padding: 16px;
  }

  .carousel-info-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
    margin-bottom: 12px;
  }

  .carousel-info h3 {
    font-size: 16px;
    width: 100%;
  }

  .carousel-actions {
    margin-left: 0;
    width: 100%;
    justify-content: flex-end;
  }

  .carousel-tags {
    gap: 6px;
  }

  .carousel-tags .el-tag {
    font-size: 11px;
    padding: 2px 8px;
  }

  .carousel-hint,
  .carousel-hint-mobile {
    margin-top: 12px;
    padding: 10px;
    font-size: 11px;
  }

  .carousel-hint-mobile {
    display: block;
  }

  .carousel-container :deep(.el-carousel__arrow) {
    width: 36px;
    height: 36px;
  }

  .carousel-container :deep(.el-carousel__arrow i) {
    font-size: 16px;
  }

  .carousel-container :deep(.el-carousel__arrow--left) {
    left: 8px;
  }

  .carousel-container :deep(.el-carousel__arrow--right) {
    right: 8px;
  }

  .carousel-container :deep(.el-carousel__indicators) {
    margin-top: 12px;
  }

  .carousel-container :deep(.el-carousel__indicator) {
    width: 8px;
    height: 8px;
    margin: 0 4px;
  }

  .carousel-container :deep(.el-carousel__indicator.is-active) {
    width: 20px;
  }

  .welcome-message,
  .no-results {
    min-height: 300px;
  }

  .welcome-icon {
    font-size: 48px;
  }

  .welcome-message h3 {
    font-size: 18px;
  }

  .welcome-message p {
    font-size: 14px;
    padding: 0 20px;
  }
}

@media (max-width: 480px) {
  .page-container {
    padding: 10px 8px;
  }

  .image-search-container {
    padding: 16px 12px;
    border-radius: 12px;
  }

  .search-header {
    margin-bottom: 16px;
  }

  .search-header h1 {
    font-size: 20px;
  }

  .search-header p {
    font-size: 12px;
  }

  .search-card {
    padding: 12px;
    border-radius: 10px;
  }

  .search-controls {
    gap: 10px;
  }

  .control-group :deep(.el-form-item__label) {
    font-size: 12px;
  }

  .action-buttons {
    flex-direction: row;
    gap: 6px;
  }

  .action-buttons .el-button {
    flex: 1;
    padding: 10px 8px;
    font-size: 13px;
    min-width: 0;
  }

  .search-btn {
    flex: 1.3;
  }

  .reset-btn {
    flex: 0.7;
  }

  .results-header h3 {
    font-size: 14px;
  }

  .images-grid {
    grid-template-columns: repeat(2, 1fr);
    gap: 6px;
  }

  /* 图片卡片在小屏幕的优化 */
  .images-grid :deep(.image-card) {
    border-radius: 6px;
    width: 100% !important;
    height: auto !important;
    margin: 0 !important;
    min-height: 180px;
  }

  .images-grid :deep(.image-card .image-body-container) {
    height: 140px !important;
  }

  .images-grid :deep(.image-card .card-header) {
    padding: 6px 10px;
  }

  .images-grid :deep(.image-card .card-title) {
    font-size: 12px;
  }

  .images-grid :deep(.image-card .tag-item) {
    font-size: 9px;
    padding: 1px 4px;
    margin: 0;
  }

  .images-grid :deep(.image-card .tags-container) {
    max-height: 36px;
    gap: 3px;
  }

  .images-grid :deep(.image-card .toolbar) {
    top: 2px;
    right: 2px;
    gap: 1px;
    padding: 1px;
    opacity: 0.8; /* 移动端稍微显示工具栏 */
    border-radius: 2px;
    background-color: rgba(0, 0, 0, 0.5);
  }

  .images-grid :deep(.image-card:hover .toolbar),
  .images-grid :deep(.image-card:active .toolbar) {
    opacity: 1;
  }

  .images-grid :deep(.image-card .toolbar .el-button) {
    width: 18px;
    height: 18px;
    padding: 0;
    min-width: 18px;
  }

  .images-grid :deep(.image-card .toolbar .el-button.is-circle) {
    border-radius: 2px;
  }

  .images-grid :deep(.image-card .toolbar .el-icon) {
    font-size: 9px;
  }

  .pagination-footer {
    margin-top: 16px;
  }

  .pagination-footer :deep(.el-pagination) {
    font-size: 12px;
  }

  .pagination-footer :deep(.el-pagination .btn-prev),
  .pagination-footer :deep(.el-pagination .btn-next),
  .pagination-footer :deep(.el-pagination .el-pager li) {
    min-width: 28px;
    height: 28px;
    line-height: 28px;
    font-size: 12px;
  }

  /* 对话框响应式 */
  :deep(.preview-dialog.mobile-dialog .el-dialog__header) {
    padding: 8px 10px;
  }

  :deep(.preview-dialog.mobile-dialog .el-dialog__title) {
    font-size: 16px;
  }

  :deep(.preview-dialog.mobile-dialog .el-dialog__body) {
    padding: 10px;
  }

  :deep(.preview-dialog.mobile-dialog .el-dialog__footer) {
    padding: 8px 10px;
  }

  .mobile-layout .preview-hint {
    padding: 8px;
    font-size: 11px;
    margin-bottom: 8px;
  }

  .mobile-layout .preview-image-wrapper {
    padding: 8px;
  }

  .mobile-footer .el-button {
    height: 42px;
    font-size: 15px;
  }

  :deep(.preview-dialog .el-dialog),
  :deep(.carousel-dialog .el-dialog) {
    width: 98% !important;
    margin: 1vh auto;
    border-radius: 10px;
  }

  :deep(.preview-dialog .el-dialog__header),
  :deep(.carousel-dialog .el-dialog__header) {
    padding: 12px;
  }

  :deep(.preview-dialog .el-dialog__title),
  :deep(.carousel-dialog .el-dialog__title) {
    font-size: 14px;
  }

  /* 小屏幕关闭按钮优化 */
  :deep(.carousel-dialog .el-dialog__headerbtn) {
    top: 10px;
    right: 10px;
    width: 32px;
    height: 32px;
    background-color: rgba(255, 255, 255, 0.95);
    border-radius: 50%;
    box-shadow: 0 2px 6px rgba(0, 0, 0, 0.2);
  }

  :deep(.carousel-dialog .el-dialog__headerbtn:active) {
    transform: scale(0.95);
  }

  :deep(.carousel-dialog .el-dialog__close) {
    font-size: 16px;
    font-weight: 600;
  }

  /* 确保计数器不遮挡关闭按钮 */
  .carousel-header-info {
    top: -40px;
    right: 45px;
  }

  .preview-hint {
    padding: 10px;
    font-size: 12px;
  }

  .preview-image-wrapper {
    min-height: 250px;
    max-height: 350px;
    padding: 8px;
  }

  .preview-image {
    max-height: 320px;
  }

  .dialog-footer {
    padding: 10px 12px;
    flex-direction: column;
    gap: 8px;
  }

  .dialog-footer .el-button {
    width: 100%;
    padding: 12px;
  }

  /* 走马灯响应式 */
  .carousel-container :deep(.el-carousel) {
    height: 50vh !important;
  }

  .carousel-header-info {
    top: -40px;
  }

  .carousel-counter {
    padding: 3px 10px;
    font-size: 11px;
  }

  .carousel-image-wrapper {
    padding: 8px;
  }

  .carousel-info {
    padding: 12px;
  }

  .carousel-info h3 {
    font-size: 14px;
  }

  .carousel-tags {
    gap: 4px;
  }

  .carousel-tags .el-tag {
    font-size: 10px;
    padding: 2px 6px;
  }

  .carousel-hint {
    display: none; /* 小屏幕隐藏桌面提示 */
  }

  .carousel-hint-mobile {
    margin-top: 8px;
    padding: 8px;
    font-size: 10px;
    display: block;
  }

  .carousel-container :deep(.el-carousel__arrow) {
    width: 32px;
    height: 32px;
  }

  .carousel-container :deep(.el-carousel__arrow i) {
    font-size: 14px;
  }

  .carousel-container :deep(.el-carousel__arrow--left) {
    left: 4px;
  }

  .carousel-container :deep(.el-carousel__arrow--right) {
    right: 4px;
  }

  .carousel-container :deep(.el-carousel__indicators) {
    margin-top: 8px;
  }

  .carousel-container :deep(.el-carousel__indicator) {
    width: 6px;
    height: 6px;
    margin: 0 3px;
  }

  .carousel-container :deep(.el-carousel__indicator.is-active) {
    width: 16px;
  }

  .welcome-message,
  .no-results {
    min-height: 250px;
  }

  .welcome-icon {
    font-size: 40px;
  }

  .welcome-message h3 {
    font-size: 16px;
  }

  .welcome-message p {
    font-size: 12px;
    padding: 0 16px;
  }

  .image-loading {
    min-height: 300px;
    font-size: 14px;
  }

  .image-loading .el-icon {
    font-size: 40px;
  }

  .time-picker {
    font-size: 10px;
  }

  .time-picker :deep(.el-input__wrapper) {
    padding: 2px 4px !important;
    min-height: 28px !important;
  }

  .time-picker :deep(.el-input__inner) {
    font-size: 10px !important;
    line-height: 1.3;
    padding: 0 2px;
  }

  .time-picker :deep(.el-input__suffix) {
    font-size: 12px;
  }

  .time-picker :deep(.el-input__suffix .el-icon) {
    font-size: 12px;
  }
}

/* 超小屏幕优化 (小于 360px) */
@media (max-width: 360px) {
  .page-container {
    padding: 8px 4px;
  }

  .image-search-container {
    padding: 12px 8px;
  }

  .search-header h1 {
    font-size: 18px;
  }

  .images-grid {
    grid-template-columns: 1fr;
    gap: 8px;
  }

  /* 图片卡片在超小屏幕的优化 */
  .images-grid :deep(.image-card) {
    border-radius: 6px;
    width: 100% !important;
    height: auto !important;
    margin: 0 !important;
    min-height: 200px;
  }

  .images-grid :deep(.image-card .image-body-container) {
    height: 180px !important;
  }

  .images-grid :deep(.image-card .card-header) {
    padding: 8px 12px;
  }

  .images-grid :deep(.image-card .card-title) {
    font-size: 14px;
  }

  .images-grid :deep(.image-card .toolbar) {
    top: 2px;
    right: 2px;
    gap: 1px;
    padding: 1px;
    opacity: 0.8;
    border-radius: 2px;
    background-color: rgba(0, 0, 0, 0.5);
  }

  .images-grid :deep(.image-card .toolbar .el-button) {
    width: 18px;
    height: 18px;
    padding: 0;
    min-width: 18px;
  }

  .images-grid :deep(.image-card .toolbar .el-button.is-circle) {
    border-radius: 2px;
  }

  .images-grid :deep(.image-card .toolbar .el-icon) {
    font-size: 9px;
  }

  .action-buttons {
    gap: 6px;
  }

  .action-buttons .el-button {
    padding: 10px 6px;
    font-size: 12px;
  }

  .carousel-container :deep(.el-carousel) {
    height: 45vh !important;
  }

  .carousel-info {
    padding: 10px;
  }

  .carousel-info h3 {
    font-size: 13px;
  }
}

/* 触摸设备优化 */
@media (hover: none) and (pointer: coarse) {
  .carousel-container :deep(.el-carousel__arrow) {
    width: 44px;
    height: 44px;
    background-color: rgba(64, 158, 255, 0.9);
  }

  .action-buttons .el-button {
    min-height: 44px; /* 触摸友好的按钮高度 */
  }

  .search-btn,
  .reset-btn {
    min-height: 44px;
  }

  .carousel-image {
    cursor: default; /* 移动端移除鼠标样式 */
  }

  /* 触摸设备上工具栏始终可见 */
  .images-grid :deep(.image-card .toolbar) {
    opacity: 0.7 !important;
  }

  .images-grid :deep(.image-card:active .toolbar) {
    opacity: 1 !important;
  }
}
</style>