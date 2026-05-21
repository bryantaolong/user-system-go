package com.bryan.system.controller.user

import com.bryan.system.domain.converter.UserConverter
import com.bryan.system.domain.dto.UserProfileUpdateDTO
import com.bryan.system.domain.dto.UserUpdateDTO
import com.bryan.system.domain.entity.SysUser
import com.bryan.system.domain.entity.UserProfile
import com.bryan.system.domain.enums.HttpStatus
import com.bryan.system.domain.request.user.ChangePasswordRequest
import com.bryan.system.domain.request.user.ChangeRoleRequest
import com.bryan.system.domain.request.user.UserCreateRequest
import com.bryan.system.domain.request.user.UserExportRequest
import com.bryan.system.domain.request.user.UserSearchRequest
import com.bryan.system.domain.request.user.UserUpdateRequest
import com.bryan.system.domain.response.PageResult
import com.bryan.system.domain.response.Result
import com.bryan.system.domain.vo.UserProfileVO
import com.bryan.system.domain.vo.UserRoleOptionVO
import com.bryan.system.service.auth.AuthService
import com.bryan.system.service.user.UserExportService
import com.bryan.system.service.user.UserProfileService
import com.bryan.system.service.user.UserRoleService
import com.bryan.system.service.user.UserService
import jakarta.servlet.http.HttpServletResponse
import jakarta.validation.Valid
import org.springframework.security.access.prepost.PreAuthorize
import org.springframework.validation.annotation.Validated
import org.springframework.web.bind.annotation.DeleteMapping
import org.springframework.web.bind.annotation.GetMapping
import org.springframework.web.bind.annotation.PathVariable
import org.springframework.web.bind.annotation.PostMapping
import org.springframework.web.bind.annotation.PutMapping
import org.springframework.web.bind.annotation.RequestBody
import org.springframework.web.bind.annotation.RequestMapping
import org.springframework.web.bind.annotation.RequestParam
import org.springframework.web.bind.annotation.RestController
import org.springframework.web.multipart.MultipartFile

@Validated
@RestController
@RequestMapping("/api/users")
class UserController(
    private val userService: UserService,
    private val userProfileService: UserProfileService
) {
    @PostMapping
    @PreAuthorize("hasRole('ADMIN')")
    fun createUser(@RequestBody @Valid req: UserCreateRequest): Result<SysUser> {
        val created = userService.createUser(req)
        userProfileService.createUserProfile(UserProfile(userId = created.id))
        return Result.success(created)
    }

    @GetMapping
    @PreAuthorize("hasRole('ADMIN')")
    fun listUsers(@RequestParam(defaultValue = "1") pageNum: Int, @RequestParam(defaultValue = "10") pageSize: Int): Result<PageResult<SysUser>> =
        Result.success(userService.getAllUsers(pageNum, pageSize))

    @GetMapping("/{userId}")
    @PreAuthorize("hasRole('ADMIN')")
    fun getUserById(@PathVariable userId: Long): Result<SysUser> = Result.success(userService.getUserById(userId))

    @GetMapping("/username/{username}")
    @PreAuthorize("hasRole('ADMIN')")
    fun getUserByUsername(@PathVariable username: String): Result<SysUser?> = Result.success(userService.getUserByUsername(username))

    @PostMapping("/search")
    @PreAuthorize("hasRole('ADMIN')")
    fun queryUsers(@RequestBody searchRequest: UserSearchRequest, @RequestParam(defaultValue = "1") pageNum: Int, @RequestParam(defaultValue = "10") pageSize: Int): Result<PageResult<SysUser>> =
        Result.success(userService.queryUsers(searchRequest, pageNum, pageSize))

    @PutMapping("/{userId}")
    @PreAuthorize("hasRole('ADMIN') or (#userId == authentication.principal.id)")
    fun updateUser(@PathVariable userId: Long, @RequestBody @Valid req: UserUpdateRequest): Result<SysUser> =
        Result.success(userService.updateUser(userId, UserUpdateDTO(phone = req.phone, email = req.email)))

    @PutMapping("/roles/{userId}")
    @PreAuthorize("hasRole('ADMIN')")
    fun changeRoleByIds(@PathVariable userId: Long, @RequestBody @Valid req: ChangeRoleRequest): Result<SysUser> =
        Result.success(userService.changeRoleByIds(userId, req))

    @PutMapping("/password/{userId}")
    @PreAuthorize("hasRole('ADMIN')")
    fun resetPassword(@PathVariable userId: Long, @RequestBody req: ChangePasswordRequest): Result<SysUser> =
        Result.success(userService.resetPassword(userId, req.newPassword))

    @PutMapping("/block/{userId}")
    @PreAuthorize("hasRole('ADMIN')")
    fun blockUser(@PathVariable userId: Long): Result<SysUser> = Result.success(userService.blockUser(userId))

    @PutMapping("/unblock/{userId}")
    @PreAuthorize("hasRole('ADMIN')")
    fun unblockUser(@PathVariable userId: Long): Result<SysUser> = Result.success(userService.unblockUser(userId))

    @DeleteMapping("/{userId}")
    @PreAuthorize("hasRole('ADMIN')")
    fun deleteUser(@PathVariable userId: Long): Result<Long> = Result.success(userService.deleteUser(userId))
}

@Validated
@RestController
@RequestMapping("/api/user-profiles")
class UserProfileController(
    private val userProfileService: UserProfileService,
    private val userService: UserService,
    private val authService: AuthService
) {
    @PostMapping("/avatar")
    @PreAuthorize("isAuthenticated()")
    fun uploadAvatar(@RequestParam("file") file: MultipartFile): Result<String> {
        if (file.isEmpty) return Result.error(HttpStatus.BAD_REQUEST, "上传文件不能为空")
        return Result.success(userProfileService.updateAvatar(authService.getCurrentUserId(), file))
    }

    @GetMapping("/{userId}")
    @PreAuthorize("permitAll()")
    fun getUserProfileByUserId(@PathVariable userId: Long): Result<UserProfileVO> =
        Result.success(UserConverter.toUserProfileVO(userService.getUserById(userId), userProfileService.getUserProfileByUserId(userId)))

    @GetMapping("/name/{realName}")
    @PreAuthorize("isAuthenticated()")
    fun getUserProfileByRealName(@PathVariable realName: String): Result<UserProfileVO> {
        val profile = userProfileService.getUserProfileByRealName(realName)
        return Result.success(UserConverter.toUserProfileVO(userService.getUserById(profile.userId!!), profile))
    }

    @GetMapping("/me")
    @PreAuthorize("isAuthenticated()")
    fun getCurrentUserProfile(): Result<UserProfileVO> {
        val userId = authService.getCurrentUserId()
        return Result.success(UserConverter.toUserProfileVO(userService.getUserById(userId), userProfileService.getUserProfileByUserIdOrEmpty(userId)))
    }

    @PutMapping
    @PreAuthorize("isAuthenticated()")
    fun updateUserProfile(@RequestBody req: UserUpdateRequest): Result<UserProfileVO> {
        val userId = authService.getCurrentUserId()
        val dto = UserProfileUpdateDTO(req.realName, req.gender, req.birthday, req.avatar)
        return Result.success(UserConverter.toUserProfileVO(authService.getCurrentUser(), userProfileService.updateUserProfile(userId, dto)))
    }
}

@Validated
@RestController
@RequestMapping("/api/user-roles")
class UserRoleController(private val userRoleService: UserRoleService) {
    @GetMapping
    @PreAuthorize("hasRole('ADMIN')")
    fun listRoles(): Result<List<UserRoleOptionVO>> =
        Result.success(userRoleService.listAll().map { UserRoleOptionVO(it.id, it.roleName) })
}

@Validated
@RestController
@RequestMapping("/api/users/export")
class UserExportController(private val userExportService: UserExportService) {
    @GetMapping
    @PreAuthorize("hasRole('ADMIN')")
    fun exportAllUsers(response: HttpServletResponse, @RequestParam(defaultValue = "1") pageNum: Int, @RequestParam(defaultValue = "1000") pageSize: Int) {
        userExportService.exportAllUsers(UserExportRequest(), response, pageNum, pageSize)
    }
}
