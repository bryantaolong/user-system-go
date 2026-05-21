package com.bryan.system.domain.response

import com.bryan.system.domain.enums.HttpStatus
import java.io.Serializable

data class Result<T>(
    var code: Int = 0,
    var message: String? = null,
    var data: T? = null
) {
    companion object {
        fun <T> success(data: T?): Result<T> = Result(HttpStatus.SUCCESS.code, HttpStatus.SUCCESS.msg, data)
        fun <T> error(httpStatus: HttpStatus): Result<T> = Result(httpStatus.code, httpStatus.msg, null)
        fun <T> error(httpStatus: HttpStatus, message: String?): Result<T> = Result(httpStatus.code, message, null)
    }
}

data class PageResult<T>(
    var rows: List<T> = emptyList(),
    var total: Long = 0,
    var pageNum: Long = 1,
    var pageSize: Long = 10
) : Serializable {
    val pages: Long
        get() = if (total == 0L || pageSize == 0L) 0 else (total + pageSize - 1) / pageSize

    fun isEmpty(): Boolean = rows.isEmpty()

    companion object {
        fun <T> of(rows: List<T>, total: Long, pageNum: Long, pageSize: Long): PageResult<T> =
            PageResult(rows, total, pageNum, pageSize)
    }
}
