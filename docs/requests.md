### Endpoints da API

- **POST `/items`**: Cria um novo item no sistema.
- **POST `/interactions`**: Registra uma interação de um usuário com um item.
- **GET `/recommendations/:user_id`**: Obtém recomendações de itens para um usuário específico.

### 1. Criar um Item

**Endpoint**: `POST /items`

**Descrição**: Este endpoint cria um novo item com um nome e uma categoria.

**Exemplo de Requisição com `curl`:**

```bash
curl -X POST http://localhost:8080/items \
-H "Content-Type: application/json" \
-d '{
    "name": "Produto Exemplo",
    "category": "Eletrônicos"
}'
```

**Explicação:**

- **Método**: POST
- **URL**: `http://localhost:8080/items`
- **Cabeçalho**: `Content-Type: application/json` para indicar que o corpo da requisição está em formato JSON.
- **Corpo da Requisição**:

  ```json
  {
    "name": "Produto Exemplo",
    "category": "Eletrônicos"
  }
  ```

**Resposta Esperada (201 Created):**

```json
{
  "message": "Item criado com sucesso",
  "item_id": "ID único do item criado"
}
```

**Exemplo de Resposta:**

```json
{
  "message": "Item criado com sucesso",
  "item_id": "d290f1ee-6c54-4b01-90e6-d701748f0851"
}
```

### 2. Registrar uma Interação

**Endpoint**: `POST /interactions`

**Descrição**: Registra uma interação de um usuário com um item, incluindo o tipo de ação realizada (por exemplo, "view", "purchase").

**Exemplo de Requisição com `curl`:**

```bash
curl -X POST http://localhost:8080/interactions \
-H "Content-Type: application/json" \
-d '{
    "user_id": "ID do usuário",
    "item_id": "ID do item",
    "action_type": "view"
}'
```

**Explicação:**

- **Método**: POST
- **URL**: `http://localhost:8080/interactions`
- **Cabeçalho**: `Content-Type: application/json`
- **Corpo da Requisição**:

  ```json
  {
    "user_id": "ID do usuário",
    "item_id": "ID do item",
    "action_type": "view"
  }
  ```

**Onde:**

- **`user_id`**: UUID representando o usuário. Você pode gerar um UUID se não tiver um.
- **`item_id`**: UUID do item com o qual o usuário interagiu. Este ID é retornado ao criar um item.
- **`action_type`**: Tipo de ação realizada. Pode ser "view", "purchase", "add_to_cart", etc.

**Exemplo Completo:**

Assumindo que você tem os seguintes IDs:

- **`user_id`**: `7b9e8f10-5e4a-11ec-bf63-0242ac130002`
- **`item_id`**: `d290f1ee-6c54-4b01-90e6-d701748f0851` (obtido ao criar o item)

Requisição:

```bash
curl -X POST http://localhost:8080/interactions \
-H "Content-Type: application/json" \
-d '{
    "user_id": "7b9e8f10-5e4a-11ec-bf63-0242ac130002",
    "item_id": "d290f1ee-6c54-4b01-90e6-d701748f0851",
    "action_type": "view"
}'
```

**Resposta Esperada (200 OK):**

```json
{
  "message": "Interaction logged successfully"
}
```

### 3. Obter Recomendações

**Endpoint**: `GET /recommendations/:user_id`

**Descrição**: Retorna uma lista de itens recomendados para o usuário especificado pelo `user_id`.

**Exemplo de Requisição com `curl`:**

```bash
curl -X GET http://localhost:8080/recommendations/7b9e8f10-5e4a-11ec-bf63-0242ac130002
```

**Explicação:**

- **Método**: GET
- **URL**: `http://localhost:8080/recommendations/{user_id}`
- **Substitua** `{user_id}` pelo UUID do usuário.

**Resposta Esperada (200 OK):**

```json
{
  "recommendations": [
    {
      "item_id": "UUID do item recomendado",
      "name": "Nome do item",
      "category": "Categoria do item"
    },
    {
      "item_id": "UUID do item recomendado",
      "name": "Nome do item",
      "category": "Categoria do item"
    }
    // Até 5 itens recomendados
  ]
}
```

**Exemplo de Resposta:**

```json
{
  "recommendations": [
    {
      "item_id": "a1b2c3d4-5e6f-7a8b-9c0d-1e2f3a4b5c6d",
      "name": "Produto Recomendado 1",
      "category": "Eletrônicos"
    },
    {
      "item_id": "b2c3d4e5-6f7a-8b9c-0d1e-2f3a4b5c6d7e",
      "name": "Produto Recomendado 2",
      "category": "Eletrônicos"
    }
    // Outros itens
  ]
}
```

### 4. Exemplo de Fluxo Completo

Para facilitar, vamos simular um fluxo completo de utilização da API.

#### Passo 1: Criar Itens

**Criando dois itens:**

- **Item 1**

  ```bash
  curl -X POST http://localhost:8080/items \
  -H "Content-Type: application/json" \
  -d '{
      "name": "Smartphone XYZ",
      "category": "Eletrônicos"
  }'
  ```

- **Item 2**

  ```bash
  curl -X POST http://localhost:8080/items \
  -H "Content-Type: application/json" \
  -d '{
      "name": "Livro ABC",
      "category": "Livros"
  }'
  ```

**Anote os `item_id` retornados em cada resposta.**

#### Passo 2: Registrar Interações

**Gerar um `user_id`** (se ainda não tiver um):

```bash
uuidgen
```

**Exemplo de `user_id`:**

```
e7b1c3d2-4f5a-678b-90cd-1e2f3a4b5c6d
```

**Registrar interações:**

- **Usuário visualiza o Smartphone XYZ**

  ```bash
  curl -X POST http://localhost:8080/interactions \
  -H "Content-Type: application/json" \
  -d '{
      "user_id": "e7b1c3d2-4f5a-678b-90cd-1e2f3a4b5c6d",
      "item_id": "item_id_do_Smartphone_XYZ",
      "action_type": "view"
  }'
  ```

- **Usuário compra o Livro ABC**

  ```bash
  curl -X POST http://localhost:8080/interactions \
  -H "Content-Type: application/json" \
  -d '{
      "user_id": "e7b1c3d2-4f5a-678b-90cd-1e2f3a4b5c6d",
      "item_id": "item_id_do_Livro_ABC",
      "action_type": "purchase"
  }'
  ```

#### Passo 3: Obter Recomendações

**Requisição:**

```bash
curl -X GET http://localhost:8080/recommendations/e7b1c3d2-4f5a-678b-90cd-1e2f3a4b5c6d
```

**Resposta Esperada:**

Uma lista de até 5 itens recomendados, preferencialmente em categorias com as quais o usuário mais interagiu, excluindo os itens já interagidos.

### 5. Gerar UUIDs

Se você precisar gerar UUIDs para `user_id` ou `item_id`, você pode:

- **No Linux/MacOS**, usar o comando:

  ```bash
  uuidgen
  ```

- **Em Go**, utilizar o pacote `github.com/google/uuid`:

  ```go
  import "github.com/google/uuid"

  newUUID := uuid.New().String()
  ```

### 6. Dicas Adicionais

- **Testar com Postman ou Insomnia**: Ferramentas como Postman ou Insomnia podem facilitar o teste da API, permitindo salvar requisições e visualizar respostas de forma organizada.

- **Headers**: Sempre inclua o header `Content-Type: application/json` em requisições POST com corpo JSON.

- **Porta da API**: Certifique-se de que a API está rodando na porta correta (no exemplo, usei `8080`). Ajuste a URL se necessário.

### 7. Lidar com Erros

**Exemplos de possíveis erros e como resolvê-los:**

- **Erro 400 Bad Request**:

  - **Causa**: Corpo da requisição inválido ou campos ausentes.
  - **Solução**: Verifique se todos os campos necessários estão presentes e se o JSON está formatado corretamente.

- **Erro 500 Internal Server Error**:

  - **Causa**: Problemas internos no servidor, como falha ao acessar o banco de dados.
  - **Solução**: Verifique os logs do servidor para identificar o problema. Certifique-se de que o Cassandra está em execução e acessível.

### 8. Exemplo de Script para Testes

Se quiser automatizar os testes, você pode criar um script em **Bash**:

```bash
#!/bin/bash

API_URL="http://localhost:8080"

# Criar item
item_response=$(curl -s -X POST $API_URL/items \
-H "Content-Type: application/json" \
-d '{
    "name": "Produto de Teste",
    "category": "Testes"
}')

item_id=$(echo $item_response | jq -r '.item_id')

echo "Item criado com ID: $item_id"

# Gerar user_id
user_id=$(uuidgen)

echo "Usuário criado com ID: $user_id"

# Registrar interação
curl -s -X POST $API_URL/interactions \
-H "Content-Type: application/json" \
-d "{
    \"user_id\": \"$user_id\",
    \"item_id\": \"$item_id\",
    \"action_type\": \"view\"
}"

echo "Interação registrada."

# Obter recomendações
recommendations=$(curl -s -X GET $API_URL/recommendations/$user_id)

echo "Recomendações para o usuário $user_id:"
echo $recommendations | jq
```

**Observações:**

- Certifique-se de ter o `jq` instalado para processar JSON no terminal.
- Dê permissão de execução ao script:

  ```bash
  chmod +x test_api.sh
  ```

- Execute o script:

  ```bash
  ./test_api.sh
  ```

### 9. Testando Erros

Para garantir que sua API lida corretamente com erros, você pode testar cenários como:

- **Enviar uma requisição sem um campo obrigatório**.
- **Usar um `user_id` ou `item_id` inválido**.
- **Tentar registrar uma interação com um `item_id` que não existe**.

**Exemplo de Requisição com `item_id` inválido:**

```bash
curl -X POST http://localhost:8080/interactions \
-H "Content-Type: application/json" \
-d '{
    "user_id": "e7b1c3d2-4f5a-678b-90cd-1e2f3a4b5c6d",
    "item_id": "id-invalido",
    "action_type": "view"
}'
```

**Resposta Esperada:**

```json
{
  "error": "Failed to retrieve item category"
}
```