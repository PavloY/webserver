package api

import (
	"WEB_SERVER/internal/app/middleware"
	"WEB_SERVER/internal/app/models"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"time"

	"github.com/form3tech-oss/jwt-go"
)

type Message struct {
	StatusCode int `json:"status_code"`
	Message string `json:"message"`
	IsError bool `json:"is_error"`
}

func initHeaders(writer http.ResponseWriter){
	writer.Header().Set("Content-Type", "application/json")
}

func (api *API) GetAllProducts(writer http.ResponseWriter, req *http.Request) {
	initHeaders(writer)
	api.logger.Info("GET /api/v1/products")
	products, err := api.storage.Product().SelectAll()
	if err != nil {
		api.logger.Info("Error while Product().SelectAll():", err)
		msg := Message {
			StatusCode: 501,
			Message: "We have some troubles to accessing database.",
			IsError: true,
		}
		writer.WriteHeader(501)
		json.NewEncoder(writer).Encode(msg)
		return
	}
	writer.WriteHeader(200)
	json.NewEncoder(writer).Encode(products)
}

func (api *API) GetProductById(writer http.ResponseWriter, req *http.Request) {
	initHeaders(writer)
	api.logger.Info("Get product by ID /api/v1/products/{id}")
	id, err := strconv.Atoi(mux.Vars(req)["id"])
	if err != nil {
	  api.logger.Info("Troubles while parsing {id} param:", err)
	  msg := Message{
		StatusCode: 400,
		Message:    "Unapropriate id value. don't use ID as uncasting to int value.",
		IsError:    true,
	  }
	  writer.WriteHeader(400)
	  json.NewEncoder(writer).Encode(msg)
	  return
	}
	product, ok, err := api.storage.Product().FindProductById(id)
	if err != nil {
	  api.logger.Info("Troubles while accessing database table (products) with id. err:", err)
	  msg := Message{
		StatusCode: 500,
		Message:    "Have some troubles to accessing database. Try again",
		IsError:    true,
	  }
	  writer.WriteHeader(500)
	  json.NewEncoder(writer).Encode(msg)
	  return
	}
	if !ok {
	  api.logger.Info("Can not find product with that ID in database")
	  msg := Message{
		StatusCode: 404,
		Message:    "Product with that ID does not exists in database.",
		IsError:    true,
	  }
  
	  writer.WriteHeader(404)
	  json.NewEncoder(writer).Encode(msg)
	  return
	}
	writer.WriteHeader(200)
	json.NewEncoder(writer).Encode(product)
  
  }

  func (api *API) UpdateProductById(writer http.ResponseWriter, req *http.Request) {
	initHeaders(writer)
	api.logger.Info("Update product by Id PUT /api/v1/products/{id}")
	id, err := strconv.Atoi(mux.Vars(req)["id"])
	if err != nil {
		api.logger.Info("Troubles while parsing {id} param:", err)
		msg := Message{
			StatusCode: 400,
			Message:    "Inappropriate id value. Don't use ID that cannot be cast to int.",
			IsError:    true,
		}
		writer.WriteHeader(400)
		json.NewEncoder(writer).Encode(msg)
		return
	}

	product, ok, err := api.storage.Product().FindProductById(id)
	if err != nil {
		api.logger.Info("Troubles while accessing database table (products) with id. err:", err)
		msg := Message{
			StatusCode: 500,
			Message:    "We have some troubles accessing the database. Try again.",
			IsError:    true,
		}
		writer.WriteHeader(500)
		json.NewEncoder(writer).Encode(msg)
		return
	}

	if !ok {
		api.logger.Info("Cannot find product with that ID in database")
		msg := Message{
			StatusCode: 404,
			Message:    "Product with that ID does not exist in the database.",
			IsError:    true,
		}
		writer.WriteHeader(404)
		json.NewEncoder(writer).Encode(msg)
		return
	}

	var updatedProduct models.Product
	err = json.NewDecoder(req.Body).Decode(&updatedProduct)
	if err != nil {
		api.logger.Info("Invalid JSON format for updating product. err:", err)
		msg := Message{
			StatusCode: 400,
			Message:    "Invalid JSON format for updating product.",
			IsError:    true,
		}
		writer.WriteHeader(400)
		json.NewEncoder(writer).Encode(msg)
		return
	}

	updatedProduct.ID = product.ID

	_,err = api.storage.Product().UpdateProductById(id, &updatedProduct)
	if err != nil {
		api.logger.Info("Troubles while updating product in the database. err:", err)
		msg := Message{
			StatusCode: 500,
			Message:    "We have some troubles updating the product in the database. Try again.",
			IsError:    true,
		}
		writer.WriteHeader(500)
		json.NewEncoder(writer).Encode(msg)
		return
	}

	writer.WriteHeader(200)
	msg := Message{
		StatusCode: 200,
		Message:    fmt.Sprintf("Product with ID %d successfully updated.", id),
		IsError:    false,
	}
	json.NewEncoder(writer).Encode(msg)
}
  
func (api *API) DeleteProductById(writer http.ResponseWriter, req *http.Request) {
	initHeaders(writer)
  api.logger.Info("Delete product by Id DELETE /api/v1/products/{id}")
  id, err := strconv.Atoi(mux.Vars(req)["id"])
  if err != nil {
    api.logger.Info("Troubles while parsing {id} param:", err)
    msg := Message{
      StatusCode: 400,
      Message:    "Unapropriate id value. don't use ID as uncasting to int value.",
      IsError:    true,
    }
    writer.WriteHeader(400)
    json.NewEncoder(writer).Encode(msg)
    return
  }

  _, ok, err := api.storage.Product().FindProductById(id)
  if err != nil {
    api.logger.Info("Troubles while accessing database table (products) with id. err:", err)
    msg := Message{
      StatusCode: 500,
      Message:    "We have some troubles to accessing database. Try again",
      IsError:    true,
    }
    writer.WriteHeader(500)
    json.NewEncoder(writer).Encode(msg)
    return
  }

  if !ok {
    api.logger.Info("Can not find product with that ID in database")
    msg := Message{
      StatusCode: 404,
      Message:    "Product with that ID does not exists in database.",
      IsError:    true,
    }

    writer.WriteHeader(404)
    json.NewEncoder(writer).Encode(msg)
    return
}
_, err = api.storage.Product().DeleteById(id)
  if err != nil {
    api.logger.Info("Troubles while deleting database elemnt from table (products) with id. err:", err)
    msg := Message{
      StatusCode: 501,
      Message:    "We have some troubles to accessing database. Try again",
      IsError:    true,
    }
    writer.WriteHeader(501)
    json.NewEncoder(writer).Encode(msg)
    return
  }
  writer.WriteHeader(202)
  msg := Message{
    StatusCode: 202,
    Message:    fmt.Sprintf("Product with ID %d successfully deleted.", id),
    IsError:    false,
  }
  json.NewEncoder(writer).Encode(msg)
}

func (api *API) CreateProduct(writer http.ResponseWriter, req *http.Request) {
	initHeaders(writer)
	api.logger.Info("POST /api/v1/products")
	var product models.Product
	err := json.NewDecoder(req.Body).Decode(&product)
	if err !=nil{
		api.logger.Info("invalid json")
		msg := Message{
			StatusCode: 400,
			Message: "json is invalid",
			IsError: true,
		}
		writer.WriteHeader(400)
		json.NewEncoder(writer).Encode(msg)
		return
	}
	p, err := api.storage.Product().Create(&product)
	if err !=nil {
		api.logger.Info("Error creating new product:", err)
		msg := Message{
			StatusCode: 501,
			Message :  "Tronbles to accessing database",
			IsError: true,
		}
		writer.WriteHeader(501)
		json.NewEncoder(writer).Encode(msg)
	}
	writer.WriteHeader(201)
	json.NewEncoder(writer).Encode(p)
	}

func (api *API) RegisterUser(writer http.ResponseWriter, req *http.Request) {
	initHeaders(writer)
	api.logger.Info("Post User Register POST /api/v1/user/register")
	var user models.User
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
	  api.logger.Info("Invalid json recieved from client")
	  msg := Message{
		StatusCode: 400,
		Message:    "Provided json is invalid",
		IsError:    true,
	  }
	  writer.WriteHeader(400)
	  json.NewEncoder(writer).Encode(msg)
	  return
	}
  
	_, ok, err := api.storage.User().FindByLogin(user.Login)
	if err != nil {
	  api.logger.Info("Troubles while accessing database table (users) with id. err:", err)
	  msg := Message{
		StatusCode: 500,
		Message:    "We have some troubles to accessing database. Try again",
		IsError:    true,
	  }
	  writer.WriteHeader(500)
	  json.NewEncoder(writer).Encode(msg)
	  return
	}
  
	if ok {
	  api.logger.Info("User with that ID already exists")
	  msg := Message{
		StatusCode: 400,
		Message:    "User with that login already exists in database",
		IsError:    true,
	  }
	  writer.WriteHeader(400)
	  json.NewEncoder(writer).Encode(msg)
	  return
	}

	userAdded, err := api.storage.User().Create(&user)
	if err != nil {
	  api.logger.Info("Troubles while accessing database table (users) with id. err:", err)
	  msg := Message{
		StatusCode: 500,
		Message:    "We have some troubles to accessing database. Try again",
		IsError:    true,
	  }
	  writer.WriteHeader(500)
	  json.NewEncoder(writer).Encode(msg)
	  return
	}
  
	msg := Message{
	  StatusCode: 201,
	  Message:    fmt.Sprintf("User {login:%s} successfully registered!", userAdded.Login),
	  IsError:    false,
	}
	writer.WriteHeader(201)
	json.NewEncoder(writer).Encode(msg)
  
  }

  func (api *API) PostToAuth(writer http.ResponseWriter, req *http.Request) {
	initHeaders(writer)
	api.logger.Info("Post to Auth POST /api/v1/user/auth")
	var userFromJson models.User
	err := json.NewDecoder(req.Body).Decode(&userFromJson)

	if err != nil {
	  api.logger.Info("Invalid json recieved from client")
	  msg := Message{
		StatusCode: 400,
		Message:    "Provided json is invalid",
		IsError:    true,
	  }
	  writer.WriteHeader(400)
	  json.NewEncoder(writer).Encode(msg)
	  return
	}

	userInDB, ok, err := api.storage.User().FindByLogin(userFromJson.Login)

	if err != nil {
	  api.logger.Info("Can not make user search in database:", err)
	  msg := Message{
		StatusCode: 500,
		Message:    "We have some troubles while accessing database",
		IsError:    true,
	  }
	  writer.WriteHeader(500)
	  json.NewEncoder(writer).Encode(msg)
	  return
	}
  
	if !ok {
	  api.logger.Info("User with that login does not exists")
	  msg := Message{
		StatusCode: 400,
		Message:    "User with that login does not exists in database. Try register first",
		IsError:    true,
	  }
	  writer.WriteHeader(400)
	  json.NewEncoder(writer).Encode(msg)
	  return
	}

	token := jwt.New(jwt.SigningMethodHS256)             
	claims := token.Claims.(jwt.MapClaims)               
	claims["exp"] = time.Now().Add(time.Hour * 2).Unix() 
	claims["admin"] = true
	claims["name"] = userInDB.Login
	tokenString, err := token.SignedString(middleware.SecretKey)

	if err != nil {
	  api.logger.Info("Can not claim jwt-token")
	  msg := Message{
		StatusCode: 500,
		Message:    "We have some troubles. Try again",
		IsError:    true,
	  }
	  writer.WriteHeader(500)
	  json.NewEncoder(writer).Encode(msg)
	  return
	}

	msg := Message{
	  StatusCode: 201,
	  Message:    tokenString,
	  IsError:    false,
	}
	writer.WriteHeader(201)
	json.NewEncoder(writer).Encode(msg)
  
  }
