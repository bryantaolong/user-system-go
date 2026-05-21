package com.bryan.system.service.file

import org.slf4j.LoggerFactory
import org.springframework.beans.factory.annotation.Value
import org.springframework.stereotype.Service
import org.springframework.web.multipart.MultipartFile
import java.io.InputStream
import java.nio.file.Files
import java.nio.file.Path
import java.nio.file.Paths
import java.nio.file.StandardCopyOption

@Service
class LocalFileService {
    private val log = LoggerFactory.getLogger(javaClass)

    @Value("\${file.upload-dir}")
    private lateinit var uploadDir: String

    fun storeFile(file: MultipartFile, subDirectory: String): String {
        val uploadPath = Paths.get(uploadDir, subDirectory).toAbsolutePath().normalize()
        if (!Files.exists(uploadPath)) Files.createDirectories(uploadPath)
        val original = file.originalFilename?.takeIf { it.isNotBlank() } ?: throw java.io.IOException("文件名不能为空")
        val detected = detectContentType(file.inputStream)
        if (detected !in ALLOWED_CONTENT_TYPES) throw java.io.IOException("不支持的文件类型")
        val base = original.substringBeforeLast('.', original)
        val corrected = if (original.lowercase().endsWith(".png")) original else "$base.png"
        val fileName = "${System.currentTimeMillis()}_$corrected"
        Files.copy(file.inputStream, uploadPath.resolve(fileName), StandardCopyOption.REPLACE_EXISTING)
        return Paths.get(subDirectory, fileName).toString()
    }

    fun loadFileAsBytes(filePath: String): ByteArray =
        Files.readAllBytes(Paths.get(uploadDir, filePath).toAbsolutePath().normalize())

    fun deleteFile(filePath: String): Boolean = try {
        Files.deleteIfExists(Paths.get(uploadDir, filePath).toAbsolutePath().normalize())
    } catch (e: Exception) {
        log.error("Failed to delete file: {}", filePath, e)
        false
    }

    private fun detectContentType(input: InputStream): String? {
        val header = ByteArray(12)
        val read = input.read(header)
        if (read >= 8 && header.take(8).map { it.toInt().toByte() } == PNG_MAGIC.toList()) return "image/png"
        if (read >= 2 && header[0] == 0xFF.toByte() && header[1] == 0xD8.toByte()) return "image/jpeg"
        if (read >= 6 && String(header, 0, 6) == "GIF87a" || read >= 6 && String(header, 0, 6) == "GIF89a") return "image/gif"
        if (read >= 12 && String(header, 0, 4) == "RIFF" && String(header, 8, 4) == "WEBP") return "image/webp"
        return null
    }

    companion object {
        private val ALLOWED_CONTENT_TYPES = setOf("image/png", "image/jpeg", "image/gif", "image/webp")
        private val PNG_MAGIC = byteArrayOf(0x89.toByte(), 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A)
    }
}
