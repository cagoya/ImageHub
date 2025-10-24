// 日期时间工具函数

/**
 * 格式化日期字符串
 * @param dateString 日期字符串
 * @returns 格式化后的日期字符串
 */
export const formatDate = (dateString: string): string => {
  const date = new Date(dateString);
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  });
};

/**
 * 格式化时间字符串
 * @param timeString 时间字符串
 * @returns 格式化后的时间字符串
 */
export const formatTime = (timeString?: string): string => {
  if (!timeString) return '无时间信息';
  const date = new Date(timeString);
  return date.toLocaleString('zh-CN');
};

/**
 * 获取当前时间字符串
 * @returns 当前时间的ISO字符串
 */
export const getCurrentTimeString = (): string => {
  return new Date().toISOString();
};
