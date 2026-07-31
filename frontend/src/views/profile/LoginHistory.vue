<template>
  <a-table :columns="columns" :data="history" style="width: 100%">
    <template #render-column="{ column, record }">
      <template v-if="column.key === 'loginTime'">
        {{ record.loginTime }}
      </template>
      <template v-else-if="column.key === 'ipAddress'">
        {{ record.ipAddress }}
      </template>
      <template v-else-if="column.key === 'location'">
        {{ record.location || '-' }}
      </template>
      <template v-else-if="column.key === 'device'">
        {{ record.device }}
      </template>
    </template>
  </a-table>
</template>

<script setup lang="ts">
defineProps<{
  history: Array<{ loginTime: string; ipAddress: string; location?: string; device: string }>
}>()

const columns = [
  { title: '登录时间', dataIndex: 'loginTime', key: 'loginTime', width: 180 },
  { title: 'IP地址', dataIndex: 'ipAddress', key: 'ipAddress', width: 140 },
  { title: '登录地点', dataIndex: 'location', key: 'location' },
  { title: '设备信息', dataIndex: 'device', key: 'device' },
]
</script>
