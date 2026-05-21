package com.bryan.system.service.system

import com.bryan.system.exception.BusinessException
import org.slf4j.LoggerFactory
import org.springframework.beans.factory.annotation.Value
import org.springframework.stereotype.Service
import java.nio.charset.StandardCharsets
import java.nio.file.Files
import java.nio.file.Path
import java.nio.file.Paths
import java.util.zip.GZIPInputStream

@Service
class LogService {
    private val log = LoggerFactory.getLogger(javaClass)

    @Value("\${logging.file.name:logs/platform.log}")
    private lateinit var logFileName: String

    fun listLatestLogs(maxLines: Int): List<String> = listLatestLogs(null, maxLines)

    fun listLatestLogs(fileName: String?, maxLines: Int): List<String> {
        val limit = maxLines.coerceIn(1, 2000)
        val path = resolveLogPath(fileName)
        if (!Files.exists(path)) throw BusinessException("日志文件不存在，请检查日志配置")
        val lines = readAllLines(path)
        return lines.takeLast(limit)
    }

    fun listLogFiles(): List<String> {
        val defaultLogPath = resolveDefaultLogPath()
        val dir = logsDirectory()
        if (!Files.isDirectory(dir)) return listOf(defaultLogPath.fileName.toString())
        return runCatching {
            Files.list(dir).use { stream ->
                stream.filter(Files::isRegularFile)
                    .map { it.fileName.toString() }
                    .filter { it.endsWith(".log") || it.endsWith(".gz") }
                    .sorted()
                    .toList()
                    .ifEmpty { listOf(defaultLogPath.fileName.toString()) }
            }
        }.getOrElse {
            log.warn("List log files failed: {}", dir.toAbsolutePath(), it)
            listOf(defaultLogPath.fileName.toString())
        }
    }

    private fun resolveLogPath(fileName: String?): Path {
        if (fileName.isNullOrBlank()) return resolveDefaultLogPath()
        val dir = logsDirectory()
        val path = dir.resolve(fileName).normalize()
        if (!path.startsWith(dir)) throw BusinessException("非法的日志文件路径")
        return path
    }

    private fun resolveDefaultLogPath(): Path {
        var path = Paths.get(logFileName)
        if (!path.isAbsolute) path = Paths.get(System.getProperty("user.dir")).resolve(path).normalize()
        return path
    }

    private fun logsDirectory(): Path = resolveDefaultLogPath().parent ?: Paths.get(System.getProperty("user.dir"))

    private fun readAllLines(path: Path): List<String> =
        if (path.fileName.toString().endsWith(".gz")) {
            GZIPInputStream(Files.newInputStream(path)).bufferedReader(StandardCharsets.UTF_8).useLines { it.toList() }
        } else {
            Files.readAllLines(path, StandardCharsets.UTF_8)
        }
}
