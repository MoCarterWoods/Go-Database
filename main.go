package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	_ "github.com/lib/pq"
	"github.com/gofiber/fiber/v2"
)

const (
    host     = "localhost"
    port     = 5432
    user     = "myuser"
    password = "mypassword"
    dbname   = "mydatabase"
)

var db *sql.DB

type Product struct {
	ID    int `json:"id"`
	Name  string `json:"name"`
	Price int `json:"price"`
	Supplier int `json:"supplier"`
}


func main() {
	// สร้าง connection string สำหรับ PostgreSQL
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// เชื่อมต่อฐานข้อมูล PostgreSQL
	sdb, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	db = sdb
	defer db.Close()

	// ทดสอบการเชื่อมต่อฐานข้อมูล
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}

	fmt.Println("Successfully connected to PostgreSQL!")


	//-----------------------Update Function----------------------------------
	// err = updateProduct(1, &Product{Name: "Carter", Price: 222, Supplier: 1})
	// if err != nil {
	// 	log.Fatalf("Error updating product: %v", err)
	// }
	// 	fmt.Println("Product updated successfully!")


	//-----------------------Delete Function----------------------------------
	// err = deleteProduct(5)
	// if err != nil {
	// 	log.Fatalf("Error deleting product: %v", err)
	// }
	// 	fmt.Println("Product deleted successfully!")
	


	

	//-----------------------Get Function----------------------------------
	// product, err := getProduct(1)
	// if err != nil {
	// 	log.Fatalf("Error getting product: %v", err)
	// }
	// if product == nil {
	// 	fmt.Println("Product not found")
	// } else {
	// 	fmt.Printf("Product: %v\n", product)
	// }


	//-----------------------Get All Function----------------------------
	products, err := getProducts()
	if err != nil {
		log.Fatalf("Error getting products: %v", err)
	}

	if len(products) == 0 {
		fmt.Println("No products found")
	} else {
		for _, p := range products {
			fmt.Printf("ID: %d, Name: %s, Price: %d, Supplier: %d\n", p.ID, p.Name, p.Price, p.Supplier)
		}
	}

	//-----------------------Create Function-------------------------------
	// err = createProduct(&Product{Name: "Product 1", Price: 100})
	// if err != nil {
	// 	log.Fatalf("Error creating product: %v", err)
	// }

	// fmt.Println("Product created successfully!")
	

	app := fiber.New()


	app.Get("/product/:id",getProductHandler)
	app.Post("/product",createProductHandler)
	app.Put("/product/:id",updateProductHandler)
	app.Delete("/product/:id",deleteProductHandler)
	app.Get("/products",getProductsHandler)

	app. Listen(":8080")

}

// -------------------------Fiber Get Function----------------------------------
func getProductHandler(c *fiber.Ctx) error {
	productId, err := strconv.Atoi(c.Params("id"))

	product, err := getProduct(productId)
	if err != nil {
		return err
	}
	if product == nil {	
		return c.Status(404).SendString("Product not found")
	}	

	return c.JSON(product)
}


func getProduct(id int) (*Product, error) {
	var p Product
	row := db.QueryRow("SELECT id, name, price, supplier_id FROM products WHERE id = $1", id)
	err := row.Scan(&p.ID, &p.Name, &p.Price, &p.Supplier)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

// -------------------------Fiber Create Function----------------------------

func createProductHandler(c *fiber.Ctx) error {
	var p Product
	if err := c.BodyParser(&p); err != nil {
		return c.Status(500).SendString(err.Error())
	}
	_, err := db.Exec("INSERT INTO products (name, price, supplier_id) VALUES ($1, $2, $3)", p.Name, p.Price, p.Supplier)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	return c.JSON(p)
}

// -------------------------Fiber Update Function----------------------------
func updateProductHandler(c *fiber.Ctx) error {
	// อ่านพารามิเตอร์ id จาก URL และแปลงเป็น integer
	productId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	// แปลงข้อมูล JSON จาก request body มาเป็น struct Product
	var p Product
	if err := c.BodyParser(&p); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	// ทำการอัปเดตข้อมูลในฐานข้อมูล
	_, err = db.Exec("UPDATE products SET name = $1, price = $2, supplier_id = $3 WHERE id = $4", p.Name, p.Price, p.Supplier, productId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	// ส่งข้อมูลที่อัปเดตแล้วกลับเป็น JSON
	return c.JSON(p)
}


// -------------------------Fiber Delete Function----------------------------
func deleteProductHandler(c *fiber.Ctx) error {
	// อ่านพารามิเตอร์ id จาก URL และแปลงเป็น integer
	productId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	// ทำการลบข้อมูลในฐานข้อมูล
	_, err = db.Exec("DELETE FROM products WHERE id = $1", productId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	// ส่งข้อมูลที่ลบแล้วกลับเป็น JSON
	return c.JSON(fiber.Map{"status": "success"})
}

// -------------------------Fiber Gets Function----------------------------------
func getProductsHandler(c *fiber.Ctx) error {
	products, err := getProducts()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.JSON(products)
}