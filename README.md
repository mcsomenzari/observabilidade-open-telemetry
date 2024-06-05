
## Passos para execução:

1. Subir container Docker
   docker-compose up --build

3. Testar api do servicoA
   Utilize o arquivo teste.http

5. Acompanhar os traces
   Abra um navegador e acesse o Zipkin UI em http://localhost:9411. Você poderá visualizar os traces gerados pelas requisições feitas aos serviços ao clicar em RUN QUERY.
 

## Requisitos 

### Serviço A (responsável pelo input):

O sistema deve receber um input de 8 dígitos via POST, através do schema:  { "cep": "29902555" }
O sistema deve validar se o input é valido (contem 8 dígitos) e é uma STRING
Caso seja válido, será encaminhado para o Serviço B via HTTP.

# Caso não seja válido, deve retornar:
- Código HTTP: 422
- Mensagem: CEP inválido.

### Serviço B (responsável pela orquestração):

O sistema deve receber um CEP válido de 8 digitos
O sistema deve realizar a pesquisa do CEP e encontrar o nome da localização, a partir disso, deverá retornar as temperaturas e formata-lás em: Celsius, Fahrenheit, Kelvin juntamente com o nome da localização.
O sistema deve responder adequadamente nos seguintes cenários:

# Em caso de sucesso
- Código HTTP: 200
- Response Body: { "city: "São Paulo", "temp_C": 28.5, "temp_F": 28.5, "temp_K": 28.5 }

# Em caso de falha, caso o CEP não seja válido (com formato correto)
- Código HTTP: 422
- Mensagem: CEP inválido.

# Em caso de falha, caso o CEP não seja encontrado
- Código HTTP: 404
- Mensagem: CEP não encontrado.
