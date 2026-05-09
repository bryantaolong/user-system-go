package com.bryan.system.domain.converter;

import com.bryan.system.domain.entity.SysUser;
import com.bryan.system.domain.entity.UserProfile;
import com.bryan.system.domain.enums.user.UserStatusEnum;
import com.bryan.system.domain.vo.UserExportVO;
import com.bryan.system.domain.vo.UserProfileVO;
import com.bryan.system.domain.vo.UserVO;

/**
 * 用户实体与值对象转换器
 * 负责 SysUser、UserProfile 与各层级 VO 之间的相互转换。
 *
 * @author Bryan Long
 */
public class UserConverter {

    /**
     * 将用户实体与资料实体合并为用户资料 VO
     *
     * @param user    用户实体
     * @param profile 用户资料实体
     * @return 用户资料 VO；若两者均为 null 则返回 null
     */
    public static UserProfileVO toUserProfileVO(SysUser user, UserProfile profile) {
        if (user == null && profile == null) {
            return null;
        }

        assert user != null;
        return UserProfileVO.builder()
                .userId(user.getId())
                .username(user.getUsername())
                .phone(user.getPhone())
                .email(user.getEmail())
                .realName(profile.getRealName())
                .gender(profile.getGender())
                .birthday(profile.getBirthday())
                .avatar(profile.getAvatar())
                .build();
    }

    /**
     * 将用户实体转换为对外展示的用户 VO
     *
     * @param user 用户实体
     * @return 用户 VO；若入参为 null 则返回 null
     */
    public static UserVO toUserVO(SysUser user) {
        if (user == null) {
            return null;
        }

        return UserVO.builder()
                .id(user.getId())
                .username(user.getUsername())
                .email(user.getEmail())
                .phone(user.getPhone())
                .status(user.getStatus() != null ? user.getStatus().name() : null)
                .createdAt(user.getCreatedAt())
                .lastLoginAt(user.getLastLoginAt())
                .lastLoginIp(user.getLastLoginIp())
                .lastLoginDevice(user.getLastLoginDevice())
                .roles(user.getRoles())
                .build();
    }

    /**
     * 将用户实体转换为导出用的用户导出 VO
     *
     * @param user 用户实体
     * @return 用户导出 VO；若入参为 null 则返回 null
     */
    public static UserExportVO toExportVO(SysUser user) {
        if (user == null) {
            return null;
        }

        return UserExportVO.builder()
                .id(user.getId())
                .username(user.getUsername())
                .phone(user.getPhone())
                .email(user.getEmail())
                .status(convertStatus(user.getStatus()))
                .roles(user.getRoles())
                .lastLoginAt(user.getLastLoginAt())
                .lastLoginIp(user.getLastLoginIp())
                .lastLoginDevice(user.getLastLoginDevice())
                .passwordResetAt(user.getPasswordResetAt())
                .deleted(convertDeletedStatus(user.getDeleted()))
                .createdAt(user.getCreatedAt())
                .updatedAt(user.getUpdatedAt())
                .createdBy(user.getCreatedBy())
                .updatedBy(user.getUpdatedBy())
                .build();
    }

    /**
     * 将用户状态枚举转换为中文描述
     *
     * @param status 用户状态枚举
     * @return 中文描述
     */
    private static String convertStatus(UserStatusEnum status) {
        if (status == null) return "";
        return switch (status) {
            case NORMAL -> "正常";
            case BANNED -> "封禁";
            case LOCKED -> "锁定";
            // 如果以后再加枚举，保留 default 分支
            default -> "未知";
        };
    }

    /**
     * 将删除标记转换为中文描述
     *
     * @param deleted 删除标记（0：未删除；其他：已删除）
     * @return 中文描述
     */
    private static String convertDeletedStatus(Integer deleted) {
        if (deleted == null) return "";
        return deleted == 0 ? "未删除" : "已删除";
    }
}
