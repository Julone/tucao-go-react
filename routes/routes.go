package routes

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/klauspost/compress/zip"
	"github.com/segmentio/kafka-go"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
	"tuxiaocao/pkg/logger"
	controllers2 "tuxiaocao/routes/controllers"
)

// PublicRoutes func for describe group of public routes.
func PublicRoutes(app *fiber.App) {

	pubRoute := app.Group("/api/v1")
	pubRoute.Get("/", func(ctx *fiber.Ctx) error {
		//ctx.Response().Header.Set("Cache-Control", "no-cache")
		str := time.Now().Format(time.DateTime)

		ctx.Send([]byte(str))
		return nil
	})
	// Routes for GET method:
	pubRoute.Get("/products", controllers2.Getproducts)   // get list of all products
	pubRoute.Get("/product/:id", controllers2.Getproduct) // get one product by ID
	// Routes for POST method:
	pubRoute.Post("/user/sign/up", controllers2.UserSignUp) // register app new user
	pubRoute.Post("/user/sign/in", controllers2.UserSignIn) // auth, return Access & Refresh tokens

	// Create routes group.
	route := app.Group("/api/v1")
	// Routes for POST method:
	route.Post("/product", controllers2.Createproduct)     // create app new product
	route.Post("/user/sign/out", controllers2.UserSignOut) // de-authorization user
	route.Post("/token/renew", controllers2.RenewTokens)   // renew Access & Refresh tokens
	// Routes for PUT method:
	route.Put("/product", controllers2.Updateproduct) // update one product by ID
	// Routes for DELETE method:
	route.Delete("/product", controllers2.Deleteproduct) // delete one product by ID

	route.Get("/kafka", func(ctx *fiber.Ctx) error {
		topic := "my-topic"
		partition := 0
		// 连接至Kafka集群的Leader节点
		conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:29092", topic, partition)
		if err != nil {
			logger.Log.Fatal("failed to dial leader:", err)
		}
		// 设置发送消息的超时时间
		conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
		// 发送消息
		_, err = conn.WriteMessages(
			kafka.Message{Value: []byte("one!")},
			kafka.Message{Value: []byte("two!")},
			kafka.Message{Value: []byte("three!")},
		)
		if err != nil {
			logger.Log.Fatal("failed to write messages:", err)
		}

		return nil

	})

	route.Put("/update/Frontend", func(ctx *fiber.Ctx) error {
		f, err := ctx.FormFile("file")
		if err != nil || (f != nil && !strings.HasSuffix(f.Filename, ".zip")) {
			return ctx.JSON(fiber.Map{"error": true, "msg": "nedd receive zip"})
		}
		uploadFileName := strings.ReplaceAll(time.Now().Format(time.DateTime), ":", "_")
		targetFileFile := "./saving/" + uploadFileName + ".zip"
		targetFileDir := "./saving/" + uploadFileName
		frontendDir := "./frontendSource"

		err = ctx.SaveFile(f, targetFileFile)
		if err != nil {
			return ctx.JSON(fiber.Map{"error2": true, "msg": err.Error()})
		}
		err = unzip(targetFileFile, targetFileDir)
		if err != nil {
			return ctx.JSON(fiber.Map{"error2": true, "msg": err.Error()})
		}

		cmd := exec.Command("powershell", "rm", "-r", "./frontendSource")
		result, err := cmd.Output()
		//err = os.RemoveAll("./frontendSource")
		logger.Log.Info(string(result))
		if err != nil {
			return ctx.JSON(fiber.Map{"error2": true, "msg": string(result) + err.Error()})
		}
		err = os.Rename(targetFileDir+"/dist", frontendDir)
		if err != nil {
			return ctx.JSON(fiber.Map{"error2": true, "msg": err.Error()})
		}
		defer func() {
			os.RemoveAll(targetFileFile)
			os.RemoveAll(targetFileDir)
		}()
		return nil
	}) // delete one product by ID

	// Create routes group.
	swagRoute := app.Group("/swagger")
	// Routes for GET method:
	swagRoute.Get("*", swagger.HandlerDefault) // get one user by ID
	app.Use("*", NotFoundHandler)
}

func unzip(source, destination string) error {
	zipReader, err := zip.OpenReader(source)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	for _, file := range zipReader.File {
		filePath := filepath.Join(destination, file.Name)

		if file.FileInfo().IsDir() {
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return err
		}

		zipFile, err := file.Open()
		if err != nil {
			return err
		}
		defer zipFile.Close()

		targetFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, zipFile); err != nil {
			return err
		}
	}

	return nil
}

// NotFoundHandler func for describe 404 Error route.
func NotFoundHandler(c *fiber.Ctx) error {
	// Register new special route.
	// Return HTTP 404 status and JSON response.
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"error": true,
		"msg":   "sorry, endpoint is not found",
	})
}
