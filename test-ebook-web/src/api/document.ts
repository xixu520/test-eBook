import request from '@/utils/request'

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
  publisher?: string
  implementation_date?: string
  implementation_status?: string
  verify_status?: string
  sync_status?: string
  created_at: string
  field_values?: Array<{
    field_id: number
    value: string
    field?: {
      label: string
      field_key: string
    }
  }>
}

export interface DocumentQuery {
  page?: number
  page_size?: number
  keyword?: string
  category_id?: number
  publisher?: string
  implementation_status?: string
  start_date?: string
  end_date?: string
  [key: string]: any // 支持 filter[id] 这种动态 key
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

export function updateDocument(id: number, data: any) {
  return request({
    url: `/documents/${id}`,
    method: 'put',
    data
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

export interface Announcement {
  id: number
  title: string
  content: string
  active: boolean
  created_at: string
  updated_at: string
}

export function getActiveAnnouncement() {
  return request<Announcement>({
    url: '/announcements/active',
    method: 'get',
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

// --- Sync Tasks ---

export function retrySync(docId: number) {
  return request({
    url: `/documents/${docId}/retry-sync`,
    method: 'post'
  })
}

export function getUploadTasks() {
  return request({
    url: '/admin/upload-tasks',
    method: 'get'
  })
}
