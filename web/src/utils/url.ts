// 获取图片 URL（缩略图），自动从 localStorage 获取 userID
export const getImageUrl = (imagePath: string) => {
    const userID = localStorage.getItem('userID');
    if (!userID) {
        console.warn('未找到用户 ID');
        return '';
    }
    return `http://localhost:8888/public/${userID}/thumbnail/${imagePath}`;
}

// 获取原图 URL，自动从 localStorage 获取 userID  
export const getOriginImageUrl = (imagePath: string) => {
    const userID = localStorage.getItem('userID');
    if (!userID) {
        console.warn('未找到用户 ID');
        return '';
    }
    return `http://localhost:8888/public/${userID}/origin/${imagePath}`;
}