package controllers

import (
    "log"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/gocql/gocql"
    "cassandra_go_recommendation_api/db"
    "cassandra_go_recommendation_api/models"
)

func CreateItem(c *gin.Context) {
    var item models.Item

    if err := c.ShouldBindJSON(&item); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido fornecido"})
        return
    }

    itemID, err := gocql.RandomUUID()
    if err != nil {
        log.Println("Falha ao gerar UUID:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao criar item"})
        return
    }

    item.ItemID = itemID

    batch := db.Session.NewBatch(gocql.LoggedBatch)

    batch.Query(`
        INSERT INTO items (item_id, name, category) VALUES (?, ?, ?)`,
        item.ItemID, item.Name, item.Category)

    batch.Query(`
        INSERT INTO items_by_category (category, item_id, name) VALUES (?, ?, ?)`,
        item.Category, item.ItemID, item.Name)

    if err := db.Session.ExecuteBatch(batch); err != nil {
        log.Println("Falha ao inserir item no Cassandra:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao criar item"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "message": "Item criado com sucesso",
        "item_id": item.ItemID,
    })
}
