package controllers

import (
	"log"
	"net/http"
	"time"

	"cassandra_go_recommendation_api/db"
	"cassandra_go_recommendation_api/models"

	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
)

func LogInteraction(c *gin.Context) {
    var interaction models.UserInteraction
    if err := c.ShouldBindJSON(&interaction); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    timestamp := time.Now()

    var category string
    if err := db.Session.Query(`
        SELECT category FROM items WHERE item_id = ?`,
        interaction.ItemID).Scan(&category); err != nil {
        log.Println("Erro ao obter categoria do item:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve item category"})
        return
    }

    batch := db.Session.NewBatch(gocql.LoggedBatch)

    batch.Query(`
        INSERT INTO user_interactions (user_id, timestamp, item_id, category, action_type)
        VALUES (?, ?, ?, ?, ?)`,
        interaction.UserID, timestamp, interaction.ItemID, category, interaction.ActionType)

    batch.Query(`
        INSERT INTO user_interacted_items (user_id, item_id)
        VALUES (?, ?)`,
        interaction.UserID, interaction.ItemID)

    if err := db.Session.ExecuteBatch(batch); err != nil {
        log.Println("Erro ao registrar interação:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to log interaction"})
        return
    }

    if err := db.Session.Query(`
        UPDATE user_category_interactions
        SET interaction_count = interaction_count + 1
        WHERE user_id = ? AND category = ?`,
        interaction.UserID, category).Exec(); err != nil {
        log.Println("Erro ao atualizar contador em user_category_interactions:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user category interactions"})
        return
    }

    if err := db.Session.Query(`
        UPDATE item_popularity_by_category
        SET popularity_counter = popularity_counter + 1
        WHERE category = ? AND item_id = ?`,
        category, interaction.ItemID).Exec(); err != nil {
        log.Println("Erro ao atualizar contador em item_popularity_by_category:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update item popularity"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Interaction logged successfully"})
}

