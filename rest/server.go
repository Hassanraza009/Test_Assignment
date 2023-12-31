package rest

import (
	"bytes"
	"encoding/json"
	"io"
	"time"

	"test/logger"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type HttpServer struct {
	addr           string
	userController UserController
	middleware     MiddleWare
}

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

// NewHttpServer create server instance
func NewHttpServer(addr string,
) *HttpServer {
	return &HttpServer{
		addr: addr,
	}
}

func loggingMiddleware(next gin.HandlerFunc) gin.HandlerFunc {

	return func(c *gin.Context) {
		start := time.Now()
		w := &responseBodyWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = w

		c.Next()

		responseString := w.body.String()
		var response StandardResponse
		err := json.Unmarshal([]byte(responseString), &response)

		if err != nil {
			logger.Instance().Debug("Failed to record the response")
		}

		errs := c.Errors.JSON()

		if c.Errors != nil {
			if len(c.Errors) > 0 {
				errs = c.Errors[0].JSON()
			}
		}
		corId := ""

		logrusFields := logrus.Fields{
			"corelationId":        corId,
			"request_method":      c.Request.Method,
			"request_path":        c.Request.URL.Path,
			"client_ip":           c.ClientIP(),
			"latency_nanoseconds": time.Since(start),
			"response":            response,
			"error":               errs,
		}
		logger.Instance().WithFields(logrusFields).Debug("Request Details")

	}
}

func (server *HttpServer) Start() {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard

	r := gin.Default()
	r.Use(loggingMiddleware(r.HandleContext))

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	AttachRoutes(r, server)

	err := r.Run(server.addr)
	if err != nil {
		panic(err)
	}

}
