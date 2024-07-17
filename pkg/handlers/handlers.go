package handlers

import (
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v4"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"test_task/config"
	"test_task/models"
	"test_task/pkg/db"
)

var cfg *config.Config

func InitHandlers(c *config.Config) {
	cfg = c
}

func Hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func GetCats(c echo.Context) error {
	rows, err := db.DB.Query(context.Background(), "SELECT id, name, years_of_experience, breed, salary FROM cats")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	defer rows.Close()

	var cats []models.Cat
	for rows.Next() {
		var cat models.Cat
		err := rows.Scan(&cat.ID, &cat.Name, &cat.YearsOfExperience, &cat.Breed, &cat.Salary)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		cats = append(cats, cat)
	}

	return c.JSON(http.StatusOK, cats)
}

func CreateCat(c echo.Context) error {
	cat := new(models.Cat)
	if err := c.Bind(cat); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if !isValidBreed(cat.Breed) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid breed"})
	}

	err := db.DB.QueryRow(context.Background(),
		"INSERT INTO cats (name, years_of_experience, breed, salary) VALUES ($1, $2, $3, $4) RETURNING id",
		cat.Name, cat.YearsOfExperience, cat.Breed, cat.Salary).Scan(&cat.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, cat)
}
func isValidBreed(breed string) bool {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.thecatapi.com/v1/breeds", nil)
	if err != nil {
		return false
	}

	req.Header.Set("x-api-key", cfg.APIKey)
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false
	}

	var breeds []struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&breeds); err != nil {
		return false
	}

	for _, b := range breeds {
		if b.Name == breed {
			return true
		}
	}

	return false
}

func GetMissions(c echo.Context) error {
	rows, err := db.DB.Query(context.Background(), "SELECT id, cat_id, complete FROM missions")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	defer rows.Close()

	var missions []models.Mission
	for rows.Next() {
		var mission models.Mission
		err := rows.Scan(&mission.ID, &mission.CatID, &mission.Complete)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		missions = append(missions, mission)
	}

	return c.JSON(http.StatusOK, missions)
}

func CreateMission(c echo.Context) error {
	mission := new(models.Mission)
	if err := c.Bind(mission); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	// Check if the cat already has an incomplete mission
	var existingMissionID int
	err := db.DB.QueryRow(context.Background(), "SELECT id FROM missions WHERE cat_id = $1 AND complete = FALSE", mission.CatID).Scan(&existingMissionID)
	if err == nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Cat already has an incomplete mission"})
	} else if err != pgx.ErrNoRows {
		return c.JSON(http.StatusInternalServerError, err)
	}

	// Insert the mission
	err = db.DB.QueryRow(context.Background(), "INSERT INTO missions (cat_id, complete) VALUES ($1, $2) RETURNING id", mission.CatID, mission.Complete).Scan(&mission.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	// Insert the targets
	for _, target := range mission.Targets {
		_, err := db.DB.Exec(context.Background(),
			"INSERT INTO targets (mission_id, name, country, notes, complete) VALUES ($1, $2, $3, $4, $5)",
			mission.ID, target.Name, target.Country, target.Notes, target.Complete)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
	}

	return c.JSON(http.StatusCreated, mission)
}

func UpdateMission(c echo.Context) error {
	id := c.Param("id")
	mission := new(models.Mission)
	if err := c.Bind(mission); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	_, err := db.DB.Exec(context.Background(),
		"UPDATE missions SET cat_id = $1, complete = $2 WHERE id = $3",
		mission.CatID, mission.Complete, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, mission)
}

func UpdateCatSalary(c echo.Context) error {
	catID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid cat ID"})
	}

	type SalaryUpdate struct {
		Salary int `json:"salary"`
	}
	var newSalary SalaryUpdate
	if err := c.Bind(&newSalary); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid data format"})
	}

	_, err = db.DB.Exec(context.Background(),
		"UPDATE cats SET salary = $1 WHERE id = $2",
		newSalary.Salary, catID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	var updatedCat models.Cat
	err = db.DB.QueryRow(context.Background(),
		"SELECT id, name, years_of_experience, breed, salary FROM cats WHERE id = $1", catID).
		Scan(&updatedCat.ID, &updatedCat.Name, &updatedCat.YearsOfExperience, &updatedCat.Breed, &updatedCat.Salary)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, updatedCat)
}

func DeleteMission(c echo.Context) error {
	id := c.Param("id")

	var catID int
	err := db.DB.QueryRow(context.Background(), "SELECT cat_id FROM missions WHERE id = $1", id).Scan(&catID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Mission not found"})
		}
		return c.JSON(http.StatusInternalServerError, err)
	}

	if catID != 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Mission cannot be deleted because it is assigned to a cat"})
	}

	_, err = db.DB.Exec(context.Background(), "DELETE FROM missions WHERE id = $1", id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusNoContent)
}

func UpdateTarget(c echo.Context) error {
	id := c.Param("id")
	target := new(models.Target)
	if err := c.Bind(target); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	var complete bool
	err := db.DB.QueryRow(context.Background(), "SELECT complete FROM targets WHERE id = $1", id).Scan(&complete)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	if complete {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Cannot update a completed target"})
	}

	_, err = db.DB.Exec(context.Background(),
		"UPDATE targets SET name = $1, country = $2, notes = $3, complete = $4 WHERE id = $5",
		target.Name, target.Country, target.Notes, target.Complete, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, target)
}
func UpdateMissionComplete(c echo.Context) error {
	id := c.Param("id")

	var complete bool
	err := db.DB.QueryRow(context.Background(), "SELECT complete FROM missions WHERE id = $1", id).Scan(&complete)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	if complete {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Mission is already completed"})
	}

	_, err = db.DB.Exec(context.Background(),
		"UPDATE missions SET complete = $1 WHERE id = $2",
		true, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, map[string]string{"result": "Mission completed"})
}
func DeleteTarget(c echo.Context) error {
	id := c.Param("id")

	// Check if the target is complete
	var complete bool
	err := db.DB.QueryRow(context.Background(), "SELECT complete FROM targets WHERE id = $1", id).Scan(&complete)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	if complete {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Cannot delete a completed target"})
	}

	// Check if the mission is completed
	var missionComplete bool
	err = db.DB.QueryRow(context.Background(), "SELECT m.complete FROM missions m JOIN targets t ON m.id = t.mission_id WHERE t.id = $1", id).Scan(&missionComplete)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	if missionComplete {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Cannot delete a target from a completed mission"})
	}

	// Delete the target
	_, err = db.DB.Exec(context.Background(), "DELETE FROM targets WHERE id = $1", id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusNoContent)
}

func GetTargets(c echo.Context) error {
	rows, err := db.DB.Query(context.Background(), "SELECT id, mission_id, name, country, notes, complete FROM targets")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	defer rows.Close()

	var targets []models.Target
	for rows.Next() {
		var target models.Target
		err := rows.Scan(&target.ID, &target.MissionID, &target.Name, &target.Country, &target.Notes, &target.Complete)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		targets = append(targets, target)
	}

	return c.JSON(http.StatusOK, targets)
}
func UpdateNotes(c echo.Context) error {
	id := c.Param("id")
	note := new(models.Note)
	if err := c.Bind(note); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	var targetComplete bool
	err := db.DB.QueryRow(context.Background(), "SELECT complete FROM targets WHERE id = $1", id).Scan(&targetComplete)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	if targetComplete {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Cannot update notes for a completed target"})
	}

	var missionComplete bool
	err = db.DB.QueryRow(context.Background(), "SELECT m.complete FROM missions m JOIN targets t ON m.id = t.mission_id WHERE t.id = $1", id).Scan(&missionComplete)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	if missionComplete {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Cannot update notes for a target in a completed mission"})
	}

	_, err = db.DB.Exec(context.Background(),
		"UPDATE targets SET notes = $1 WHERE id = $2",
		note.Notes, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, note)
}
func AddTargetToMission(c echo.Context) error {
	missionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid mission ID"})
	}

	// Перевірка, чи місія вже виконана
	var missionComplete bool
	err = db.DB.QueryRow(context.Background(), "SELECT complete FROM missions WHERE id = $1", missionID).Scan(&missionComplete)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	if missionComplete {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Cannot add a target to a completed mission"})
	}

	target := new(models.Target)
	if err := c.Bind(target); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	_, err = db.DB.Exec(context.Background(),
		"INSERT INTO targets (mission_id, name, country, notes, complete) VALUES ($1, $2, $3, $4, $5)",
		missionID, target.Name, target.Country, target.Notes, target.Complete)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	targetInfo := map[string]interface{}{
		"name":     target.Name,
		"country":  target.Country,
		"notes":    target.Notes,
		"complete": target.Complete,
	}

	return c.JSON(http.StatusCreated, targetInfo)
}

func AssignCatToMission(c echo.Context) error {

	type AssignmentRequest struct {
		CatID     int `json:"cat_id"`
		MissionID int `json:"mission_id"`
	}

	req := new(AssignmentRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	var catExists bool
	err := db.DB.QueryRow(context.Background(), "SELECT EXISTS (SELECT 1 FROM cats WHERE id = $1)", req.CatID).Scan(&catExists)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	if !catExists {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Cat not found"})
	}

	var missionCompleted bool
	err = db.DB.QueryRow(context.Background(), "SELECT complete FROM missions WHERE id = $1", req.MissionID).Scan(&missionCompleted)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	if missionCompleted {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Cannot assign a cat to a completed mission"})
	}

	_, err = db.DB.Exec(context.Background(),
		"UPDATE missions SET cat_id = $1 WHERE id = $2",
		req.CatID, req.MissionID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Cat assigned to mission successfully"})
}
