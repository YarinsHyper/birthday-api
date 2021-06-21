package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	bpb "github.com/yarinBenisty/birthday-service/proto"
	"google.golang.org/grpc"
)

//router that's connecting to the client
//functions in the bd-service proto file
type Router struct {
	client bpb.BirthdayFunctionsClient
}

//function that creates a new client and connects
//to the bd-service client through the proto file
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

//function that inserts a birthday object
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

//function that read/get all birthday objects
func (r *Router) GetBirthday(c *gin.Context) {
	personalNumber := c.Query("personalNumber")

	request := &bpb.GetBirthdayRequest{PersonalNumber: personalNumber}
	res, err := r.client.GetBirthday(c, request)
	if err != nil {
		fmt.Println("get birthday method failed. error: ", err)
	}

	c.JSON(200, res)
}

//function that updates a birthday object by personalNumber
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

//function that deletes a certain birthday object by personal number
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
