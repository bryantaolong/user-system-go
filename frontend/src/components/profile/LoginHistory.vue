<template>
  <a-table :data="history" style="width: 100%" :pagination="false">
    <a-table-column title="登录时间" data-index="loginTime" :width="180" />
    <a-table-column title="IP地址" data-index="ipAddress" :width="140" />
    <a-table-column title="登录地点" data-index="location" />
    <a-table-column title="设备信息" data-index="device" />
  </a-table>
</template>

<script setup lang="ts">
defineProps<{
  history: Array<{
    loginTime: string
    ipAddress: string
    location?: string
    device: string
  }>
}>()
</script>
