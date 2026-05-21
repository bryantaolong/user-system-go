package com.bryan.system.domain.entity

import com.bryan.system.domain.enums.GenderEnum
import com.bryan.system.domain.enums.UserStatusEnum
import com.fasterxml.jackson.annotation.JsonIgnore
import org.springframework.security.core.GrantedAuthority
import org.springframework.security.core.authority.SimpleGrantedAuthority
import org.springframework.security.core.userdetails.UserDetails
import java.io.Serializable
import java.time.LocalDateTime

class SysUser(
    var id: Long? = null,
    username: String? = null,
    password: String? = null,
    var phone: String? = null,
    var email: String? = null,
    var status: UserStatusEnum? = null,
    var roles: String? = null,
    var lastLoginAt: LocalDateTime? = null,
    var lastLoginIp: String? = null,
    var lastLoginDevice: String? = null,
    var passwordResetAt: LocalDateTime? = null,
    var loginFailCount: Int? = null,
    var lockedAt: LocalDateTime? = null,
    var deleted: Int? = null,
    var version: Int? = null,
    var createdAt: LocalDateTime? = null,
    var updatedAt: LocalDateTime? = null,
    var createdBy: String? = null,
    var updatedBy: String? = null
) : Serializable, UserDetails {
    private var usernameValue: String? = username

    @JsonIgnore
    private var passwordValue: String? = password

    override fun getUsername(): String = usernameValue ?: ""
    fun setUsername(value: String?) { usernameValue = value }

    @JsonIgnore
    override fun getPassword(): String = passwordValue ?: ""
    fun setPassword(value: String?) { passwordValue = value }

    @JsonIgnore
    override fun getAuthorities(): MutableCollection<out GrantedAuthority> =
        roles?.takeIf { it.isNotBlank() }
            ?.split(",")
            ?.map { SimpleGrantedAuthority(it.trim()) }
            ?.toMutableList()
            ?: mutableListOf()

    @JsonIgnore
    override fun isAccountNonExpired() = true

    @JsonIgnore
    override fun isAccountNonLocked(): Boolean {
        if (status == UserStatusEnum.NORMAL) return true
        if (status == UserStatusEnum.LOCKED && lockedAt != null) {
            return LocalDateTime.now().isAfter(lockedAt!!.plusHours(1))
        }
        return false
    }

    @JsonIgnore
    override fun isCredentialsNonExpired() = true

    @JsonIgnore
    override fun isEnabled(): Boolean = status != UserStatusEnum.BANNED && deleted == 0
}

data class UserProfile(
    var userId: Long? = null,
    var realName: String? = null,
    var gender: GenderEnum? = null,
    var birthday: LocalDateTime? = null,
    var avatar: String? = null,
    var deleted: Int? = null,
    var version: Int? = null,
    var createdAt: LocalDateTime? = null,
    var updatedAt: LocalDateTime? = null,
    var createdBy: String? = null,
    var updatedBy: String? = null
)

data class UserRole(
    var id: Long? = null,
    var roleName: String? = null,
    var isDefault: Boolean? = null,
    var deleted: Int? = null,
    var version: Int? = null,
    var createdAt: LocalDateTime? = null,
    var updatedAt: LocalDateTime? = null,
    var createdBy: String? = null,
    var updatedBy: String? = null
) : Serializable
