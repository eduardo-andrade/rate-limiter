
# Projeto Rate Limiter

Este projeto em Go implementa um serviço de rate limiting (limitação de taxa de requisições) utilizando Redis como backend para armazenamento e controle. O sistema permite limitar requisições por IP e por token, ajudando a proteger APIs e aplicações contra abuso e excesso de chamadas.

## Estrutura de Diretórios
```
rate-limiter/
├── config/
├── limiter/
├── middleware/
├── storage/
├── web/
│   └── test.html
├── main.go
├── go.mod
├── go.sum
├── Dockerfile
├── docker-compose.yml
└── README.md
```

## Requisitos
- Go 1.21 ou superior
- Docker e Docker Compose (opcional, para rodar via containers)
- Redis (pode ser local ou via container Docker)

## Instalação

1. **Clone o repositório**:
```bash
git clone <url-do-repositorio>
cd rate-limiter
```

2. **Baixe as dependências Go**:
```bash
go mod tidy
```

## Execução

### Rodando localmente

Configure a variável de ambiente `REDIS_ADDR` apontando para sua instância Redis (exemplo: `localhost:6379`):

```bash
export REDIS_ADDR=localhost:6379
go run main.go
```

O servidor vai iniciar na porta 8080.

### Rodando via Docker Compose

Construa e suba os containers com:

```bash
docker compose up --build
```

Isso vai iniciar o Redis e o serviço Go, expondo a aplicação na porta 8080.

## Testando o Rate Limiter

- Acesse o endpoint `/test` para abrir a página HTML de testes.
- Acesse `/test/run` com os parâmetros `testType` (ip ou token), `requests`, `interval`, `maxAllowed`, `ip` e `token` para realizar testes programáticos.
- Exemplo:
```
http://localhost:8080/test/run?testType=ip&requests=10&interval=100&maxAllowed=5&ip=127.0.0.1
```
- O Rate Limiter é muito robusto, permitindo testes com valores altos e baixo intervalo entre as requisições.

## Possíveis Erros e Soluções

- **Erro ao conectar no Redis localmente:**
  - Certifique-se de que o Redis está rodando na máquina ou ajuste `REDIS_ADDR` para apontar corretamente.
  - Exemplo: `export REDIS_ADDR=localhost:6379`

- **Erro ao rodar `docker compose`:**
  - Use o comando `docker compose` sem hífen (não `docker-compose`), pois versões recentes do Docker mudaram o comando.
  - Confirme que o Docker está rodando e que você tem permissão para usar o Docker.

- **Arquivo `test.html` não encontrado dentro do container:**
  - Certifique-se que o Dockerfile copia a pasta `web` corretamente para o container.
  - Veja o Dockerfile no repositório para garantir que `COPY web ./web` está presente.

## Considerações Finais

Este projeto serve para controlar a taxa de requisições em APIs, evitando abusos e sobrecarga. Pode ser estendido para diferentes estratégias de limitação e integração com sistemas maiores.
