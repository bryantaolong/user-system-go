package com.bryan.system.config;

import org.springframework.cache.annotation.EnableCaching;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.data.redis.connection.RedisConnectionFactory;
import org.springframework.data.redis.core.RedisTemplate;
import org.springframework.data.redis.core.StringRedisTemplate;
import org.springframework.data.redis.serializer.RedisSerializer;

/**
 * Redis 全局配置类
 * 统一配置 RedisTemplate 的序列化策略，确保 key/value 可读且跨语言兼容。
 *
 * @author Bryan Long
 */
@Configuration
@EnableCaching
public class RedisConfig {

    /**
     * 注册通用 RedisTemplate
     * 采用 String 序列化 key，JSON 序列化 value，支持任意对象存取。
     *
     * @param factory SpringBoot 自动配置的连接工厂
     * @return 配置完成的 RedisTemplate<String, Object>
     */
    @Bean
    public RedisTemplate<String, Object> redisTemplate(RedisConnectionFactory factory) {
        RedisTemplate<String, Object> template = new RedisTemplate<>();
        template.setConnectionFactory(factory);

        // Key 使用 String 序列化
        template.setKeySerializer(RedisSerializer.string());

        // Value 使用 JSON 序列化
        template.setValueSerializer(RedisSerializer.json());

        // Hash Key/Value 序列化
        template.setHashKeySerializer(RedisSerializer.string());
        template.setHashValueSerializer(RedisSerializer.json());

        template.afterPropertiesSet();
        return template;
    }

    /**
     * 注册专用 StringRedisTemplate
     * 仅处理 String 类型，性能更高，适合缓存简单字符串或计数器场景。
     *
     * @param factory SpringBoot 自动配置的连接工厂
     * @return 配置完成的 StringRedisTemplate
     */
    @Bean
    public StringRedisTemplate stringRedisTemplate(RedisConnectionFactory factory) {
        return new StringRedisTemplate(factory);
    }
}
