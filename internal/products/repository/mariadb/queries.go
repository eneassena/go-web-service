package repository

const (
	SqlGetAll     = `SELECT id,name,type,count,price FROM products`
	SqlGetOne     = `SELECT id,name,type,count,price FROM products WHERE id=?`
	SqlStore      = `INSERT INTO products (name, type, count, price) VALUES (?, ?, ?, ?)`
	SqlLastID     = `SELECT MAX(id) as last_id FROM products`
	SqlUpdate     = `UPDATE products SET name=?, type=?, count=?, price=? WHERE id=?`
	SqlUpdateName = `UPDATE products SET name=? WHERE id=?`
	SqlDelete     = `DELETE FROM products WHERE id=?`
)
