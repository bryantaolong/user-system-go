package com.bryan.system.domain.converter

import com.bryan.system.domain.entity.SysUser
import com.bryan.system.domain.entity.UserProfile
import com.bryan.system.domain.vo.UserExportVO
import com.bryan.system.domain.vo.UserProfileVO
import com.bryan.system.domain.vo.UserVO

object UserConverter {
    @JvmStatic
    fun toUserVO(user: SysUser?): UserVO? {
        if (user == null) return null
        return UserVO(
            id = user.id,
            username = user.getUsername(),
            phone = user.phone,
            email = user.email,
            status = user.status,
            roles = user.roles,
            lastLoginAt = user.lastLoginAt,
            lastLoginIp = user.lastLoginIp,
            lastLoginDevice = user.lastLoginDevice,
            createdAt = user.createdAt,
            updatedAt = user.updatedAt
        )
    }

    @JvmStatic
    fun toUserProfileVO(user: SysUser?, profile: UserProfile?): UserProfileVO =
        UserProfileVO(
            userId = user?.id ?: profile?.userId,
            username = user?.getUsername(),
            phone = user?.phone,
            email = user?.email,
            realName = profile?.realName,
            gender = profile?.gender,
            birthday = profile?.birthday,
            avatar = profile?.avatar
        )

    @JvmStatic
    fun toExportVO(user: SysUser): UserExportVO =
        UserExportVO(
            id = user.id,
            username = user.getUsername(),
            phone = user.phone,
            email = user.email,
            status = user.status?.name,
            roles = user.roles,
            lastLoginAt = user.lastLoginAt,
            lastLoginIp = user.lastLoginIp,
            lastLoginDevice = user.lastLoginDevice,
            createdAt = user.createdAt,
            updatedAt = user.updatedAt
        )
}
