<template>
  <el-table :data="history" style="width: 100%">
    <el-table-column prop="loginTime" label="登录时间" width="180"/>
    <el-table-column prop="ipAddress" label="IP地址" width="140"/>
    <el-table-column prop="location" label="登录地点"/>
    <el-table-column prop="device" label="设备信息"/>
  </el-table>
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

