package databases

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Code  string `gorm:"index"`
	Price int
}

func ternary[T any](cond bool, ifTrue T, ifNotTrue T) T {
	if cond {
		return ifTrue
	} else {
		return ifNotTrue
	}
}

func (p Product) PrettyPrint() {
	fmt.Printf("ID: %d\nCode: %s\nPrice: %d\nCreatedAt: %s\nUpdatedAt: %s\nDeletedAt: %s\n",
		p.ID,
		p.Code,
		p.Price,
		p.CreatedAt.Format(time.RFC3339),
		p.UpdatedAt.Format(time.RFC3339),
		ternary(p.DeletedAt.Valid, p.DeletedAt.Time.Format(time.RFC3339), "<nil>"),
	)
}

func LearnGorm() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Panicln("Unable to connect to the database!", err)
	}

	fmt.Println("migrating models...")
	db.AutoMigrate(&Product{})

	fmt.Println("inserting product into table...")
	db.Create(&Product{Code: "p123", Price: 1000})

	fmt.Println("getting the product...")
	var product Product
	db.First(&product, "code = ?", "p123")
	product.PrettyPrint()

	fmt.Println("updating the product...")
	db.Model(&product).Update("price", 500)
	product.PrettyPrint()

	fmt.Println("deleting the product...")
	db.Delete(&product, 1)
}
