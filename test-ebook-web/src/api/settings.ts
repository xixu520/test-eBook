import request from '@/utils/request'

export function getSettings() {
  return request({
    url: '/settings',
    method: 'get',
  })
}

export function saveSettings(data: any) {
  return request({
    url: '/settings',
    method: 'post',
    data,
  })
}

export function getSystemStatus() {
  return request({
    url: '/system/status',
    method: 'get',
  })
}
