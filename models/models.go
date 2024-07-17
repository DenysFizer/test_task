package models

type Cat struct {
	ID                int     `json:"id"`
	Name              string  `json:"name"`
	YearsOfExperience int     `json:"years_of_experience"`
	Breed             string  `json:"breed"`
	Salary            float64 `json:"salary"`
}

type Mission struct {
	ID       int      `json:"id"`
	CatID    int      `json:"cat_id"`
	Complete bool     `json:"complete"`
	Targets  []Target `json:"targets"`
}

type Target struct {
	ID        int    `json:"id"`
	MissionID int    `json:"mission_id"`
	Name      string `json:"name"`
	Country   string `json:"country"`
	Notes     string `json:"notes"`
	Complete  bool   `json:"complete"`
}
type Note struct {
	Notes string `json:"notes"`
}
