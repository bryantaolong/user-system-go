package com.bryan.system.domain.dto;

import com.bryan.system.domain.enums.user.GenderEnum;
import lombok.Builder;
import lombok.Data;

import java.time.LocalDateTime;

/**
 * UserProfileUpdateDTO
 *
 * @author Bryan Long
 */
@Data
@Builder
public class UserProfileUpdateDTO {

    private Long userId;

    private String realName;

    private GenderEnum gender;

    private LocalDateTime birthday;

    private String avatar;
}
