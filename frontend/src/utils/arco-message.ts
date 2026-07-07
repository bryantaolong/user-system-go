import { h, ref } from 'vue'
import { Message, Modal } from '@arco-design/web-vue'

export const ArcoMessage = {
  success: (content: string) => Message.success(content),
  error: (content: string) => Message.error(content),
  warning: (content: string) => Message.warning(content),
  info: (content: string) => Message.info(content),
  loading: (content?: string) => Message.loading(content),
}

export const ArcoMessageBox = {
  confirm: (content: string, title: string = '提示', options?: { confirmText?: string; cancelText?: string; type?: 'info' | 'success' | 'warning' | 'error' }) => {
    return Modal.confirm({
      title,
      content,
      okText: options?.confirmText || '确定',
      cancelText: options?.cancelText || '取消',
      ...(options?.type && { type: options.type })
    })
  },
  prompt: (title: string, content?: string, options?: { confirmText?: string; cancelText?: string; placeholder?: string; defaultValue?: string }) => {
    return new Promise((resolve, reject) => {
      const inputValue = ref(options?.defaultValue || '')

      Modal.open({
        title,
        content: () => h('div', {}, [
          content ? h('p', { style: { marginBottom: '8px', color: 'var(--color-text-2)' } }, content) : null,
          h('a-input', {
            placeholder: options?.placeholder || '',
            modelValue: inputValue.value,
            'onUpdate:modelValue': (val: string) => { inputValue.value = val },
            style: { width: '100%' },
            allowClear: true
          })
        ]),
        okText: options?.confirmText || '确定',
        cancelText: options?.cancelText || '取消',
        onOk: () => {
          if (!inputValue.value) {
            return false
          }
          resolve({ value: inputValue.value })
        },
        onCancel: () => {
          reject('cancel')
        }
      })
    })
  },
}
