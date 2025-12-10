import { backendURL } from '@/api';

// 获取图片 URL（缩略图）
// @param imagePath 图片路径
// @param userID 图片所属用户的ID（可选，如果不提供则从 localStorage 获取当前登录用户的ID）
export const getImageUrl = (imagePath: string, userID?: number | string) => {
    let userId: string;
    if (userID !== undefined) {
        userId = String(userID);
    } else {
        userId = localStorage.getItem('userID') || '';
        if (!userId) {
            console.warn('未找到用户 ID');
            return '';
        }
    }
    return `${backendURL}/public/${userId}/thumbnail/${imagePath}`;
}

// 获取原图 URL
// @param imagePath 图片路径
// @param userID 图片所属用户的ID（可选，如果不提供则从 localStorage 获取当前登录用户的ID）
export const getOriginImageUrl = (imagePath: string, userID?: number | string) => {
    let userId: string;
    if (userID !== undefined) {
        userId = String(userID);
    } else {
        userId = localStorage.getItem('userID') || '';
        if (!userId) {
            console.warn('未找到用户 ID');
            return '';
        }
    }
    return `${backendURL}/public/${userId}/origin/${imagePath}`;
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