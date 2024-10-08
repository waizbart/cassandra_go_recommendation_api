package main

import (
    "cassandra_go_recommendation_api/db"
    "cassandra_go_recommendation_api/routes"
    "log"
)

func main() {
    db.InitCassandra()
    defer db.CloseCassandra()

    r := routes.SetupRouter()
    if err := r.Run(); err != nil {
        log.Fatal("Erro ao iniciar o servidor:", err)
    }
}
