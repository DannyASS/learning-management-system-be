package user_usecase

import (
	"errors"
	"fmt"
	"time"

	"github.com/DannyAss/users/config"
	"github.com/DannyAss/users/internal/database"
	user_model "github.com/DannyAss/users/internal/models/database_model/user"
	auth_usercase "github.com/DannyAss/users/internal/repositories/auth"
	users_repository "github.com/DannyAss/users/internal/repositories/users"
	"github.com/DannyAss/users/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

type userUsecase struct {
	userRepo users_repository.IuserRepos
	cfg      *config.ConfigEnv
	crypto   *utils.CryptoService
	db       *database.DBManager
}

type IuserUsecae interface {
	Register(register user_model.RegisterDTO) error
	Login(login user_model.LoginDTO) (*user_model.LoginResponse, *fiber.Cookie, error)
	RefreshToken(userToken string) (*user_model.LoginResponse, error)
	GetUsers(req user_model.Pagination) (*map[string]interface{}, error)
	Updateuser(req user_model.RequestUserGetList) error
	Logout(refreshToken string) (*fiber.Cookie, error)
	GetAllTeacher() ([]map[string]interface{}, error)
}

func NewUserCaseUser(userRepos users_repository.IuserRepos, cfg *config.ConfigEnv, db *database.DBManager) IuserUsecae {
	fmt.Println("cek cfg apke key :", cfg.AppKey)
	cs, err := utils.NewCryptoService([]byte(cfg.AppKey))
	if err != nil {
		panic(err)
	}
	return &userUsecase{userRepo: userRepos, cfg: cfg, crypto: cs, db: db}
}

func (u *userUsecase) Register(register user_model.RegisterDTO) error {
	tx := u.db.GetDB().Begin()

	var errTx error // <-- FLAG ERROR

	defer func() {
		if errTx != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	// create repo with transaction
	userRepo := u.userRepo.WithTx(tx)
	// cek username eksis
	filterUser := user_model.FilterUser{
		Username: register.Username,
	}
	_, UEerr := userRepo.GetUser(filterUser, "")
	if UEerr == nil {
		errTx = errors.New("username already exists")
		return errTx
	}

	if register.Password == "" {
		errTx = errors.New("password is mandatory")
		return errTx
	}
	// cek role
	filterRole := user_model.RoleFilter{
		Name: register.Role,
	}

	role, Rerr := userRepo.GetRole(filterRole, "")
	if Rerr != nil {
		errTx = errors.New("role doesn't exists")
		return errTx
	}

	// encript semua data credential
	email, _ := u.crypto.Encrypt([]byte(register.Email))
	phone, _ := u.crypto.Encrypt([]byte(register.Phone))
	password, _ := u.crypto.HashString(register.Password)
	// insert user
	userModel := user_model.User{
		Name:       register.Name,
		Email:      email,
		Username:   register.Username,
		Phone:      phone,
		Password:   password,
		IsVerified: false,
		Gender:     register.Gender,
	}

	user, CUerr := userRepo.CreateUser(userModel)
	if CUerr != nil {
		errTx = CUerr
		return errTx
	}

	// insert user role
	userRole := user_model.UserRole{
		UserID: user.ID,
		RoleID: role.ID,
	}

	_, USerr := userRepo.CreateUserRole(userRole)

	if USerr != nil {
		errTx = USerr
		return errTx
	}

	return nil
}

func (u *userUsecase) Login(login user_model.LoginDTO) (*user_model.LoginResponse, *fiber.Cookie, error) {
	tx := u.db.GetDB().Begin()

	var errTx error // <-- FLAG ERROR

	defer func() {
		if errTx != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	userRepos := u.userRepo.WithTx(tx)

	// GET USER
	user, err := userRepos.GetUser(user_model.FilterUser{
		Username: login.Username,
	}, "")
	if err != nil {
		errTx = errors.New("account doesn't exist")
		return nil, nil, errTx
	}

	// GET ROLE
	roleId, err := userRepos.RoleUserWithRole(user_model.FilterUser{
		Id: user.ID,
	})
	if err != nil {
		errTx = err
		return nil, nil, errTx
	}

	// VERIFY PASSWORD
	_, err = u.crypto.VerifyString(login.Password, user.Password)
	if err != nil {
		errTx = errors.New("invalid password")
		return nil, nil, errTx
	}

	// DECRYPT EMAIL
	uEmail, err := u.crypto.Decrypt(user.Email)
	if err != nil {
		errTx = err
		return nil, nil, errTx
	}

	// GENERATE ACCESS TOKEN
	accessTTL := 15 * time.Minute
	rolesId := roleId.ID

	access, _, err := auth_usercase.GenerateAccessToken(
		[]byte(u.cfg.JWTSecretKey),
		uint(user.ID),
		string(uEmail),
		user.Username,
		uint(rolesId),
		accessTTL,
	)
	if err != nil {
		errTx = err
		return nil, nil, errTx
	}

	// GENERATE REFRESH TOKEN
	refreshToken, _, err := auth_usercase.GenerateRefreshToken()
	if err != nil {
		errTx = err
		return nil, nil, errTx
	}

	cookies := utils.CookieConfig(u.cfg.AppEnv, u.cfg.AppDomain, refreshToken)

	// SAVE REFRESH TOKEN
	err = userRepos.CreateUserTokens(user_model.UserToken{
		UserID:       uint(user.ID),
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour),
	})
	if err != nil {
		errTx = err
		return nil, nil, errTx
	}

	// ENCRYPT ACCESS TOKEN
	chipperAccess, err := u.crypto.Encrypt([]byte(access))
	if err != nil {
		errTx = err
		return nil, nil, errTx
	}

	// SUCCESS: errTx tetap nil → Commit di defer

	resp := user_model.LoginResponse{
		Token: chipperAccess,
		User: user_model.UserDTO{
			Name:     user.Name,
			Email:    string(uEmail),
			Username: user.Username,
			Role:     roleId,
		},
	}

	return &resp, &cookies, nil
}

func (u *userUsecase) RefreshToken(userToken string) (*user_model.LoginResponse, error) {
	filterToken := user_model.FilterUserToken{
		RefreshToken: userToken,
		Revoked:      false,
	}
	getUserId, err := u.userRepo.GetRefreshToken(filterToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	if time.Now().After(getUserId.ExpiresAt) {
		return nil, errors.New("session expired")
	}

	filterUser := user_model.FilterUser{
		Id: uint64(getUserId.UserID),
	}

	user, err1 := u.userRepo.GetUser(filterUser, "")
	if err1 != nil {
		return nil, err1
	}
	email, _ := u.crypto.Decrypt(user.Email)

	var rolesId uint

	filterUser = user_model.FilterUser{
		UserId: user.ID,
	}

	getRolesId, err2 := u.userRepo.GetListUserRoles(filterUser)

	if err2 != nil {
		return nil, err2
	}

	rolesId = uint(getRolesId[0].ID)

	accessTTL := time.Duration(15) * time.Minute
	accessToken, _, err3 := auth_usercase.GenerateAccessToken([]byte(u.cfg.JWTSecretKey), uint(user.ID), string(email), user.Username, rolesId, accessTTL)
	if err3 != nil {
		return nil, err3
	}

	chipperAccess, _ := u.crypto.Encrypt([]byte(accessToken))

	response := user_model.LoginResponse{
		Token: chipperAccess,
	}

	return &response, nil
}

func (u *userUsecase) GetUsers(req user_model.Pagination) (*map[string]interface{}, error) {
	if req.Page < 1 {
		req.Page = 1
	}

	if req.Perpage < 1 {
		req.Perpage = 10
	}

	data, page, err := u.userRepo.GetUserAndRole(req)
	if err != nil {
		return nil, err
	}

	for c := range data {
		Email, err1 := u.crypto.Decrypt(data[c].Email)
		if err1 != nil {
			data[c].Email = ""
		} else {
			data[c].Email = string(Email)
		}

		Phone, err2 := u.crypto.Decrypt(data[c].Phone)
		if err2 != nil {
			data[c].Phone = ""
		} else {
			data[c].Phone = string(Phone)
		}
	}

	response := map[string]interface{}{
		"data":       data,
		"page":       page.Page,
		"perpage":    page.Perpage,
		"total_data": page.TotalData,
		"total_page": page.TotalPage,
	}

	return &response, nil
}

func (u *userUsecase) Updateuser(req user_model.RequestUserGetList) error {
	tx := u.db.GetDB().Begin()

	if tx.Error != nil {
		return tx.Error
	}

	userRepos := u.userRepo.WithTx(tx)

	defer func() {
		if tx.Error != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	Eemail, err := u.crypto.Encrypt([]byte(req.Email))
	if err != nil {
		return err
	}

	Ephone, err1 := u.crypto.Encrypt([]byte(req.Phone))
	if err1 != nil {
		return err1
	}

	model := user_model.User{
		Name:  req.Name,
		Email: Eemail,
		Phone: Ephone,
	}

	filter := user_model.FilterUser{
		Id: uint64(req.Id),
	}

	err2 := userRepos.UpdateUser(model, filter)
	if err2 != nil {
		return err2
	}

	err3 := userRepos.UpdateUserRole(req.RoleId, filter)
	if err3 != nil {
		return err3
	}

	return nil
}

func (u *userUsecase) Logout(refreshToken string) (*fiber.Cookie, error) {
	if refreshToken == "" {
		return nil, errors.New("refresh token not found")
	}

	// revoke token
	err := u.userRepo.RevokeToken(refreshToken)

	cookies := utils.CookieConfig(u.cfg.AppEnv, u.cfg.AppDomain, "")

	if err != nil {
		return nil, err
	}

	return &cookies, nil
}

func (u *userUsecase) GetAllTeacher() ([]map[string]interface{}, error) {
	var result []map[string]interface{}

	data, err := u.userRepo.AllGetTeacher()

	if err != nil {
		return nil, err
	}

	for _, c := range data {
		dto := map[string]interface{}{
			"teacher_id": c.ID,
			"name":       c.Name,
		}

		result = append(result, dto)
	}

	return result, nil
}
