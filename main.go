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
	defer conn.Close()

	client := bpb.NewBirthdayFunctionsClient(conn)

	return client
}

func (r *Router) CreateBirthday(c *gin.Context) {
	name := c.Request.FormValue("name")
	personalNumber := c.Request.FormValue("personalNumber")
	date := c.Request.FormValue("date")

	request := &bpb.CreateBirthdayRequest{
		PersonalNumber: personalNumber, Name: name, Date: date}
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
	personalNumber := c.Request.FormValue("personalNumber")
	name := c.Request.FormValue("name")
	date := c.Request.FormValue("date")

	request := &bpb.UpdateBirthdayRequest{PersonalNumber: personalNumber, Name: name, Date: date}
	res, err := r.client.UpdateBirthday(c, request)
	if err != nil {
		fmt.Println("update birthday method failed. error: ", err)
	}
	c.JSON(200, res)
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
		log.Fatalln("failed to run api-gateway router. error: ", err)
	}
}
