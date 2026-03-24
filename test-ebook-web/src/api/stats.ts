import request from '@/utils/request'

export function getDashboardStats() {
  return request({
    url: '/admin/dashboard',
    method: 'get',
  })
}
