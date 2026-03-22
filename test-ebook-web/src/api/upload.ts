import request from '@/utils/request'

export function uploadFile(formData: FormData, onProgress?: (event: any) => void) {
  return request({
    url: '/upload',
    method: 'post',
    data: formData,
    headers: { 'Content-Type': 'multipart/form-data' },
    onUploadProgress: onProgress,
  })
}

export function getOcrTasks() {
  return request({
    url: '/ocr/tasks',
    method: 'get',
  })
}
