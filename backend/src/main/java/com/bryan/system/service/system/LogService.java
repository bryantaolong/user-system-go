package com.bryan.system.service.system;

import com.bryan.system.exception.BusinessException;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStream;
import java.io.InputStreamReader;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import java.util.zip.GZIPInputStream;

/**
 * 日志读取服务
 * 提供读取应用日志文件（支持 .gz 压缩）及列表能力，供后台管理端使用。
 *
 * @author Bryan Long
 */
@Slf4j
@Service
public class LogService {

    /**
     * 日志文件路径，支持 Spring Boot 配置项 logging.file.name
     * 默认相对路径：logs/platform.log
     */
    @Value("${logging.file.name:logs/platform.log}")
    private String logFileName;

    /**
     * 读取默认日志文件最近 N 行
     *
     * @param maxLines 最大行数（1~2000）
     * @return 日志行列表（按时间正序）
     */
    public List<String> listLatestLogs(int maxLines) {
        return this.listLatestLogs(null, maxLines);
    }

    /**
     * 读取指定日志文件最近 N 行
     *
     * @param fileName 日志文件名（可选），为空时使用默认日志文件
     * @param maxLines 最大行数（1~2000）
     * @return 日志行列表（按时间正序）
     */
    public List<String> listLatestLogs(String fileName, int maxLines) {
        int limit = Math.max(1, Math.min(maxLines, 2000));
        Path path = this.resolveLogPath(fileName);

        if (!Files.exists(path)) {
            log.warn("日志文件不存在：{}", path.toAbsolutePath());
            throw new BusinessException("日志文件不存在，请检查日志配置");
        }

        try {
            List<String> allLines = this.readAllLines(path);
            if (allLines.isEmpty()) {
                return Collections.emptyList();
            }
            int fromIndex = Math.max(0, allLines.size() - limit);
            return allLines.subList(fromIndex, allLines.size());
        } catch (IOException e) {
            log.error("读取日志文件失败：{}", path.toAbsolutePath(), e);
            throw new BusinessException("读取日志文件失败，请稍后重试", e);
        }
    }

    /**
     * 列出日志目录下所有可用日志文件（.log / .gz）
     *
     * @return 文件名列表（仅文件名，不含路径）
     */
    public List<String> listLogFiles() {
        Path defaultLogPath = this.resolveDefaultLogPath();
        Path logDir = this.getLogsDirectory();

        if (!Files.exists(logDir) || !Files.isDirectory(logDir)) {
            return List.of(defaultLogPath.getFileName().toString());
        }

        try (var stream = Files.list(logDir)) {
            List<String> files = stream
                    .filter(Files::isRegularFile)
                    .map(Path::getFileName)
                    .map(Path::toString)
                    .filter(name -> name.endsWith(".log") || name.endsWith(".gz"))
                    .sorted()
                    .toList();
            return files.isEmpty() ? List.of(defaultLogPath.getFileName().toString()) : files;
        } catch (IOException e) {
            log.warn("列出日志文件失败：{}", logDir.toAbsolutePath(), e);
            return List.of(defaultLogPath.getFileName().toString());
        }
    }

    /* -------------------- 私有工具方法 -------------------- */

    /**
     * 根据文件名解析日志路径，并做路径穿越防护
     */
    private Path resolveLogPath(String fileName) {
        if (fileName == null || fileName.isBlank()) {
            return this.resolveDefaultLogPath();
        }
        Path logDir = this.getLogsDirectory();
        Path path = logDir.resolve(fileName).normalize();
        if (!path.startsWith(logDir)) {
            throw new BusinessException("非法的日志文件路径");
        }
        return path;
    }

    /**
     * 解析默认日志文件路径（支持相对路径）
     */
    private Path resolveDefaultLogPath() {
        Path path = Paths.get(logFileName);
        if (!path.isAbsolute()) {
            Path userDir = Paths.get(System.getProperty("user.dir"));
            path = userDir.resolve(path).normalize();
        }
        return path;
    }

    /**
     * 获取日志目录（默认日志文件的父目录）
     */
    private Path getLogsDirectory() {
        Path defaultLogPath = this.resolveDefaultLogPath();
        Path parent = defaultLogPath.getParent();
        return parent == null ? Paths.get(System.getProperty("user.dir")) : parent;
    }

    /**
     * 读取文件全部行，支持普通文本与 .gz 压缩
     */
    private List<String> readAllLines(Path path) throws IOException {
        String fileName = path.getFileName().toString();
        if (fileName.endsWith(".gz")) {
            List<String> lines = new ArrayList<>();
            try (InputStream in = Files.newInputStream(path);
                 GZIPInputStream gzip = new GZIPInputStream(in);
                 BufferedReader reader = new BufferedReader(new InputStreamReader(gzip, StandardCharsets.UTF_8))) {
                String line;
                while ((line = reader.readLine()) != null) {
                    lines.add(line);
                }
            }
            return lines;
        }
        return Files.readAllLines(path, StandardCharsets.UTF_8);
    }
}
