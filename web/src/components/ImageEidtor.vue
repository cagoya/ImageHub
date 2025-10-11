<template>
    <el-dialog
      v-model="dialogVisible"
      title="图片裁剪"
      width="600px"
      :close-on-click-modal="false"
      @closed="handleDialogClosed"
    >
      <div class="cropper-container" :style="{ height: containerHeight + 'px' }">
        <img ref="imageRef" :src="imageUrl" alt="待裁剪图片" />
      </div>
  
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" :loading="isCropping" @click="handleCrop">
            确定裁剪
          </el-button>
        </span>
      </template>
    </el-dialog>
  </template>
  
  <script setup lang="ts">
  import { ref, watch, onMounted, nextTick } from 'vue';
  import 'cropperjs/dist/cropper.css';
  import Cropper from 'cropperjs';
  import { ElMessage } from 'element-plus';
  
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
  const cropper = ref<Cropper | null>(null);
  const isCropping = ref(false);
  
  // 双向绑定 visible
  watch(
    () => props.visible,
    (val) => {
      dialogVisible.value = val;
      if (val && props.imageUrl) {
        console.log(props.imageUrl);
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
      ElMessage.success('图片裁剪成功！');
    } catch (error) {
      ElMessage.error('图片裁剪失败！');
      console.error('裁剪错误:', error);
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
  .cropper-container {
    /* 必须设置高度，否则图片容器不会显示 */
    max-width: 100%;
    /* 确保 img 元素在容器内 */
    overflow: hidden;
    /* 裁剪器需要 */
    position: relative;
  }
  
  /* cropperjs 的样式在全局引入了，这里只需要保证 img 元素是可用的 */
  .cropper-container img {
    /* cropperjs 库的要求 */
    display: block;
    max-width: 100%;
  }
  </style>