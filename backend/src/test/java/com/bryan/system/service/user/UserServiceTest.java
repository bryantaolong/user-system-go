package com.bryan.system.service.user;

import com.bryan.system.domain.dto.UserUpdateDTO;
import com.bryan.system.domain.entity.SysUser;
import com.bryan.system.domain.entity.UserRole;
import com.bryan.system.domain.enums.user.UserStatusEnum;
import com.bryan.system.domain.request.user.ChangeRoleRequest;
import com.bryan.system.exception.ResourceNotFoundException;
import com.bryan.system.mapper.UserMapper;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import org.springframework.security.crypto.password.PasswordEncoder;

import java.util.ArrayList;
import java.util.List;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;
import static org.mockito.Mockito.never;
import static org.mockito.Mockito.times;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

@ExtendWith(MockitoExtension.class)
class UserServiceTest {

    @Mock
    private UserMapper userMapper;

    @Mock
    private PasswordEncoder passwordEncoder;

    @Mock
    private UserRoleService userRoleService;

    @InjectMocks
    private UserService userService;

    @Test
    void shouldReturnEmptyWhenGetUsersByIdsInputEmpty() {
        List<SysUser> result = userService.getUsersByIds(List.of());
        assertTrue(result.isEmpty());
        verify(userMapper, never()).selectByIdList(List.of());
    }

    @Test
    void shouldReturnExistsByIdTrueWhenUserFound() {
        when(userMapper.selectById(1L)).thenReturn(SysUser.builder().id(1L).build());
        assertTrue(userService.existsById(1L));
    }

    @Test
    void shouldReturnExistsByIdFalseWhenUserNotFound() {
        when(userMapper.selectById(1L)).thenReturn(null);
        assertFalse(userService.existsById(1L));
    }

    @Test
    void shouldThrowWhenGetUserByIdNotFound() {
        when(userMapper.selectById(99L)).thenReturn(null);
        assertThrows(ResourceNotFoundException.class, () -> userService.getUserById(99L));
    }

    @Test
    void shouldUpdateUserPhoneAndEmail() {
        SysUser user = SysUser.builder().id(1L).phone("old").email("old@x.com").build();
        UserUpdateDTO dto = UserUpdateDTO.builder().phone("new").email("new@x.com").build();
        when(userMapper.selectById(1L)).thenReturn(user);

        SysUser actual = userService.updateUser(1L, dto);

        assertEquals("new", actual.getPhone());
        assertEquals("new@x.com", actual.getEmail());
        verify(userMapper).update(user);
    }

    @Test
    void shouldBlockAndUnblockUser() {
        SysUser user = SysUser.builder().id(1L).status(UserStatusEnum.NORMAL).build();
        when(userMapper.selectById(1L)).thenReturn(user);

        SysUser blocked = userService.blockUser(1L);
        assertEquals(UserStatusEnum.BANNED, blocked.getStatus());

        SysUser unblocked = userService.unblockUser(1L);
        assertEquals(UserStatusEnum.NORMAL, unblocked.getStatus());

        verify(userMapper, times(2)).update(user);
    }

    @Test
    void shouldChangeRoleByIds() {
        SysUser user = SysUser.builder().id(1L).roles("ROLE_OLD").build();
        ChangeRoleRequest req = new ChangeRoleRequest();
        req.setRoleIds(new ArrayList<>(List.of(1L, 2L)));
        List<UserRole> roles = List.of(
                UserRole.builder().id(1L).roleName("ROLE_USER").build(),
                UserRole.builder().id(2L).roleName("ROLE_ADMIN").build()
        );
        when(userRoleService.listByIds(req.getRoleIds())).thenReturn(roles);
        when(userMapper.selectById(1L)).thenReturn(user);

        SysUser actual = userService.changeRoleByIds(1L, req);

        assertEquals("ROLE_USER,ROLE_ADMIN", actual.getRoles());
        verify(userMapper).update(user);
    }

    @Test
    void shouldThrowWhenChangeRoleByIdsContainsMissingRole() {
        ChangeRoleRequest req = new ChangeRoleRequest();
        req.setRoleIds(new ArrayList<>(List.of(1L, 2L)));
        when(userRoleService.listByIds(req.getRoleIds()))
                .thenReturn(List.of(UserRole.builder().id(1L).roleName("ROLE_USER").build()));

        assertThrows(IllegalArgumentException.class, () -> userService.changeRoleByIds(1L, req));
    }
}
