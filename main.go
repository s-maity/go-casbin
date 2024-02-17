package main

import (
	"log"
	"net/http"

	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/go-casbin/controllers"
	"github.com/go-casbin/database"
	"github.com/labstack/echo/v4"
)

func main() {
	LoadAppConfig()

	// Initialize Database
	database.Connect(AppConfig.ConnectionString)
	database.Migrate()

	e := echo.New()
	adapter, err := gormadapter.NewAdapterByDB(DB)
	if err != nil {
		log.Fatalln("not able to connect")
	}
	// e.Use(Authenticate(adapter))
	e.Use(AuthenticateUser(adapter))

	// Users
	e.POST("/users", func(con echo.Context) error {
		return controllers.CreateUser(con)
	})
	e.GET("/users/groups", func(con echo.Context) error {
		return controllers.ListUsersGroup(con)
	})
	e.GET("/group/users", func(con echo.Context) error {
		return controllers.ListGroupUsers(con)
	})

	// Documents
	e.POST("/documents", func(con echo.Context) error {
		return controllers.CreateDocuent(con)
	})

	e.GET("/documents/by-users", func(con echo.Context) error {
		return controllers.GetAllDocumentsByUser(con)
	})

	e.GET("/documents", func(con echo.Context) error {
		return controllers.GetDocumentById(con)
	})

	e.POST("/documents/make-public", func(con echo.Context) error {
		return controllers.MakeDocumentsPublic(con)
	})

	// Permissions
	e.POST("/create-group", func(c echo.Context) error {
		return c.JSON(http.StatusOK, controllers.CreateGroup(c))
	})

	e.POST("/grant-access", func(c echo.Context) error {
		return c.JSON(http.StatusOK, controllers.AddPermission(c))
	})

	e.Logger.Fatal(e.Start("0.0.0.0:3000"))
}
