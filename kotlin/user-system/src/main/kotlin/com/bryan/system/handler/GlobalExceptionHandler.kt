package com.bryan.system.handler

import com.bryan.system.domain.enums.HttpStatus
import com.bryan.system.domain.response.Result
import com.bryan.system.exception.BusinessException
import com.bryan.system.exception.OptimisticLockException
import com.bryan.system.exception.ResourceNotFoundException
import com.bryan.system.exception.UnauthorizedException
import jakarta.servlet.http.HttpServletRequest
import org.mybatis.spring.MyBatisSystemException
import org.slf4j.LoggerFactory
import org.springframework.dao.DataIntegrityViolationException
import org.springframework.security.access.AccessDeniedException
import org.springframework.validation.FieldError
import org.springframework.web.bind.MethodArgumentNotValidException
import org.springframework.web.bind.annotation.ExceptionHandler
import org.springframework.web.bind.annotation.RestControllerAdvice

@RestControllerAdvice
class GlobalExceptionHandler {
    private val log = LoggerFactory.getLogger(javaClass)

    @ExceptionHandler(MyBatisSystemException::class)
    fun handleMyBatisSystemException(request: HttpServletRequest, e: MyBatisSystemException): Result<String> {
        log.error("请求URL: {}, MyBatis 系统异常: {}", request.requestURL, e.message, e)
        return Result.error(HttpStatus.INTERNAL_ERROR, "数据库操作异常，请联系管理员")
    }

    @ExceptionHandler(AccessDeniedException::class)
    fun handleAccessDeniedException(request: HttpServletRequest, e: AccessDeniedException): Result<String> {
        log.warn("请求URL: {}, 权限拒绝: {}", request.requestURL, e.message)
        return Result.error(HttpStatus.FORBIDDEN, "权限不足，无法访问此资源")
    }

    @ExceptionHandler(MethodArgumentNotValidException::class)
    fun handleValidationException(e: MethodArgumentNotValidException): Result<String> {
        val msg = e.bindingResult.fieldErrors.joinToString(", ") { obj: FieldError -> obj.defaultMessage.orEmpty() }
        return Result.error(HttpStatus.BAD_REQUEST, msg)
    }

    @ExceptionHandler(ResourceNotFoundException::class)
    fun handleNotFoundException(e: ResourceNotFoundException): Result<String> = Result.error(HttpStatus.NOT_FOUND)

    @ExceptionHandler(DataIntegrityViolationException::class)
    fun handleDataIntegrityViolationException(e: DataIntegrityViolationException): Result<String> =
        Result.error(HttpStatus.INTERNAL_ERROR, "数据操作失败，请检查数据格式")

    @ExceptionHandler(BusinessException::class)
    fun handleBusinessException(e: BusinessException): Result<String> =
        Result.error(HttpStatus.INTERNAL_ERROR, e.message)

    @ExceptionHandler(OptimisticLockException::class)
    fun handleOptimisticLockException(e: OptimisticLockException): Result<String> =
        Result.error(HttpStatus.CONFLICT, e.message)

    @ExceptionHandler(UnauthorizedException::class)
    fun handleUnauthorizedException(e: UnauthorizedException): Result<String> =
        Result.error(HttpStatus.UNAUTHORIZED, e.message)

    @ExceptionHandler(RuntimeException::class)
    fun handleRuntimeException(request: HttpServletRequest, e: RuntimeException): Result<String> {
        log.error("请求URL: {}, 运行时异常类型: {}", request.requestURL, e.javaClass.simpleName)
        return Result.error(HttpStatus.INTERNAL_ERROR, "服务繁忙，请稍后重试")
    }
}
