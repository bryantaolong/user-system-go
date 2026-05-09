package com.bryan.system.exception;

/**
 * 乐观锁冲突异常
 * 当数据版本号不匹配时抛出，表示数据已被其他事务修改。
 *
 * @author Bryan Long
 */
public class OptimisticLockException extends RuntimeException {

    public OptimisticLockException(String message) {
        super(message);
    }

    public OptimisticLockException(String message, Throwable cause) {
        super(message, cause);
    }
}
