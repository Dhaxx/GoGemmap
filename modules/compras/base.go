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

	tx, err := cnxFdb.Begin()
	if err != nil {
		panic(fmt.Sprintf("erro ao iniciar transação: %v", err))
	}
	defer tx.Commit()

	insert, err := tx.Prepare("insert into cadunimedida(sigla,descricao) values(?,?)")
	if err != nil {
		panic(fmt.Sprintf("erro ao preparar insert: %v", err.Error()))
	}

	query := `SELECT SIGLA, NOME AS descricao
	FROM (
		SELECT 
			SIGLA, 
			NOME, 
			ROW_NUMBER() OVER (PARTITION BY SIGLA ORDER BY NOME) AS rn
		FROM SYSTEM.D_UNID_MED
	) sub
	WHERE rn = 1`
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
	modules.NewCol("cadsubgr", "key_subgrupo")
	modules.NewCol("cadsubgr", "base")

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
	defer tx.Commit()

	query := `SELECT
		distinct
		G.FLG_TIPO || G.NRO AS key_grupo ,
		SUBSTR(S.NRO,2,3) AS subgrupo,
		'N' ocultar,
		S.FLG_TIPO||S.NRO as key_subgrupo,
		S.NOME
	FROM
		SYSTEM.D_GRUPO_MS S
			INNER JOIN SYSTEM.D_GRUPO_MS G ON
					S.GRMS_FLG_TIPO = G.FLG_TIPO
				AND S.GRMS_NRO = G.NRO
			left join (select GRMS_NRO nro, GRMS_FLG_TIPO tipo, row_number()
																over (partition by d.GRMS_FLG_TIPO, d.GRMS_NRO
																	order by  d.GRMS_FLG_TIPO, d.GRMS_NRO, d.nro) sequencia
					From system.D_MATSERV d
					order by  d.GRMS_FLG_TIPO, d.GRMS_NRO, d.nro) prod on s.nro = prod.nro
			and s.FLG_TIPO = prod.tipo
	WHERE
			S.FLG_NIVEL > 0
	union all
	select
		'9' || nro,
		'001',
		'N',
		'9'||nro,
		nome
	From system.e_GRUPO_PS
	ORDER BY
		key_grupo, subgrupo`

	cacheGrupo := make(map[string]string)
	rowsGrupo, err := cnxFdb.Query("SELECT grupo, conv_tipo||conv_nro key_grupo FROM cadgrupo")
	if err != nil {
		panic(fmt.Sprintf("erro ao executar query grupo: %v", err.Error()))
	}
	for rowsGrupo.Next() {	
		var grupo, key string

		if err := rowsGrupo.Scan(&grupo, &key); err != nil {
			panic(fmt.Sprintf("erro ao escanear registro grupo: %v", err.Error()))
		}
		cacheGrupo[key] = grupo
	}

	insertSubgrupo, err := tx.Prepare("insert into cadsubgr(grupo,subgrupo,nome,ocultar,key_subgrupo, base) values(?,?,?,?,?,?)")
	if err != nil {
		panic(fmt.Sprintf("erro ao preparar insert cadsubgr: %v", err.Error()))
	}

	bar := modules.NewProgressBar(p, 1, "Cadsubgr")

	rowsSubgrupo, err := cnxOra.Queryx(query)
	if err != nil {
		panic(fmt.Sprintf("erro ao executar query cadsubgr: %v", err.Error()))
	}

	for rowsSubgrupo.Next() {
		var registro ModelCadsubgr

		if err := rowsSubgrupo.StructScan(&registro); err != nil {
			panic(fmt.Sprintf("erro ao escanear registro cadsubgr: %v", err.Error()))
		}

		registro.Grupo = cacheGrupo[registro.Grupo]
		if registro.Grupo == "" {
			panic(fmt.Sprintf("grupo não encontrado para chave: %s", registro.Grupo))
		}

		registro.Nome.String, err = modules.DecodeToWin1252(registro.Nome.String); if err != nil {
			panic("erro ao decodificar nome: " + err.Error())
		}

		if _, err := insertSubgrupo.Exec(registro.Grupo, registro.Subgrupo, registro.Nome, registro.Ocultar, registro.KeySubgrupo, "S"); err != nil {
			panic(fmt.Sprintf("erro ao executar insert cadsubgr: %v", err.Error()))
		}
		bar.Increment()
	}
}

func Cadest(p *mpb.Progress) {
	modules.LimpaTabela([]string{"cadsubgr where base = 'N'", "cadest"})

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
	defer tx.Commit()

	query := `SELECT
		a.nro codreduz,
		--b.GRMS_FLG_TIPO || b.GRMS_NRO key_grupo,
		a.GRMS_FLG_TIPO || a.GRMS_NRO key_subgrupo,
		ROW_NUMBER() OVER (PARTITION BY a.GRMS_FLG_TIPO || a.GRMS_NRO
	ORDER BY
		a.GRMS_FLG_TIPO || a.GRMS_NRO,
		a.nro) codigo,
		a.nome disc1,
		est_min quanmin,
		est_max quanmax,
		CASE
			FLG_ATIVO WHEN 'S' THEN 'N'
			ELSE 'S'
		END AS ocultar,
		UM_SIGLA AS unid1,
		CASE
			WHEN b.GRMS_NRO = 1000
			AND b.GRMS_FLG_TIPO = '1' THEN 'P'
			WHEN b.GRMS_NRO = 2000
			AND b.GRMS_FLG_TIPO = '1' THEN 'E'
			WHEN b.GRMS_NRO = 1000
			AND b.GRMS_FLG_TIPO = '2' THEN 'P'
			WHEN b.GRMS_NRO = 1000
			AND b.GRMS_FLG_TIPO = '3' THEN 'S'
			WHEN b.GRMS_NRO = 1000
			AND b.GRMS_FLG_TIPO = '4' THEN 'S'
			WHEN b.GRMS_NRO = 2000
			AND b.GRMS_FLG_TIPO = '4' THEN 'S'
			WHEN b.GRMS_NRO = 3000
			AND b.GRMS_FLG_TIPO = '4' THEN 'S'
			ELSE 'P'
		END AS TIPOPRO ,
		CASE
			WHEN b.GRMS_FLG_TIPO = '2' THEN 'P'
			ELSE 'C'
		END AS USOPRO
	FROM
		system.D_MATSERV a
	JOIN system.D_GRUPO_MS b 
		ON
		a.GRMS_FLG_TIPO = b.FLG_TIPO
		AND a.GRMS_NRO = b.NRO
	UNION ALL
	SELECT
		nro,
		--'9'||a.egrps_nro,
		'9'||a.egrps_nro,	
		row_number() OVER (PARTITION BY '9'||a.egrps_nro ORDER BY '9'||a.egrps_nro, nro),
		nome,
		0,
		0,
		'N',
		um_sigla,
		'P',
		'C'
	FROM
		system.E_prod_serv a`
	
	insert, err := tx.Prepare("insert into cadest(grupo,subgrupo,codigo,cadpro,codreduz,disc1,quanmin,quanmax,ocultar,unid1,tipopro,usopro) values(?,?,?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		panic(fmt.Sprintf("erro ao preparar insert cadest: %v", err.Error()))
	}

	cacheSubgrupo := make(map[string][]string)
	querySubgrupo := `SELECT distinct grupo, SUBGRUPO, key_subgrupo FROM cadsubgr`
	rowsSubgrupo, err := cnxFdb.Query(querySubgrupo)
	if err != nil {
		panic(fmt.Sprintf("erro ao executar query grupo: %v", err.Error()))
	}
	for rowsSubgrupo.Next() {
		var grupo, subgrupo, key string
		if err := rowsSubgrupo.Scan(&grupo, &subgrupo, &key); err != nil {
			panic(fmt.Sprintf("erro ao escanear registro subgrupo: %v", err.Error()))
		}
		cacheSubgrupo[key] = []string{grupo, subgrupo} // Armazena o nome do subgrupo e o código
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

	subgruposItens := make(map[string]int64)

	for rows.Next() {
		var registro ModelCadest
		if err := rows.StructScan(&registro); err != nil {
			panic(fmt.Sprintf("erro ao escanear registro cadest: %v", err.Error()))
		}

		keySubgrupo := registro.Subgrupo
		subgruposItens[registro.Subgrupo]++

		registro.Grupo = cacheSubgrupo[registro.Subgrupo][0] // Obtém o grupo a partir do cache
		if registro.Grupo == "" {
			panic(fmt.Sprintf("grupo não encontrado para chave: %s", registro.Subgrupo))
		}
		registro.Subgrupo = cacheSubgrupo[registro.Subgrupo][1] // Obtém o subgrupo a partir do cache
		if registro.Subgrupo == "" {
			panic(fmt.Sprintf("subgrupo não encontrado para chave: %s", registro.Subgrupo))
		}
		
		if subgruposItens[keySubgrupo] > 999 {
			subgrupoNovo := modules.ExtourouSubgrupo(keySubgrupo)
			cacheSubgrupo[keySubgrupo] = []string{registro.Grupo, subgrupoNovo} // Atualiza o cache com o novo subgrupo
			registro.Subgrupo = subgrupoNovo // Atualiza o subgrupo
			registro.Codigo = registro.Codigo%100
			subgruposItens[keySubgrupo] = int64(registro.Codigo)
		}

		codigoFormatado := fmt.Sprintf("%03d", subgruposItens[keySubgrupo])
		
		registro.Cadpro = fmt.Sprintf("%s.%s.%s", registro.Grupo, registro.Subgrupo, codigoFormatado)
		
		registro.Disc1.String, err = modules.DecodeToWin1252(registro.Disc1.String); if err != nil {
			panic("erro ao decodificar disc1: "+err.Error())
		}
		registro.Unid1.String, err = modules.DecodeToWin1252(registro.Unid1.String); if err != nil {
			panic("erro ao decodificar unid1: "+err.Error())
		}

		if _, err := insert.Exec(registro.Grupo, registro.Subgrupo, codigoFormatado, registro.Cadpro, registro.Codreduz,
			registro.Disc1, registro.Quanmin, registro.Quanmax, registro.Ocultar, registro.Unid1, registro.Tipopro, registro.Usopro); err != nil {
			panic(fmt.Sprintf("erro ao executar insert cadest: %v", err.Error()))
		}

		bar.Increment()
	}
	tx.Commit()

	cadpros, err := cnxFdb.Query(`select cadpro, 
		cast(codreduz as integer) material, 
		case cast(g.conv_tipo as integer) when 9 then 9 else 1 end tipo
		From cadest t join cadgrupo g on g.GRUPO = t.GRUPO`)
	if err != nil {
		panic("Falha ao executar consulta: " + err.Error())
	}
	defer cadpros.Close()

	modules.Cache.Cadpros = make(map[string]string)
	for cadpros.Next() {
		var cadpro string
		var material int
		var tipo int
		if err := cadpros.Scan(&cadpro, &material, &tipo); err != nil {
			panic("Falha ao ler resultados da consulta: " + err.Error())
		}
		modules.Cache.Cadpros[fmt.Sprintf("%d|%d", tipo, material)] = cadpro
	}
}

func CentroCusto(p *mpb.Progress) {
	modules.LimpaTabela([]string{"centrocusto"})

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
	defer tx.Commit()

	if _, err = tx.Exec(`insert into centrocusto(poder,orgao,destino,ccusto,codccusto, empresa) 
		select first 1 poder,orgao,'000000001','001',0, empresa from 
		desdis where empresa =(select empresa from cadcli)`); err != nil {
		panic(fmt.Sprintf("erro ao executar insert centrocusto: %v", err.Error()))
	}

	insert, err := tx.Prepare(`insert
		into
		centrocusto(codccusto,
		descr,
		ccusto,
		ocultar,
		responsa,
		empresa,
		poder,
		orgao,
		unidade,
		destino)
	values(?,?,?,?,?,?,?,?,?,?)`)
	if err != nil {
		panic(fmt.Sprintf("erro ao preparar insert centrocusto: %v", err.Error()))
	}

	query := fmt.Sprintf(`SELECT
		nro codccusto,
		substr(nome,1,60) descr,
		'001' ccusto,
		case
			FLG_ATIVO when 'S' then 'N'
			else 'S'
		end ocultar,
		null responsa,
		%v empresa,
		coalesce((
		SELECT
			lpad(SECEST_NRO,
			9,
			'0')
		FROM
			SYSTEM.D_REQUISICAO r
		where
			r.DEPSEC_NRO = s.NRO
			and ROWNUM = 1),
		'000000001') destino
	FROM
		SYSTEM.DEPTO_SECAO s`, modules.Cache.Empresa)

	totalRows, err := modules.CountRows(query)
	if err != nil {
		panic(fmt.Sprintf("erro ao contar linhas centrocusto: %v", err.Error()))
	}
	bar := modules.NewProgressBar(p, totalRows, "CentroCusto")
	
	rows, err := cnxOra.Queryx(query)
	if err != nil {
		panic(fmt.Sprintf("erro ao executar query centrocusto: %v", err.Error()))
	}

	var poder, orgao string
	tx.QueryRow("select first 1 poder,orgao From desdis where empresa = (select empresa from cadcli)").Scan(&poder, &orgao)

	for rows.Next() {
		var registro ModelCentroCusto

		if err := rows.StructScan(&registro); err != nil {
			panic(fmt.Sprintf("erro ao escanear registro centrocusto: %v", err.Error()))
		}

		registro.Descr, err = modules.DecodeToWin1252(registro.Descr); if err != nil {
			panic("erro ao decodificar descr: " + err.Error())
		}
		registro.Responsa.String, err = modules.DecodeToWin1252(registro.Responsa.String); if err != nil {
			panic("erro ao decodificar responsa: " + err.Error())
		}

		if _, err := insert.Exec(registro.Codccusto, registro.Descr, registro.Ccusto, registro.Ocultar, registro.Responsa, registro.Empresa, poder, orgao, nil, registro.Destino); err != nil {
			panic(fmt.Sprintf("erro ao executar insert centrocusto: %v", err.Error()))
		}

		bar.Increment()
	}
}

func Destino(p *mpb.Progress) {
	modules.LimpaTabela([]string{"destino"})

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

	insert, err := tx.Prepare("insert into destino(cod,desti,empresa) values(?,?,?)")
	if err != nil {
		panic(fmt.Sprintf("erro ao preparar insert destino: %v", err.Error()))
	}

	query := fmt.Sprintf(`SELECT
		lpad(NRO, 9, '0') cod,
		NOME desti,
		%v EMPRESA
	FROM
		SYSTEM.D_SECR_ESTOQ`, modules.Cache.Empresa)

	totalRows, err := modules.CountRows(query)
	if err != nil {
		panic(fmt.Sprintf("erro ao contar linhas destino: %v", err.Error()))
	}
	bar := modules.NewProgressBar(p, totalRows, "Destino")

	rows, err := cnxOra.Queryx(query)
	if err != nil {
		panic(fmt.Sprintf("erro ao executar query destino: %v", err.Error()))
	}

	for rows.Next() {
		var registro ModelDestino

		if err := rows.StructScan(&registro); err != nil {
			panic(fmt.Sprintf("erro ao escanear registro destino: %v", err.Error()))
		}

		registro.Descr, err = modules.DecodeToWin1252(registro.Descr); if err != nil {
			panic("erro ao decodificar descr: " + err.Error())
		}

		if _, err := insert.Exec(registro.Destino, registro.Descr, registro.Empresa); err != nil {
			panic(fmt.Sprintf("erro ao executar insert destino: %v", err.Error()))
		}

		bar.Increment()
	}
	tx.Commit()
}