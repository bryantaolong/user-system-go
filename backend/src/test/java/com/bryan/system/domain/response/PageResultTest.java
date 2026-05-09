package com.bryan.system.domain.response;

import org.junit.jupiter.api.Test;

import java.util.List;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertTrue;

class PageResultTest {

    @Test
    void shouldComputePagesUsingCeilingDivision() {
        PageResult<String> pageResult = PageResult.of(List.of("a", "b"), 11, 1, 10);
        assertEquals(2, pageResult.getPages());
    }

    @Test
    void shouldReturnZeroPagesWhenTotalIsZero() {
        PageResult<String> pageResult = PageResult.of(List.of(), 0, 1, 10);
        assertEquals(0, pageResult.getPages());
    }

    @Test
    void shouldReturnZeroPagesWhenPageSizeIsZero() {
        PageResult<String> pageResult = PageResult.of(List.of("a"), 10, 1, 0);
        assertEquals(0, pageResult.getPages());
    }

    @Test
    void shouldTreatNullRowsAsEmpty() {
        PageResult<String> pageResult = PageResult.<String>builder()
                .rows(null)
                .total(1)
                .pageNum(1)
                .pageSize(10)
                .build();

        assertTrue(pageResult.isEmpty());
    }

    @Test
    void shouldTreatNonEmptyRowsAsNotEmpty() {
        PageResult<String> pageResult = PageResult.of(List.of("x"), 1, 1, 10);
        assertFalse(pageResult.isEmpty());
    }
}
