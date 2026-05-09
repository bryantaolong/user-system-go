import axios from 'axios';

const instance = axios.create({
    timeout: 10000,
});

// 请求拦截器
instance.interceptors.request.use((config) => {
    const token = localStorage.getItem('token');
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
});

// 响应拦截器
instance.interceptors.response.use(
    (response) => {
        // 如果是 blob 类型，返回完整的 response 对象，以便于处理文件下载
        if (response.config.responseType === 'blob') {
            return response;
        }
        return response.data;
    },
    (error) => {
        if (error.response?.status === 401) {
            localStorage.removeItem('token');
            window.location.href = '/login';
        }
        return Promise.reject(error.response?.data || error);
    }
);

export default instance;