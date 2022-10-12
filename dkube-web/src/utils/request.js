import axios from 'axios';

const httpClient = axios.create({
    validateStatus(status) {
        return status >= 200 && status < 504
    },
    timeout: 10000
});

httpClient.defaults.retry = 3 
httpClient.defaults.retryDelay = 1000
httpClient.defaults.shouldRetry = true

httpClient.interceptors.request.use (
    config => {
        config.headers['Content-Type'] = 'application/json'
        config.headers['Accept-Language'] = 'zh-CN'
        config.headers['Authorization'] = localStorage.getItem('token') 
        if (config.method === 'post') {
            if (!config.data) {
                config.data = {}
            }
        }
        return config
    },
    err => {
        Promise.reject(err)
    }
);

httpClient.interceptors.response.use (
    response => {
        if (response.status !== 200) {
            return Promise.reject(response.data)
        } else {
            return response.data
        }
    },
    err => {
        return Promise.reject(err)
    }
);

export default httpClient;