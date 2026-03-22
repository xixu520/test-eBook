import { MockMethod } from 'vite-plugin-mock'

export default [
  {
    url: '/api/v1/auth/login',
    method: 'post',
    response: ({ body }: { body: any }) => {
      const { username, password } = body
      const users: Record<string, any> = {
        admin: { id: 1, role: 'admin' },
        editor: { id: 2, role: 'editor' },
        user: { id: 3, role: 'user' },
      }

      if (users[username] && password === '123456') {
        const user = users[username]
        return {
          code: 200,
          message: '登录成功',
          data: {
            token: `mock-token-${username}-${Date.now()}`,
            user: {
              id: user.id,
              username: username,
              role: user.role,
              theme: 'light',
            },
          },
        }
      } else {
        return {
          code: 400,
          message: '用户名或密码错误',
          data: null,
        }
      }
    },
  },
  {
    url: '/api/v1/auth/me',
    method: 'get',
    response: ({ headers }: any) => {
      const token = headers.authorization || ''
      let role = 'admin'
      let username = 'admin'
      let id = 1

      if (token.includes('editor')) {
        role = 'editor'
        username = 'editor'
        id = 2
      } else if (token.includes('user')) {
        role = 'user'
        username = 'user'
        id = 3
      }

      return {
        code: 200,
        message: '获取成功',
        data: {
          id,
          username,
          role,
          permissions: role === 'admin' ? ['upload', 'download', 'delete', 'manage_category'] : (role === 'editor' ? ['upload', 'download'] : ['download']),
        },
      }
    },
  },
] as MockMethod[]
