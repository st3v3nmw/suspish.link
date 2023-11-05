package main

import (
	"math/rand"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetHttpScheme(c *gin.Context) string {
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	return scheme
}

func GenerateRandomString(n int) string {
	alphabet := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	var builder strings.Builder
	for i := 0; i < n; i++ {
		builder.WriteString(string(alphabet[rand.Int()%62]))
	}
	return builder.String()
}
