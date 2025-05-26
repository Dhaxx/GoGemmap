package compras

import (
	"GoGemmap/connection"
	"GoGemmap/modules"
	"fmt"
	"github.com/vbauerster/mpb"
)

func Cadunimedida(p *mpb.Progress) {
	modules.LimpaTabela([]string{"cadunimedida"})

	cnxFdb, cnxOra, err := connection.GetConexoes()
    if err != nil {
		panic(fmt.Sprintf("erro ao obter conexões: %v", err.Error()))
    }
    defer cnxFdb.Close()
    defer cnxOra.Close()

	insert, err := cnxFdb.Prepare("insert into cadunimedida(sigla,descricao) values(?,?)")
	if err != nil {
		panic(fmt.Sprintf("erro ao preparar insert: %v", err.Error()))
	}

	query := `SELECT SIGLA , NOME descricao FROM SYSTEM.D_UNID_MED`
	rows, err := cnxOra.Queryx(query)
	if err != nil {
		panic(fmt.Sprintf("erro ao executar query: %v", err.Error()))
	}

	totalRows, err := modules.CountRows(query)
	if err != nil {
		panic(fmt.Sprintf("erro ao contar linhas: %v", err.Error()))
	}
	bar := modules.NewProgressBar(p, totalRows, "Cadunimedida")
	
	for rows.Next() {
		var registro ModelCadunimedida
		if err := rows.StructScan(&registro); err != nil {
			panic(fmt.Sprintf("erro ao escanear registro: %v", err.Error()))
		}

		registro.Descricao, err = modules.DecodeToWin1252(registro.Descricao); if err != nil {
			panic("erro ao decodificar: "+err.Error())
		}
		registro.Sigla, err = modules.DecodeToWin1252(registro.Sigla); if err != nil {
			panic("erro ao decodificar: "+err.Error())
		}

		if _, err = insert.Exec(registro.Sigla, registro.Descricao); err != nil {
			panic(fmt.Sprintf("erro ao executar insert: %v", err.Error()))
		}

		bar.Increment()
	}
	cnxFdb.Close()
}

func Grupo(p *mpb.Progress) {
	modules.LimpaTabela([]string{"cadsubgr", "cadgrupo"})
	modules.NewCol("cadgrupo", "conv_tipo")
	modules.NewCol("cadgrupo", "conv_nro")

	cnxFdb, cnxOra, err := connection.GetConexoes()
	if err != nil {
		panic(fmt.Sprintf("erro ao obter conexões: %v", err.Error()))
	}
	defer cnxFdb.Close()
	defer cnxOra.Close()

	tx, err := cnxFdb.Begin()
	if err != nil {
		fmt.Printf("erro ao iniciar transação: %v", err)
	}

	query := `SELECT
		nome,
		'N' ocultar,
		nro,
		FLG_TIPO tipo,
		null balco_tce,
		null balco_tce_saida
	FROM
		SYSTEM.D_GRUPO_MS
	WHERE
			FLG_NIVEL = 0
	union all
	select nome,
		'N',
		nro,
		'9',
		null,
		null
	From system.e_GRUPO_PS`
	
	insert, err := tx.Prepare("insert into cadgrupo(grupo,nome,balco_tce,balco_tce_saida,ocultar,conv_tipo,conv_nro) values(?,?,?,?,?,?,?)")
	if err != nil {
		panic(fmt.Sprintf("erro ao preparar insert cadgrupo: %v", err.Error()))
	}

	rows, err := cnxOra.Queryx(query)
	if err != nil {
		panic(fmt.Sprintf("erro ao executar query cadgrupo: %v", err.Error()))
	}

	totalRows, err := modules.CountRows(query)
	if err != nil {
		panic(fmt.Sprintf("erro ao contar linhas cadgrupo: %v", err.Error()))
	}
	bar := modules.NewProgressBar(p, totalRows, "Cadgrupo")

	var sequencia int64

	for rows.Next() {
		sequencia += 1
		var registro ModelCadgrupo
		if err := rows.StructScan(&registro); err != nil {
			panic(fmt.Sprintf("erro ao escanear registro cadgrupo: %v", err.Error()))
		}

		registro.Grupo = fmt.Sprintf("%03d", sequencia)

		registro.Nome, err = modules.DecodeToWin1252(registro.Nome); if err != nil {
			panic("erro ao decodificar nome: "+err.Error())
		}

		if _, err := insert.Exec(registro.Grupo, registro.Nome, registro.BalcoTce, registro.BalcoTceSaida, registro.Ocultar, registro.ConvTipo, registro.ConvNro); err != nil {
			panic(fmt.Sprintf("erro ao executar insert cadgrupo: %v", err.Error()))
		}

		bar.Increment()
	}
	tx.Commit()
}

func Subgrupo(p *mpb.Progress) {
	modules.LimpaTabela([]string{"cadsubgr"})

	cnxFdb, cnxOra, err := connection.GetConexoes()
    if err != nil {
		panic(fmt.Sprintf("erro ao obter conexões: %v", err.Error()))
    }
    defer cnxFdb.Close()
    defer cnxOra.Close()

	tx, err := cnxFdb.Begin()
	if err != nil {
		fmt.Printf("erro ao iniciar transação: %v", err)
	}

	query := `SELECT
		DISTINCT
		G.FLG_TIPO AS tipo_grupo ,
		G.NRO AS nro_grupo ,
		SUBSTR(S.NRO, 2, 3) AS subgrupo,
		lpad(round(COALESCE(sequencia, 0) / 999), 3, '0') sequencia,
		'N' ocultar,
		S.FLG_TIPO AS tipo_subgrupo ,
		S.NRO AS nro_subgrupo,
		S.NOME
	FROM
		SYSTEM.D_GRUPO_MS S
	INNER JOIN SYSTEM.D_GRUPO_MS G ON
		S.GRMS_FLG_TIPO = G.FLG_TIPO
		AND S.GRMS_NRO = G.NRO
	LEFT JOIN (
		SELECT
			GRMS_NRO nro,
			GRMS_FLG_TIPO tipo,
			ROW_NUMBER() OVER (PARTITION BY d.GRMS_FLG_TIPO,
			d.GRMS_NRO
		ORDER BY
			d.GRMS_FLG_TIPO,
			d.GRMS_NRO,
			d.nro) sequencia
		FROM
			system.D_MATSERV d
		ORDER BY
			d.GRMS_FLG_TIPO,
			d.GRMS_NRO,
			d.nro) prod ON
		s.nro = prod.nro
		AND s.FLG_TIPO = prod.tipo
	WHERE
		S.FLG_NIVEL > 0
	UNION ALL
	SELECT
		'9',
		nro,
		'001',
		'000',
		'N',
		'9',
		nro,
		nome
	FROM
		system.e_GRUPO_PS
	ORDER BY
		tipo_grupo ,
		nro_grupo`
	
	insert, err := tx.Prepare("insert into cadsubgr(grupo,subgrupo,nome,ocultar,conv_tipo,conv_nro) values(?,?,?,?,?)")
	if err != nil {
		panic(fmt.Sprintf("erro ao preparar insert cadsubgr: %v", err.Error()))
	}

	cacheGrupo := make(map[string]string)
	queryGrupo := `SELECT grupo, conv_tipo||conv_nro key FROM cadgrupo`
	rowsGrupo, err := cnxFdb.Query(queryGrupo)
	if err != nil {
		panic(fmt.Sprintf("erro ao executar query cadgrupo: %v", err.Error()))
	}
	for rowsGrupo.Next() {
		var grupo, key string
		if err := rowsGrupo.Scan(&grupo, &key); err != nil {
			panic(fmt.Sprintf("erro ao escanear registro cadgrupo: %v", err.Error()))
		}
		cacheGrupo[key] = grupo
	}

	rows, err := cnxOra.Queryx(query)
	if err != nil {
		panic(fmt.Sprintf("erro ao executar query cadsubgr: %v", err.Error()))
	}

	totalRows, err := modules.CountRows(query)
	if err != nil {
		panic(fmt.Sprintf("erro ao contar linhas cadsubgr: %v", err.Error()))
	}
	bar := modules.NewProgressBar(p, totalRows, "Cadsubgr")

	for rows.Next() {
		var registro ModelCadsubgr
		if err := rows.StructScan(&registro); err != nil {
			panic(fmt.Sprintf("erro ao escanear registro cadsubgr: %v", err.Error()))
		}

		key := fmt.Sprintf("%s%s", registro.ConvTipoGrupo, registro.ConvNroGrupo)
		registro.Grupo = cacheGrupo[key]
		if registro.Grupo == "" {
			panic(fmt.Sprintf("grupo não encontrado para chave: %s", key))
		}

		if _, err := insert.Exec(registro.Grupo, registro.Subgrupo, registro.Nome, registro.Ocultar, registro.ConvTipoSubgr, registro.ConvNroSubgr); err != nil {
			panic(fmt.Sprintf("erro ao executar insert cadsubgr: %v", err.Error()))
		}

		bar.Increment()
	}
}

func Cadest(p *mpb.Progress) {
	modules.LimpaTabela([]string{"cadest"})

	cnxFdb, cnxOra, err := connection.GetConexoes()
	if err != nil {
		panic(fmt.Sprintf("erro ao obter conexões: %v", err.Error()))
	}
	defer cnxFdb.Close()
	defer cnxOra.Close()

	tx, err := cnxFdb.Begin()
	if err != nil {
		panic(fmt.Sprintf("erro ao iniciar transação: %v", err))
	}

	query := `select
		lpad(row_number() over (partition by TIPO_SUBGRUPO , NRO_SUBGRUPO,seqcont order by TIPO_SUBGRUPO, NRO_SUBGRUPO,seqcont),3,'0') codigo,
		itens.*
	From (
		select
			lpad(round(coalesce(sequencia,0) / 999),3,'0') seqcont,
			produtos.*
		from (
				SELECT
							row_number() over (partition by d.GRMS_FLG_TIPO, d.GRMS_NRO order by  d.GRMS_FLG_TIPO, d.GRMS_NRO, d.nro) sequencia,
							D.GRMS_FLG_TIPO AS TIPO_SUBGRUPO,
							D.GRMS_NRO AS NRO_SUBGRUPO,
							substr(D.GRMS_NRO,2,3) AS subgrupo,
							D.NRO AS codreduz,
							D.NOME disc1,
							EST_MIN quanmin,
							EST_MAX quanmax,
							case FLG_ATIVO when 'S' then 'N' else 'S' end AS ocultar,
							UM_SIGLA AS unid1,
							CASE
								WHEN S.GRMS_NRO = 1000
									AND S.GRMS_FLG_TIPO = '1' THEN 'P'
								WHEN S.GRMS_NRO = 2000
									AND S.GRMS_FLG_TIPO = '1' THEN 'E'
								WHEN S.GRMS_NRO = 1000
									AND S.GRMS_FLG_TIPO = '2' THEN 'P'
								WHEN S.GRMS_NRO = 1000
									AND S.GRMS_FLG_TIPO = '3' THEN 'S'
								WHEN S.GRMS_NRO = 1000
									AND S.GRMS_FLG_TIPO = '4' THEN 'S'
								WHEN S.GRMS_NRO = 2000
									AND S.GRMS_FLG_TIPO = '4' THEN 'S'
								WHEN S.GRMS_NRO = 3000
									AND S.GRMS_FLG_TIPO = '4' THEN 'S'
								ELSE 'P'
								END AS TIPOPRO ,
							CASE
								WHEN S.GRMS_FLG_TIPO = '2' THEN 'P'
								ELSE 'C'
								END AS USOPRO
				FROM
					SYSTEM.D_MATSERV D
						INNER JOIN SYSTEM.D_GRUPO_MS S ON
								D.GRMS_NRO = S.NRO
							AND D.GRMS_FLG_TIPO = S.FLG_TIPO
				union all
				SELECT
					0,
					'9',
					EGRPS_NRO,
					lpad(EGRPS_NRO,3,'0'),
					nro,
					nome,
					0,
					0,
					'N',
					'UN',
					'P',
					'C' fROM SYSTEM.E_PROD_SERV
				order by TIPO_SUBGRUPO, NRO_SUBGRUPO, codreduz) produtos) itens`
	
	insert, err := tx.Prepare("insert into cadest(grupo,subgrupo,codigo,cadpro,codreduz,disc1,quanmin,quanmax,ocultar,unid1,tipopro,usopro) values(?,?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		panic(fmt.Sprintf("erro ao preparar insert cadest: %v", err.Error()))
	}

	cacheSubgrupo := make(map[string]string)
	querySubgrupo := `SELECT distinct GRUPO, CONV_TIPO||CONV_NRO AS key FROM CADSUBGR`
	rowsSubgrupo, err := cnxOra.Queryx(querySubgrupo)
	if err != nil {
		panic(fmt.Sprintf("erro ao executar query cadsubgr: %v", err.Error()))
	}
	for rowsSubgrupo.Next() {
		var grupo, key string
		if err := rowsSubgrupo.Scan(&grupo, &key); err != nil {
			panic(fmt.Sprintf("erro ao escanear registro cadsubgr: %v", err.Error()))
		}
		cacheSubgrupo[key] = grupo
	}

	rows, err := cnxOra.Queryx(query)
	if err != nil {
		panic(fmt.Sprintf("erro ao executar query cadest: %v", err.Error()))
	}
	totalRows, err := modules.CountRows(query)
	if err != nil {
		panic(fmt.Sprintf("erro ao contar linhas cadest: %v", err.Error()))
	}
	bar := modules.NewProgressBar(p, totalRows, "Cadest")
	
	for rows.Next() {
		var registro ModelCadest
		if err := rows.StructScan(&registro); err != nil {
			panic(fmt.Sprintf("erro ao escanear registro cadest: %v", err.Error()))
		}

		key := fmt.Sprintf("%s%s", registro.ConvTipo, registro.ConvNro)
		registro.Grupo = cacheSubgrupo[key]
		if registro.Grupo == "" {
			panic(fmt.Sprintf("subgrupo não encontrado para chave: %s", key))
		}

		registro.Cadpro = fmt.Sprintf("%s.%s.%s", registro.Grupo, registro.Subgrupo, registro.Codigo)
		
		registro.Disc1.String, err = modules.DecodeToWin1252(registro.Disc1.String); if err != nil {
			panic("erro ao decodificar disc1: "+err.Error())
		}
		registro.Unid1.String, err = modules.DecodeToWin1252(registro.Unid1.String); if err != nil {
			panic("erro ao decodificar unid1: "+err.Error())
		}

		if _, err := insert.Exec(registro.Grupo, registro.Subgrupo, registro.Codigo, registro.Cadpro, registro.Codreduz, registro.Disc1, registro.Quanmin, registro.Quanmax, registro.Ocultar, registro.Unid1, registro.Tipopro, registro.Usopro); err != nil {
			panic(fmt.Sprintf("erro ao executar insert cadest: %v", err.Error()))
		}

		bar.Increment()
	}
}

func CentroCusto(p *mpb.Progress) {
	modules.LimpaTabela([]string{"centrocusto"})
}