package controllers

import (
	"log"
	"net/http"
	"time"

	"cassandra_go_recommendation_api/db"

	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
)

func GetRecommendations(c *gin.Context) {
    userID := c.Param("user_id")

    // Buscar as categorias mais frequentes nas interações do usuário
    queryInteractions := `SELECT category, COUNT(*) as interactions_count
                          FROM user_interactions WHERE user_id = ? AND timestamp > ? 
                          GROUP BY category ORDER BY interactions_count DESC LIMIT 3`
    iter := db.Session.Query(queryInteractions, userID, time.Now().AddDate(0, -1, 0)).Iter() // Exemplo: Interações no último mês

    var category string
    var interactionsCount int
    userCategories := []string{}

    for iter.Scan(&category, &interactionsCount) {
        userCategories = append(userCategories, category)
    }

    if err := iter.Close(); err != nil {
        log.Println("Erro ao buscar interações do usuário:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user interactions"})
        return
    }

    if len(userCategories) == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "No interactions found for user"})
        return
    }

    // Para cada categoria, buscar os itens mais populares
    recommendations := []map[string]interface{}{}

    for _, category := range userCategories {
        queryItems := `SELECT item_id, name FROM items WHERE category = ? LIMIT 5`
        iterItems := db.Session.Query(queryItems, category).Iter()

        var itemID gocql.UUID
        var name string

        for iterItems.Scan(&itemID, &name) {
            recommendations = append(recommendations, map[string]interface{}{
                "category": category,
                "item_id":  itemID,
                "name":     name,
            })
        }

        if err := iterItems.Close(); err != nil {
            log.Println("Erro ao buscar itens populares por categoria:", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get popular items"})
            return
        }
    }

    if len(recommendations) == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "No recommendations available for user"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"recommendations": recommendations})
}

func GetPopularItems(c *gin.Context) {
    query := `SELECT item_id, popularity_counter FROM item_popularity ORDER BY popularity_counter DESC LIMIT 5`
    iter := db.Session.Query(query).Iter()

    var itemID gocql.UUID
    var popularityCounter int64
    popularItems := []map[string]interface{}{}

    for iter.Scan(&itemID, &popularityCounter) {
        popularItems = append(popularItems, map[string]interface{}{
            "item_id":          itemID,
            "popularity_counter": popularityCounter,
        })
    }

    if err := iter.Close(); err != nil {
        log.Println("Erro ao obter popularidade dos itens:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get popular items"})
        return
    }

    if len(popularItems) == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "No popular items found"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"popular_items": popularItems})
}
