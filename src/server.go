package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	Router *gin.Engine
	Server *http.Server
}

func NewServer() *Server {
	router := gin.New()

	router.Use(gin.Logger(), gin.Recovery())

	host := os.Getenv("HOST")
	if host == "" {
		host = "0.0.0.0"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

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
	r := server.Router.Group("/health")
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
	r.Static("/assets", "./assets")
	r.LoadHTMLGlob("templates/**/*.html")
	r.GET("/", func (c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"Title": "HelloGIN",
		})
	})
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
