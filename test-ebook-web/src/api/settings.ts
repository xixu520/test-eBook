import request from '@/utils/request'

export interface SystemSetting {
  id?: number
  key: string
  value: string
  description?: string
}

// 将扁平的 Key-Value 数组还原为嵌套对象
function unflatten(settings: SystemSetting[]): Record<string, any> {
	const result: Record<string, any> = {}
	for (const item of settings) {
		const keys = item.key.split('.')
		let current = result
		for (let i = 0; i < keys.length - 1; i++) {
			if (!(keys[i] in current)) {
				current[keys[i]] = {}
			}
			current = current[keys[i]]
		}
		// 类型转换：尝试还原 boolean / number / array 类型
		const lastKey = keys[keys.length - 1]
		const val = item.value
		if (val === 'true') {
			current[lastKey] = true
		} else if (val === 'false') {
			current[lastKey] = false
		} else if (/^\d+$/.test(val)) {
			current[lastKey] = Number(val)
		} else if (val.startsWith('[') || val.startsWith('{')) {
			try {
				current[lastKey] = JSON.parse(val)
			} catch {
				current[lastKey] = val
			}
		} else {
			current[lastKey] = val
		}
	}
	return result
}

export async function getSettings(): Promise<Record<string, any>> {
	const list = await request<SystemSetting[]>({
		url: '/settings',
		method: 'get'
	})
	return unflatten(list as any as SystemSetting[])
}

export function saveSettings(data: any) {
	// Convert nested object to flattened Key-Value list for backend storage
	const settingsList: SystemSetting[] = []

	const flatten = (obj: any, prefix = '') => {
		for (const key in obj) {
			const fullKey = prefix ? `${prefix}.${key}` : key
			if (typeof obj[key] === 'object' && !Array.isArray(obj[key]) && obj[key] !== null) {
				flatten(obj[key], fullKey)
			} else {
				settingsList.push({
					key: fullKey,
					value: typeof obj[key] === 'object' ? JSON.stringify(obj[key]) : String(obj[key])
				})
			}
		}
	}

	flatten(data)

	return request({
		url: '/settings',
		method: 'put',
		data: settingsList
	})
}

export function testOcrConnection(ocrSettings: any) {
	return request({
		url: '/settings/ocr-test',
		method: 'post',
		data: {
			engine: ocrSettings.engine,
			config: ocrSettings.engine === 'paddleocr' ? ocrSettings.paddle_config : ocrSettings.baidu_config
		}
	})
}

export function testStorageConnection(storageSettings: any) {
	return request({
		url: '/settings/storage-test',
		method: 'post',
		data: {
			type: storageSettings.type,
			config: storageSettings
		}
	})
}

export function orphanScan() {
	return request({
		url: '/settings/orphan-scan',
		method: 'post'
	})
}
