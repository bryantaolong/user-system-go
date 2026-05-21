package com.bryan.system.service.redis

import org.slf4j.LoggerFactory
import org.springframework.data.redis.core.RedisTemplate
import org.springframework.data.redis.core.StringRedisTemplate
import org.springframework.stereotype.Service
import java.time.Duration

@Service
class RedisStringService(private val stringRedisTemplate: StringRedisTemplate) {
    private val log = LoggerFactory.getLogger(javaClass)
    fun set(key: String, value: String): Boolean = runCatching { stringRedisTemplate.opsForValue().set(key, value) }.onFailure { log.error("Redis set failed: {}", key, it) }.isSuccess
    fun set(key: String, value: String, seconds: Long): Boolean = if (seconds > 0) runCatching { stringRedisTemplate.opsForValue().set(key, value, Duration.ofSeconds(seconds)) }.isSuccess else set(key, value)
    fun setExpire(key: String, seconds: Long): Boolean = seconds > 0 && stringRedisTemplate.expire(key, Duration.ofSeconds(seconds)) == true
    fun get(key: String): String? = runCatching { stringRedisTemplate.opsForValue().get(key) }.getOrNull()
    fun delete(key: String): Boolean = runCatching { stringRedisTemplate.delete(key) == true }.getOrDefault(false)
    fun hasKey(key: String): Boolean = runCatching { stringRedisTemplate.hasKey(key) == true }.getOrDefault(false)
    fun increment(key: String, delta: Long): Long? = runCatching { stringRedisTemplate.opsForValue().increment(key, delta) }.getOrNull()
    fun decrement(key: String, delta: Long): Long? = runCatching { stringRedisTemplate.opsForValue().decrement(key, delta) }.getOrNull()
}

@Service
class RedisHashService(private val redisTemplate: RedisTemplate<String, Any>) {
    fun set(key: String, value: Map<String, Any>): Boolean = runCatching { redisTemplate.opsForHash<String, Any>().putAll(key, value) }.isSuccess
    fun set(key: String, hashKey: String, value: Any): Boolean = runCatching { redisTemplate.opsForHash<String, Any>().put(key, hashKey, value) }.isSuccess
    fun get(key: String, hashKey: String): Any? = runCatching { redisTemplate.opsForHash<String, Any>().get(key, hashKey) }.getOrNull()
    fun delete(key: String, hashKey: String): Boolean = runCatching { redisTemplate.opsForHash<String, Any>().delete(key, hashKey) > 0 }.getOrDefault(false)
    fun hasKey(key: String, hashKey: String): Boolean = runCatching { redisTemplate.opsForHash<String, Any>().hasKey(key, hashKey) }.getOrDefault(false)
    fun getAll(key: String): Map<String, Any> = runCatching { redisTemplate.opsForHash<String, Any>().entries(key) }.getOrDefault(emptyMap())
    fun keys(key: String): Set<String> = runCatching { redisTemplate.opsForHash<String, Any>().keys(key) }.getOrDefault(emptySet())
    fun values(key: String): List<Any> = runCatching { redisTemplate.opsForHash<String, Any>().values(key) }.getOrDefault(emptyList())
    fun size(key: String): Long = runCatching { redisTemplate.opsForHash<String, Any>().size(key) }.getOrDefault(0)
}

@Service
class RedisListService(private val redisTemplate: RedisTemplate<String, Any>) {
    fun leftPush(key: String, value: Any): Boolean = runCatching { redisTemplate.opsForList().leftPush(key, value) }.isSuccess
    fun leftPushAll(key: String, vararg values: Any): Boolean = runCatching { redisTemplate.opsForList().leftPushAll(key, *values) }.isSuccess
    fun rightPush(key: String, value: Any): Boolean = runCatching { redisTemplate.opsForList().rightPush(key, value) }.isSuccess
    fun rightPushAll(key: String, vararg values: Any): Boolean = runCatching { redisTemplate.opsForList().rightPushAll(key, *values) }.isSuccess
    fun leftPop(key: String): Any? = runCatching { redisTemplate.opsForList().leftPop(key) }.getOrNull()
    fun rightPop(key: String): Any? = runCatching { redisTemplate.opsForList().rightPop(key) }.getOrNull()
    fun range(key: String, start: Long, end: Long): List<Any> = runCatching { redisTemplate.opsForList().range(key, start, end) ?: emptyList() }.getOrDefault(emptyList())
    fun size(key: String): Long = runCatching { redisTemplate.opsForList().size(key) ?: 0 }.getOrDefault(0)
    fun index(key: String, index: Long): Any? = runCatching { redisTemplate.opsForList().index(key, index) }.getOrNull()
    fun set(key: String, index: Long, value: Any): Boolean = runCatching { redisTemplate.opsForList().set(key, index, value) }.isSuccess
    fun remove(key: String, count: Long, value: Any): Long = runCatching { redisTemplate.opsForList().remove(key, count, value) ?: 0 }.getOrDefault(0)
    fun trim(key: String, start: Long, end: Long): Boolean = runCatching { redisTemplate.opsForList().trim(key, start, end) }.isSuccess
}
