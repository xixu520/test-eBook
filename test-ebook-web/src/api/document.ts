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
  standard_no: string
  name: string
  version?: string
  is_latest?: boolean
  publisher: string
  status: 'current' | 'obsolete' | 'upcoming'
  category_id: number
  category_name: string
  issue_date: string
  implement_date: string
  obsolete_date?: string | null
  ocr_status: 'pending' | 'processing' | 'completed' | 'failed'
  verify_status: 'pending' | 'pass' | 'retry'
  uploader_name: string
  upload_time: string
}

export function getDocuments(params: DocumentQuery) {
  return request({
    url: '/documents',
    method: 'get',
    params,
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
