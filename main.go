package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	bpb "github.com/yarinBenisty/birthday-service/proto"
	"google.golang.org/grpc"
)

const (
	name           = "name"
	personalNumber = "personalNumber"
	date           = "date"
)

// Router that's connecting to the client
// functions in the bd-service proto file
type Router struct {
	client bpb.BirthdayFunctionsClient
	bpb.UnimplementedBirthdayFunctionsServer
}

// Func that creates a new client and connects
// to the bd-service client through the proto file
func initClientConnection() bpb.BirthdayFunctionsClient {

	envError := godotenv.Load()
	if envError != nil {
		log.Fatal("Error loading .env file! error: ", envError)
	}

	address := os.Getenv("ADDRESS") + os.Getenv("PORT")
	conn, err := grpc.Dial(
		address,
		grpc.WithInsecure(),
		grpc.FailOnNonTempDialError(true),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatalln("connection error: ", err)
		os.Exit(4)
	}
	// defer conn.Close()

	client := bpb.NewBirthdayFunctionsClient(conn)

	return client
}

func corsRouterConfig() cors.Config {
	corsConfig := cors.DefaultConfig()
	corsConfig.AddExposeHeaders("x-uploadid")
	corsConfig.AllowAllOrigins = false
	corsConfig.AllowWildcard = true
	corsConfig.AllowOrigins = strings.Split("http://localhost*", ",")
	corsConfig.AllowCredentials = true
	corsConfig.AddAllowHeaders(
		"x-content-length",
		"authorization",
		"cache-control",
		"x-requested-with",
		"content-disposition",
		"content-range",
		"destination",
		"fileID",
	)

	return corsConfig
}

//CreateBirthday inserts a birthday object
func (r *Router) CreateBirthday(c *gin.Context) {
	var code = 201
	request := &bpb.CreateBirthdayRequest{
		PersonalNumber: c.Request.FormValue(personalNumber),
		Name:           c.Request.FormValue(name),
		Date:           c.Request.FormValue(date),
	}
	res, err := r.client.CreateBirthday(c, request)
	if err != nil {
		code = 400
		fmt.Println("create birthday error: ", err)
		c.String(code, "create birthday method failed. \nerror: %s", err)
	}
	c.JSON(code, res)
}

// GetBirthday get a birthday object
func (r *Router) GetBirthday(c *gin.Context) {
	var code = 200
	request := &bpb.GetBirthdayRequest{PersonalNumber: c.Query(personalNumber)}
	res, err := r.client.GetBirthday(c, request)
	if err != nil {
		code = 404
		fmt.Println("get birthday error: ", err)
		c.String(code, "get birthday method failed. \nerror: %s", err)
	}
	c.JSON(code, res)
}

//GetAllBirthday gets all birthday objects
func (r *Router) GetAllBirthday(c *gin.Context) {
	var code = 200
	request := &bpb.GetAllBirthdayRequest{}
	res, err := r.client.GetAllBirthday(c, request)
	if err != nil {
		code = 404
		fmt.Println("get all birthday error.\n", err)
		c.String(code, "get all birthday method failed. \nerror: %s", err)
	}
	c.JSON(code, res)
}

// UpdateBirthday updates a birthday object by personalNumber
func (r *Router) UpdateBirthday(c *gin.Context) {
	var code = 201
	request := &bpb.UpdateBirthdayRequest{PersonalNumber: personalNumber, Name: name, Date: date}
	res, err := r.client.UpdateBirthday(c, request)
	if err != nil {
		code = 400
		fmt.Println("update birthday error: ", err)
		c.String(code, "update birthday method failed. \nerror: %s", err)
	}
	c.JSON(code, res)
}

// DeleteBirthday deletes a certain birthday object by personal number
func (r *Router) DeleteBirthday(c *gin.Context) {
	var code = 204
	request := &bpb.DeleteBirthdayRequest{PersonalNumber: personalNumber}
	res, err := r.client.DeleteBirthday(c, request)
	if err != nil {
		code = 404
		fmt.Println("delete birthday error: ", err)
		c.String(code, "delete birthday method failed. \nerror: %s", err)
	}
	c.JSON(code, res)
}

func main() {

	envError := godotenv.Load()
	if envError != nil {
		log.Fatal("Error loading .env file! error: ", envError)
	}
	routerPort := os.Getenv("ROUTER_PORT")

	r := &Router{}
	r.client = initClientConnection()

	mainRouter := gin.Default()
	mainRouter.POST("/api/createBirthday", r.CreateBirthday)
	mainRouter.GET("/api/getBirthday", r.GetBirthday)
	mainRouter.GET("/api/getAllBirthday", r.GetAllBirthday)
	mainRouter.POST("/api/updateBirthday", r.UpdateBirthday)
	mainRouter.DELETE("/api/deleteBirthday", r.DeleteBirthday)

	mainRouter.Use(
		cors.New(corsRouterConfig()),
	)
	err := mainRouter.Run(":" + routerPort)
	if err != nil {
		log.Fatalln("failed to run api-gateway router. \nerror: ", err)
	}
}
