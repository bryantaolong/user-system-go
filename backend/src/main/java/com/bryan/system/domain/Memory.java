package com.bryan.system.domain;

import com.bryan.system.util.math.ArithmeticUtil;

/**
 * Memory 內存相关信息
 *
 * @author Bryan Long
 */
public class Memory {
    /**
     * 内存总量
     */
    private double totalMemory;

    /**
     * 已用内存
     */
    private double usedMemory;

    /**
     * 剩余内存
     */
    private double remainingMemory;

    public double getTotalMemory() {
        return ArithmeticUtil.div(totalMemory, (1024 * 1024 * 1024), 2);
    }

    public void setTotalMemory(long totalMemory) {
        this.totalMemory = totalMemory;
    }

    public double getUsedMemory() {
        return ArithmeticUtil.div(usedMemory, (1024 * 1024 * 1024), 2);
    }

    public void setUsedMemory(long usedMemory) {
        this.usedMemory = usedMemory;
    }

    public double getRemainingMemory() {
        return ArithmeticUtil.div(remainingMemory, (1024 * 1024 * 1024), 2);
    }

    public void setRemainingMemory(long remainingMemory) {
        this.remainingMemory = remainingMemory;
    }

    public double getUsage() {
        return ArithmeticUtil.mul(ArithmeticUtil.div(usedMemory, remainingMemory, 4), 100);
    }
}
