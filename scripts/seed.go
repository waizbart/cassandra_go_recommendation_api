package main

import (
    "cassandra_go_recommendation_api/db"
    "cassandra_go_recommendation_api/models"
    "fmt"
    "log"
    "math/rand"
    "time"

    "github.com/gocql/gocql"
)

func main() {
    session := db.InitCassandra()
    defer db.CloseCassandra()

    numUsers := 100
    numItems := 1000
    numInteractions := 10000

    categories := []string{"Eletrônicos", "Livros", "Roupas", "Jogos", "Móveis", "Filmes", "Séries"}
    itemIDs := make([]gocql.UUID, numItems)

    for i := 0; i < numItems; i++ {
        itemID, err := gocql.RandomUUID()
        if err != nil {
            log.Fatal("Falha ao gerar UUID para item:", err)
        }

        item := models.Item{
            ItemID:   itemID,
            Name:     fmt.Sprintf("Produto %d", i+1),
            Category: categories[rand.Intn(len(categories))],
        }

        batch := session.NewBatch(gocql.LoggedBatch)
        batch.Query(`
            INSERT INTO items (item_id, name, category) VALUES (?, ?, ?)`,
            item.ItemID, item.Name, item.Category)
        batch.Query(`
            INSERT INTO items_by_category (category, item_id, name) VALUES (?, ?, ?)`,
            item.Category, item.ItemID, item.Name)

        if err := session.ExecuteBatch(batch); err != nil {
            log.Println("Falha ao inserir item no Cassandra:", err)
        }

        itemIDs[i] = itemID
    }

    fmt.Printf("Inseridos %d itens.\n", numItems)

    userIDs := make([]gocql.UUID, numUsers)
    for i := 0; i < numUsers; i++ {
        userID, err := gocql.RandomUUID()
        if err != nil {
            log.Fatal("Falha ao gerar UUID para usuário:", err)
        }
        userIDs[i] = userID
    }

    fmt.Printf("Gerados %d usuários.\n", numUsers)

    for i := 0; i < numInteractions; i++ {
        userID := userIDs[rand.Intn(len(userIDs))]
        itemID := itemIDs[rand.Intn(len(itemIDs))]
        actionTypes := []string{"view", "purchase", "add_to_cart", "click"}
        actionType := actionTypes[rand.Intn(len(actionTypes))]
        timestamp := time.Now().Add(time.Duration(-rand.Intn(100000)) * time.Second)

        var category string
        if err := session.Query(`
            SELECT category FROM items WHERE item_id = ?`,
            itemID).Scan(&category); err != nil {
            log.Println("Erro ao obter categoria do item:", err)
            continue
        }

        batch := session.NewBatch(gocql.LoggedBatch)
        batch.Query(`
            INSERT INTO user_interactions (user_id, timestamp, item_id, category, action_type)
            VALUES (?, ?, ?, ?, ?)`,
            userID, timestamp, itemID, category, actionType)

        batch.Query(`
            INSERT INTO user_interacted_items (user_id, item_id)
            VALUES (?, ?)`,
            userID, itemID)

        if err := session.ExecuteBatch(batch); err != nil {
            log.Println("Erro ao registrar interação:", err)
            continue
        }

        if err := session.Query(`
            UPDATE user_category_interactions
            SET interaction_count = interaction_count + 1
            WHERE user_id = ? AND category = ?`,
            userID, category).Exec(); err != nil {
            log.Println("Erro ao atualizar contador em user_category_interactions:", err)
            continue
        }

        if err := session.Query(`
            UPDATE item_popularity_by_category
            SET popularity_counter = popularity_counter + 1
            WHERE category = ? AND item_id = ?`,
            category, itemID).Exec(); err != nil {
            log.Println("Erro ao atualizar contador em item_popularity_by_category:", err)
            continue
        }

        if (i+1)%1000 == 0 {
            fmt.Printf("%d interações registradas.\n", i+1)
        }
    }

    fmt.Printf("Inseridas %d interações.\n", numInteractions)
}
