import { MockMethod } from 'vite-plugin-mock'

export default [
  // 获取系统设置
  {
    url: '/api/v1/settings',
    method: 'get',
    response: () => {
      return {
        code: 200,
        message: '获取成功',
        data: {
          ocr: {
            engine: 'paddleocr',
            language: ['ch', 'en'],
            threads: 4,
            use_gpu: true,
          },
          storage: {
            type: 'local',
            local_path: '/data/ebook/storage',
            max_size: 50, // MB
          },
        },
      }
    },
  },
  // 保存系统设置
  {
    url: '/api/v1/settings',
    method: 'post',
    timeout: 1000,
    response: ({ body }: any) => {
      return {
        code: 200,
        message: '设置已更新',
        data: body,
      }
    },
  },
  // 获取系统资源状态
  {
    url: '/api/v1/system/status',
    method: 'get',
    response: () => {
      return {
        code: 200,
        message: '获取成功',
        data: {
          cpu: Math.floor(Math.random() * 40) + 10,
          memory: Math.floor(Math.random() * 30) + 40,
          disk: 65,
          uptime: '15d 4h 32m',
        },
      }
    },
  },
] as MockMethod[]
