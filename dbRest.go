package main


import (
    "fmt"
    "log"
    "net/http"
	"encoding/json"
	"github.com/gorilla/mux"
    "io/ioutil"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    
)

type user struct {
	Id      string `json:"Id"`
    Name string `json:"Name"`
    City string `json:"City"`
}

type users struct {
    users []user `json:"users"`
}

func dbConn() (db *sql.DB) {
    dbDriver := "mysql"
    dbUser := "root"
    dbPass := "root"
    dbName := "demo"
    db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp(127.0.0.1:3306)/"+dbName)
    if err != nil {
        panic(err.Error())
    }
    return db
}

func Requests() {
    // creates a new instance of a mux router
    myRouter := mux.NewRouter().StrictSlash(true)
    // replace http.HandleFunc with myRouter.HandleFunc
    myRouter.HandleFunc("/", home)
	myRouter.HandleFunc("/all", AllUsers)
	myRouter.HandleFunc("/user/{id}", returnSingleUser)
    myRouter.HandleFunc("/newUser", createNewUser).Methods("POST")
    myRouter.HandleFunc("/newSingleUser", createSingleNewUser).Methods("POST")
    
    // finally, instead of passing in nil, we want
    // to pass in our newly created router as the second
    // argument
    log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func main() {
    fmt.Println("Rest API v2.0 - Mux Routers")
    
    Requests()
}

func AllUsers(w http.ResponseWriter, r *http.Request){
    fmt.Println("Endpoint Hit: returnAllArticles")
    db := dbConn()
    selDB, err := db.Query("Select id,name,city from users WHERE id>0 Order by id ASC")
    if err != nil {
        panic(err.Error())
    }
    emp := user{}
    res := []user{}
    for selDB.Next() {
        var id,name, city string
        err = selDB.Scan(&id, &name, &city)
        if err != nil {
            panic(err.Error())
        }
        emp.Id = id
        emp.Name = name
        emp.City = city
        res = append(res, emp)
    }
   // tmpl.ExecuteTemplate(w, "Index", res)
    defer db.Close()
    json.NewEncoder(w).Encode(res)
}

func returnSingleUser(w http.ResponseWriter, r *http.Request){
    vars := mux.Vars(r)
    key := vars["id"]

    db := dbConn()
    selDB, err := db.Query("Select id,name,city from users WHERE id=" + key + " Order by id ASC")
    if err != nil {
        panic(err.Error())
    }
    emp := user{}
    res := []user{}
    for selDB.Next() {
        var id,name, city string
        err = selDB.Scan(&id, &name, &city)
        if err != nil {
            panic(err.Error())
        }
        emp.Id = id
        emp.Name = name
        emp.City = city
        res = append(res, emp)
    }
   // tmpl.ExecuteTemplate(w, "Index", res)
    defer db.Close()
    json.NewEncoder(w).Encode(res)
}

func home(w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w, "Welcome to the HomePage!")
    fmt.Println("Endpoint Hit: homePage")
}

func createNewUser(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Startpoint Hit: NewUser")
    w.Header().Add("content-type","application/json")
    var newUsers users
    
    reqBody, _ := ioutil.ReadAll(r.Body)
    //newUser := user{}
    //json.NewDecoder(r.Body).Decode(&newUsers)

   json.Unmarshal(reqBody, &newUsers)

   
    db := dbConn()

    for i := 0; i < len(newUsers.users); i++ {
        _, err := db.Query("INSERT INTO users(id,name,city) VALUES('" + newUsers.users[i].Id + "','" + newUsers.users[i].Name + "','" + newUsers.users[i].City + "');")
        if err != nil {
            panic(err.Error())
        }
     //   fmt.Println("INSERT INTO users(id,name,city) VALUES('" + newUsers.users[i].Id + "','" + newUsers.users[i].Name + "','" + newUsers.users[i].City + "');")
    }

    //fmt.Println(newUsers.users[0].Id)

    //fmt.Println(newUsers.users[0].Name)

   // fmt.Println(newUsers.users[0].City)

   fmt.Println(newUsers)

    fmt.Println(string(reqBody))

    fmt.Println(r.Body)

    //fmt.Println(len(newUsers.users))

    fmt.Println("Endpoint Hit: NewUser")

    json.NewEncoder(w).Encode("Not Working")

    defer db.Close()
   
}


func createSingleNewUser(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Startpoint Hit: NewUser")
    var newUsers user
    reqBody, _ := ioutil.ReadAll(r.Body)
    json.Unmarshal(reqBody, &newUsers)
    db := dbConn()

    _, err := db.Query("INSERT INTO users(id,name,city) VALUES('" + newUsers.Id + "','" + newUsers.Name + "','" + newUsers.City + "');")
    if err != nil {
        panic(err.Error())
    }
   fmt.Println(newUsers)

    fmt.Println(string(reqBody))

    fmt.Println(r.Body)

    //fmt.Println(len(newUsers.users))

    fmt.Println("Endpoint Hit: NewUser")

    defer db.Close()
   
}