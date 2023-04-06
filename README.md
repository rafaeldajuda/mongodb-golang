# MongoDB - Golang

Projeto golang utilizando a biblioteca MongoDB.

## MongoDB Docker

O comando abaixo irá baixar e rodar a imagem padrão do MongoDB. Neste comando foi definido o usuário, senha a e porta de acesso ao MongoDB. O nome do container será definido como 'my-mongo'.

```command
docker run -d -e MONGO_INITDB_ROOT_USERNAME=admin -e MONGO_INITDB_ROOT_PASSWORD=admin -p 27017:27017 --name my-mongo mongo
``` 

## Acesso ao Container

Para acessar o container basta rodar o comando abaixo passando o usuário e senha que foi definido anteriormente.

```command
docker exec -it my-mongo mongosh -u admin -p admin
``` 

## Documentos Para Testes

Aqui terá alguns comandos para alimentar o banco para os testes.

* Cria um banco.
```mongodb
use dev
```

* Cria uma coleção. 
```mongodb
db.createCollection("pessoa")
```

* Inseri dados na coleção.
```mongodb
db.pessoa.insertOne({"nome":"Rafael","idade":26,"cidade":"Taubaté"})
db.pessoa.insertOne({"nome":"José","idade":54,"cidade":"Taubaté"})
db.pessoa.insertOne({"nome":"Maria","idade":30,"cidade":"Tremembé"})
```

## Começando o Projeto

Iremos começar o projeto com o comando **go mod init github.com/rafaeldajuda/mongodb-golang** (você pode usar o nome que quiser). 

Criaremos a seguinte estrutura e variáveis globais.

(main.go)
```go
type MongoConfig struct {
	User       string
	Password   string
	Host       string
	Port       string
	Database   string
	Collection string
}

var mongoConfig MongoConfig
var collection *mongo.Collection
var ctx = context.TODO()
```

* mongoConfig: estrutura que conterá os dados de acesso ao MongoDB.
* collection: estrutura que terá acesso a coleção do banco.
* ctx: contexto que será utilizado nas operações.

Iremos utilizar a função **init**. Nela iremos carregar as variáveis de ambiente para alimentar a entidade **mongoConfig**. Para criar as variáveis de ambiente iremos utilizar a biblioteca **github.com/joho/godotenv**. Também precisamos criar um arquivo **.env** com as variáveis de ambiente que terá os acessos ao MongoDB.

(.env)
```env
MONGO_USER=admin
MONGO_PASSWORD=admin
MONGO_HOST=localhost
MONGO_PORT=27017
MONGO_DATABASE="dev"
MONGO_COLLECTION="pessoa"
```

(main.go)
```go
func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mongoConfig.User = os.Getenv("MONGO_USER")
	mongoConfig.Password = os.Getenv("MONGO_PASSWORD")
	mongoConfig.Host = os.Getenv("MONGO_HOST")
	mongoConfig.Port = os.Getenv("MONGO_PORT")
	mongoConfig.Database = os.Getenv("MONGO_DATABASE")
	mongoConfig.Collection = os.Getenv("MONGO_COLLECTION")
}
```

A função **godotenv.Load()** por padrão carrega o arquivo **.env**, mas é possível carregar outros arquivos.

```golang
godotenv.Load("arquivoqualquer")
godotenv.Load("arquivoum.env", "arquivodois.env")
```

## Conectando ao MongoDB

Na função **main** primeiro iremos montar a URI da conexão ao banco.

(main.go)
```golang
// mongodb uri format
uri := fmt.Sprintf("mongodb://%s:%s@%s:%s/",
	mongoConfig.User, mongoConfig.Password, mongoConfig.Host, mongoConfig.Port)
if uri == "" {
	log.Fatal("You must set your 'uri' variable.")
}
```

Iremos passar a URI para uma entidade **clientOpt** do tipo ***options.ClientOptions**. O ClientOptions possui várias configurações de conexão ao MongoDB, caso precise configurar uma conexão com regras específica será aqui que será adicionado.

Após isso iremos chamar a função **mongo.Connect()** passando o contexto e o clientOpt. Também terá uma função para encerrar a conexão e outra para dar um ping no banco. 

(main.go)
```go
// connection
clientOpt := options.Client().ApplyURI(uri)
client, err := mongo.Connect(ctx, clientOpt)
if err != nil {
	log.Fatal("MongoDB connection error: " + err.Error())
}

defer func() {
	if err := client.Disconnect(ctx); err != nil {
		log.Fatal("MongoDB disconnection error: " + err.Error())
	}
}()

// check connection
err = client.Ping(ctx, nil)
if err != nil {
	log.Fatal(err)
}
```

Com a conexão feita, iremos definir qual banco e coleção irá ser manipulada.

(main.go)
```golang
collection = client.Database(mongoConfig.Database).Collection(mongoConfig.Collection)
filter := bson.D{{Key: "nome", Value: "Rafael"}}
```

A variável **filter** é responsável por filtrar a consulta da query.

Para realizar a consulta é preciso chamar a função **FindOne()** que pertence ao **collection**. É preciso passar o contexto, o filtro e uma estrutura do tipo **bson.M** para receber o resultado.

A função **FindOne()** retorna somente um documento do banco a partir de um filtro.

(main.go)
```golang
// select one document
var result bson.M
err = collection.FindOne(ctx, filter).Decode(&result)
if err != nil {
	if err == mongo.ErrNoDocuments {
		fmt.Println("Document not found")
		return
	}
	log.Fatal(err)
}

// print result
resultB, _ := json.Marshal(result)
fmt.Println(string(resultB))
```

## Código Final

(.env)
```env
MONGO_USER=admin
MONGO_PASSWORD=admin
MONGO_HOST=localhost
MONGO_PORT=27017
MONGO_DATABASE="dev"
MONGO_COLLECTION="pessoa"
```

(main.go)
```golang
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConfig struct {
	User       string
	Password   string
	Host       string
	Port       string
	Database   string
	Collection string
}

var mongoConfig MongoConfig
var collection *mongo.Collection
var ctx = context.TODO()

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mongoConfig.User = os.Getenv("MONGO_USER")
	mongoConfig.Password = os.Getenv("MONGO_PASSWORD")
	mongoConfig.Host = os.Getenv("MONGO_HOST")
	mongoConfig.Port = os.Getenv("MONGO_PORT")
	mongoConfig.Database = os.Getenv("MONGO_DATABASE")
	mongoConfig.Collection = os.Getenv("MONGO_COLLECTION")
}

func main() {
	// mongodb uri format
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s/",
		mongoConfig.User, mongoConfig.Password, mongoConfig.Host, mongoConfig.Port)
	if uri == "" {
		log.Fatal("You must set your 'uri' variable.")
	}

	// connection
	clientOpt := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOpt)
	if err != nil {
		log.Fatal("MongoDB connection error: " + err.Error())
	}

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal("MongoDB disconnection error: " + err.Error())
		}
	}()

	// check connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	// set database and collection
	collection = client.Database(mongoConfig.Database).Collection(mongoConfig.Collection)
	filter := bson.D{{Key: "nome", Value: "Rafael"}}

	// select one document
	var result bson.M
	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("Document not found")
			return
		}
		log.Fatal(err)
	}

	// print result
	resultB, _ := json.Marshal(result)
	fmt.Println(string(resultB))

}
```

## Mais Exemplos

No seguinte link existe mais exemplos de como utilizar o MongoDB.

https://github.com/rafaeldajuda/mongodb-golang/examples
