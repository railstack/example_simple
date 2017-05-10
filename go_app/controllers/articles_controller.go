package controller

import (
	"fmt"
	"log"
	"net/http"

	m "../models"
	"gopkg.in/gin-gonic/gin.v1"
)

func ArticlesIndex(c *gin.Context) {
	articles, err := m.AllArticles()
	if err != nil {
		msg := fmt.Sprintf("Get article index error: %v", err)
		c.JSON(http.StatusOK, BuildResp("400", msg, nil))
		return
	}
	resp := BuildResp("200", "Get article index success", articles)
	c.JSON(http.StatusOK, resp)
}

func ArticlesShow(c *gin.Context) {
	id, err := ToInt(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, BuildResp("400", "Parsing id error!", nil))
		return
	}
	article, err := m.FindArticle(id)
	if err != nil {
		msg := fmt.Sprintf("Get article error: %v", err)
		c.JSON(http.StatusOK, BuildResp("400", msg, nil))
		return
	}
	resp := BuildResp("200", "Get article success", article)
	c.JSON(http.StatusOK, resp)
}

func ArticlesNew(c *gin.Context) {
}

func ArticlesEdit(c *gin.Context) {
}

func ArticlesCreate(c *gin.Context) {
	var ar m.Article
	if c.BindJSON(&ar) != nil {
		c.JSON(http.StatusOK, BuildResp("400", "Parsing id error!", nil))
		return
	}
	err := ar.Create()
	if err != nil {
		msg := fmt.Sprintf("Create article error: %v", err)
		log.Println(msg)
		c.JSON(http.StatusOK, BuildResp("400", msg, nil))
		return
	}
	c.JSON(http.StatusOK, BuildResp("200", "Create article success", nil))
}

func ArticlesUpdate(c *gin.Context) {
	id, err := ToInt(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, BuildResp("400", "Parsing id error!", nil))
		return
	}
	ar, err := m.FindArticle(id)
	if err != nil {
		msg := fmt.Sprintf("Update article error: %v", err)
		log.Println(msg)
		c.JSON(http.StatusOK, BuildResp("400", msg, nil))
	}
	am := map[string]interface{}{}
	var json m.Article
	if c.BindJSON(&json) == nil {
		if json.Title != "" {
			am["title"] = json.Title
		}
		if json.Text != "" {
			am["text"] = json.Text
		}
	}
	err = ar.Update(am)
	if err != nil {
		msg := fmt.Sprintf("Update article error: %v", err)
		log.Println(msg)
		c.JSON(http.StatusOK, BuildResp("400", msg, nil))
		return
	}
	resp := BuildResp("200", "Update article success", nil)
	c.JSON(http.StatusOK, resp)
}

func ArticlesDestroy(c *gin.Context) {
	id, err := ToInt(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, BuildResp("400", "Params error!", nil))
		return
	}
	err = m.DestroyArticle(id)
	if err != nil {
		fmt.Println(err)
	}
	resp := BuildResp("200", "Article destroied", nil)
	c.JSON(http.StatusOK, resp)
}
