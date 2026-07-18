import { useState, useEffect } from 'react'
import { Card, Select, InputNumber, Button, Empty } from '@arco-design/web-react'
import message from '@arco-design/web-react/es/Message'
import * as logApi from '@/api/system/log'

const SystemLog = () => {
  const [loading, setLoading] = useState(false)
  const [logs, setLogs] = useState<string[]>([])
  const [lineCount, setLineCount] = useState(200)
  const [logFiles, setLogFiles] = useState<string[]>([])
  const [selectedFile, setSelectedFile] = useState<string>('')

  const logsText = logs.join('\n')

  const loadFiles = async () => {
    try {
      const res = await logApi.listLogFiles()
      if (res.code === 200) {
        setLogFiles(res.data || [])
        if (!selectedFile && (res.data || []).length > 0) {
          setSelectedFile(res.data![0])
        }
      } else {
        message.error(res.message || '加载日志文件列表失败')
      }
    } catch (error) {
      console.error('加载日志文件列表失败:', error)
      message.error('加载日志文件列表失败，请稍后重试')
    }
  }

  const loadLogs = async () => {
    setLoading(true)
    try {
      const res = await logApi.listLatestLogs(lineCount, selectedFile || undefined)
      if (res.code === 200) {
        setLogs(res.data || [])
      } else {
        message.error(res.message || '加载日志失败')
      }
    } catch (error) {
      console.error('加载系统日志失败:', error)
      message.error('加载系统日志失败，请稍后重试')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadFiles()
  }, [])

  useEffect(() => {
    if (selectedFile) {
      loadLogs()
    }
  }, [selectedFile, lineCount])

  return (
    <div className="system-log">
      <Card className="log-card">
        <div className="log-header">
          <div className="title-section">
            <h2>系统日志</h2>
            <p className="subtitle">查看后台应用运行日志，仅管理员可访问</p>
          </div>
          <div className="actions">
            <Select
              value={selectedFile}
              placeholder="选择日志文件"
              size="small"
              className="file-select"
              onChange={setSelectedFile}
              style={{ width: 200 }}
            >
              {logFiles.map((file) => (
                <Select.Option key={file} value={file}>
                  {file}
                </Select.Option>
              ))}
            </Select>
            <span className="lines-label">行数：</span>
            <InputNumber
              value={lineCount}
              onChange={setLineCount}
              min={50}
              max={2000}
              step={50}
              size="small"
            />
            <Button type="primary" loading={loading} onClick={loadLogs}>
              刷新
            </Button>
          </div>
        </div>

        <div style={{ margin: '16px 0', borderTop: '1px solid var(--color-border-2)' }} />

        {!loading && logs.length === 0 ? (
          <Empty description="暂无日志数据" />
        ) : (
          <div style={{ maxHeight: 600, backgroundColor: '#1e1e1e', borderRadius: 8, padding: 12, overflow: 'auto' }}>
            <pre style={{ margin: 0, fontFamily: 'Consolas, Menlo, Monaco, Courier New, monospace', fontSize: 12, color: '#d4d4d4', lineHeight: 1.5, whiteSpace: 'pre-wrap' }}>
              <code>{logsText}</code>
            </pre>
          </div>
        )}
      </Card>
    </div>
  )
}

export default SystemLog
