### 

Este projeto consiste em dois serviços em Go:

- **Serviço A**: Recebe um CEP via POST, valida o formato e encaminha para o Serviço B.
- **Serviço B**: Recebe um CEP válido, consulta a localização através do serviço ViaCEP, e consulta a temperatura atual através do serviço WeatherAPI. Retorna a temperatura em Celsius, Fahrenheit e Kelvin, juntamente com o nome da cidade.


### Instruções para Execução

1. **Construir as Imagens**: Na raiz do projeto, execute:

   ```
   docker-compose build
   ```

2. **Executar o Docker Compose**: Na raiz do projeto (onde está o arquivo `docker-compose.yml`), execute:

   ```
   docker-compose up -d
   ```


3. **Testar a Aplicação**: 
Após iniciar os serviços, utilize o arquivo api.http ou 

- Para o serviço A:

   Teste CEP VÁLIDO:
   curl -X POST -H "Content-Type: application/json" -d "{\"cep\":\"02712080\"}" http://localhost:8080/cep
   
   Teste CEP INVÁLIDO:
   curl -X POST -H "Content-Type: application/json" -d "{\"cep\":\"00000000\"}" http://localhost:8080/cep
   
   
- Para o serviço B (substitua `CEP` por um código postal válido de 8 dígitos):
     ```
    curl -X GET "http://localhost:8081/weather?cep=02712080" -H "accept: application/json" 
     ```


4. **Testar o zipkin**:

- Abra o link no seu navegador:

   ```
   http://localhost:9411
   ```
   
![Clipboard02](https://github.com/mcsomenzari/observabilidade-open-telemetry/assets/29438629/dad5278b-a27a-46c7-b92c-a62fbbc988b6)
