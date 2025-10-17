package compras

import (
	"GoGemmap/connection"
	"GoGemmap/modules"
	"fmt"
	"github.com/vbauerster/mpb"
)

func Cadorc(p *mpb.Progress) {
	modules.LimpaTabela([]string{"cadorc"})

	func() {
		cnxFdb, _, err := connection.GetConexoes()
		if err != nil {
			panic("Falha ao conectar com o banco de destino: " + err.Error())
		}
		defer cnxFdb.Close()

		cadpros, err := cnxFdb.Query(`select cadpro, 
			cast(codreduz as integer)||case cast(g.conv_tipo as integer) when 9 then 9 else 1 end key
			From cadest t join cadgrupo g on g.GRUPO = t.GRUPO`)
		if err != nil {
			panic("Falha ao executar consulta: " + err.Error())
		}
		defer cadpros.Close()

		modules.Cache.Cadpros = make(map[string]string)
		for cadpros.Next() {
			var cadpro, key string
			if err := cadpros.Scan(&cadpro, &key); err != nil {
				panic("Falha ao ler resultados da consulta: " + err.Error())
			}
			modules.Cache.Cadpros[key] = cadpro
		}
	}()

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
	defer tx.Commit()

	query := fmt.Sprintf(`SELECT
		ROW_NUMBER() OVER (PARTITION BY num
	ORDER BY
		num,
		item) sequencia,
		itens.*
	FROM
		(
		SELECT
			D.NRO id_cadorc,
			D.EX_ANO ano ,
			lpad(D.NRO_DOC, 5, '0') num,
			lpad(D.NRO_DOC, 5, '0') || '/' || substr(d.EX_ANO, 3, 2) numorc ,
			D.DT_EMISSAO dtorc,
			D.LIC_NRO numlic,
			D.OBJETO descr,
			d.OBJETO obs,
			0 codccusto,
			'EC' status,
			'S' liberado,
			NULL proclic,
			NULL registropreco,
			NULL ficha,
			NULL desdobro,
			I.MTSV_NRO material,
			sum(I.QUANT) qtd,
			I.MARCA,
			0 valor,
			I.NROSEQ item,
			i.NROSEQ itemorc,
			%v empresa
		FROM
			SYSTEM.D_COTACAO D
		INNER JOIN SYSTEM.D_COT_ITENS I ON
			I.COT_NRO = D.NRO
		WHERE
			EX_ANO = %v
		GROUP BY
			D.NRO,
			D.EX_ANO ,
			d.nro_doc,
			D.DT_EMISSAO,
			D.LIC_NRO ,
			D.OBJETO ,
			d.OBJETO ,
			I.MTSV_NRO,
			I.MARCA,
			i.NROSEQ
		ORDER BY
			id_cadorc,
			item) itens`, modules.Cache.Empresa, modules.Cache.Ano)
	
	insertCadorc, err := tx.Prepare(`INSERT
		INTO
		cadorc(numorc,
		num,
		ano,
		dtorc,
		descr,
		obs,
		codccusto,
		status,
		liberado,
		id_cadorc,
		empresa,
		proclic,
		registropreco,
		numlic)
	VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		panic(fmt.Sprintf("erro ao preparar insert: %v", err.Error()))
	}

	insertIcadorc, err := tx.Prepare(`INSERT
		INTO
		icadorc(numorc,
		item,
		itemorc,
		valor,
		cadpro,
		qtd,
		id_cadorc,
		ficha,
		codccusto)
	VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		panic(fmt.Sprintf("erro ao preparar insert: %v", err.Error()))
	}

	rows, err := cnxOra.Queryx(query)
	if err != nil {
		panic(fmt.Sprintf("erro ao executar query cadorc: %v", err.Error()))
	}
	defer rows.Close()

	totalRows, err := modules.CountRows(query)
	if err != nil {
		panic(fmt.Sprintf("erro ao contar linhas cadorc: %v", err.Error()))
	}
	bar := modules.NewProgressBar(p, totalRows, "Cadorc & Icadorc")

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
		var registro ModelCadorc
		if err := rows.StructScan(&registro); err != nil {
			panic(fmt.Sprintf("erro ao ler resultados da consulta: %v", err.Error()))
		}

		if registro.Sequencia == 1 {
			registro.Descr.String, err = modules.DecodeToWin1252(registro.Descr.String)
			if err != nil {
				panic(fmt.Sprintf("erro ao decodificar descrição: %v", err.Error()))
			}
			registro.Obs.String, err = modules.DecodeToWin1252(registro.Obs.String)
			if err != nil {
				panic(fmt.Sprintf("erro ao decodificar observação: %v", err.Error()))
			}

			_, err = insertCadorc.Exec(
				registro.Numorc,
				registro.Num,
				registro.Ano,
				registro.DtOrc,
				registro.Descr.String,
				registro.Obs.String,
				registro.CodCcusto,
				registro.Status,
				registro.Liberado,
				registro.IdCadorc,
				modules.Cache.Empresa,
				registro.ProcLic,
				registro.RegistroPreco,
				registro.NumLic,
			)
			if err != nil {
				panic(fmt.Sprintf("erro ao inserir cadorc: %v", err.Error()))
			}
			bar.Increment()
		}

		cadpro, ok := cacheCadest[registro.Codreduz.String]
		if !ok {
			panic(fmt.Sprintf("cadpro não encontrado para codreduz: %s", registro.Codreduz.String))
		}
		registro.Cadpro.String = cadpro

		_, err = insertIcadorc.Exec(
			registro.Numorc,
			registro.Item,
			registro.ItemOrc,
			registro.Valor,
			registro.Cadpro.String,
			registro.Qtd,
			registro.IdCadorc,
			nil,
			registro.CodCcusto,
		)
		if err != nil {
			panic(fmt.Sprintf("erro ao inserir icadorc: %v", err.Error()))
		}
	}
}

func Fcadorc(p *mpb.Progress) {
	modules.LimpaTabela([]string{"fcadorc"})

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

	query := fmt.Sprintf(`SELECT
		DISTINCT D.NRO id_cadorc,
				lpad(D.NRO_DOC,5,'0') || '/' || substr(d.EX_ANO,3,2) numorc ,
				A.nro codif,
				substr(trim(a.nome),1,70) nome,
				0 valorc
	FROM
		SYSTEM.D_COT_ITENS I
			INNER JOIN SYSTEM.D_COTACAO D ON
				I.COT_NRO = D.NRO
			INNER JOIN SYSTEM.PESSOA A ON
				A.NRO = I.FO_PES_NRO
	WHERE
			D.EX_ANO = %v`, modules.Cache.Ano)
	
	insert, err := tx.Prepare(`INSERT
		INTO
		fcadorc(id_cadorc,
		numorc,
		codif,
		nome,
		valorc)
	VALUES(?, ?, ?, ?, ?)`)
	if err != nil {
		panic(fmt.Sprintf("erro ao preparar insert: %v", err.Error()))
	}

	rows, err := cnxOra.Queryx(query)
	if err != nil {
		panic(fmt.Sprintf("erro ao executar query fcadorc: %v", err.Error()))
	}
	defer rows.Close()

	totalRows, err := modules.CountRows(query)
	if err != nil {
		panic(fmt.Sprintf("erro ao contar linhas fcadorc: %v", err.Error()))
	}
	bar := modules.NewProgressBar(p, totalRows, "Fcadorc")

	for rows.Next() {
		var registro ModelFcadorc
		if err := rows.StructScan(&registro); err != nil {
			panic(fmt.Sprintf("erro ao ler resultados da consulta: %v", err.Error()))
		}

		if len(registro.Nome) > 70 {
			registro.Nome = registro.Nome[:70]
		}

		_, err = insert.Exec(
			registro.IdCadorc,
			registro.Numorc,
			registro.Codif,
			registro.Nome,
			registro.Valorc,
		)
		if err != nil {
			panic(fmt.Sprintf("erro ao inserir fcadorc: %v", err.Error()))
		}
		bar.Increment()
	}
}

func Vcadorc(p *mpb.Progress) {
	modules.LimpaTabela([]string{"vcadorc"})

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

	query := fmt.Sprintf(`SELECT
				lpad(D.NRO_DOC,5,'0') || '/' || substr(d.EX_ANO,3,2) numorc ,
				D.NRO AS id_cadorc,
				i.FO_PES_NRO codif,
				i.nroseq AS item,
				i.valor vlruni,
				i.quant * i.valor AS vlrtot,
				'GL' classe,
				MARCA,
				gn.ganhou,
				gn.vlrganhou
	FROM
		SYSTEM.D_COT_ITENS I
			INNER JOIN SYSTEM.D_COTACAO D ON
				I.COT_NRO = D.NRO
			left join (select g.COT_NRO ,nroseq , FO_PES_NRO ganhou, valor vlrganhou, row_number() OVER (PARTITION BY COT_NRO, nroseq ORDER BY valor) seq
					from SYSTEM.D_COT_ITENS g where FLG_MENOR_PRECO = 1) gn on gn.COT_NRO = i.COT_NRO
			and gn.nroseq = i.NROSEQ AND i.FO_PES_NRO = ganhou
	WHERE
			D.EX_ANO = %v`, modules.Cache.Ano)

	insert, err := tx.Prepare(`INSERT
		INTO
		vcadorc(numorc,
		id_cadorc,
		codif,
		item,
		vlruni,
		ganhou,
		vlrtot,
		classe,
		marca,
		vlrganhou)
	VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		panic(fmt.Sprintf("erro ao preparar insert: %v", err.Error()))
	}

	rows, err := cnxOra.Queryx(query)
	if err != nil {
		panic(fmt.Sprintf("erro ao executar query vcadorc: %v", err.Error()))
	}
	defer rows.Close()

	totalRows, err := modules.CountRows(query)
	if err != nil {
		panic(fmt.Sprintf("erro ao contar linhas vcadorc: %v", err.Error()))
	}
	bar := modules.NewProgressBar(p, totalRows, "Vcadorc")
	
	for rows.Next() {
		var registro ModelVcadorc
		if err := rows.StructScan(&registro); err != nil {
			panic(fmt.Sprintf("erro ao ler resultados da consulta: %v", err.Error()))
		}
		_, err = insert.Exec(
			registro.Numorc,
			registro.IdCadorc,
			registro.Codif,
			registro.Item,
			registro.VlrUni,
			registro.Ganhou,
			registro.VlrTot,
			registro.Classe,
			registro.Marca,
			registro.VlrGanhou,
		)
		if err != nil {
			panic(fmt.Sprintf("erro ao inserir vcadorc: %v", err.Error()))
		}
		bar.Increment()
	}
	tx.Commit()

	if _, err := cnxFdb.Exec(`update fcadorc f set f.valorc = (select sum(v.vlrtot) from vcadorc v 
                where f.codif = v.codif and f.id_cadorc = v.id_cadorc)`); err != nil {
		panic(fmt.Sprintf("erro ao atualizar fcadorc: %v", err.Error()))
	}

	_, err = cnxFdb.Exec(`MERGE INTO icadorc_cot c
	USING (
	SELECT
		FIRST 1 item,
		id_cadorc,
		ganhou,
		vlrganhou
	FROM vcadorc
	GROUP BY item, id_cadorc, ganhou, vlrganhou
	) v
	ON c.item = v.item AND c.id_cadorc = v.id_cadorc
	WHEN MATCHED THEN
	UPDATE SET
		c.codif = v.ganhou,
		c.valunt = v.vlrganhou`)
	if err != nil {
		panic(fmt.Sprintf("erro ao atualizar icadorc_cot: %v", err.Error()))
	}

	_, err = cnxFdb.Exec(`UPDATE icadorc_cot set valtot = valunt * qtd, tipo = 'M', flg_aceito = 'S'`)
	if err != nil {
		panic(fmt.Sprintf("erro ao atualizar icadorc_cot: %v", err.Error()))
	}

	_, err = cnxFdb.Exec(`UPDATE icadorc_cot set valtot = valunt * qtd, valunt = 0 where valunt is null or valtot is null`)
	if err != nil {
		panic(fmt.Sprintf("erro ao atualizar icadorc_cot: %v", err.Error()))
	}
}