package com.bryan.system.controller.auth

import com.bryan.system.domain.converter.UserConverter
import com.bryan.system.domain.entity.UserProfile
import com.bryan.system.domain.request.auth.LoginRequest
import com.bryan.system.domain.request.auth.RegisterRequest
import com.bryan.system.domain.request.user.ChangePasswordRequest
import com.bryan.system.domain.response.Result
import com.bryan.system.domain.vo.UserVO
import com.bryan.system.service.auth.AuthService
import com.bryan.system.service.user.UserProfileService
import jakarta.validation.Valid
import org.springframework.security.access.prepost.PreAuthorize
import org.springframework.security.core.annotation.AuthenticationPrincipal
import org.springframework.security.core.userdetails.UserDetails
import org.springframework.validation.annotation.Validated
import org.springframework.web.bind.annotation.DeleteMapping
import org.springframework.web.bind.annotation.GetMapping
import org.springframework.web.bind.annotation.PostMapping
import org.springframework.web.bind.annotation.PutMapping
import org.springframework.web.bind.annotation.RequestBody
import org.springframework.web.bind.annotation.RequestMapping
import org.springframework.web.bind.annotation.RequestParam
import org.springframework.web.bind.annotation.RestController

@Validated
@RestController
@RequestMapping("/api/auth")
class AuthController(
    private val authService: AuthService,
    private val userProfileService: UserProfileService
) {
    @PostMapping("/register")
    @PreAuthorize("permitAll()")
    fun register(@RequestBody @Valid req: RegisterRequest): Result<UserVO?> {
        val registered = authService.register(req)
        userProfileService.createUserProfile(UserProfile(userId = registered.id))
        return Result.success(UserConverter.toUserVO(registered))
    }

    @PostMapping("/login")
    @PreAuthorize("permitAll()")
    fun login(@RequestBody @Valid req: LoginRequest): Result<String> = Result.success(authService.login(req))

    @GetMapping("/me")
    @PreAuthorize("isAuthenticated()")
    fun getCurrentUser(): Result<UserVO?> = Result.success(UserConverter.toUserVO(authService.getCurrentUser()))

    @GetMapping("/validate")
    @PreAuthorize("permitAll()")
    fun validate(@RequestParam token: String, @AuthenticationPrincipal userDetails: UserDetails?): Result<String> {
        if (!authService.validateToken(token)) return Result.success("Invalid token")
        if (userDetails != null) {
            if (!userDetails.isAccountNonLocked) return Result.success("Account locked")
            if (!userDetails.isAccountNonExpired) return Result.success("Account expired")
            if (!userDetails.isEnabled) return Result.success("Account disabled")
        }
        return Result.success("Validation passed")
    }

    @PutMapping("/password")
    @PreAuthorize("isAuthenticated()")
    fun changePassword(@RequestBody @Valid req: ChangePasswordRequest): Result<UserVO?> =
        Result.success(UserConverter.toUserVO(authService.changePassword(req.oldPassword, req.newPassword)))

    @DeleteMapping
    @PreAuthorize("isAuthenticated()")
    fun deleteAccount(): Result<UserVO?> = Result.success(UserConverter.toUserVO(authService.deleteAccount()))

    @GetMapping("/logout")
    @PreAuthorize("isAuthenticated()")
    fun logout(): Boolean = authService.logout()
}
