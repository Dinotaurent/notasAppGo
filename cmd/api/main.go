package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"notas-app-go/internal/api"
	"notas-app-go/internal/data"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {

	// Cargar .env
	//err := godotenv.Load()
	//if err != nil {
	//	slog.Error("Error loading .env file", "error", err)
	//	return
	//}

	_ = godotenv.Load()
	webPort := os.Getenv("PORT")

	// Conexion con la db
	var uri string
	if os.Getenv("MONGODB_URL") == "" {
		mongoHost := os.Getenv("MONGO_HOST")
		mongoUser := os.Getenv("MONGO_USER")
		mongoPassword := os.Getenv("MONGO_PASSWORD")
		mongoDb := os.Getenv("MONGO_DB")
		mongoPort := os.Getenv("MONGO_PORT")
		uri = fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?authSource=admin", mongoUser, mongoPassword, mongoHost, mongoPort, mongoDb)
	} else {
		uri = os.Getenv("MONGODB_URL")
	}

	//log.Printf("CONTENIDO DE LA URI: %s", uri)

	client, err := mongo.Connect(options.Client().ApplyURI(uri))

	// Verificar problemas con la cadena de conexion
	if err != nil {
		//log.Printf("CONTENIDO DE LA URI: %s", uri)
		slog.Error("Error en la cadena de conexion con al db", "Error:", err.Error())
		return
	}

	// Timeout para la conexion con la db
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	// Verificiar la conexion
	err = client.Ping(ctx, nil)
	if err != nil {
		//log.Printf("CONTENIDO DE LA URI: %s", uri)
		slog.Error("Error de conexion con la db, seguro que esta arriba ?", "Error:", err.Error())
		return
	}
	log.Println("Conectado con MongoDB!")
	log.Printf("Servicio iniciado en el puerto: %s\n", webPort)

	app := &api.Application{
		Notas: data.NotaModel{Client: client},
	}

	// Definir el servidor http
	srv := http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.Routes(),
	}

	// Iniciar el servidor
	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

}
