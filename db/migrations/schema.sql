-- Criação do keyspace
CREATE KEYSPACE IF NOT EXISTS recommendation_keyspace
WITH replication = {'class': 'SimpleStrategy', 'replication_factor': '1'};

-- Uso do keyspace
USE recommendation_keyspace;

-- Tabela 'items' para armazenar informações básicas dos itens
CREATE TABLE IF NOT EXISTS items (
    item_id UUID PRIMARY KEY,
    name TEXT,
    category TEXT
);

-- Tabela 'items_by_category' para recuperar itens por categoria
CREATE TABLE IF NOT EXISTS items_by_category (
    category TEXT,
    item_id UUID,
    name TEXT,
    PRIMARY KEY (category, item_id)
);

-- Tabela 'user_interactions' para registrar interações dos usuários com os itens
CREATE TABLE IF NOT EXISTS user_interactions (
    user_id UUID,
    timestamp TIMESTAMP,
    item_id UUID,
    category TEXT,
    action_type TEXT,
    PRIMARY KEY (user_id, timestamp)
) WITH CLUSTERING ORDER BY (timestamp DESC);

-- Tabela 'user_category_interactions' para manter a contagem de interações por usuário e categoria
CREATE TABLE IF NOT EXISTS user_category_interactions (
    user_id UUID,
    category TEXT,
    interaction_count COUNTER,
    PRIMARY KEY (user_id, category)
);

-- Tabela 'item_popularity_by_category' para manter a popularidade dos itens por categoria
CREATE TABLE IF NOT EXISTS item_popularity_by_category (
    category TEXT,
    item_id UUID,
    popularity_counter COUNTER,
    PRIMARY KEY (category, item_id)
);

-- Tabela 'user_interacted_items' para registrar itens com os quais o usuário já interagiu
CREATE TABLE IF NOT EXISTS user_interacted_items (
    user_id UUID,
    item_id UUID,
    PRIMARY KEY (user_id, item_id)
);
