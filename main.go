package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	graphql "github.com/graph-gophers/graphql-go"
	"github.com/itsjamie/gin-cors"
)

type resolverNotFoundError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e resolverNotFoundError) Error() string {
	return fmt.Sprintf("Error [%s]: %s", e.Code, e.Message)
}

func (e resolverNotFoundError) Extensions() map[string]interface{} {
	return map[string]interface{}{
		"code":    e.Code,
		"message": e.Message,
	}
}

type query struct{}

func (*query) Today() graphql.Time {
	// 用 graphql.Time 包過再回傳
	return graphql.Time{Time: time.Now()}
}

func (*query) DistanceOfTimeToNowInWords(input struct{ From graphql.Time }) string {
	// 接收到的 input.From 為 graphql.Time 的型別，能夠 input.From.Time 拿到 time.Time 型別的數值。
	return fmt.Sprintf("%d second(s)", time.Since(input.From.Time)/time.Second)
}

func (*query) TestError() (string, error) {
	return "", resolverNotFoundError{
		Code:    "NotFound",
		Message: "This is not the droid you are looking for",
	}
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	// Apply the middleware to the router (works with groups too)
	r.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}))

	s := `
	schema {
		query: Query
	}
	type Query {
		# 測試回傳數值為 type Time，框架會如何去處理
		today: Time!

		# 該時間距今的距離
		# 測試拿到 argument 為 type Time，框架會如何去處理
		distanceOfTimeToNowInWords(from: Time!): String!
		testError: String
	}
	
	# represent time format in ISO8601
	scalar Time
	`
	schema := graphql.MustParseSchema(s, &query{})

	r.POST("/query", func(c *gin.Context) {
		var params struct {
			Query         string                 `json:"query"`
			OperationName string                 `json:"operationName"`
			Variables     map[string]interface{} `json:"variables"`
		}

		if err := c.BindJSON(&params); err != nil {
			log.Fatal(err)
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		c.JSON(
			http.StatusOK,
			schema.Exec(c, params.Query, params.OperationName, params.Variables),
		)
	})
	return r
}

func main() {
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
