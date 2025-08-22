package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"id-backend-grpc/internal/app/config"
	"id-backend-grpc/internal/models"
	"id-backend-grpc/internal/models/consts"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type VKAuthResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	IdToken      string `json:"id_token"`
	Scope        string `json:"scope"`
	State        string `json:"state"`
	UserId       int64  `json:"user_id"`
}

func (is *IdentityService) VKAuthCallback(c echo.Context) error {
	var r VKAuthResponse
	if err := c.Bind(r); err != nil || r.UserId == 0 {
		return ErrorResponse(c, "Bad request", 400)
	}
	vku, err := getUserInfo(r.AccessToken)
	if err != nil {
		return ErrorResponse(c, "Internal Server Error", 500)
	}
	// TODO я ваще не понял, какого-то хуя ид в вк это строка, пись пись как как
	id, _ := strconv.ParseUint(vku.User.ID, 10, 64)
	u, err := is.userService.GetUserByVKID(uint(id))
	if err != nil {
		return ErrorResponse(c, "Your vk account is not attached", http.StatusMethodNotAllowed)
	}
	if !is.tokenService.AccessTokenIsExists(u.ID) {
		// TODO OTPService
		ctx := c.Request().Context()
		token, _ := is.tokenService.GenerateToken(ctx, &models.CreateToken{
			Ipaddress: c.RealIP(),
			UserID:    u.ID,
			Type:      consts.TokenTemporaryAccess,
			ExpiresAt: consts.TokenTemporaryAccessLifeTime})
		return c.JSON(http.StatusCreated, models.TokenResponse(u.Email, token))
	}
	if !u.ConfirmEmail {
		// TODO OTPService confirm email
		return c.JSON(http.StatusLocked, &models.JsonResponse{Msg: "Please, confirm your email. А new letter has been sent on your email"})
	}
	t, _ := is.tokenService.GetTokenByUserIDAndType(u.ID, consts.TokenAccess)
	return c.JSON(http.StatusOK, models.TokenResponse(u.Email, t))
}

func (is *IdentityService) VKIntegration(c echo.Context) error {
	t, _ := is.tokenService.GetToken(GetToken(c))
	if t.Type != consts.TokenAccess {
		return ErrorResponse(c, "Locked", http.StatusLocked)
	}
	var r VKAuthResponse
	if err := c.Bind(r); err != nil {
		return ErrorResponse(c, "Bad request", http.StatusBadRequest)
	}
	vku, err := getUserInfo(r.AccessToken)
	if err != nil {
		return ErrorResponse(c, "Internal Server Error", 500)
	}
	u, _ := is.userService.GetUserByID(t.UserID)
	VKid, err := strconv.ParseInt(vku.User.ID, 10, 64)
	u = is.userService.SetVKID(u, VKid)
	err = is.userService.SaveUser(u)
	if err != nil {
		return ErrorResponse(c, "Internal Server Error", http.StatusInternalServerError)
	}
	return c.JSON(http.StatusCreated, models.JsonResponse{Msg: "OK. VK account is integrated"})
}

func getUserInfo(accessToken string) (*models.VKUser, error) { //TODO скопировал из старого проекта, переписать подумать как улучшить|уменьшить код!!!!
	cfg, _ := config.LoadVKConfig()
	data := url.Values{}
	data.Set("access_token", accessToken)
	data.Set("client_id", cfg.AppID)

	req, err := http.NewRequest("POST", "https://id.vk.com/oauth2/user_info", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var vkidUser models.VKUser
	if err := json.Unmarshal(body, &vkidUser); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &vkidUser, nil
}
