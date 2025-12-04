import axios from 'axios';
import { type UploadRawFile } from 'element-plus';

// 后端的地址，从环境变量获取，如果没有则使用相对路径（适用于 Docker nginx 代理）
// 开发环境可以通过 .env 文件设置 VITE_BACKEND_URL
// Docker 构建时可以通过构建参数设置
const getBackendURL = (): string => {
  // 优先使用环境变量
  const envURL = import.meta.env.VITE_BACKEND_URL;
  if (envURL) {
    return envURL;
  }
  // 如果没有设置环境变量，使用相对路径（适用于 Docker 中的 nginx 代理）
  // 在开发环境中，如果需要直接访问后端，可以在 .env 文件中设置 VITE_BACKEND_URL
  return '';
};

export const backendURL = getBackendURL();

// 构建 API URL：如果 backendURL 为空，则使用相对路径 /api（适用于 nginx 代理）
// 否则使用完整的 URL
const apiURL = backendURL ? `${backendURL}/api` : '/api';

// 创建axios实例
const instance = axios.create({
    baseURL: apiURL,
    timeout: 5000,
    headers: {
      'Content-Type': 'application/json'
    }
  });


// 请求拦截器添加认证信息
instance.interceptors.request.use(
    (config) => {
      // 检查是否为不需要认证的接口
      const isAuthRequest = config.url?.includes('/login') || 
                           config.url?.includes('/register') ||
                           config.url?.includes('/forget-password');
      
      // 如果不是认证相关接口，则添加token
      if (!isAuthRequest) {
        const token = localStorage.getItem('userToken');
        if (token) {
          config.headers.Authorization = `Bearer ${token}`;
        }
      }
      return config;
    },
    (error) => {
      return Promise.reject(error);
    }
  );
  
  // 响应拦截器处理错误
  instance.interceptors.response.use(
    (response) => {
    // 1. 获取响应头
    const headers = response.headers;

    // 2. 检查特定的 Token 头部
    // 后端使用 'x-auth-token' 这个自定义头部来发送新的 Token
    // 注意：axios 会将响应头名称转换为小写
    const newToken = headers['x-auth-token'] as string | undefined;

    if (newToken) {
      console.log('在响应头中捕获到新的 Token:', newToken);
      
      // 3. 将新的 Token 存储起来
      localStorage.setItem('userToken', newToken);
    }

      return response;
    },
    async (error) => {
      const originalRequest = error.config;
      
      // 检查是否为401未授权错误
      if (error.response?.status === 401) {
        console.log('收到401响应，清除用户信息并跳转到登录页');
        
        // 清除localStorage中的用户信息
        localStorage.removeItem('userToken');
        localStorage.removeItem('userInfo');
        localStorage.removeItem('userID');
        
        // 跳转到登录页
        window.location.href = '/login';
        return Promise.reject(error.response?.data || '请求失败');
      }
      
      // 检查是否为403禁止访问错误
      if (error.response?.status === 403) {
        console.log('收到403响应，尝试刷新token');
        
        // 检查是否已经重试过，避免无限循环
        if (originalRequest._retry) {
          console.log('已经重试过，不再重试');
          return Promise.reject(error.response?.data || '请求失败');
        }
        
        originalRequest._retry = true;
        
        // 暂存原来的token等信息
        const oldToken = localStorage.getItem('userToken');
        const oldUserInfo = localStorage.getItem('userInfo');
        const oldUserID = localStorage.getItem('userID');
        const currentPath = window.location.pathname;
        
        try {
          // 调用刷新token的API
          const refreshResponse = await instance.post('/user/refresh-token');
          const newToken = refreshResponse.data.data.token;
          const newUserInfo = refreshResponse.data.data.user;
          
          if (newToken) {
            console.log('Token刷新成功，更新本地存储');
            // 更新token
            localStorage.setItem('userToken', newToken);
            // 更新用户信息
            localStorage.setItem('userInfo', JSON.stringify(newUserInfo));
            localStorage.setItem('userID', String(newUserInfo.id));
            
            // 触发自定义事件，通知其他组件用户信息已更新
            window.dispatchEvent(new CustomEvent('userInfoUpdated', { detail: newUserInfo }));
          }
        } catch (refreshError) {
          console.error('刷新token失败:', refreshError);
          // 如果刷新token失败，恢复原来的token等信息
          if (oldToken) {
            localStorage.setItem('userToken', oldToken);
          } else {
            localStorage.removeItem('userToken');
          }
          if (oldUserInfo) {
            localStorage.setItem('userInfo', oldUserInfo);
          } else {
            localStorage.removeItem('userInfo');
          }
          if (oldUserID) {
            localStorage.setItem('userID', oldUserID);
          } else {
            localStorage.removeItem('userID');
          }
        } finally {
          // 无论是否成功都重新访问当前页面
          window.location.href = currentPath;
        }
        
        return Promise.reject(error.response?.data || '请求失败');
      }
      
      return Promise.reject(error.response?.data || '请求失败');
    }
  );

  // 用户相关接口
  const UserAPI = {
    /**
     * 用户登录
     * @param username 用户名
     * @param password 密码
     */
    login: (username: string, password: string) => {
        return instance.post('/user/login', {
            username: username,
            password: password
        });
    },
    /**
     * 用户注册
     * @param username 用户名
     * @param email 邮箱
     * @param password 密码
     */
    register: (username: string, email: string, password: string) => {
        return instance.post('/user/register', {
            username: username,
            email: email,
            password: password
        });
    },
    /**
     * 找回密码
     * @param username 用户名
     * @param email 邮箱
     * @param newPassword 新密码
     */
    forgetPassword: (username: string, email: string, newPassword: string) => {
        return instance.post('/user/forget-password', {
          username: username,
          email: email,
          new_password: newPassword
        })
    },
    /**
     * 根据ID获取用户信息
     * @param id 用户ID
     */
    getUserByID: (id: number) => {
      return instance.get(`/user/id?id=${id}`)
    },
    /**
     * 根据用户名获取用户信息
     * @param username 用户名
     */
    getUserByUsername: (username: string) => {
      return instance.get(`/user/username?username=${encodeURIComponent(username)}`)
    },
    /**
     * 获取所有用户
     * @param page 页码
     * @param pageSize 页大小
     */
    getAllUsers: (page: number, pageSize: number) => {
      return instance.get(`/user/all?page=${page}&pageSize=${pageSize}`)
    },
    /**
     * 更新用户信息
     * @param id 用户ID
     * @param username 用户名
     * @param email 邮箱
     * @param password 密码（可选）
     */
    updateUser: (id: number, username: string, email: string, password?: string) => {
      const data: any = {
        id: id,
        username: username,
        email: email
      }
      if (password) {
        data.password = password
      }
      return instance.put('/user/update', data)
    },
    /**
     * 授予管理员权限
     * @param id 用户ID
     * @returns 
     */
    grantAdminPermission: (id: number) => {
      return instance.put(`/user/grant?id=${id}`)
    },
    /**
     * 撤销管理员权限
     * @param id 用户ID
     * @returns 
     */
    revokeAdminPermission: (id: number) => {
      return instance.put(`/user/revoke?id=${id}`)
    },
    /**
     * 封禁用户
     * @param id 用户ID
     * @param expiration 封禁时长（分钟）
     * @returns 
     */
    suspendUser: (id: number, expiration: number) => {
      return instance.post(`/user/suspend?id=${id}&expiration=${expiration}`)
    },
    /**
     * 解封用户
     * @param id 用户ID
     * @returns 
     */
    unsuspendUser: (id: number) => {
      return instance.post(`/user/unsuspend?id=${id}`)
    },
    /**
     * 设置是否解析 EXIF
     * @param parseExif 是否解析 EXIF
     * @returns 
     */
    setParseExif: (parseExif: boolean) => {
      return instance.put(`/user/parse-exif?parseExif=${parseExif}`)
    },
  }
  /**
   * 图片相关接口
   */
  const ImageAPI = {
    /**
     * 上传图片
     * @param image 图片
     */
    upload(image: UploadRawFile) {
      const formData = new FormData();
      formData.append('image', image);
      return instance.post(`/image/upload`, formData, {
        headers: {
          'Content-Type': 'multipart/form-data'
        }
      })
    },
    /**
     * 更新图片
     * @param imageID 图片ID
     * @param image 图片文件
     */
    update(imageID: number, image: File) {
      const formData = new FormData();
      formData.append('image', image);
      return instance.put(`/image/update?imageID=${imageID}`, formData, {
        headers: {
          'Content-Type': 'multipart/form-data'
        }
      })
    },
    /**
     * 删除图片
     * @param imageID 图片ID
     * @param userID 用户ID
     */
    delete(imageID: number, userID: number) {
      return instance.delete(`/image/delete?imageID=${imageID}&userID=${userID}`)
    },
    /**
     * 获取图片详细信息
     * @param imageID 图片ID
     */
    getDetail(imageID: number) {
      return instance.get(`/image/detail?imageID=${imageID}`)
    },
    /**
     * 获取图片文件
     * @param imageID 图片ID
     */
    getFile(imageID: number) {
      return instance.get(`/image/file?imageID=${imageID}`, {
        responseType: 'blob' // 告诉axios返回的是二进制数据
      })
    },
    /**
     * 获取所有图片（仅管理员）
     * @param page 页码
     * @param pageSize 页大小
     */
    getAllImages(page: number, pageSize: number) {
      return instance.get(`/image/all?page=${page}&pageSize=${pageSize}`)
    },
    /**
     * 根据用户ID获取图片
     * @param userID 用户ID
     * @param page 页码
     * @param pageSize 页大小
     */
    getImagesByUserID(userID: number, page: number, pageSize: number) {
      return instance.get(`/image/user-id?userID=${userID}&page=${page}&pageSize=${pageSize}`)
    },
    /**
     * 获取个人所有图片
     * @param page 页码
     * @param pageSize 页大小
     */
    getPersonalAllImages(page: number, pageSize: number) {
      // 从localStorage获取当前用户ID
      const currentUserID = localStorage.getItem('userID');
      if (!currentUserID) {
        return Promise.reject('用户未登录');
      }
      return instance.get(`/image/user-id?userID=${currentUserID}&page=${page}&pageSize=${pageSize}`)
    },
    /**
     * 根据标签查询图片
     * @param tag 标签
     * @param page 页码
     * @param pageSize 页大小
     */
    getImagesByTag(tag: string, page: number, pageSize: number) {
      return instance.get(`/image/tag?tag=${encodeURIComponent(tag)}&page=${page}&pageSize=${pageSize}`)
    },
    /**
     * 根据时间范围查询图片
     * @param startTime 开始时间 (格式: "2006-01-02 15:04:05")
     * @param endTime 结束时间 (格式: "2006-01-02 15:04:05")
     * @param page 页码
     * @param pageSize 页大小
     */
    getImagesByTimeRange(startTime: string, endTime: string, page: number, pageSize: number) {
      return instance.get(`/image/time-range?startTime=${encodeURIComponent(startTime)}&endTime=${encodeURIComponent(endTime)}&page=${page}&pageSize=${pageSize}`)
    },
    /**
     * 根据地区查询图片
     * @param province 省份（可选）
     * @param city 城市（可选）
     * @param district 区县（可选）
     * @param page 页码
     * @param pageSize 页大小
     */
    getImagesByArea(province?: string, city?: string, district?: string, page: number = 1, pageSize: number = 10) {
      let url = `/image/area?page=${page}&pageSize=${pageSize}`;
      if (province) {
        url += `&province=${encodeURIComponent(province)}`;
      }
      if (city) {
        url += `&city=${encodeURIComponent(city)}`;
      }
      if (district) {
        url += `&district=${encodeURIComponent(district)}`;
      }
      return instance.get(url)
    },
    /**
     * 设置图片标签
     * @param imageID 图片ID
     * @param aiTags AI标签数组
     * @param manualTags 手动标签数组
     */
    setTags(imageID: number, aiTags: string[], manualTags: string[]) {
      return instance.put(`/image/tags`, {
        image_id: imageID,
        ai_tags: aiTags,
        manual_tags: manualTags
      })
    },
    /**
     * 由AI生成图片标签
     * @param imageID 图片ID
     */
    generateTags(imageID: number) {
      return instance.get(`/image/generate-tags?imageID=${imageID}`)
    }
  }

  /**
   * MCP 相关接口
   */
  const MCPAPI = {
    /**
     * 获取用户的会话列表
     */
    getChats: () => {
      return instance.get('/mcp/chats')
    },
    /**
     * 获取会话消息
     * @param chatID 会话ID
     */
    getChatMessages: (chatID: number) => {
      return instance.get(`/mcp/messages?chatID=${chatID}`)
    }
  }

// 导出所有API模块
export { UserAPI, ImageAPI, MCPAPI };