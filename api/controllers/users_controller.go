package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Grimm7301/apigo/api/auth"
	"github.com/Grimm7301/apigo/api/models"
	"github.com/Grimm7301/apigo/api/utils/formaterror"
	"github.com/gin-gonic/gin"
)

func (server *Server) CreateUser(c *gin.Context) {
	newUser := models.User{}
	if err := c.BindJSON(&newUser); err != nil {
		c.IndentedJSON(http.StatusUnprocessableEntity, gin.H{"Error": err})
		return
	}
	newUser.Prepare()
	err := newUser.Validate("")
	if err != nil {
		c.IndentedJSON(http.StatusUnprocessableEntity, gin.H{"Error": err})
		return
	}
	userCreated, err := newUser.CreateUser(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"Error": formattedError})
		return
	}
	c.IndentedJSON(http.StatusCreated, userCreated)
}

func (server *Server) GetUsers(c *gin.Context) {

	user := models.User{}

	users, err := user.FindAllUsers(server.DB)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"Error": err})
		return
	}
	c.IndentedJSON(http.StatusOK, users)
}

func (server *Server) GetUser(c *gin.Context) {

	uid, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"Error": err})
		return
	}
	user := models.User{}
	userGotten, err := user.FindUserByID(server.DB, uint32(uid))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"Error": err})
		return
	}
	c.IndentedJSON(http.StatusOK, userGotten)
}

func (server *Server) UpdateUser(c *gin.Context) {

	uid, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"Error": err})
		return
	}
	user := models.User{}
	if err := c.BindJSON(&user); err != nil {
		c.IndentedJSON(http.StatusUnprocessableEntity, gin.H{"Error": err})
		return
	}

	tokenID, err := auth.ExtractTokenID(c.Request)
	if err != nil {
		c.IndentedJSON(http.StatusUnprocessableEntity, gin.H{"Error": errors.New("Unauthorized")})
		return
	}
	if tokenID != uint32(uid) {
		c.IndentedJSON(http.StatusUnprocessableEntity, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	user.Prepare()
	err = user.Validate("update")
	if err != nil {
		c.IndentedJSON(http.StatusUnprocessableEntity, gin.H{"Error": err})
		return
	}
	updatedUser, err := user.UpdateUser(server.DB, uint32(uid))
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, formattedError)
		return
	}
	c.IndentedJSON(http.StatusOK, updatedUser)
}

func (server *Server) DeleteUser(c *gin.Context) {

	uid, err := strconv.ParseUint(c.Param("id"), 10, 32)
	user := models.User{}
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"Error": err})
		return
	}
	tokenID, err := auth.ExtractTokenID(c.Request)
	if err != nil {
		c.IndentedJSON(http.StatusUnprocessableEntity, gin.H{"Error": errors.New("Unauthorized")})
		return
	}
	if tokenID != uint32(uid) {
		c.IndentedJSON(http.StatusUnprocessableEntity, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	_, err = user.DeleteUser(server.DB, uint32(uid))
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
		return
	}
	c.IndentedJSON(http.StatusNoContent, gin.H{"DeletetedUserID": uid})
}
