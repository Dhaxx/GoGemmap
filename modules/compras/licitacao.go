package compras

import (
	"GoGemmap/connection"
	"GoGemmap/modules"
	"fmt"
	"github.com/vbauerster/mpb"
)

func Cadlic(p *mpb.Progress) {
	modules.LimpaTabela([]string{"cadlic"})
	modules.NewCol("cadlic", "localiza")

	cnxFdb, cnxOra, err := connection.GetConexoes()
	if err != nil {
		panic(fmt.Sprintf("erro ao obter conexões: %v", err.Error()))
	}
	defer cnxFdb.Close()
	defer cnxOra.Close()

	tx, err := cnxFdb.Begin()
	if err != nil {
		panic(fmt.Sprintf("erro ao iniciar transação: %v", err.Error()))
	}

	insert, err := tx.Prepare(`INSERT
		INTO
		cadlic(
		numlic,
		proclic,
		numero,
		ano,
		comp,
		licnova,
		liberacompra,
		discr,
		detalhe,
		registropreco,
		microempresa,
		numpro,
		discr7,
		datae,
		processo_data,
		horabe,
		horreal,
		tipopubl,
		dtadj,
		dthom,
		localiza,
		codtce,
		anomod,
		modlic,
		licit,
		codmod,
		dtpub,
		dtenc,
		horenc,
		valor,
		empresa,
		processo,
		processo_ano,
		dtreal)
	VALUES ( ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
	?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `)
	if err != nil {
		panic(fmt.Sprintf("erro ao preparar insert cadlic: %v", err.Error()))
	}

	query := fmt.Sprintf(`SELECT
		NRO AS NUMLIC,
		lpad(row_number() over (partition by EX_ANO order by  EX_ANO, NRO_PROC),6,'0')  || '/' ||substr(ex_ano,3,2) sequencia,
		lpad(NRO_PROC,6,'0')  || '/' ||substr(ex_ano,3,2) proclic,
		NRO_DOC numero,
		EX_ANO ano,
		3 comp,
		1 licnova,
		'S' liberacompra,
		substr(OBJETO,1,512) discr,
		OBJETO detalhe,
		FLG_REGPRECO AS REGISTROPRECO,
		0 microempresa,
		NRO_PROC numpro,
		'Menor Preco Unitario' discr7,
		DT_LICITACAO AS datae,
		DT_LICITACAO AS dtreal,
		DT_ABER_PROP AS processo_data,
		substr(HR_ABER_PROP,1,4) AS horabe,
		substr(HR_ENTR_PROP,1,4) AS horreal,
		'Outros' tipopubl,
		DT_HOMOLOGACAO dtadj,
		DT_HOMOLOGACAO dthom,
		substr(REPARTICAO,1,100) local,
		COD_LICITACAO codtce,
		EX_ANO anomod,
		MODLIC_NRO AS modalidade,
		DT_PUBLICACAO dtpub,
		DT_ENCERR dtenc,
		substr(HR_ENCERR,1,4) horenc,
		VLR_GLOBAL_LIC AS VALOR,
		--'00000001' lotelic,
		--1 sessao,
		%v empresa
	FROM
		SYSTEM.D_LICITACAO`, modules.Cache.Empresa)

	totalRows, err := modules.CountRows(query)
	if err != nil {
		panic(fmt.Sprintf("erro ao contar linhas: %v", err.Error()))
	}
	bar := modules.NewProgressBar(p, totalRows, "Cadlic")

	rows, err := cnxOra.Queryx(query)
	if err != nil {
		panic(fmt.Sprintf("erro ao executar query cadlic: %v", err.Error()))
	}
	defer rows.Close()

	for rows.Next() {
		var registro ModelCadlic
		if err := rows.StructScan(&registro); err != nil {
			panic(fmt.Sprintf("erro ao ler resultados da consulta cadlic: %v", err.Error()))
		}

		var modalidade *modules.ProcessoLicitatorio
		for _, m := range modules.Modalidades {
			if m.Nro == registro.Modalidade {
				modalidade = &m
				break
			}
		}

		registro.Discr.String, err = modules.DecodeToWin1252(registro.Discr.String)
		if err != nil {
			panic(fmt.Sprintf("erro ao decodificar discr: %v", err.Error()))
		}
		registro.Discr7.String, err = modules.DecodeToWin1252(registro.Discr7.String)
		if err != nil {
			panic(fmt.Sprintf("erro ao decodificar discr7: %v", err.Error()))
		}
		registro.Detalhe.String, err = modules.DecodeToWin1252(registro.Detalhe.String)
		if err != nil {
			panic(fmt.Sprintf("erro ao decodificar detalhe: %v", err.Error()))
		}

		_, err = insert.Exec(
			registro.Numlic,
			registro.Proclic,
			registro.Numero,
			registro.Ano,
			registro.Comp,
			registro.Licnova,
			registro.Liberacompra,
			registro.Discr,
			registro.Detalhe,
			registro.Registropreco,
			registro.Microempresa,
			registro.Numpro,
			registro.Discr7,
			registro.Datae,
			registro.ProcessoData,
			registro.Horabe,
			registro.Horreal,
			registro.Tipopubl,
			registro.Dtadj,
			registro.Dthom,
			registro.Local,
			registro.Codtce,
			registro.Anomod,
			modalidade.Modlic,
			modalidade.Licit,
			modalidade.Codmod,
			registro.Dtpub,
			registro.Dtenc,
			registro.Horenc,
			registro.Valor,
			registro.Empresa,
			registro.Processo,
			registro.Ano,
			registro.Dtreal,
		)
		if err != nil {
			panic(fmt.Sprintf("erro ao inserir registro cadlic: %v", err.Error()))
		}
		bar.Increment()
	}
	tx.Commit()

	cnxFdb.Exec(`EXECUTE BLOCK AS
	DECLARE VARIABLE DESCMOD VARCHAR(1024);
	DECLARE VARIABLE CODMOD INTEGER;
	BEGIN
		FOR
			SELECT CODMOD, DESCMOD FROM MODLIC INTO :CODMOD, :DESCMOD
		DO
		BEGIN
			UPDATE CADLIC SET LICIT = :DESCMOD where CODMOD = :CODMOD;
		END
	END`)

	cnxFdb.Exec(`INSERT
		INTO
		MODLICANO (ULTNUMPRO,
		CODMOD,
		ANOMOD,
		EMPRESA)
	SELECT
		COALESCE(MAX(NUMPRO), 0),
		CODMOD,
		COALESCE(ANO, 0) ANO,
		EMPRESA
	FROM
		CADLIC c
	WHERE
		CODMOD IS NOT NULL
	GROUP BY
		2,
		3,
		4
	ORDER BY
		ano,
		codmod`)
	
	cnxFdb.Exec(`UPDATE CADLIC SET FK_MODLICANO = (SELECT PK_MODLICANO FROM MODLICANO WHERE CODMOD = CADLIC.CODMOD AND ANOMOD = CADLIC.ANO AND CADLIC.EMPRESA = MODLICANO.EMPRESA) WHERE CODMOD IS NOT NULL`)
	
	cnxFdb.Exec(`INSERT INTO CADLIC_SESSAO (NUMLIC, SESSAO, DTREAL, HORREAL, COMP, DTENC, HORENC, SESSAOPARA, MOTIVO) 
	SELECT L.NUMLIC, CAST(1 AS INTEGER), L.DTREAL, L.HORREAL, L.COMP, L.DTENC, L.HORENC, CAST('T' AS VARCHAR(1)), CAST('O' AS VARCHAR(1)) FROM CADLIC L 
	WHERE numlic not in (SELECT FIRST 1 S.NUMLIC FROM CADLIC_SESSAO S WHERE S.NUMLIC = L.NUMLIC)`)
}

func Prolics(p *mpb.Progress) {
	modules.LimpaTabela([]string{"prolics","prolic"})

	cnxFdb, cnxOra, err := connection.GetConexoes()
	if err != nil {
		panic(fmt.Sprintf("erro ao obter conexões: %v", err.Error()))
	}
	defer cnxFdb.Close()
	defer cnxOra.Close()

	tx, err := cnxFdb.Begin()
	if err != nil {
		panic(fmt.Sprintf("erro ao iniciar transação: %v", err.Error()))
	}
	defer tx.Commit()

	insertProlic, err := tx.Prepare(`insert into prolic(numlic,codif,nome,status) values(?,?,?,?)`)
	if err != nil {
		panic(fmt.Sprintf("erro ao preparar insert prolic: %v", err.Error()))
	}
	defer insertProlic.Close()
	
	insertProlics, err := tx.Prepare(`insert into prolics(sessao,numlic,codif,habilitado,status,cpf,representante) values(?,?,?,?,?,?,?)`)
	if err != nil {
		panic(fmt.Sprintf("erro ao preparar insert prolics: %v", err.Error()))
	}
	defer insertProlics.Close()

	query := `SELECT
		1 sessao,
		f.FO_PES_NRO codif,
		LIC_NRO AS numlic,
		'S' habilitado,
		'A' status,
		substr(P.NOME, 1, 40) nome,
		substr(REPRES_NOME, 1, 40) representante,
		REPRES_CPF AS CPF
	FROM
		SYSTEM.D_LIC_FOR F
	INNER JOIN SYSTEM.PESSOA P ON
		P.NRO = F.FO_PES_NRO
	INNER JOIN SYSTEM.D_LICITACAO L ON
		L.NRO = F.LIC_NRO`

	totalRows, err := modules.CountRows(query)
	if err != nil {
		panic(fmt.Sprintf("erro ao contar linhas: %v", err.Error()))
	}
	bar := modules.NewProgressBar(p, totalRows, "Prolics")
	rows, err := cnxOra.Queryx(query)
	if err != nil {
		panic(fmt.Sprintf("erro ao executar query prolics: %v", err.Error()))
	}
	defer rows.Close()

	for rows.Next() {
		var registro ModelProlics
		if err := rows.StructScan(&registro); err != nil {
			panic(fmt.Sprintf("erro ao ler resultados da consulta prolics: %v", err.Error()))
		}

		registro.Nome.String, err = modules.DecodeToWin1252(registro.Nome.String)
		if err != nil {
			panic(fmt.Sprintf("erro ao decodificar nome: %v", err.Error()))
		}
		registro.Representante.String, err = modules.DecodeToWin1252(registro.Representante.String)
		if err != nil {
			panic(fmt.Sprintf("erro ao decodificar representante: %v", err.Error()))
		}

		if _, err = insertProlic.Exec(
			registro.Numlic,
			registro.Codif,
			registro.Nome,
			registro.Status,
		); err != nil {
			panic(fmt.Sprintf("erro ao inserir prolic: %v", err.Error()))
		}

		if _, err = insertProlics.Exec(
			registro.Sessao,
			registro.Numlic,
			registro.Codif,
			registro.Habilitado,
			registro.Status,
			registro.Cpf,
			registro.Representante,
		); err != nil {
			panic(fmt.Sprintf("erro ao inserir prolics: %v", err.Error()))
		}
		bar.Increment()
	}
}

func Cadprolic(p *mpb.Progress) {
	modules.LimpaTabela([]string{"cadprolic_detalhe_fic", "cadprolic_detalhe", "cadprolic"})

	cnxFdb, cnxOra, err := connection.GetConexoes()
	if err != nil {
		panic(fmt.Sprintf("erro ao obter conexões: %v", err.Error()))
	}
	defer cnxFdb.Close()
	defer cnxOra.Close()

	tx, err := cnxFdb.Begin()
	if err != nil {
		panic(fmt.Sprintf("erro ao iniciar transação: %v", err.Error()))
	}
	defer tx.Commit()

	insertCadprolic, err := tx.Prepare(`INSERT
		INTO
		cadprolic(item,
		item_mask,
		itemorc,
		cadpro,
		quan1,
		vamed1,
		vatomed1,
		codccusto,
		ficha,
		reduz,
		ordnumorc,
		numlic,
		id_cadorc,
		lotelic,
		item_lote)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		panic(fmt.Sprintf("erro ao preparar insert cadprolic: %v", err.Error()))
	}
	defer insertCadprolic.Close()

	insertCadprolicDetalhe, err := tx.Prepare(`INSERT
		INTO
		cadprolic_detalhe(numlic,
		item,
		ordnumorc,
		numorc,
		itemorc,
		cadpro,
		quan1,
		vamed1,
		vatomed1,
		codccusto,
		ficha,
		item_cadprolic,
		id_cadorc)
	VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		panic(fmt.Sprintf("erro ao preparar insert cadprolic_detalhe: %v", err.Error()))
	}
	defer insertCadprolicDetalhe.Close()

	insertCadprolicDetalheFic, err := tx.Prepare(`INSERT
		INTO
		cadprolic_detalhe_fic(numlic,
		item,
		codigo,
		ficha,
		qtd,
		valor,
		qtdadt,
		valoradt,
		codccusto,
		qtdmed,
		valormed,
		tipo)
	VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		panic(fmt.Sprintf("erro ao preparar insert cadprolic_detalhe_fic: %v", err.Error()))
	}
	defer insertCadprolicDetalheFic.Close()

	query := fmt.Sprintf(`WITH ranked_items AS (
		SELECT
			dm.NOME,
			dli.NROSEQ,
			dlv.NROSEQ AS SEQ_VLR,
			dlv.LICIT_MTSV_NRO,
			dli.FLG_ANEXO,
			dlv.LICIT_LIC_NRO AS numlic,
			dlv.QUANT,
			dlv.VLR_UNIT,
			dlv.VLR_TOTAL,
			dlv.LICFOR_FO_PES_NRO codif, 
			COUNT(*) OVER (PARTITION BY dlv.LICIT_LIC_NRO, dlv.NROSEQ) AS QTD_DUPLICADAS,
			ROW_NUMBER() OVER (
				PARTITION BY dlv.LICIT_LIC_NRO, dlv.NROSEQ
				ORDER BY dlv.VLR_UNIT ASC
			) AS rn
		FROM
			system.D_LICIT_VLR dlv
		JOIN system.D_MATSERV dm ON dlv.LICIT_MTSV_NRO = dm.NRO
		JOIN system.D_LIC_ITENS dli ON dlv.LICIT_LIC_NRO = dli.LIC_NRO AND dlv.LICIT_MTSV_NRO = dli.MTSV_NRO AND dli.flg_anexo = 'N'
		JOIN system.D_LICITACAO dl ON dl.NRO = dlv.LICIT_LIC_NRO --AND dl.EX_ANO >= %v
		--WHERE dlv.flg_menor_preco = 1
	)
	SELECT
		SEQ_VLR item,
		numlic,
		dl.MTSV_NRO material,
		ranked_items.QUANT QUAN1,
		dl.VLR_MEDIO vamed1,
		dl.VLR_TOTAL valor,
		--codif,
		QTD_DUPLICADAS QUANTIDADE,
		'C' TIPO,
		0 codccusto,
		'N' reduz,
		NULL ficha,
		'00000001' lotelic
	FROM
		ranked_items
	LEFT JOIN system.V_LICITACAO_ITENS dl ON 
		dl.id_lic = NUMLIC 
		AND dl.MTSV_NRO  = ranked_items.LICIT_MTSV_NRO 
	WHERE
		rn = 1
	ORDER BY
		numlic,
		SEQ_VLR`, modules.Cache.Ano-5)

	totalRows, err := modules.CountRows(query)
	if err != nil {
		panic(fmt.Sprintf("erro ao contar linhas: %v", err.Error()))
	}
	bar := modules.NewProgressBar(p, totalRows, "Cadprolic")
	
	rows, err := cnxOra.Queryx(query)
	if err != nil {
		panic(fmt.Sprintf("erro ao executar query cadprolic: %v", err.Error()))
	}
	defer rows.Close()

	cacheCadest := make(map[string]string)
	queryCadest, err := cnxFdb.Query(`select codreduz, cadpro from cadest`)
	if err != nil {
		panic("Erro ao consultar cadest: " + err.Error())
	}
	defer queryCadest.Close()
	for queryCadest.Next() {
		var codreduz, cadpro string
		if err := queryCadest.Scan(&codreduz, &cadpro); err != nil {
			panic("Erro ao ler cadest: " + err.Error())
		}
		cacheCadest[codreduz] = cadpro
	}

	for rows.Next() {
		var registro ModelCadprolic
		if err := rows.StructScan(&registro); err != nil {
			panic(fmt.Sprintf("erro ao ler resultados da consulta cadprolic: %v", err.Error()))
		}

		registro.Cadpro = cacheCadest[registro.Codreduz.String]
		if registro.Cadpro == "" {
			panic(fmt.Sprintf("cadpro não encontrado para o código: %s", registro.Codreduz.String))
		}

		if _, err = insertCadprolic.Exec(
			registro.Item,
			registro.Item,
			nil,
			registro.Cadpro,
			registro.Quan1,
			registro.Vamed1,
			registro.Valor,
			registro.Codccusto,
			nil,
			registro.Reduz,
			nil,
			registro.Numlic,
			nil,
			registro.Lotelic,
			registro.Item,
		); err != nil {
			panic(fmt.Sprintf("erro ao inserir cadprolic: %v", err.Error()))
		}

		if _, err = insertCadprolicDetalhe.Exec(
			registro.Numlic,
			registro.Item,
			nil, // ordnumorc
			nil, // numorc
			nil, // itemorc
			registro.Cadpro,
			registro.Quan1,
			registro.Vamed1,
			registro.Valor,
			registro.Codccusto,
			nil,
			registro.Item,
			nil, // id_cadorc
		); err != nil {
			panic(fmt.Sprintf("erro ao inserir cadprolic_detalhe: %v", err.Error()))
		}

		if _, err = insertCadprolicDetalheFic.Exec(
			registro.Numlic,
			registro.Item,
			registro.Item, // codigo
			nil,
			registro.Quan1, // qtd
			registro.Valor, // valor
			registro.Quan1, // qtdadt
			registro.Valor, // valoradt
			registro.Codccusto,
			registro.Quan1, // qtdmed
			registro.Valor, // valormed
			registro.Tipo, // tipo
		); err != nil {
			panic(fmt.Sprintf("erro ao inserir cadprolic_detalhe_fic: %v", err.Error()))
		}

		bar.Increment()
	}
	tx.Commit()

	cnxFdb.Exec(`INSERT
		INTO
		cadlotelic (descr,
		lotelic,
		numlic) 
	SELECT distinct
		'Lote ' || lotelic,
		lotelic,
		numlic
	FROM
		cadprolic a
	WHERE
		lotelic IS NOT NULL
		AND NOT EXISTS (
		SELECT
			1
		FROM
			CADLOTELIC c
		WHERE
			c.numlic = a.numlic
			AND a.lotelic = c.lotelic)`)

	modules.Trigger("TBIU_CADPRO_STATUS", false)
	cnxFdb.Exec(`INSERT INTO cadpro_status (numlic, sessao, itemp, telafinal)
		SELECT b.NUMLIC, 1 AS sessao, a.item, 'I_ENCERRAMENTO'
	FROM CADPROLIC a
	JOIN cadlic b ON a.NUMLIC = b.NUMLIC
	WHERE NOT EXISTS (
			SELECT 1
			FROM cadpro_status c
			WHERE a.numlic = c.numlic and a.item = c.item)`)
	modules.Trigger("TBIU_CADPRO_STATUS", true)
}

func CadproProposta(p *mpb.Progress) {
	modules.LimpaTabela([]string{"cadpro","cadpro_final","cadpro_proposta"})
	modules.NewCol("cadpro_proposta", "cadpro")
	modules.NewCol("cadpro_proposta", "tpcontrole_saldo")
	modules.NewCol("cadpro_proposta", "qtdadt")
	modules.NewCol("cadpro_proposta", "vaunadt")

	cnxFdb, cnxOra, err := connection.GetConexoes()
	if err != nil {
		panic(fmt.Sprintf("erro ao obter conexões: %v", err.Error()))
	}
	defer cnxFdb.Close()
	defer cnxOra.Close()

	tx, err := cnxFdb.Begin()
	if err != nil {
		panic(fmt.Sprintf("erro ao iniciar transação: %v", err.Error()))
	}

	insert, err := tx.Prepare(`INSERT
		INTO
		cadpro_proposta(sessao,
		codif,
		item,
		itemp,
		quan1,
		vaun1,
		vato1,
		numlic,
		status,
		subem,
		marca,
		itemlance,
		lotelic,
		cadpro,
		tpcontrole_saldo,
		qtdadt,
		vaunadt)
	VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`)
	if err != nil {
		panic(fmt.Sprintf("erro ao preparar insert cadpro_proposta: %v", err.Error()))
	}
	defer insert.Close()

	// PROPOSTA DE LICITANTES
	query := `WITH max_rodada AS (
	SELECT
		LICIT_LIC_NRO AS numlic,
		LICFOR_FO_PES_NRO AS codif,
		LICIT_MTSV_NRO AS material,
		MAX(nro_rodada) AS nro_rodada
	FROM
		system.D_LANCES
	GROUP BY
		LICIT_LIC_NRO,
		LICFOR_FO_PES_NRO,
		LICIT_MTSV_NRO
	)
	SELECT 
		1 sessao, 
		p.numlic, 
		p.codif, 
		p.subem, 
		p.material, 
		p.item, 
		p.quan1, 
		p.vaun1, 
		p.vato1, 
		p.status,
		'S' item_lance,
		'00000001' lotelic,
		CASE WHEN p.quan1 = 1 THEN 'V' ELSE 'Q' END tpcontrole_saldo,
		p.nro_rodada
	FROM (
		SELECT
			dl.nro_rodada,
			dl.LICIT_LIC_NRO AS numlic,
			dl.LICFOR_FO_PES_NRO AS codif,
			CASE WHEN dl.FLG_VENCEDOR = 'S' THEN 1 ELSE 0 END AS subem,
			dl.LICIT_MTSV_NRO AS material,
			b.item,
			b.quan1,
			dl.vlr_lance AS vaun1,
			b.quan1 * dl.VLR_LANCE AS vato1,
			b.status
		FROM
			system.D_LANCES dl
		JOIN max_rodada mr 
			ON dl.LICIT_LIC_NRO = mr.numlic 
			AND dl.LICFOR_FO_PES_NRO = mr.codif 
			AND dl.LICIT_MTSV_NRO = mr.material
			AND dl.nro_rodada = mr.nro_rodada
		JOIN (
			WITH ranked_items AS (
					SELECT
						dm.NOME,
						dli.NROSEQ,
						dlv.NROSEQ AS SEQ_VLR,
						dlv.LICIT_MTSV_NRO,
						dli.FLG_ANEXO,
						dlv.LICIT_LIC_NRO AS numlic,
						dlv.QUANT,
						dlv.VLR_UNIT,
						dlv.VLR_TOTAL,
						dlv.LICFOR_FO_PES_NRO AS codif,
						CASE 
							WHEN dli.FLG_FRACASSADO = 'S' THEN 'D' 
							ELSE 'C' 
						END AS status,
						COUNT(*) OVER (PARTITION BY dlv.LICIT_LIC_NRO, dlv.NROSEQ) AS QTD_DUPLICADAS,
						ROW_NUMBER() OVER (
							PARTITION BY dlv.LICIT_LIC_NRO, dlv.NROSEQ
							ORDER BY dlv.VLR_UNIT ASC
						) AS rn
					FROM
						system.D_LICIT_VLR dlv
					JOIN system.D_MATSERV dm 
						ON dlv.LICIT_MTSV_NRO = dm.NRO
					JOIN system.D_LIC_ITENS dli 
						ON dlv.LICIT_LIC_NRO = dli.LIC_NRO 
						AND dlv.LICIT_MTSV_NRO = dli.MTSV_NRO 
						AND dli.flg_anexo = 'N'
					JOIN system.D_LICITACAO dl 
						ON dl.NRO = dlv.LICIT_LIC_NRO 
				)
				SELECT
					SEQ_VLR AS item,
					numlic,
					dl.MTSV_NRO AS material,
					ranked_items.QUANT AS quan1,
					status
				FROM
					ranked_items
				LEFT JOIN system.V_LICITACAO_ITENS dl 
					ON dl.id_lic = numlic 
					AND dl.MTSV_NRO = ranked_items.LICIT_MTSV_NRO 
				WHERE
					rn = 1
		) b 
			ON dl.LICIT_LIC_NRO = b.numlic 
			AND dl.LICIT_MTSV_NRO = b.material
	) p
	ORDER BY 
		p.numlic, 
		p.item`

	totalRows, err := modules.CountRows(query)
	if err != nil {
		panic(fmt.Sprintf("erro ao contar linhas: %v", err.Error()))
	}
	bar := modules.NewProgressBar(p, totalRows, "Cadpro Proposta")

	cacheCadest := make(map[string]string)
	queryCadest, err := cnxFdb.Query(`select codreduz, cadpro from cadest`)
	if err != nil {
		panic("Erro ao consultar cadest: " + err.Error())
	}
	defer queryCadest.Close()
	for queryCadest.Next() {
		var codreduz, cadpro string
		if err := queryCadest.Scan(&codreduz, &cadpro); err != nil {
			panic("Erro ao ler cadest: " + err.Error())
		}
		cacheCadest[codreduz] = cadpro
	}

	rows, err := cnxOra.Queryx(query)
	if err != nil {
		panic(fmt.Sprintf("erro ao executar query cadpro_proposta: %v", err.Error()))
	}
	defer rows.Close()

	for rows.Next() {
		var registro ModelProposta

		if err := rows.StructScan(&registro); err != nil {
			panic(fmt.Sprintf("erro ao ler resultados da consulta cadpro_proposta: %v", err.Error()))
		}

		registro.Cadpro = cacheCadest[registro.Codreduz.String]
		if registro.Cadpro == "" {
			panic(fmt.Sprintf("cadpro não encontrado para o código: %s", registro.Codreduz.String))
		}

		if _, err = insert.Exec(
			registro.Sessao,
			registro.Codif,
			registro.Item,
			registro.Item,
			registro.Quan1,
			registro.Vaun1,
			registro.Vato1,
			registro.Numlic,
			registro.Status,
			registro.Subem,
			registro.Marca,
			registro.ItemLance,
			registro.Lotelic,
			registro.Cadpro,
			registro.Tpcontrole,
			registro.QtdAdt,
			registro.VaunAdt,
		); err != nil {
			panic(fmt.Sprintf("erro ao inserir cadpro_proposta: %v", err.Error()))
		}
		bar.Increment()
	}

	// PROPOSTA DE DISPENSAS
	query = fmt.Sprintf(`WITH ranked_items AS (
	SELECT
		dm.NOME,
		dli.NROSEQ,
		dlv.NROSEQ AS SEQ_VLR,
		dlv.LICIT_MTSV_NRO codreduz,
		dli.FLG_ANEXO,
		dlv.LICIT_LIC_NRO AS numlic,
		dlv.QUANT,
		dlv.VLR_UNIT,
		dlv.VLR_TOTAL,
		dlv.LICFOR_FO_PES_NRO codif, 
		FLG_MENOR_PRECO subem,
		COUNT(*) OVER (PARTITION BY dlv.LICIT_LIC_NRO, dlv.NROSEQ) AS QTD_DUPLICADAS,
		ROW_NUMBER() OVER (
			PARTITION BY dlv.LICIT_LIC_NRO, dlv.NROSEQ
			ORDER BY dlv.VLR_UNIT ASC
		) AS rn
	FROM
		system.D_LICIT_VLR dlv
	JOIN system.D_MATSERV dm ON dlv.LICIT_MTSV_NRO = dm.NRO
	JOIN system.D_LIC_ITENS dli ON dlv.LICIT_LIC_NRO = dli.LIC_NRO AND dlv.LICIT_MTSV_NRO = dli.MTSV_NRO AND dli.flg_anexo = 'N'
	JOIN system.D_LICITACAO dl ON dl.NRO = dlv.LICIT_LIC_NRO AND dl.EX_ANO >= %v
	--WHERE dlv.flg_menor_preco = 1
	)
	SELECT
		'1' sessao, 
		codif,
		SEQ_VLR item,
		ranked_items.QUANT QUAN1,
		ranked_items.VLR_UNIT VAUN1,
		ranked_items.VLR_TOTAL VATO1,
		numlic,
		'C' status,
		subem,
		NULL marca,
		'S' item_lance,
		'00000001' lotelic,
		CODREDUZ material,
		CASE WHEN ranked_items.QUANT = 1 THEN 'V' ELSE 'Q' END tpcontrole_saldo,
		1 nro_rodada
	FROM
		ranked_items
	WHERE SUBEM = 1 AND  NOT EXISTS (SELECT 1
        FROM system.D_LANCES l
        WHERE 
            l.LICIT_LIC_NRO = ranked_items.numlic
            AND l.LICIT_MTSV_NRO = ranked_items.CODREDUZ )
	ORDER BY
		numlic,
		SEQ_VLR`, modules.Cache.Ano-6)
	totalRows, err = modules.CountRows(query)
	if err != nil {
		panic(fmt.Sprintf("erro ao contar linhas: %v", err.Error()))
	}
	bar = modules.NewProgressBar(p, totalRows, "Cadpro Proposta - DISPENSAS")

	rows, err = cnxOra.Queryx(query)
	if err != nil {
		panic(fmt.Sprintf("erro ao executar query cadpro_proposta - dispensas: %v", err.Error()))
	}
	defer rows.Close()

	for rows.Next() {
		var registro ModelProposta

		if err := rows.StructScan(&registro); err != nil {
			panic(fmt.Sprintf("erro ao ler resultados da consulta cadpro_proposta - dispensas: %v", err.Error()))
		}

		registro.Cadpro = cacheCadest[registro.Codreduz.String]
		if registro.Cadpro == "" {
			panic(fmt.Sprintf("cadpro não encontrado para o código: %s", registro.Codreduz.String))
		}

		if _, err = insert.Exec(
			registro.Sessao,
			registro.Codif,
			registro.Item,
			registro.Item,
			registro.Quan1,
			registro.Vaun1,
			registro.Vato1,
			registro.Numlic,
			registro.Status,
			registro.Subem,
			registro.Marca,
			registro.ItemLance,
			registro.Lotelic,
			registro.Cadpro,
			registro.Tpcontrole,
			registro.QtdAdt,
			registro.VaunAdt,
		); err != nil {
			panic(fmt.Sprintf("erro ao inserir cadpro_proposta: %v", err.Error()))
		}
		bar.Increment()
	}	
	tx.Commit()

	cnxFdb.Exec(`insert into cadpro_lance (sessao, rodada, codif, itemp, vaunl, vatol, status, subem, numlic)
	SELECT sessao, 1 rodada, CODIF, ITEMP, VAUN1, VATO1, 'F' status, SUBEM, numlic FROM CADPRO_PROPOSTA cp where subem = 1 and not exists
	(select 1 from cadpro_lance cl where cp.codif = cl.codif and cl.itemp = cp.itemp and cl.numlic = cp.numlic)`)

	cnxFdb.Exec(`INSERT into cadpro_final (numlic, ult_sessao, codif, itemp, vaunf, vatof, STATUS, subem)
	SELECT numlic, sessao, codif, itemp, vaun1, vato1, CASE WHEN status = 'F' THEN 'C' ELSE status end, subem FROM cadpro_proposta
	WHERE NOT EXISTS (SELECT 1 FROM cadpro_final f WHERE f.numlic = cadpro_proposta.numlic AND f.itemp = cadpro_proposta.itemp AND f.codif = cadpro_proposta.codif)`)

	cnxFdb.Exec(`INSERT INTO CADPRO (
		CODIF,
		CADPRO,
		QUAN1,
		VAUN1,
		VATO1,
		SUBEM,
		STATUS,
		ITEM,
		NUMORC,
		ITEMORCPED,
		CODCCUSTO,
		FICHA,
		ELEMENTO,
		DESDOBRO,
		NUMLIC,
		ULT_SESSAO,
		ITEMP,
		QTDADT,
		QTDPED,
		VAUNADT,
		VATOADT,
		PERC,
		QTDSOL,
		ID_CADORC,
		VATOPED,
		VATOSOL,
		TPCONTROLE_SALDO,
		QTDPED_FORNECEDOR_ANT,
		VATOPED_FORNECEDOR_ANT,
		marca
	)
	SELECT
		a.CODIF,
		c.CADPRO,
		CASE WHEN a.VAUNL <> 0 THEN ROUND((a.vatol / a.VAUNL), 2) ELSE 0 END qtdunit,
		a.VAUNL,
		CASE WHEN a.VAUNL <> 0 THEN ROUND((a.vatol / a.VAUNL), 2) * a.VAUNL ELSE 0 END VATOTAL,
		1,
		'C',
		c.ITEM,
		c.NUMORC,
		c.ITEM,
		c.CODCCUSTO,
		c.FICHA,
		c.ELEMENTO,
		c.DESDOBRO,
		a.NUMLIC,
		1,
		b.ITEMP,
		CASE WHEN a.VAUNL <> 0 THEN ROUND((a.vatol / a.VAUNL), 2) ELSE 0 END qtdunit,
		0,
		a.VAUNL,
		CASE WHEN a.VAUNL <> 0 THEN ROUND((a.vatol / a.VAUNL), 2) * a.VAUNL ELSE 0 END vatoadt,
		0,
		0,
		c.ID_CADORC,
		0,
		0,
		'Q',
		0,
		0,
		p.marca
	FROM
		CADPRO_LANCE a
	INNER JOIN CADPRO_STATUS b ON
		b.NUMLIC = a.NUMLIC AND a.ITEMP = b.ITEMP AND a.SESSAO = b.SESSAO
	INNER JOIN CADPROLIC_DETALHE c ON
		c.NUMLIC = a.NUMLIC AND b.ITEMP = c.ITEM_CADPROLIC
	INNER JOIN CADLIC D ON
		D.NUMLIC = A.NUMLIC
	inner join cadpro_proposta p on 
		p.numlic = a.numlic and p.itemp = a.itemp and p.codif = a.codif
	WHERE
		a.SUBEM = 1 AND a.STATUS = 'F'
		AND NOT EXISTS (
			SELECT 1 
			FROM CADPRO cp
			WHERE cp.NUMLIC = a.NUMLIC 
			AND cp.ITEM = c.ITEM 
			AND cp.CODIF = a.CODIF
		)`)

	cnxFdb.Exec(`
	EXECUTE BLOCK AS  
		BEGIN  
		INSERT INTO REGPRECODOC (NUMLIC, CODATUALIZACAO, DTPRAZO, ULTIMA)  
		SELECT DISTINCT A.NUMLIC, 0, DATEADD(1 YEAR TO A.DTHOM), 'S'  
		FROM CADLIC A WHERE A.REGISTROPRECO = 'S' AND A.DTHOM IS NOT NULL  
		AND NOT EXISTS(SELECT 1 FROM REGPRECODOC X  
		WHERE X.NUMLIC = A.NUMLIC);  

		INSERT INTO REGPRECO (COD, DTPRAZO, NUMLIC, CODIF, CADPRO, CODCCUSTO, ITEM, CODATUALIZACAO, QUAN1, VAUN1, VATO1, QTDENT, SUBEM, STATUS, ULTIMA)  
		SELECT B.ITEM, DATEADD(1 YEAR TO A.DTHOM), B.NUMLIC, B.CODIF, B.CADPRO, B.CODCCUSTO, B.ITEM, 0, B.QUAN1, B.VAUN1, B.VATO1, 0, B.SUBEM, B.STATUS, 'S'  
		FROM CADLIC A INNER JOIN CADPRO B ON (A.NUMLIC = B.NUMLIC) WHERE A.REGISTROPRECO = 'S' AND A.DTHOM IS NOT NULL  
		AND NOT EXISTS(SELECT 1 FROM REGPRECO X  
		WHERE X.NUMLIC = B.NUMLIC AND X.CODIF = B.CODIF AND X.CADPRO = B.CADPRO AND X.CODCCUSTO = B.CODCCUSTO AND X.ITEM = B.ITEM);  

		INSERT INTO REGPRECOHIS (NUMLIC, CODIF, CADPRO, CODCCUSTO, ITEM, CODATUALIZACAO, QUAN1, VAUN1, VATO1, SUBEM, STATUS, MOTIVO, MARCA, NUMORC, ULTIMA)  
		SELECT B.NUMLIC, B.CODIF, B.CADPRO, B.CODCCUSTO, B.ITEM, 0, B.QUAN1, B.VAUN1, B.VATO1, B.SUBEM, B.STATUS, B.MOTIVO, B.MARCA, B.NUMORC, 'S'  
		FROM CADLIC A INNER JOIN CADPRO B ON (A.NUMLIC = B.NUMLIC) WHERE A.REGISTROPRECO = 'S' AND A.DTHOM IS NOT NULL  
		AND NOT EXISTS(SELECT 1 FROM REGPRECOHIS X  
		WHERE X.NUMLIC = B.NUMLIC AND X.CODIF = B.CODIF AND X.CADPRO = B.CADPRO AND X.CODCCUSTO = B.CODCCUSTO AND X.ITEM = B.ITEM);  
	END;`)
}

func Aditamento(p *mpb.Progress) {
	modules.Trigger("TBU_CADPRO", false)
	cnxFdb, cnxOra, err := connection.GetConexoes()
	if err != nil {
		panic(fmt.Sprintf("erro ao obter conexões: %v", err.Error()))
	}
	defer cnxFdb.Close()
	defer cnxOra.Close()

	tx, err := cnxFdb.Begin()
	if err != nil {
		panic(fmt.Sprintf("erro ao iniciar transação: %v", err.Error()))
	}
	defer tx.Commit()

	if _, err := tx.Exec(`update cadpro set qtdadt = quan1, vaunadt = vaun1, vatoadt = vato1 where qtdadt <> quan1`); err != nil {
		panic(fmt.Sprintf("erro ao atualizar cadpro: %v", err.Error()))
	}

	update, err := tx.Prepare(`update cadpro set qtdadt = ?+quan1, vatoadt = (?+quan1)*vaun1 where numlic = ? and item = ? and subem = 1`);
	if err != nil {
		panic(fmt.Sprintf("erro ao preparar update cadpro: %v", err.Error()))
	}
	defer update.Close()

	query := `SELECT
		ped.LIC_NRO ,
		i.MTSV_NRO ,
		COALESCE(sum(i.QUANT), 0) qtd
	FROM
		SYSTEM.D_PED_ITENS i
	JOIN SYSTEM.D_PED_LIC ped ON
		ped.PED_NRO = i.PED_NRO
	JOIN system.D_PEDIDO o ON
		o.NRO = i.PED_NRO
		AND o.FLG_TIPO = 2
	GROUP BY
		ped.LIC_NRO ,
		i.MTSV_NRO`
	
	totalRows, err := modules.CountRows(query)
	if err != nil {
		panic(fmt.Sprintf("erro ao contar linhas: %v", err.Error()))
	}
	bar := modules.NewProgressBar(p, totalRows, "Aditamento")
	
	cacheCadprolic := make(map[string]int)
	queryCadprolic, err := cnxFdb.Query("select numlic||'-'||codreduz key, item from cadprolic a join cadest b using (cadpro)")
	if err != nil {
		panic(fmt.Sprintf("erro ao consultar cadprolic: %v", err.Error()))
	}
	defer queryCadprolic.Close()
	for queryCadprolic.Next() {
		var key string
		var item int
		if err := queryCadprolic.Scan(&key, &item); err != nil {
			panic(fmt.Sprintf("erro ao ler cadprolic: %v", err.Error()))
		}
		cacheCadprolic[key] = item
	}

	rows, err := cnxOra.Query(query)
	if err != nil {
		panic(fmt.Sprintf("erro ao executar query aditamento: %v", err.Error()))
	}
	defer rows.Close()

	for rows.Next() {
		var (
			numlic, codreduz int
			qtd              float64
		)
		if err := rows.Scan(&numlic, &codreduz, &qtd); err != nil {
			panic(fmt.Sprintf("erro ao ler resultado da query aditamento: %v", err.Error()))
		}

		key := fmt.Sprintf("%d-%d", numlic, codreduz)
		if item, ok := cacheCadprolic[key]; ok {
			if _, err := update.Exec(qtd, qtd, numlic, item); err != nil {
				panic(fmt.Sprintf("erro ao executar update cadpro: %v", err.Error()))
			} 
		} else {
			panic(fmt.Sprintf("cadprolic não encontrado para a chave: %s", key))
		}
		bar.Increment()
	}
	tx.Commit()

	query = `WITH ranked_items AS (
		SELECT
			licit_lic_nro,
			LICIT_MTSV_NRO,
			dla.QUANT,
			VLR_UNIT,
			DATA_INICIO,
			ROW_NUMBER() OVER (PARTITION BY licit_lic_nro, LICIT_MTSV_NRO ORDER BY DATA_INICIO DESC) AS rn
		FROM
			system.D_LICIT_AJUST dla)
	SELECT
		licit_lic_nro,
		LICIT_MTSV_NRO,
		--quant,
		VLR_UNIT
	FROM
		ranked_items
	WHERE
		rn = 1
	ORDER BY
		LICIT_MTSV_NRO`
	
	totalRows, err = modules.CountRows(query)
	if err != nil {
		panic(fmt.Sprintf("erro ao contar linhas: %v", err.Error()))
	}

	bar = modules.NewProgressBar(p, totalRows, "Aditamento - Valores")

	rows, err = cnxOra.Query(query)
	if err != nil {
		panic(fmt.Sprintf("erro ao executar query aditamento - valores: %v", err.Error()))
	}

	for rows.Next() {
		var (
			numlic int
			codreduz int
			vlr_unit float64
		)
		if err := rows.Scan(&numlic, &codreduz, &vlr_unit); err != nil {
			panic(fmt.Sprintf("erro ao ler resultado da query aditamento - valores: %v", err.Error()))
		}

		key := fmt.Sprintf("%d-%d", numlic, codreduz)
		if item, ok := cacheCadprolic[key]; ok {
			if _, err := cnxFdb.Exec(fmt.Sprintf(`update cadpro set vaunadt = %v, vatoadt = qtdadt*%v where numlic = %v and item = %v and subem = 1`, vlr_unit, vlr_unit, numlic, item)); err != nil {
				panic(fmt.Sprintf("erro ao executar update cadpro: %v", err.Error()))
			}
		} else {
			panic(fmt.Sprintf("cadprolic não encontrado para a chave: %s", key))
		}
		bar.Increment()
	}

	modules.Trigger("TBU_CADPRO", true)
}