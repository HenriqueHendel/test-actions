package main

import (
	"database/sql"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
)

type Product struct {
	Id    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Price float64   `json:"price"`
}

func main() {
	e := echo.New()

	e.GET("/products", listProducts)
	e.POST("/products", createProduct)

	e.Logger.Fatal(e.Start(":3333"))
}

func listProducts(c echo.Context) error {
	db, err := sql.Open("sqlite3", "database.db")

	if err != nil {
		return err
	}

	rows, err := db.Query("Select * from products")

	if err != nil {
		return err
	}

	err = rows.Err()

	if err != nil {
		return err
	}

	defer rows.Close()

	return c.JSON(200, rows)
}

func createProduct(c echo.Context) error {
	uuid := uuid.New()
	product := Product{Id: uuid}
	c.Bind(&product)

	err := persistProduct(product)

	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	return c.JSON(http.StatusCreated, product)
}

func persistProduct(product Product) error {
	db, err := sql.Open("sqlite3", "database.db")

	if err != nil {
		return err
	}

	stmt, err := db.Prepare("Insert into products(id, name, price) values($1, $2, $3)")

	if err != nil {
		return err
	}

	defer stmt.Close()
	_, err = stmt.Exec(product.Id, product.Name, product.Price)

	if err != nil {
		return err
	}

	return nil
}
