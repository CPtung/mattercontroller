package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type APIServer interface {
	Start() error
	Stop() error
	Router() *gin.Engine
}

type APIServerImpl struct {
	http.Server
	logger *logrus.Entry
}

var (
	runtimePath = "/var/run/matter"
	unixPath    = runtimePath + "/matter.sock"
)

func Address() string {
	return unixPath
}

// recovery is gin middleware to recovery from panic
// To be refactor
func recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				switch v := err.(type) {
				case HttpError:
					c.JSON(v.Status(), gin.H{"error": gin.H{"code": v.Code(), "message": v.Error()}})
				case string:
					c.JSON(http.StatusExpectationFailed, gin.H{"error": gin.H{"code": 500, "message": v}})
				default:
					httprequest, _ := httputil.DumpRequest(c.Request, false)
					reset := string([]byte{27, 91, 48, 109})
					fmt.Printf("panic recovered:\n\n%s%v\n\n%s", httprequest, err, reset)
					c.JSON(http.StatusInternalServerError, gin.H{"message": "internal error"})
				}
				c.Abort()
			}
		}()
		c.Next() // execute all the handlers
	}
}

func NewAPIServer() APIServer {
	if _, err := os.Stat(runtimePath); os.IsNotExist(err) {
		os.MkdirAll(runtimePath, 0750)
	}
	return &APIServerImpl{
		Server: http.Server{
			Addr:    "unix://" + unixPath,
			Handler: gin.Default(),
		},
		logger: logrus.WithFields(logrus.Fields{"origin": "matter"}),
	}
}

func (s *APIServerImpl) Start() error {
	// Add recovery middleware
	if ginEngine, ok := s.Handler.(*gin.Engine); ok {
		ginEngine.Use(recovery())
	}
	// Remove the socket file if it already exists
	if _, err := os.Stat(unixPath); os.IsExist(err) {
		os.Remove(unixPath)
	}
	listener, err := net.Listen("unix", unixPath)
	if err != nil {
		return err
	}
	// Set the socket permission
	os.Chmod(unixPath, 750)
	// Start the server in a goroutine
	go func() {
		defer listener.Close()
		if err := s.Serve(listener); err != nil && err != http.ErrServerClosed {
			s.logger.Errorf("listen: %s", err)
			return
		}
	}()
	return nil
}

func (s *APIServerImpl) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil && err != http.ErrServerClosed {
		return err
	}
	s.logger.Info("shutdown server...")
	return nil
}

func (s *APIServerImpl) Router() *gin.Engine {
	if ginEngine, ok := s.Handler.(*gin.Engine); ok {
		return ginEngine
	}
	return nil
}
