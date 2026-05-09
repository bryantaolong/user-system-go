package com.bryan.system.domain.request.user;

import com.bryan.system.domain.enums.user.GenderEnum;
import jakarta.validation.constraints.Email;
import jakarta.validation.constraints.Past;
import jakarta.validation.constraints.Pattern;
import jakarta.validation.constraints.Size;
import lombok.Getter;

import java.time.LocalDateTime;

/**
 * 用户更新请求对象
 *
 * @author Bryan Long
 */
@Getter
public class UserUpdateRequest {
    /**
     * 电话号码。
     * 在更新时是可选的，如果提供则更新，如果不提供则保持不变。
     */
    @Pattern(regexp = "^1[3-9]\\d{9}$", message = "手机号格式不正确")
    private String phone;

    /**
     * 邮箱。
     * 在更新时是可选的，如果提供则更新。会校验邮箱格式。
     */
    @Email(message = "邮箱格式不正确")
    private String email;

    /**
     * 真实姓名。
     * 在更新时是可选的，如果提供则更新，如果不提供则保持不变。
     */
    @Size(min = 2, max = 20, message = "用户名长度应在2-20个字符之间")
    private String realName;

    /**
     * 性别。
     * 在更新时是可选的，如果提供则更新，如果不提供则保持不变。
     */
    private GenderEnum gender;

    /**
     * 出生日期。
     * 在更新时是可选的，如果提供则更新。
     * 必须是过去的日期，不能是未来日期。
     */
    @Past(message = "出生日期必须是过去的日期")
    private LocalDateTime birthday;

    /**
     * 头像 URL。
     * 在更新时是可选的，如果提供则更新，如果不提供则保持不变。
     * 长度限制为 500字符，需符合URL格式。
     */
    @Size(max = 500, message = "头像URL长度不能超过500个字符")
    @Pattern(regexp = "^(https?://.*)?$", message = "头像URL格式不正确，需以http://或https://开头")
    private String avatar;
}
