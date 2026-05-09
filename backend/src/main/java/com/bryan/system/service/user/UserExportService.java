package com.bryan.system.service.user;

import com.alibaba.excel.EasyExcel;
import com.alibaba.excel.ExcelWriter;
import com.alibaba.excel.write.handler.CellWriteHandler;
import com.alibaba.excel.write.metadata.WriteSheet;
import com.bryan.system.domain.converter.UserConverter;
import com.bryan.system.domain.entity.SysUser;
import com.bryan.system.domain.request.user.UserExportRequest;
import com.bryan.system.domain.vo.UserExportVO;
import com.bryan.system.exception.BusinessException;
import com.bryan.system.mapper.UserMapper;
import jakarta.servlet.http.HttpServletResponse;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.apache.poi.ss.usermodel.*;
import org.springframework.stereotype.Service;
import org.springframework.util.CollectionUtils;

import java.io.IOException;
import java.net.URLEncoder;
import java.nio.charset.StandardCharsets;
import java.util.List;
import java.util.Optional;

/**
 * 用户导出业务服务
 * 支持全量导出与按指定字段导出，分批写入防止内存溢出。
 *
 * @author Bryan Long
 */
@Slf4j
@Service
@RequiredArgsConstructor
public class UserExportService {

    private final UserMapper userMapper;

    /* -------------------- 公开导出入口 -------------------- */

    /**
     * 全量导出：所有字段
     *
     * @param exportRequest 文件名称、状态过滤条件
     * @param response      响应流
     */
    public void exportAllUsers(UserExportRequest exportRequest, HttpServletResponse response,
                               int pageNum, int pageSize) throws IOException {
        try {
            String fileName = Optional.ofNullable(exportRequest.getFileName()).orElse("用户数据全量导出");
            this.setupResponse(response, fileName);

            ExcelWriter excelWriter = EasyExcel.write(response.getOutputStream())
                    .head(UserExportVO.class)
                    .registerWriteHandler(new CustomCellWriteHandler())
                    .build();

            this.executeBatchExport(excelWriter,
                    EasyExcel.writerSheet("用户列表").build(),
                    exportRequest, pageNum, pageSize);
        } catch (IOException e) {
            throw new BusinessException("全量导出失败，请检查系统资源", e);
        }
    }

    /* -------------------- 私有辅助 -------------------- */

    /**
     * 分批查询 + 写入，防止内存溢出
     */
    private void executeBatchExport(ExcelWriter excelWriter,
                                    WriteSheet writeSheet,
                                    UserExportRequest exportRequest,
                                    int pageNum, int pageSize) throws IOException {
        int total    = 0;

        while (true) {
            int offset = (pageNum - 1) * pageSize;
            List<SysUser> records = userMapper.selectExportPage(offset, pageSize, exportRequest);
            if (CollectionUtils.isEmpty(records)) {
                break;
            }
            List<UserExportVO> vos = records.stream()
                    .map(UserConverter::toExportVO)
                    .toList();
            excelWriter.write(vos, writeSheet);
            total += vos.size();
            log.info("已导出 {} 条数据", total);
            pageNum++;
        }
        excelWriter.finish();
        log.info("导出完成，总计 {} 条数据", total);
    }

    /**
     * 设置响应头：文件名、编码、Content-Type
     */
    private void setupResponse(HttpServletResponse response, String fileName) throws IOException {
        String encoded = URLEncoder.encode(fileName, StandardCharsets.UTF_8).replace("+", "%20");
        response.setContentType("application/vnd.openxmlformats-officedocument.spreadsheetml.sheet");
        response.setCharacterEncoding("utf-8");
        response.setHeader("Content-disposition", "attachment;filename*=utf-8''" + encoded + ".xlsx");
    }

    /**
     * 自定义单元格样式：表头浅绿背景 + 内容居中
     */
    private static class CustomCellWriteHandler implements CellWriteHandler {
        @Override
        public void afterCellCreate(com.alibaba.excel.write.metadata.holder.WriteSheetHolder writeSheetHolder,
                                    com.alibaba.excel.write.metadata.holder.WriteTableHolder writeTableHolder,
                                    Cell cell,
                                    com.alibaba.excel.metadata.Head head,
                                    Integer relativeRowIndex,
                                    Boolean isHead) {
            Workbook wb = cell.getSheet().getWorkbook();
            if (isHead) {
                CellStyle style = wb.createCellStyle();
                style.setFillForegroundColor(IndexedColors.LIGHT_GREEN.getIndex());
                style.setFillPattern(FillPatternType.SOLID_FOREGROUND);
                style.setAlignment(HorizontalAlignment.CENTER);
                cell.setCellStyle(style);
            } else {
                CellStyle style = wb.createCellStyle();
                style.setAlignment(HorizontalAlignment.CENTER);
                cell.setCellStyle(style);
            }
        }
    }
}
