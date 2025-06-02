package compras

import (
	"GoGemmap/connection"
	"GoGemmap/modules"
	"fmt"

	"github.com/vbauerster/mpb"
)

func SaldoInicial(p *mpb.Progress) {
	modules.Trigger("TD_ICADREQ", false)
	modules.LimpaTabela([]string{"icadreq where requi containing '000000/'", "requi where requi containing '000000/'"})
	modules.Trigger("TI_ICADREQ", false)

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

	insertIcadreq, err := tx.Prepare(`insert into icadreq (id_requi, requi, codccusto, empresa, item, quan1, quan2, vaun1, vaun2, vato1, vato2, cadpro, destino) values (?,?,?,?,?,?,?,?,?,?,?,?,?)`)
	if err != nil {
		panic("Erro ao preparar insert: " + err.Error())
	}
	defer insertIcadreq.Close()

	query := fmt.Sprintf(`SELECT
		0 as id_requi,
		'000000/%v' AS requi,
		0 codccusto,
		%v AS empresa,
		ROW_NUMBER() OVER (ORDER BY mtsv_nro) AS item,
		TO_CHAR(secest_nro, 'fm00000000') AS destino,
		mtsv_nro AS codreduz,
		saldo AS quan1,
		0 quan2,
		preco_medio AS vaun1,
		0 vaun2,
		saldo * preco_medio AS vato1,
		0 vato2
	FROM
		system.D_INVENTARIO di
	WHERE
		ex_ano = %v
		AND di.mes = 12
		AND saldo * preco_medio <> 0`, modules.Cache.Ano%2000, modules.Cache.Empresa, modules.Cache.Ano-1)

	totalLinhas, err := modules.CountRows(query)
	if err != nil {
		panic("Erro ao contar linhas: " + err.Error())
	}
	bar := modules.NewProgressBar(p, totalLinhas, "Saldo Inicial")

	_, err = tx.Exec(`INSERT
	INTO
	requi (empresa,
	id_requi,
	requi,
	num,
	ano,
	destino,
	codccusto,
	datae,
	dtlan,
	entr,
	said,
	comp,
	codif)
	VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)`,
		modules.Cache.Empresa,
		0,
		fmt.Sprintf("000000/%v", modules.Cache.Ano%2000),
		"000000",
		modules.Cache.Ano,
		"000000000",
		0,
		fmt.Sprintf("31.12.%d", modules.Cache.Ano-1),
		fmt.Sprintf("31.12.%d", modules.Cache.Ano-1),
		"S",
		"N",
		"P",
		nil)
	if err != nil {
		panic("Erro ao inserir na tabela requi: " + err.Error())
	}

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
		panic("Erro ao executar consulta: " + err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var registro ModelIcadreq
		if err := rows.StructScan(&registro); err != nil {
			panic("Erro ao ler registro: " + err.Error())
		}

		registro.Cadpro = cacheCadest[registro.Cadpro]

		if _, err = insertIcadreq.Exec(registro.Id_requi, registro.Requi, registro.Codccusto, registro.Empresa, registro.Item, registro.Quan1, registro.Quan2, registro.Vaun1, registro.Vaun2, registro.Vato1, registro.Vato2, registro.Cadpro, registro.Destino); err != nil {
			panic("Erro ao inserir icadreq: " + err.Error())
		}
		bar.Increment()
	}
	modules.Trigger("TD_ICADREQ", true)
	modules.Trigger("TI_ICADREQ", true)
}

func Requi(p *mpb.Progress) {
	modules.Trigger("TD_ICADREQ", false)
	modules.Trigger("TAU_ESTOQUE_DESTINO", false)
	modules.Trigger("TI_ICADREQ", false)
	modules.LimpaTabela([]string{"icadreq where requi not containing '000000/'", "requi where requi not containing '000000/'"})

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

	insertRequi, err := tx.Prepare(`INSERT
		INTO
		requi (empresa,
		id_requi,
		requi,
		num,
		ano,
		destino,
		codccusto,
		datae,
		dtlan,
		dtpag,
		entr,
		said,
		comp,
		codif,
		entr_said)
	VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`)
	if err != nil {
		panic("Erro ao preparar insert: " + err.Error())
	}
	defer insertRequi.Close()

	insertIcadreq, err := tx.Prepare(`insert into icadreq (id_requi, requi, codccusto, empresa, item, quan1, quan2, vaun1, vaun2, vato1, vato2, cadpro, destino) values (?,?,?,?,?,?,?,?,?,?,?,?,?)`)
	if err != nil {
		panic("Erro ao preparar insert: " + err.Error())
	}
	defer insertIcadreq.Close()

	query := fmt.Sprintf(`
	WITH base AS (
	select
			%v entidade,
			m.nro id_requi,
			lpad(m.nro,6,'0') || '/' || substr(EX_ANO,3,2) requi,
			lpad(m.nro,6,'0') num,
			ex_ano ano,
			TO_CHAR(secest_nro, 'fm00000000') AS destino, 
			0 codccusto,
			dt_emissao dtlan,
			dt_emissao datae,
			dt_emissao dtpag,
			'X' tipomov,
			'P' comp,
			m.FO_PES_NRO codif,
			nro_doc||'/'||serie docum,
			i.NROSEQ item,
			i.MTSV_NRO codreduz,
			i.quant quantidade,
			i.VAL_UNIT valorunitario,
			i.VAL_TOTAL valortotal,
			--NULL motorista,
			--NULL km,
			--NULL placa,
			(select case when aut.NRO_DOC is not null then lpad(aut.NRO_DOC,5,'0') || '/' || substr(aut.EX_ANO,3,2) else null end numped
				from system.D_MOV_ATC ped
				join system.D_AUTCOMPR aut on ped.ATC_NRO = aut.NRO and aut.ATC_NRO is null
			where  m.NRO = ped.MOV_NRO and rownum = 1
			) numped
		from system.D_MOVTO m
				join system.D_MOV_ITENS i on i.MOV_NRO = m.NRO
		where extract(year from m.DT_EMISSAO) = %v AND i.FLG_SAI_DIR = 'S'
		UNION all
		--Entradas
		SELECT
			%v entidade,
			a.nro id_requi,
			lpad(a.nro,6,'0') || '/' || substr(EX_ANO,3,2) requi,
			lpad(a.nro,6,'0') num,
			ex_ano ano,
			TO_CHAR(secest_nro, 'fm00000000') AS destino, 
			0 codccusto,
			dt_emissao dtlan,
			dt_emissao datae,
			null dtpag,
			'E' tipomov,
			'P' comp,
			a.FO_PES_NRO codif,
			nro_doc||'/'||serie docum,
			b.NROSEQ item,
			b.MTSV_NRO codreduz,
			b.quant quantidade,
			b.VAL_UNIT vaun1,
			b.VAL_TOTAL vato1,
			--NULL motorista,
			--NULL km,
			--NULL placa,
			(select case when aut.NRO_DOC is not null then lpad(aut.NRO_DOC,5,'0') || '/' || substr(aut.EX_ANO,3,2) else null end numped
				from system.D_MOV_ATC ped
				join system.D_AUTCOMPR aut on ped.ATC_NRO = aut.NRO and aut.ATC_NRO is null
			where  a.NRO = ped.MOV_NRO and rownum = 1
			) numped
		FROM
			system.D_MOVTO a
		JOIN system.D_MOV_ITENS b ON a.nro = b.MOV_NRO AND b.FLG_SAI_DIR <> 'S' AND a.EX_ANO = %v
		UNION all
		--Saidas
		SELECT
			%v entidade,
			dri.REQ_NRO id_requi,
			lpad(dr.nro,6,'0') || '/' || substr(EX_ANO,3,2) requi,
			lpad(dr.nro,6,'0') num,
			ex_ano ano,
			TO_CHAR(secest_nro, 'fm00000000') AS destino,
			DEPSEC_NRO codccusto,
			dt_emissao dtlan,
			NULL datae,
			dt_emissao dtpag,
			'S' tipomov,
			'P' comp,
			null codif,
			null docum,
			dri.nroseq item,
			dri.MTSV_NRO,
			dri.quant,
			0,
			0,
			--NULL motorista,
			--NULL km,
			--NULL placa,
			NULL numped
		FROM
			system.D_REQUISICAO dr
		JOIN system.D_REQ_ITENS dri ON
			dr.nro = dri.req_nro
		WHERE ex_ano = %v AND dr.FLG_SAI_DIR <> 'S')
	SELECT
	entidade,
	DENSE_RANK() OVER (ORDER BY num, tipomov) AS id_requi,
	ano,
	destino,
	codccusto,
	dtlan,
	datae,
	dtpag,
	tipomov,
	comp,
	codif,
	docum,
	item,
	CODREDUZ,
	QUANTIDADE,
	VALORUNITARIO,
	VALORTOTAL
	FROM base`, modules.Cache.Empresa, modules.Cache.Ano, modules.Cache.Empresa, modules.Cache.Ano, modules.Cache.Empresa, modules.Cache.Ano)

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

	totalLinhas, err := modules.CountRows(query)
	if err != nil {
		panic("Erro ao contar linhas: " + err.Error())
	}
	bar := modules.NewProgressBar(p, totalLinhas, "Requisição")

	idRequiAnterior := 0

	rows, err := cnxOra.Queryx(query)
	if err != nil {
		panic("Erro ao executar consulta: " + err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var (
			registro ModelRequi
			entrSaid = "N"
		)

		if err := rows.StructScan(&registro); err != nil {
			panic("Erro ao ler registro: " + err.Error())
		}

		registro.Requi = fmt.Sprintf("%06d/%v", registro.Id_requi, registro.Ano%1000)
		registro.Num = fmt.Sprintf("%06d", registro.Id_requi)

		if registro.Id_requi != idRequiAnterior {
			if registro.Tipo == "E" {
				registro.Entr = "S"
				registro.Said = "N"
			} else if registro.Tipo == "S" {
				registro.Entr = "N"
				registro.Said = "S"
			} else {
				registro.Entr = "S"
				registro.Said = "S"
				entrSaid = "S"
			}
			_, err = insertRequi.Exec(registro.Empresa, registro.Id_requi, registro.Requi, registro.Num, registro.Ano, registro.Destino, registro.Codccusto, registro.Datae, registro.Dtlan, registro.Dtpag, registro.Entr, registro.Said, registro.Comp, registro.Codif, entrSaid)
			if err != nil {
				panic("Erro ao inserir na tabela requi: " + err.Error())
			}
			idRequiAnterior = registro.Id_requi
		}

		registro.Cadpro = cacheCadest[registro.Cadpro]
		if registro.Tipo == "E" {
			registro.Quan1 = registro.Quantidade
			registro.Vaun1 = registro.ValorUnit
			registro.Vato1 = registro.ValorTotal

			registro.Quan2 = 0
			registro.Vaun2 = 0
			registro.Vato2 = 0
		} else if registro.Tipo == "S" {
			registro.Quan1 = 0
			registro.Vaun1 = 0
			registro.Vato1 = 0

			registro.Quan2 = registro.Quantidade
			registro.Vaun2 = registro.ValorUnit
			registro.Vato2 = registro.ValorTotal
		} else {
			registro.Quan1 = registro.Quantidade
			registro.Vaun1 = registro.ValorUnit
			registro.Vato1 = registro.ValorTotal

			registro.Quan2 = registro.Quantidade
			registro.Vaun2 = registro.ValorUnit
			registro.Vato2 = registro.ValorTotal
		}
		_, err = insertIcadreq.Exec(registro.Id_requi, registro.Requi, registro.Codccusto, registro.Empresa, registro.Item, registro.Quan1, registro.Quan2, registro.Vaun1, registro.Vaun2, registro.Vato1, registro.Vato2, registro.Cadpro, registro.Destino)
		if err != nil {
			panic("Erro ao inserir icadreq: " + err.Error())
		}
		bar.Increment()
	}
	tx.Commit()

	modules.Trigger("TD_ICADREQ", true)
	modules.Trigger("TAU_ESTOQUE_DESTINO", true)
	modules.Trigger("TI_ICADREQ", true)
}
