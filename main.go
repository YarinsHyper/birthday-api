package main

import (
	"log"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/yarinBenisty/api-gateway/util"
	bpb "github.com/yarinBenisty/birthday-service/proto"
	"google.golang.org/grpc"
)

const (
	// Name parameter
	Name = "name"
	// PersonalNumber parameter
	PersonalNumber = "personalNumber"
	// Date parameter
	Date = "date"
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

	// Loading the dotenv parameters
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	address := config.BirthdayServiceAddress
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

// CreateBirthday inserts a birthday object if it
// doesn't exist. if it does, its being overritten
func (r *Router) CreateBirthday(c *gin.Context) {
	request := &bpb.CreateBirthdayRequest{
		PersonalNumber: c.Request.FormValue(PersonalNumber),
		Name:           c.Request.FormValue(Name),
		Date:           c.Request.FormValue(Date),
	}
	res, err := r.client.CreateBirthday(c, request)
	if err != nil {
		log.Fatal("create birthday error: ", err)
		c.String(400, "create birthday method failed. \nerror: %s", err)
	}
	c.JSON(200, res)
}

// GetBirthday returns a birthday object
func (r *Router) GetBirthday(c *gin.Context) {
	request := &bpb.GetBirthdayRequest{PersonalNumber: c.Query(PersonalNumber)}
	res, err := r.client.GetBirthday(c, request)
	if err != nil {
		log.Fatal("get birthday error: ", err)
		c.String(404, "get birthday method failed. \nerror: %s", err)
	}
	c.JSON(200, res)
}

//GetAllBirthdays returns all birthday objects
func (r *Router) GetAllBirthdays(c *gin.Context) {
	request := &bpb.GetAllBirthdaysRequest{}
	res, err := r.client.GetAllBirthdays(c, request)
	if err != nil {
		log.Fatal("get all birthdays error.\n", err)
		c.String(404, "get all birthday method failed. \nerror: %s", err)
	}
	c.JSON(200, res)
}

// DeleteBirthday deletes a birthday object by personal number
func (r *Router) DeleteBirthday(c *gin.Context) {
	request := &bpb.DeleteBirthdayRequest{PersonalNumber: PersonalNumber}
	res, err := r.client.DeleteBirthday(c, request)
	if err != nil {
		log.Fatal("delete birthday error: ", err)
		c.String(404, "delete birthday method failed. \nerror: %s", err)
	}
	c.JSON(204, res)
}

func main() {

	// Loading the dotenv parameters
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	routerPort := config.GrpcRouterPort

	r := &Router{}
	r.client = initClientConnection()

	mainRouter := gin.Default()
	mainRouter.Use(
		cors.New(corsRouterConfig()),
	)

	mainRouter.POST("/api/createBirthday", r.CreateBirthday)
	mainRouter.GET("/api/getBirthday", r.GetBirthday)
	mainRouter.GET("/api/getAllBirthdays", r.GetAllBirthdays)
	mainRouter.POST("/api/updateBirthday", r.CreateBirthday)
	mainRouter.DELETE("/api/deleteBirthday", r.DeleteBirthday)

	err = mainRouter.Run(":" + routerPort)
	if err != nil {
		log.Fatalln("failed to run api-gateway router. \nerror: ", err)
	}
}
