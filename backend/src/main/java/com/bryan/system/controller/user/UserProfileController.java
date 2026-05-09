package com.bryan.system.controller.user;

import com.bryan.system.domain.converter.UserConverter;
import com.bryan.system.domain.dto.UserProfileUpdateDTO;
import com.bryan.system.domain.entity.SysUser;
import com.bryan.system.domain.entity.UserProfile;
import com.bryan.system.domain.enums.HttpStatus;
import com.bryan.system.domain.request.user.UserUpdateRequest;
import com.bryan.system.domain.response.Result;
import com.bryan.system.domain.vo.UserProfileVO;
import com.bryan.system.service.auth.AuthService;
import com.bryan.system.service.user.UserProfileService;
import com.bryan.system.service.user.UserService;
import lombok.RequiredArgsConstructor;
import org.springframework.security.access.prepost.PreAuthorize;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.multipart.MultipartFile;

/**
 * 用户资料控制器
 * 提供用户资料的查询、更新等接口。
 *
 * @author Bryan Long
 */
@Validated
@RestController
@RequestMapping("/api/user-profiles")
@RequiredArgsConstructor
public class UserProfileController {

    private final UserProfileService userProfileService;
    private final UserService userService;
    private final AuthService authService;

    /**
     * 上传当前用户头像
     *
     * @param file 头像文件
     * @return 统一响应结果，包含头像相对路径
     */
    @PostMapping("/avatar")
    @PreAuthorize("isAuthenticated()")
    public Result<String> uploadAvatar(@RequestParam("file") MultipartFile file) {
        if (file.isEmpty()) {
            return Result.error(HttpStatus.BAD_REQUEST, "上传文件不能为空");
        }
        Long userId = authService.getCurrentUserId();
        String avatarPath = userProfileService.updateAvatar(userId, file);
        return Result.success(avatarPath);
    }

    /**
     * 根据用户主键查询用户资料（公开访问，用于展示用户信息）
     *
     * @param userId 用户主键
     * @return 用户资料 VO
     */
    @GetMapping("/{userId}")
    @PreAuthorize("permitAll()")
    public Result<UserProfileVO> getUserProfileByUserId(@PathVariable Long userId) {
        UserProfile profile = userProfileService.getUserProfileByUserId(userId);
        SysUser user = userService.getUserById(userId);
        return Result.success(UserConverter.toUserProfileVO(user, profile));
    }

    /**
     * 根据真实姓名查询用户资料
     *
     * @param realName 真实姓名
     * @return 用户资料 VO
     */
    @GetMapping("/name/{realName}")
    @PreAuthorize("isAuthenticated()")
    public Result<UserProfileVO> getUserProfileByRealName(@PathVariable String realName) {
        UserProfile profile = userProfileService.getUserProfileByRealName(realName);
        SysUser user = userService.getUserById(profile.getUserId());
        return Result.success(UserConverter.toUserProfileVO(user, profile));
    }

    /**
     * 获取当前登录用户的资料
     *
     * @return 用户资料 VO，如果用户资料不存在则返回空VO
     */
    @GetMapping("/me")
    @PreAuthorize("isAuthenticated()")
    public Result<UserProfileVO> getCurrentUserProfile() {
        Long userId = authService.getCurrentUserId();
        UserProfile profile = userProfileService.getUserProfileByUserIdOrEmpty(userId);
        SysUser user = userService.getUserById(userId);
        return Result.success(UserConverter.toUserProfileVO(user, profile));
    }

    /**
     * 更新当前用户资料
     *
     * @param req 更新请求参数
     * @return 更新后的用户资料 VO
     */
    @PutMapping
    @PreAuthorize("isAuthenticated()")
    public Result<UserProfileVO> updateUserProfile(
            @RequestBody UserUpdateRequest req) {
        Long userId = authService.getCurrentUserId();
        UserProfileUpdateDTO dto = UserProfileUpdateDTO.builder()
                .realName(req.getRealName())
                .gender(req.getGender())
                .birthday(req.getBirthday())
                .avatar(req.getAvatar())
                .build();
        UserProfileVO vo = UserConverter.toUserProfileVO(authService.getCurrentUser(),
                userProfileService.updateUserProfile(userId, dto));
        return Result.success(vo);
    }
}
