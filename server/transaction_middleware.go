package server

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func TransactionMiddleware(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tx := db.MustBegin()
		defer tx.Rollback()

		c.Set("tx", tx)

		c.Next()

		if _, ok := c.Get("txCommitted"); !ok {
			tx.Rollback()
		}
	}
}
