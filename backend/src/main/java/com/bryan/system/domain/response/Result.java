package com.bryan.system.domain.response;

import com.bryan.system.domain.enums.HttpStatus;
import lombok.Data;
import lombok.NoArgsConstructor;

/**
 * 统一响应封装
 * 提供成功/失败快速构造方法，确保所有接口返回格式一致。
 *
 * @author Bryan Long
 */
@Data
@NoArgsConstructor
public class Result<T> {

    /**
     * 业务状态码
     * 200 表示成功，其余见 HttpStatus 枚举
     */
    private int code;

    /**
     * 响应描述
     */
    private String message;

    /**
     * 实际业务数据
     */
    private T data;

    /**
     * 快速构造成功响应
     *
     * @param data 业务数据
     * @param <T>  数据类型
     * @return 成功响应对象
     */
    public static <T> Result<T> success(T data) {
        return new Result<>(HttpStatus.SUCCESS.getCode(), HttpStatus.SUCCESS.getMsg(), data);
    }

    /**
     * 构造失败响应（使用枚举消息）
     *
     * @param httpStatus 业务错误码枚举
     * @param <T>        数据类型
     * @return 失败响应对象
     */
    public static <T> Result<T> error(HttpStatus httpStatus) {
        return new Result<>(httpStatus.getCode(), httpStatus.getMsg(), null);
    }

    /**
     * 构造失败响应（自定义消息）
     *
     * @param httpStatus 业务错误码枚举
     * @param message    自定义错误描述
     * @param <T>        数据类型
     * @return 失败响应对象
     */
    public static <T> Result<T> error(HttpStatus httpStatus, String message) {
        return new Result<>(httpStatus.getCode(), message, null);
    }

    /**
     * 全参私有构造
     * 强制调用静态工厂方法，保证语义明确
     *
     * @param code    状态码
     * @param message 描述
     * @param data    数据
     */
    private Result(int code, String message, T data) {
        this.code = code;
        this.message = message;
        this.data = data;
    }
}
