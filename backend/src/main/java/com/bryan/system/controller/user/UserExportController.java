package com.bryan.system.controller.user;

import com.bryan.system.domain.request.user.UserExportRequest;
import com.bryan.system.service.user.UserExportService;
import jakarta.servlet.http.HttpServletResponse;
import lombok.RequiredArgsConstructor;
import org.springframework.security.access.prepost.PreAuthorize;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

import java.io.IOException;

/**
 * UserExportController
 *
 * @author Bryan Long
 */
@Validated
@RestController
@RequestMapping("/api/users/export")
@RequiredArgsConstructor
public class UserExportController {

    private final UserExportService userExportService;

    /**
     * 导出所有用户数据为 Excel 文件。
     * <p>仅管理员可操作。</p>
     *
     * @param response HttpServletResponse
     */
    @GetMapping
    @PreAuthorize("hasRole('ADMIN')")
    public void exportAllUsers(HttpServletResponse response,
                               @RequestParam(defaultValue = "1") int pageNum,
                               @RequestParam(defaultValue = "1000") int pageSize) throws IOException {
        userExportService.exportAllUsers(new UserExportRequest(), response, pageNum, pageSize);
    }
}
