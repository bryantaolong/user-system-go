package com.bryan.system.domain.vo

import com.alibaba.excel.annotation.ExcelProperty
import com.bryan.system.domain.enums.GenderEnum
import com.bryan.system.domain.enums.UserStatusEnum
import java.time.LocalDateTime

data class UserVO(
    var id: Long? = null,
    var username: String? = null,
    var phone: String? = null,
    var email: String? = null,
    var status: UserStatusEnum? = null,
    var roles: String? = null,
    var lastLoginAt: LocalDateTime? = null,
    var lastLoginIp: String? = null,
    var lastLoginDevice: String? = null,
    var createdAt: LocalDateTime? = null,
    var updatedAt: LocalDateTime? = null
)

data class UserProfileVO(
    var userId: Long? = null,
    var username: String? = null,
    var phone: String? = null,
    var email: String? = null,
    var realName: String? = null,
    var gender: GenderEnum? = null,
    var birthday: LocalDateTime? = null,
    var avatar: String? = null
)

data class UserRoleOptionVO(
    var id: Long? = null,
    var roleName: String? = null
)

data class UserExportVO(
    @ExcelProperty("ID") var id: Long? = null,
    @ExcelProperty("Username") var username: String? = null,
    @ExcelProperty("Phone") var phone: String? = null,
    @ExcelProperty("Email") var email: String? = null,
    @ExcelProperty("Status") var status: String? = null,
    @ExcelProperty("Roles") var roles: String? = null,
    @ExcelProperty("Last Login At") var lastLoginAt: LocalDateTime? = null,
    @ExcelProperty("Last Login IP") var lastLoginIp: String? = null,
    @ExcelProperty("Last Login Device") var lastLoginDevice: String? = null,
    @ExcelProperty("Created At") var createdAt: LocalDateTime? = null,
    @ExcelProperty("Updated At") var updatedAt: LocalDateTime? = null
)
