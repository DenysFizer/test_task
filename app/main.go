package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"test_task/config"
	"test_task/pkg/db"
	"test_task/pkg/handlers"
)

func main() {

	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		panic("Error to load config file")
	}

	db.InitDB(cfg)
	defer db.CloseDB()

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	handlers.InitHandlers(cfg)

	e.GET("/hello", handlers.Hello)
	e.GET("/cats", handlers.GetCats)
	e.POST("/cats", handlers.CreateCat)
	e.PUT("/cats/:id", handlers.UpdateCatSalary)
	e.GET("/missions", handlers.GetMissions)
	e.POST("/missions", handlers.CreateMission)
	e.PUT("/missions/:id", handlers.UpdateMission)
	e.DELETE("/missions/:id", handlers.DeleteMission)
	e.PUT("/targets/:id", handlers.UpdateTarget)
	e.PUT("/complete/:id", handlers.UpdateMissionComplete)
	e.DELETE("/targets/:id", handlers.DeleteTarget)
	e.GET("/targets", handlers.GetTargets)
	e.PUT("/updatenotes/:id", handlers.UpdateNotes)
	e.POST("/addtarget/:id", handlers.AddTargetToMission)
	e.PUT("/assigncat", handlers.AssignCatToMission)

	e.Logger.Fatal(e.Start(":1323"))
}
