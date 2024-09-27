package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/steinfletcher/apitest"
	"jordi.tfg.rewrite/internal/data"
)

func TestGetUserInfo(t *testing.T) {
	// Crear una instancia de Application
	db, err := OpenDB()
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}

	app := &Application{
		Models: data.NewModels(db),
		Logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}

	// Crea un router que utilice la función GetUserInfo
	router := httprouter.New()
	router.GET("/user/:id", app.GetUserInfo)

	tests := []struct {
		id       int
		wantcode int
	}{
		// {id: 8, wantcode: http.StatusOK},                   // Suponiendo que el usuario con ID 8 existe
		{id: -1, wantcode: http.StatusUnprocessableEntity}, // ID inválido
	}

	// Realizar la prueba con apitest
	for _, tt := range tests {
		url := fmt.Sprintf("/user/%d", tt.id) // Cambiado %q por %d
		apitest.New().
			Handler(router).
			Get(url). // Simular una petición GET con un ID de ejemplo
			Expect(t).
			Status(tt.wantcode). // Esperar el código de estado deseado
			End()
	}
}
