package com.bryan.system.controller.user;

import com.bryan.system.domain.dto.UserUpdateDTO;
import com.bryan.system.domain.entity.SysUser;
import com.bryan.system.domain.entity.UserProfile;
import com.bryan.system.domain.request.user.*;
import com.bryan.system.domain.response.PageResult;
import com.bryan.system.domain.response.Result;
import com.bryan.system.service.user.UserProfileService;
import com.bryan.system.service.user.UserService;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.security.access.prepost.PreAuthorize;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.*;

/**
 * 用户控制器：提供用户相关的 RESTful API 接口。
 * 包括用户信息查询、更新、角色变更、密码修改、逻辑删除及用户数据导出等功能。
 * 依赖 Spring Security 进行权限控制，支持管理员和普通用户不同权限的接口访问。
 *
 * @author Bryan Long
 */
@Validated
@RestController
@RequestMapping("/api/users")
@RequiredArgsConstructor
public class UserController {

    private final UserService userService;
    private final UserProfileService userProfileService;

    /**
     * 创建用户
     *
     * @param req 用户创建请求
     * @return 创建后的用户实体
     */
    @PostMapping
    @PreAuthorize("hasRole('ADMIN')")
    public Result<SysUser> createUser(@RequestBody @Valid UserCreateRequest req) {
        // 1. 调用创建服务，创建新用户
        SysUser created = userService.createUser(req);

        // 2. 初始化 UserProfile
        UserProfile userProfile = UserProfile.builder().userId(created.getId()).build();
        UserProfile profile = userProfileService.createUserProfile(userProfile);

        return Result.success(created);
    }

    /**
     * 获取所有用户列表（不分页）。
     * 仅允许拥有 ADMIN 角色的用户访问。
     *
     * @return 包含所有用户数据的分页对象（目前不支持分页参数，建议后续优化）。
     */
    @GetMapping
    @PreAuthorize("hasRole('ADMIN')")
    public Result<PageResult<SysUser>> listUsers(
            @RequestParam(defaultValue = "1") int pageNum,
            @RequestParam(defaultValue = "10") int pageSize) {
        // 1. 调用服务层获取所有用户列表
        return Result.success(userService.getAllUsers(pageNum, pageSize));
    }

    /**
     * 根据用户 ID 查询用户信息。
     * 仅允许 ADMIN 角色访问。
     *
     * @param userId 目标用户 ID
     * @return 对应用户实体
     */
    @GetMapping("/{userId}")
    @PreAuthorize("hasRole('ADMIN')")
    public Result<SysUser> getUserById(@PathVariable Long userId) {
        // 1. 调用服务获取用户信息
        return Result.success(userService.getUserById(userId));
    }

    /**
     * 根据用户名查询用户信息。
     * 允许未认证用户访问，用于公开信息展示。
     *
     * @param username 用户名
     * @return 对应用户实体
     */
    @GetMapping("/username/{username}")
    @PreAuthorize("hasRole('ADMIN')")
    public Result<SysUser> getUserByUsername(@PathVariable String username) {
        // 1. 调用服务获取用户信息
        return Result.success(userService.getUserByUsername(username));
    }

    /**
     * 用户搜索接口，支持多条件模糊查询和分页。
     *
     * @param searchRequest 搜索条件
     * @return 用户分页结果
     */
    @PostMapping("/search")
    @PreAuthorize("hasRole('ADMIN')")
    public Result<PageResult<SysUser>> queryUsers(
            @RequestBody UserSearchRequest searchRequest,
            @RequestParam(defaultValue = "1") int pageNum,
            @RequestParam(defaultValue = "10") int pageSize) {
        PageResult<SysUser> page = userService.queryUsers(searchRequest, pageNum, pageSize);
        return Result.success(page);
    }

    /**
     * 更新用户基本信息。
     * 允许管理员更新任意用户信息，或用户本人更新自己的信息。
     *
     * @param userId             目标用户 ID
     * @param req  包含需要更新的信息（用户名、邮箱等）
     * @return 更新后的用户实体
     * @throws IllegalArgumentException 当权限校验失败时抛出
     */
    @PutMapping("/{userId}")
    @PreAuthorize("hasRole('ADMIN') or (#userId == authentication.principal.id)")
    public Result<SysUser> updateUser(
            @PathVariable Long userId,
            @RequestBody @Valid UserUpdateRequest req) {
        // 1. 调用服务更新用户信息
        UserUpdateDTO dto = UserUpdateDTO.builder()
                .phone(req.getPhone())
                .email(req.getEmail())
                .build();
        return Result.success(userService.updateUser(userId, dto));
    }

    /**
     * 修改用户角色。
     * 仅管理员可操作。
     *
     * @param userId 目标用户 ID
     * @param req  新角色字符串，逗号分隔（如 "ROLE_USER,ROLE_ADMIN"）
     * @return 更新后的用户实体
     */
    @PutMapping("/roles/{userId}")
    @PreAuthorize("hasRole('ADMIN')")
    public Result<SysUser> changeRoleByIds(@PathVariable Long userId,
                                           @Valid @RequestBody ChangeRoleRequest req) {
        return Result.success(userService.changeRoleByIds(userId, req));
    }

    /**
     * 强制修改用户密码。
     * 管理员可强制修改任意用户密码。
     *
     * @param userId                   目标用户 ID
     * @param changePasswordRequest    修改密码请求
     * @return 更新后的用户实体
     */
    @PutMapping("/password/{userId}")
    @PreAuthorize("hasRole('ADMIN')")
    public Result<SysUser> resetPassword(
            @PathVariable Long userId,
            @RequestBody ChangePasswordRequest changePasswordRequest) {
        // 1. 调用服务层执行密码修改
        return Result.success(userService.resetPassword(userId, changePasswordRequest.getNewPassword()));
    }

    /**
     * 封禁用户。
     * 仅管理员可操作。
     *
     * @param userId 目标用户 ID
     * @return 更新后的用户实体
     */
    @PutMapping("/block/{userId}")
    @PreAuthorize("hasRole('ADMIN')")
    public Result<SysUser> blockUser(
            @PathVariable Long userId) {
        // 1. 调用服务封禁用户
        return Result.success(userService.blockUser(userId));
    }

    /**
     * 解封用户。
     * 仅管理员可操作。
     *
     * @param userId 目标用户 ID
     * @return 更新后的用户实体
     */
    @PutMapping("/unblock/{userId}")
    @PreAuthorize("hasRole('ADMIN')")
    public Result<SysUser> unblockUser(
            @PathVariable Long userId) {
        // 1. 调用服务解封用户
        return Result.success(userService.unblockUser(userId));
    }

    /**
     * 删除用户（逻辑删除）。
     * 仅管理员可执行。
     *
     * @param userId 目标用户 ID
     * @return 被删除的用户实体
     */
    @DeleteMapping("/{userId}")
    @PreAuthorize("hasRole('ADMIN')")
    public Result<Long> deleteUser(@PathVariable Long userId) {
        // 1. 调用服务执行逻辑删除
        return Result.success(userService.deleteUser(userId));
    }
}
