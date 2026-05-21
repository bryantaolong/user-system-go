package com.bryan.system.controller.admin

import com.bryan.system.domain.response.Result
import com.bryan.system.service.system.LogService
import org.springframework.security.access.prepost.PreAuthorize
import org.springframework.validation.annotation.Validated
import org.springframework.web.bind.annotation.GetMapping
import org.springframework.web.bind.annotation.RequestMapping
import org.springframework.web.bind.annotation.RequestParam
import org.springframework.web.bind.annotation.RestController

@Validated
@RestController
@RequestMapping("/api/admin/logs")
class SystemLogController(private val logService: LogService) {
    @GetMapping
    @PreAuthorize("hasRole('ADMIN')")
    fun listLatestLogs(@RequestParam(defaultValue = "200") lines: Int, @RequestParam(required = false) file: String?): Result<List<String>> =
        Result.success(logService.listLatestLogs(file, lines))

    @GetMapping("/files")
    @PreAuthorize("hasRole('ADMIN')")
    fun listLogFiles(): Result<List<String>> = Result.success(logService.listLogFiles())
}
