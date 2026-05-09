package com.bryan.system.domain;

import lombok.Getter;
import lombok.Setter;

/**
 * System 系统相关信息
 *
 * @author Bryan Long
 */
@Getter
@Setter
public class System {
    /**
     * 服务器名称
     */
    private String computerName;

    /**
     * 服务器 IP
     */
    private String computerIp;

    /**
     * 项目路径
     */
    private String userDirectory;

    /**
     * 操作系统
     */
    private String osName;

    /**
     * 系统架构
     */
    private String osArchitecture;

}
