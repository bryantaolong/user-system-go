package com.bryan.system.util.math;

import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;

class ArithmeticUtilTest {

    @Test
    void shouldAddPrecisely() {
        assertEquals(0.3, ArithmeticUtil.add(0.1, 0.2), 1e-12);
    }

    @Test
    void shouldSubtractPrecisely() {
        assertEquals(0.2, ArithmeticUtil.sub(0.5, 0.3), 1e-12);
    }

    @Test
    void shouldMultiplyPrecisely() {
        assertEquals(0.02, ArithmeticUtil.mul(0.1, 0.2), 1e-12);
    }

    @Test
    void shouldDivideWithDefaultScale() {
        assertEquals(0.3333333333, ArithmeticUtil.div(1, 3), 1e-12);
    }

    @Test
    void shouldReturnZeroWhenDividendIsZero() {
        assertEquals(0.0, ArithmeticUtil.div(0, 3), 1e-12);
    }

    @Test
    void shouldThrowWhenScaleIsNegativeOnDivide() {
        assertThrows(IllegalArgumentException.class, () -> ArithmeticUtil.div(1, 3, -1));
    }

    @Test
    void shouldRoundHalfUp() {
        assertEquals(2.35, ArithmeticUtil.round(2.345, 2), 1e-12);
        assertEquals(2.34, ArithmeticUtil.round(2.344, 2), 1e-12);
    }

    @Test
    void shouldThrowWhenScaleIsNegativeOnRound() {
        assertThrows(IllegalArgumentException.class, () -> ArithmeticUtil.round(1.23, -1));
    }
}
