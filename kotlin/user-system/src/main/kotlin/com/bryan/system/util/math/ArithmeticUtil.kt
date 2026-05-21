package com.bryan.system.util.math

import java.math.BigDecimal
import java.math.RoundingMode

object ArithmeticUtil {
    @JvmStatic
    fun add(v1: Double, v2: Double): Double = BigDecimal.valueOf(v1).add(BigDecimal.valueOf(v2)).toDouble()

    @JvmStatic
    fun sub(v1: Double, v2: Double): Double = BigDecimal.valueOf(v1).subtract(BigDecimal.valueOf(v2)).toDouble()

    @JvmStatic
    fun mul(v1: Double, v2: Double): Double = BigDecimal.valueOf(v1).multiply(BigDecimal.valueOf(v2)).toDouble()

    @JvmStatic
    fun div(v1: Double, v2: Double, scale: Int = 2): Double =
        BigDecimal.valueOf(v1).divide(BigDecimal.valueOf(v2), scale, RoundingMode.HALF_UP).toDouble()
}
