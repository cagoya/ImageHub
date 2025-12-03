  <template>
    <el-dialog
      v-model="dialogVisible"
      title="图片编辑"
      :width="isMobile ? '100%' : '900px'"
      :close-on-click-modal="false"
      :fullscreen="isMobile"
      :show-close="!isMobile"
      @closed="handleDialogClosed"
      class="image-editor-dialog"
      :class="{ 'mobile-dialog': isMobile }"
    >
      <div class="editor-layout" :class="{ 'mobile-layout': isMobile }">
        <!-- 裁剪区域 -->
        <div class="cropper-section">
          <div class="cropper-container" :style="{ height: isMobile ? cropperHeight + 'px' : containerHeight + 'px' }">
            <canvas ref="canvasRef" style="max-width: 100%; display: none;"></canvas>
            <img ref="imageRef" :src="displayImageUrl" alt="待编辑图片" />
          </div>
        </div>

        <!-- 色调调整面板 -->
        <div class="adjustment-panel">
          <div class="panel-title">色调调整</div>
          
          <div class="adjustment-item">
            <div class="adjustment-label">
              <span>亮度</span>
              <span class="adjustment-value">{{ formatValue(adjustments.brightness) }}</span>
            </div>
            <el-slider
              v-model="adjustments.brightness"
              :min="-100"
              :max="100"
              :step="0.1"
              :show-tooltip="false"
              @input="applyAdjustments"
            />
          </div>

          <div class="adjustment-item">
            <div class="adjustment-label">
              <span>对比度</span>
              <span class="adjustment-value">{{ formatValue(adjustments.contrast) }}</span>
            </div>
            <el-slider
              v-model="adjustments.contrast"
              :min="-100"
              :max="100"
              :step="0.1"
              :show-tooltip="false"
              @input="applyAdjustments"
            />
          </div>

          <div class="adjustment-item">
            <div class="adjustment-label">
              <span>饱和度</span>
              <span class="adjustment-value">{{ formatValue(adjustments.saturation) }}</span>
            </div>
            <el-slider
              v-model="adjustments.saturation"
              :min="-100"
              :max="100"
              :step="0.1"
              :show-tooltip="false"
              @input="applyAdjustments"
            />
          </div>

          <div class="adjustment-item">
            <div class="adjustment-label">
              <span>色相</span>
              <span class="adjustment-value">{{ formatValue(adjustments.hue) }}°</span>
            </div>
            <el-slider
              v-model="adjustments.hue"
              :min="0"
              :max="360"
              :step="0.1"
              :show-tooltip="false"
              @input="applyAdjustments"
            />
          </div>

          <el-button 
            type="default" 
            @click="resetAdjustments" 
            class="reset-button"
          >
            重置调整
          </el-button>
        </div>
      </div>
  
      <template #footer>
        <div class="dialog-footer" :class="{ 'mobile-footer': isMobile }">
          <el-button @click="dialogVisible = false" :size="isMobile ? 'large' : 'default'">
            {{ isMobile ? '取消' : '取消' }}
          </el-button>
          <el-button 
            type="primary" 
            :loading="isCropping" 
            @click="handleCrop"
            :size="isMobile ? 'large' : 'default'"
          >
            确定
          </el-button>
        </div>
      </template>
    </el-dialog>
  </template>
  
  <script setup lang="ts">
import { ref, watch, onMounted, nextTick, computed, onUnmounted } from 'vue';
import 'cropperjs/dist/cropper.css';
import Cropper from 'cropperjs';
import { ElMessage } from 'element-plus';
import { rgbToHsl, hslToRgb } from '@/utils/color';

// 声明 CamanJS 全局类型
declare global {
  interface Window {
    Caman: any;
  }
}
  
  // --- Props & Emits ---
  interface Props {
    // 图片的 URL 或 Base64 字符串
    imageUrl: string;
    // 容器高度，方便控制裁剪框大小
    containerHeight?: number;
    // 是否可见（通过 v-model:visible 绑定）
    visible: boolean;
    // 裁剪的宽高比（可选，例如：1/1, 16/9, 4/3）
    aspectRatio?: number;
  }
  
  const props = withDefaults(defineProps<Props>(), {
    containerHeight: 400,
    aspectRatio: NaN, // NaN 表示自由比例
  });
  
  interface Emits {
    (e: 'update:visible', value: boolean): void;
    // 裁剪成功后，将裁剪后的 base64 数据传给父组件
    (e: 'cropped', base64: string): void;
  }
  const emit = defineEmits<Emits>();
  
  // --- State ---
  const dialogVisible = ref(props.visible);
  const imageRef = ref<HTMLImageElement | null>(null);
  const canvasRef = ref<HTMLCanvasElement | null>(null);
  const cropper = ref<Cropper | null>(null);
  const isCropping = ref(false);
  const displayImageUrl = ref(props.imageUrl);
  const originalImageUrl = ref(props.imageUrl);
  
  // CamanJS 实例
  let camanInstance: any = null;
  
  // 防抖定时器
  let debounceTimer: ReturnType<typeof setTimeout> | null = null;
  
  // 性能优化：缓存原始图片数据
  let cachedOriginalImage: HTMLImageElement | null = null;
  let cachedImageDimensions: { width: number; height: number } | null = null;
  
  // 预览分辨率限制（最大宽度或高度，保持宽高比）
  // PC端使用较高分辨率，移动端使用较低分辨率以提升性能
  const PREVIEW_MAX_SIZE_PC = 1920; // PC端预览时最大尺寸
  const PREVIEW_MAX_SIZE_MOBILE = 800; // 移动端预览时最大尺寸（屏幕小，性能差）
  const isProcessing = ref(false);
  let rafId: number | null = null;
  
  // 移动端检测
  const isMobile = computed(() => {
    if (typeof window === 'undefined') return false;
    return window.innerWidth <= 768;
  });
  
  // 移动端裁剪区域高度（动态计算）
  const cropperHeight = ref(250);
  
  // 计算移动端裁剪区域高度
  const calculateCropperHeight = () => {
    if (typeof window === 'undefined') return;
    const viewportHeight = window.innerHeight;
    // 减去对话框头部、调整面板、底部按钮和间距的高度
    // 头部约 45px，调整面板约 320px（更紧凑的布局），底部按钮约 65px，间距约 12px，对话框内边距约 24px
    const reservedHeight = 466;
    const availableHeight = viewportHeight - reservedHeight;
    // 确保最小高度和最大高度，移动端裁剪区域应该更小，保证内容能完整显示
    cropperHeight.value = Math.max(180, Math.min(260, availableHeight));
  };
  
  // 窗口大小变化监听
  const handleResize = () => {
    if (isMobile.value) {
      calculateCropperHeight();
    }
  };
  
  onMounted(() => {
    window.addEventListener('resize', handleResize);
    if (isMobile.value) {
      calculateCropperHeight();
    }
  });
  
  onUnmounted(() => {
    window.removeEventListener('resize', handleResize);
    // 清理防抖定时器
    if (debounceTimer) {
      clearTimeout(debounceTimer);
      debounceTimer = null;
    }
    // 清理 requestAnimationFrame
    if (rafId !== null) {
      cancelAnimationFrame(rafId);
      rafId = null;
    }
    // 清理 CamanJS 实例
    if (camanInstance) {
      try {
        camanInstance.revert();
      } catch (e) {
        // 忽略清理错误
      }
      camanInstance = null;
    }
    // 清理缓存
    cachedOriginalImage = null;
    cachedImageDimensions = null;
  });
  
  // 监听移动端状态变化
  watch(isMobile, (val) => {
    if (val) {
      nextTick(() => {
        calculateCropperHeight();
      });
    }
  });
  
  // 色调调整参数（使用小数）
  const adjustments = ref({
    brightness: 0,    // -100 到 100，支持小数
    contrast: 0,      // -100 到 100，支持小数
    saturation: 0,    // -100 到 100，支持小数
    hue: 0,           // 0 到 360，支持小数
  });
  
  // 格式化数值显示（保留一位小数，如果是整数则不显示小数点）
  const formatValue = (value: number): string => {
    const rounded = Math.round(value * 10) / 10;
    // 如果是整数，返回整数格式；否则返回一位小数
    if (rounded % 1 === 0) {
      return rounded.toString();
    }
    return rounded.toFixed(1);
  };
  
  // 双向绑定 visible
  watch(
    () => props.visible,
    (val) => {
      dialogVisible.value = val;
      if (val && props.imageUrl) {
        // 清除保存的裁剪数据（新会话）
        savedCropData.value = null;
        // 重置调整参数和图片URL
        originalImageUrl.value = props.imageUrl;
        displayImageUrl.value = props.imageUrl;
        // 清除缓存，重新加载图片
        cachedOriginalImage = null;
        cachedImageDimensions = null;
        resetAdjustments();
        // 对话框显示时，确保DOM更新完毕后初始化 Cropper
        nextTick(() => {
          initCropper();
        });
      }
    }
  );
  watch(dialogVisible, (val) => {
    emit('update:visible', val);
  });
  
  // 保存裁剪框数据，用于在重新初始化后恢复位置
  const savedCropData = ref<{
    x: number;
    y: number;
    width: number;
    height: number;
  } | null>(null);

  /**
   * 初始化 Cropper 实例
   */
  const initCropper = (restoreCropData = false) => {
    if (!imageRef.value) {
      ElMessage.error('图片DOM元素未找到。');
      return;
    }
  
    // 如果已经有实例，先保存裁剪框数据
    if (cropper.value && restoreCropData) {
      try {
        const cropData = cropper.value.getData();
        savedCropData.value = {
          x: cropData.x,
          y: cropData.y,
          width: cropData.width,
          height: cropData.height,
        };
      } catch (e) {
        // 如果获取数据失败，忽略
        console.warn('Failed to get crop data:', e);
      }
    }
  
    // 如果已经有实例，先销毁
    if (cropper.value) {
      cropper.value.destroy();
      cropper.value = null;
    }
  
    // 初始化 Cropper
    cropper.value = new Cropper(imageRef.value, {
      aspectRatio: props.aspectRatio,
      viewMode: 1, // 视图模式，1：限制裁剪框不能超出画布
      // 裁剪框可以移动
      movable: true,
      // 图片可以移动
      cropBoxMovable: true,
      // 裁剪框可以调整大小
      cropBoxResizable: true,
      // 自动显示裁剪框
      autoCrop: true,
      // 鼠标拖动模式：'move' | 'crop'
      dragMode: 'move',
      // 缩放相关
      zoomable: true,
      zoomOnTouch: true,
      zoomOnWheel: !isMobile.value, // 移动端禁用滚轮缩放
      // 旋转相关（移动端通常不需要）
      rotatable: false,
      // 响应式
      responsive: true,
      // 恢复裁剪区域
      restore: false,
      // 检查方向
      checkOrientation: true,
      // 模态框背景
      modal: true,
      // 引导线
      guides: true,
      // 中心指示器
      center: true,
      // 高亮
      highlight: true,
      // 背景
      background: true,
      // 自动裁剪区域
      autoCropArea: 0.8,
      // 移动端优化：调整最小裁剪区域
      minCropBoxWidth: isMobile.value ? 50 : 0,
      minCropBoxHeight: isMobile.value ? 50 : 0,
      // 初始化完成后恢复裁剪框位置
      ready: () => {
        if (restoreCropData && savedCropData.value && cropper.value) {
          try {
            // getData() 和 setData() 都使用相对于原始图片的坐标
            // 由于图片尺寸不变，可以直接恢复
            cropper.value.setData(savedCropData.value);
          } catch (e) {
            // 如果恢复失败，忽略
            console.warn('Failed to restore crop data:', e);
          }
        }
      },
    });
  };
  
  /**
   * 手动应用色相调整（CamanJS可能不支持直接的hue方法）
   * 移动端优化：分批处理像素以避免阻塞主线程
   */
  const applyHueAdjustment = (canvas: HTMLCanvasElement, hueRotation: number) => {
    const ctx = canvas.getContext('2d');
    if (!ctx) return;

    const imageData = ctx.getImageData(0, 0, canvas.width, canvas.height);
    const data = imageData.data;
    const totalPixels = data.length / 4;

    // 移动端使用分批处理以避免长时间阻塞主线程
    if (isMobile.value && totalPixels > 100000) {
      // 大图分批处理
      const batchSize = 50000; // 每批处理5万像素
      let currentIndex = 0;

      const processBatch = () => {
        const endIndex = Math.min(currentIndex + batchSize, totalPixels);
        
        for (let i = currentIndex; i < endIndex; i++) {
          const pixelIndex = i * 4;
          const r = data[pixelIndex] ?? 0;
          const g = data[pixelIndex + 1] ?? 0;
          const b = data[pixelIndex + 2] ?? 0;

          // 转换到 HSL (返回格式: [h: 0-360, s: 0-100, l: 0-100])
          const hsl = rgbToHsl(r, g, b);
          
          // 应用色相旋转
          let newHue = (hsl[0] + hueRotation) % 360;
          if (newHue < 0) newHue += 360;
          
          // 转换回 RGB (输入格式: [h: 0-360, s: 0-100, l: 0-100])
          const rgb = hslToRgb(newHue, hsl[1], hsl[2]);
          
          data[pixelIndex] = Math.max(0, Math.min(255, rgb[0]));
          data[pixelIndex + 1] = Math.max(0, Math.min(255, rgb[1]));
          data[pixelIndex + 2] = Math.max(0, Math.min(255, rgb[2]));
        }

        currentIndex = endIndex;

        if (currentIndex < totalPixels) {
          // 使用 setTimeout 让出主线程，继续处理下一批
          setTimeout(processBatch, 0);
        } else {
          // 处理完成，更新canvas
          ctx.putImageData(imageData, 0, 0);
        }
      };

      processBatch();
    } else {
      // PC端或小图直接处理
      for (let i = 0; i < data.length; i += 4) {
        const r = data[i] ?? 0;
        const g = data[i + 1] ?? 0;
        const b = data[i + 2] ?? 0;

        // 转换到 HSL (返回格式: [h: 0-360, s: 0-100, l: 0-100])
        const hsl = rgbToHsl(r, g, b);
        
        // 应用色相旋转
        let newHue = (hsl[0] + hueRotation) % 360;
        if (newHue < 0) newHue += 360;
        
        // 转换回 RGB (输入格式: [h: 0-360, s: 0-100, l: 0-100])
        const rgb = hslToRgb(newHue, hsl[1], hsl[2]);
        
        data[i] = Math.max(0, Math.min(255, rgb[0]));
        data[i + 1] = Math.max(0, Math.min(255, rgb[1]));
        data[i + 2] = Math.max(0, Math.min(255, rgb[2]));
      }

      ctx.putImageData(imageData, 0, 0);
    }
  };

  /**
   * 计算预览尺寸（保持宽高比，限制最大尺寸）
   * 移动端使用更小的预览尺寸以提升性能
   */
  const calculatePreviewSize = (width: number, height: number): { width: number; height: number } => {
    const maxSize = isMobile.value ? PREVIEW_MAX_SIZE_MOBILE : PREVIEW_MAX_SIZE_PC;
    
    if (width <= maxSize && height <= maxSize) {
      return { width, height };
    }
    
    const ratio = Math.min(maxSize / width, maxSize / height);
    return {
      width: Math.round(width * ratio),
      height: Math.round(height * ratio)
    };
  };

  /**
   * 加载并缓存原始图片
   */
  const loadAndCacheOriginalImage = (): Promise<HTMLImageElement> => {
    return new Promise((resolve, reject) => {
      // 如果已缓存，直接返回
      if (cachedOriginalImage && cachedOriginalImage.complete) {
        resolve(cachedOriginalImage);
        return;
      }

      const img = new Image();
      img.crossOrigin = 'anonymous';
      
      img.onload = () => {
        cachedOriginalImage = img;
        cachedImageDimensions = { width: img.width, height: img.height };
        resolve(img);
      };
      
      img.onerror = () => {
        reject(new Error('图片加载失败'));
      };
      
      img.src = originalImageUrl.value;
    });
  };

  /**
   * 使用 CamanJS 应用色调调整（性能优化版本）
   */
  const doApplyAdjustments = () => {
    if (!imageRef.value || typeof window === 'undefined' || !window.Caman) {
      console.warn('CamanJS not available');
      return;
    }

    // 如果正在处理中，跳过本次请求
    if (isProcessing.value) {
      return;
    }

    // 如果所有调整都是0，直接使用原始图片
    if (
      adjustments.value.brightness === 0 &&
      adjustments.value.contrast === 0 &&
      adjustments.value.saturation === 0 &&
      adjustments.value.hue === 0
    ) {
      displayImageUrl.value = originalImageUrl.value;
      nextTick(() => {
        initCropper(true);
      });
      return;
    }

    isProcessing.value = true;

    // 使用 requestAnimationFrame 优化渲染时机
    // 移动端可以延迟一帧以让UI更流畅
    if (rafId !== null) {
      cancelAnimationFrame(rafId);
    }

    const processCallback = () => {
      rafId = null;
      processImageWithCaman();
    };

    // 移动端延迟一帧处理，让UI更流畅
    if (isMobile.value) {
      rafId = requestAnimationFrame(() => {
        rafId = requestAnimationFrame(processCallback);
      });
    } else {
      rafId = requestAnimationFrame(processCallback);
    }
  };

  /**
   * 使用 CamanJS 处理图片（实际处理函数）
   */
  const processImageWithCaman = async () => {
    try {
      // 加载并缓存原始图片
      const originalImg = await loadAndCacheOriginalImage();
      
      if (!canvasRef.value || !cachedImageDimensions) {
        isProcessing.value = false;
        return;
      }

      const canvas = canvasRef.value;
      const ctx = canvas.getContext('2d');
      if (!ctx) {
        isProcessing.value = false;
        return;
      }

      // 计算预览尺寸（大图缩小处理，提升性能）
      const previewSize = calculatePreviewSize(
        cachedImageDimensions.width,
        cachedImageDimensions.height
      );

      // 设置 canvas 尺寸为预览尺寸
      canvas.width = previewSize.width;
      canvas.height = previewSize.height;

      // 使用高质量缩放绘制到预览尺寸
      // 移动端使用medium质量以提升性能，PC端使用high质量
      ctx.imageSmoothingEnabled = true;
      ctx.imageSmoothingQuality = isMobile.value ? 'medium' : 'high';
      ctx.drawImage(originalImg, 0, 0, previewSize.width, previewSize.height);

      // 销毁之前的 CamanJS 实例
      if (camanInstance) {
        try {
          camanInstance.revert();
        } catch (e) {
          // 忽略错误
        }
      }

      // 使用 CamanJS 处理
      window.Caman(canvas, function(this: any) {
        this.revert(false); // 先恢复到原始状态
        
        // 应用调整
        // CamanJS 的 brightness 和 contrast 范围是 -100 到 100
        if (adjustments.value.brightness !== 0) {
          this.brightness(adjustments.value.brightness);
        }
        
        if (adjustments.value.contrast !== 0) {
          this.contrast(adjustments.value.contrast);
        }
        
        // CamanJS 的 saturation 范围是 -100 到 100
        if (adjustments.value.saturation !== 0) {
          this.saturation(adjustments.value.saturation);
        }
        
        // 如果 hue 不为 0，需要手动处理色相
        const needsHueAdjustment = adjustments.value.hue !== 0;
        
          // 渲染处理结果
          this.render(() => {
            // 如果需要进行色相调整，在渲染后手动处理
            if (needsHueAdjustment) {
              applyHueAdjustment(canvas, adjustments.value.hue);
            }
            
            // 处理完成，更新显示的图片 URL
            // 移动端使用较低的JPEG质量（0.85）以减少文件大小和内存占用，提升性能
            // PC端使用PNG格式保证质量
            const imageFormat = isMobile.value ? 'image/jpeg' : 'image/png';
            const imageQuality = isMobile.value ? 0.85 : undefined;
            displayImageUrl.value = canvas.toDataURL(imageFormat, imageQuality);
            
            // 保存 CamanJS 实例
            camanInstance = this;
            
            isProcessing.value = false;
            
            // 重新初始化 cropper，并恢复裁剪框位置
            nextTick(() => {
              initCropper(true);
            });
          });
      });
    } catch (error) {
      console.error('CamanJS processing error:', error);
      ElMessage.error('图片处理失败');
      isProcessing.value = false;
    }
  };

  // 保存上一次的调整参数，用于检测变化
  const lastAdjustments = ref({
    brightness: 0,
    contrast: 0,
    saturation: 0,
    hue: 0,
  });

  /**
   * 计算动态防抖时间（根据图片大小和设备类型调整）
   * 移动端使用更长的防抖时间以提升性能和减少卡顿
   */
  const getDebounceTime = (): number => {
    if (!cachedImageDimensions) {
      // 移动端默认更长防抖时间
      return isMobile.value ? 300 : 150;
    }
    
    const pixelCount = cachedImageDimensions.width * cachedImageDimensions.height;
    
    // 移动端使用更长的防抖时间
    if (isMobile.value) {
      if (pixelCount > 2000000) { // 大于200万像素
        return 400;
      } else if (pixelCount > 1000000) { // 大于100万像素
        return 300;
      } else {
        return 250;
      }
    } else {
      // PC端防抖时间
      if (pixelCount > 2000000) { // 大于200万像素
        return 200;
      } else if (pixelCount > 1000000) { // 大于100万像素
        return 150;
      } else {
        return 100;
      }
    }
  };

  /**
   * 应用色调调整（优化版防抖，动态调整延迟时间）
   */
  const applyAdjustments = () => {
    // 检查参数是否真的变化了
    const current = adjustments.value;
    const last = lastAdjustments.value;
    
    const hasChanged = 
      current.brightness !== last.brightness ||
      current.contrast !== last.contrast ||
      current.saturation !== last.saturation ||
      current.hue !== last.hue;
    
    if (!hasChanged) {
      return; // 参数没有变化，不需要处理
    }
    
    // 更新上一次的参数值
    lastAdjustments.value = {
      brightness: current.brightness,
      contrast: current.contrast,
      saturation: current.saturation,
      hue: current.hue,
    };
    
    // 清除之前的定时器
    if (debounceTimer) {
      clearTimeout(debounceTimer);
    }
    
    // 根据图片大小动态调整防抖时间
    const debounceTime = getDebounceTime();
    
    // 设置新的定时器
    debounceTimer = setTimeout(() => {
      doApplyAdjustments();
      debounceTimer = null;
    }, debounceTime);
  };

  /**
   * 重置色调调整
   */
  const resetAdjustments = () => {
    adjustments.value = {
      brightness: 0,
      contrast: 0,
      saturation: 0,
      hue: 0,
    };
    // 同时重置上一次的参数值
    lastAdjustments.value = {
      brightness: 0,
      contrast: 0,
      saturation: 0,
      hue: 0,
    };
    displayImageUrl.value = originalImageUrl.value;
    
    // 清理 CamanJS 实例
    if (camanInstance) {
      try {
        camanInstance.revert();
      } catch (e) {
        // 忽略错误
      }
      camanInstance = null;
    }
    
    nextTick(() => {
      // 如果已经有 Cropper 实例，重置时也保持裁剪框位置
      const shouldRestore = cropper.value !== null;
      initCropper(shouldRestore);
    });
  };

  /**
   * 处理确定裁剪按钮点击（使用原图进行最终处理）
   */
  const handleCrop = async () => {
    if (!cropper.value) {
      ElMessage.error('裁剪器未初始化！');
      return;
    }
  
    isCropping.value = true;
    try {
      // 如果有调整参数，需要先对原图应用调整，然后再裁剪
      const hasAdjustments = 
        adjustments.value.brightness !== 0 ||
        adjustments.value.contrast !== 0 ||
        adjustments.value.saturation !== 0 ||
        adjustments.value.hue !== 0;

      if (hasAdjustments && cachedOriginalImage && cachedImageDimensions) {
        // 使用原图进行最终处理
        await processOriginalImageForCrop();
      } else {
        // 没有调整，直接裁剪预览图
        const base64 = cropper.value.getCroppedCanvas().toDataURL('image/png');
        emit('cropped', base64);
        dialogVisible.value = false;
        ElMessage.success('图片编辑成功！');
      }
    } catch (error) {
      ElMessage.error('图片编辑失败！');
      console.error('编辑错误:', error);
    } finally {
      isCropping.value = false;
    }
  };

  /**
   * 处理原图以进行最终裁剪（使用原图尺寸）
   */
  const processOriginalImageForCrop = async () => {
    if (!cachedOriginalImage || !cachedImageDimensions || !canvasRef.value) {
      ElMessage.error('图片数据不可用');
      return;
    }

    try {
      // 创建临时canvas用于处理原图
      const tempCanvas = document.createElement('canvas');
      tempCanvas.width = cachedImageDimensions.width;
      tempCanvas.height = cachedImageDimensions.height;
      const tempCtx = tempCanvas.getContext('2d');
      if (!tempCtx) {
        ElMessage.error('无法创建画布');
        return;
      }

      // 绘制原图
      tempCtx.drawImage(cachedOriginalImage, 0, 0);

      // 使用 CamanJS 处理原图
      await new Promise<void>((resolve, reject) => {
        window.Caman(tempCanvas, function(this: any) {
          this.revert(false);
          
          if (adjustments.value.brightness !== 0) {
            this.brightness(adjustments.value.brightness);
          }
          
          if (adjustments.value.contrast !== 0) {
            this.contrast(adjustments.value.contrast);
          }
          
          if (adjustments.value.saturation !== 0) {
            this.saturation(adjustments.value.saturation);
          }
          
          const needsHueAdjustment = adjustments.value.hue !== 0;
          
          this.render(() => {
            if (needsHueAdjustment) {
              applyHueAdjustment(tempCanvas, adjustments.value.hue);
            }
            
            // 获取处理后的图片URL
            const processedUrl = tempCanvas.toDataURL('image/png');
            
            // 创建新的图片元素用于裁剪
            const processedImg = new Image();
            processedImg.crossOrigin = 'anonymous';
            processedImg.onload = () => {
              // 临时更新图片源，然后获取裁剪数据
              const originalDisplayUrl = displayImageUrl.value;
              displayImageUrl.value = processedUrl;
              
              nextTick(() => {
                // 重新初始化cropper以使用处理后的原图
                if (imageRef.value) {
                  const shouldRestore = false; // 不恢复裁剪框，使用新的图片
                  initCropper(shouldRestore);
                  
                  // 等待cropper初始化后获取裁剪数据
                  setTimeout(() => {
                    if (cropper.value) {
                      try {
                        const base64 = cropper.value.getCroppedCanvas().toDataURL('image/png');
                        emit('cropped', base64);
                        dialogVisible.value = false;
                        ElMessage.success('图片编辑成功！');
                        // 恢复显示URL
                        displayImageUrl.value = originalDisplayUrl;
                      } catch (error) {
                        console.error('裁剪错误:', error);
                        ElMessage.error('裁剪失败');
                        displayImageUrl.value = originalDisplayUrl;
                      }
                    }
                  }, 100);
                }
              });
              
              resolve();
            };
            processedImg.onerror = () => {
              reject(new Error('处理后的图片加载失败'));
            };
            processedImg.src = processedUrl;
          });
        });
      });
    } catch (error) {
      console.error('原图处理错误:', error);
      ElMessage.error('图片处理失败');
    }
  };
  
  /**
   * 对话框关闭后的清理工作
   */
  const handleDialogClosed = () => {
    // 清除防抖定时器
    if (debounceTimer) {
      clearTimeout(debounceTimer);
      debounceTimer = null;
    }
    // 清理 requestAnimationFrame
    if (rafId !== null) {
      cancelAnimationFrame(rafId);
      rafId = null;
    }
    // 重置处理状态
    isProcessing.value = false;
    // 清理 CamanJS 实例
    if (camanInstance) {
      try {
        camanInstance.revert();
      } catch (e) {
        // 忽略清理错误
      }
      camanInstance = null;
    }
    // 销毁 Cropper 实例，释放资源
    if (cropper.value) {
      cropper.value.destroy();
      cropper.value = null;
    }
  };
  
  // 暴露一些方法给父组件调用（如果需要）
  defineExpose({
    getCroppedDataURL: () => {
      if (cropper.value) {
        return cropper.value.getCroppedCanvas().toDataURL('image/png');
      }
      return null;
    }
  });
  </script>
  
  <style scoped>
  .editor-layout {
    display: flex;
    gap: 20px;
  }

  .editor-layout.mobile-layout {
    flex-direction: column;
    gap: 12px;
  }

  .cropper-section {
    flex: 1;
    min-width: 0;
  }

  .cropper-container {
    /* 必须设置高度，否则图片容器不会显示 */
    max-width: 100%;
    /* 确保 img 元素在容器内 */
    overflow: hidden;
    /* 裁剪器需要 */
    position: relative;
    background: #f5f5f5;
    border-radius: 4px;
  }
  
  /* cropperjs 的样式在全局引入了，这里只需要保证 img 元素是可用的 */
  .cropper-container img {
    /* cropperjs 库的要求 */
    display: block;
    max-width: 100%;
  }

  .adjustment-panel {
    width: 250px;
    padding: 20px;
    background: #fafafa;
    border-radius: 4px;
    flex-shrink: 0;
  }

  .mobile-layout .adjustment-panel {
    width: 100%;
    padding: 12px;
    background: #fff;
  }

  .panel-title {
    font-size: 16px;
    font-weight: bold;
    margin-bottom: 20px;
    color: #333;
  }

  .mobile-layout .panel-title {
    font-size: 14px;
    margin-bottom: 12px;
  }

  .adjustment-item {
    margin-bottom: 24px;
  }

  .mobile-layout .adjustment-item {
    margin-bottom: 14px;
  }

  .adjustment-label {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 8px;
    font-size: 14px;
    color: #606266;
  }

  .mobile-layout .adjustment-label {
    font-size: 12px;
    margin-bottom: 6px;
  }

  .adjustment-value {
    font-weight: bold;
    color: #409eff;
    min-width: 50px;
    text-align: right;
    font-size: 14px;
  }

  .mobile-layout .adjustment-value {
    font-size: 13px;
    min-width: 45px;
  }

  .reset-button {
    width: 100%;
    margin-top: 10px;
  }

  .mobile-layout .reset-button {
    margin-top: 8px;
    padding: 8px 16px;
    font-size: 13px;
  }

  .dialog-footer {
    display: flex;
    justify-content: flex-end;
    gap: 10px;
  }

  .dialog-footer.mobile-footer {
    flex-direction: row;
    justify-content: space-between;
    gap: 10px;
  }

  .dialog-footer.mobile-footer .el-button {
    flex: 1;
    margin: 0;
    height: 44px;
    font-size: 16px;
    max-width: 45%;
  }

  /* 移动端 slider 优化 */
  .mobile-layout :deep(.el-slider) {
    margin: 8px 0;
  }

  .mobile-layout :deep(.el-slider__runway) {
    height: 6px;
    background-color: #e4e7ed;
  }

  .mobile-layout :deep(.el-slider__button) {
    width: 20px;
    height: 20px;
    border: 2px solid #409eff;
    background-color: #fff;
  }

  .mobile-layout :deep(.el-slider__button-wrapper) {
    width: 36px;
    height: 36px;
    top: -18px;
    cursor: pointer;
    -webkit-tap-highlight-color: transparent;
  }

  .mobile-layout :deep(.el-slider__button-wrapper:hover) {
    transform: scale(1.1);
  }

  /* 对话框移动端优化 */
  :deep(.image-editor-dialog .el-dialog) {
    margin: 5vh auto;
  }

  /* 移动端全屏对话框 */
  :deep(.image-editor-dialog.mobile-dialog .el-dialog) {
    width: 100% !important;
    height: 100vh !important;
    max-height: 100vh !important;
    margin: 0 !important;
    border-radius: 0 !important;
    display: flex;
    flex-direction: column;
  }

  :deep(.image-editor-dialog.mobile-dialog .el-dialog__header) {
    padding: 10px 12px;
    border-bottom: 1px solid #ebeef5;
    flex-shrink: 0;
    position: relative;
  }

  :deep(.image-editor-dialog.mobile-dialog .el-dialog__title) {
    font-size: 17px;
    font-weight: 600;
  }

  :deep(.image-editor-dialog.mobile-dialog .el-dialog__headerbtn) {
    top: 10px;
    right: 12px;
    width: 32px;
    height: 32px;
  }

  :deep(.image-editor-dialog.mobile-dialog .el-dialog__close) {
    font-size: 20px;
  }

  :deep(.image-editor-dialog.mobile-dialog .el-dialog__body) {
    padding: 12px;
    flex: 1;
    overflow-y: auto;
    -webkit-overflow-scrolling: touch;
    overscroll-behavior: contain;
  }

  :deep(.image-editor-dialog.mobile-dialog .el-dialog__footer) {
    padding: 10px 12px;
    flex-shrink: 0;
    border-top: 1px solid #ebeef5;
    background-color: #fff;
  }

  /* 移动端裁剪区域优化 */
  .mobile-layout .cropper-container {
    border-radius: 8px;
    touch-action: none;
  }

  /* 移动端 cropperjs 触摸优化 */
  .mobile-layout :deep(.cropper-container .cropper-view-box) {
    outline: 2px solid #409eff;
    outline-offset: -2px;
  }

  .mobile-layout :deep(.cropper-container .cropper-face) {
    background-color: rgba(64, 158, 255, 0.1);
  }

  .mobile-layout :deep(.cropper-container .cropper-line) {
    background-color: #409eff;
  }

  .mobile-layout :deep(.cropper-container .cropper-point) {
    width: 10px;
    height: 10px;
    background-color: #409eff;
    border: 2px solid #fff;
  }

  /* 小屏幕优化 */
  @media (max-width: 480px) {
    :deep(.image-editor-dialog.mobile-dialog .el-dialog__body) {
      padding: 12px;
    }

    .mobile-layout .adjustment-panel {
      padding: 12px;
    }

    .mobile-layout .cropper-container {
      border-radius: 4px;
    }
  }

  /* 横屏优化 */
  @media (max-width: 768px) and (orientation: landscape) {
    .mobile-layout .cropper-container {
      height: 200px !important;
    }
  }
  </style>