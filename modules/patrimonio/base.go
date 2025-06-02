package patrimonio

import (
	"GoGemmap/connection"
	"GoGemmap/modules"
	"fmt"

	"github.com/vbauerster/mpb"
)

func TipoMov(p *mpb.Progress) {
	modules.LimpaTabela([]string{"pt_tipomov"})

	cnxFdb, cnxOra, err := connection.GetConexoes()
	if err != nil {
		panic("Falha ao conectar com o banco de destino: " + err.Error())
	}
	defer cnxFdb.Close()
	defer cnxOra.Close()

	tx, err := cnxFdb.Begin()
	if err != nil {
		panic("Erro ao iniciar transação: " + err.Error())
	}
	defer tx.Commit()

	valores := map[string]string{
		"A": "Aquisição",
		"B": "Baixa",
		"T": "Transferência",
		"R": "Procedimento Contábil",
		"P": "Transferência de Plano Contábil",
	}

	insert, err := tx.Prepare("INSERT INTO PT_TIPOMOV (codigo_tmv, descricao_tmv) VALUES (?, ?)")
	if err != nil {
		fmt.Printf("Erro ao preparar insert: %v", err)
	}

	barTipoMov := modules.NewProgressBar(p, 1, "TIPOMOV")

	for sigla, descricao := range valores {
		descricaoConvertido1252, err := modules.DecodeToWin1252(descricao)
		if err != nil {
			fmt.Printf("Erro ao converter descrição para Win1252: %v", err)
		}

		_, err = insert.Exec(sigla, descricaoConvertido1252)
		if err != nil {
			fmt.Printf("Erro ao inserir valores: %v", err)
		}
	}
	barTipoMov.Completed()
}

func Cadajuste(p *mpb.Progress) {
	modules.LimpaTabela([]string{"pt_cadajuste"})

	cnxFdb, cnxOra, err := connection.GetConexoes()
	if err != nil {
		panic("Falha ao conectar com o banco de destino: " + err.Error())
	}
	defer cnxFdb.Close()
	defer cnxOra.Close()

	tx, err := cnxFdb.Begin()
	if err != nil {
		panic("Erro ao iniciar transação: " + err.Error())
	}
	defer tx.Commit()

	descricaoConvertido1252, err := modules.DecodeToWin1252("REAVALIAÇÃO (ANTES DO CORTE)")
	if err != nil {
		fmt.Printf("Erro ao converter descrição para Win1252: %v", err)
	}

	barCadAjuste := modules.NewProgressBar(p, 1, "CADAJUSTE")
	cnxFdb.Exec("INSERT INTO PT_CADAJUSTE (CODIGO_AJU, EMPRESA_AJU, DESCRICAO_AJU) VALUES (1, ?, ?)", modules.Cache.Empresa, descricaoConvertido1252)
	barCadAjuste.Completed()
}

func Cadbai(p *mpb.Progress) {
	modules.LimpaTabela([]string{"pt_cadbai"})

	cnxFdb, cnxOra, err := connection.GetConexoes()
	if err != nil {
		panic("Falha ao conectar com o banco de destino: " + err.Error())
	}
	defer cnxFdb.Close()
	defer cnxOra.Close()

	tx, err := cnxFdb.Begin()
	if err != nil {
		panic("Erro ao iniciar transação: " + err.Error())
	}
	defer tx.Commit()

	barCadBai := modules.NewProgressBar(p, 1, "CADBAI")

	query := `SELECT DISTINCT P.FLG_MODAL_BX codigo, CASE WHEN COALESCE(bx_valor,0) <> 0 THEN 'ALIENAÇÃO' ELSE 'BAIXA' END descricao_bai   FROM SYSTEM.D3_BEM_PATR P WHERE (P.BX_MOTIVO IS NOT NULL) `
	rows, err := cnxOra.Query(query)
	if err != nil {
		fmt.Printf("Erro ao executar query: %v", err)
	}

	for rows.Next() {
		var codigo, descricao string
		err := rows.Scan(&codigo, &descricao)
		if err != nil {
			fmt.Printf("Erro ao escanear valores: %v", err)
		}

		cnxFdb.Exec("INSERT INTO PT_CADBAI (CODIGO_BAI, EMPRESA_BAI, DESCRICAO_BAI) VALUES (?, ?, ?)", codigo, modules.Cache.Empresa, descricao)
	}
	barCadBai.Completed()
}

func Cadsit(p *mpb.Progress) {
	modules.LimpaTabela([]string{"pt_cadsit"})

	cnxFdb, cnxOra, err := connection.GetConexoes()
	if err != nil {
		panic("Falha ao conectar com o banco de destino: " + err.Error())
	}
	defer cnxFdb.Close()
	defer cnxOra.Close()

	tx, err := cnxFdb.Begin()
	if err != nil {
		panic("Erro ao iniciar transação: " + err.Error())
	}
	defer tx.Commit()

	query := `SELECT C.NRO, C.NOME FROM SYSTEM.D3_EST_CONSERV C ORDER BY C.NRO`

	totalRows, _ := modules.CountRows(query)
	barCadSit := modules.NewProgressBar(p, totalRows, "CADSIT")

	insert, err := tx.Prepare("INSERT INTO PT_CADSIT (CODIGO_SIT, EMPRESA_SIT, DESCRICAO_SIT) VALUES (?, ?, ?)")
	if err != nil {
		fmt.Printf("Erro ao preparar insert: %v", err)
	}

	rows, err := cnxOra.Query(query)
	if err != nil {
		fmt.Printf("Erro ao executar query: %v", err)
	}

	for rows.Next() {
		var codigo, descricao string
		err := rows.Scan(&codigo, &descricao)
		if err != nil {
			fmt.Printf("Erro ao escanear valores: %v", err)
		}

		if _, err := insert.Exec(codigo, modules.Cache.Empresa, descricao); err != nil {
			fmt.Printf("Erro ao inserir valores: %v", err)
		}
	}
	barCadSit.Completed()
}

func Cadtip(p *mpb.Progress) {
	modules.LimpaTabela([]string{"pt_cadtip"})

	cnxFdb, cnxOra, err := connection.GetConexoes()
	if err != nil {
		panic("Falha ao conectar com o banco de destino: " + err.Error())
	}
	defer cnxFdb.Close()
	defer cnxOra.Close()

	tx, err := cnxFdb.Begin()
	if err != nil {
		panic("Erro ao iniciar transação: " + err.Error())
	}
	defer tx.Commit()

	query := `SELECT E.NRO, E.NOME FROM SYSTEM.D3_BP_ESPECIE E ORDER BY E.NRO`

	totalRows, _ := modules.CountRows(query)
	barCadTip := modules.NewProgressBar(p, totalRows, "CADTIP")
	
	insert, err := tx.Prepare("INSERT INTO PT_CADTIP (CODIGO_TIP, EMPRESA_TIP, DESCRICAO_TIP) VALUES (?, ?, ?)")
	if err != nil {
		fmt.Printf("Erro ao preparar insert: %v", err)
	}

	rows, err := cnxOra.Query(query)
	if err != nil {
		fmt.Printf("Erro ao executar query: %v", err)
	}

	for rows.Next() {
		var codigo, descricao string
		err := rows.Scan(&codigo, &descricao)
		if err != nil {
			fmt.Printf("Erro ao escanear valores: %v", err)
		}

		if _, err := insert.Exec(codigo, modules.Cache.Empresa, descricao); err != nil {
			fmt.Printf("Erro ao inserir valores: %v", err)
		}
		barCadTip.Increment()
	}
}

func Cadpatd(p *mpb.Progress) {
	modules.LimpaTabela([]string{"pt_cadpatd"})

	cnxFdb, cnxOra, err := connection.GetConexoes()
	if err != nil {
		panic("Falha ao conectar com o banco de destino: " + err.Error())
	}
	defer cnxFdb.Close()
	defer cnxOra.Close()

	tx, err := cnxFdb.Begin()
	if err != nil {
		panic("Erro ao iniciar transação: " + err.Error())
	}
	defer tx.Commit()

	query := `SELECT
        DS.NRO,
        substr(DS.NOME,1,60) as NOME,
        --DS.FLG_DEPSEC,
        case when DS.FLG_ATIVO = 'S' then 'N' else 'S' end as OCULTAR_DES
        --DS.DEPSEC_NRO
    FROM
	    SYSTEM.DEPTO_SECAO DS
    WHERE
	    DS.DEPSEC_NRO IS NULL
    ORDER BY
	    DS.NRO`

	totalRows, _ := modules.CountRows(query)
	barCadPatd := modules.NewProgressBar(p, totalRows, "CADPATD")

	insert, err := tx.Prepare("INSERT INTO PT_CADPATD (EMPRESA_DES, CODIGO_DES, NAUNI_DES, OCULTAR_DES) VALUES (?, ?, ?, ?)")
	if err != nil {
		panic(fmt.Sprintf("Erro ao preparar insert: %v", err))
	}

	rows, err := cnxOra.Query(query)
	if err != nil {
		panic(fmt.Sprintf("Erro ao executar query: %v", err))
	}

	for rows.Next() {
		var (
			codigo, descricao, ocultar string
		)
		if err := rows.Scan(&codigo, &descricao, &ocultar); err != nil {
			panic(fmt.Sprintf("Erro ao escanear valores: %v", err))
		}

		descricaoConvertido1252, err := modules.DecodeToWin1252(descricao)
		if err != nil {
			panic(fmt.Sprintf("Erro ao converter descrição para Win1252: %v", err))
		}

		if _, err := insert.Exec(modules.Cache.Empresa, codigo, descricaoConvertido1252, ocultar); err != nil {
			panic(fmt.Sprintf("Erro ao inserir valores: %v", err))
		}

		barCadPatd.Increment()
	}
}

func Cadpats(p *mpb.Progress) {
	modules.LimpaTabela([]string{"pt_cadpats"})

	cnxFdb, cnxOra, err := connection.GetConexoes()
	if err != nil {
		panic("Falha ao conectar com o banco de destino: " + err.Error())
	}
	defer cnxFdb.Close()
	defer cnxOra.Close()

	tx, err := cnxFdb.Begin()
	if err != nil {
		panic("Erro ao iniciar transação: " + err.Error())
	}
	defer tx.Commit()

	query := `SELECT
        DS.NRO,
        substr(DS.NOME,1,60) as NOME,
        --DS.FLG_DEPSEC,
        case when DS.FLG_ATIVO = 'S' then 'N' else 'S' end as OCULTAR_DES,
        DS.DEPSEC_NRO
    FROM
	    SYSTEM.DEPTO_SECAO DS
    WHERE
	    DS.DEPSEC_NRO IS NOT NULL AND EXISTS (SELECT 1 FROM system.D3_BP_DS x where x.DEPSEC_NRO = ds.NRO)
    ORDER BY
	    DS.NRO`

	totalRows, _ := modules.CountRows(query)
	barCadPats := modules.NewProgressBar(p, totalRows, "CADPATS")

	insert, err := tx.Prepare("insert into pt_cadpats (codigo_set, empresa_set, codigo_des_set, noset_set, ocultar_set) values (?,?,?,?,?)")
	if err != nil {
		fmt.Printf("Erro ao preparar insert: %v", err)
	}

	rows, err := cnxOra.Query(query)
	if err != nil {
		fmt.Printf("Erro ao executar query: %v", err)
	}

	for rows.Next() {
		var (
			codigo_set, codigo_des_set int
			noset_set, ocultar_set string
		)

		if err := rows.Scan(&codigo_set, &noset_set, &ocultar_set, &codigo_des_set); err != nil {
			fmt.Printf("Erro ao escanear valores: %v", err)
		}

		nosetConvertido1252, err := modules.DecodeToWin1252(noset_set)
		if err != nil {
			fmt.Printf("Erro ao converter descrição para Win1252: %v", err)
		}

		if _, err := insert.Exec(codigo_set, modules.Cache.Empresa, codigo_des_set, nosetConvertido1252, ocultar_set); err != nil {
			fmt.Printf("Erro ao inserir valores: %v", err)
		}
		barCadPats.Increment()
	}
}

func Cadpatg(p *mpb.Progress) {
	modules.LimpaTabela([]string{"pt_cadpatg"})

	cnxFdb, cnxOra, err := connection.GetConexoes()
	if err != nil {
		panic("Falha ao conectar com o banco de destino: " + err.Error())
	}
	defer cnxFdb.Close()
	defer cnxOra.Close()

	if _, err := cnxFdb.Exec(fmt.Sprintf("INSERT INTO PT_CADPATG (CODIGO_GRU,EMPRESA_GRU,NOGRU_GRU) VALUES (1,%v,'GERAL')", modules.Cache.Empresa)); err != nil {
		fmt.Printf("Erro ao inserir valores: %v", err)
	}
}