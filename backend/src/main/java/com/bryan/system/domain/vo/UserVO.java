package com.bryan.system.domain.vo;

import lombok.Builder;
import lombok.Data;

import java.time.LocalDateTime;

/**
 * 用户数据传输对象，防止敏感信息暴露
 *
 * @author Bryan
 */
@Data
@Builder
public class UserVO {

    private Long id;
    
    private String username;
    
    private String email;
    
    private String phone;
    
    private String status;

    private String roles;
    
    private LocalDateTime lastLoginAt;

    private String lastLoginIp;

    private String lastLoginDevice;

    private LocalDateTime createdAt;

}
