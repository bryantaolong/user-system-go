package com.bryan.system.service.user;

import com.bryan.system.domain.entity.UserRole;
import com.bryan.system.mapper.UserRoleMapper;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;

import java.util.Collection;
import java.util.List;

/**
 * 用户角色业务服务
 * 提供角色列表查询及按 ID 批量查询能力。
 *
 * @author Bryan Long
 */
@Service
@RequiredArgsConstructor
public class UserRoleService {

    private final UserRoleMapper userRoleMapper;

    /**
     * 查询全部角色
     *
     * @return 角色实体列表
     */
    public List<UserRole> listAll() {
        return userRoleMapper.selectAll();
    }

    /**
     * 获取默认角色
     *
     * @return 默认角色；若未配置则返回 null
     */
    public UserRole getDefaultRole() {
        return userRoleMapper.selectOneByIsDefaultTrue();
    }

    /**
     * 根据角色主键列表批量查询
     *
     * @param ids 角色主键集合
     * @return 角色实体列表
     */
    public List<UserRole> listByIds(Collection<Long> ids) {
        return userRoleMapper.selectByIdList(ids);
    }
}
