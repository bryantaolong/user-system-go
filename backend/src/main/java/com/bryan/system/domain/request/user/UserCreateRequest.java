package com.bryan.system.domain.request.user;

import jakarta.validation.constraints.Email;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.Pattern;
import jakarta.validation.constraints.Size;
import lombok.Data;

import java.util.List;

/**
 * UserCreateRequest
 *
 * @author Bryan Long
 */
@Data
public class UserCreateRequest {

    @NotBlank
    @Size(min = 2, max = 20, message = "用户名长度应在2-20个字符之间")
    private String username;

    @NotBlank
    @Size(min = 6, message = "密码至少6位")
    private String password;

    @Pattern(regexp = "^1[3-9]\\d{9}$", message = "手机号格式不正确")
    private String phone;

    @Email(message = "邮箱格式不正确")
    private String email;

    private List<Long> roleIds;
}
