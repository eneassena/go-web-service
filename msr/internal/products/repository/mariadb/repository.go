package repository

import (
	"context"
	"database/sql"
	"errors"

	domain "go-web-service/internal/products/domain"
)

type mariaDBRepository struct {
	db *sql.DB
}

func NewMariaDBRepository(db *sql.DB) domain.Repository {
	return &mariaDBRepository{db: db}
}

func (rep mariaDBRepository) GetAll(ctx context.Context) ([]domain.Produtos, error) {
	var produtosList []domain.Produtos = []domain.Produtos{}
	rows, err := rep.db.QueryContext(ctx, SqlGetAll)
	if err != nil {
		return produtosList, err
	}

	defer rows.Close()
	for rows.Next() {
		var product domain.Produtos
		err := rows.Scan(&product.ID, &product.Name, &product.Type, &product.Count, &product.Price)
		if err != nil {
			return produtosList, err
		}
		produtosList = append(produtosList, product)
	}

	return produtosList, nil
}

func (rep mariaDBRepository) GetOne(id int) (domain.Produtos, error) {
	var produto domain.Produtos
	stmt := rep.db.QueryRow(SqlGetOne, id)

	err := stmt.Scan(
		&produto.ID,
		&produto.Name,
		&produto.Type,
		&produto.Count,
		&produto.Price,
	)
	if err != nil {
		return domain.Produtos{}, errors.New("produto não está registrado")
	}
	return produto, nil
}

func (rep *mariaDBRepository) Store(produto domain.Produtos) (domain.Produtos, error) {
	stmt, err := rep.db.Prepare(SqlStore)
	if err != nil {
		return produto, err
	}

	defer stmt.Close()

	res, err := stmt.Exec(
		&produto.Name,
		&produto.Type,
		&produto.Count,
		&produto.Price,
	)
	if err != nil {
		return produto, err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return produto, err
	}
	produto.ID = int(lastID)

	return produto, nil
}

func (rep *mariaDBRepository) Update(produto domain.Produtos) (domain.Produtos, error) {
	stmt, err := rep.db.Prepare(SqlUpdate)
	if err != nil {
		return produto, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(
		&produto.Name,
		&produto.Type,
		&produto.Count,
		&produto.Price,
		produto.ID,
	)
	if err != nil {
		return produto, err
	}

	return produto, nil
}

func (rep *mariaDBRepository) UpdateName(id int, name string) (string, error) {
	stmt, err := rep.db.Prepare(SqlUpdateName)
	if err != nil {
		return name, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(&name, &id)
	if err != nil {
		return name, err
	}

	return name, nil
}

func (rep *mariaDBRepository) Delete(id int) error {
	result, err := rep.db.Exec(SqlDelete, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("produto não foi removido")
	}
	return nil
}
