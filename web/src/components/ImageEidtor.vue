  <template>
    <el-dialog
      v-model="dialogVisible"
      title="图片编辑"
      :width="isMobile ? '95%' : '900px'"
      :close-on-click-modal="false"
      @closed="handleDialogClosed"
      class="image-editor-dialog"
    >
      <div class="editor-layout" :class="{ 'mobile-layout': isMobile }">
        <!-- 裁剪区域 -->
        <div class="cropper-section">
          <div class="cropper-container" :style="{ height: isMobile ? '300px' : containerHeight + 'px' }">
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
            style="width: 100%; margin-top: 10px;"
          >
            重置调整
          </el-button>
        </div>
      </div>
  
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" :loading="isCropping" @click="handleCrop">
            确定
          </el-button>
        </span>
      </template>
    </el-dialog>
  </template>
  
  <script setup lang="ts">
import { ref, watch, onMounted, nextTick, computed, onUnmounted } from 'vue';
import 'cropperjs/dist/cropper.css';
import Cropper from 'cropperjs';
import { ElMessage } from 'element-plus';
import { rgbToHsl as utilRgbToHsl, hslToRgb as utilHslToRgb } from '@/utils/color';
  
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
  const originalImage = ref<HTMLImageElement | null>(null);
  
  // 防抖定时器
  let debounceTimer: ReturnType<typeof setTimeout> | null = null;
  
  // 移动端检测
  const isMobile = computed(() => {
    if (typeof window === 'undefined') return false;
    return window.innerWidth <= 768;
  });
  
  // 窗口大小变化监听
  const handleResize = () => {
    // 触发响应式更新
  };
  
  onMounted(() => {
    window.addEventListener('resize', handleResize);
  });
  
  onUnmounted(() => {
    window.removeEventListener('resize', handleResize);
    // 清理防抖定时器
    if (debounceTimer) {
      clearTimeout(debounceTimer);
      debounceTimer = null;
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
        console.log(props.imageUrl);
        // 重置调整参数和图片URL
        displayImageUrl.value = props.imageUrl;
        resetAdjustments();
        // 对话框显示时，确保DOM更新完毕后初始化 Cropper
        nextTick(() => {
          loadOriginalImage();
        });
      }
    }
  );
  watch(dialogVisible, (val) => {
    emit('update:visible', val);
  });

  /**
   * 加载原始图片
   */
  const loadOriginalImage = () => {
    const img = new Image();
    img.crossOrigin = 'anonymous'; // 处理跨域问题
    img.onload = () => {
      originalImage.value = img;
      initCropper();
    };
    img.onerror = () => {
      ElMessage.error('图片加载失败');
    };
    img.src = props.imageUrl;
  };
  
  /**
   * 初始化 Cropper 实例
   */
  const initCropper = () => {
    if (!imageRef.value) {
      ElMessage.error('图片DOM元素未找到。');
      return;
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
      // ... 其他配置项
    });
  };
  
  /**
   * 应用色调调整（实际执行函数）
   */
  const doApplyAdjustments = () => {
    if (!originalImage.value || !canvasRef.value) return;

    const canvas = canvasRef.value;
    const ctx = canvas.getContext('2d');
    if (!ctx) return;

    // 设置 canvas 尺寸
    canvas.width = originalImage.value.width;
    canvas.height = originalImage.value.height;

    // 绘制原始图片
    ctx.drawImage(originalImage.value, 0, 0);

    // 获取图片数据
    const imageData = ctx.getImageData(0, 0, canvas.width, canvas.height);
    const data = imageData.data;

    // 计算调整系数
    const brightnessValue = adjustments.value.brightness * 2.55; // 转换到 -255 到 255
    const contrastFactor = (adjustments.value.contrast + 100) / 100;
    const saturationFactor = (adjustments.value.saturation + 100) / 100;
    const hueRotation = adjustments.value.hue;

    // 处理每个像素
    for (let i = 0; i < data.length; i += 4) {
      let r = data[i] ?? 0;
      let g = data[i + 1] ?? 0;
      let b = data[i + 2] ?? 0;

      // 应用亮度
      r += brightnessValue;
      g += brightnessValue;
      b += brightnessValue;

      // 应用对比度
      r = ((r / 255 - 0.5) * contrastFactor + 0.5) * 255;
      g = ((g / 255 - 0.5) * contrastFactor + 0.5) * 255;
      b = ((b / 255 - 0.5) * contrastFactor + 0.5) * 255;

      // 转换到 HSL 来应用饱和度和色相
      if (saturationFactor !== 1 || hueRotation !== 0) {
        const hsl = convertRgbToHsl(r, g, b);
        
        // 应用饱和度
        hsl[1] = Math.max(0, Math.min(1, hsl[1] * saturationFactor));
        
        // 应用色相旋转
        if (hueRotation !== 0) {
          hsl[0] = (hsl[0] + hueRotation / 360) % 1;
          if (hsl[0] < 0) hsl[0] += 1;
        }
        
        // 转换回 RGB
        const rgb = convertHslToRgb(hsl[0], hsl[1], hsl[2]);
        r = rgb[0];
        g = rgb[1];
        b = rgb[2];
      }

      // 限制在有效范围内
      data[i] = Math.max(0, Math.min(255, r));
      data[i + 1] = Math.max(0, Math.min(255, g));
      data[i + 2] = Math.max(0, Math.min(255, b));
    }

    // 将处理后的数据放回 canvas
    ctx.putImageData(imageData, 0, 0);

    // 更新显示的图片 URL
    displayImageUrl.value = canvas.toDataURL('image/png');

    // 重新初始化 cropper
    nextTick(() => {
      initCropper();
    });
  };

  /**
   * 应用色调调整（防抖版本，50ms延迟）
   */
  const applyAdjustments = () => {
    // 清除之前的定时器
    if (debounceTimer) {
      clearTimeout(debounceTimer);
    }
    
    // 设置新的定时器，50ms后执行
    debounceTimer = setTimeout(() => {
      doApplyAdjustments();
      debounceTimer = null;
    }, 50);
  };

  /**
   * RGB 转 HSL（使用工具函数）
   */
  const convertRgbToHsl = (r: number, g: number, b: number): [number, number, number] => {
    return utilRgbToHsl(r, g, b);
  };

  /**
   * HSL 转 RGB（使用工具函数）
   */
  const convertHslToRgb = (h: number, s: number, l: number): [number, number, number] => {
    return utilHslToRgb(h, s, l);
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
    displayImageUrl.value = props.imageUrl;
    nextTick(() => {
      if (originalImage.value) {
        initCropper();
      }
    });
  };

  /**
   * 处理确定裁剪按钮点击
   */
  const handleCrop = () => {
    if (!cropper.value) {
      ElMessage.error('裁剪器未初始化！');
      return;
    }
  
    isCropping.value = true;
    try {
      // 获取裁剪后的 Base64 数据
      const base64 = cropper.value.getCroppedCanvas().toDataURL('image/png');
  
      // 触发事件给父组件
      emit('cropped', base64);
      
      // 关闭对话框
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
   * 对话框关闭后的清理工作
   */
  const handleDialogClosed = () => {
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
    gap: 16px;
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
    padding: 16px;
  }

  .panel-title {
    font-size: 16px;
    font-weight: bold;
    margin-bottom: 20px;
    color: #333;
  }

  .mobile-layout .panel-title {
    font-size: 15px;
    margin-bottom: 16px;
  }

  .adjustment-item {
    margin-bottom: 24px;
  }

  .mobile-layout .adjustment-item {
    margin-bottom: 20px;
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
    font-size: 13px;
    margin-bottom: 10px;
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

  .dialog-footer {
    display: flex;
    justify-content: flex-end;
    gap: 10px;
  }

  /* 移动端 slider 优化 */
  .mobile-layout :deep(.el-slider) {
    margin: 8px 0;
  }

  .mobile-layout :deep(.el-slider__runway) {
    height: 4px;
  }

  .mobile-layout :deep(.el-slider__button) {
    width: 16px;
    height: 16px;
  }

  .mobile-layout :deep(.el-slider__button-wrapper) {
    width: 32px;
    height: 32px;
    top: -14px;
  }

  /* 对话框移动端优化 */
  :deep(.image-editor-dialog .el-dialog) {
    margin: 5vh auto;
  }

  @media (max-width: 768px) {
    :deep(.image-editor-dialog .el-dialog) {
      width: 95% !important;
      margin: 2vh auto;
      max-height: 90vh;
      display: flex;
      flex-direction: column;
    }

    :deep(.image-editor-dialog .el-dialog__body) {
      padding: 16px;
      flex: 1;
      overflow-y: auto;
    }

    :deep(.image-editor-dialog .el-dialog__footer) {
      padding: 12px 16px;
      flex-shrink: 0;
    }

    .dialog-footer {
      flex-direction: column;
    }

    .dialog-footer .el-button {
      width: 100%;
      margin: 0;
    }

    .dialog-footer .el-button + .el-button {
      margin-top: 10px;
    }
  }

  @media (max-width: 480px) {
    :deep(.image-editor-dialog .el-dialog) {
      width: 98% !important;
      margin: 1vh auto;
    }

    :deep(.image-editor-dialog .el-dialog__header) {
      padding: 12px 16px;
    }

    :deep(.image-editor-dialog .el-dialog__title) {
      font-size: 16px;
    }

    :deep(.image-editor-dialog .el-dialog__body) {
      padding: 12px;
    }

    .adjustment-panel {
      padding: 12px;
    }
  }
  </style>