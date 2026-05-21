package com.bryan.system.service.user

import com.alibaba.excel.EasyExcel
import com.alibaba.excel.write.handler.CellWriteHandler
import com.alibaba.excel.write.metadata.holder.WriteSheetHolder
import com.alibaba.excel.write.metadata.holder.WriteTableHolder
import com.bryan.system.domain.converter.UserConverter
import com.bryan.system.domain.dto.UserProfileUpdateDTO
import com.bryan.system.domain.dto.UserUpdateDTO
import com.bryan.system.domain.entity.SysUser
import com.bryan.system.domain.entity.UserProfile
import com.bryan.system.domain.entity.UserRole
import com.bryan.system.domain.enums.UserStatusEnum
import com.bryan.system.domain.request.user.ChangeRoleRequest
import com.bryan.system.domain.request.user.UserCreateRequest
import com.bryan.system.domain.request.user.UserExportRequest
import com.bryan.system.domain.request.user.UserSearchRequest
import com.bryan.system.domain.response.PageResult
import com.bryan.system.domain.vo.UserExportVO
import com.bryan.system.exception.BusinessException
import com.bryan.system.exception.ResourceNotFoundException
import com.bryan.system.mapper.UserMapper
import com.bryan.system.mapper.UserProfileMapper
import com.bryan.system.mapper.UserRoleMapper
import com.bryan.system.service.file.LocalFileService
import com.bryan.system.util.jwt.JwtUtils
import jakarta.servlet.http.HttpServletResponse
import org.apache.poi.ss.usermodel.Cell
import org.apache.poi.ss.usermodel.FillPatternType
import org.apache.poi.ss.usermodel.HorizontalAlignment
import org.apache.poi.ss.usermodel.IndexedColors
import org.slf4j.LoggerFactory
import org.springframework.security.crypto.password.PasswordEncoder
import org.springframework.stereotype.Service
import org.springframework.transaction.annotation.Transactional
import org.springframework.web.multipart.MultipartFile
import java.net.URLEncoder
import java.nio.charset.StandardCharsets
import java.time.LocalDateTime

@Service
class UserRoleService(private val userRoleMapper: UserRoleMapper) {
    fun listAll(): List<UserRole> = userRoleMapper.selectAll()
    fun getDefaultRole(): UserRole? = userRoleMapper.selectOneByIsDefaultTrue()
    fun listByIds(ids: Collection<Long>): List<UserRole> = userRoleMapper.selectByIdList(ids)
}

@Service
class UserService(
    private val userMapper: UserMapper,
    private val passwordEncoder: PasswordEncoder,
    private val userRoleService: UserRoleService
) {
    private val log = LoggerFactory.getLogger(javaClass)

    @Transactional
    fun createUser(req: UserCreateRequest): SysUser {
        if (userMapper.selectByUsername(req.username) != null) throw BusinessException("用户名已存在")
        val roleIds = req.roleIds?.toMutableList() ?: mutableListOf()
        if (roleIds.isEmpty()) roleIds.add(userRoleService.getDefaultRole()?.id ?: throw BusinessException("系统未配置默认角色"))
        val roles = userRoleService.listByIds(roleIds)
        if (roles.size != roleIds.size) throw IllegalArgumentException("角色不存在")
        val user = SysUser(
            username = req.username,
            password = passwordEncoder.encode(req.password),
            phone = req.phone,
            email = req.email,
            roles = roles.joinToString(",") { it.roleName.orEmpty() },
            status = UserStatusEnum.NORMAL,
            loginFailCount = 0,
            passwordResetAt = LocalDateTime.now()
        )
        fillInsert(user)
        if (userMapper.insert(user) == 0) throw BusinessException("插入数据库失败")
        return user
    }

    fun getAllUsers(pageNum: Int, pageSize: Int): PageResult<SysUser> {
        val offset = (pageNum - 1) * pageSize
        return PageResult.of(userMapper.selectPage(offset, pageSize, null), userMapper.count(null), pageNum.toLong(), pageSize.toLong())
    }

    fun getUserById(userId: Long): SysUser =
        userMapper.selectById(userId) ?: throw ResourceNotFoundException("用户不存在")

    fun getUserByUsername(username: String): SysUser? = userMapper.selectByUsername(username)

    fun queryUsers(searchRequest: UserSearchRequest, pageNum: Int, pageSize: Int): PageResult<SysUser> {
        val offset = (pageNum - 1) * pageSize
        return PageResult.of(userMapper.selectPage(offset, pageSize, searchRequest), userMapper.count(searchRequest), pageNum.toLong(), pageSize.toLong())
    }

    fun getUsersByIds(userIds: List<Long>?): List<SysUser> = if (userIds.isNullOrEmpty()) emptyList() else userMapper.selectByIdList(userIds)
    fun existsById(userId: Long): Boolean = userMapper.selectById(userId) != null

    fun updateUser(userId: Long, dto: UserUpdateDTO): SysUser {
        val user = getUserById(userId)
        dto.phone?.let { user.phone = it }
        dto.email?.let { user.email = it }
        fillUpdate(user)
        if (userMapper.update(user) != 1) throw BusinessException("更新失败，可能数据已变更")
        return user
    }

    @Transactional
    fun changeRoleByIds(userId: Long, req: ChangeRoleRequest): SysUser {
        val ids = req.roleIds
        val roles = userRoleService.listByIds(ids)
        if (roles.size != ids.size) throw IllegalArgumentException("角色不存在")
        val user = getUserById(userId)
        user.roles = roles.joinToString(",") { it.roleName.orEmpty() }
        fillUpdate(user)
        userMapper.update(user)
        return user
    }

    fun resetPassword(userId: Long, newPassword: String?): SysUser {
        val user = getUserById(userId)
        user.setPassword(passwordEncoder.encode(newPassword))
        user.passwordResetAt = LocalDateTime.now()
        fillUpdate(user)
        userMapper.update(user)
        return user
    }

    fun blockUser(userId: Long): SysUser = setStatus(userId, UserStatusEnum.BANNED)
    fun unblockUser(userId: Long): SysUser = setStatus(userId, UserStatusEnum.NORMAL)

    private fun setStatus(userId: Long, status: UserStatusEnum): SysUser {
        val user = getUserById(userId)
        user.status = status
        fillUpdate(user)
        userMapper.update(user)
        return user
    }

    fun deleteUser(userId: Long): Long {
        if (userMapper.deleteById(userId, LocalDateTime.now(), JwtUtils.getCurrentUsername()) == 0) {
            throw ResourceNotFoundException("用户不存在或已被删除")
        }
        return userId
    }

    private fun fillInsert(user: SysUser) {
        val now = LocalDateTime.now()
        val operator = JwtUtils.getCurrentUserId().toString()
        user.deleted = 0
        user.version = 0
        user.createdAt = now
        user.updatedAt = now
        user.createdBy = operator
        user.updatedBy = operator
    }

    private fun fillUpdate(user: SysUser) {
        user.version = (user.version ?: 0) + 1
        user.updatedAt = LocalDateTime.now()
        user.updatedBy = JwtUtils.getCurrentUserId().toString()
    }
}

@Service
class UserProfileService(
    private val userProfileMapper: UserProfileMapper,
    private val localFileService: LocalFileService
) {
    private val log = LoggerFactory.getLogger(javaClass)

    fun createUserProfile(record: UserProfile): UserProfile {
        fillInsert(record)
        if (userProfileMapper.insert(record) <= 0) throw BusinessException("创建用户信息失败")
        return record
    }

    fun getUserProfileByUserId(userId: Long?): UserProfile =
        userProfileMapper.selectByUserId(userId) ?: throw ResourceNotFoundException("用户信息不存在")

    fun getUserProfileByUserIdOrEmpty(userId: Long?): UserProfile =
        userProfileMapper.selectByUserId(userId) ?: UserProfile(userId = userId)

    fun getUserProfileByRealName(realName: String): UserProfile =
        userProfileMapper.selectByRealName(realName) ?: throw ResourceNotFoundException("用户信息不存在")

    fun updateUserProfile(userId: Long, dto: UserProfileUpdateDTO): UserProfile {
        val profile = getUserProfileByUserId(userId)
        dto.realName?.let { profile.realName = it }
        dto.gender?.let { profile.gender = it }
        dto.birthday?.let { profile.birthday = it }
        dto.avatar?.let { profile.avatar = it }
        fillUpdate(profile)
        if (userProfileMapper.update(profile) == 0) throw BusinessException("用户信息更新失败")
        return profile
    }

    fun updateAvatar(userId: Long, file: MultipartFile): String {
        val profile = getUserProfileByUserId(userId)
        val avatarPath = localFileService.storeFile(file, "avatars")
        profile.avatar?.takeIf { it.isNotBlank() }?.let { localFileService.deleteFile(it) }
        profile.avatar = avatarPath
        fillUpdate(profile)
        if (userProfileMapper.update(profile) == 0) throw BusinessException("头像更新失败")
        return avatarPath
    }

    private fun fillInsert(record: UserProfile) {
        val now = LocalDateTime.now()
        val operator = runCatching { JwtUtils.getCurrentUserId().toString() }.getOrDefault(record.userId?.toString() ?: "SYSTEM")
        record.deleted = 0
        record.version = 0
        record.createdAt = now
        record.updatedAt = now
        record.createdBy = operator
        record.updatedBy = operator
    }

    private fun fillUpdate(profile: UserProfile) {
        profile.version = (profile.version ?: 0) + 1
        profile.updatedAt = LocalDateTime.now()
        profile.updatedBy = JwtUtils.getCurrentUserId().toString()
    }
}

@Service
class UserExportService(private val userMapper: UserMapper) {
    fun exportAllUsers(exportRequest: UserExportRequest, response: HttpServletResponse, pageNum: Int, pageSize: Int) {
        val fileName = exportRequest.fileName ?: "用户数据全量导出"
        val encoded = URLEncoder.encode(fileName, StandardCharsets.UTF_8).replace("+", "%20")
        response.contentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
        response.characterEncoding = "utf-8"
        response.setHeader("Content-disposition", "attachment;filename*=utf-8''$encoded.xlsx")

        val writer = EasyExcel.write(response.outputStream)
            .head(UserExportVO::class.java)
            .registerWriteHandler(CustomCellWriteHandler())
            .build()
        try {
            val sheet = EasyExcel.writerSheet("用户列表").build()
            var current = pageNum
            while (true) {
                val records = userMapper.selectExportPage((current - 1) * pageSize, pageSize, exportRequest)
                if (records.isEmpty()) break
                writer.write(records.map(UserConverter::toExportVO), sheet)
                current++
            }
        } finally {
            writer.finish()
        }
    }

    private class CustomCellWriteHandler : CellWriteHandler {
        override fun afterCellCreate(writeSheetHolder: WriteSheetHolder, writeTableHolder: WriteTableHolder?, cell: Cell, head: com.alibaba.excel.metadata.Head?, relativeRowIndex: Int?, isHead: Boolean?) {
            val style = cell.sheet.workbook.createCellStyle()
            if (isHead == true) {
                style.fillForegroundColor = IndexedColors.LIGHT_GREEN.index
                style.fillPattern = FillPatternType.SOLID_FOREGROUND
            }
            style.alignment = HorizontalAlignment.CENTER
            cell.cellStyle = style
        }
    }
}
