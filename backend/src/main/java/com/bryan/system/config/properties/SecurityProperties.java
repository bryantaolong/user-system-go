package com.bryan.system.config.properties;

import lombok.Getter;
import lombok.Setter;
import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.stereotype.Component;

/**
 * 安全配置属性类
 * 用于从配置文件中读取安全相关配置，避免硬编码。
 *
 * @author Bryan Long
 */
@Setter
@Getter
@Component
@ConfigurationProperties(prefix = "security")
public class SecurityProperties {

    /**
     * 登录失败次数限额
     * 超过此次数将锁定账号
     * 默认 5 次
     */
    private Integer loginFailLimit = 5;

    /**
     * 登录失败计数重置时间（分钟）
     * 超过此时间未登录失败，计数将重置
     * 默认 30 分钟
     */
    private Integer loginFailResetMinutes = 30;

    /**
     * 账号锁定时间（分钟）
     * 锁定后需要等待此时间才能再次尝试登录
     * 默认 30 分钟
     */
    private Integer accountLockDurationMinutes = 30;
}
