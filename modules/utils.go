package modules

import (
	"GoGemmap/connection"
	"bytes"
	"database/sql"
	"fmt"

	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
)

var Cache struct {
	Empresa int
	Ano     int
}

func init() {
	cnxFdb, _, err := connection.GetConexoes()
	if err != nil {
		panic("Falha ao conectar com o banco de destino: " + err.Error())
	}
	defer cnxFdb.Close()

	cnxFdb.QueryRow("Select empresa from cadcli").Scan(&Cache.Empresa)
	cnxFdb.QueryRow("Select mexer from cadcli").Scan(&Cache.Ano)
}

func LimpaTabela(tabelas []string) {
	cnxFdb, _, err := connection.GetConexoes()
	if err != nil {
		fmt.Printf("Falha ao conectar com o banco de destino: %v", err)
	}
	defer cnxFdb.Close()

	tx, err := cnxFdb.Begin()
	if err != nil {
		fmt.Printf("erro ao iniciar transação: %v", err)
	}

	for _, tabela := range tabelas {
		if _, err = tx.Exec(fmt.Sprintf("DELETE FROM %v", tabela)); err != nil {
			fmt.Printf("erro ao limpar tabela: %v", err)
			tx.Rollback()
		}
	}
	tx.Commit()
}

func CountRows(q string, args ...any) (int64, error) {
	_, cnxPg, err := connection.GetConexoes()
	if err != nil {
		fmt.Printf("Falha ao conectar com o banco de destino: %v", err)
	}
	defer cnxPg.Close()

	var count int64
	query := fmt.Sprintf("SELECT count(*) FROM (%v) as subquery", q)

	if err := cnxPg.QueryRow(query).Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("nenhuma linha recuperada: %v", sql.ErrNoRows.Error())
		}
		return 0, fmt.Errorf("erro ao contar registros: %v", err)
	}
	return count, nil
}

func NewProgressBar(p *mpb.Progress, total int64, label string) *mpb.Bar {
	return p.AddBar(total,
		mpb.BarWidth(60),
		mpb.BarStyle("[██████░░░░░░]"),
		mpb.PrependDecorators(
			decor.Name(label+": "),
			decor.CountersNoUnit("%d / %d"),
		),
		mpb.AppendDecorators(
			decor.Percentage(),
			decor.EwmaETA(decor.ET_STYLE_GO, 60),
		),
	)
}

func NewCol(table string, colName string) {
	cnxFdb, _, err := connection.GetConexoes()
	if err != nil {
		fmt.Printf("Falha ao conectar com o banco de destino: %v", err)
	}
	defer cnxFdb.Close()

	tx, err := cnxFdb.Begin()
	if err != nil {
		fmt.Printf("erro ao iniciar transação: %v", err)
	}

	_, err = tx.Exec(fmt.Sprintf("ALTER TABLE %v ADD %v varchar(50)", table, colName))
	if err != nil {
		tx.Rollback()
		fmt.Printf("erro ao criar coluna %v: %v", colName, err)
	}

	tx.Commit()
}

func DecodeToWin1252(input string) (string, error) {
	// Define uma tabela de caracteres válidos no Windows-1252
	validChars := charmap.Windows1252

	// Remove ou substitui caracteres inválidos
	t := transform.Chain(
		runes.Remove(runes.Predicate(func(r rune) bool {
			// Remove caracteres que não são válidos no Windows-1252
			_, ok := validChars.EncodeRune(r)
			return !ok
		})),
		validChars.NewEncoder(),
	)

	// Transforma a string
	var buf bytes.Buffer
	writer := transform.NewWriter(&buf, t)

	_, err := writer.Write([]byte(input))
	if err != nil {
		return "", fmt.Errorf("erro ao codificar para Windows-1252: %w", err)
	}

	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("erro ao finalizar o writer: %w", err)
	}

	return buf.String(), nil
}

func LimpaLicitacoes() {
	cnxAux, _, err := connection.GetConexoes()
	if err != nil {
		panic("Falha ao conectar com o banco de destino: " + err.Error())
	}
	defer cnxAux.Close()

	_, err = cnxAux.Exec(`execute block as
		begin
		DELETE FROM regpreco;
		DELETE FROM regprecohis;
		DELETE FROM regprecodoc;
		DELETE FROM CADPROLIC_DETALHE_FIC;
		DELETE FROM CADPRO;
		DELETE FROM CADPRO_FINAL;
		DELETE FROM CADPRO_LANCE;
		DELETE FROM CADPRO_PROPOSTA;
		DELETE FROM PROLICS;
		DELETE FROM PROLIC;
		DELETE FROM CADPRO_STATUS;
		DELETE FROM CADLIC_SESSAO;
		DELETE FROM CADPROLIC_DETALHE;
		DELETE FROM CADPROLIC;
		DELETE FROM CADLIC;
		end;`)
	if err != nil {
		panic("Falha ao executar delete: " + err.Error())
	}
}

func LimpaCompras() {
	cnxAux, _, err := connection.GetConexoes()
	if err != nil {
		panic("Falha ao conectar com o banco de destino: " + err.Error())
	}
	defer cnxAux.Close()

	Trigger("TD_ICADREQ", false)

	_, err = cnxAux.Exec(`execute block as
		begin
		DELETE FROM ICADREQ;
		DELETE FROM REQUI;
		DELETE FROM ICADPED;
		DELETE FROM CADPED;
		DELETE FROM regpreco;
		DELETE FROM regprecohis;
		DELETE FROM regprecodoc;
		DELETE FROM CADPRO_SALDO_ANT;
		DELETE FROM CADPROLIC_DETALHE_FIC;
		DELETE FROM CADPRO;
		DELETE FROM CADPRO_FINAL;
		DELETE FROM CADPRO_LANCE;
		DELETE FROM CADPRO_PROPOSTA;
		DELETE FROM PROLICS;
		DELETE FROM PROLIC;
		DELETE FROM CADPRO_STATUS;
		DELETE FROM CADLIC_SESSAO;
		DELETE FROM CADPROLIC_DETALHE;
		DELETE FROM CADPROLIC;
		DELETE FROM CADLIC;
		DELETE FROM VCADORC;
		DELETE FROM FCADORC;
		DELETE FROM ICADORC;
		DELETE FROM CADORC;
		DELETE FROM CADEST;
		DELETE FROM CENTROCUSTO;
		DELETE FROM DESTINO;
		DELETE FROM DESFORCRC_PADRAO;
		end;`)
	if err != nil {
		panic("Falha ao executar delete: " + err.Error())
	}
}

func LimpaPatrimonio() {
	cnxAux, _, err := connection.GetConexoes()
	if err != nil {
		panic("Falha ao conectar com o banco de destino: " + err.Error())
	}
	defer cnxAux.Close()

	_, err = cnxAux.Exec(`execute block as
		begin
		DELETE FROM PT_CADPAT_EMPEN;
		DELETE FROM PT_MOVBEM;
		DELETE FROM PT_CADPAT;
		DELETE FROM PT_CADPATS;
		DELETE FROM PT_CADPATD;
		DELETE FROM PT_CADPATG;
		DELETE FROM PT_CADTIP;
		DELETE FROM PT_CADSIT;
		DELETE FROM PT_CADBAI;
		DELETE FROM PT_CADAJUSTE;
		DELETE FROM PT_TIPOMOV;
		DELETE FROM PT_CADRESPONSAVEL;
		end;`)
	if err != nil {
		panic("Falha ao executar delete: " + err.Error())
	}
}

func Trigger(trigger string, status bool) {
	cnxFdb, _, err := connection.GetConexoes()
	if err != nil {
		panic("Falha ao conectar com o banco de destino: " + err.Error())
	}
	defer cnxFdb.Close()

	tx, err := cnxFdb.Begin()
	if err != nil {
		panic("erro ao iniciar transação: " + err.Error())
	}
	defer tx.Commit()

	var statusStr string
	if status {
		statusStr = "ACTIVE"
	} else {
		statusStr = "INACTIVE"
	}

	tx.Exec(fmt.Sprintf("ALTER TRIGGER %s %s", trigger, statusStr))
}