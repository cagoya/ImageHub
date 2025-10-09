import axios from 'axios';

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
    register: (username: string, password: string, email: string) => {
        return instance.post('/user/register', {
            Username: username,
            Password: password,
            Email: email
        });
    },
  }


// 导出所有API模块
export { UserAPI };