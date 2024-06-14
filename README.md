## Passos para execução:
Execute o comando docker-compose:  
docker-compose up -d

Instale o 'REST Client' no 'VS Code' e execute o teste da pasta:  
./api/api.http  

Url jaeger:  
http://localhost:16686/  
Pesquisar pelos serviços serviceA e/ou serviceB

Url zipkin:  
http://localhost:9411/


## Tecnologias:
go 1.22
 - Router [chi](https://github.com/go-chi/chi)
 - Opentelemetry [otel](https://opentelemetry.io/docs/languages/go/getting-started/)
 - Opentelemetry - Span [otel-span](https://opentelemetry.io/docs/languages/go/instrumentation/#creating-spans)
 - Opentelemetry - Collector[otel-collector](https://opentelemetry.io/docs/collector/quick-start/)
 - Zipkin [zipkin] (https://zipkin.io/)
 - Prometheus
 - Jaeger
 
