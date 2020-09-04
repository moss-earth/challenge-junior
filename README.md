# Pacotes usados

"gorilla/mux" para gerenciar as rotas que irei implementar

"lib/pq" para rodar o database postgres

"joho/godotenv" para leitura de arquivos .env que irão manter dados seguros

# Banco de dados
Tentei usar o PostgreSQL pelo método Docker Image mas acabei descobrindo que o docker não instala na máquina que tenho disponivel comigo, anexarei o print do erro na pasta raiz da API para caso possam me auxiliar posteriormente

Optei pelo método Cloud Based utilizando o serviço do ElephantSQL que é bem fácil de settar (Além de ser grátis sem colocar cartão de crédito) seguindo essa documentação:
"https://www.elephantsql.com/docs/index.html" 
Anexarei a operação SQL realizada para criar as duas tables do projeto na pasta raiz.


# Usaremos o Postman para os testes
Primeiro precisamos startar nosso server

Com seu terminal aberto em nossa pasta raiz execute:

```javascript
$ go run main.go
```

Com o postman aberto, podemos agora começar nossos testes
# Testes para a table projects

##### Usaremos primeiro o _POST_ para criar nosso project

(OBS: o conteudo do Body que utilizaremos deve ser raw/JSON em todos os testes com body)

URL: http://localhost:8080/api/newproject

Body(*trocadilho ínfame* :tw-1f47d:):

```javascript
{
    "name": "mossgo"
}
```
é interessante fazer mais de um

##### Usaremos agora o método _GET_ para capturar um project

URL: http://localhost:8080/api/project/1

O projectid é passado como parametro na url

##### Usaremos agora o _GetALL_ para capturar todos os project da table projects usando o método GET

URL: http://localhost:8080/api/project

##### Usaremos agora o método _PUT_ para atualizar o nome de um project

URL: http://localhost:8080/api/project/1

use o body exemplo

```javascript
{
    "name": "gomoss"
}
```
##### Por último o método _DELETE_ para atualizar o nome de um project

URL: http://localhost:8080/api/deleteproject/1

# Agora podemos iniciar os testes com os lot e a table Lots

##### Usaremos primeiro o _POST_ para criar nosso primeiro lot

Para fazer o teste de criar lot no Postman é necessário usar o Pre-request script para capturar o "current time"

```javascript
{
    var current_timestamp = new Date();
	postman.setEnvironmentVariable("current_timestamp",
	current_timestamp.toISOString());
}
```

Use a URL abaixo 
http://localhost:8080/api/newlot

Preencha o Body no Postman com numeros no price e quantity, buytime captura o tempo atual e tem que se atentar a colocar um projectID existente

```javascript
{
    "price":15,
    "quantity":100,
    "buytime": "{{current_timestamp}}",
    "projectID":2
}
```

##### Agora o método GET para pegar um Lot pelo ID do mesmo. (Aqui usaremos 1 para capturar o lotid 1)

URL: http://localhost:8080/api/lot/1

##### Agora o método GET para pegar todo Lot de um project pelo projectID do mesmo. (Aqui usaremos 2 para capturar os lots de projectID 2)

URL: http://localhost:8080/api/lot/2

##### Agora o método DELETE para excluir um Lot

URL: http://localhost:8080/api/deletelot/1

##### Por último e não menos importante, agradeço pela equipe da Moss por esse desafio, espero poder atender as expectativas com esse projeto, mas já adianto que foi um grande aprendizado e uma bela oportunidade para acelerar minha entrada no mundo do Go. Abraços e bom fim de semana!