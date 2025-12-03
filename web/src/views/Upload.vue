<template>
  <AppHeader />
  <div class="image-uploader-container">
    <div class="uploader-header">
      <h2>批量图片上传</h2>
      <p>支持上传多张图片，系统将自动处理您的图片</p>
    </div>

    <div class="uploader-card">
      <el-upload
        class="avatar-uploader"
        action="#"
        :limit="10"
        :show-file-list="true"
        :on-change="handleFileChange"
        :on-exceed="handleExceed"
        :on-remove="handleRemove"
        :on-preview="handlePictureCardPreview"
        :before-upload="beforeUpload"
        :http-request="customUpload"
        :file-list="fileList"
        list-type="picture-card"
        accept=".jpg, .jpeg, .png"
        multiple>
        <div class="upload-area">
          <el-icon class="upload-icon"><Plus /></el-icon>
          <div class="upload-text">点击上传图片</div>
          <div class="upload-subtext">或将图片拖拽到此处</div>
        </div>
        <template #tip>
          <div class="el-upload__tip">
            最多上传10张JPG/PNG图片，单张大小不能超过20MB
          </div>
        </template>
      </el-upload>

      <div class="button-wrapper">
        <el-button
          class="submit-btn"
          type="primary"
          :loading="isUploading"
          @click="submitUpload"
          :disabled="fileList.length === 0 || isUploading"
        >
          <el-icon v-if="!isUploading" class="btn-icon"><Upload /></el-icon>
          <span class="btn-text">{{ isUploading ? '上传中...' : '提交上传' }}</span>
          <span v-if="fileList.length > 0" class="file-count">({{ fileList.length }})</span>
        </el-button>
      </div>

      <div v-if="uploadSummary.successCount + uploadSummary.failCount > 0" class="upload-summary">
        <div class="summary-content">
          <el-icon class="success-icon"><SuccessFilled /></el-icon>
          <span class="summary-text">
            上传完成：
            <span class="success-count">成功 {{ uploadSummary.successCount }} 张</span>，
            <span class="fail-count">失败 {{ uploadSummary.failCount }} 张</span>
          </span>
        </div>
      </div>
    </div>
  </div>
  
  <el-image-viewer
    v-if="imageViewerVisible"
    :url-list="[dialogImageUrl]"
    @close="closeImageViewer"
  />
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { ElMessage, type UploadProps, type UploadRawFile, type UploadFile, type UploadUserFile } from 'element-plus';
import { Plus, Upload, SuccessFilled } from '@element-plus/icons-vue';
import { ImageAPI } from '@/api';
import AppHeader from '@/components/AppHeader.vue';

const dialogImageUrl = ref('');
const imageViewerVisible = ref(false);

const fileList = ref<UploadUserFile[]>([]); 
const isUploading = ref(false);
const uploadSummary = ref({
  successCount: 0,
  failCount: 0,
});

// 文件状态改变时的钩子
const handleFileChange: UploadProps['onChange'] = (uploadFile, list) => {
  fileList.value = list;
};

const handlePictureCardPreview: UploadProps['onPreview'] = (uploadFile: UploadFile) => {
  if (uploadFile.url) {
    dialogImageUrl.value = uploadFile.url;
    imageViewerVisible.value = true;
  }
};

// 关闭图片预览
const closeImageViewer = () => {
  imageViewerVisible.value = false;
};

const handleRemove: UploadProps['onRemove'] = (uploadFile: UploadFile, uploadFiles: UploadUserFile[]) => {
  fileList.value = uploadFiles;
  if (uploadFile.url && uploadFile.raw) {
    URL.revokeObjectURL(uploadFile.url);
  }
};

// 超过文件数量限制时的钩子
const handleExceed: UploadProps['onExceed'] = () => {
  ElMessage.warning('最多只能选择 10 张图片进行上传！');
};

// 上传前的校验钩子 (格式和大小)
const beforeUpload: UploadProps['beforeUpload'] = (rawFile) => {
  const isTypeValid = rawFile.type.startsWith('image/'); 
  const isLt20M = rawFile.size / 1024 / 1024 < 20; 

  if (!isTypeValid) {
    ElMessage.error('文件格式错误，只能上传图片格式!');
    return false;
  }
  if (!isLt20M) {
    ElMessage.error('图片大小不能超过 20MB!');
    return false;
  }
  return true;
};

// 覆盖默认上传行为的钩子
const customUpload: UploadProps['httpRequest'] = async () => {
  return;
};

// 批量提交
const submitUpload = async () => {
  const filesToUpload = fileList.value.map(f => f.raw as UploadRawFile);

  if (filesToUpload.length === 0) {
    ElMessage.warning('没有需要上传的文件');
    return;
  }

  isUploading.value = true;
  uploadSummary.value = { successCount: 0, failCount: 0 };
  
  const uploadPromises = filesToUpload.map((file) => 
    ImageAPI.upload(file).then(response => ({ file, response }))
  );

  try {
    const responses = await Promise.allSettled(uploadPromises);
    const filesToRemove: UploadRawFile[] = [];

    responses.forEach(item => {
      if (item.status === 'fulfilled') {
        const { file, response } = item.value;
        if (response.data.code === 200) {
          uploadSummary.value.successCount++;
          filesToRemove.push(file);
        } else {
          uploadSummary.value.failCount++;
        }
      } else {
        uploadSummary.value.failCount++;
      }
    });

    filesToRemove.forEach(file => {
      const fileInList = fileList.value.find(f => f.raw === file);
      if (fileInList?.url) {
        URL.revokeObjectURL(fileInList.url);
      }
    });
    
    fileList.value = fileList.value.filter(f => !filesToRemove.includes(f.raw as UploadRawFile));
    
    if (uploadSummary.value.failCount === 0) {
      ElMessage.success(`所有 ${uploadSummary.value.successCount} 张图片上传成功！`);
    } else if (uploadSummary.value.successCount > 0) {
      ElMessage.warning(`部分图片上传完成：成功 ${uploadSummary.value.successCount} 张，失败 ${uploadSummary.value.failCount} 张`);
    } else {
      ElMessage.error(`所有 ${uploadSummary.value.failCount} 张图片上传失败！`);
    }

  } catch (error) {
    ElMessage.error('图片上传失败，请检查网络或后端服务');
  } finally {
    isUploading.value = false;
  }
};
</script>

<style scoped>
.image-uploader-container {
  min-height: calc(100vh - 60px);
  background: linear-gradient(135deg, #f0f8ff 0%, #e6f3ff 50%, #dcedff 100%);
  padding: 30px 20px;
  display: flex;
  flex-direction: column;
  align-items: center;
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
}

.uploader-header {
  text-align: center;
  margin-bottom: 30px;
}

.uploader-header h2 {
  color: #2c3e50;
  font-size: 32px;
  font-weight: 700;
  margin: 0 0 8px 0;
  background: linear-gradient(135deg, #409EFF, #66B3FF);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.uploader-header p {
  color: #7f8c8d;
  font-size: 16px;
  margin: 0;
  font-weight: 400;
}

.uploader-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 20px;
  padding: 40px;
  box-shadow: 
    0 8px 32px rgba(64, 158, 255, 0.15),
    0 2px 8px rgba(64, 158, 255, 0.1);
  border: 1px solid rgba(64, 158, 255, 0.1);
  max-width: 1000px;
  width: 100%;
}

.avatar-uploader {
  width: 100%;
}

.upload-area {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 15px 15px;
  border: 2px dashed #409EFF;
  border-radius: 16px;
  background: #f8fbff;
  transition: all 0.3s ease;
  cursor: pointer;
}

.upload-area:hover {
  border-color: #66B3FF;
  background: #f0f7ff;
  transform: translateY(-2px);
  box-shadow: 0 4px 15px rgba(64, 158, 255, 0.2);
}

.upload-icon {
  font-size: 48px;
  color: #409EFF;
  margin-bottom: 12px;
}

.upload-text {
  font-size: 18px;
  font-weight: 600;
  color: #409EFF;
  margin-bottom: 4px;
}

.upload-subtext {
  font-size: 12px;
  color: #7f8c8d;
}

:deep(.el-upload-list--picture-card) {
  display: flex;
  justify-content: center;
  flex-wrap: wrap;
  gap: 12px;
  margin-top: 20px;
}

:deep(.el-upload-list--picture-card .el-upload-list__item) {
  border-radius: 12px;
  border: 2px solid #e6f3ff;
  transition: all 0.3s ease;
}

:deep(.el-upload-list--picture-card .el-upload-list__item:hover) {
  border-color: #409EFF;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(64, 158, 255, 0.2);
}

:deep(.el-upload__tip) {
  text-align: center;
  color: #7f8c8d;
  font-size: 14px;
  margin-top: 16px;
}

.button-wrapper {
  display: flex; /* 启用 Flex 布局 */
  justify-content: center; /* 水平居中 */
  width: 100%; /* 确保容器宽度足够 */
  margin-top: 30px; /* 保持与原 submit-btn 的上边距一致 */
}

.submit-btn {
  padding: 14px 32px;
  background: linear-gradient(135deg, #409EFF, #66B3FF);
  border: none;
  border-radius: 16px;
  font-weight: 600;
  font-size: 16px;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 8px;
  box-shadow: 0 4px 15px rgba(64, 158, 255, 0.3);
}

.submit-btn:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(64, 158, 255, 0.4);
  background: linear-gradient(135deg, #66B3FF, #409EFF);
}

.submit-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

.btn-icon {
  font-size: 18px;
}

.btn-text {
  font-weight: 600;
}

.file-count {
  background: rgba(255, 255, 255, 0.2);
  padding: 2px 8px;
  border-radius: 12px;
  font-size: 12px;
  margin-left: 4px;
}

.upload-summary {
  margin-top: 20px;
  padding: 16px;
  background: linear-gradient(135deg, #f8fbff, #f0f7ff);
  border-radius: 12px;
  border: 1px solid rgba(64, 158, 255, 0.1);
}

.summary-content {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.success-icon {
  color: #67C23A;
  font-size: 18px;
}

.summary-text {
  color: #2c3e50;
  font-weight: 500;
}

.success-count {
  color: #67C23A;
  font-weight: 600;
}

.fail-count {
  color: #F56C6C;
  font-weight: 600;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .image-uploader-container {
    padding: 20px 15px;
  }
  
  .uploader-card {
    padding: 30px 20px;
  }
  
  .uploader-header h2 {
    font-size: 28px;
  }
  
  .upload-area {
    padding: 30px 15px;
  }
  
  .upload-icon {
    font-size: 40px;
  }
  
  .upload-text {
    font-size: 16px;
  }
}

@media (max-width: 480px) {
  .uploader-header h2 {
    font-size: 24px;
  }
  
  .uploader-card {
    padding: 20px 15px;
  }
  
  .upload-area {
    padding: 20px 10px;
  }
}
</style>