# find-binlog

## 📖 **Finalidade**

`find-binlog` é uma ferramenta de linha de comando que varre os binlogs de um
servidor MySQL e retorna o primeiro evento ocorrido a partir de uma data
informada. O resultado exibe o arquivo e a posição do binlog que podem ser
utilizados para iniciar processos de replicação, investigações ou recuperações a
partir de um ponto específico. Se não houver evento na data exata, o programa
informa o mais próximo encontrado.

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
Execute o binário gerado especificando as credenciais e a data desejada:
```bash
./binlog-finder find-binlog \
  --host <host> \
  --port 3306 \
  --user <user> \
  --password <password> \
  --date YYYY-MM-DD \
  [--frameshot]
```
Use um usuário com permissão de leitura de binlog. O argumento `--date` deve
seguir o formato `YYYY-MM-DD`. O programa exibirá o nome do arquivo de log,
posição e data do evento encontrado, informações que podem ser usadas como
ponto inicial para replicação ou auditorias. Caso nenhum evento exista na data
informada, será mostrado o registro mais próximo. O parâmetro opcional
`--frameshot` salva em um arquivo as 100 linhas antes e depois do evento para
auxiliar inspeções manuais.

## 🧠 **Notas sobre compatibilidade**
* Este projeto foi testado com MySQL 8.0 e Aurora 3.08.2.
* Os binários `mysql` e `mysqlbinlog` que acompanham este repositório estão nessa mesma versão.
* Caso precise de outra versão, baixe os binários em [https://dev.mysql.com/downloads/mysql/](https://dev.mysql.com/downloads/mysql/) e substitua os arquivos em `pkg/bin/`.
