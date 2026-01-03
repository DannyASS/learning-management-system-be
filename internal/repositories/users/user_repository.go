package users_repository

import (
	"math"
	"strings"

	"github.com/DannyAss/users/internal/database"
	user_model "github.com/DannyAss/users/internal/models/database_model/user"
	"gorm.io/gorm"
)

type userRepos struct {
	db *database.DBManager
	tx *gorm.DB
}

type IuserRepos interface {
	WithTx(tx *gorm.DB) IuserRepos
	CreateUser(user user_model.User) (*user_model.User, error)
	GetlistUser(user user_model.FilterUser, mode string) ([]user_model.User, error)
	GetUser(user user_model.FilterUser, mode string) (*user_model.User, error)
	GetRole(filter user_model.RoleFilter, mode string) (*user_model.Role, error)
	CreateRole(role user_model.Role) (*user_model.Role, error)
	CreateUserRole(models user_model.UserRole) (*user_model.UserRole, error)
	RoleUserWithRole(filter user_model.FilterUser) (*user_model.Role, error)
	CreateUserTokens(refrestoken user_model.UserToken) error
	GetListUserRoles(filter user_model.FilterUser) ([]user_model.UserRole, error)
	GetRefreshToken(filter user_model.FilterUserToken) (*user_model.UserToken, error)
	GetUserAndRole(page user_model.Pagination) ([]user_model.RequestUserGetList, *user_model.Pagination, error)
	UpdateUser(model user_model.User, filter user_model.FilterUser) error
	UpdateUserRole(role_id int, filter user_model.FilterUser) error
	RevokeToken(token string) error
	AllGetTeacher() ([]user_model.User, error)
}

func NewReposUser(db *database.DBManager) IuserRepos {
	return &userRepos{db: db}
}

func (u *userRepos) WithTx(tx *gorm.DB) IuserRepos {
	return &userRepos{db: u.db, tx: tx}
}

func (u *userRepos) getDB() *gorm.DB {
	if u.tx != nil {
		return u.tx
	}

	return u.db.GetDB()
}

func (u *userRepos) CreateUser(user user_model.User) (*user_model.User, error) {
	tx := u.getDB()

	err := tx.Create(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *userRepos) GetlistUser(user user_model.FilterUser, mode string) ([]user_model.User, error) {
	query := u.getDB()
	var users []user_model.User
	if strings.ToLower(mode) == "full" {
		query = query.Where(user)
	} else {
		query = query.Where(&user)
	}

	err := query.Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (u *userRepos) GetUser(user user_model.FilterUser, mode string) (*user_model.User, error) {
	query := u.getDB().Debug()
	var users user_model.User
	if strings.ToLower(mode) == "full" {
		query = query.Where(user)
	} else {
		query = query.Where(&user)
	}

	err := query.First(&users).Error
	if err != nil {
		return nil, err
	}

	return &users, nil
}

func (u *userRepos) GetListRole(role user_model.RoleFilter, mode string) ([]user_model.Role, error) {
	tx := u.getDB()

	var roles []user_model.Role

	if strings.ToLower(mode) == "full" {
		tx = tx.Where(role)
	} else {
		tx = tx.Where(&role)
	}

	err := tx.Find(&roles).Error
	if err != nil {
		return nil, err
	}

	return roles, nil
}

func (u *userRepos) GetRole(filter user_model.RoleFilter, mode string) (*user_model.Role, error) {
	tx := u.getDB()
	var role user_model.Role
	if strings.ToLower(mode) == "full" {
		tx = tx.Where(filter)
	} else {
		tx = tx.Where(&filter)
	}

	err := tx.First(&role).Error
	if err != nil {
		return nil, err
	}

	return &role, nil
}

func (u *userRepos) CreateRole(role user_model.Role) (*user_model.Role, error) {
	tx := u.getDB()

	err := tx.Create(&role).Error
	if err != nil {
		return nil, err
	}

	return &role, nil
}

func (u *userRepos) CreateUserRole(models user_model.UserRole) (*user_model.UserRole, error) {
	tx := u.getDB()

	err := tx.Create(&models).Error
	if err != nil {
		return nil, err
	}

	return &models, nil
}

func (u *userRepos) RoleUserWithRole(filter user_model.FilterUser) (*user_model.Role, error) {
	tx := u.getDB()

	var result user_model.Role

	err := tx.Table("users as a").
		Select("c.*").
		Joins("join user_roles as b on b.user_id = a.id").
		Joins("join roles as c on c.id = b.role_id").
		Where("a.id = ?", filter.Id).
		First(&result).Error
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (u *userRepos) CreateUserTokens(refrestoken user_model.UserToken) error {
	if err := u.getDB().Create(&refrestoken).Error; err != nil {
		return err
	}

	return nil
}

func (u *userRepos) GetListUserRoles(filter user_model.FilterUser) ([]user_model.UserRole, error) {
	tx := u.getDB()
	var userRole []user_model.UserRole

	if err := tx.Where(&filter).Find(&userRole).Error; err != nil {
		return nil, err
	}

	return userRole, nil
}

func (u *userRepos) GetRefreshToken(filter user_model.FilterUserToken) (*user_model.UserToken, error) {
	tx := u.getDB()
	var userRole user_model.UserToken

	if err := tx.Where(&filter).First(&userRole).Error; err != nil {
		return nil, err
	}

	return &userRole, nil
}

func (u *userRepos) GetUserAndRole(page user_model.Pagination) ([]user_model.RequestUserGetList, *user_model.Pagination, error) {
	db := u.getDB()

	var totalData int64
	data := []user_model.RequestUserGetList{}

	// Base query (tanpa limit)
	baseQuery := db.Table("users a").
		Select(`
			a.id as id,
			a.name as name,
			a.username as username,
			a.phone as phone,
			a.email as email,
			c.description as role_name,
			c.id as role_id
		`).
		Joins("join user_roles b on a.id = b.user_id").
		Joins("join roles c on c.id = b.role_id")

	// Search (jika ada)
	if page.Search != "" {
		baseQuery = baseQuery.Where("a.name ILIKE ?", "%"+page.Search+"%")
	}

	// Count
	if err := baseQuery.Count(&totalData).Error; err != nil {
		return nil, nil, err
	}

	// Pagination
	if err := baseQuery.
		Limit(page.Perpage).
		Offset((page.Page - 1) * page.Perpage).
		Scan(&data).Error; err != nil {
		return nil, nil, err
	}

	page.TotalData = int(totalData)
	page.TotalPage = int(math.Ceil(float64(totalData) / float64(page.Perpage)))

	return data, &page, nil
}

func (u *userRepos) UpdateUser(model user_model.User, filter user_model.FilterUser) error {
	tx := u.getDB()

	if err := tx.Where(&filter).Updates(&model).Error; err != nil {
		return err
	}

	return nil
}

func (u *userRepos) UpdateUserRole(role_id int, filter user_model.FilterUser) error {
	tx := u.getDB().Model(&user_model.UserRole{})

	if err := tx.Where("user_id = ?", filter.Id).Update("role_id = ?", role_id).Error; err != nil {
		return err
	}

	return nil
}

func (r *userRepos) RevokeToken(token string) error {
	return r.getDB().Model(&user_model.UserToken{}).
		Where("refresh_token = ?", token).
		Update("revoked", 1).Error
}

func (r *userRepos) AllGetTeacher() ([]user_model.User, error) {
	var user []user_model.User
	tx := r.getDB()

	tx = tx.Table("users as a").
		Joins("join user_roles as b on b.user_id = a.id").
		Joins("join roles as c on b.role_id = c.id")

	tx = tx.Select(`
		distinct
		a.name as name,
		a.username as username,
		a.id as id
	`)

	if err := tx.Where("c.name = ?", "TC").Scan(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}
