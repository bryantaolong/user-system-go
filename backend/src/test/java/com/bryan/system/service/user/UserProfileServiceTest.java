package com.bryan.system.service.user;

import com.bryan.system.domain.dto.UserProfileUpdateDTO;
import com.bryan.system.domain.entity.UserProfile;
import com.bryan.system.exception.BusinessException;
import com.bryan.system.exception.ResourceNotFoundException;
import com.bryan.system.mapper.UserProfileMapper;
import com.bryan.system.service.file.LocalFileService;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import org.springframework.mock.web.MockMultipartFile;

import java.io.IOException;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.never;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

@ExtendWith(MockitoExtension.class)
class UserProfileServiceTest {

    @Mock
    private UserProfileMapper userProfileMapper;

    @Mock
    private LocalFileService localFileService;

    @InjectMocks
    private UserProfileService userProfileService;

    @Test
    void shouldCreateUserProfileWhenInsertSuccess() {
        UserProfile profile = UserProfile.builder().userId(10L).build();
        when(userProfileMapper.insert(profile)).thenReturn(1);

        UserProfile actual = userProfileService.createUserProfile(profile);

        assertEquals(profile, actual);
        verify(userProfileMapper).insert(profile);
    }

    @Test
    void shouldThrowWhenCreateUserProfileInsertFailed() {
        UserProfile profile = UserProfile.builder().userId(10L).build();
        when(userProfileMapper.insert(profile)).thenReturn(0);

        assertThrows(BusinessException.class, () -> userProfileService.createUserProfile(profile));
    }

    @Test
    void shouldGetUserProfileByUserId() {
        UserProfile profile = UserProfile.builder().userId(10L).realName("Tom").build();
        when(userProfileMapper.selectByUserId(10L)).thenReturn(profile);

        UserProfile actual = userProfileService.getUserProfileByUserId(10L);

        assertEquals(profile, actual);
    }

    @Test
    void shouldThrowWhenGetUserProfileByUserIdNotFound() {
        when(userProfileMapper.selectByUserId(10L)).thenReturn(null);

        assertThrows(ResourceNotFoundException.class, () -> userProfileService.getUserProfileByUserId(10L));
    }

    @Test
    void shouldReturnEmptyProfileWhenUserProfileMissing() {
        when(userProfileMapper.selectByUserId(10L)).thenReturn(null);

        UserProfile actual = userProfileService.getUserProfileByUserIdOrEmpty(10L);

        assertNotNull(actual);
        assertEquals(10L, actual.getUserId());
    }

    @Test
    void shouldUpdateUserProfileFields() {
        UserProfile profile = UserProfile.builder().userId(10L).realName("Old").avatar("old.png").build();
        UserProfileUpdateDTO dto = UserProfileUpdateDTO.builder()
                .realName("New")
                .avatar("new.png")
                .build();
        when(userProfileMapper.selectByUserId(10L)).thenReturn(profile);
        when(userProfileMapper.update(profile)).thenReturn(1);

        UserProfile actual = userProfileService.updateUserProfile(10L, dto);

        assertEquals("New", actual.getRealName());
        assertEquals("new.png", actual.getAvatar());
        verify(userProfileMapper).update(profile);
    }

    @Test
    void shouldThrowWhenUpdateUserProfileAffectsNoRows() {
        UserProfile profile = UserProfile.builder().userId(10L).build();
        UserProfileUpdateDTO dto = UserProfileUpdateDTO.builder().realName("New").build();
        when(userProfileMapper.selectByUserId(10L)).thenReturn(profile);
        when(userProfileMapper.update(profile)).thenReturn(0);

        assertThrows(BusinessException.class, () -> userProfileService.updateUserProfile(10L, dto));
    }

    @Test
    void shouldUpdateAvatarAndDeleteOldAvatar() throws IOException {
        UserProfile profile = UserProfile.builder().userId(10L).avatar("avatars/old.png").build();
        MockMultipartFile file = new MockMultipartFile("file", "a.png", "image/png", new byte[]{1, 2});
        when(userProfileMapper.selectByUserId(10L)).thenReturn(profile);
        when(localFileService.storeFile(file, "avatars")).thenReturn("avatars/new.png");
        when(userProfileMapper.update(profile)).thenReturn(1);

        String path = userProfileService.updateAvatar(10L, file);

        assertEquals("avatars/new.png", path);
        assertEquals("avatars/new.png", profile.getAvatar());
        verify(localFileService).deleteFile("avatars/old.png");
        verify(userProfileMapper).update(profile);
    }

    @Test
    void shouldThrowBusinessExceptionWhenAvatarStoreFails() throws IOException {
        UserProfile profile = UserProfile.builder().userId(10L).avatar("avatars/old.png").build();
        MockMultipartFile file = new MockMultipartFile("file", "a.png", "image/png", new byte[]{1, 2});
        when(userProfileMapper.selectByUserId(10L)).thenReturn(profile);
        when(localFileService.storeFile(any(), any())).thenThrow(new IOException("io"));

        assertThrows(BusinessException.class, () -> userProfileService.updateAvatar(10L, file));
        verify(userProfileMapper, never()).update(any());
    }
}
