package service_products

import (
	"context"
	"fmt"

	"go-web-service/internal/products/domain"
)

// endpoint <-> controller <-> service <-> repository <-> db
// arquivo chama as regras de negocio

type service struct {
	repository domain.Repository
}

func NewService(repository domain.Repository) domain.Service {
	return &service{repository: repository}
}

func (s service) GetAll(ctx context.Context) ([]domain.Produtos, error) {
	produtos, err := s.repository.GetAll(ctx)
	if err != nil {
		return []domain.Produtos{}, err
	}
	return produtos, nil
}

func (s *service) GetOne(id int) (domain.Produtos, error) {
	produto, erro := s.repository.GetOne(id)
	if erro != nil {
		return domain.Produtos{}, erro
	}
	return produto, nil
}

func (s *service) Store(produto domain.ProdutoRequest) (domain.Produtos, error) {
	var prod domain.Produtos = domain.Produtos{
		Name:  produto.Name,
		Type:  produto.Type,
		Count: produto.Count,
		Price: produto.Price,
	}

	newProduto, err := s.repository.Store(prod)
	if err != nil {
		messageErr := fmt.Errorf("error: falha ao registra um novo produto, %w", err)
		return newProduto, messageErr
	}
	return newProduto, nil
}

func (s *service) UpdateName(id int, name string) (string, error) {
	if ok, err := productExists(s, id); !ok {
		return name, err
	}
	produto, err := s.repository.UpdateName(id, name)
	if err != nil {
		return name, err
	}
	return produto, nil
}

func (s *service) Update(id int, produto domain.ProdutoRequest) (domain.Produtos, error) {
	var prod domain.Produtos = domain.Produtos{
		ID:    id,
		Name:  produto.Name,
		Type:  produto.Type,
		Count: produto.Count,
		Price: produto.Price,
	}

	if ok, err := productExists(s, id); !ok {
		return prod, err
	}

	produt, err := s.repository.Update(prod)
	if err != nil {
		return produt, err
	}

	return produt, nil
}

func (s *service) Delete(id int) error {
	if ok, err := productExists(s, id); !ok {
		return err
	}

	if err := s.repository.Delete(id); err != nil {
		return err
	}
	return nil
}

func productExists(s *service, id int) (bool, error) {
	_, err := s.GetOne(id)
	if err != nil {
		return false, err
	}
	return true, nil
}
