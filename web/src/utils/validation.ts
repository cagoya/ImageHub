// 表单验证工具函数

/**
 * 验证邮箱格式
 * @param email 邮箱地址
 * @returns 是否为有效邮箱
 */
export const isValidEmail = (email: string): boolean => {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  return emailRegex.test(email);
};

/**
 * 验证密码强度
 * @param password 密码
 * @returns 密码是否有效（长度至少6位）
 */
export const isValidPassword = (password: string): boolean => {
  return password.length >= 6;
};

/**
 * 验证用户名
 * @param username 用户名
 * @returns 用户名是否有效
 */
export const isValidUsername = (username: string): boolean => {
  return username.length >= 1 && username.length <= 20;
};
