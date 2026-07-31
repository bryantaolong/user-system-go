<template>
  <div class="system-log">
    <a-card>
      <div class="log-header">
        <div class="title-section">
          <h2>系统日志</h2>
          <p class="subtitle">查看后台应用运行日志，仅管理员可访问</p>
        </div>
        <div class="actions">
          <a-select v-model="selectedFile" placeholder="选择日志文件" allow-clear style="width: 200px" @change="loadLogs">
            <a-option v-for="file in logFiles" :key="file" :label="file" :value="file" />
          </a-select>
          <span class="lines-label">行数：</span>
          <a-input-number v-model="lineCount" :min="1" :max="2000" :step="50" @change="loadLogs" />
          <a-button type="primary" :loading="loading" @click="loadLogs">刷新</a-button>
        </div>
      </div>

      <a-divider />

      <div v-if="!loading && logs.length === 0" class="empty-block">
        <a-empty description="暂无日志数据" />
      </div>
      <div v-else class="log-content">
        <pre>{{ logsText }}</pre>
      </div>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'
import * as logApi from '@/api/log'

const loading = ref(false)
const logs = ref<string[]>([])
const lineCount = ref(200)
const logFiles = ref<string[]>([])
const selectedFile = ref('')

const logsText = computed(() => logs.value.join('\n'))

async function loadFiles() {
  try {
    const res = await logApi.listLogFiles()
    if (res.code === 200) {
      logFiles.value = res.data || []
      if (!selectedFile.value && res.data && res.data.length > 0) {
        selectedFile.value = res.data[0]
      }
    } else {
      Message.error(res.message || '加载日志文件列表失败')
    }
  } catch (error) {
    Message.error('加载日志文件列表失败，请稍后重试')
  }
}

async function loadLogs() {
  loading.value = true
  try {
    const res = await logApi.listLatestLogs(lineCount.value, selectedFile.value || undefined)
    if (res.code === 200) {
      logs.value = res.data || []
    } else {
      Message.error(res.message || '加载日志失败')
    }
  } catch (error) {
    Message.error('加载系统日志失败，请稍后重试')
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  await loadFiles()
  await loadLogs()
})
</script>

<style scoped>
.log-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.title-section h2 {
  font-size: 20px;
  color: #1d2129;
  margin-bottom: 4px;
}

.subtitle {
  font-size: 14px;
  color: #86909c;
}

.actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.lines-label {
  font-size: 14px;
  color: #4e5969;
}

.log-content {
  max-height: 600px;
  background-color: #1e1e1e;
  border-radius: 8px;
  padding: 12px;
  overflow: auto;
}

.log-content pre {
  margin: 0;
  font-family: Consolas, Menlo, Monaco, 'Courier New', monospace;
  font-size: 12px;
  color: #d4d4d4;
  line-height: 1.5;
  white-space: pre-wrap;
}

.empty-block {
  padding: 40px 0;
}
</style>
