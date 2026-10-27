package repository

import (
	"strings"

	"github.com/MrVirul/icecream-inventory/internal/models"
	"gorm.io/gorm"
)

type ItemRepository interface {
	List(search, category string, lowStock bool, page, limit int, sort string, desc bool) ([]models.InventoryItem, int64, error)
	GetByID(id uint) (*models.InventoryItem, error)
	GetByName(name string) (*models.InventoryItem, error)
	Create(item *models.InventoryItem) error
	Update(item *models.InventoryItem) error
	Delete(id uint) error
}

type itemRepository struct {
	db *gorm.DB
}

func NewItemRepository(db *gorm.DB) ItemRepository {
	return &itemRepository{db: db}
}

func (r *itemRepository) List(search, category string, lowStock bool, page, limit int, sort string, desc bool) ([]models.InventoryItem, int64, error) {
	var items []models.InventoryItem
	var total int64

	q := r.db.Model(&models.InventoryItem{})
	if search != "" {
		s := strings.ToLower(search)
		// escape % and _
		s = strings.ReplaceAll(s, "\\", "\\\\")
		s = strings.ReplaceAll(s, "%", "\\%")
		s = strings.ReplaceAll(s, "_", "\\_")
		q = q.Where("LOWER(name) LIKE ? ESCAPE '\\'", "%"+s+"%")
	}
	if category != "" {
		q = q.Where("category = ?", category)
	}
	if lowStock {
		q = q.Where("quantity <= reorder_level")
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

func (r *itemRepository) GetByID(id uint) (*models.InventoryItem, error) {
	var item models.InventoryItem
	if err := r.db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *itemRepository) GetByName(name string) (*models.InventoryItem, error) {
	var item models.InventoryItem
	if err := r.db.Where("LOWER(name) = LOWER(?)", strings.TrimSpace(name)).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *itemRepository) Create(item *models.InventoryItem) error {
	return r.db.Create(item).Error
}

func (r *itemRepository) Update(item *models.InventoryItem) error {
	return r.db.Save(item).Error
}

func (r *itemRepository) Delete(id uint) error {
	return r.db.Delete(&models.InventoryItem{}, id).Error
}
