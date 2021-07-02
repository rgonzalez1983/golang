package middleware

import (
	//"github.com/go-chi/chi"
	//httpSwagger "github.com/swaggo/http-swagger"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go_project/db"
	"go_project/internal/entity"
	"go_project/internal/persistance"
	//"go_project/internal/service"
	"encoding/json"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"net/http"
	"os"
	"time"
)

type App struct {
	Router           *mux.Router
	DB               *db.MongoConnection
	Logg             log.Logger
	PersonRepository persistance.PersonRepository
}

func (a *App) Run(addr string) error {
	err := http.ListenAndServe(addr, a.Router)
	return err
}

// CORS Middleware
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set headers
		w.Header().Set("Access-Control-Allow-Headers:", "*")
		w.Header().Set("Access-Control-Allow-Origin", os.Getenv("ORIGIN_ALLOWED"))
		w.Header().Set("Access-Control-Allow-Methods", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
		return
	})
}

func (a *App) Initialize(_user, _password string) (err error) {
	fmt.Println("Starting the application....")
	host := os.Getenv("MONGO_HOST") + ":27017"
	dbs := os.Getenv("MONGO_DATABASE")
	info := &mgo.DialInfo{
		Addrs:    []string{host},
		Timeout:  60 * time.Hour,
		Database: dbs,
		Username: _user,
		Password: _password,
	}
	a.DB, _ = db.NewConnection(info)
	fmt.Println("Connected to MongoDB!")
	muxObj := mux.NewRouter()
	muxObj.Use(CORS)
	a.Router = muxObj
	a.initializeLogger()
	a.initializeRoutes()
	a.initializeRepository()
	return err
}

// routing
func (a *App) initializeRoutes() {
	//a.Router.HandleFunc("/api/doc/index.html", httpSwagger.WrapHandler).Methods("GET")
	//a.Router.PathPrefix("/api").Handler(httpSwagger.WrapHandler)
	a.Router.HandleFunc("/index", a.getIndex).Methods("GET")
	a.Router.HandleFunc("/create_person", a.CreatePerson).Methods("POST")
	a.Router.HandleFunc("/update_person", a.UpdatePerson).Methods("POST")
	a.Router.HandleFunc("/get_person", a.GetPerson).Methods("POST")
	a.Router.HandleFunc("/delete_person", a.DeletePerson).Methods("POST")
	a.Router.HandleFunc("/list_persons", a.ListPersons).Methods("POST")
}

// Logger
func (a *App) initializeLogger() {
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}
	a.Logg = logger
}

// Repository
func (a *App) initializeRepository() {
	a.PersonRepository = persistance.NewPersonRepository(a.DB)
}

// Inicio
func (a *App) getIndex(w http.ResponseWriter, r *http.Request) {

	item := &entity.JsonResponse{
		Message:    "API MS-Golang Basic CRUD",
		StatusCode: http.StatusOK,
	}
	respondWithJSON(w, http.StatusOK, item)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = w.Write(response)
}
