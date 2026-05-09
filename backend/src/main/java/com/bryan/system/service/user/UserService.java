package com.bryan.system.service.user;

import com.bryan.system.domain.dto.UserUpdateDTO;
import com.bryan.system.domain.entity.SysUser;
import com.bryan.system.domain.entity.UserRole;
import com.bryan.system.domain.enums.user.UserStatusEnum;
import com.bryan.system.domain.request.user.ChangeRoleRequest;
import com.bryan.system.domain.request.user.UserCreateRequest;
import com.bryan.system.domain.request.user.UserSearchRequest;
import com.bryan.system.domain.response.PageResult;
import com.bryan.system.exception.BusinessException;
import com.bryan.system.exception.ResourceNotFoundException;
import com.bryan.system.mapper.UserMapper;
import com.bryan.system.util.jwt.JwtUtils;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.time.LocalDateTime;
import java.util.ArrayList;
import java.util.List;
import java.util.Optional;
import java.util.Set;
import java.util.stream.Collectors;

/**
 * 用户服务实现类，处理用户注册、登录、信息管理、导出等业务逻辑。
 *
 * @author Bryan Long
 */
@Slf4j
@Service
@RequiredArgsConstructor
public class UserService {

    private final UserMapper userMapper;
    private final PasswordEncoder passwordEncoder;
    private final UserRoleService userRoleService;

    /**
     * 管理员创建用户
     *
     * @param req 创建请求
     * @return 新增的用户实体
     */
    @Transactional
    public SysUser createUser(UserCreateRequest req) {
        if (userMapper.selectByUsername(req.getUsername()) != null) {
            throw new BusinessException("用户名已存在");
        }

        List<Long> roleIds = req.getRoleIds() == null
                ? new ArrayList<>()
                : new ArrayList<>(req.getRoleIds());
        if (roleIds.isEmpty()) {
            UserRole defaultRole = userRoleService.getDefaultRole();
            if (defaultRole == null) {
                throw new BusinessException("系统未配置默认角色");
            }
            roleIds.add(defaultRole.getId());
        }

        List<UserRole> roles = userRoleService.listByIds(roleIds);
        if (roles.size() != roleIds.size()) {
            Set<Long> exist = roles.stream()
                    .map(UserRole::getId)
                    .collect(Collectors.toSet());
            List<Long> missing = new ArrayList<>(roleIds);
            missing.removeAll(exist);
            throw new IllegalArgumentException("角色不存在：" + missing);
        }

        String roleNames = roles.stream()
                .map(UserRole::getRoleName)
                .collect(Collectors.joining(","));

        LocalDateTime now = LocalDateTime.now();

        SysUser sysUser = SysUser.builder()
                .username(req.getUsername())
                .password(passwordEncoder.encode(req.getPassword()))
                .phone(req.getPhone())
                .email(req.getEmail())
                .roles(roleNames)
                .status(UserStatusEnum.NORMAL)
                .loginFailCount(0)
                .passwordResetAt(now)
                .build();
        this.fillInsert(sysUser);

        int saved = userMapper.insert(sysUser);
        if (saved == 0) {
            throw new BusinessException("插入数据库失败");
        }

        log.info("管理员创建用户成功: id: {}, username: {} ", sysUser.getId(), sysUser.getUsername());
        return sysUser;
    }

    /**
     * 获取所有用户列表（不分页）。
     *
     * @return 包含所有用户的分页对象（Page）。
     */
    public PageResult<SysUser> getAllUsers(int pageNum,
                                           int pageSize) {
        int offset = (pageNum - 1) * pageSize;
        List<SysUser> rows = userMapper.selectPage(
                offset,
                pageSize,
                null);
        long total = userMapper.count(null);
        return PageResult.of(rows, pageNum, pageSize, total);
    }

    /**
     * 根据用户ID获取用户信息。
     *
     * @param userId 用户 ID
     * @return 用户实体对象
     * @throws ResourceNotFoundException 用户不存在时抛出
     */
    public SysUser getUserById(Long userId) {
        return Optional.ofNullable(userMapper.selectById(userId))
                .orElseThrow(() -> new ResourceNotFoundException("用户不存在"));
    }

    /**
     * 根据用户名获取用户信息。
     *
     * @param username 用户名
     * @return 用户实体对象
     */
    public SysUser getUserByUsername(String username) {
        return userMapper.selectByUsername(username);
    }

    /**
     * 通用用户搜索，支持多条件模糊查询和分页。
     *
     * @param searchRequest 搜索请求
     * @return 符合查询条件的分页对象（Page）
     */
    public PageResult<SysUser> queryUsers(UserSearchRequest searchRequest,
                                           int pageNum,
                                           int pageSize) {
        int offset = (pageNum - 1) * pageSize;
        List<SysUser> rows = userMapper.selectPage(
                offset,
                pageSize,
                searchRequest);
        long total = userMapper.count(searchRequest);
        return PageResult.of(rows, pageNum, pageSize, total);
    }

    /**
     * 根据用户 ID 列表批量查询用户信息
     *
     * @param userIds 用户 ID 列表
     * @return 用户信息列表
     */
    public List<SysUser> getUsersByIds(List<Long> userIds) {
        if (userIds == null || userIds.isEmpty()) {
            return new ArrayList<>();
        }
        return userMapper.selectByIdList(userIds);
    }

    /**
     * 检查用户是否存在
     *
     * @param userId 用户 ID
     * @return 用户是否存在
     */
    public boolean existsById(Long userId) {
        return userMapper.selectById(userId) != null;
    }

    /**
     * 更新用户基础信息（用户名和邮箱）。
     *
     * @param userId                     用户 ID
     * @param dto                        用户更新 DTO
     * @return                           更新后的用户对象
     * @throws ResourceNotFoundException 用户不存在时抛出
     * @throws BusinessException         用户名重复时抛出
     */
    public SysUser updateUser(Long userId, UserUpdateDTO dto) {
        SysUser user = this.getUserById(userId);
        if (user == null) {
            throw new BusinessException("用户不存在: " + userId);
        }

        if (dto.getPhone() != null) {
            user.setPhone(dto.getPhone());
        }
        if (dto.getEmail() != null) {
            user.setEmail(dto.getEmail());
        }

        this.fillUpdate(user);

        int rows = userMapper.update(user);
        if (rows != 1) {
            throw new BusinessException("更新失败，可能数据已变更");
        }

        log.info("用户ID: {} 的信息更新成功", userId);
        return user;
    }

    /**
     * 修改用户角色。
     *
     * @param userId 用户 ID
     * @param req  新角色字符串（多个角色用逗号分隔）
     * @return 更新后的用户对象
     * @throws ResourceNotFoundException 用户不存在时抛出
     */
    @Transactional
    public SysUser changeRoleByIds(Long userId, ChangeRoleRequest req) {
        List<Long> ids = req.getRoleIds();
        List<UserRole> roles = userRoleService.listByIds(ids);

        if (roles.size() != ids.size()) {
            Set<Long> exist = roles.stream()
                    .map(UserRole::getId)
                    .collect(Collectors.toSet());
            ids.removeAll(exist);
            throw new IllegalArgumentException("角色不存在：" + ids);
        }
        String roleNames = roles.stream()
                .map(UserRole::getRoleName)
                .collect(Collectors.joining(","));
        SysUser user = this.getUserById(userId);
        user.setRoles(roleNames);

        this.fillUpdate(user);

        userMapper.update(user);
        return user;
    }

    /**
     * 重置用户密码（管理员）。
     *
     * @param userId      用户 ID
     * @param newPassword 新密码（明文）
     * @return 更新后的用户对象
     * @throws ResourceNotFoundException 用户不存在时抛出
     * @throws BusinessException         旧密码验证失败时抛出
     */
    public SysUser resetPassword(Long userId, String newPassword) {
        SysUser user = this.getUserById(userId);
        user.setPassword(passwordEncoder.encode(newPassword));
        user.setPasswordResetAt(LocalDateTime.now());

        this.fillUpdate(user);

        userMapper.update(user);
        log.info("用户ID: {} 的密码强制修改成功", userId);
        return user;
    }

    /**
     * 封禁指定用户。
     *
     * @param userId 用户 ID
     * @return 更新后的用户对象
     * @throws ResourceNotFoundException 用户不存在时抛出
     */
    public SysUser blockUser(Long userId) {
        SysUser user = this.getUserById(userId);
        user.setStatus(UserStatusEnum.BANNED);

        this.fillUpdate(user);

        userMapper.update(user);
        log.info("用户ID: {} 封禁成功", userId);
        return user;
    }

    /**
     * 解封指定用户。
     *
     * @param userId 用户 ID
     * @return 更新后的用户对象
     * @throws ResourceNotFoundException 用户不存在时抛出
     */
    public SysUser unblockUser(Long userId) {
        SysUser user = this.getUserById(userId);
        user.setStatus(UserStatusEnum.NORMAL);

        this.fillUpdate(user);

        userMapper.update(user);
        log.info("用户ID: {} 解封成功", userId);
        return user;
    }

    /**
     * 删除用户（逻辑删除）。
     *
     * @param userId 用户 ID
     * @return 被删除的用户对象
     * @throws ResourceNotFoundException 用户不存在时抛出
     */
    public Long deleteUser(Long userId) {
        int rows = userMapper.deleteById(userId, LocalDateTime.now(), JwtUtils.getCurrentUsername());
        if (rows == 0) {
            throw new ResourceNotFoundException("用户不存在或已被删除");
        }
        log.info("用户ID: {} 删除成功 (逻辑删除)", userId);
        return userId;
    }

    private void fillInsert(SysUser user) {
        LocalDateTime now = LocalDateTime.now();
        Long operator = JwtUtils.getCurrentUserId();

        user.setDeleted(0);
        user.setVersion(0);
        user.setCreatedAt(now);
        user.setUpdatedAt(now);
        user.setUpdatedBy(operator.toString());
        user.setCreatedBy(operator.toString());

    }

    private  void fillUpdate(SysUser user) {
        LocalDateTime now = LocalDateTime.now();
        Long operator = JwtUtils.getCurrentUserId();

        user.setVersion(user.getVersion() + 1);
        user.setUpdatedAt(now);
        user.setUpdatedBy(operator.toString());
    }
}
