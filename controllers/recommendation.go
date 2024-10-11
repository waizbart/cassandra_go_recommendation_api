package controllers

import (
	"fmt"
	"log"
	"net/http"
	"sort"

	"cassandra_go_recommendation_api/db"
	"cassandra_go_recommendation_api/models"

	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
)

func GetRecommendations(c *gin.Context) {
    userIDStr := c.Param("user_id")
    userID, err := gocql.ParseUUID(userIDStr); 
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
        return
    }

    iter := db.Session.Query(`
        SELECT category FROM user_category_interactions
        WHERE user_id = ?`, userID).Iter()

    var categories []string
    var category string
    for iter.Scan(&category) {
        categories = append(categories, category)
    }

    fmt.Println("categories", categories)

    if err := iter.Close(); err != nil {
        log.Println("Failed to fetch user categories:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user data"})
        return
    }

    if len(categories) == 0 {
        c.JSON(http.StatusOK, gin.H{"recommendations": []models.Item{}})
        return
    }

    type ItemScore struct {
        Item   models.Item
        Score  int64
    }

    var recommendations []ItemScore
    interactedItems := make(map[gocql.UUID]struct{})

    userItemsIter := db.Session.Query(`
        SELECT item_id FROM user_interacted_items
        WHERE user_id = ?`, userID).Iter()

    var interactedItemID gocql.UUID
    for userItemsIter.Scan(&interactedItemID) {
        interactedItems[interactedItemID] = struct{}{}
    }

    fmt.Println("interacted items", interactedItems)

    if err := userItemsIter.Close(); err != nil {
        log.Println("Failed to fetch user interactions:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user interactions"})
        return
    }

    for _, category := range categories {
        fmt.Println(category)
        itemsIter := db.Session.Query(`
            SELECT item_id, popularity_counter FROM item_popularity_by_category
            WHERE category = ?`, category).Iter()

        var itemID gocql.UUID
        var popularityCounter int64

        for itemsIter.Scan(&itemID, &popularityCounter) {
            fmt.Println(itemID, popularityCounter)
            if _, exists := interactedItems[itemID]; exists {
                continue
            }

            var item models.Item
            if err := db.Session.Query(`
                SELECT item_id, name, category FROM items WHERE item_id = ?`,
                itemID).Scan(&item.ItemID, &item.Name, &item.Category); err != nil {
                log.Println("Failed to fetch item details:", err)
                continue
            }

            recommendations = append(recommendations, ItemScore{
                Item:  item,
                Score: popularityCounter,
            })
        }

        if err := itemsIter.Close(); err != nil {
            log.Println("Failed to fetch items for category:", category, err)
            continue
        }
    }

    fmt.Println("recomendations", recommendations)

    // 3. Ordenar recomendações por score (popularidade) e limitar a N itens
    sort.Slice(recommendations, func(i, j int) bool {
        return recommendations[i].Score > recommendations[j].Score
    })

    // Limitar a 5 recomendações
    maxRecommendations := 5
    if len(recommendations) < maxRecommendations {
        maxRecommendations = len(recommendations)
    }

    finalRecommendations := make([]models.Item, maxRecommendations)
    for i := 0; i < maxRecommendations; i++ {
        finalRecommendations[i] = recommendations[i].Item
    }

    c.JSON(http.StatusOK, gin.H{"recommendations": finalRecommendations})
}
