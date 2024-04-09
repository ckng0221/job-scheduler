package controllers

import (
	"context"
	"errors"
	"fmt"
	"job-scheduler/api/config"
	"job-scheduler/api/initializers"
	"job-scheduler/api/models"
	"job-scheduler/utils"
	"net/http"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	// User login based on access token
	var body struct {
		Code  string
		State string
		Nonce string
	}

	cookieState, err := c.Cookie("state")
	if err != nil {
		c.AbortWithStatusJSON(401, "State not found")
		return
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	if body.Code == "" {
		c.AbortWithStatusJSON(401, "No access token")
		return
	}

	tokenClaims, idToken, err := getTokenClaimJwtFromLogin(body.Code, body.State, cookieState, body.Nonce)
	if err != nil {
		c.AbortWithError(400, err)
		return
	}

	// register user if not exist
	var user models.User
	err = initializers.Db.Where("sub = ?", tokenClaims.Sub).Limit(1).Find(&user).Error
	fmt.Println(err)
	if user.ID == 0 {
		fmt.Println("user not found")
		initializers.Db.Create(&models.User{
			Name:       tokenClaims.Name,
			Email:      tokenClaims.Email,
			Sub:        tokenClaims.Sub,
			ProfilePic: tokenClaims.Picture,
		})
		fmt.Println("user created")
	}

	// Set cookies
	// c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", idToken, 3600*24*30, "", "", false, false)

	// respond
	c.JSON(http.StatusOK, gin.H{
		"name":         tokenClaims.Name,
		"access_token": idToken,
	})
}

func GoogleLogin(c *gin.Context) {

	state, err := utils.RandString(16)
	if err != nil {
		c.AbortWithStatus(500)
		return
	}
	nonce, err := utils.RandString(16)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(500)
		return
	}
	// c.SetCookie("state", state, 3600*24*30, "", "", false, false)
	// c.SetCookie("nonce", nonce, 3600*24*30, "", "", false, false)

	config, err := config.GoogleConfig()
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(500)
		return
	}

	url := config.AuthCodeURL(state, oidc.Nonce(nonce))

	c.JSON(http.StatusOK, gin.H{
		"state": state,
		"nonce": nonce,
		"url":   url,
	})
}

func getTokenClaimJwtFromLogin(code, state, cookieState, nonce string) (config.IDTokenClaims, string, error) {
	if state != cookieState {
		return config.IDTokenClaims{}, "", errors.New("state does not match")
	}
	ctx := context.Background()
	authConfig, _ := config.GoogleConfig()
	verifier := config.GetVerifier()

	oauth2Token, err := authConfig.Exchange(ctx, code)
	if err != nil {
		return config.IDTokenClaims{}, "", err
	}

	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		fmt.Println("No id_token field in oauth2 token.")
		return config.IDTokenClaims{}, "", err
	}

	// JWT token from identify provider
	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		fmt.Println("Failed to verify id token", err)
		return config.IDTokenClaims{}, "", err
	}

	if idToken.Nonce != nonce {
		return config.IDTokenClaims{}, "", errors.New("nonce does not match")
	}

	var tokenClaims config.IDTokenClaims
	if err := idToken.Claims(&tokenClaims); err != nil {
		// handle error
		fmt.Println("Failed to unmarshal claim")
		return config.IDTokenClaims{}, "", err
	}

	return tokenClaims, rawIDToken, nil
}

func Validate(c *gin.Context) {
	user, _ := c.Get("user")

	c.JSON(http.StatusOK, user)
}

func Logout(c *gin.Context) {
	// Set cookies
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", "", -1, "", "", false, false)

	// respond
	c.JSON(http.StatusOK, gin.H{
		"message": "logged out",
	})
}

func requireOwner(c *gin.Context, ownerId uint) error {
	requestUser, _ := c.Get("user")
	if requestUser.(models.User).ID != ownerId && requestUser.(models.User).Role != "admin" {
		return errors.New("forbidden")
	}
	return nil
}

func requireAdmin(c *gin.Context) error {
	requestUser, _ := c.Get("user")
	if requestUser.(models.User).Role != "admin" {
		return errors.New("forbidden")
	}
	return nil
}
