package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	
	// import Package pq is a pure Go Postgres
	_ "github.com/lib/pq"
)

type User struct {
	ID     int    `json:"id"`
	F_name string `json:"F_name"`
	L_name  string `json:"L_name"`
}

func newRouter(db *sql.DB) http.Handler {
	humans := "/humans"
	humansId := "/humans/{id}"
	//create router
	router := mux.NewRouter()
	router.HandleFunc("/", root()).Methods("GET")
	router.HandleFunc(humans, getUsers(db)).Methods("GET")
	router.HandleFunc(humansId, getUser(db)).Methods("GET")
	router.HandleFunc(humans, createUser(db)).Methods("POST")
	router.HandleFunc(humansId, updateUser(db)).Methods("PUT")
	router.HandleFunc(humansId, deleteUser(db)).Methods("DELETE")

	return jsonContentTypeMiddleware(router)
}

func main() {
	//connect to database
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//create the table if it doesn't exist
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS humans (id SERIAL PRIMARY KEY, F_name TEXT, L_name TEXT)")

	if err != nil {
		log.Fatal(err)
	}

	//start server
	log.Fatal(http.ListenAndServe(":8000", newRouter(db)))
}


func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
// root /
func root()http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("Create 			POST : localhost:8000/humans")
	json.NewEncoder(w).Encode("Read all human 	GET : localhost:8000/humans")
	json.NewEncoder(w).Encode("select human{id} GET : localhost:8000/humans/{id}")
	json.NewEncoder(w).Encode("Update 			PUT : localhost:8000/humans/{id}")
	json.NewEncoder(w).Encode("Delete 			DELETE : localhost:8000/humans/{id}")
	}
}
// get all users
func getUsers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT * FROM humans")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		users := []User{}
		for rows.Next() {
			var u User
			if err := rows.Scan(&u.ID, &u.F_name, &u.L_name); err != nil {
				log.Fatal(err)
			}
			users = append(users, u)
		}
		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(w).Encode(users)
	}
}

// get user by id
func getUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var u User
		err := db.QueryRow("SELECT * FROM humans WHERE id = $1", id).Scan(&u.ID, &u.F_name, &u.L_name)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(u)
	}
}

// create user
func createUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u User
		json.NewDecoder(r.Body).Decode(&u)

		err := db.QueryRow("INSERT INTO humans (F_name, L_name) VALUES ($1, $2) RETURNING id", u.F_name, u.L_name).Scan(&u.ID)
		if err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(w).Encode(u)
	}
}

// update user
func updateUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u User
		json.NewDecoder(r.Body).Decode(&u)

		vars := mux.Vars(r)
		id := vars["id"]

		_, err := db.Exec("UPDATE humans SET F_name = $1, L_name = $2 WHERE id = $3", u.F_name, u.L_name, id)
		if err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(w).Encode(u)
	}
}

// delete user
func deleteUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var u User
		err := db.QueryRow("SELECT * FROM humans WHERE id = $1", id).Scan(&u.ID, &u.F_name, &u.L_name)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			_, err := db.Exec("DELETE FROM humans WHERE id = $1", id)
			if err != nil {
				// Todo : fix error handling
				w.WriteHeader(http.StatusNotFound)
				return
			}

			json.NewEncoder(w).Encode("Humans deleted")
		}
	}
}
