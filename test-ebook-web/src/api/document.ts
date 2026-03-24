import request from '@/utils/request'

export interface DocumentQuery {
  page?: number
  size?: number
  keyword?: string
  category_id?: number
  publisher?: string
  status?: string
  start_date?: string
  end_date?: string
}

export interface Document {
  id: number
  number: string // standard_no -> number
  title: string  // name -> title
  year?: string
  version?: string
  category_id: number
  file_path: string
  file_size: number
  status: number // 0: processing, 1: processed, 2: failed
  ocr_content?: string
  created_at: string
}

export function getDocuments(params: DocumentQuery) {
  return request({
    url: '/documents',
    method: 'get',
    params,
  })
}

export function uploadFile(formData: FormData) {
  return request({
    url: '/documents/upload',
    method: 'post',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

export function getDocumentDetail(id: number) {
  return request({
    url: `/documents/${id}`,
    method: 'get'
  })
}

export function deleteDocument(id: number) {
  return request({
    url: `/documents/${id}`,
    method: 'delete'
  })
}

// --- OCR Tasks ---

export function getTaskStatus(taskId: string) {
  return request({
    url: `/tasks/${taskId}/status`,
    method: 'get'
  })
}

export function retryOCR(docId: number) {
  return request({
    url: `/documents/${docId}/ocr/retry`,
    method: 'post'
  })
}

export function getActiveAnnouncement() {
  return request({
    url: '/announcements/active',
    method: 'get',
  })
}

export function getDocumentHistory(standardNo: string) {
  return request({
    url: '/documents/history',
    method: 'get',
    params: { standard_no: standardNo },
  })
}

export function getRecycleBinDocuments() {
  return request({
    url: '/recycle-bin/documents',
    method: 'get'
  })
}

export function restoreDocuments(document_ids: number[]) {
  return request({
    url: '/recycle-bin/documents/restore',
    method: 'put',
    data: { document_ids }
  })
}

export function batchDeleteDocuments(document_ids: number[], empty_all = false) {
  return request({
    url: '/recycle-bin/documents/batch-delete',
    method: 'post',
    data: { document_ids, empty_all }
  })
}
