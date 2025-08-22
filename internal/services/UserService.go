package services

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"id-backend-grpc/internal/models"
	"id-backend-grpc/internal/repositories"
	"io"
	"mime/multipart"
	"os"
	"sort"
	"strconv"
	"strings"
)

type UserService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	u, err := s.repo.GetUserByEmail(email)
	return u, err
}

func (s *UserService) CreateUser(Ctx context.Context, cu *models.CreateUser) (*models.User, error) {
	salt, hashedPassword, err := generateSaltForPassword(cu.Password)
	if err != nil {
		return nil, fmt.Errorf("error while hash password: %s", err.Error())
	}
	u := models.NewUser(cu.Email, hashedPassword, salt, cu.Name)
	if cu.TelegramID != 0 {
		u.TelegramId = cu.TelegramID
	}
	if cu.VKid != 0 {
		u.VkId = cu.VKid
	}
	_ = s.repo.Create(Ctx, u)
	return u, nil
}

func (s *UserService) SaveUser(u *models.User) error {
	return s.repo.Save(u)
}

func (s *UserService) GetUserByID(id uint) (*models.User, error) { // TODO не нравится название метода, переделать
	return s.repo.GetUserByID(id)
}

func (s *UserService) GetUserByTelegramID(id uint) (*models.User, error) {
	return s.repo.GetUserByTelegramID(id)
}

func (s *UserService) GetUserByVKID(id uint) (*models.User, error) {
	return s.repo.GetUserByVKID(id)
}

func (s *UserService) EditPassword(u *models.User, newPassword string) (*models.User, error) {
	// TODO дописать валидатор (проверка что пароль не одинаков со старым и другие приколямбы)
	// TODO так же раскидать модель USER на слой DTO(возможно)
	salt, hashedPassword, err := generateSaltForPassword(newPassword)
	if err != nil {
		return nil, fmt.Errorf("error while hash password: %s", err.Error())
	}
	u.Password = hashedPassword
	u.Salt = salt
	return u, nil
}

func (s *UserService) VerifyPasswordWithSalt(hexedPassword, salt, inputPassword string) bool {
	saltedPassword := inputPassword + salt
	hashedPassword, err := hex.DecodeString(hexedPassword)
	if err != nil {
		return false
	}
	fmt.Println(bcrypt.CompareHashAndPassword(hashedPassword, []byte(saltedPassword)))
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(saltedPassword))

	return err == nil
}

func (s *UserService) ConfirmEmailUserByID(id uint) error {
	u, err := s.repo.GetUserByID(id)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if u == nil {
		return fmt.Errorf("user with ID %d not found", id)
	}
	u.ConfirmEmail = true
	if err := s.repo.Save(u); err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}
	return nil
}

func (s *UserService) ValidateTelegramData(tu *models.TelegramUser, botToken string) bool { // TODO возможно переписать, пока что оставил так
	data := map[string]string{
		"id":         strconv.FormatInt(tu.ID, 10),
		"first_name": tu.FirstName,
		"auth_date":  strconv.FormatInt(tu.AuthDate, 10),
	}
	if tu.LastName != "" {
		data["last_name"] = tu.LastName
	}
	if tu.Username != "" {
		data["username"] = tu.Username
	}
	if tu.PhotoURL != "" {
		data["photo_url"] = tu.PhotoURL
	}

	var dataCheckArr []string
	for key, value := range data {
		dataCheckArr = append(dataCheckArr, fmt.Sprintf("%s=%s", key, value))
	}

	sort.Strings(dataCheckArr)
	dataCheckString := strings.Join(dataCheckArr, "\n")

	secretKey := sha256.Sum256([]byte(botToken))
	h := hmac.New(sha256.New, secretKey[:])
	h.Write([]byte(dataCheckString))
	computedHash := hex.EncodeToString(h.Sum(nil))

	return computedHash == tu.Hash
}

func (s *UserService) UploadAvatar(user *models.User, form *multipart.Form) bool {
	files := form.File["files"]
	if len(files) == 0 {
		return false
	}
	fileUpload := files[0]
	if !strings.HasPrefix(fileUpload.Header.Get("Content-Type"), "image/") {
		return false
	}
	src, err := fileUpload.Open()
	if err != nil {
		return false
	}
	defer src.Close()

	fileBytes, err := io.ReadAll(src)
	if err != nil {
		return false
	}

	set := setAvatar(fileBytes, user)
	if !set {
		return false
	}

	return true
}

func (s *UserService) GetAvatar(user *models.User) *os.File {
	//TODO дописать, пакет поставил
	return nil
}

func (s *UserService) SetTelegramID(user *models.User, id int64) *models.User {
	user.TelegramId = id
	return user
}

func (s *UserService) SetVKID(user *models.User, id int64) *models.User {
	user.VkId = id
	return user
}

func setAvatar(fileData []byte, user *models.User) bool {
	//TODO дописать, пакет поставил
	//if len(fileData) == 0 {
	//	return false
	//}
	//userKeyAvatar := sha256.Sum256([]byte(user.Email))
	//err := connections.Minio.Set(hex.EncodeToString(userKeyAvatar[:]), fileData, 0)
	//if err != nil {
	//	return false
	//}

	return false
}

// возвращает соль и хэшированный пароль
func generateSaltForPassword(password string) (string, string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", "", err
	}
	saltedPassword := password + hex.EncodeToString(salt)
	hashed, _ := bcrypt.GenerateFromPassword([]byte(saltedPassword), bcrypt.DefaultCost)
	return hex.EncodeToString(salt), hex.EncodeToString(hashed[:]), nil
}
