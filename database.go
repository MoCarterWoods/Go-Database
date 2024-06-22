package main

// -------------------------Create Function-------------------------------
// func createProduct(p *Product) error {
// 	_, err := db.Exec("INSERT INTO products (name, price) VALUES ($1, $2)", p.Name, p.Price)
// 	return err
// }


// -------------------------Get Function--------------------------------
// func getProduct(id int) (*Product, error) {
// 	var p Product
// 	row := db.QueryRow("SELECT id, name, price FROM products WHERE id = $1", id)
// 	err := row.Scan(&p.ID, &p.Name, &p.Price)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return nil, nil
// 		}
// 		return nil, err
// 	}
// 	return &p, nil
// }

// -------------------------Get All Function----------------------------
func getProducts() ([]Product, error) {
	rows, err := db.Query("SELECT id, name, price, supplier_id FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var products []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Supplier); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return products, nil

}

// -------------------------Update Function----------------------------
// func updateProduct(id int, p *Product) error {
// 	_, err := db.Exec("UPDATE products SET name = $1, price = $2, supplier_id = $3 WHERE id = $4", p.Name, p.Price, p.Supplier, id)
// 	return err
// }

// // -------------------------Delete Function----------------------------
// func deleteProduct(id int) error {
// 	_, err := db.Exec("DELETE FROM products WHERE id = $1", id)
// 	return err
// }