import request from '@/utils/request'

export function getAuditLogs(params: any) {
  return request({
    url: '/audit-logs',
    method: 'get',
    params
  })
}
