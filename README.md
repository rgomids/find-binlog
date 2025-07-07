# find-binlog

## 📦 **Pré-requisitos**
- Go 1.23+
- Os binários `mysql` e `mysqlbinlog` já estão incluídos em `pkg/bin/` (versão 8.0). Se precisar de outra versão, substitua-os manualmente.

## 🚀 **Instalação**
Para compilar o projeto:
```bash
make build
```

## 🧪 **Rodando os testes**
```bash
make test
```

## 📂 **Gerando pacote final**
```bash
make package
```
* Resultado será colocado na pasta `dist/`

## 🕵️ **Uso**
```bash
./binlog-finder find-binlog \
  --host <host> \
  --user <user> \
  --password <password> \
  --date 2025-03-01
```

## 🧠 **Notas sobre compatibilidade**
* Este projeto foi testado com MySQL 8.0 e Aurora 3.08.2.
* Os binários `mysql` e `mysqlbinlog` que acompanham este repositório estão nessa mesma versão.
* Caso precise de outra versão, baixe os binários em [https://dev.mysql.com/downloads/mysql/](https://dev.mysql.com/downloads/mysql/) e substitua os arquivos em `pkg/bin/`.
