# find-binlog

## 📦 **Pré-requisitos**
- Go 1.23+
- Binário do `mysqlbinlog` compatível com sua versão do MySQL (ex: `8.0.mysql_aurora.3.08.2`)
- Posicione o `mysqlbinlog` em `pkg/bin/` com permissão de execução

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
* Este projeto foi testado com MySQL 8.0 e Aurora 3.08.2
* O binário `mysqlbinlog` deve ser da mesma versão do servidor
* Você pode baixá-lo manualmente em [https://dev.mysql.com/downloads/mysql/](https://dev.mysql.com/downloads/mysql/)

Para baixar o binário `mysqlbinlog` compatível com o **MySQL 8.0**, siga os passos abaixo:

---

### 📦 Baixar `mysqlbinlog` (MySQL 8.0)

#### 1. Acesse o site oficial da Oracle:

[https://dev.mysql.com/downloads/mysql/](https://dev.mysql.com/downloads/mysql/)

#### 2. Escolha:

* **Versão:** 8.0.36 (ou mais próxima da usada no seu Aurora)
* **OS:** Linux - Generic

#### 3. Baixe o pacote tar:

Exemplo para Linux x86_64:

```bash
wget https://dev.mysql.com/get/Downloads/MySQL-8.0/mysql-8.0.36-linux-glibc2.28-x86_64.tar.xz
```

#### 4. Extraia o binário:

```bash
tar -xf mysql-8.0.36-linux-glibc2.28-x86_64.tar.xz
```

#### 5. Copie o `mysqlbinlog` para seu projeto:

```bash
cp mysql-8.0.36-linux-glibc2.28-x86_64/bin/mysqlbinlog ./pkg/bin/
chmod +x ./pkg/bin/mysqlbinlog
```
