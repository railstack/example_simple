package controllers

import (
	"fmt"
	"log"
	"net/http"

	m "../models"
	"gopkg.in/gin-gonic/gin.v1"
)

// GET /articles/1/comments
func CommentsIndex(c *gin.Context) {
	id, err := ToInt(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, BuildResp("400", "Parsing id error!", nil))
		return
	}
	Comments, err := m.FindCommentsBy("article_id", id)
	if err != nil {
		msg := fmt.Sprintf("Get Comment index error: %v", err)
		c.JSON(http.StatusOK, BuildResp("400", msg, nil))
		return
	}
	resp := BuildResp("200", "Get Comment index success", Comments)
	c.JSON(http.StatusOK, resp)
}

// GET /comments/1
func CommentsShow(c *gin.Context) {
	id, err := ToInt(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, BuildResp("400", "Parsing id error!", nil))
		return
	}
	Comment, err := m.FindComment(id)
	if err != nil {
		msg := fmt.Sprintf("Get Comment error: %v", err)
		c.JSON(http.StatusOK, BuildResp("400", msg, nil))
		return
	}
	resp := BuildResp("200", "Get Comment success", Comment)
	c.JSON(http.StatusOK, resp)
}

func CommentsNew(c *gin.Context) {
}

func CommentsEdit(c *gin.Context) {
}

// POST /comments
func CommentsCreate(c *gin.Context) {
	var ar m.Comment
	if c.BindJSON(&ar) != nil {
		c.JSON(http.StatusOK, BuildResp("400", "Parsing id error!", nil))
		return
	}
	id, err := ar.Create()
	if err != nil {
		msg := fmt.Sprintf("Create Comment error: %v", err)
		log.Println(msg)
		c.JSON(http.StatusOK, BuildResp("400", msg, nil))
		return
	}
	c.JSON(http.StatusOK, BuildResp("200", "Create Comment success", map[string]int64{"id": id}))
}

// PUT /comments/1
func CommentsUpdate(c *gin.Context) {
	id, err := ToInt(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, BuildResp("400", "Parsing id error!", nil))
		return
	}
	ar, err := m.FindComment(id)
	if err != nil {
		msg := fmt.Sprintf("Update Comment error: %v", err)
		log.Println(msg)
		c.JSON(http.StatusOK, BuildResp("400", msg, nil))
	}
	am := map[string]interface{}{}
	var json m.Comment
	if c.BindJSON(&json) == nil {
		if json.Commenter != "" {
			am["commenter"] = json.Commenter
		}
		if json.Body != "" {
			am["body"] = json.Body
		}
		if json.ArticleId != 0 {
			am["article_id"] = json.ArticleId
		}
	}
	err = ar.Update(am)
	if err != nil {
		msg := fmt.Sprintf("Update Comment error: %v", err)
		log.Println(msg)
		c.JSON(http.StatusOK, BuildResp("400", msg, nil))
		return
	}
	resp := BuildResp("200", "Update Comment success", nil)
	c.JSON(http.StatusOK, resp)
}

// DELETE /comments/1
func CommentsDestroy(c *gin.Context) {
	id, err := ToInt(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, BuildResp("400", "Params error!", nil))
		return
	}
	err = m.DestroyComment(id)
	if err != nil {
		fmt.Println(err)
	}
	resp := BuildResp("200", "Comment destroied", nil)
	c.JSON(http.StatusOK, resp)
}
