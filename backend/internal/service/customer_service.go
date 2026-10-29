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

type CustomerService interface {
	List(search string, page, limit int, sort string, desc bool) ([]dto.CustomerResponse, httpx.Meta, error)
	Get(id uint) (*dto.CustomerResponse, error)
	Create(req dto.CustomerCreateRequest, userID uint) (*dto.CustomerResponse, error)
	Replace(id uint, req dto.CustomerUpdateRequest, userID uint) (*dto.CustomerResponse, error)
	Patch(id uint, req dto.CustomerPatchRequest, userID uint) (*dto.CustomerResponse, error)
	Delete(id uint) error
}

type customerService struct {
	customers repository.CustomerRepository
}

func NewCustomerService(customers repository.CustomerRepository) CustomerService {
	return &customerService{customers: customers}
}

func toCustomerResponse(m models.Customer) dto.CustomerResponse {
	var email *string
	if m.Email != nil && *m.Email != "" {
		e := *m.Email
		email = &e
	}
	return dto.CustomerResponse{
		ID:          m.ID,
		Name:        m.Name,
		Email:       email,
		Phone:       m.Phone,
		Address:     m.Address,
		Notes:       m.Notes,
		CreatedByID: m.CreatedByID,
		UpdatedByID: m.UpdatedByID,
		CreatedAt:   m.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   m.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func (s *customerService) List(search string, page, limit int, sort string, desc bool) ([]dto.CustomerResponse, httpx.Meta, error) {
	items, total, err := s.customers.List(search, page, limit, sort, desc)
	if err != nil {
		return nil, httpx.Meta{}, err
	}
	res := make([]dto.CustomerResponse, 0, len(items))
	for _, it := range items {
		res = append(res, toCustomerResponse(it))
	}
	meta := httpx.Meta{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: httpx.TotalPages(total, limit),
	}
	return res, meta, nil
}

func (s *customerService) Get(id uint) (*dto.CustomerResponse, error) {
	cust, err := s.customers.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, httpx.NotFound("customer")
		}
		return nil, err
	}
	r := toCustomerResponse(*cust)
	return &r, nil
}

func (s *customerService) Create(req dto.CustomerCreateRequest, userID uint) (*dto.CustomerResponse, error) {
	name := strings.TrimSpace(req.Name)
	var emailPtr *string
	if req.Email != nil && strings.TrimSpace(*req.Email) != "" {
		e := strings.ToLower(strings.TrimSpace(*req.Email))
		emailPtr = &e
		if existing, err := s.customers.GetByEmail(e); err == nil && existing != nil {
			return nil, httpx.Conflict("customer email already exists")
		} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	cust := &models.Customer{
		Name:        name,
		Email:       emailPtr,
		Phone:       strings.TrimSpace(req.Phone),
		Address:     strings.TrimSpace(req.Address),
		Notes:       strings.TrimSpace(req.Notes),
		CreatedByID: &userID,
		UpdatedByID: &userID,
	}
	if err := s.customers.Create(cust); err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "unique") || strings.Contains(strings.ToLower(err.Error()), "duplicate") {
			return nil, httpx.Conflict("customer email already exists")
		}
		return nil, err
	}
	r := toCustomerResponse(*cust)
	return &r, nil
}

func (s *customerService) Replace(id uint, req dto.CustomerUpdateRequest, userID uint) (*dto.CustomerResponse, error) {
	cust, err := s.customers.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, httpx.NotFound("customer")
		}
		return nil, err
	}
	name := strings.TrimSpace(req.Name)
	var emailPtr *string
	if req.Email != nil && strings.TrimSpace(*req.Email) != "" {
		e := strings.ToLower(strings.TrimSpace(*req.Email))
		emailPtr = &e
		if existing, err := s.customers.GetByEmail(e); err == nil && existing != nil && existing.ID != id {
			return nil, httpx.Conflict("customer email already exists")
		} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	cust.Name = name
	cust.Email = emailPtr
	cust.Phone = strings.TrimSpace(req.Phone)
	cust.Address = strings.TrimSpace(req.Address)
	cust.Notes = strings.TrimSpace(req.Notes)
	cust.UpdatedByID = &userID

	if err := s.customers.Update(cust); err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "unique") || strings.Contains(strings.ToLower(err.Error()), "duplicate") {
			return nil, httpx.Conflict("customer email already exists")
		}
		return nil, err
	}
	r := toCustomerResponse(*cust)
	return &r, nil
}

func (s *customerService) Patch(id uint, req dto.CustomerPatchRequest, userID uint) (*dto.CustomerResponse, error) {
	cust, err := s.customers.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, httpx.NotFound("customer")
		}
		return nil, err
	}
	if req.Name != nil {
		cust.Name = strings.TrimSpace(*req.Name)
	}
	if req.Email != nil {
		if strings.TrimSpace(*req.Email) == "" {
			cust.Email = nil
		} else {
			e := strings.ToLower(strings.TrimSpace(*req.Email))
			if existing, err := s.customers.GetByEmail(e); err == nil && existing != nil && existing.ID != id {
				return nil, httpx.Conflict("customer email already exists")
			} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, err
			}
			cust.Email = &e
		}
	}
	if req.Phone != nil {
		cust.Phone = strings.TrimSpace(*req.Phone)
	}
	if req.Address != nil {
		cust.Address = strings.TrimSpace(*req.Address)
	}
	if req.Notes != nil {
		cust.Notes = strings.TrimSpace(*req.Notes)
	}
	cust.UpdatedByID = &userID

	if err := s.customers.Update(cust); err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "unique") || strings.Contains(strings.ToLower(err.Error()), "duplicate") {
			return nil, httpx.Conflict("customer email already exists")
		}
		return nil, err
	}
	r := toCustomerResponse(*cust)
	return &r, nil
}

func (s *customerService) Delete(id uint) error {
	cust, err := s.customers.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return httpx.NotFound("customer")
		}
		return err
	}
	if err := s.customers.Delete(cust.ID); err != nil {
		return err
	}
	return nil
}
