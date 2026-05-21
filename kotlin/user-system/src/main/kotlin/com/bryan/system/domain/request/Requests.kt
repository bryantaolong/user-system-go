package com.bryan.system.domain.request

import com.bryan.system.domain.enums.GenderEnum
import com.bryan.system.domain.enums.UserStatusEnum
import jakarta.validation.constraints.Email
import jakarta.validation.constraints.NotBlank
import jakarta.validation.constraints.NotEmpty
import jakarta.validation.constraints.Past
import jakarta.validation.constraints.Pattern
import jakarta.validation.constraints.Size
import java.time.LocalDateTime

data class LoginRequest(
    @field:NotBlank var username: String? = null,
    @field:NotBlank var password: String? = null
)

data class RegisterRequest(
    @field:NotBlank @field:Size(min = 3, max = 20) var username: String? = null,
    @field:NotBlank @field:Size(min = 6, max = 100) var password: String? = null,
    @field:Pattern(regexp = "^1[3-9]\\d{9}$") var phone: String? = null,
    @field:Email var email: String? = null
)

data class ChangePasswordRequest(
    var oldPassword: String? = null,
    @field:NotBlank @field:Size(min = 6, max = 100) var newPassword: String? = null
)

data class ChangeRoleRequest(
    @field:NotEmpty var roleIds: MutableList<Long> = mutableListOf()
)

data class UserCreateRequest(
    @field:NotBlank @field:Size(min = 3, max = 20) var username: String? = null,
    @field:NotBlank @field:Size(min = 6, max = 100) var password: String? = null,
    @field:Pattern(regexp = "^1[3-9]\\d{9}$") var phone: String? = null,
    @field:Email var email: String? = null,
    var roleIds: MutableList<Long>? = null
)

data class UserExportRequest(
    var fileName: String? = null,
    var status: UserStatusEnum? = null,
    var keyword: String? = null
)

data class UserSearchRequest(
    var keyword: String? = null,
    var username: String? = null,
    var phone: String? = null,
    var email: String? = null,
    var status: UserStatusEnum? = null,
    var role: String? = null
)

data class UserUpdateRequest(
    @field:Pattern(regexp = "^1[3-9]\\d{9}$") var phone: String? = null,
    @field:Email var email: String? = null,
    @field:Size(min = 2, max = 20) var realName: String? = null,
    var gender: GenderEnum? = null,
    @field:Past var birthday: LocalDateTime? = null,
    @field:Size(max = 500)
    @field:Pattern(regexp = "^(https?://.*)?$")
    var avatar: String? = null
)
