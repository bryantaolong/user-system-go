package com.bryan.system.config;

import org.mybatis.spring.annotation.MapperScan;
import org.springframework.context.annotation.Configuration;
import org.springframework.transaction.annotation.EnableTransactionManagement;

/**
 * MyBatis 全局配置类
 * 负责注册 Mapper 扫描路径、事务管理。
 *
 * @author Bryan Long
 */
@Configuration
@MapperScan("com.bryan.system.mapper")
@EnableTransactionManagement
public class MyBatisConfig {
}
