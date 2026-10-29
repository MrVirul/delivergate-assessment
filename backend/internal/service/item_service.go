package service

import (
	"errors"
	"strings"

	"github.com/MrVirul/icecream-inventory/internal/dto"
	"github.com/MrVirul/icecream-inventory/internal/httpx"
	"github.com/MrVirul/icecream-inventory/internal/models"
	"github.com/MrVirul/icecream-inventory/internal/repository"
	"gorm.io/gorm"
)

type ItemService interface {
	List(search, category string, lowStock bool, page, limit int, sort string, desc bool) ([]dto.ItemResponse, httpx.Meta, error)
	Get(id uint) (*dto.ItemResponse, error)
	Create(req dto.ItemCreateRequest, userID uint) (*dto.ItemResponse, error)
	Replace(id uint, req dto.ItemUpdateRequest, userID uint) (*dto.ItemResponse, error)
	Patch(id uint, req dto.ItemPatchRequest, userID uint) (*dto.ItemResponse, error)
	Delete(id uint) error
}

type itemService struct {
	items repository.ItemRepository
}

func NewItemService(items repository.ItemRepository) ItemService {
	return &itemService{items: items}
}

func toItemResponse(m models.InventoryItem) dto.ItemResponse {
	return dto.ItemResponse{
		ID:           m.ID,
		Name:         m.Name,
		Category:     m.Category,
		Description:  m.Description,
		Quantity:     m.Quantity,
		Unit:         m.Unit,
		PriceCents:   m.PriceCents,
		ReorderLevel: m.ReorderLevel,
		LowStock:     m.Quantity <= m.ReorderLevel,
		CreatedByID:  m.CreatedByID,
		UpdatedByID:  m.UpdatedByID,
		CreatedAt:    m.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:    m.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func (s *itemService) List(search, category string, lowStock bool, page, limit int, sort string, desc bool) ([]dto.ItemResponse, httpx.Meta, error) {
	items, total, err := s.items.List(search, category, lowStock, page, limit, sort, desc)
	if err != nil {
		return nil, httpx.Meta{}, err
	}
	res := make([]dto.ItemResponse, 0, len(items))
	for _, it := range items {
		res = append(res, toItemResponse(it))
	}
	meta := httpx.Meta{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: httpx.TotalPages(total, limit),
	}
	return res, meta, nil
}

func (s *itemService) Get(id uint) (*dto.ItemResponse, error) {
	item, err := s.items.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, httpx.NotFound("item")
		}
		return nil, err
	}
	r := toItemResponse(*item)
	return &r, nil
}

func (s *itemService) Create(req dto.ItemCreateRequest, userID uint) (*dto.ItemResponse, error) {
	name := strings.TrimSpace(req.Name)
	// duplicate check
	if existing, err := s.items.GetByName(name); err == nil && existing != nil {
		return nil, httpx.Conflict("item name already exists")
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	qty := 0
	if req.Quantity != nil {
		qty = *req.Quantity
	}
	pc := int64(0)
	if req.PriceCents != nil {
		pc = *req.PriceCents
	}
	rl := 0
	if req.ReorderLevel != nil {
		rl = *req.ReorderLevel
	}

	item := &models.InventoryItem{
		Name:         name,
		Category:     req.Category,
		Description:  strings.TrimSpace(req.Description),
		Quantity:     qty,
		Unit:         req.Unit,
		PriceCents:   pc,
		ReorderLevel: rl,
		CreatedByID:  &userID,
		UpdatedByID:  &userID,
	}
	if err := s.items.Create(item); err != nil {
		// try to map common unique errors
		if strings.Contains(strings.ToLower(err.Error()), "unique") || strings.Contains(strings.ToLower(err.Error()), "duplicate") {
			return nil, httpx.Conflict("item name already exists")
		}
		return nil, err
	}
	r := toItemResponse(*item)
	return &r, nil
}

func (s *itemService) Replace(id uint, req dto.ItemUpdateRequest, userID uint) (*dto.ItemResponse, error) {
	item, err := s.items.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, httpx.NotFound("item")
		}
		return nil, err
	}
	name := strings.TrimSpace(req.Name)
	if existing, err := s.items.GetByName(name); err == nil && existing != nil && existing.ID != id {
		return nil, httpx.Conflict("item name already exists")
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	qty := 0
	if req.Quantity != nil {
		qty = *req.Quantity
	}
	pc := int64(0)
	if req.PriceCents != nil {
		pc = *req.PriceCents
	}
	rl := 0
	if req.ReorderLevel != nil {
		rl = *req.ReorderLevel
	}

	item.Name = name
	item.Category = req.Category
	item.Description = strings.TrimSpace(req.Description)
	item.Quantity = qty
	item.Unit = req.Unit
	item.PriceCents = pc
	item.ReorderLevel = rl
	item.UpdatedByID = &userID

	if err := s.items.Update(item); err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "unique") || strings.Contains(strings.ToLower(err.Error()), "duplicate") {
			return nil, httpx.Conflict("item name already exists")
		}
		return nil, err
	}
	r := toItemResponse(*item)
	return &r, nil
}

func (s *itemService) Patch(id uint, req dto.ItemPatchRequest, userID uint) (*dto.ItemResponse, error) {
	item, err := s.items.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, httpx.NotFound("item")
		}
		return nil, err
	}
	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if existing, err := s.items.GetByName(name); err == nil && existing != nil && existing.ID != id {
			return nil, httpx.Conflict("item name already exists")
		} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		item.Name = name
	}
	if req.Category != nil {
		item.Category = *req.Category
	}
	if req.Description != nil {
		item.Description = strings.TrimSpace(*req.Description)
	}
	if req.Quantity != nil {
		item.Quantity = *req.Quantity
	}
	if req.Unit != nil {
		item.Unit = *req.Unit
	}
	if req.PriceCents != nil {
		item.PriceCents = *req.PriceCents
	}
	if req.ReorderLevel != nil {
		item.ReorderLevel = *req.ReorderLevel
	}
	item.UpdatedByID = &userID

	if err := s.items.Update(item); err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "unique") || strings.Contains(strings.ToLower(err.Error()), "duplicate") {
			return nil, httpx.Conflict("item name already exists")
		}
		return nil, err
	}
	r := toItemResponse(*item)
	return &r, nil
}

func (s *itemService) Delete(id uint) error {
	item, err := s.items.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return httpx.NotFound("item")
		}
		return err
	}
	if err := s.items.Delete(item.ID); err != nil {
		return err
	}
	return nil
}
