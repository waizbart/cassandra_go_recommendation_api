### **Visão Geral**

Em uma aplicação de recomendação que utiliza o Cassandra, as tabelas são projetadas para otimizar consultas específicas, seguindo uma abordagem orientada a consultas. Isso significa que, ao invés de normalizar dados como em bancos relacionais, no Cassandra você desnormaliza e cria tabelas que atendem diretamente às necessidades de consulta da sua aplicação, garantindo alto desempenho.

---

### **1. Tabela `items`**

**Propósito:**

- Armazena detalhes de cada item disponível no sistema, como `item_id`, `name` e `category`.
- Serve como a fonte de verdade para informações básicas sobre os itens.

**Estrutura:**

```sql
CREATE TABLE items (
    item_id UUID PRIMARY KEY,
    name TEXT,
    category TEXT
);
```

**Por que existe:**

- **Recuperação de detalhes do item**: Quando você precisa exibir informações sobre um item específico, como nome e categoria.
- **Referência centralizada**: Outros componentes do sistema podem referenciar esta tabela para obter informações atualizadas sobre os itens.

---

### **2. Tabela `items_by_category`**

**Propósito:**

- Indexa itens por categoria, permitindo consultas eficientes para listar itens de uma categoria específica.
- Facilita a obtenção de todos os itens pertencentes a uma determinada categoria.

**Estrutura:**

```sql
CREATE TABLE items_by_category (
    category TEXT,
    item_id UUID,
    name TEXT,
    PRIMARY KEY (category, item_id)
);
```

**Por que existe:**

- **Consultas otimizadas por categoria**: Permite recuperar rapidamente todos os itens de uma categoria sem a necessidade de varrer toda a tabela `items`.
- **Desempenho**: Evita operações de filtragem custosas em grandes volumes de dados.

---

### **3. Tabela `user_interactions`**

**Propósito:**

- Armazena detalhes completos sobre cada interação de um usuário com itens, incluindo `timestamp`, `item_id`, `category` e `action_type`.
- Permite rastrear o histórico de interações dos usuários, útil para análises e funcionalidades baseadas no histórico.

**Estrutura:**

```sql
CREATE TABLE user_interactions (
    user_id UUID,
    timestamp TIMESTAMP,
    item_id UUID,
    category TEXT,
    action_type TEXT,
    PRIMARY KEY (user_id, timestamp)
) WITH CLUSTERING ORDER BY (timestamp DESC);
```

**Por que existe:**

- **Histórico de atividades**: Permite recuperar todas as interações de um usuário em ordem cronológica.
- **Análises comportamentais**: Útil para entender o comportamento do usuário ao longo do tempo.
- **Consultas por intervalo de tempo**: Possibilita consultas em um período específico.

---

### **4. Tabela `user_interacted_items`**

**Propósito:**

- Mantém um registro simples dos itens com os quais o usuário já interagiu.
- Otimiza verificações rápidas para saber se um usuário já interagiu com um determinado item.

**Estrutura:**

```sql
CREATE TABLE user_interacted_items (
    user_id UUID,
    item_id UUID,
    PRIMARY KEY (user_id, item_id)
);
```

**Por que existe:**

- **Filtragem eficiente**: Ao gerar recomendações, é importante não sugerir itens que o usuário já viu. Esta tabela permite essa filtragem de forma rápida.
- **Desempenho**: Evita consultas pesadas na tabela `user_interactions`, que pode ter um volume grande de dados.

---

### **5. Tabela `user_category_interactions`**

**Propósito:**

- Armazena contadores de interações de um usuário em cada categoria.
- Permite identificar as categorias com as quais o usuário mais interage.

**Estrutura:**

```sql
CREATE TABLE user_category_interactions (
    user_id UUID,
    category TEXT,
    interaction_count COUNTER,
    PRIMARY KEY (user_id, category)
);
```

**Por que existe:**

- **Personalização de recomendações**: Ao saber quais categorias o usuário mais interage, é possível priorizar recomendações nessas categorias.
- **Análise de preferências**: Ajuda a entender as preferências do usuário em termos de categorias de produtos.

---

### **6. Tabela `item_popularity_by_category`**

**Propósito:**

- Mantém um contador de popularidade para cada item dentro de uma categoria.
- Auxilia na identificação dos itens mais populares em cada categoria.

**Estrutura:**

```sql
CREATE TABLE item_popularity_by_category (
    category TEXT,
    item_id UUID,
    popularity_counter COUNTER,
    PRIMARY KEY (category, item_id)
);
```

**Por que existe:**

- **Recomendações populares**: Permite recomendar itens populares dentro das categorias de interesse do usuário.
- **Desempenho**: Otimiza a consulta dos itens mais populares em uma categoria específica.

---

### **7. Tabela `item_popularity`** (se estiver em uso)

**Propósito:**

- Mantém um contador de popularidade geral de cada item, independente da categoria.
- Ajuda a identificar os itens mais populares em todo o catálogo.

**Estrutura:**

```sql
CREATE TABLE item_popularity (
    item_id UUID PRIMARY KEY,
    popularity_counter COUNTER
);
```

**Por que existe:**

- **Recomendações gerais**: Permite recomendar itens populares para novos usuários ou usuários sem histórico suficiente.
- **Análises de tendências**: Útil para identificar itens que estão em alta no momento.

---

### **8. Tabela `user_category_preference`** (se estiver em uso)

**Propósito:**

- Armazena as preferências do usuário por categoria, possivelmente calculadas a partir das interações.
- Pode ser utilizada para ajustar o peso das categorias na geração de recomendações.

**Estrutura:**

```sql
CREATE TABLE user_category_preference (
    user_id UUID,
    category TEXT,
    preference_score DOUBLE,
    PRIMARY KEY (user_id, category)
);
```

**Por que existe:**

- **Personalização avançada**: Refina as recomendações considerando não apenas o número de interações, mas também a relevância ou interesse do usuário em cada categoria.
- **Algoritmos de recomendação**: Pode ser usada em modelos mais sofisticados que consideram scores de preferência.

---

### **Razão Geral para Múltiplas Tabelas**

Em bancos de dados relacionais, você normalmente teria uma estrutura normalizada com tabelas e relacionamentos, e utilizaria joins para recuperar dados complexos. No Cassandra, porém, não há suporte eficiente para joins, e as operações de leitura e escrita são otimizadas para acesso direto a linhas com base na chave primária.

Portanto, para atender às diferentes necessidades de consulta da aplicação, você cria tabelas separadas, cada uma projetada para um tipo específico de consulta. Isso leva à desnormalização dos dados, mas permite que você tenha:

- **Consultas eficientes e rápidas**: Cada tabela é otimizada para um padrão de acesso específico.
- **Escalabilidade**: O Cassandra é altamente escalável e pode lidar com grandes volumes de dados distribuídos.
- **Alta disponibilidade**: A replicação e distribuição de dados garantem disponibilidade e tolerância a falhas.

---

### **Como as Tabelas se Integram na Aplicação**

Ao **registrar uma interação** do usuário:

1. **`user_interactions`**: Armazena o detalhe completo da interação, incluindo data e hora, tipo de ação, etc.
2. **`user_interacted_items`**: Registra que o usuário interagiu com o item, para filtragem futura nas recomendações.
3. **`user_category_interactions`**: Atualiza o contador de interações do usuário na categoria do item, ajudando a identificar preferências.
4. **`item_popularity_by_category`**: Atualiza o contador de popularidade do item dentro da categoria, para saber quais itens são mais populares.

Ao **gerar recomendações**:

1. **Identifica as categorias de interesse do usuário** usando `user_category_interactions`.
2. **Recupera os itens mais populares** nessas categorias através de `item_popularity_by_category`.
3. **Filtra os itens já interagidos** utilizando `user_interacted_items`, para não recomendar itens que o usuário já conhece.
4. **Obtém detalhes dos itens** da tabela `items` para compor a recomendação final.

---

### **Benefícios Desta Abordagem**

- **Performance**: Consultas rápidas e eficientes, essenciais para uma boa experiência do usuário.
- **Escalabilidade**: Capacidade de lidar com crescimento de usuários e interações sem degradação significativa de performance.
- **Flexibilidade**: Possibilidade de ajustar e otimizar o esquema conforme as necessidades de consulta evoluem.