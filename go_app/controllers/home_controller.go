package controller

import (
	"net/http"

	"gopkg.in/gin-gonic/gin.v1"
	// you can import models
	//m "../models"
)

func HomeHandler(c *gin.Context) {
	// you can use model functions to do CRUD
	//
	// user, _ := m.FindUser(1)
	// u, err := json.Marshal(user)
	// if err != nil {
	// 	log.Printf("JSON encoding error: %v\n", err)
	// 	u = []byte("Get data error!")
	// }

	type Envs struct {
		GoOnRailsVer string
		GolangVer    string
	}

	gorVer := "0.1.4"
	golangVer := "go version go1.7.4 darwin/amd64"

	envs := Envs{GoOnRailsVer: gorVer, GolangVer: golangVer}
	c.HTML(http.StatusOK, "index.tmpl", envs)
}
