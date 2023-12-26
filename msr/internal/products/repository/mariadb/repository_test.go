package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"regexp"
	"testing"

	domain "web-service-gin/internal/products/domain"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	return db, mock
}

var product = domain.Produtos{
	ID:    1,
	Name:  "iphone",
	Type:  "Eletronicos",
	Count: 1,
	Price: 10000,
}

var mockProducts = []domain.Produtos{
	{
		ID:    1,
		Name:  "Playstation 5",
		Type:  "Eletrônicos",
		Count: 1,
		Price: 4500,
	},
	{
		ID:    2,
		Name:  "XBOX Series X",
		Type:  "Eletrônicos",
		Count: 1,
		Price: 4500,
	},
}

func TestGetAll(t *testing.T) {
	t.Run("teste select, (sucesso)", func(t *testing.T) {
		db, mock := NewMock() // cria um mock do bando de dados

		defer db.Close() // fecha conexão

		// cria os retorno que esperamos do banco de dados
		rows := sqlmock.NewRows([]string{
			"id", "name", "type", "count", "price",
		}).AddRow(
			mockProducts[0].ID,
			mockProducts[0].Name,
			mockProducts[0].Type,
			mockProducts[0].Count,
			mockProducts[0].Price,
		).AddRow(
			mockProducts[1].ID,
			mockProducts[1].Name,
			mockProducts[1].Type,
			mockProducts[1].Count,
			mockProducts[1].Price,
		)

		// cria a query que será executada
		query := "SELECT id,name,type,count,price FROM products"
		// executa a query criada
		mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(rows)
		// cria um repository com db mockado
		productsRepo := NewMariaDBRepository(db)
		// obtem o resultado do metodo testado
		ctx := context.Background()
		result, err := productsRepo.GetAll(ctx)
		// verifica o error de retorno do metodo
		assert.NoError(t, err)
		// verifica a lista de produtos retornada pelo metodo testado
		assert.Equal(t, result[0].Name, "Playstation 5")
		assert.Equal(t, result[1].Name, "XBOX Series X")
	})
}

func TestGetAllFailSelect(t *testing.T) {
	t.Run("test listar products (error)", func(t *testing.T) {
		db, mock := NewMock()
		defer db.Close()
		query := "SELECT id,name,type,count,price FROM products"
		mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnError(sql.ErrNoRows)
		productsRepo := NewMariaDBRepository(db)
		ctx := context.Background()
		_, err := productsRepo.GetAll(ctx)
		assert.Error(t, err)
	})
}

func TestStore(t *testing.T) {
	db, mock := NewMock() // cria um mock do bando de dados

	defer db.Close() // fecha conexão
	query := `INSERT INTO products (name, type, count, price) VALUES (?, ?, ?, ?)`
	sqlMock := mock.ExpectPrepare(regexp.QuoteMeta(query))
	sqlMock.ExpectExec().WithArgs(&product.Name, &product.Type, &product.Count, &product.Price).
		WillReturnResult(sqlmock.NewResult(1, 1))
	repoProduct := NewMariaDBRepository(db)
	result, err := repoProduct.Store(product)
	assert.NoError(t, err)
	assert.ObjectsAreEqual(product, result)
}

func TestStoreError(t *testing.T) {
	db, mock := NewMock() // cria um mock do bando de dados

	defer db.Close() // fecha conexão
	expectError := errors.New(`ExecQuery 'INSERT INTO products (name, type, count, price) VALUES (?, ?, ?, ?)', arguments do not match: argument 0 expected [int64 - 20] does not match actual [string - iphone]`)
	query := `INSERT INTO products (name, type, count, price) VALUES (?, ?, ?, ?)`
	sqlMock := mock.ExpectPrepare(regexp.QuoteMeta(query))
	name := 20
	sqlMock.ExpectExec().WithArgs(&name, &product.Type, &product.Count, &product.Price).
		WillReturnError(expectError)
	repoProduct := NewMariaDBRepository(db)
	objProduct, err := repoProduct.Store(product)
	assert.Equal(t, expectError, err)
	assert.Equal(t, product, objProduct)
}

func TestGetID(t *testing.T) {
	db, mock := NewMock() // cria um mock do bando de dados
	defer db.Close()      // fecha conexão
	rows := sqlmock.NewRows([]string{
		"id", "name", "type", "count", "price",
	}).AddRow(
		mockProducts[0].ID,
		mockProducts[0].Name,
		mockProducts[0].Type,
		mockProducts[0].Count,
		mockProducts[0].Price,
	)
	query := `SELECT id,name,type,count,price FROM products WHERE id=?`
	mock.ExpectQuery(query).WithArgs(2).WillReturnRows(rows)
	repoProduct := NewMariaDBRepository(db)
	objProduct, err := repoProduct.GetOne(2)
	assert.NoError(t, err)
	assert.NotEmpty(t, objProduct)
}

func TestUpdate(t *testing.T) {
	db, mock := NewMock() // cria um mock do bando de dados
	defer db.Close()      // fecha conexão

	mockProduct := domain.Produtos{
		ID:    1,
		Name:  "iphone",
		Type:  "Eletronicos",
		Count: 1,
		Price: 10000,
	}

	query := regexp.QuoteMeta(`UPDATE products SET name=?, type=?, count=?, price=? WHERE id=?`)
	stmtM := mock.ExpectPrepare(query)
	stmtM.ExpectExec().WithArgs(
		mockProduct.Name, mockProduct.Type, mockProduct.Count, mockProduct.Price, mockProduct.ID,
	).WillReturnResult(sqlmock.NewResult(0, 1))
	repoProduct := NewMariaDBRepository(db)
	objProduct, err := repoProduct.Update(product)
	assert.NoError(t, err)
	assert.Equal(t, mockProduct.Name, objProduct.Name)
}

func TestUpdateName(t *testing.T) {
	db, mock := NewMock() // cria um mock do bando de dados
	defer db.Close()      // fecha conexão

	mockProduct := domain.Produtos{
		ID:    1,
		Name:  "iphone 10",
		Type:  "Eletronicos",
		Count: 1,
		Price: 10000,
	}
	name := "iphone 10"
	query := regexp.QuoteMeta(`UPDATE products SET name=? WHERE id=?`)
	stmtM := mock.ExpectPrepare(query)
	stmtM.ExpectExec().WithArgs(
		&name, &mockProduct.ID,
	).WillReturnResult(sqlmock.NewResult(0, 1))
	repoProduct := NewMariaDBRepository(db)
	name, err := repoProduct.UpdateName(mockProduct.ID, mockProduct.Name)
	assert.NoError(t, err)
	assert.Equal(t, mockProduct.Name, name)
}

func TestDelete(t *testing.T) {
	db, mock := NewMock() // cria um mock do bando de dados
	defer db.Close()      // fecha conexão
	mockProductId := 2
	query := regexp.QuoteMeta(`DELETE FROM products WHERE id=?`)
	stmtM := mock.ExpectExec(query)
	stmtM.WithArgs(&mockProductId).WillReturnResult(sqlmock.NewResult(0, 1))
	repoProduct := NewMariaDBRepository(db)
	objProductError := repoProduct.Delete(mockProductId)
	assert.NoError(t, objProductError)
}

func TestDeleteError(t *testing.T) {
	db, mock := NewMock() // cria um mock do bando de dados
	defer db.Close()      // fecha conexão
	mockProductId := 2
	expectError := errors.New("produto não foi removido")
	query := regexp.QuoteMeta(`DELETE FROM products WHERE id=?`)
	stmtM := mock.ExpectExec(query)
	stmtM.WithArgs(&mockProductId).WillReturnResult(sqlmock.NewErrorResult(expectError))
	repoProduct := NewMariaDBRepository(db)
	objProductError := repoProduct.Delete(mockProductId)
	assert.Error(t, objProductError)
}

func TestDeleteErrorRows(t *testing.T) {
	db, mock := NewMock() // cria um mock do bando de dados
	defer db.Close()      // fecha conexão
	mockProductId := 2
	expectError := errors.New("produto não foi removido")
	query := regexp.QuoteMeta(`DELETE FROM products WHERE id=?`)
	stmtM := mock.ExpectExec(query)
	stmtM.WithArgs(&mockProductId).WillReturnResult(sqlmock.NewResult(0, 0))
	repoProduct := NewMariaDBRepository(db)
	objProductError := repoProduct.Delete(mockProductId)
	assert.Error(t, objProductError)
	assert.Equal(t, expectError, objProductError)
}
