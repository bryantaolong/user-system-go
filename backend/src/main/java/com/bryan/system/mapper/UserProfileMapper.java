package com.bryan.system.mapper;

import com.bryan.system.domain.entity.UserProfile;
import org.apache.ibatis.annotations.Mapper;

/**
 * UserProfileMapper
 *
 * @author Bryan Long
 */
@Mapper
public interface UserProfileMapper {

    int insert(UserProfile record);

    UserProfile selectByUserId(Long userId);

    UserProfile selectByRealName(String realName);

    int update(UserProfile record);
}
