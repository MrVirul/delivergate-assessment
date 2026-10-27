package handler

import (
	"strconv"
	"strings"

	"github.com/MrVirul/icecream-inventory/internal/dto"
	"github.com/MrVirul/icecream-inventory/internal/httpx"
	"github.com/MrVirul/icecream-inventory/internal/middleware"
	"github.com/MrVirul/icecream-inventory/internal/service"
	"github.com/gofiber/fiber/v2"
)

type ItemHandler struct {
	items service.ItemService
}

func NewItemHandler(items service.ItemService) *ItemHandler {
	return &ItemHandler{items: items}
}

func parseID(c *fiber.Ctx) (uint, error) {
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id == 0 {
		return 0, httpx.BadRequest("invalid id")
	}
	return uint(id), nil
}

func (h *ItemHandler) List(c *fiber.Ctx) error {
	search := c.Query("search")
	category := c.Query("category")
	lowStock := strings.ToLower(c.Query("low_stock")) == "true"

	allowed := map[string]bool{"name": true, "quantity": true, "price_cents": true, "created_at": true}
	params := httpx.ParseListParams(c, allowed, "name")

	items, meta, err := h.items.List(search, category, lowStock, params.Page, params.Limit, params.Sort, params.Desc)
	if err != nil {
		return err
	}
	return httpx.OKList(c, items, meta)
}

func (h *ItemHandler) Get(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return err
	}
	item, err := h.items.Get(id)
	if err != nil {
		return err
	}
	return httpx.OK(c, item)
}

func (h *ItemHandler) Create(c *fiber.Ctx) error {
	var req dto.ItemCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return httpx.BadRequest("bad request")
	}
	if err := httpx.ValidateStruct(req); err != nil {
		return err
	}
	uid := c.Locals(middleware.CtxUserID).(uint)
	item, err := h.items.Create(req, uid)
	if err != nil {
		return err
	}
	return httpx.Created(c, item)
}

func (h *ItemHandler) Replace(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return err
	}
	var req dto.ItemUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return httpx.BadRequest("bad request")
	}
	if err := httpx.ValidateStruct(req); err != nil {
		return err
	}
	uid := c.Locals(middleware.CtxUserID).(uint)
	item, err := h.items.Replace(id, req, uid)
	if err != nil {
		return err
	}
	return httpx.OK(c, item)
}

func (h *ItemHandler) Patch(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return err
	}
	var req dto.ItemPatchRequest
	if err := c.BodyParser(&req); err != nil {
		return httpx.BadRequest("bad request")
	}
	// check if any field set
	if req.Name == nil && req.Category == nil && req.Description == nil && req.Quantity == nil && req.Unit == nil && req.PriceCents == nil && req.ReorderLevel == nil {
		return httpx.ValidationError([]httpx.FieldError{{Field: "body", Message: "at least one field is required"}})
	}
	if err := httpx.ValidateStruct(req); err != nil {
		return err
	}
	uid := c.Locals(middleware.CtxUserID).(uint)
	item, err := h.items.Patch(id, req, uid)
	if err != nil {
		return err
	}
	return httpx.OK(c, item)
}

func (h *ItemHandler) Delete(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return err
	}
	if err := h.items.Delete(id); err != nil {
		return err
	}
	return httpx.NoContent(c)
}
