package com.bryan.system.controller.auth;

import com.bryan.system.domain.converter.UserConverter;
import com.bryan.system.domain.entity.SysUser;
import com.bryan.system.domain.entity.UserProfile;
import com.bryan.system.domain.request.auth.LoginRequest;
import com.bryan.system.domain.request.auth.RegisterRequest;
import com.bryan.system.domain.request.user.ChangePasswordRequest;
import com.bryan.system.domain.response.Result;
import com.bryan.system.domain.vo.UserVO;
import com.bryan.system.service.auth.AuthService;
import com.bryan.system.service.user.UserProfileService;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.security.access.prepost.PreAuthorize;
import org.springframework.security.core.annotation.AuthenticationPrincipal;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.*;

/**
 * 控制器：认证与授权接口
 * 提供用户注册、登录、当前用户信息获取及 Token 校验等接口。
 *
 * @author Bryan Long
 */
@Validated
@RestController
@RequestMapping("/api/auth")
@RequiredArgsConstructor
public class AuthController {

    private final AuthService authService;
    private final UserProfileService userProfileService;

    /**
     * 用户注册接口
     *
     * @param registerRequest 用户注册信息，包括用户名、密码、邮箱等
     * @return 注册成功的用户对象封装在统一响应结构中
     */
    @PostMapping("/register")
    @PreAuthorize("permitAll()")
    public Result<UserVO> register(@RequestBody @Valid RegisterRequest registerRequest) {
        // 1. 调用注册服务，创建新用户
        SysUser registered = authService.register(registerRequest);

        // 2. 初始化 UserProfile
        UserProfile userProfile = UserProfile.builder().userId(registered.getId()).build();
        UserProfile profile = userProfileService.createUserProfile(userProfile);

        // 2. 返回注册结果
        return Result.success(UserConverter.toUserVO(registered));
    }

    /**
     * 用户登录接口
     *
     * @param loginRequest 用户登录信息，包括用户名和密码
     * @return 登录成功后生成的 JWT Token 字符串封装在统一响应结构中
     */
    @PostMapping("/login")
    @PreAuthorize("permitAll()")
    public Result<String> login(@RequestBody @Valid LoginRequest loginRequest) {
        // 1. 调用认证服务进行登录验证
        String token = authService.login(loginRequest);

        // 2. 返回生成的 JWT Token
        return Result.success(token);
    }

    /**
     * 获取当前认证用户信息
     *
     * @return 当前登录用户的 VO 对象封装在统一响应结构中
     */
    @GetMapping("/me")
    @PreAuthorize("isAuthenticated()")
    public Result<UserVO> getCurrentUser() {
        // 1. 从 Spring Security 上下文中获取当前登录用户
        SysUser currentSysUser = authService.getCurrentUser();

        // 2. 使用 UserConverter 转换为DTO，排除敏感信息
        UserVO userVO = UserConverter.toUserVO(currentSysUser);

        // 3. 返回用户信息
        return Result.success(userVO);
    }

    /**
     * 验证 Token 合法性及用户状态
     *
     * @param token JWT Token 字符串
     * @param userDetails 当前已认证的用户信息（通过 Spring Security 注入）
     * @return 验证结果字符串封装在统一响应结构中
     */
    @GetMapping("/validate")
    @PreAuthorize("permitAll()")
    public Result<String> validate(
            @RequestParam String token,
            @AuthenticationPrincipal UserDetails userDetails) {
        // 1. 验证 Token 格式与签名是否合法
        if (!authService.validateToken(token)) {
            return Result.success("Invalid token");
        }

        // 2. 校验用户账户状态（仅在已认证情况下执行）
        if (userDetails != null) {
            // 2.1 检查账户是否被锁定
            if (!userDetails.isAccountNonLocked()) {
                return Result.success("Account locked");
            }

            // 2.2 检查账户是否过期
            if (!userDetails.isAccountNonExpired()) {
                return Result.success("Account expired");
            }

            // 2.3 检查账户是否被禁用
            if (!userDetails.isEnabled()) {
                return Result.success("Account disabled");
            }
        }

        // 3. 校验通过
        return Result.success("Validation passed");
    }

    /**
     * 修改用户密码。
     * 管理员可修改任意用户密码，用户本人可修改自己的密码。
     *
     * @param changePasswordRequest 包含旧密码和新密码的请求体
     * @return 更新后的用户实体
     */
    @PutMapping("/password")
    @PreAuthorize("isAuthenticated()")
    public Result<UserVO> changePassword(
            @RequestBody @Valid ChangePasswordRequest changePasswordRequest) {
        // 1. 调用服务层执行密码修改
        SysUser updated = authService.changePassword(changePasswordRequest.getOldPassword(),
                changePasswordRequest.getNewPassword());
        return Result.success(UserConverter.toUserVO(updated));
    }

    /**
     * 注销用户。
     *
     * @return 注销结果
     */
    @DeleteMapping
    @PreAuthorize("isAuthenticated()")
    public Result<UserVO> deleteAccount() {
        // 1. 调用服务层执行密码修改
        return Result.success(UserConverter.toUserVO(authService.deleteAccount()));
    }

    /**
     * 退出登录
     *
     * @return boolean 是否退出登录
     */
    @GetMapping("/logout")
    @PreAuthorize("isAuthenticated()")
    public boolean logout() {
        return authService.logout();
    }
}
