package handler

import (
	"github.com/MrVirul/icecream-inventory/internal/dto"
	"github.com/MrVirul/icecream-inventory/internal/httpx"
	"github.com/MrVirul/icecream-inventory/internal/middleware"
	"github.com/MrVirul/icecream-inventory/internal/service"
	"github.com/gofiber/fiber/v2"
)

type CustomerHandler struct {
	customers service.CustomerService
}

func NewCustomerHandler(customers service.CustomerService) *CustomerHandler {
	return &CustomerHandler{customers: customers}
}

func (h *CustomerHandler) List(c *fiber.Ctx) error {
	search := c.Query("search")
	allowed := map[string]bool{"name": true, "created_at": true}
	params := httpx.ParseListParams(c, allowed, "name")
	items, meta, err := h.customers.List(search, params.Page, params.Limit, params.Sort, params.Desc)
	if err != nil {
		return err
	}
	return httpx.OKList(c, items, meta)
}

func (h *CustomerHandler) Get(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return err
	}
	item, err := h.customers.Get(id)
	if err != nil {
		return err
	}
	return httpx.OK(c, item)
}

func (h *CustomerHandler) Create(c *fiber.Ctx) error {
	var req dto.CustomerCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return httpx.BadRequest("bad request")
	}
	if err := httpx.ValidateStruct(req); err != nil {
		return err
	}
	// phone validation
	if !isValidPhone(req.Phone) {
		return httpx.ValidationError([]httpx.FieldError{{Field: "phone", Message: "must contain only digits and + - ( ) space"}})
	}
	uid := c.Locals(middleware.CtxUserID).(uint)
	item, err := h.customers.Create(req, uid)
	if err != nil {
		return err
	}
	return httpx.Created(c, item)
}

func (h *CustomerHandler) Replace(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return err
	}
	var req dto.CustomerUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return httpx.BadRequest("bad request")
	}
	if err := httpx.ValidateStruct(req); err != nil {
		return err
	}
	if !isValidPhone(req.Phone) {
		return httpx.ValidationError([]httpx.FieldError{{Field: "phone", Message: "must contain only digits and + - ( ) space"}})
	}
	uid := c.Locals(middleware.CtxUserID).(uint)
	item, err := h.customers.Replace(id, req, uid)
	if err != nil {
		return err
	}
	return httpx.OK(c, item)
}

func (h *CustomerHandler) Patch(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return err
	}
	var req dto.CustomerPatchRequest
	if err := c.BodyParser(&req); err != nil {
		return httpx.BadRequest("bad request")
	}
	if req.Name == nil && req.Email == nil && req.Phone == nil && req.Address == nil && req.Notes == nil {
		return httpx.ValidationError([]httpx.FieldError{{Field: "body", Message: "at least one field is required"}})
	}
	if err := httpx.ValidateStruct(req); err != nil {
		return err
	}
	if req.Phone != nil && !isValidPhone(*req.Phone) {
		return httpx.ValidationError([]httpx.FieldError{{Field: "phone", Message: "must contain only digits and + - ( ) space"}})
	}
	uid := c.Locals(middleware.CtxUserID).(uint)
	item, err := h.customers.Patch(id, req, uid)
	if err != nil {
		return err
	}
	return httpx.OK(c, item)
}

func (h *CustomerHandler) Delete(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return err
	}
	if err := h.customers.Delete(id); err != nil {
		return err
	}
	return httpx.NoContent(c)
}

func isValidPhone(p string) bool {
	if len(p) < 7 || len(p) > 30 {
		return false
	}
	for _, r := range p {
		if (r >= '0' && r <= '9') || r == '+' || r == '-' || r == '(' || r == ')' || r == ' ' {
			continue
		}
		return false
	}
	return true
}
