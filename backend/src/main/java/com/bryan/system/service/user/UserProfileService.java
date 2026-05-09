package com.bryan.system.service.user;

import com.bryan.system.domain.dto.UserProfileUpdateDTO;
import com.bryan.system.domain.entity.UserProfile;
import com.bryan.system.exception.BusinessException;
import com.bryan.system.exception.ResourceNotFoundException;
import com.bryan.system.mapper.UserProfileMapper;
import com.bryan.system.service.file.LocalFileService;
import com.bryan.system.util.jwt.JwtUtils;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import org.springframework.web.multipart.MultipartFile;

import java.io.IOException;
import java.time.LocalDateTime;

/**
 * 用户资料业务服务
 * 提供用户资料的创建、查询、更新能力。
 *
 * @author Bryan Long
 */
@Slf4j
@Service
@RequiredArgsConstructor
public class UserProfileService {

    private final UserProfileMapper userProfileMapper;
    private final LocalFileService localFileService;

    /**
     * 创建用户资料
     *
     * @param record 用户资料实体
     * @return 已持久化的实体
     * @throws BusinessException 数据库插入失败
     */
    public UserProfile createUserProfile(UserProfile record) {
        this.fillInsert(record);

        int inserted = userProfileMapper.insert(record);
        if (inserted <= 0) {
            throw new BusinessException("创建用户信息失败");
        }
        log.info("用户信息创建成功，用户ID: {}", record.getUserId());
        return record;
    }

    /**
     * 根据用户主键查询用户资料
     *
     * @param userId 用户主键
     * @return 用户资料实体
     * @throws ResourceNotFoundException 用户资料不存在
     */
    public UserProfile getUserProfileByUserId(Long userId) {
        UserProfile profile = userProfileMapper.selectByUserId(userId);
        if (profile == null) {
            throw new ResourceNotFoundException("用户信息不存在");
        }
        return profile;
    }

    /**
     * 根据用户主键查询用户资料，如果不存在则返回空实体
     *
     * @param userId 用户主键
     * @return 用户资料实体，如果不存在则返回空的实体
     */
    public UserProfile getUserProfileByUserIdOrEmpty(Long userId) {
        UserProfile profile = userProfileMapper.selectByUserId(userId);
        if (profile == null) {
            log.warn("用户资料不存在，用户ID: {}, 返回空实体", userId);
            return UserProfile.builder().userId(userId).build();
        }
        return profile;
    }

    /**
     * 根据真实姓名查询用户资料
     *
     * @param realName 真实姓名
     * @return 用户资料实体
     * @throws ResourceNotFoundException 用户资料不存在
     */
    public UserProfile getUserProfileByRealName(String realName) {
        UserProfile profile = userProfileMapper.selectByRealName(realName);
        if (profile == null) {
            throw new ResourceNotFoundException("用户信息不存在");
        }
        return profile;
    }

    /**
     * 更新用户资料（字段可选）
     *
     * @param userId 用户主键
     * @param dto    待更新字段封装
     * @return 更新后的用户资料实体
     * @throws BusinessException 更新失败
     */
    public UserProfile updateUserProfile(Long userId, UserProfileUpdateDTO dto) {
        UserProfile profile = this.getUserProfileByUserId(userId);

        if (dto.getRealName() != null) {
            profile.setRealName(dto.getRealName());
        }
        if (dto.getGender() != null) {
            profile.setGender(dto.getGender());
        }
        if (dto.getBirthday() != null) {
            profile.setBirthday(dto.getBirthday());
        }
        if (dto.getAvatar() != null) {
            profile.setAvatar(dto.getAvatar());
        }

        this.fillUpdate(profile);

        int updated = userProfileMapper.update(profile);
        if (updated == 0) {
            throw new BusinessException("用户信息更新失败");
        }
        log.info("用户信息更新成功，用户ID: {}", userId);
        return profile;
    }

    /**
     * 上传并更新用户头像
     *
     * @param userId 用户主键
     * @param file   头像文件
     * @return 更新后的头像相对路径
     * @throws BusinessException 文件存储或数据库更新失败
     */
    public String updateAvatar(Long userId, MultipartFile file) {
        UserProfile profile = this.getUserProfileByUserId(userId);

        try {
            // 1. 存储新头像文件
            String avatarPath = localFileService.storeFile(file, "avatars");

            // 2. 如果原有头像存在，尝试删除旧文件（可选，根据业务需求）
            if (profile.getAvatar() != null && !profile.getAvatar().isEmpty()) {
                localFileService.deleteFile(profile.getAvatar());
            }

            // 3. 更新数据库
            profile.setAvatar(avatarPath);

            this.fillUpdate(profile);

            int updated = userProfileMapper.update(profile);
            if (updated == 0) {
                throw new BusinessException("头像更新失败");
            }

            log.info("用户头像更新成功，用户ID: {}, 路径: {}", userId, avatarPath);
            return avatarPath;
        } catch (IOException e) {
            log.error("用户头像上传失败，用户ID: {}", userId, e);
            throw new BusinessException("头像上传失败: " + e.getMessage());
        }
    }

    private void fillInsert(UserProfile record) {
        LocalDateTime now = LocalDateTime.now();
        Long operator = JwtUtils.getCurrentUserId();

        record.setDeleted(0);
        record.setVersion(0);
        record.setCreatedAt(now);
        record.setUpdatedAt(now);
        record.setUpdatedBy(operator.toString());
        record.setCreatedBy(operator.toString());
    }

    private void fillUpdate(UserProfile profile) {
        LocalDateTime now = LocalDateTime.now();
        Long operator = JwtUtils.getCurrentUserId();

        profile.setVersion(profile.getVersion() + 1);
        profile.setUpdatedAt(now);
        profile.setUpdatedBy(operator.toString());
    }
}
