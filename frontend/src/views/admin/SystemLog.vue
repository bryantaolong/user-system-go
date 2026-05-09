<template>
  <div class="system-log">
    <el-card class="log-card">
      <div class="log-header">
        <div class="title-section">
          <h2>系统日志</h2>
          <p class="subtitle">查看后台应用运行日志，仅管理员可访问</p>
        </div>
        <div class="actions">
          <el-select
            v-model="selectedFile"
            placeholder="选择日志文件"
            size="small"
            class="file-select"
            @change="loadLogs"
          >
            <el-option
              v-for="file in logFiles"
              :key="file"
              :label="file"
              :value="file"
            />
          </el-select>
          <span class="lines-label">行数：</span>
          <el-input-number
            v-model="lineCount"
            :min="50"
            :max="2000"
            :step="50"
            size="small"
          />
          <el-button type="primary" :loading="loading" @click="loadLogs">
            刷新
          </el-button>
        </div>
      </div>

      <el-divider />

      <el-empty v-if="!loading && logs.length === 0" description="暂无日志数据" />

      <el-scrollbar v-else class="log-content">
        <pre><code>{{ logsText }}</code></pre>
      </el-scrollbar>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import * as logApi from '@/api/system/log'

const loading = ref(false)
const logs = ref<string[]>([])
const lineCount = ref(200)
const logFiles = ref<string[]>([])
const selectedFile = ref<string>('')

const logsText = computed(() => logs.value.join('\n'))

const loadFiles = async () => {
  try {
    const res = await logApi.listLogFiles()
    if (res.code === 200) {
      logFiles.value = res.data || []
      if (!selectedFile.value && logFiles.value.length > 0) {
        selectedFile.value = logFiles.value[0]
      }
    } else {
      ElMessage.error(res.message || '加载日志文件列表失败')
    }
  } catch (error) {
    console.error('加载日志文件列表失败:', error)
    ElMessage.error('加载日志文件列表失败，请稍后重试')
  }
}

const loadLogs = async () => {
  loading.value = true
  try {
    const res = await logApi.listLatestLogs(lineCount.value, selectedFile.value || undefined)
    if (res.code === 200) {
      logs.value = res.data || []
    } else {
      ElMessage.error(res.message || '加载日志失败')
    }
  } catch (error) {
    console.error('加载系统日志失败:', error)
    ElMessage.error('加载系统日志失败，请稍后重试')
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
.system-log {
  padding: 20px;
}

.log-card {
  border-radius: 12px;
}

.log-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
}

.title-section h2 {
  margin: 0 0 8px 0;
  font-size: 22px;
  color: #303133;
}

.subtitle {
  margin: 0;
  font-size: 13px;
  color: #909399;
}

.actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.file-select {
  min-width: 200px;
}

.lines-label {
  font-size: 13px;
  color: #606266;
}

.log-content {
  margin-top: 12px;
  max-height: 600px;
  background-color: #1e1e1e;
  border-radius: 8px;
  padding: 12px;
}

pre {
  margin: 0;
  font-family: Consolas, Menlo, Monaco, 'Courier New', monospace;
  font-size: 12px;
  color: #d4d4d4;
  line-height: 1.5;
  white-space: pre-wrap;
}
</style>
