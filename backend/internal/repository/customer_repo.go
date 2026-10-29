package repository

import (
	"strings"

	"github.com/MrVirul/icecream-inventory/internal/models"
	"gorm.io/gorm"
)

type CustomerRepository interface {
	List(search string, page, limit int, sort string, desc bool) ([]models.Customer, int64, error)
	GetByID(id uint) (*models.Customer, error)
	GetByEmail(email string) (*models.Customer, error)
	Create(cust *models.Customer) error
	Update(cust *models.Customer) error
	Delete(id uint) error
}

type customerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) CustomerRepository {
	return &customerRepository{db: db}
}

func (r *customerRepository) List(search string, page, limit int, sort string, desc bool) ([]models.Customer, int64, error) {
	var items []models.Customer
	var total int64

	q := r.db.Model(&models.Customer{})
	if search != "" {
		s := strings.ToLower(search)
		s = strings.ReplaceAll(s, "\\", "\\\\")
		s = strings.ReplaceAll(s, "%", "\\%")
		s = strings.ReplaceAll(s, "_", "\\_")
		like := "%" + s + "%"
		q = q.Where("LOWER(name) LIKE ? ESCAPE '\\' OR LOWER(COALESCE(email,'')) LIKE ? ESCAPE '\\' OR phone LIKE ?", like, like, "%"+search+"%")
	}

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if sort == "" {
		sort = "name"
	}
	order := sort
	if desc {
		order = sort + " DESC"
	} else {
		order = sort + " ASC"
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	offset := (page - 1) * limit

	if err := q.Order(order).Offset(offset).Limit(limit).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *customerRepository) GetByID(id uint) (*models.Customer, error) {
	var cust models.Customer
	if err := r.db.First(&cust, id).Error; err != nil {
		return nil, err
	}
	return &cust, nil
}

func (r *customerRepository) GetByEmail(email string) (*models.Customer, error) {
	if strings.TrimSpace(email) == "" {
		return nil, gorm.ErrRecordNotFound
	}
	var cust models.Customer
	if err := r.db.Where("LOWER(email) = LOWER(?)", strings.TrimSpace(email)).First(&cust).Error; err != nil {
		return nil, err
	}
	return &cust, nil
}

func (r *customerRepository) Create(cust *models.Customer) error {
	return r.db.Create(cust).Error
}

func (r *customerRepository) Update(cust *models.Customer) error {
	return r.db.Save(cust).Error
}

func (r *customerRepository) Delete(id uint) error {
	return r.db.Delete(&models.Customer{}, id).Error
}
