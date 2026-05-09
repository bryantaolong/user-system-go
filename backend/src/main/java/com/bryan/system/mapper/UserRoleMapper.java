package com.bryan.system.mapper;

import com.bryan.system.domain.entity.UserRole;
import org.apache.ibatis.annotations.Mapper;
import org.apache.ibatis.annotations.Param;

import java.util.Collection;
import java.util.List;

/**
 * UserRoleMapper
 *
 * @author Bryan Long
 */
@Mapper
public interface UserRoleMapper {

    UserRole selectOneByIsDefaultTrue();

    List<UserRole> selectAll();

    List<UserRole> selectByIdList(@Param("ids") Collection<Long> ids);
}
