import request from '@/utils/request'

export interface AuditLog {
  id: number
  created_at: string
  user_id: number
  username: string
  action: string
  details: string
  ip: string
}

export function getAuditLogs(params: { page: number; page_size: number; action?: string }) {
  return request<{ list: AuditLog[]; total: number }>({
    url: '/audit-logs',
    method: 'get',
    params
  })
}
