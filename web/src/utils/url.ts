import { backendURL } from '@/api';

// 获取图片 URL（缩略图），自动从 localStorage 获取 userID
export const getImageUrl = (imagePath: string) => {
    const userID = localStorage.getItem('userID');
    if (!userID) {
        console.warn('未找到用户 ID');
        return '';
    }
    return `${backendURL}/public/${userID}/thumbnail/${imagePath}`;
}

// 获取原图 URL，自动从 localStorage 获取 userID  
export const getOriginImageUrl = (imagePath: string) => {
    const userID = localStorage.getItem('userID');
    if (!userID) {
        console.warn('未找到用户 ID');
        return '';
    }
    return `${backendURL}/public/${userID}/origin/${imagePath}`;
}

// 获取AI图片 URL（用于AI回复中的图片）
// imagePath 格式为 "ai/{filename}" 或 "{filename}"
export const getAIImageUrl = (imagePath: string) => {
    if (!imagePath) {
        return '';
    }
    // 如果路径已经包含 "ai/"，直接使用
    if (imagePath.startsWith('ai/')) {
        return `${backendURL}/public/${imagePath}`;
    }
    // 否则，假设是文件名，添加 "ai/" 前缀
    return `${backendURL}/public/ai/${imagePath}`;
}