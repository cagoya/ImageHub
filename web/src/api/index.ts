import axios from 'axios';
import { type UploadRawFile } from 'element-plus';

// 后端的地址，方便调整
export const backendURL = 'http://localhost:8888/api';

// 创建axios实例
const instance = axios.create({
    baseURL: backendURL,
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
                           config.url?.includes('/register');
      
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
    (error) => {
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
            Username: username,
            Password: password
        });
    },
    /**
     * 用户注册
     * @param username 用户名
     * @param password 密码
     * @param email 邮箱
     */
    register: (username: string, email: string, password: string) => {
        return instance.post('/user/register', {
            Username: username,
            Email: email,
            Password: password
        });
    },
    /**
     * 找回密码
     * @param username 用户名
     * @param newPassword 新密码
     * @param email 邮箱
     */
    forgetPassword: (username: string, email: string, newPassword: string) => {
        return instance.post('user/forgetPassword', {
          username: username,
          email: email,
          new_password: newPassword
        })
    }
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
      return instance.post('/image/upload', formData, {
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
     */
    delete(imageID: number, userID: number) {
      return instance.delete(`/image/delete?imageID=${imageID}&userID=${userID}`)
    },
    /**
     * 获取图片详细信息
     * @param imageID 图片ID
     */
    getDetail(imageID: number) {
      return instance.get(`/image/getDetail?imageID=${imageID}`)
    },
    /**
     * 获取图片文件
     * @param imageID 图片ID
     */
    getFile(imageID: number) {
      return instance.get(`/image/getFile?imageID=${imageID}`, {
        responseType: 'blob' // 告诉axios返回的是二进制数据
      })
    },
    /**
     * 获取所有图片
     * @param page 页码
     * @param pageSize 页大小
     */
    getAllImages(page: number, pageSize: number) {
      return instance.get(`/image/getAllImages?page=${page}&pageSize=${pageSize}`)
    },
    /**
     * 获取个人所有图片
     * @param page 页码
     * @param pageSize 页大小
     */
    getPersonalAllImages(page: number, pageSize: number) {
      return instance.get(`/image/getByUserID?page=${page}&pageSize=${pageSize}`)
    },
    /**
     * 根据标签查询图片
     * @param tag 标签
     * @param page 页码
     * @param pageSize 页大小
     */
    getImagesByTag(tag: string, page: number, pageSize: number) {
      return instance.get(`/image/getByTag?tag=${encodeURIComponent(tag)}&page=${page}&pageSize=${pageSize}`)
    },
    /**
     * 根据时间范围查询图片
     * @param startTime 开始时间
     * @param endTime 结束时间
     * @param page 页码
     * @param pageSize 页大小
     */
    getImagesByTimeRange(startTime: string, endTime: string, page: number, pageSize: number) {
      return instance.get(`/image/getByTimeRange?startTime=${encodeURIComponent(startTime)}&endTime=${encodeURIComponent(endTime)}&page=${page}&pageSize=${pageSize}`)
    },
    /**
     * 根据省份查询图片
     * @param province 省份
     * @param page 页码
     * @param pageSize 页大小
     */
    getImagesByProvince(province: string, page: number, pageSize: number) {
      return instance.get(`/image/getByProvince?province=${encodeURIComponent(province)}&page=${page}&pageSize=${pageSize}`)
    },
    /**
     * 根据城市查询图片
     * @param city 城市
     * @param page 页码
     * @param pageSize 页大小
     */
    getImagesByCity(city: string, page: number, pageSize: number) {
      return instance.get(`/image/getByCity?city=${encodeURIComponent(city)}&page=${page}&pageSize=${pageSize}`)
    },
    /**
     * 根据区县查询图片
     * @param district 区县
     * @param page 页码
     * @param pageSize 页大小
     */
    getImagesByDistrict(district: string, page: number, pageSize: number) {
      return instance.get(`/image/getByDistrict?district=${encodeURIComponent(district)}&page=${page}&pageSize=${pageSize}`)
    },
    /**
     * 设置图片标签
     * @param imageID 图片ID
     * @param tags 标签
     */
    setTags(imageID: number, tags: string[]) {
      return instance.put(`/image/setTags`, {
        image_id: imageID,
        tags: tags
      })
    }
  }


// 导出所有API模块
export { UserAPI, ImageAPI };