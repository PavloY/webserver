package storage

import (
	"database/sql"
	"log"
	_ "github.com/lib/pq"
)

type Storage struct {
	config *Config
	db     *sql.DB
	userRepository *UserRepository
	productRepository *ProductRepository
}

func New(config *Config) *Storage {
	return &Storage{
		config: config,
	}
}

func (storage *Storage) Open() error {
	db, err := sql.Open("postgres", storage.config.DatabaseURI)
	if err !=nil {
		return err
	}
	if err := db.Ping(); err != nil{
		return err
	}
	storage.db = db
	log.Println("Database connection created successfully!")
	return nil
}

func (storage *Storage) Close() {
	storage.db.Close()
}

func (s *Storage) User() *UserRepository{
	if s.userRepository !=nil {
		return s.userRepository
	}
	s.userRepository = &UserRepository{storage: s,}
	return s.userRepository
}

func (s *Storage) Product() *ProductRepository{
	if s.productRepository !=nil {
		return s.productRepository
	}
	s.productRepository = &ProductRepository{storage: s,}
	return s.productRepository
}