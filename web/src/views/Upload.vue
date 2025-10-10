<template>
  <AppHeader />
    <div class="image-uploader-container">
      <h2>批量图片上传</h2>
  
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
        multiple >
        <el-icon><Plus /></el-icon>
        <template #tip>
          <div class="el-upload__tip">
            最多上传10 张JPG/PNG 图片，单张大小不能超过20MB。
          </div>
        </template>
      </el-upload>
  
      <el-button
        type="primary"
        :loading="isUploading"
        @click="submitUpload"
        style="margin-top: 20px"
        :disabled="fileList.length === 0 || isUploading"
      >
        {{ isUploading ? '上传中...' : '提交上传' }}
      </el-button>
  
      <div v-if="uploadSummary.successCount + uploadSummary.failCount > 0" style="margin-top: 15px;">
          <p>
              上传完成：
              <span style="color: green;">成功 {{ uploadSummary.successCount }} 张</span>，
              <span style="color: red;">失败 {{ uploadSummary.failCount }} 张</span>
          </p>
      </div>
    </div>
    <el-image-viewer
    v-if="imageViewerVisible"       :url-list="[dialogImageUrl]"    @close="closeImageViewer"       />
  </template>
  
  <script setup lang="ts">
  import { ref } from 'vue';
  import { ElMessage, type UploadProps, type UploadRawFile, type UploadFile, type UploadUserFile } from 'element-plus';
  import { Plus } from '@element-plus/icons-vue';
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
      // 同步文件列表
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
    // 更新 fileList：Element Plus 默认会从 uploadFiles 中移除，我们只需同步
    fileList.value = uploadFiles;

    // 释放本地 URL 资源 (针对用户刚选择，但还没上传的图片)
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
      // 允许任何以 'image/' 开头的 MIME 类型
      const isTypeValid = rawFile.type.startsWith('image/'); 
      // 20MB 限制
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
  
  // 覆盖默认上传行为的钩子，由于我们使用 submitUpload 手动触发，这里不做处理
  const customUpload: UploadProps['httpRequest'] = async () => {
      return;
  };
  
  
  // 批量提交
  const submitUpload = async () => {
      // 获取所有待上传文件
      const filesToUpload = fileList.value.map(f => f.raw as UploadRawFile);
  
      if (filesToUpload.length === 0) {
          ElMessage.warning('没有需要上传的文件');
          return;
      }
  
      isUploading.value = true;
      uploadSummary.value = { successCount: 0, failCount: 0 };
      
      // 使用 Promise.allSettled 并发上传所有文件
      const uploadPromises = filesToUpload.map((file) => 
          ImageAPI.upload(file).then(response => ({ file, response }))
      );
  
      try {
          const responses = await Promise.allSettled(uploadPromises);
  
          // 收集需要移除的文件
          const filesToRemove: UploadRawFile[] = [];
  
          responses.forEach(item => {
              if (item.status === 'fulfilled') {
                  const { file, response } = item.value;
  
                  if (response.data.code === 200) {
                      uploadSummary.value.successCount++;
                      // 上传成功，标记为待移除
                      filesToRemove.push(file);
                  } else {
                      uploadSummary.value.failCount++;
                  }
              } else {
                  // 请求本身失败 (网络错误等)
                  uploadSummary.value.failCount++;
              }
          });
  
          // 释放被移除文件的 URL 资源
          filesToRemove.forEach(file => {
              const fileInList = fileList.value.find(f => f.raw === file);
              if (fileInList?.url) {
                  URL.revokeObjectURL(fileInList.url);
              }
          });
          
          // 移除上传成功的文件
          fileList.value = fileList.value.filter(f => !filesToRemove.includes(f.raw as UploadRawFile));
          
          // 最终提示
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
    padding: 20px;
  }
  /* 样式保持不变 */
  .avatar-uploader .el-upload {
    border: 1px dashed var(--el-border-color);
    border-radius: 6px;
    cursor: pointer;
    position: relative;
    overflow: hidden;
    transition: var(--el-transition-duration-fast);
  }
  
  .avatar-uploader .el-upload:hover {
    border-color: var(--el-color-primary);
  }
  
  .el-icon.avatar-uploader-icon {
    font-size: 28px;
    color: #8c939d;
    width: 148px; 
    height: 148px;
    text-align: center;
  }
  </style>