package service_products

import (
	"context"
	"errors"
	"fmt"
	"testing"

	domain "web-service-gin/internal/products/domain"
	mocks "web-service-gin/internal/products/domain/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var productsList []domain.Produtos = []domain.Produtos{
	{
		ID:    1,
		Name:  "Tenis",
		Type:  "Calçados",
		Count: 1,
		Price: 1000.01,
	},
	{
		ID:    2,
		Name:  "Sandalha",
		Type:  "Calçados",
		Count: 1,
		Price: 1000.01,
	},
}

var pdRequest = domain.ProdutoRequest{
	Name:  productsList[0].Name,
	Type:  productsList[0].Type,
	Count: productsList[0].Count,
	Price: productsList[0].Price,
}

func TestServiceGetAll(t *testing.T) {
	mockRep := new(mocks.Repository)

	t.Run("test de integração service e repository, metodo GetAll, caso de sucesso", func(t *testing.T) {
		mockRep.On("GetAll", mock.Anything).
			Return(productsList, nil).
			Once()

		service := NewService(mockRep)
		ctx := context.Background()
		pList, err := service.GetAll(ctx)

		assert.Nil(t, err)
		assert.Equal(t, productsList, pList)

		mockRep.AssertExpectations(t)
	})
	t.Run("test de integração service e repository, metodo GetAll, caso de error", func(t *testing.T) {
		var (
			products    []domain.Produtos = []domain.Produtos{}
			expectedErr error             = errors.New("não há produtos registrados")
		)

		mockRep.On("GetAll", mock.AnythingOfType("*context.emptyCtx")).
			Return(products, expectedErr).
			Once()

		service := NewService(mockRep)
		ctx := context.Background()
		pList, err := service.GetAll(ctx)

		assert.NotNil(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Equal(t, products, pList)

		mockRep.AssertExpectations(t)
	})
}

func TestServiceGetOne(t *testing.T) {
	mockRep := new(mocks.Repository)

	searchProduct := domain.Produtos{
		ID:    1,
		Name:  "Tenis",
		Type:  "Calçados",
		Count: 1,
		Price: 100.12,
	}

	t.Run("test de integração em repository e service no metodo GetOne, caso de sucesso", func(t *testing.T) {
		mockRep.On("GetOne", mock.AnythingOfType("int")).
			Return(searchProduct, nil).
			Once()

		service := NewService(mockRep)

		productEncontrado, err := service.GetOne(searchProduct.ID)

		assert.Nil(t, err)
		assert.ObjectsAreEqual(searchProduct, productEncontrado)

		mockRep.AssertExpectations(t)
	})
	t.Run("test de integração em repository e service no metodo GetOne, caso de error", func(t *testing.T) {
		expectedError := fmt.Errorf("produto não esta registrado")
		expectedProduct := domain.Produtos{}

		mockRep.On("GetOne", mock.AnythingOfType("int")).
			Return(expectedProduct, expectedError).
			Once()

		service := NewService(mockRep)

		productEncontrado, err := service.GetOne(30)

		assert.NotNil(t, err)
		assert.Equal(t, expectedError, err)
		assert.ObjectsAreEqual(expectedProduct, productEncontrado)

		mockRep.AssertExpectations(t)
	})
}

func TestServiceStore(t *testing.T) {
	mockRep := new(mocks.Repository)

	t.Run("test de integração service e repository, caso de sucesso", func(t *testing.T) {
		newProduct := domain.Produtos{
			ID:    1,
			Name:  "Mause",
			Type:  "Informatica",
			Count: 1,
			Price: 645.445,
		}
		mockRep.On("Store",
			mock.AnythingOfType("domain.Produtos"),
		).
			Return(newProduct, nil).
			Once()

		service := NewService(mockRep)

		productCriado, err := service.Store(pdRequest)

		assert.Nil(t, err)
		assert.ObjectsAreEqual(newProduct, productCriado)
	})
	t.Run("test de integração service e repository caso de error", func(t *testing.T) {
		productVazio := domain.Produtos{}
		expectedErr := fmt.Errorf("falha ao criar um novo produto")

		mockRep.On("Store",
			mock.AnythingOfType("domain.Produtos"),
		).
			Return(productVazio, expectedErr).
			Once()

		service := NewService(mockRep)

		_, err := service.Store(pdRequest)
		assert.NotNil(t, err)
		assert.ObjectsAreEqual(expectedErr, err)
		mockRep.AssertExpectations(t)
	})

	t.Run("test de integração service e repository caso de error no LastID", func(t *testing.T) {
		productVazio := domain.Produtos{}
		expectedErr := fmt.Errorf("falha ao criar um novo produto")

		mockRep.On("Store",
			mock.AnythingOfType("domain.Produtos"),
		).
			Return(productVazio, expectedErr).
			Once()

		service := NewService(mockRep)

		_, err := service.Store(pdRequest)

		assert.NotNil(t, err)
		assert.ObjectsAreEqual(expectedErr, err)
	})
}

func TestServiceUpdate(t *testing.T) {
}

func TestServiceUpdateName(t *testing.T) {
}

func TestServiceDelete(t *testing.T) {
}
