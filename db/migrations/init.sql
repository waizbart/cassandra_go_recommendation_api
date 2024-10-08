CREATE KEYSPACE IF NOT EXISTS recommendation_keyspace WITH replication = {'class': 'SimpleStrategy', 'replication_factor': '1'};

USE recommendation_keyspace;

-- 1. Tabela para armazenar interações dos usuários
CREATE TABLE IF NOT EXISTS user_interactions (
    user_id UUID,
    item_id UUID,
    action_type TEXT,
    category TEXT,
    timestamp TIMESTAMP,
    PRIMARY KEY (user_id, timestamp)
)
WITH CLUSTERING ORDER BY (timestamp DESC);

-- 2. Tabela para armazenar os itens
CREATE TABLE IF NOT EXISTS items (
    item_id UUID PRIMARY KEY,
    name TEXT,
    category TEXT
);

-- 3. Tabela para armazenar contadores de popularidade
CREATE TABLE IF NOT EXISTS item_popularity (
    item_id UUID PRIMARY KEY,
    popularity_counter COUNTER
) WITH CLUSTERING ORDER BY (popularity_counter DESC);
