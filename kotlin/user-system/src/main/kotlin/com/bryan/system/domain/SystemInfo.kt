package com.bryan.system.domain

data class Cpu(
    var cpuNum: Int? = null,
    var total: Double? = null,
    var sys: Double? = null,
    var used: Double? = null,
    var wait: Double? = null,
    var free: Double? = null
)

data class Memory(
    var total: Double? = null,
    var used: Double? = null,
    var free: Double? = null,
    var usage: Double? = null
)

data class System(
    var computerName: String? = null,
    var computerIp: String? = null,
    var userDir: String? = null,
    var osName: String? = null,
    var osArch: String? = null
)
