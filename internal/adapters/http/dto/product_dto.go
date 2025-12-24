package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/jadiazinf/inventory/internal/core/domain"
)

// ProductRequest represents the request to create/update a product
type ProductRequest struct {
	SKU          string               `json:"sku" validate:"required"`
	Barcode      *string              `json:"barcode,omitempty"`
	Name         string               `json:"name" validate:"required"`
	Description  *string              `json:"description,omitempty"`
	CategoryID   *uuid.UUID           `json:"category_id,omitempty"`
	SellingPrice float64              `json:"selling_price" validate:"required,gt=0"`
	CostPrice    *float64             `json:"cost_price,omitempty"`
	MinStock int            `json:"min_stock_level,omitempty"`
	MaxStock int            `json:"max_stock_level,omitempty"`
	Status       domain.ProductStatus `json:"status,omitempty"`
	ImageURL     *string              `json:"image_url,omitempty"`
}

// ProductResponse represents a product in API responses
type ProductResponse struct {
	ProductID     uuid.UUID            `json:"product_id"`
	SKU           string               `json:"sku"`
	Barcode       *string              `json:"barcode,omitempty"`
	Name          string               `json:"name"`
	Description   *string              `json:"description,omitempty"`
	CategoryID    *uuid.UUID           `json:"category_id,omitempty"`
	SellingPrice  float64              `json:"selling_price"`
	CostPrice     *float64             `json:"cost_price,omitempty"`
	MinStock int             `json:"min_stock_level,omitempty"`
	MaxStock int             `json:"max_stock_level,omitempty"`
	Status        domain.ProductStatus `json:"status"`
	ImageURL      *string              `json:"image_url,omitempty"`
	CreatedAt     time.Time            `json:"created_at"`
	UpdatedAt     time.Time            `json:"updated_at"`
}

// ProductListResponse represents paginated product list
type ProductListResponse struct {
	Products []ProductResponse `json:"products"`
	Total    int64             `json:"total"`
	Limit    int               `json:"limit"`
	Offset   int               `json:"offset"`
}

// UpdatePriceRequest represents a price update request
type UpdatePriceRequest struct {
	NewPrice float64             `json:"new_price" validate:"required,gt=0"`
	Currency domain.CurrencyCode `json:"currency" validate:"required"`
	Reason   string              `json:"reason,omitempty"`
}

// ToProductDomain converts ProductRequest to domain.Product
func (r *ProductRequest) ToProductDomain() *domain.Product {
	return &domain.Product{
		ProductID:     uuid.New(),
		SKU:           r.SKU,
		Barcode:       r.Barcode,
		Name:          r.Name,
		Description:   r.Description,
		CategoryID:    r.CategoryID,
		SellingPrice:  r.SellingPrice,
		CostPrice:     r.CostPrice,
		MinStock: r.MinStock,
		MaxStock: r.MaxStock,
		Status:        r.Status,
		ImageURL:      r.ImageURL,
	}
}

// ToProductResponse converts domain.Product to ProductResponse
func ToProductResponse(p *domain.Product) ProductResponse {
	return ProductResponse{
		ProductID:     p.ProductID,
		SKU:           p.SKU,
		Barcode:       p.Barcode,
		Name:          p.Name,
		Description:   p.Description,
		CategoryID:    p.CategoryID,
		SellingPrice:  p.SellingPrice,
		CostPrice:     p.CostPrice,
		MinStock: p.MinStock,
		MaxStock: p.MaxStock,
		Status:        p.Status,
		ImageURL:      p.ImageURL,
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
	}
}

// ToProductListResponse converts product slice to list response
func ToProductListResponse(products []domain.Product, total int64, limit, offset int) ProductListResponse {
	responses := make([]ProductResponse, len(products))
	for i, p := range products {
		responses[i] = ToProductResponse(&p)
	}
	return ProductListResponse{
		Products: responses,
		Total:    total,
		Limit:    limit,
		Offset:   offset,
	}
}
