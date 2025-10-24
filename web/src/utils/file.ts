// 文件处理工具函数

/**
 * 将base64字符串转换为File对象
 * @param base64 base64字符串
 * @param filename 文件名
 * @returns File对象
 */
export const base64ToFile = (base64: string, filename: string): File => {
  const arr = base64.split(',');
  const match = arr[0]?.match(/:(.*?);/);
  const mime = match?.[1] || 'image/jpeg';
  const bstr = atob(arr[1] || '');
  let n = bstr.length;
  const u8arr = new Uint8Array(n);
  while (n--) {
    u8arr[n] = bstr.charCodeAt(n);
  }
  return new File([u8arr], filename, { type: mime });
};

/**
 * 获取文件扩展名
 * @param filename 文件名
 * @returns 文件扩展名
 */
export const getFileExtension = (filename: string): string => {
  return filename.split('.').pop()?.toLowerCase() || '';
};

/**
 * 检查文件类型是否为图片
 * @param file 文件对象
 * @returns 是否为图片文件
 */
export const isImageFile = (file: File): boolean => {
  const imageTypes = ['image/jpeg', 'image/jpg', 'image/png', 'image/gif', 'image/webp'];
  return imageTypes.includes(file.type);
};

/**
 * 检查文件大小
 * @param file 文件对象
 * @param maxSizeMB 最大大小（MB）
 * @returns 是否超过大小限制
 */
export const isFileSizeValid = (file: File, maxSizeMB: number = 10): boolean => {
  const maxSizeBytes = maxSizeMB * 1024 * 1024;
  return file.size <= maxSizeBytes;
};
