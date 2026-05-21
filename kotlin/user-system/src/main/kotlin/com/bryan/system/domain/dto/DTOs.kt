package com.bryan.system.domain.dto

import com.bryan.system.domain.enums.GenderEnum
import java.time.LocalDateTime

data class UserUpdateDTO(
    var phone: String? = null,
    var email: String? = null
)

data class UserProfileUpdateDTO(
    var realName: String? = null,
    var gender: GenderEnum? = null,
    var birthday: LocalDateTime? = null,
    var avatar: String? = null
)
