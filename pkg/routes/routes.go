package routes

import (
	"encoding/json"
	"io"
	"log"

	"github.com/adaggerboy/genesis-academy-case-app/pkg/3rd/openexchangeapi"
	"github.com/adaggerboy/genesis-academy-case-app/pkg/database"
	"github.com/gin-gonic/gin"
)

type StructResponseJSON struct {
	Rate float32 `json:"rate"`
}

type StructSubscriptionRequest struct {
	Email string `json:"email"`
}

func DeployRoutes(r gin.IRouter) {
	r.GET("rate", func(ctx *gin.Context) {
		rate, err := openexchangeapi.RequestUSDPairCached("UAH")
		if err != nil {
			log.Printf("Request /rate: error: %s", err)
			ctx.String(500, "internal server error")
			return
		}
		ctx.JSON(200, StructResponseJSON{
			Rate: rate,
		})
	})
	r.POST("subscribe", func(ctx *gin.Context) {
		sub := StructSubscriptionRequest{}
		bytes, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			log.Printf("Request /subscribe: error: %s", err)
			ctx.String(400, "bad request")
			return
		}
		json.Unmarshal(bytes, &sub)
		if err != nil {
			log.Printf("Request /subscribe: error: %s", err)
			ctx.String(400, "bad request")
			return
		}
		created, err := database.GetDatabase().CreateSubscription(sub.Email)
		if err != nil {
			log.Printf("Request /subscribe: error: %s", err)
			ctx.String(500, "internal server error")
			return
		}
		if created {
			ctx.Status(201)
		} else {
			ctx.Status(409)
		}

	})
}
