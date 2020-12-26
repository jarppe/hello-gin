package main

import (
	"context"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/jarppe/hello-gin/assets"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"
)

var (
	host, port, resources string
)

func init() {
	if host = os.Getenv("HOST"); host == "" {
		host = "0.0.0.0"
	}
	if port = os.Getenv("PORT"); port == "" {
		port = "8080"
	}
	if resources = os.Getenv("RESOURCES"); resources == "" {
		resources = "."
	}
}

type Server struct {
	Router *gin.Engine
	Server *http.Server
}

func NewServer() *Server {
	router := gin.New()

	router.Use(gin.Logger(), gin.Recovery())

	server := &http.Server{
		Addr:    host + ":" + port,
		Handler: router,
	}

	return &Server{
		Router: router,
		Server: server,
	}
}

func (server *Server) Run() {
	go func() {
		if err := server.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
}

func (server *Server) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
}


func (server *Server) Routes() {
	server.HealthRoutes()
	server.ApiRoutes()
	server.AppRoutes()
}

func (server *Server) HealthRoutes() {
	r := server.Router
	r.GET("/livez", server.Live())
	r.GET("/readyz", server.Ready())
}

func (server *Server) Live() gin.HandlerFunc {
	return func (c *gin.Context) {
		c.Status(200)
	}
}

func (server *Server) Ready() gin.HandlerFunc {
	return func (c *gin.Context) {
		c.Status(200)
	}
}

func (server *Server) ApiRoutes() {
	r := server.Router.Group("/api")
	r.GET("/ping", server.Ping())
	r.POST("/hello", server.Hello())
}

func (server *Server) Ping() gin.HandlerFunc {
	return func (c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	}
}

func (server *Server) Hello() gin.HandlerFunc {
	type FormA struct {
		Foo string `json:"foo"`
	}
	return func (c *gin.Context) {
		foo := &FormA{}
		if err := c.ShouldBind(foo); err != nil {
			c.String(http.StatusBadRequest, `the body should be formA`)
			return
		}
		name := c.Param("name")
		c.String(http.StatusOK, "%s %s", foo.Foo, name)
	}
}

func (server *Server) AppRoutes() {
	r := server.Router

	r.LoadHTMLGlob(path.Join(resources, "templates/**/*.html"))
	r.GET("/", gzip.Gzip(1), func (c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"Title": "HelloGIN",
		})
	})

	r.GET("/assets/*asset", assets.NewAssetsHandler(path.Join(resources, "assets")))
}

func (server *Server) WaitSignal() {
	sigCh := make(chan os.Signal)
	signal.Notify(sigCh,
		syscall.SIGINT,
		syscall.SIGHUP,
		syscall.SIGTERM,
		syscall.SIGUSR2)
	sig := <-sigCh
	signal.Reset(sig)

	log.Printf("Got signal %q, terminating...", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
}
