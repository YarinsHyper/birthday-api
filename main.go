package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func CreateBirthday(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "recieved and created a birthday succuesfully",
	})
}

func GetBirthday(c *gin.Context) {
	personalNumber := c.Query("personalNumber")

	c.JSON(200, gin.H{
		"recieved personal number (to read) succesfully. p-n": personalNumber,
	})
}

func UpdateBirthday(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "recieved and updated a birthday succesfully",
	})
}

func DeleteBirthday(c *gin.Context) {
	personalNumber := c.Query("personalNumber")

	c.JSON(200, gin.H{
		"recieved personal number (to delete) succesfully. p-n": personalNumber,
	})
}

func main() {
	fmt.Println("api gateway is running succesfully")

	r := gin.Default()

	r.POST("/api/createBirthday", CreateBirthday)
	r.GET("/api/getBirthday/query", GetBirthday)
	r.POST("/api/updateBirthday", UpdateBirthday)
	r.DELETE("/api/deleteBirthday/query", DeleteBirthday)

	r.Run()
}
