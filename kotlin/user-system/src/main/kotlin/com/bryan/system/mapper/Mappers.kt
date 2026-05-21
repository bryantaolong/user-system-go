package com.bryan.system.mapper

import com.bryan.system.domain.entity.SysUser
import com.bryan.system.domain.entity.UserProfile
import com.bryan.system.domain.entity.UserRole
import com.bryan.system.domain.enums.UserStatusEnum
import com.bryan.system.domain.request.user.UserExportRequest
import com.bryan.system.domain.request.user.UserSearchRequest
import org.apache.ibatis.annotations.Mapper
import org.apache.ibatis.annotations.Param
import java.time.LocalDateTime

@Mapper
interface UserMapper {
    fun insert(user: SysUser): Int
    fun selectById(id: Long?): SysUser?
    fun selectByUsername(username: String?): SysUser?
    fun selectPage(@Param("offset") offset: Int, @Param("pageSize") pageSize: Int, @Param("req") search: UserSearchRequest?): List<SysUser>
    fun selectExportPage(@Param("offset") offset: Int, @Param("pageSize") pageSize: Int, @Param("export") export: UserExportRequest?): List<SysUser>
    fun selectByIdList(@Param("ids") ids: Collection<Long>): List<SysUser>
    fun selectByStatus(@Param("status") status: UserStatusEnum?): SysUser?
    fun count(@Param("req") search: UserSearchRequest?): Long
    fun update(user: SysUser): Int
    fun deleteById(@Param("id") id: Long?, @Param("updatedAt") updatedAt: LocalDateTime, @Param("updatedBy") updatedBy: String?): Int
}

@Mapper
interface UserProfileMapper {
    fun insert(record: UserProfile): Int
    fun selectByUserId(userId: Long?): UserProfile?
    fun selectByRealName(realName: String?): UserProfile?
    fun update(record: UserProfile): Int
}

@Mapper
interface UserRoleMapper {
    fun selectAll(): List<UserRole>
    fun selectOneByIsDefaultTrue(): UserRole?
    fun selectByIdList(@Param("ids") ids: Collection<Long>): List<UserRole>
}
