package com.bryan.system.service.file;

import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;
import org.springframework.web.multipart.MultipartFile;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.nio.file.StandardCopyOption;
import java.util.Set;

/**
 * 文件存储服务
 * 提供文件上传、删除与读取能力。
 */
@Slf4j
@Service
public class LocalFileService {

    // 允许的文件类型（MIME类型白名单）
    private static final Set<String> ALLOWED_CONTENT_TYPES = Set.of(
            "image/png",
            "image/jpeg",
            "image/gif",
            "image/webp"
    );

    // PNG 文件的魔数 (文件头签名)
    private static final byte[] PNG_MAGIC = new byte[]{(byte) 0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A};
    // JPEG 文件的魔数
    private static final byte[] JPEG_MAGIC = new byte[]{(byte) 0xFF, (byte) 0xD8};

    @Value("${file.upload-dir}")
    private String uploadDir;

    /**
     * 存储上传文件并返回相对路径
     *
     * @param file 上传文件
     * @param subDirectory 子目录
     * @return uploads 下的相对路径
     * @throws IOException 文件读写异常
     */
    public String storeFile(MultipartFile file, String subDirectory) throws IOException {
        // 构建上传目录的绝对路径，包括子目录
        Path uploadPath = Paths.get(uploadDir, subDirectory).toAbsolutePath().normalize();

        // 如果目录不存在，则创建
        if (!Files.exists(uploadPath)) {
            Files.createDirectories(uploadPath);
        }

        // 获取原始文件名
        String originalFilename = file.getOriginalFilename();
        if (originalFilename == null || originalFilename.isEmpty()) {
            throw new IOException("文件名不能为空");
        }

        // 验证文件内容类型（通过魔数检测）
        String detectedContentType = detectContentType(file.getInputStream());
        if (detectedContentType == null || !ALLOWED_CONTENT_TYPES.contains(detectedContentType)) {
            throw new IOException("不支持的文件类型，仅允许 PNG、JPEG、GIF、WebP 格式");
        }

        // 确保文件名以.png结尾
        String correctedFilename = originalFilename;
        if (!correctedFilename.toLowerCase().endsWith(".png")) {
            // 如果已经有其他扩展名，去掉后再添加.png
            int lastDotIndex = correctedFilename.lastIndexOf('.');
            if (lastDotIndex > 0) {
                correctedFilename = correctedFilename.substring(0, lastDotIndex);
            }
            correctedFilename += ".png";
        }

        // 生成唯一的文件名，以时间戳开头，加上修正后的原始文件名
        String fileName = System.currentTimeMillis() + "_" + correctedFilename;
        // 构建文件的完整物理路径
        Path filePath = uploadPath.resolve(fileName);

        // 将文件输入流复制到目标路径，如果文件已存在则替换
        Files.copy(file.getInputStream(), filePath, StandardCopyOption.REPLACE_EXISTING);

        // 返回文件在 uploads 目录下的相对路径
        return Paths.get(subDirectory, fileName).toString();
//        ！！！此是备用的方案，返回文件在 uploads 目录下的相对路径，使用正斜杠（URL格式）
//        return Paths.get(subDirectory, fileName).toString().replace("\\", "/");
    }

    /**
     * 通过文件头魔数检测文件真实类型
     *
     * @param inputStream 文件输入流
     * @return 检测到的MIME类型，如果无法识别则返回null
     * @throws IOException IO异常
     */
    private String detectContentType(java.io.InputStream inputStream) throws IOException {
        // 读取文件头8个字节
        byte[] header = new byte[8];
        int bytesRead = inputStream.read(header);
        if (bytesRead < 2) {
            return null;
        }

        // 检测PNG
        if (bytesRead >= 8 && this.isMagicMatch(header, PNG_MAGIC)) {
            return "image/png";
        }

        // 检测JPEG
        if (bytesRead >= 2 && this.isMagicMatch(header, JPEG_MAGIC)) {
            return "image/jpeg";
        }

        // 尝试使用Files.probeContentType进一步检测
        try {
            // 创建一个临时文件来检测类型
            Path tempFile = Files.createTempFile("img_detect_", ".tmp");
            try {
                Files.copy(inputStream, tempFile, StandardCopyOption.REPLACE_EXISTING);
                String contentType = Files.probeContentType(tempFile);
                return contentType;
            } finally {
                Files.deleteIfExists(tempFile);
            }
        } catch (Exception e) {
            return null;
        }
    }

    /**
     * 检查字节数组是否以指定的魔数开头
     *
     * @param data 要检查的数据
     * @param magic 魔数
     * @return 是否匹配
     */
    private boolean isMagicMatch(byte[] data, byte[] magic) {
        for (int i = 0; i < magic.length; i++) {
            if (data[i] != magic[i]) {
                return false;
            }
        }
        return true;
    }

    /**
     * 读取指定相对路径文件为字节数组
     *
     * @param filePath uploads 下的相对路径
     * @return 文件字节数组
     * @throws IOException 文件读写异常
     */
    public byte[] loadFileAsBytes(String filePath) throws IOException {
        // 构建文件的完整物理路径
        Path fullPath = Paths.get(uploadDir, filePath).toAbsolutePath().normalize();
        // 读取文件所有字节
        return Files.readAllBytes(fullPath);
    }

    /**
     * 删除指定相对路径的文件
     *
     * @param filePath uploads 下的相对路径
     * @return 是否删除成功
     */
    public boolean deleteFile(String filePath) {
        try {
            // 构建文件的完整物理路径
            Path fullPath = Paths.get(uploadDir, filePath).toAbsolutePath().normalize();
            // 删除文件，如果文件不存在则返回 false
            return Files.deleteIfExists(fullPath);
        } catch (IOException e) {
            // 记录异常信息到日志，但返回 false 表示删除失败
            log.error("Failed to delete file: {}, error: {}", filePath, e.getMessage(), e);
            return false;
        }
    }
}
