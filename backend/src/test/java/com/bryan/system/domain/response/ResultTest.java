package com.bryan.system.domain.response;

import com.bryan.system.domain.enums.HttpStatus;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNull;

class ResultTest {

    @Test
    void shouldCreateSuccessResult() {
        Result<String> result = Result.success("payload");

        assertEquals(HttpStatus.SUCCESS.getCode(), result.getCode());
        assertEquals(HttpStatus.SUCCESS.getMsg(), result.getMessage());
        assertEquals("payload", result.getData());
    }

    @Test
    void shouldCreateErrorResultFromStatus() {
        Result<Void> result = Result.error(HttpStatus.BAD_REQUEST);

        assertEquals(HttpStatus.BAD_REQUEST.getCode(), result.getCode());
        assertEquals(HttpStatus.BAD_REQUEST.getMsg(), result.getMessage());
        assertNull(result.getData());
    }

    @Test
    void shouldCreateErrorResultWithCustomMessage() {
        Result<Void> result = Result.error(HttpStatus.BAD_REQUEST, "custom");

        assertEquals(HttpStatus.BAD_REQUEST.getCode(), result.getCode());
        assertEquals("custom", result.getMessage());
        assertNull(result.getData());
    }
}
