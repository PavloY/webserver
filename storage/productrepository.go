package storage

import (
	"WEB_SERVER/internal/app/models"
	"fmt"
	"log"
)

type ProductRepository struct {
	storage *Storage
}

var (
	tableProduct string = "products"
)

func (pr *ProductRepository) Create(p *models.Product)(*models.Product, error) {
	query :=fmt.Sprintf("INSERT INTO %s (name) VALUES($1) RETURNING id", tableProduct)
	if err := pr.storage.db.QueryRow(query, p.Name).Scan(&p.ID); err != nil{
	return nil, err
}
return p, nil
}

func (pr *ProductRepository) UpdateProductById(id int, updatedProduct *models.Product) (*models.Product, error) {

	product, ok, err := pr.FindProductById(id)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("product with id %d not found", id)
	}

	query := fmt.Sprintf("UPDATE %s SET name=$1 WHERE id=$2", tableProduct)
	_, err = pr.storage.db.Exec(query, updatedProduct.Name, id)
	if err != nil {
		return nil, err
	}

	product.Name = updatedProduct.Name

	return product, nil
}


func (pr *ProductRepository) DeleteById(id int)(*models.Product, error) {
	product, ok, err := pr.FindProductById(id)
	if err != nil {
		return nil, err
	}
	if ok {
		query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", tableProduct)
		_, err := pr.storage.db.Exec(query, id)
		if err !=nil {
			return nil, err
		}
	}
	return product, nil
}

func (pr *ProductRepository) FindProductById(id int)(*models.Product, bool, error) {
	products, err := pr.SelectAll()
	var founded bool
	if err != nil {
		return nil, founded, err
	}
	var productFinded *models.Product
	for _, p := range products{
		if p.ID == id {
			productFinded = p
			founded = true
			break
		}
	}
	return productFinded, founded, nil
}

func (pr *ProductRepository) SelectAll()([]*models.Product, error) {
	query :=fmt.Sprintf("SELECT * FROM %s", tableProduct)
	rows, err := pr.storage.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products :=make([]*models.Product, 0)
	for rows.Next(){
		p:= models.Product{}
		err:=rows.Scan(&p.ID, &p.Name)
		if err !=nil{
			log.Println(err)
			continue
		}
		products = append(products, &p)
	}
	return products, nil
}