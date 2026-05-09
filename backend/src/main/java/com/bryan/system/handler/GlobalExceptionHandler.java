package com.bryan.system.handler;

import com.bryan.system.domain.enums.HttpStatus;
import com.bryan.system.domain.response.Result;
import com.bryan.system.exception.BusinessException;
import com.bryan.system.exception.OptimisticLockException;
import com.bryan.system.exception.ResourceNotFoundException;
import com.bryan.system.exception.UnauthorizedException;
import jakarta.servlet.http.HttpServletRequest;
import lombok.extern.slf4j.Slf4j;
import org.mybatis.spring.MyBatisSystemException;
import org.springframework.dao.DataIntegrityViolationException;
import org.springframework.security.access.AccessDeniedException;
import org.springframework.validation.FieldError;
import org.springframework.web.bind.MethodArgumentNotValidException;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.bind.annotation.RestControllerAdvice;

import java.util.stream.Collectors;

/**
 * 全局异常处理器
 * 统一捕获所有控制器层异常，转换为标准 Result 响应，避免敏感信息泄露。
 *
 * @author Bryan Long
 */
@Slf4j
@RestControllerAdvice
public class GlobalExceptionHandler {

    /**
     * 处理 MyBatis 系统异常
     * 通常由 SQL 语法错误、数据库连接问题或映射错误引起
     *
     * @param request 当前请求
     * @param e       MyBatis 系统异常
     * @return 统一错误响应
     */
    @ExceptionHandler(MyBatisSystemException.class)
    public Result<String> handleMyBatisSystemException(HttpServletRequest request, MyBatisSystemException e) {
        log.error("请求URL: {}, MyBatis 系统异常: {}",
                request.getRequestURL(), e.getMessage(), e);
        return Result.error(HttpStatus.INTERNAL_ERROR, "数据库操作异常，请联系管理员");
    }

    /**
     * 处理权限拒绝异常
     * 当用户尝试访问其没有权限的资源时抛出
     *
     * @param request 当前请求
     * @param e       权限拒绝异常
     * @return 统一错误响应
     */
    @ExceptionHandler(AccessDeniedException.class)
    public Result<String> handleAccessDeniedException(HttpServletRequest request, AccessDeniedException e) {
        log.warn("请求URL: {}, 权限拒绝: {}",
                request.getRequestURL(), e.getMessage());
        return Result.error(HttpStatus.FORBIDDEN, "权限不足，无法访问此资源");
    }

    /**
     * 处理运行时异常兜底
     * 生产环境仅返回友好提示，不暴露堆栈
     *
     * @param request 当前请求
     * @param e       异常
     * @return 统一错误响应
     */
    @ExceptionHandler(RuntimeException.class)
    public Result<String> handleRuntimeException(HttpServletRequest request, RuntimeException e) {
        log.error("请求URL: {}, 业务异常类型: {}",
                request.getRequestURL(), e.getClass().getSimpleName());
        return Result.error(HttpStatus.INTERNAL_ERROR, "服务繁忙，请稍后重试");
    }

    /**
     * 处理参数校验异常
     * 提取全部字段错误信息并拼接返回
     *
     * @param e 校验异常
     * @return 统一错误响应
     */
    @ExceptionHandler(MethodArgumentNotValidException.class)
    public Result<String> handleValidationException(MethodArgumentNotValidException e) {
        String errorMsg = e.getBindingResult()
                .getFieldErrors()
                .stream()
                .map(FieldError::getDefaultMessage)
                .collect(Collectors.joining(", "));
        log.warn("参数校验失败: {}", errorMsg);
        return Result.error(HttpStatus.BAD_REQUEST, errorMsg);
    }

    /**
     * 处理资源不存在异常
     *
     * @param e 自定义 404 异常
     * @return 统一错误响应
     */
    @ExceptionHandler(ResourceNotFoundException.class)
    public Result<String> handleNotFoundException(ResourceNotFoundException e) {
        log.warn("资源不存在: {}", e.getMessage());
        return Result.error(HttpStatus.NOT_FOUND);
    }

    /**
     * 处理数据完整性异常
     *
     * @param e 数据完整性异常
     * @return 统一错误响应
     */
    @ExceptionHandler(DataIntegrityViolationException.class)
    public Result<String> handleDataIntegrityViolationException(DataIntegrityViolationException e) {
        log.error("数据完整性异常: {}", e.getMessage(), e);
        return Result.error(HttpStatus.INTERNAL_ERROR, "数据操作失败，请检查数据格式");
    }

    /**
     * 处理业务逻辑异常
     *
     * @param e 业务异常
     * @return 统一错误响应
     */
    @ExceptionHandler(BusinessException.class)
    public Result<String> handleBusinessException(BusinessException e) {
        log.warn("业务异常: {}", e.getMessage());
        // 业务异常消息可能是敏感信息，根据实际情况决定是否返回给客户端
        // 这里保留原始消息，因为业务异常通常包含用户友好的提示
        return Result.error(HttpStatus.INTERNAL_ERROR, e.getMessage());
    }

    /**
     * 处理乐观锁冲突异常
     * 当数据版本号不匹配时抛出，提示用户刷新后重试
     *
     * @param e 乐观锁冲突异常
     * @return 统一错误响应
     */
    @ExceptionHandler(OptimisticLockException.class)
    public Result<String> handleOptimisticLockException(OptimisticLockException e) {
        log.warn("乐观锁冲突: {}", e.getMessage());
        return Result.error(HttpStatus.CONFLICT, e.getMessage());
    }

    /**
     * 处理未授权异常
     *
     * @param e 授权异常
     * @return 统一错误响应
     */
    @ExceptionHandler(UnauthorizedException.class)
    public Result<String> handleUnauthorizedException(UnauthorizedException e) {
        log.warn("未授权访问: {}", e.getMessage());
        return Result.error(HttpStatus.UNAUTHORIZED, e.getMessage());
    }
}
