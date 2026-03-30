import request from '@/utils/request'

export function login(data: any) {
  return request({
    url: '/auth/login',
    method: 'post',
    data,
  })
}

export function getInfo() {
  return request({
    url: '/auth/me',
    method: 'get',
  })
}

export function register(data: any) {
  return request({
    url: '/auth/register',
    method: 'post',
    data,
  })
}

