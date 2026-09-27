package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Peserta struct {
	ID           int    `json:"id" gorm:"primaryKey"`
	Nama         string `json:"nama"`
	Asal_Sekolah string `json:"asal_sekolah"`
}

var db *gorm.DB

func initDB() {
	dsn := "root:@tcp(127.0.0.1:3306)/db_peserta?charset=utf8mb4&parseTime=True&loc=Local"

	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Gagal konek ke database: " + err.Error())
	}

	db.AutoMigrate(&Peserta{})

	fmt.Println("Berhasil konek ke MySQL")
}

func main() {
	initDB()

	router := gin.Default()

	v1 := router.Group("/peserta")
	{
		v1.GET("/", getPesertas)
		v1.POST("/", createPeserta)
		v1.GET("/:id", getPeserta)
		v1.PUT("/:id", updatePeserta)
		v1.DELETE("/:id", deletePeserta)
	}

	router.Run(":8080")
}

func getPesertas(c *gin.Context) {
	var pesertas []Peserta
	db.Find(&pesertas)
	c.JSON(http.StatusOK, pesertas)
}

func getPeserta(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var peserta Peserta

	if err := db.First(&peserta, id).Error; err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, peserta)
}

func createPeserta(c *gin.Context) {
	var peserta Peserta
	if err := c.ShouldBindJSON(&peserta); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Create(&peserta)
	c.JSON(http.StatusCreated, peserta)
}

func updatePeserta(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var peserta Peserta

	if err := db.First(&peserta, id).Error; err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	var input Peserta
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	peserta.Nama = input.Nama
	peserta.Asal_Sekolah = input.Asal_Sekolah
	db.Save(&peserta)

	c.JSON(http.StatusOK, peserta)
}

func deletePeserta(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var peserta Peserta

	if err := db.First(&peserta, id).Error; err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	db.Delete(&peserta)
	c.Status(http.StatusNoContent)
}