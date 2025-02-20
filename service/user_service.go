package service

import (
	"context"

	"github.com/PhantomX7/dhamma/dto"
	"github.com/PhantomX7/dhamma/dto/request"
	"github.com/PhantomX7/dhamma/dto/response"
	"github.com/PhantomX7/dhamma/entity"
	"github.com/PhantomX7/dhamma/helpers"
	"github.com/PhantomX7/dhamma/repository"
	"github.com/jinzhu/copier"
	"golang.org/x/crypto/bcrypt"
)

type (
	UserService interface {
		Register(ctx context.Context, req request.UserCreateRequest) (response.UserResponse, error)
		GetAllUserWithPagination(ctx context.Context, req request.PaginationRequest) (response.UserPaginationResponse, error)
		GetUserById(ctx context.Context, userId string) (response.UserResponse, error)
		GetUserByEmail(ctx context.Context, email string) (response.UserResponse, error)
		Update(ctx context.Context, req request.UserUpdateRequest, userId string) (response.UserUpdateResponse, error)
		Delete(ctx context.Context, userId string) error
		Verify(ctx context.Context, req request.UserLoginRequest) (response.UserLoginResponse, error)
		// SendVerificationEmail(ctx context.Context, req dto.SendVerificationEmailRequest) error
		// VerifyEmail(ctx context.Context, req dto.VerifyEmailRequest) (dto.VerifyEmailResponse, error)
	}

	userService struct {
		userRepo   repository.UserRepository
		jwtService JWTService
	}
)

func NewUserService(userRepo repository.UserRepository, jwtService JWTService) UserService {
	return &userService{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

const (
	LOCAL_URL          = "http://localhost:3000"
	VERIFY_EMAIL_ROUTE = "register/verify_email"
)

func (s *userService) Register(ctx context.Context, req request.UserCreateRequest) (res response.UserResponse, err error) {

	// if req.Image != nil {
	// 	imageId := uuid.New()
	// 	ext := utils.GetExtensions(req.Image.Filename)

	// 	filename = fmt.Sprintf("profile/%s.%s", imageId, ext)
	// 	if err := utils.UploadFile(req.Image, filename); err != nil {
	// 		return response.UserResponse{}, err
	// 	}
	// }

	user := entity.User{}

	_ = copier.Copy(&user, &req)

	password, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		err = dto.ErrCreateUser
		return
	}
	user.Password = string(password)

	userReg, err := s.userRepo.RegisterUser(ctx, nil, user)
	if err != nil {
		return response.UserResponse{}, dto.ErrCreateUser
	}

	return response.UserResponse{
		ID:       userReg.ID.String(),
		Username: userReg.Username,
		IsAdmin:  userReg.IsAdmin,
	}, nil
}

func (s *userService) GetAllUserWithPagination(ctx context.Context, req request.PaginationRequest) (response.UserPaginationResponse, error) {
	dataWithPaginate, err := s.userRepo.GetAllUserWithPagination(ctx, nil, req)
	if err != nil {
		return response.UserPaginationResponse{}, err
	}

	var datas []response.UserResponse
	for _, user := range dataWithPaginate.Users {
		data := response.UserResponse{
			ID: user.ID.String(),
		}

		datas = append(datas, data)
	}

	return response.UserPaginationResponse{
		Data: datas,
		PaginationResponse: response.PaginationResponse{
			Page:    dataWithPaginate.Page,
			PerPage: dataWithPaginate.PerPage,
			MaxPage: dataWithPaginate.MaxPage,
			Count:   dataWithPaginate.Count,
		},
	}, nil
}

func (s *userService) GetUserById(ctx context.Context, userId string) (response.UserResponse, error) {
	user, err := s.userRepo.GetUserById(ctx, nil, userId)
	if err != nil {
		return response.UserResponse{}, dto.ErrGetUserById
	}

	return response.UserResponse{
		ID:       user.ID.String(),
		Username: user.Username,
		IsAdmin:  user.IsAdmin,
	}, nil
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (response.UserResponse, error) {
	emails, err := s.userRepo.GetUserByEmail(ctx, nil, email)
	if err != nil {
		return response.UserResponse{}, dto.ErrGetUserByEmail
	}

	return response.UserResponse{
		ID:       emails.ID.String(),
		Username: emails.Username,
		IsAdmin:  emails.IsAdmin,
	}, nil
}

func (s *userService) Update(ctx context.Context, req request.UserUpdateRequest, userId string) (response.UserUpdateResponse, error) {
	user, err := s.userRepo.GetUserById(ctx, nil, userId)
	if err != nil {
		return response.UserUpdateResponse{}, dto.ErrUserNotFound
	}

	_ = copier.Copy(&user, &req)

	userUpdate, err := s.userRepo.UpdateUser(ctx, nil, user)
	if err != nil {
		return response.UserUpdateResponse{}, dto.ErrUpdateUser
	}

	return response.UserUpdateResponse{
		ID:       userUpdate.ID.String(),
		Username: userUpdate.Username,
		IsAdmin:  userUpdate.IsAdmin,
	}, nil
}

func (s *userService) Delete(ctx context.Context, userId string) error {
	user, err := s.userRepo.GetUserById(ctx, nil, userId)
	if err != nil {
		return dto.ErrUserNotFound
	}

	err = s.userRepo.DeleteUser(ctx, nil, user.ID.String())
	if err != nil {
		return dto.ErrDeleteUser
	}

	return nil
}

func (s *userService) Verify(ctx context.Context, req request.UserLoginRequest) (response.UserLoginResponse, error) {
	check, flag, err := s.userRepo.CheckEmail(ctx, nil, req.Email)
	if err != nil || !flag {
		return response.UserLoginResponse{}, dto.ErrEmailNotFound
	}

	checkPassword, err := helpers.CheckPassword(check.Password, []byte(req.Password))
	if err != nil || !checkPassword {
		return response.UserLoginResponse{}, dto.ErrPasswordNotMatch
	}

	token := s.jwtService.GenerateToken(check.ID.String(), "admin")

	return response.UserLoginResponse{
		Token: token,
		Role:  "admin",
	}, nil
}

// func makeVerificationEmail(receiverEmail string) (map[string]string, error) {
// 	expired := time.Now().Add(time.Hour * 24).Format("2006-01-02 15:04:05")
// 	plainText := receiverEmail + "_" + expired
// 	token, err := utils.AESEncrypt(plainText)
// 	if err != nil {
// 		return nil, err
// 	}

// 	verifyLink := LOCAL_URL + "/" + VERIFY_EMAIL_ROUTE + "?token=" + token

// 	readHtml, err := os.ReadFile("utils/email-template/base_mail.html")
// 	if err != nil {
// 		return nil, err
// 	}

// 	data := struct {
// 		Email  string
// 		Verify string
// 	}{
// 		Email:  receiverEmail,
// 		Verify: verifyLink,
// 	}

// 	tmpl, err := template.New("custom").Parse(string(readHtml))
// 	if err != nil {
// 		return nil, err
// 	}

// 	var strMail bytes.Buffer
// 	if err := tmpl.Execute(&strMail, data); err != nil {
// 		return nil, err
// 	}

// 	draftEmail := map[string]string{
// 		"subject": "Cakno - Go Gin Template",
// 		"body":    strMail.String(),
// 	}

// 	return draftEmail, nil
// }

// func (s *userService) SendVerificationEmail(ctx context.Context, req dto.SendVerificationEmailRequest) error {
// 	user, err := s.userRepo.GetUserByEmail(ctx, nil, req.Email)
// 	if err != nil {
// 		return dto.ErrEmailNotFound
// 	}

// 	draftEmail, err := makeVerificationEmail(user.Email)
// 	if err != nil {
// 		return err
// 	}

// 	err = utils.SendMail(user.Email, draftEmail["subject"], draftEmail["body"])
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (s *userService) VerifyEmail(ctx context.Context, req dto.VerifyEmailRequest) (dto.VerifyEmailResponse, error) {
// 	decryptedToken, err := utils.AESDecrypt(req.Token)
// 	if err != nil {
// 		return dto.VerifyEmailResponse{}, dto.ErrTokenInvalid
// 	}

// 	if !strings.Contains(decryptedToken, "_") {
// 		return dto.VerifyEmailResponse{}, dto.ErrTokenInvalid
// 	}

// 	decryptedTokenSplit := strings.Split(decryptedToken, "_")
// 	email := decryptedTokenSplit[0]
// 	expired := decryptedTokenSplit[1]

// 	now := time.Now()
// 	expiredTime, err := time.Parse("2006-01-02 15:04:05", expired)
// 	if err != nil {
// 		return dto.VerifyEmailResponse{}, dto.ErrTokenInvalid
// 	}

// 	if expiredTime.Sub(now) < 0 {
// 		return dto.VerifyEmailResponse{
// 			Email:      email,
// 			IsVerified: false,
// 		}, dto.ErrTokenExpired
// 	}

// 	user, err := s.userRepo.GetUserByEmail(ctx, nil, email)
// 	if err != nil {
// 		return dto.VerifyEmailResponse{}, dto.ErrUserNotFound
// 	}

// 	if user.IsVerified {
// 		return dto.VerifyEmailResponse{}, dto.ErrAccountAlreadyVerified
// 	}

// 	updatedUser, err := s.userRepo.UpdateUser(ctx, nil, entity.User{
// 		ID:         user.ID,
// 		IsVerified: true,
// 	})
// 	if err != nil {
// 		return dto.VerifyEmailResponse{}, dto.ErrUpdateUser
// 	}

// 	return dto.VerifyEmailResponse{
// 		Email:      email,
// 		IsVerified: updatedUser.IsVerified,
// 	}, nil
// }
