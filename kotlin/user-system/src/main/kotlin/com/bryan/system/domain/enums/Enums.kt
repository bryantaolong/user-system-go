package com.bryan.system.domain.enums

enum class HttpStatus(val code: Int, val msg: String) {
    SUCCESS(200, "success"),
    BAD_REQUEST(400, "bad request"),
    UNAUTHORIZED(401, "unauthorized"),
    FORBIDDEN(403, "forbidden"),
    NOT_FOUND(404, "not found"),
    CONFLICT(409, "conflict"),
    INTERNAL_ERROR(500, "internal server error")
}

enum class GenderEnum {
    UNKNOWN, MALE, FEMALE
}

enum class UserStatusEnum {
    NORMAL, LOCKED, BANNED
}
