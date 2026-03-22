import { MockMethod } from 'vite-plugin-mock'

export default [
  // 获取用户列表
  {
    url: '/api/v1/admin/users',
    method: 'get',
    response: () => {
      return {
        code: 200,
        message: '获取成功',
        data: [
          { id: 1, username: 'admin', role: 'admin', status: 1, last_login: '2024-03-22 10:00:00' },
          { id: 2, username: 'editor_zhang', role: 'editor', status: 1, last_login: '2024-03-21 15:30:00' },
          { id: 3, username: 'user_li', role: 'user', status: 1, last_login: '2024-03-22 09:12:00' },
          { id: 4, username: 'test_blocked', role: 'user', status: 0, last_login: '2024-03-10 11:20:00' },
        ],
      }
    },
  },
  // 修改用户状态
  {
    url: '/api/v1/admin/users/:id/status',
    method: 'put',
    response: ({ body }: any) => {
      return {
        code: 200,
        message: '状态已更新',
        data: body,
      }
    },
  },
  // 修改用户角色
  {
    url: '/api/v1/admin/users/:id/role',
    method: 'put',
    response: ({ body }: any) => {
      return {
        code: 200,
        message: '角色已更新',
        data: body,
      }
    },
  },
  // 删除用户
  {
    url: '/api/v1/admin/users/:id',
    method: 'delete',
    response: () => {
      return {
        code: 200,
        message: '用户已删除',
      }
    },
  },
] as MockMethod[]
