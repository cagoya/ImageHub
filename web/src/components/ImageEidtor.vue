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
      <div v-if="dialogVisible" class="editor-layout" :class="{ 'mobile-layout': isMobile }">
        <div class="cropper-section">
          <div class="cropper-container" :style="{ height: isMobile ? cropperHeight + 'px' : containerHeight + 'px' }">
            <canvas ref="canvasRef" style="max-width: 100%; display: none;"></canvas>
            <img ref="imageRef" :src="displayImageUrl" alt="待编辑图片" />
            <div v-if="isProcessing || isCropping" class="loading-overlay">
                <el-icon class="is-loading" size="2em"><Loading /></el-icon>
                <div class="loading-text">{{ isCropping ? '正在处理高分辨率图片...' : '正在应用色调调整...' }}</div>
            </div>
          </div>
        </div>

        <div class="adjustment-panel">
          <div class="panel-title">色调调整</div>
          
          <div class="adjustment-item">
            <div class="adjustment-label">
              <span>亮度</span>
              <span class="adjustment-value">{{ formatValue(adjustments.brightness) }}</span>
            </div>
            <el-slider
              v-model="adjustments.brightness"
              :min="-50"
              :max="50"
              :step="0.1"
              :show-tooltip="false"
              @input="applyAdjustments"
              :disabled="isProcessing || isCropping"
            />
          </div>

          <div class="adjustment-item">
            <div class="adjustment-label">
              <span>对比度</span>
              <span class="adjustment-value">{{ formatValue(adjustments.contrast) }}</span>
            </div>
            <el-slider
              v-model="adjustments.contrast"
              :min="-50"
              :max="50"
              :step="0.1"
              :show-tooltip="false"
              @input="applyAdjustments"
              :disabled="isProcessing || isCropping"
            />
          </div>

          <div class="adjustment-item">
            <div class="adjustment-label">
              <span>饱和度</span>
              <span class="adjustment-value">{{ formatValue(adjustments.saturation) }}</span>
            </div>
            <el-slider
              v-model="adjustments.saturation"
              :min="-50"
              :max="50"
              :step="0.1"
              :show-tooltip="false"
              @input="applyAdjustments"
              :disabled="isProcessing || isCropping"
            />
          </div>

          <div class="adjustment-item">
            <div class="adjustment-label">
              <span>色相</span>
              <span class="adjustment-value">{{ formatValue(adjustments.hue) }}°</span>
            </div>
            <el-slider
              v-model="adjustments.hue"
              :min="-180"
              :max="180"
              :step="0.1"
              :show-tooltip="false"
              @input="applyAdjustments"
              :disabled="isProcessing || isCropping"
            />
          </div>

          <el-button 
            type="default" 
            @click="resetAdjustments" 
            class="reset-button"
            :disabled="isProcessing || isCropping"
          >
            重置调整
          </el-button>
        </div>
      </div>
  
      <template #footer>
        <div class="dialog-footer" :class="{ 'mobile-footer': isMobile }">
          <el-button 
            @click="dialogVisible = false" 
            :size="isMobile ? 'large' : 'default'"
            :disabled="isProcessing || isCropping"
          >
            {{ isMobile ? '取消' : '取消' }}
          </el-button>
          <el-button 
            type="primary" 
            :loading="isCropping" 
            :disabled="isProcessing"
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
import { Loading } from '@element-plus/icons-vue'; 
import { rgbToHsl, hslToRgb } from '@/utils/color'; 

// 声明 CamanJS 全局类型
declare global {
  interface Window {
    Caman: any;
  }
}
  
  // --- Props & Emits ---
  interface Props {
    imageUrl: string;
    containerHeight?: number;
    visible: boolean;
    aspectRatio?: number;
  }
  
  const props = withDefaults(defineProps<Props>(), {
    containerHeight: 400,
    aspectRatio: NaN,
  });
  
  interface Emits {
    (e: 'update:visible', value: boolean): void;
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
  
  let camanInstance: any = null;
  let debounceTimer: ReturnType<typeof setTimeout> | null = null;
  let cachedOriginalImage: HTMLImageElement | null = null;
  let cachedImageDimensions: { width: number; height: number } | null = null;
  
  const PREVIEW_MAX_SIZE_PC = 1920; 
  const PREVIEW_MAX_SIZE_MOBILE = 800;
  const isProcessing = ref(false); 
  let rafId: number | null = null;
  
  const isMobile = computed(() => {
    if (typeof window === 'undefined') return false;
    return window.innerWidth <= 768;
  });
  
  const cropperHeight = ref(250);
  
  // --- 生命周期 & 响应式 ---
  const calculateCropperHeight = () => {
    if (typeof window === 'undefined') return;
    const viewportHeight = window.innerHeight;
    const reservedHeight = 514; 
    const availableHeight = viewportHeight - reservedHeight;
    cropperHeight.value = Math.max(180, Math.min(260, availableHeight));
  };
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
    // 初始加载和缓存
    loadAndCacheOriginalImage().catch(e => console.error(e));
  });
  onUnmounted(() => {
    window.removeEventListener('resize', handleResize);
    // 组件卸载时，确保所有定时器和实例被清理
    if (debounceTimer) {
      clearTimeout(debounceTimer);
    }
    if (rafId !== null) {
      cancelAnimationFrame(rafId);
    }
    if (camanInstance) {
      try { camanInstance.revert(); } catch (e) {}
      camanInstance = null;
    }
    // 理论上由于 v-if，cropper 也会被销毁，这里只是保险
    if (cropper.value) {
      cropper.value.destroy();
      cropper.value = null;
    }
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

  // 双向绑定 visible
  watch(
    () => props.visible,
    (val) => {
      dialogVisible.value = val;
      if (val && props.imageUrl) {
        // 由于使用了 v-if，DOM 会被重建，不需要手动销毁 cropper，只需重置状态
        savedCropData.value = null;
        currentDisplayImageSize = null;
        originalImageUrl.value = props.imageUrl;
        displayImageUrl.value = props.imageUrl; // 重置为新的图片 URL
        cachedOriginalImage = null;
        cachedImageDimensions = null;
        
        resetAdjustmentsState(); 
        if (camanInstance) {
          try { camanInstance.revert(); } catch (e) {}
          camanInstance = null;
        }
        
        // 预加载、缓存并初始化 Cropper
        loadAndCacheOriginalImage()
          .then(() => {
             nextTick(() => {
                // DOM 已经重建，现在可以安全地初始化 Cropper
                initCropper();
             });
          })
          .catch(() => {
             ElMessage.error('原始图片加载失败');
          });
      }
    }
  );
  
  watch(dialogVisible, (val) => {
    emit('update:visible', val);
  });

  // --- 裁剪状态管理 ---
  const savedCropData = ref<{
    x: number;
    y: number;
    width: number;
    height: number;
  } | null>(null);
  
  // 保存当前显示图片的尺寸，用于坐标缩放
  let currentDisplayImageSize: { width: number; height: number } | null = null;

  /**
   * 初始化 Cropper 实例
   * **已优化：确保销毁前保存数据，并处理图片尺寸变化**
   */
  const initCropper = (restoreCropData = false) => {
    if (!imageRef.value) return;
  
    // 1. **核心优化：** 无论是否恢复，只要实例存在，就尝试保存当前裁剪框数据
    if (cropper.value) {
      try {
        const cropData = cropper.value.getData();
        // 覆盖保存，确保是最新状态
        savedCropData.value = { x: cropData.x, y: cropData.y, width: cropData.width, height: cropData.height };
        // 保存当前图片的显示尺寸（使用 Cropper 的 API）
        const imageData = cropper.value.getImageData();
        currentDisplayImageSize = {
          width: imageData.naturalWidth,
          height: imageData.naturalHeight
        };
      } catch (e) {
        console.warn('Failed to get crop data before destroy:', e);
        // 如果获取失败，则不进行恢复，使用默认裁剪框
        savedCropData.value = null;
        currentDisplayImageSize = null;
      }
      
      // 销毁现有实例
      cropper.value.destroy();
      cropper.value = null;
    }
  
    // 2. 初始化 Cropper
    cropper.value = new Cropper(imageRef.value, {
      aspectRatio: props.aspectRatio,
      viewMode: 1,
      dragMode: 'none', // 禁止移动图片，只允许调整裁剪框
      zoomable: false, // 禁止缩放
      zoomOnWheel: false, // 禁止鼠标滚轮缩放
      rotatable: false,
      responsive: true,
      restore: false,
      autoCropArea: 0.8,
      minCropBoxWidth: isMobile.value ? 50 : 0,
      minCropBoxHeight: isMobile.value ? 50 : 0,
      ready: () => {
        // 3. 恢复裁剪框位置（处理图片尺寸变化）
        if (restoreCropData && savedCropData.value && cropper.value) {
          try {
            let cropDataToRestore = savedCropData.value;
            
            // 如果图片尺寸发生了变化，需要按比例缩放裁剪框坐标
            if (currentDisplayImageSize) {
              const imageData = cropper.value.getImageData();
              const newImageSize = {
                width: imageData.naturalWidth,
                height: imageData.naturalHeight
              };
              
              // 计算缩放比例
              const scaleX = newImageSize.width / currentDisplayImageSize.width;
              const scaleY = newImageSize.height / currentDisplayImageSize.height;
              
              // 如果尺寸变化较大（超过5%），则缩放坐标
              if (Math.abs(scaleX - 1) > 0.05 || Math.abs(scaleY - 1) > 0.05) {
                cropDataToRestore = {
                  x: cropDataToRestore.x * scaleX,
                  y: cropDataToRestore.y * scaleY,
                  width: cropDataToRestore.width * scaleX,
                  height: cropDataToRestore.height * scaleY
                };
              }
              
              // 更新当前图片尺寸记录
              currentDisplayImageSize = newImageSize;
            } else if (cropper.value) {
              // 如果没有保存的尺寸，记录当前图片尺寸
              const imageData = cropper.value.getImageData();
              currentDisplayImageSize = {
                width: imageData.naturalWidth,
                height: imageData.naturalHeight
              };
            }
            
            cropper.value.setData(cropDataToRestore);
            
            // 恢复成功后，清空保存的数据，避免下次非恢复调用时误用
            savedCropData.value = null; 
          } catch (e) {
            console.warn('Failed to restore crop data:', e);
            savedCropData.value = null;
          }
        } else if (!restoreCropData && cropper.value) {
          // 非恢复模式下，记录当前图片尺寸
          const imageData = cropper.value.getImageData();
          currentDisplayImageSize = {
            width: imageData.naturalWidth,
            height: imageData.naturalHeight
          };
        }
      },
    });
  };
  
  // --- 色调调整状态管理 ---
  const adjustments = ref({
    brightness: 0,
    contrast: 0,
    saturation: 0,
    hue: 0,
  });
  
  const lastAdjustments = ref({
    brightness: 0,
    contrast: 0,
    saturation: 0,
    hue: 0,
  });

  /**
   * 格式化数值显示
   */
  const formatValue = (value: number): string => {
    const rounded = Math.round(value * 10) / 10;
    if (rounded % 1 === 0) {
      return rounded.toString();
    }
    return rounded.toFixed(1);
  };
  
  /**
   * 手动应用色相调整（基于像素）
   */
  const applyHueAdjustment = (canvas: HTMLCanvasElement, hueRotation: number) => {
    const ctx = canvas.getContext('2d');
    if (!ctx) return;

    const imageData = ctx.getImageData(0, 0, canvas.width, canvas.height);
    const data = imageData.data;
    
    for (let i = 0; i < data.length; i += 4) {
        const r = data[i] ?? 0;
        const g = data[i + 1] ?? 0;
        const b = data[i + 2] ?? 0;

        const hsl = rgbToHsl(r, g, b);
        let newHue = (hsl[0] + hueRotation) % 360;
        if (newHue < 0) newHue += 360;
        
        const rgb = hslToRgb(newHue, hsl[1], hsl[2]);
        
        data[i] = Math.max(0, Math.min(255, rgb[0]));
        data[i + 1] = Math.max(0, Math.min(255, rgb[1]));
        data[i + 2] = Math.max(0, Math.min(255, rgb[2]));
    }

    ctx.putImageData(imageData, 0, 0);
  };

  /**
   * 计算预览尺寸
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
      const currentOriginalUrl = props.imageUrl || originalImageUrl.value;
      
      if (cachedOriginalImage && cachedOriginalImage.complete && cachedOriginalImage.src === currentOriginalUrl) {
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
      
      img.src = currentOriginalUrl;
    });
  };

  /**
   * 使用 CamanJS 应用色调调整（预览处理）
   */
  const doApplyAdjustments = async () => {
    if (!imageRef.value || typeof window === 'undefined' || !window.Caman || isProcessing.value) {
      return;
    }

    isProcessing.value = true;
    
    // **关键修复：在更新图片之前，先保存当前裁剪框数据和图片尺寸**
    if (cropper.value) {
      try {
        const cropData = cropper.value.getData();
        savedCropData.value = { 
          x: cropData.x, 
          y: cropData.y, 
          width: cropData.width, 
          height: cropData.height 
        };
        // 保存当前图片的显示尺寸（使用 Cropper 的 API）
        const imageData = cropper.value.getImageData();
        currentDisplayImageSize = {
          width: imageData.naturalWidth,
          height: imageData.naturalHeight
        };
      } catch (e) {
        console.warn('Failed to save crop data before adjustment:', e);
        savedCropData.value = null;
        currentDisplayImageSize = null;
      }
    }
    
    // 如果所有调整都是0，直接使用原始图片
    if (
      adjustments.value.brightness === 0 &&
      adjustments.value.contrast === 0 &&
      adjustments.value.saturation === 0 &&
      adjustments.value.hue === 0
    ) {
        displayImageUrl.value = props.imageUrl || originalImageUrl.value;
        isProcessing.value = false;
        nextTick(() => {
            initCropper(true);
        });
        return;
    }

    if (rafId !== null) {
      cancelAnimationFrame(rafId);
    }
    rafId = requestAnimationFrame(async () => {
        try {
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

            const previewSize = calculatePreviewSize(
                cachedImageDimensions.width,
                cachedImageDimensions.height
            );

            canvas.width = previewSize.width;
            canvas.height = previewSize.height;

            ctx.imageSmoothingEnabled = true;
            ctx.imageSmoothingQuality = isMobile.value ? 'medium' : 'high';
            ctx.drawImage(originalImg, 0, 0, previewSize.width, previewSize.height);

            if (camanInstance) {
              camanInstance = null; 
            }

            window.Caman(canvas, function(this: any) {
                this.revert(false); 
                
                if (adjustments.value.brightness !== 0) this.brightness(adjustments.value.brightness);
                if (adjustments.value.contrast !== 0) this.contrast(adjustments.value.contrast);
                if (adjustments.value.saturation !== 0) this.saturation(adjustments.value.saturation);
                
                const needsHueAdjustment = adjustments.value.hue !== 0;
                
                this.render(() => {
                    if (needsHueAdjustment) {
                        applyHueAdjustment(canvas, adjustments.value.hue);
                    }
                    
                    const imageFormat = isMobile.value ? 'image/jpeg' : 'image/png';
                    const imageQuality = isMobile.value ? 0.85 : undefined;
                    displayImageUrl.value = canvas.toDataURL(imageFormat, imageQuality);
                    
                    camanInstance = this;
                    
                    isProcessing.value = false;
                    
                    nextTick(() => {
                        initCropper(true);
                    });
                });
            });
        } catch (error) {
            console.error('CamanJS processing error:', error);
            ElMessage.error('图片预览处理失败');
            isProcessing.value = false;
        } finally {
            rafId = null;
        }
    });
  };

  /**
   * 计算动态防抖时间
   */
  const getDebounceTime = (): number => {
    if (!cachedImageDimensions) {
      return isMobile.value ? 300 : 150;
    }
    const pixelCount = cachedImageDimensions.width * cachedImageDimensions.height;
    
    if (isMobile.value) {
      return pixelCount > 2000000 ? 400 : (pixelCount > 1000000 ? 300 : 250);
    } else {
      return pixelCount > 2000000 ? 200 : (pixelCount > 1000000 ? 150 : 100);
    }
  };

  /**
   * 应用色调调整（防抖）
   */
  const applyAdjustments = () => {
    const current = adjustments.value;
    const last = lastAdjustments.value;
    
    const hasChanged = 
      current.brightness !== last.brightness ||
      current.contrast !== last.contrast ||
      current.saturation !== last.saturation ||
      current.hue !== last.hue;
    
    if (!hasChanged) {
      return;
    }
    
    lastAdjustments.value = { ...current };
    
    if (debounceTimer) {
      clearTimeout(debounceTimer);
    }
    
    const debounceTime = getDebounceTime();
    
    debounceTimer = setTimeout(() => {
      doApplyAdjustments();
      debounceTimer = null;
    }, debounceTime);
  };
  
  /**
   * 仅重置参数状态（用于内部调用）
   */
  const resetAdjustmentsState = () => {
      adjustments.value = {
        brightness: 0,
        contrast: 0,
        saturation: 0,
        hue: 0,
      };
      lastAdjustments.value = {
        brightness: 0,
        contrast: 0,
        saturation: 0,
        hue: 0,
      };
  }

  /**
   * 重置色调调整（UI 调用）
   */
  const resetAdjustments = () => {
    resetAdjustmentsState();
    
    displayImageUrl.value = props.imageUrl || originalImageUrl.value;
    
    if (camanInstance) {
      try { camanInstance.revert(); } catch (e) {}
      camanInstance = null;
    }
    
    nextTick(() => {
      // 保持裁剪框位置并重新初始化
      initCropper(true);
    });
  };

  /**
   * 处理确定裁剪按钮点击（原图高分辨率处理和裁剪）
   */
  const handleCrop = async () => {
    if (!cropper.value || !cachedOriginalImage || !cachedImageDimensions) {
      ElMessage.error('裁剪器或图片数据未初始化！');
      return;
    }
  
    isCropping.value = true;
    
    try {
        // 获取裁剪数据（相对于当前显示图片的坐标）
        const cropData = cropper.value.getData();
        
        // 获取当前显示图片的自然尺寸（可能是预览尺寸，也可能是原图尺寸）
        const imageData = cropper.value.getImageData();
        const displayImageNaturalWidth = imageData.naturalWidth;
        const displayImageNaturalHeight = imageData.naturalHeight;
        
        // 获取原图尺寸
        const originalWidth = cachedImageDimensions.width;
        const originalHeight = cachedImageDimensions.height;
        
        // 计算从显示图片尺寸到原图尺寸的缩放比例
        // 如果显示图片是预览图（经过缩放的），需要将坐标缩放到原图尺寸
        const scaleX = originalWidth / displayImageNaturalWidth;
        const scaleY = originalHeight / displayImageNaturalHeight;
        
        // 将裁剪坐标从显示图片尺寸缩放到原图尺寸
        const scaledCropData = {
            x: cropData.x * scaleX,
            y: cropData.y * scaleY,
            width: cropData.width * scaleX,
            height: cropData.height * scaleY
        };
        
        const hasAdjustments = 
            adjustments.value.brightness !== 0 ||
            adjustments.value.contrast !== 0 ||
            adjustments.value.saturation !== 0 ||
            adjustments.value.hue !== 0;

        let processedSourceCanvas: HTMLCanvasElement;

        if (hasAdjustments) {
            processedSourceCanvas = document.createElement('canvas');
            processedSourceCanvas.width = cachedImageDimensions.width;
            processedSourceCanvas.height = cachedImageDimensions.height;
            const tempCtx = processedSourceCanvas.getContext('2d');
            if (!tempCtx) throw new Error('无法创建画布');
            tempCtx.drawImage(cachedOriginalImage, 0, 0);

            await new Promise<void>((resolve) => {
                window.Caman(processedSourceCanvas, function(this: any) {
                    this.revert(false);
                    
                    if (adjustments.value.brightness !== 0) this.brightness(adjustments.value.brightness);
                    if (adjustments.value.contrast !== 0) this.contrast(adjustments.value.contrast);
                    if (adjustments.value.saturation !== 0) this.saturation(adjustments.value.saturation);

                    this.render(() => {
                        if (adjustments.value.hue !== 0) {
                            applyHueAdjustment(processedSourceCanvas, adjustments.value.hue);
                        }
                        resolve();
                    });
                });
            });
        } else {
            processedSourceCanvas = document.createElement('canvas');
            processedSourceCanvas.width = cachedImageDimensions.width;
            processedSourceCanvas.height = cachedImageDimensions.height;
            const tempCtx = processedSourceCanvas.getContext('2d');
            if (tempCtx) {
                tempCtx.drawImage(cachedOriginalImage, 0, 0);
            } else {
                throw new Error('无法创建画布');
            }
        }

        // 使用缩放后的坐标创建裁剪画布
        const croppedCanvas = document.createElement('canvas');
        croppedCanvas.width = Math.round(scaledCropData.width);
        croppedCanvas.height = Math.round(scaledCropData.height);
        const croppedCtx = croppedCanvas.getContext('2d');

        if (croppedCtx) {
            croppedCtx.drawImage(
                processedSourceCanvas,
                Math.round(scaledCropData.x), Math.round(scaledCropData.y),
                Math.round(scaledCropData.width), Math.round(scaledCropData.height),
                0, 0,
                Math.round(scaledCropData.width), Math.round(scaledCropData.height)
            );
            
            const base64 = croppedCanvas.toDataURL('image/png');
            emit('cropped', base64);
        } else {
            throw new Error('无法创建裁剪画布');
        }

        clearPreviewImage(false); 
        dialogVisible.value = false;
        ElMessage.success('图片编辑成功！');
    } catch (error) {
        ElMessage.error('图片编辑失败！');
        console.error('编辑错误:', error);
    } finally {
        isCropping.value = false;
    }
  };
  
  /**
   * 清理预览图（释放base64数据URL的内存）
   */
  const clearPreviewImage = (resetToOriginal = true) => {
    // 此时 Cropper 实例的销毁依赖于 handleDialogClosed 或 v-if
    
    if (displayImageUrl.value && displayImageUrl.value.startsWith('data:')) {
      if (canvasRef.value) {
        const ctx = canvasRef.value.getContext('2d');
        if (ctx) {
            ctx.clearRect(0, 0, canvasRef.value.width, canvasRef.value.height);
        }
      }
      displayImageUrl.value = '';
    }
    
    if (resetToOriginal) {
      displayImageUrl.value = props.imageUrl || originalImageUrl.value;
    }
  };

  /**
   * 对话框关闭后的清理工作
   */
  const handleDialogClosed = () => {
    // 1. 清理状态和定时器
    if (debounceTimer) clearTimeout(debounceTimer);
    if (rafId !== null) cancelAnimationFrame(rafId);
    
    isProcessing.value = false;
    isCropping.value = false;

    if (camanInstance) {
      try { camanInstance.revert(); } catch (e) {}
      camanInstance = null;
    }
    
    // 2. 清理 Cropper 实例 (双重保险，v-if 已经移除了 DOM)
    if (cropper.value) {
      cropper.value.destroy();
      cropper.value = null;
    }

    // 3. 清理预览图状态 
    clearPreviewImage(true); 
    
    // 4. 重置调整参数状态
    resetAdjustmentsState();
    
    // 5. 清理图片尺寸记录
    currentDisplayImageSize = null;
    savedCropData.value = null;
  };
  
  // 暴露方法给父组件
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
  /* 样式保持不变 */
  .editor-layout {
    display: flex;
    gap: 20px;
    min-height: 400px;
  }

  .editor-layout.mobile-layout {
    flex-direction: column;
    gap: 12px;
    min-height: auto;
  }

  .cropper-section {
    flex: 1;
    min-width: 0;
    position: relative;
  }

  .cropper-container {
    max-width: 100%;
    overflow: hidden;
    position: relative;
    background: #f5f5f5;
    border-radius: 4px;
  }
  
  .cropper-container img {
    display: block;
    max-width: 100%;
  }

  /* 新增：加载层样式 */
  .loading-overlay {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(255, 255, 255, 0.8);
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    z-index: 10; 
    color: #409eff;
    font-size: 14px;
    border-radius: 4px;
  }
  
  .loading-text {
      margin-top: 10px;
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
  }
  </style>