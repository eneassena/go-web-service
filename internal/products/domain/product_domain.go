package domain

import "context"

type Produtos struct {
	ID    int     `json:"id" `
	Name  string  `json:"name"`
	Type  string  `json:"type"`
	Count int     `json:"count"`
	Price float64 `json:"price"`
}

type ProdutoRequest struct {
	Name  string  `json:"name" binding:"required"`
	Type  string  `json:"type" binding:"required"`
	Count int     `json:"count" binding:"required"`
	Price float64 `json:"price" binding:"required"`
}
type Repository interface {
	GetAll(ctx context.Context) ([]Produtos, error)
	GetOne(id int) (Produtos, error)
	Store(produto Produtos) (Produtos, error)
	Update(produto Produtos) (Produtos, error)
	UpdateName(id int, name string) (string, error)
	Delete(id int) error
}

type Service interface {
	GetAll(ctx context.Context) ([]Produtos, error)
	GetOne(id int) (Produtos, error)
	Store(produto ProdutoRequest) (Produtos, error)
	Update(id int, produto ProdutoRequest) (Produtos, error)
	UpdateName(id int, name string) (string, error)
	Delete(id int) error
}
