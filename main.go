package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/yarinBenisty/api-gateway/util"
	bpb "github.com/yarinBenisty/birthday-service/proto"
	"google.golang.org/grpc"
)

const (
	// ParamName parameter
	ParamName = "name"
	// ParamPersonalNumber parameter
	ParamPersonalNumber = "personalNumber"
	// ParamDate parameter
	ParamDate = "date"
)

// Birthday is struct of bd-object
type Birthday struct {
	Name           string `json:"name"`
	Date           string `json:"date"`
	PersonalNumber string `json:"personalNumber"`
}

// Router that's connecting to the client
// functions in the bd-service proto file
type Router struct {
	client bpb.BirthdayFunctionsClient
	bpb.UnimplementedBirthdayFunctionsServer
}

// Func that creates a new client and connects
// to the bd-service client through the proto file
func initClientConnection() bpb.BirthdayFunctionsClient {

	address := viper.GetString(util.BirthdayServiceAddress)
	conn, err := grpc.Dial(
		address,
		grpc.WithInsecure(),
		grpc.FailOnNonTempDialError(true),
		grpc.WithBlock(),
	)
	if err != nil {
		fmt.Println("failed to get mongo connection parameters")
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

	var birthdayFilter Birthday
	request := &bpb.CreateBirthdayRequest{
		PersonalNumber: c.Request.FormValue(ParamPersonalNumber),
		Name:           c.Request.FormValue(ParamName),
		Date:           c.Request.FormValue(ParamDate),
	}

	if bindErr := c.Bind(&birthdayFilter); bindErr != nil {
		c.String(http.StatusBadRequest, "create birthday method failed. \nerror: %s", bindErr)
		return
	}

	res, err := r.client.CreateBirthday(c, request)
	if err != nil {
		c.String(http.StatusBadRequest, "failed creating birthday")
		return
	}

	c.JSON(http.StatusOK, res)
}

// GetBirthday returns a birthday object
func (r *Router) GetBirthday(c *gin.Context) {
	request := &bpb.GetBirthdayRequest{PersonalNumber: c.Param(ParamPersonalNumber)}
	res, err := r.client.GetBirthday(c, request)
	if err != nil {
		c.String(http.StatusBadRequest, "get birthday method failed. \nerror: %s", err)
		return
	}
	c.JSON(http.StatusOK, res)
}

//GetAllBirthdays returns all birthday objects
func (r *Router) GetAllBirthdays(c *gin.Context) {
	request := &bpb.GetAllBirthdaysRequest{}
	res, err := r.client.GetAllBirthdays(c, request)
	if err != nil {
		c.String(http.StatusBadRequest, "get all birthday method failed. \nerror: %s", err)
		return
	}
	c.JSON(http.StatusOK, res)
}

// DeleteBirthday deletes a birthday object by personal number
func (r *Router) DeleteBirthday(c *gin.Context) {
	request := &bpb.DeleteBirthdayRequest{PersonalNumber: c.Param(ParamPersonalNumber)}
	res, err := r.client.DeleteBirthday(c, request)
	if err != nil {
		c.String(http.StatusConflict, "delete birthday method failed. \nerror: %s", err)
		return
	}
	c.JSON(http.StatusNoContent, res)
}

func main() {

	// Loading dotenv file parameters
	err := util.LoadConfig()
	if err != nil {
		fmt.Println("cannot load config:", err)
		return
	}
	routerPort := viper.GetString(util.GrpcRouterPort)

	r := &Router{}
	r.client = initClientConnection()

	mainRouter := gin.Default()
	mainRouter.Use(
		cors.New(corsRouterConfig()),
	)

	mainRouter.POST("/api/birthday", r.CreateBirthday)
	mainRouter.GET("/api/birthday/:personalNumber", r.GetBirthday)
	mainRouter.GET("/api/birthdays", r.GetAllBirthdays)
	mainRouter.PUT("/api/birthday", r.CreateBirthday)
	mainRouter.DELETE("/api/birthday/:personalNumber", r.DeleteBirthday)

	err = mainRouter.Run(":" + routerPort)
	if err != nil {
		fmt.Println("failed to run api gateway. \nrouter error: ", err)
		return
	}
}
