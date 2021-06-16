package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	bpb "github.com/yarinBenisty/birthday-service/proto"
	"google.golang.org/grpc"
)

type Router struct {
	client bpb.BirthdayFunctionsClient
}

// type BirthdayObject struct {
// 	Name           string
// 	Date           string
// 	PersonalNumber string
// 	ID             string `bson:"_id" json:"id,omitempty"`
// }

func initClientConnection() bpb.BirthdayFunctionsClient {

	conn, err := grpc.Dial(
		"localhost:8000",
		grpc.WithInsecure(),
		grpc.FailOnNonTempDialError(true),
		grpc.WithBlock(),
	)

	if err != nil {
		log.Fatalln("error: ", err)
	}
	// defer conn.Close()

	client := bpb.NewBirthdayFunctionsClient(conn)

	return client
}

func (r *Router) CreateBirthday(c *gin.Context) {
	// name := c.Request.Form.Get("name")
	// personalNumber := c.Request.Form.Get("personalNumber")
	// date := c.Request.Form.Get("date")

	request := &bpb.CreateBirthdayRequest{
		PersonalNumber: bpb.BirthdayObject.PersonalNumber, Name: bpb.BirthdayObject.Name, Date: bpb.BirthdayObject.Date}
	res, err := r.client.CreateBirthday(c, request)
	if err != nil {
		fmt.Println("create birthday method failed. error: ", err)
	}

	c.JSON(200, res)
}

func (r *Router) GetBirthday(c *gin.Context) {
	personalNumber := c.Query("personalNumber")

	request := &bpb.GetBirthdayRequest{PersonalNumber: personalNumber}
	res, err := r.client.GetBirthday(c, request)
	if err != nil {
		fmt.Println("get birthday method failed. error: ", err)
	}

	c.JSON(200, res)
}

func (r *Router) UpdateBirthday(c *gin.Context) {

	c.JSON(200, gin.H{
		"message": "recieved and updated a birthday succesfully",
	})
}

func (r *Router) DeleteBirthday(c *gin.Context) {
	personalNumber := c.Query("personalNumber")

	request := &bpb.DeleteBirthdayRequest{PersonalNumber: personalNumber}
	res, err := r.client.DeleteBirthday(c, request)
	if err != nil {
		fmt.Println("delete birthday method failed. error: ", err)
	}

	c.JSON(200, res)
}

func main() {
	r := &Router{}
	r.client = initClientConnection()

	mainRouter := gin.Default()

	mainRouter.POST("/api/createBirthday", r.CreateBirthday)
	mainRouter.GET("/api/getBirthday/query", r.GetBirthday)
	mainRouter.POST("/api/updateBirthday", r.UpdateBirthday)
	mainRouter.DELETE("/api/deleteBirthday/query", r.DeleteBirthday)

	err := mainRouter.Run(":9000")
	if err != nil {
		log.Fatalln("error: ", err)
	}
}
