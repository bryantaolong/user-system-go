package com.bryan.system.exception

open class BusinessException(message: String?, cause: Throwable? = null) : RuntimeException(message, cause)
class OptimisticLockException(message: String?, cause: Throwable? = null) : RuntimeException(message, cause)
class ResourceNotFoundException(message: String?, cause: Throwable? = null) : RuntimeException(message, cause)
class UnauthorizedException(message: String?, cause: Throwable? = null) : RuntimeException(message, cause)
