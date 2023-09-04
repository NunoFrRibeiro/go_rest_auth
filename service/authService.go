package Service

import (
	"fmt"

	"github.com/NunoFrRibeiro/go_rest_auth/domain"
	"github.com/NunoFrRibeiro/go_rest_auth/dto"
	errs "github.com/NunoFrRibeiro/go_rest_auth/err"
	"github.com/NunoFrRibeiro/go_rest_auth/logger"
	"github.com/golang-jwt/jwt"
)

type AuthService interface {
	Login(dto.LoginRequest) (*dto.LoginResponse, *errs.AppError)
	Verify(urlParams map[string]string) *errs.AppError
	Refresh(request dto.RefreshTokenRequest) (*dto.LoginResponse, *errs.AppError)
}

type DefaultAuthService struct {
	repo            domain.AuthRepository
	rolePermissions domain.RolePermisisons
}

func (s DefaultAuthService) Login(req dto.LoginRequest) (*dto.LoginResponse, *errs.AppError) {
	var (
		appError                  *errs.AppError
		login                     *domain.Login
		accessToken, refreshToken string
	)

	if login, appError = s.repo.FindBy(req.Username, req.Password); appError != nil {
		return nil, appError
	}

	claims := login.ClaimsForAccessToken()
	authToken := domain.NewAuthToken(claims)

	if accessToken, appError = s.repo.GenerateAndStoreRefreshTokenToStore(authToken); appError != nil {
		return nil, appError
	}

	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s DefaultAuthService) Verify(urlParams map[string]string) *errs.AppError {
	if jwtToken, err := jwtTokenFromString(urlParams["token"]); err != nil {
		return errs.AuthenticationError(err.Error())
	} else {
		if jwtToken.Valid {
			claims := jwtToken.Claims.(*domain.AccessTokenClaims)
			if claims.IsUserRole() {
				if !claims.IsRequestVerifiedWithTokenClaims(urlParams) {
					return errs.AuthenticationError("request not verified with the token claims")
				}
			}
			isAuthorized := s.rolePermissions.IsAuthorizedFor(claims.Role, urlParams["routeName"])
			if !isAuthorized {
				return errs.AuthenticationError(fmt.Sprintf("%s role is not authorized", claims.Role))
			}
			return nil
		} else {
			return errs.AuthenticationError("Invalid token")
		}
	}
}

func (s DefaultAuthService) Refresh(request dto.RefreshTokenRequest) (*dto.LoginResponse, *errs.AppError) {
	if vErr := request.IsAccessTokenValid(); vErr != nil {
		if vErr.Errors == jwt.ValidationErrorExpired {
			var (
				appError    *errs.AppError
				accessToken string
			)

			if accessToken, appError = domain.NewAccesTokenFromRefresh(request.RefreshToken); appError != nil {
				return nil, appError
			}

			return &dto.LoginResponse{
				AccessToken: accessToken,
			}, nil
		}
		return nil, errs.AuthenticationError("invalid token")
	}
	return nil, errs.AuthenticationError("cannot generate access token from refresh token")
}

func jwtTokenFromString(jwtToken string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(jwtToken, &domain.AccessTokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(domain.HMAC_SAMPLE_SECRET), nil
	})
	if err != nil {
		logger.Error("error while parsing token: %s", jwtToken)
		return nil, err
	}
	return token, nil
}

func NewLoginService(repo domain.AuthRepository, permissions domain.RolePermisisons) DefaultAuthService {
	return DefaultAuthService{
		repo,
		permissions,
	}
}
