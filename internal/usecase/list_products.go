package usecase

import "github.com/bonfimjustino7/kafka-go-example/internal/entity"

type ListProductsOutDto struct {
	ID    string
	Name  string
	Price float64
}

type ListProductUseCase struct {
	ProductRepository entity.ProductRepository
}

func NewListProductUseCase(productRepository entity.ProductRepository) *ListProductUseCase {
	return &ListProductUseCase{ProductRepository: productRepository}
}

func (u *ListProductUseCase) Execute() ([]*ListProductsOutDto, error) {
	products, err := u.ProductRepository.FindAll()
	if err != nil {
		return nil, err
	}

	var productsOutput []*ListProductsOutDto
	for _, product := range products {
		productsOutput = append(productsOutput, &ListProductsOutDto{
			ID:    product.ID,
			Name:  product.Name,
			Price: product.Price,
		})
	}
	return productsOutput, err
}
