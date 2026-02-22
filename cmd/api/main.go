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
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	webPort := os.Getenv("PORT")

	// Conexion con la db
	uri := os.Getenv("MONGODB_URL")
	client, err := mongo.Connect(options.Client().ApplyURI(uri))

	// Verificar problemas con la cadena de conexion
	if err != nil {
		slog.Error("Error en la cadena de conexion con al db", "Error:", err.Error())
		return
	}

	// Timeout para la conexion con la db
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	// Verificiar la conexion
	err = client.Ping(ctx, nil)
	if err != nil {
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
