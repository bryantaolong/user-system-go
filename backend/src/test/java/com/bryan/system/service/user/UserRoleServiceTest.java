package com.bryan.system.service.user;

import com.bryan.system.domain.entity.UserRole;
import com.bryan.system.mapper.UserRoleMapper;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import java.util.List;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertSame;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

@ExtendWith(MockitoExtension.class)
class UserRoleServiceTest {

    @Mock
    private UserRoleMapper userRoleMapper;

    @InjectMocks
    private UserRoleService userRoleService;

    @Test
    void shouldListAllRoles() {
        List<UserRole> expected = List.of(UserRole.builder().id(1L).roleName("ROLE_USER").build());
        when(userRoleMapper.selectAll()).thenReturn(expected);

        List<UserRole> actual = userRoleService.listAll();

        assertEquals(expected, actual);
        verify(userRoleMapper).selectAll();
    }

    @Test
    void shouldReturnDefaultRole() {
        UserRole expected = UserRole.builder().id(2L).roleName("ROLE_ADMIN").isDefault(true).build();
        when(userRoleMapper.selectOneByIsDefaultTrue()).thenReturn(expected);

        UserRole actual = userRoleService.getDefaultRole();

        assertSame(expected, actual);
        verify(userRoleMapper).selectOneByIsDefaultTrue();
    }

    @Test
    void shouldListRolesByIds() {
        List<Long> ids = List.of(1L, 2L);
        List<UserRole> expected = List.of(
                UserRole.builder().id(1L).roleName("ROLE_USER").build(),
                UserRole.builder().id(2L).roleName("ROLE_ADMIN").build()
        );
        when(userRoleMapper.selectByIdList(ids)).thenReturn(expected);

        List<UserRole> actual = userRoleService.listByIds(ids);

        assertEquals(expected, actual);
        verify(userRoleMapper).selectByIdList(ids);
    }
}
