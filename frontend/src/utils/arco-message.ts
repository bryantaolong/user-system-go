import Message from '@arco-design/web-react/es/Message'
import { Modal } from '@arco-design/web-react'
import type { InputRef } from '@arco-design/web-react/es/Input'

export const ArcoMessage = {
  success: (content: string) => message.success(content),
  error: (content: string) => message.error(content),
  warning: (content: string) => message.warning(content),
  info: (content: string) => message.info(content),
  loading: (content?: string) => message.loading(content),
}

export const ArcoMessageBox = {
  confirm: (content: string, title: string = '提示', options?: { confirmText?: string; cancelText?: string; type?: 'info' | 'success' | 'warning' | 'error' }) => {
    return Modal.confirm({
      title,
      content,
      okText: options?.confirmText || '确定',
      cancelText: options?.cancelText || '取消',
      ...(options?.type && { type: options.type }),
    })
  },
  prompt: (title: string, content?: string, options?: { confirmText?: string; cancelText?: string; placeholder?: string; defaultValue?: string }) => {
    return new Promise((resolve, reject) => {
      const inputValue = options?.defaultValue || ''

      Modal.open({
        title,
        content: () => (
          <div>
            {content ? <p style={{ marginBottom: '8px', color: 'var(--color-text-2)' }}>{content}</p> : null}
            <input
              placeholder={options?.placeholder || ''}
              defaultValue={inputValue}
              id="arco-prompt-input"
              style={{ width: '100%' }}
            />
          </div>
        ),
        okText: options?.confirmText || '确定',
        cancelText: options?.cancelText || '取消',
        onOk: () => {
          const input = document.getElementById('arco-prompt-input') as HTMLInputElement | null
          const val = input?.value || ''
          if (!val) {
            return false
          }
          resolve({ value: val })
        },
        onCancel: () => {
          reject('cancel')
        },
      })
    })
  },
}
