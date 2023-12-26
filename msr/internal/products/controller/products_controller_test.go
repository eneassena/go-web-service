package controller_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	productController "web-service-gin/internal/products/controller"
	"web-service-gin/internal/products/domain"
	"web-service-gin/internal/products/domain/mocks"
	"web-service-gin/pkg/web"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type ProductTest struct {
	Name  string  `json:"name"`
	Type  string  `json:"type"`
	Count int     `json:"count"`
	Price float64 `json:"price"`
}

var pList = []domain.Produtos{
	{
		ID:    1,
		Name:  "Tenis",
		Type:  "Calçados",
		Count: 1,
		Price: 342,
	}, {
		Name:  "Monito",
		Type:  "Informatica",
		Count: 1,
		Price: 100.00,
	},
}

func CreateServer(serv *mocks.Service, method, url, body string) *httptest.ResponseRecorder {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	productController.NewProduto(router, serv)
	req, rr := CreateRequestTest(method, url, body)
	router.ServeHTTP(rr, req)

	return rr
}

func CreateRequestTest(method, url, body string) (*http.Request, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(method, url, bytes.NewBuffer([]byte(body)))
	req.Header.Add("Content-Type", "application/json")
	return req, httptest.NewRecorder()
}

func Test_GetProductsAll(t *testing.T) {
	t.Run("GetAll", func(t *testing.T) {
		mockServ := &mocks.Service{}
		mockServ.On("GetAll", mock.Anything).Return(pList, nil).Once()
		req := CreateServer(mockServ, http.MethodGet, "/api/v1/products", "")
		res := web.Response{}
		if err := json.Unmarshal(req.Body.Bytes(), &res); err != nil {
			assert.Error(t, err)
		}
		assert.Equal(t, 200, res.Code)
		assert.Equal(t, "", res.Error)
		assert.NotNil(t, res.Data)
	})
}

func Test_ProductNotFound(t *testing.T) {
	t.Run("GetAll, [error]", func(t *testing.T) {
		mockServ := &mocks.Service{}
		messageErr := errors.New("not has product registered")
		mockServ.On("GetAll", mock.Anything).Return([]domain.Produtos{}, messageErr).Once()
		req := CreateServer(mockServ, http.MethodGet, "/api/v1/products", "")
		responseRequest := web.Response{}
		if err := json.Unmarshal(req.Body.Bytes(), &responseRequest); err != nil {
			assert.Error(t, err)
		}
		assert.Equal(t, http.StatusNotFound, responseRequest.Code)
		assert.Equal(t, "not has product registered", responseRequest.Error)
		assert.True(t, strings.Contains(responseRequest.Error, "not has product"))
	})
}

func Test_PostProducts(t *testing.T) {
	mockServ := &mocks.Service{}

	t.Run("create products, success 200", func(t *testing.T) {
		mockServ.On("Store", mock.AnythingOfType("domain.ProdutoRequest")).
			Return(pList[1], nil).
			Once()

		dt, _ := json.Marshal(pList[1])

		rr := CreateServer(
			mockServ,
			http.MethodPost,
			"/api/v1/products",
			string(dt),
		)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("create products, error 422", func(t *testing.T) {
		messageErr := errors.New("falha ao cria um novo produto")
		mockServ.On("Store", mock.AnythingOfType("domain.ProdutoRequest")).
			Return(pList[0], messageErr).
			Once()

		product := ProductTest{
			Name:  "Monito",
			Count: 1,
			Price: 100.00,
		}
		dt, _ := json.Marshal(product)

		responseHttptest := CreateServer(
			mockServ,
			http.MethodPost,
			"/api/v1/products",
			string(dt),
		)
		assert.Equal(t, http.StatusUnprocessableEntity, responseHttptest.Code)
	})

	t.Run("create products, error 400", func(t *testing.T) {
		messageErr := errors.New("campo price é invalido")
		mockServ.On("Store", mock.AnythingOfType("domain.ProdutoRequest")).
			Return(pList[0], messageErr).
			Once()
		product := ProductTest{
			Name:  "Monito",
			Type:  "Informatica",
			Count: 1,
			Price: 1000000.00,
		}
		dt, _ := json.Marshal(product)
		responseHttptest := CreateServer(
			mockServ,
			http.MethodPost,
			"/api/v1/products",
			string(dt),
		)
		assert.Equal(t, http.StatusBadRequest, responseHttptest.Code)
	})
}

func Test_PutProducts(t *testing.T) {
	mockServ := &mocks.Service{}

	t.Run("update products, success", func(t *testing.T) {
		currentProduct := pList[1]
		afterProduct := ProductTest{
			Name:  "Mause",
			Type:  "Informatica",
			Count: 20,
			Price: 120.00,
		}

		mockServ.On("domain.ProdutoRequest").
			Return(afterProduct, nil).
			Once()
		body, _ := json.Marshal(currentProduct)

		responseHttptest := CreateServer(mockServ, http.MethodPut, "/api/v1/products/1", string(body))
		var responseObj web.Response
		err := json.Unmarshal(responseHttptest.Body.Bytes(), &responseObj)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, responseObj.Code)
		assert.Equal(t, "", responseObj.Error)
	})

	t.Run("update products, error, not found", func(t *testing.T) {
		currentProduct := pList[1]

		afterProduct := ProductTest{
			Name:  "Monito",
			Type:  "Informatica",
			Count: 10,
			Price: 100.00,
		}
		body, _ := json.Marshal(currentProduct)

		messageErro := errors.New("product is not found")

		mockServ.On("domain.ProdutoRequest").
			Return(afterProduct, messageErro).
			Once()

		responseHttpTest := CreateServer(mockServ, http.MethodPut, "/api/v1/products/1", string(body))

		var responseObj web.Response
		err := json.Unmarshal(responseHttpTest.Body.Bytes(), &responseObj)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, responseObj.Code)
	})
}

func Test_PutProductsName(t *testing.T) {
	t.Run("update Name products, success", func(t *testing.T) {
		productName := struct {
			Name string `json:"name"`
		}{Name: "Fone Game"}

		dt, _ := json.Marshal(productName)

		mockService := new(mocks.Service)
		mockService.On("UpdateName",
			mock.AnythingOfType("int"),
			mock.AnythingOfType("string"),
		).Return(productName.Name, nil).Once()
		responseHttptest := CreateServer(
			mockService,
			http.MethodPatch,
			"/api/v1/products/1",
			string(dt),
		)
		assert.Equal(t, http.StatusOK, responseHttptest.Code)
	})
	t.Run("update Name products, error", func(t *testing.T) {
		productName := struct {
			Name string `json:"name"`
		}{}
		dt, _ := json.Marshal(productName)

		mockService := new(mocks.Service)
		mockService.On("UpdateName",
			mock.AnythingOfType("int"),
			mock.AnythingOfType("string"),
		).Return(productName.Name, nil).Once()
		responseHttptest := CreateServer(
			mockService,
			http.MethodPatch,
			"/api/v1/products/11",
			string(dt),
		)

		assert.Equal(t, http.StatusUnprocessableEntity, responseHttptest.Code)
	})
	t.Run("update Name products, error no paramentro", func(t *testing.T) {
		productName := struct {
			Name string `json:"name"`
		}{Name: "Fone"}

		messageError := errors.New("produto não encontrado")
		dt, _ := json.Marshal(productName)

		urlUpdate := fmt.Sprintf("/api/v1/products/%d", 100)
		mockService := new(mocks.Service)
		mockService.On("UpdateName",
			mock.AnythingOfType("int"),
			mock.AnythingOfType("string"),
		).Return(productName.Name, messageError).Once()
		responseHttptest := CreateServer(
			mockService,
			http.MethodPatch,
			urlUpdate,
			string(dt),
		)
		assert.Equal(t, http.StatusNotFound, responseHttptest.Code)
	})
	t.Run("update Name products, error no paramentro", func(t *testing.T) {
		productName := struct {
			Name string `json:"name"`
		}{Name: "Fone"}

		dt, _ := json.Marshal(productName)

		mockService := new(mocks.Service)
		mockService.On("UpdateName",
			mock.AnythingOfType("int"),
			mock.AnythingOfType("string"),
		).Return(productName.Name, nil).Once()

		urlUpdate := fmt.Sprintf("/api/v1/products/%s", "100s")

		responseHttptest := CreateServer(
			mockService,
			http.MethodPatch,
			urlUpdate,
			string(dt),
		)

		assert.Equal(t, http.StatusBadRequest, responseHttptest.Code)
	})
}

func Test_DeleteProducts(t *testing.T) {
	mockServe := &mocks.Service{}

	t.Run("delete products, success", func(t *testing.T) {
		mockServe.On("Delete", mock.AnythingOfType("int")).Return(nil).Once()
		responseHttptest := CreateServer(mockServe,
			http.MethodDelete,
			fmt.Sprintf("/api/v1/products/%v", 12),
			"",
		)
		assert.Equal(t, 204, responseHttptest.Code)
	})
	t.Run("delete products, error", func(t *testing.T) {
		mockServe.On("Delete", mock.AnythingOfType("int")).Return(errors.New("produto não foi removido")).Once()
		responseHttptest := CreateServer(mockServe,
			http.MethodDelete,
			fmt.Sprintf("/api/v1/products/%v", 9),
			"",
		)
		obj := web.Response{}
		dataResponseHttpTest := json.Unmarshal(responseHttptest.Body.Bytes(), &obj)
		assert.NoError(t, dataResponseHttpTest)
		assert.Equal(t, 400, obj.Code)
	})
	t.Run("delete products, error", func(t *testing.T) {
		mockServe.On("Delete", mock.AnythingOfType("int")).Return(nil).Once()
		responseHttptest := CreateServer(mockServe,
			http.MethodDelete,
			fmt.Sprintf("/api/v1/products/%v", "10s"),
			"",
		)
		assert.Equal(t, 400, responseHttptest.Code)
	})
}
