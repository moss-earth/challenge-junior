package middleware

import (
	"challenge-junior/models" // models eh a pasta que definimos Project schema
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http" // isso aqui lembra o "req, resp" do nodeJs
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // postgres driver para go
)

// formato da resposta
type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

// cria conexao com o DB que linkei no aquivo .env
func createConnection() *sql.DB {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))

	if err != nil {
		panic(err)
	}

	err = db.Ping()

	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
	return db
}

// CreateProject eh para criar project
func CreateProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	var project models.Project

	err := json.NewDecoder(r.Body).Decode(&project)

	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}

	insertID := insertProject(project)

	res := response{
		ID:      insertID,
		Message: "Project created successfully",
	}

	json.NewEncoder(w).Encode(res)
}

// GetProject vai trazer um projeto pelo id
func GetProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	project, err := getProject(int64(id))

	if err != nil {
		log.Fatalf("Unable to get project. %v", err)
	}

	json.NewEncoder(w).Encode(project)
}

// GetAllProject vai retornar todos os project
func GetAllProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	projects, err := getAllProjects()

	if err != nil {
		log.Fatalf("Unable to get all project. %v", err)
	}

	json.NewEncoder(w).Encode(projects)
}

// UpdateProject atualiza o project
func UpdateProject(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	var project models.Project

	err = json.NewDecoder(r.Body).Decode(&project)

	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}

	updatedRows := updateProject(int64(id), project)

	msg := fmt.Sprintf("Project updated successfully. Total rows/record affected %v", updatedRows)

	res := response{
		ID:      int64(id),
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)
}

// DeleteProject deleta as informacoes do project no db
func DeleteProject(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	deletedRows := deleteProject(int64(id))

	msg := fmt.Sprintf("Project updated successfully. Total rows/record affected %v", deletedRows)

	res := response{
		ID:      int64(id),
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)
}

//------------------------- funcoes handler ----------------
// inserir um project no DB
func insertProject(project models.Project) int64 {

	db := createConnection()

	defer db.Close()

	sqlStatement := `INSERT INTO projects (name) VALUES ($1) RETURNING projectid`

	var id int64

	err := db.QueryRow(sqlStatement, project.Name).Scan(&id)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	fmt.Printf("Inserted a single record %v", id)

	return id
}

// pega um project do DB pelo projectid
func getProject(id int64) (models.Project, error) {

	db := createConnection()

	defer db.Close()

	var project models.Project

	sqlStatement := `SELECT * FROM projects WHERE projectid=$1`

	row := db.QueryRow(sqlStatement, id)

	err := row.Scan(&project.ID, &project.Name)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return project, nil
	case nil:
		return project, nil
	default:
		log.Fatalf("Unable to scan the row. %v", err)
	}

	return project, err
}

// pega todos os projects do DB
func getAllProjects() ([]models.Project, error) {

	db := createConnection()

	defer db.Close()

	var projects []models.Project

	sqlStatement := `SELECT * FROM projects`

	rows, err := db.Query(sqlStatement)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var project models.Project

		err = rows.Scan(&project.ID, &project.Name)

		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}

		projects = append(projects, project)

	}

	return projects, err
}

// atualizar project no DB
func updateProject(id int64, project models.Project) int64 {

	db := createConnection()

	defer db.Close()

	sqlStatement := `UPDATE projects SET name=$2 WHERE Projectid=$1`

	res, err := db.Exec(sqlStatement, id, project.Name)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	fmt.Printf("Total rows/record affected %v", rowsAffected)

	return rowsAffected
}

// deleta project no DB
func deleteProject(id int64) int64 {

	db := createConnection()

	defer db.Close()

	sqlStatement := `DELETE FROM projects WHERE projectid=$1`

	res, err := db.Exec(sqlStatement, id)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	fmt.Printf("Total rows/record affected %v", rowsAffected)

	return rowsAffected
}

// CreateLot cria um novo lot no db
func CreateLot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	var lot models.Lot

	err := json.NewDecoder(r.Body).Decode(&lot)

	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}

	insertLot := insertLot(lot)

	res := response{
		ID:      insertLot,
		Message: "Lot created successfully",
	}

	json.NewEncoder(w).Encode(res)
}

// GetLot vai retornar um lot pelo lotid
func GetLot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	params := mux.Vars(r)

	lotid, err := strconv.Atoi(params["lotid"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	lot, err := getLot(int64(lotid))

	if err != nil {
		log.Fatalf("Unable to get lot. %v", err)
	}

	// send the response
	json.NewEncoder(w).Encode(lot)
}

// GetAllLot pega todos os lot de um projeto
func GetAllLot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	lots, err := getAllLots()

	if err != nil {
		log.Fatalf("Unable to get all lot. %v", err)
	}

	json.NewEncoder(w).Encode(lots)
}

// DeleteLot deleta um lot do db
func DeleteLot(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	params := mux.Vars(r)

	lotid, err := strconv.Atoi(params["lotid"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	deletedRows := deleteLot(int64(lotid))

	msg := fmt.Sprintf("Lot delete successfully. Total rows/record affected %v", deletedRows)

	res := response{
		ID:      int64(lotid),
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)
}

//------------------------- handler das funcoes de lot  ----------------
// insert um lot no DB
func insertLot(lot models.Lot) int64 {

	db := createConnection()

	defer db.Close()

	sqlStatement := `INSERT INTO lots (price, quantity, buydate, projectID) VALUES ($1, $2, $3, $4) RETURNING lotid`

	var lotid int64

	err := db.QueryRow(sqlStatement, lot.Price, lot.Quantity, lot.Buydate, lot.ProjectID).Scan(&lotid)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	fmt.Printf("Inserted a single record %v", lotid)

	return lotid
}

// get um lot pelo lotid
func getLot(lotid int64) (models.Lot, error) {

	db := createConnection()

	defer db.Close()

	var lot models.Lot

	sqlStatement := `SELECT * FROM lots WHERE lotid=$1`

	row := db.QueryRow(sqlStatement, lotid)

	err := row.Scan(&lot.LotID, &lot.Price, &lot.Quantity, &lot.Buydate, &lot.ProjectID)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return lot, nil
	case nil:
		return lot, nil
	default:
		log.Fatalf("Unable to scan the row. %v", err)
	}

	return lot, err
}

// get pega todos os lot de um project
func getAllLots() ([]models.Lot, error) {

	db := createConnection()

	defer db.Close()

	var lots []models.Lot

	sqlStatement := `SELECT * FROM lots WHERE projectID=$1`

	rows, err := db.Query(sqlStatement)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var lot models.Lot

		err = rows.Scan(&lot.LotID, &lot.Price, &lot.Quantity, &lot.Buydate, &lot.ProjectID)

		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}

		lots = append(lots, lot)

	}

	return lots, err
}

// delete um lot no DB
func deleteLot(lotid int64) int64 {

	db := createConnection()

	defer db.Close()

	sqlStatement := `DELETE FROM lots WHERE lotid=$1`

	res, err := db.Exec(sqlStatement, lotid)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	fmt.Printf("Total rows/record affected %v", rowsAffected)

	return rowsAffected
}
