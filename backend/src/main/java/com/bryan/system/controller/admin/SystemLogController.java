package com.bryan.system.controller.admin;

import com.bryan.system.domain.response.Result;
import com.bryan.system.service.system.LogService;
import lombok.RequiredArgsConstructor;
import org.springframework.security.access.prepost.PreAuthorize;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

import java.util.List;

/**
 * 系统日志监控控制器：提供后台管理员查看应用日志的接口。
 */
@Validated
@RestController
@RequestMapping("/api/admin/logs")
@RequiredArgsConstructor
public class SystemLogController {

    private final LogService logService;

    /**
     * 获取最新的应用日志内容（按行）。
     *
     * @param lines 返回的最大行数，默认 200，最大 2000
     * @return 日志行列表
     */
    @GetMapping
    @PreAuthorize("hasRole('ADMIN')")
    public Result<List<String>> listLatestLogs(
            @RequestParam(defaultValue = "200") int lines,
            @RequestParam(required = false) String file) {
        return Result.success(logService.listLatestLogs(file, lines));
    }

    /**
     * 获取可用的日志文件列表。
     */
    @GetMapping("/files")
    @PreAuthorize("hasRole('ADMIN')")
    public Result<List<String>> listLogFiles() {
        return Result.success(logService.listLogFiles());
    }
}
