package controllers

import (
    "log"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "cassandra_go_recommendation_api/db"
    "cassandra_go_recommendation_api/models"
)

func LogInteraction(c *gin.Context) {
    var interaction models.UserInteraction
    if err := c.ShouldBindJSON(&interaction); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Registrar a interação do usuário
    query := `INSERT INTO user_interactions (user_id, item_id, action_type, category, timestamp)
              VALUES (?, ?, ?, ?, ?)`
    if err := db.Session.Query(query, interaction.UserID, interaction.ItemID, interaction.ActionType, interaction.Category, time.Now()).Exec(); err != nil {
        log.Println("Erro ao registrar interação:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to log interaction"})
        return
    }

    // Atualizar o contador de popularidade na tabela de popularidade por item
    updateItemPopularityQuery := `UPDATE item_popularity SET popularity_counter = popularity_counter + 1 WHERE item_id = ?`
    if err := db.Session.Query(updateItemPopularityQuery, interaction.ItemID).Exec(); err != nil {
        log.Println("Erro ao atualizar popularidade do item:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update item popularity"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Interaction logged successfully"})
}
