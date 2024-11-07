CREATE KEYSPACE IF NOT EXISTS recommendation_keyspace
WITH replication = {'class': 'SimpleStrategy', 'replication_factor': '2'};

USE recommendation_keyspace;

CREATE TABLE IF NOT EXISTS items (
    item_id UUID PRIMARY KEY,
    name TEXT,
    category TEXT
);

CREATE TABLE IF NOT EXISTS items_by_category (
    category TEXT,
    item_id UUID,
    name TEXT,
    PRIMARY KEY (category, item_id)
);

CREATE TABLE IF NOT EXISTS user_interactions (
    user_id UUID,
    timestamp TIMESTAMP,
    item_id UUID,
    category TEXT,
    action_type TEXT,
    PRIMARY KEY (user_id, timestamp)
) WITH CLUSTERING ORDER BY (timestamp DESC);

CREATE TABLE IF NOT EXISTS user_category_interactions (
    user_id UUID,
    category TEXT,
    interaction_count COUNTER,
    PRIMARY KEY (user_id, category)
);

CREATE TABLE IF NOT EXISTS item_popularity_by_category (
    category TEXT,
    item_id UUID,
    popularity_counter COUNTER,
    PRIMARY KEY (category, item_id)
);

CREATE TABLE IF NOT EXISTS user_interacted_items (
    user_id UUID,
    item_id UUID,
    PRIMARY KEY (user_id, item_id)
);
