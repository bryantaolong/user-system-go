package com.bryan.system.domain.request.user;

import com.bryan.system.domain.enums.user.UserStatusEnum;
import jakarta.validation.constraints.Email;
import jakarta.validation.constraints.NotNull;
import jakarta.validation.constraints.Pattern;
import jakarta.validation.constraints.Size;
import lombok.Getter;

/**
 * UserSearchRequest
 *
 * @author Bryan Long
 * @version 1.0
 * @since 2025/7/26
 */
@Getter
public class UserSearchRequest {
    // 用户名校验
    @Size(min = 2, max = 20, message = "用户名长度应在2-20个字符之间")
    @Pattern(regexp = "^[a-zA-Z0-9_\\u4e00-\\u9fa5]+$", message = "用户名只能包含中文、字母、数字和下划线")
    private String username;

    // 手机号校验（更精确的中国手机号校验）
    @Pattern(regexp = "^(1[3-9]\\d{9})?$", message = "手机号格式不正确") // 允许空值
    private String phone;

    // 邮箱校验（更宽松的格式）
    @Email(message = "邮箱格式不正确")
    @Size(max = 100, message = "邮箱长度不能超过100个字符")
    private String email;

    @NotNull(message = "状态不能为空")
    private UserStatusEnum status;
}
