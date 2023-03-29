# MongoDB - Golang

Projeto golang utilizando a biblioteca MongoDB.

## MongoDB Docker

O comando abaixo irá baixar e rodar a imagem padrão do MongoDB. Neste comando foi definido o usuário, senha a e porta de acesso ao MongoDB. O nome do container será definido como 'my-mongo'.

```command
docker run -d -e MONGO_INITDB_ROOT_USERNAME=admin -e MONGO_INITDB_ROOT_PASSWORD=admin -p 27017:27017 --name my-mongo mongo
``` 

## Container Access

Para acessar o container basta rodar o comando abaixo passando o usuário e senha que foi definido anteriormente.

```command
docker exec -it my-mongo mongosh -u admin -p admin
``` 

## Set MongoDB Data

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




