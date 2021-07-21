package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	gwruntime "github.com/grpc-ecosystem/grpc-gateway/runtime"
	loggermiddleware "github.com/meateam/api-gateway/logger"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/yarinBenisty/api-gateway/util"
	bpb "github.com/yarinBenisty/birthday-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

const (
	// ParamName parameter
	ParamName = "name"
	// ParamPersonalNumber parameter
	ParamPersonalNumber = "personalNumber"
	// ParamDate parameter
	ParamDate = "date"
)

// Birthday is struct of birthday-object
type Birthday struct {
	Name           string `json:"name"`
	Date           string `json:"date"`
	PersonalNumber string `json:"personalNumber"`
}

// Router connecting to the client
// functions in the bd-service proto file
type Router struct {
	client bpb.BirthdayFunctionsClient
	bpb.UnimplementedBirthdayFunctionsServer
	logger *logrus.Logger
}

// initClientConnection creates a new client and connects
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
		log.Fatal("failed to get mongo connection parameters")
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

	var birthdayFilter Birthday
	if bindErr := c.Bind(&birthdayFilter); bindErr != nil {
		c.String(http.StatusBadRequest, "create birthday method failed. \nerror: %s", bindErr)
		return
	}
	birthdayFilter = Birthday{
		Name:           strings.TrimSpace(c.Request.FormValue(ParamName)),
		Date:           strings.TrimSpace(c.Request.FormValue(ParamDate)),
		PersonalNumber: strings.TrimSpace(c.Request.FormValue(ParamPersonalNumber)),
	}
	request := &bpb.CreateBirthdayRequest{
		PersonalNumber: birthdayFilter.PersonalNumber,
		Name:           birthdayFilter.Name,
		Date:           birthdayFilter.Date,
	}

	res, err := r.client.CreateBirthday(c, request)
	if err != nil {
		httpStatusCode := gwruntime.HTTPStatusFromCode(status.Code(err))
		loggermiddleware.LogError(r.logger, c.AbortWithError(httpStatusCode, err))
		return
	}

	c.JSON(http.StatusOK, res)
}

// GetBirthday returns a birthday object
func (r *Router) GetBirthday(c *gin.Context) {

	request := &bpb.GetBirthdayRequest{PersonalNumber: c.Param(ParamPersonalNumber)}
	res, err := r.client.GetBirthday(c, request)
	if err != nil {
		httpStatusCode := gwruntime.HTTPStatusFromCode(status.Code(err))
		loggermiddleware.LogError(r.logger, c.AbortWithError(httpStatusCode, err))
		return
	}
	c.JSON(http.StatusOK, res)
}

//GetAllBirthdays returns all birthday objects
func (r *Router) GetAllBirthdays(c *gin.Context) {
	request := &bpb.GetAllBirthdaysRequest{}
	res, err := r.client.GetAllBirthdays(c, request)
	if err != nil {
		httpStatusCode := gwruntime.HTTPStatusFromCode(status.Code(err))
		loggermiddleware.LogError(r.logger, c.AbortWithError(httpStatusCode, err))
		return
	}
	c.JSON(http.StatusOK, res)
}

// DeleteBirthday deletes a birthday object by personal number
func (r *Router) DeleteBirthday(c *gin.Context) {
	request := &bpb.DeleteBirthdayRequest{PersonalNumber: c.Param(ParamPersonalNumber)}
	res, err := r.client.DeleteBirthday(c, request)
	if err != nil {
		httpStatusCode := gwruntime.HTTPStatusFromCode(status.Code(err))
		loggermiddleware.LogError(r.logger, c.AbortWithError(httpStatusCode, err))
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
	mainRouter.DELETE("/api/birthday/:personalNumber", r.DeleteBirthday)

	err = mainRouter.Run(":" + routerPort)
	if err != nil {
		fmt.Println("failed to run api gateway. \nrouter error: ", err)
		return
	}
}
