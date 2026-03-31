import request from '@/utils/request'

export function getDashboardStats() {
  return request({
    url: 'admin/dashboard',
    method: 'get',
    silent: true
  })
}

export function getSystemStatus() {
  return request({
    url: 'admin/system/status',
    method: 'get',
    silent: true
  })
}
