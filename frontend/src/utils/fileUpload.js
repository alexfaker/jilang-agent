/**
 * 文件上传工具函数
 */

/**
 * 验证文件类型
 * @param {File} file - 要验证的文件
 * @param {Array<string>} allowedTypes - 允许的MIME类型数组，如 ['image/jpeg', 'image/png']
 * @returns {boolean} 如果文件类型允许则返回true
 */
export const validateFileType = (file, allowedTypes) => {
  return allowedTypes.includes(file.type);
};

/**
 * 验证文件大小
 * @param {File} file - 要验证的文件
 * @param {number} maxSizeInMB - 最大文件大小（MB）
 * @returns {boolean} 如果文件大小在限制范围内则返回true
 */
export const validateFileSize = (file, maxSizeInMB) => {
  const maxSizeInBytes = maxSizeInMB * 1024 * 1024;
  return file.size <= maxSizeInBytes;
};

/**
 * 将文件转换为Base64字符串（用于预览）
 * @param {File} file - 要转换的文件
 * @returns {Promise<string>} 返回Base64字符串的Promise
 */
export const fileToBase64 = (file) => {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.readAsDataURL(file);
    reader.onload = () => resolve(reader.result);
    reader.onerror = (error) => reject(error);
  });
};

/**
 * 创建图片预览URL
 * @param {File} file - 图片文件
 * @returns {string} 预览URL
 */
export const createImagePreview = (file) => {
  return URL.createObjectURL(file);
};

/**
 * 释放图片预览URL
 * @param {string} url - 预览URL
 */
export const revokeImagePreview = (url) => {
  URL.revokeObjectURL(url);
};

/**
 * 压缩图片
 * @param {File|Blob} file - 图片文件
 * @param {Object} options - 压缩选项
 * @param {number} options.maxWidth - 最大宽度
 * @param {number} options.maxHeight - 最大高度
 * @param {number} options.quality - 质量（0-1）
 * @returns {Promise<Blob>} 返回压缩后的图片Blob的Promise
 */
export const compressImage = (file, { maxWidth = 800, maxHeight = 800, quality = 0.8 } = {}) => {
  return new Promise((resolve, reject) => {
    const img = new Image();
    img.src = URL.createObjectURL(file);
    
    img.onload = () => {
      URL.revokeObjectURL(img.src);
      
      let width = img.width;
      let height = img.height;
      
      // 计算缩放比例
      if (width > maxWidth) {
        height = Math.round((height * maxWidth) / width);
        width = maxWidth;
      }
      
      if (height > maxHeight) {
        width = Math.round((width * maxHeight) / height);
        height = maxHeight;
      }
      
      const canvas = document.createElement('canvas');
      canvas.width = width;
      canvas.height = height;
      
      const ctx = canvas.getContext('2d');
      ctx.drawImage(img, 0, 0, width, height);
      
      canvas.toBlob(
        (blob) => {
          if (blob) {
            resolve(blob);
          } else {
            reject(new Error('Canvas to Blob conversion failed'));
          }
        },
        file.type,
        quality
      );
    };
    
    img.onerror = () => {
      URL.revokeObjectURL(img.src);
      reject(new Error('Failed to load image'));
    };
  });
};

/**
 * 格式化文件大小
 * @param {number} bytes - 文件大小（字节）
 * @returns {string} 格式化后的文件大小
 */
export const formatFileSize = (bytes) => {
  if (bytes === 0) return '0 Bytes';
  
  const k = 1024;
  const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
};

/**
 * 从文件名获取扩展名
 * @param {string} filename - 文件名
 * @returns {string} 扩展名（小写）
 */
export const getFileExtension = (filename) => {
  return filename.split('.').pop().toLowerCase();
};

/**
 * 检查文件是否为图片
 * @param {File} file - 要检查的文件
 * @returns {boolean} 如果是图片则返回true
 */
export const isImageFile = (file) => {
  const imageTypes = ['image/jpeg', 'image/png', 'image/gif', 'image/webp', 'image/svg+xml'];
  return validateFileType(file, imageTypes);
};

export default {
  validateFileType,
  validateFileSize,
  fileToBase64,
  createImagePreview,
  revokeImagePreview,
  compressImage,
  formatFileSize,
  getFileExtension,
  isImageFile
}; 