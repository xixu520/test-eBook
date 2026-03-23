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
            baidu_config: {
              app_id: 'xxxxxxx',
              api_key: 'xxxxxxx',
              secret_key: 'xxxxxxx'
            },
            paddle_config: {
              token: '02ac470b6efd58e73df5818608dfc0a59b30e6a8',
              model: 'PaddleOCR-VL-1.5',
              use_doc_orientation_classify: false,
              use_doc_unwarping: false,
              use_chart_recognition: false
            }
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
  // 测试 OCR 连接
  {
    url: '/api/v1/settings/test-ocr',
    method: 'post',
    timeout: 1500,
    response: ({ body }: any) => {
      if (body.engine === 'baidu' && (!body.baidu_config?.api_key || !body.baidu_config?.secret_key)) {
        return { code: 400, message: '百度云 API Key 和 Secret Key 不能为空' }
      }
      if (body.engine === 'paddleocr') {
        if (!body.paddle_config?.token) {
          return { code: 400, message: 'PaddleOCR 官方认证 Token 不能为空' }
        }
        if (body.paddle_config.token !== '02ac470b6efd58e73df5818608dfc0a59b30e6a8') {
          return { code: 401, message: '认证失败：无效的 Token，请检查 AI Studio 平台配置' }
        }
      }
      return { code: 200, message: '连接测试成功', data: null }
    }
  }
] as MockMethod[]
