package connection

import (
    "fmt"
    "log"
    "os"
    "path/filepath"

    "github.com/joho/godotenv"
    "github.com/jmoiron/sqlx"
    _ "github.com/godror/godror"          // Driver Oracle
    _ "github.com/nakagami/firebirdsql"   // Driver Firebird
    "database/sql"
)

var dsnFdb string
var dsnOracle string

func init() {
    envPath, err := os.Getwd()
    if err != nil {
        log.Fatalf("Erro ao obter diretório: %v", err)
    }

    if err = godotenv.Load(filepath.Join(envPath, ".env")); err != nil {
        log.Fatalf("Erro ao carregar .env: %v", err)
    }

    // Firebird
    dsnFdb = fmt.Sprintf("%s:%s@%s:%s/%s?charset=win1252&auth_plugin_name=Legacy_Auth",
        os.Getenv("FDB_USER"),
        os.Getenv("FDB_PASS"),
        os.Getenv("FDB_HOST"),
        os.Getenv("FDB_PORT"),
        os.Getenv("FDB_PATH"))

    // Oracle (usando SID)
    dsnOracle = fmt.Sprintf(`user="%s" password="%s" connectString="%s:%s/%s?sid=true"`,
        os.Getenv("ORA_USER"),
        os.Getenv("ORA_PASS"),
        os.Getenv("ORA_HOST"),
        os.Getenv("ORA_PORT"),
        os.Getenv("ORA_SERVICE"))
}

// GetConexoes retorna conexão Firebird via sql e Oracle via sqlx
func GetConexoes() (*sql.DB, *sqlx.DB, error) {
    // Conexão com Firebird
    ConexaoFdb, err := sql.Open("firebirdsql", dsnFdb)
    if err != nil {
        return nil, nil, fmt.Errorf("erro ao estabelecer conexão FDB: %v", err)
    }

    // Conexão com Oracle
    ConexaoOracle, err := sqlx.Connect("godror", dsnOracle)
    if err != nil {
        ConexaoOracle.Close()
        return nil, nil, fmt.Errorf("erro ao estabelecer conexão Oracle: %v", err)
    }

    return ConexaoFdb, ConexaoOracle, nil
}
