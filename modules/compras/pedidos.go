package compras

import (
	"GoGemmap/connection"
	"GoGemmap/modules"
	"fmt"
	"strings"

	"github.com/vbauerster/mpb"
)

func Cadped(p *mpb.Progress) {
	modules.LimpaTabela([]string{"CADPED"})

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
	defer tx.Rollback()

	insertCadped, err := tx.Prepare(`INSERT
		INTO
		cadped(numped,
		num,
		ano,
		codif,
		datped,
		ficha,
		codccusto,
		entrou,
		numlic,
		proclic,
		localentg,
		condpgto,
		prozoentrega,
		obs,
		id_cadped,
		empresa,
		aditamento,
		contrato,
		npedlicit,
		id_cadpedlicit)
	VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)
	`)
	if err != nil {
		panic(fmt.Sprintf("erro ao preparar statement: %v", err.Error()))
	}
	defer insertCadped.Close()

	insertIcadped, err := tx.Prepare(`INSERT
		INTO
		icadped(numped,
		id_cadped,
		item,
		cadpro,
		codccusto,
		qtd,
		prcunt,
		prctot,
		ficha,
		categoria,
		grupo,
		modalidade,
		elemento,
		desdobro,
		vingrupo,
		vincodigo,
		destino,
		qtdanu,
		prctotanu)
	VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`)
	if err != nil {
		panic(fmt.Sprintf("erro ao preparar statement: %v", err.Error()))
	}
	defer insertIcadped.Close()

	insertFcadped, err := tx.Prepare(`INSERT
		INTO
		fcadped(numped,
		ficha,
		valor,
		categoria,
		grupo,
		modalidade,
		elemento,
		desdobro,
		codfcadped,
		id_cadped,
		pkemp)
	VALUES(?,?,?,?,?,?,?,?,?,?,?)`)
	if err != nil {
		panic(fmt.Sprintf("erro ao preparar statement: %v", err.Error()))
	}
	defer insertFcadped.Close()

	query := fmt.Sprintf(`WITH pedidos_unicos AS (
	SELECT DISTINCT
		TO_CHAR(NRO_DOC, 'fm00000') || '/' || SUBSTR(ex_ano, 3, 2) AS numped_renumerado,
		numero_sub || CASE WHEN numero_sub < 10 THEN TO_CHAR(TO_NUMBER(a.NUMERO_DOC), 'fm0000') ELSE TO_CHAR(TO_NUMBER(a.NUMERO_DOC), 'fm00000') end || '/' || substr(ex_ano, 3, 4) numped,
		A.NRO AS id_cadped
	FROM system.D_AUTCOMPR A
	),
	sequenciado AS (
	SELECT
		numped_renumerado,
		numped,
		id_cadped,
		DENSE_RANK() OVER (PARTITION BY numped_renumerado ORDER BY id_cadped) AS sequencia
	FROM pedidos_unicos
	)
	SELECT
	SUBSTR(A.ex_ano,3,2) AS anoreduz,
	S.sequencia,
	row_number() OVER (PARTITION BY ID_CADPED ORDER BY ID_CADPED) cabecalho,
	--S.numped_renumerado,
	S.numped,
	A.NRO AS id_cadped,
	CASE WHEN A.NUMERO_SUB = 0 THEN EMPE_NRO ELSE NULL END pkemp,
	A.dt_emissao,
	A.FO_PES_NRO AS codif,
	A.COND_PAGTO,
	A.obs,
	A.local_entrega,
	A.PLDESP_NRO AS ficha,
	--A.FONREC_NRO AS fonte,
	A.CDAPLVA_CDAPLFX_NRO AS vingrupo,
	to_char(A.CDAPLVA_NRO, '000') AS vincodigo,
	A.CLAED_CATED_NRO AS categoria,
	A.CLAED_GRUPD_NRO AS grupo,
	A.CLAED_MODAD_NRO AS modalidade,
	A.CLAED_ELEMD_NRO AS elemento,
	A.CLAED_NRO AS desdobro,
	--A.FLG_ORD_GLO_EST AS "GLOBAL",
	A.LIC_NRO AS numlic,
	A.DEPSEC_NRO AS centrocusto,
	A.ATC_NRO AS idcadped_licit,
	to_char(SECEST_NRO, 'fm000000000') AS destino,
	--A.EMPE_NRO AS id_empenho,
	case when CTRATO_NRO is null then null
			else lpad(CTRATO_NRO,
			4,
			'0') || '/' || substr(A.CTRATO_EX_ANO, 3, 2)
	end contrato,
	A.EX_ANO,
	I.NROSEQ AS item,
	i.MTSV_NRO as material,
	i.quant AS qtd,
	I.VLR_UNIT as prcunt,
	I.VLR_TOTAL AS prctot,
	AN.QUANT AS qtdanu,
	AN.VALOR AS prctotanu,
	%v AS empresa,
	'S' AS entrou
	FROM
	system.D_AUTCOMPR A
	JOIN SYSTEM.D_ATC_ITENS I ON
	I.ATC_NRO = A.NRO
	JOIN sequenciado S ON
	S.id_cadped = A.NRO AND S.numped_renumerado = TO_CHAR(A.NRO_DOC, 'fm00000') || '/' || SUBSTR(A.ex_ano, 3, 2)
	LEFT JOIN D_ATCIT_ANULA AN ON
	A.NRO = AN.ATCIT_ATC_NRO AND I.NROSEQ = AN.ATCIT_NROSEQ 
	WHERE
	A.ex_ano = %v
	ORDER BY S.numped_renumerado, S.sequencia, cabecalho, item`, modules.Cache.Empresa, modules.Cache.Ano)
	totalRows, err := modules.CountRows(query)
	if err != nil {
		panic(fmt.Sprintf("erro ao contar linhas: %v", err.Error()))
	}
	bar := modules.NewProgressBar(p, totalRows, "Cadped")

	empenhos := make(map[int]float64)
	cacheEmpenhos, err := tx.Query("select PKEMP, VADEM From DESPES where ANO_RESTO is null and VADEM > 0")
	if err != nil {
		panic(fmt.Sprintf("erro ao consultar empenhos: %v", err.Error()))
	}
	defer cacheEmpenhos.Close()
	for cacheEmpenhos.Next() {
		var pkemp int
		var vadem float64
		if err := cacheEmpenhos.Scan(&pkemp, &vadem); err != nil {
			panic(fmt.Sprintf("erro ao ler empenho: %v", err.Error()))
		}
		empenhos[pkemp] = vadem
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
		panic(fmt.Sprintf("erro ao executar query: %v", err.Error()))
	}
	defer rows.Close()

	numPedidoInt := 0
	Numped := ""

	for rows.Next() {
		var registro ModelPedidos
		if err := rows.StructScan(&registro); err != nil {
			panic(fmt.Sprintf("erro ao ler registro: %v", err.Error()))
		}

		if registro.Cabecalho == 1 {
			numPedidoInt++

			Numped = fmt.Sprintf("%05d/%s", numPedidoInt, registro.Anoreduz)

			// Se for usar `num` depois
			num := strings.Split(Numped, "/")[0]
			if _, err := insertCadped.Exec(
				Numped,
				num,
				registro.Ano,
				registro.Codif,
				registro.Datped,
				registro.Ficha,
				registro.Codccusto,
				registro.Entrou,
				registro.Numlic,
				nil,
				registro.Localentg,
				registro.Condpgto,
				nil,
				registro.Obs,
				registro.IdCadped,
				registro.Empresa,
				nil,
				registro.Contrato,
				nil,
				registro.IdCadpedlicit,
			); err != nil {
				panic(fmt.Sprintf("erro ao inserir cadped: %v", err.Error()))
			}

			registro.Numped.String = Numped

			if registro.Pkemp.Valid {
				valor := empenhos[registro.Pkemp.Int]

				if valor != 0 {
					if _, err := insertFcadped.Exec(
						registro.Numped,
						registro.Ficha,
						valor,
						registro.Categoria,
						registro.Grupo,
						registro.Modalidade,
						registro.Elemento,
						registro.Desdobro,
						0, // codfcadped
						registro.IdCadped,
						registro.Pkemp,
					); err != nil {
						panic(fmt.Sprintf("erro ao inserir fcadped: %v", err.Error()))
					}
				}
			}
		}

		registro.Cadpro = cacheCadest[registro.Material.String]
		if registro.Cadpro == "" {
			panic(fmt.Sprintf("cadpro não encontrado para material: %s", registro.Material.String))
		}

		if _, err := insertIcadped.Exec(
			registro.Numped,
			registro.IdCadped,
			registro.Item,
			registro.Cadpro,
			registro.Codccusto,
			registro.Qtd,
			registro.Prcunt,
			registro.Prctot,
			registro.Ficha,
			registro.Categoria,
			registro.Grupo,
			registro.Modalidade,
			registro.Elemento,
			registro.Desdobro,
			registro.Vingrupo,
			registro.Vincodigo,
			registro.Destino,
			registro.Qtdanu,
			registro.Prctotanu,
		); err != nil {
			panic(fmt.Sprintf("erro ao inserir icadped: %v", err.Error()))
		}
		bar.Increment()
	}
	tx.Commit()

	if _, err = cnxFdb.Exec("alter table viewliq drop constraint valor__menor_anulado"); err != nil {
		panic(fmt.Sprintf("erro ao alterar tabela viewliq: %v", err.Error()))
	}

	if _, err = cnxFdb.Exec("UPDATE DESPES D SET D.ID_CADPED  = (SELECT F.ID_CADPED FROM FCADPED F WHERE F.PKEMP = D.PKEMP)"); err != nil {
		panic(fmt.Sprintf("erro ao atualizar despesas: %v", err.Error()))
	}

	if _, err = cnxFdb.Exec("UPDATE DESPES D SET D.NUMPED  = (SELECT F.NUMPED FROM FCADPED F WHERE F.PKEMP = D.PKEMP)"); err != nil {
		panic(fmt.Sprintf("erro ao atualizar despesas: %v", err.Error()))
	}

	if _, err = cnxFdb.Exec(`UPDATE cadped a SET NPEDLICIT = (SELECT numped FROM cadped x WHERE a.ID_CADPEDLICIT = x.ID_CADPED) WHERE
	a.ID_CADPEDLICIT IS NOT NULL`); err != nil {
		panic(fmt.Sprintf("erro ao atualizar despesas: %v", err.Error()))
	}

	cnxFdb.Exec(`
	UPDATE
		cadped p
	SET
		p.CODATUALIZACAO_RP = 0
	WHERE
		EXISTS(
		SELECT
			1
		FROM
			cadlic l
		WHERE
			l.numlic = p.NUMLIC
			AND l.REGISTROPRECO = 'S');`)
}

func CadproSaldoAnterior(p *mpb.Progress) {
	modules.LimpaTabela([]string{"CADPRO_SALDO_ANT","REGPRECO_SALDO_ANT"})

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

	query := fmt.Sprintf(`SELECT 
		%v ano,
		numlic,
		material,
		sum(qtd-qtdanu) qtdped,
		sum(prctot-prctotanu) vatoped
	FROM 
	(WITH pedidos_unicos AS (
		SELECT DISTINCT
			TO_CHAR(NRO_DOC, 'fm00000') || '/' || SUBSTR(ex_ano, 3, 2) AS numped,
			A.NRO AS id_cadped
		FROM system.D_AUTCOMPR A
		),
		sequenciado AS (
		SELECT
			numped,
			id_cadped,
			DENSE_RANK() OVER (PARTITION BY numped ORDER BY id_cadped) AS sequencia
		FROM pedidos_unicos
		)
		SELECT
		SUBSTR(A.ex_ano,3,2) AS anoreduz,
		S.sequencia,
		row_number() OVER (PARTITION BY ID_CADPED ORDER BY ID_CADPED) cabecalho,
		S.numped,
		A.NRO AS id_cadped,
		CASE WHEN A.NUMERO_SUB = 0 THEN EMPE_NRO ELSE NULL END pkemp,
		A.dt_emissao,
		A.FO_PES_NRO AS codif,
		A.COND_PAGTO,
		A.obs,
		A.local_entrega,
		A.PLDESP_NRO AS ficha,
		--A.FONREC_NRO AS fonte,
		A.CDAPLVA_CDAPLFX_NRO AS vingrupo,
		to_char(A.CDAPLVA_NRO, '000') AS vincodigo,
		A.CLAED_CATED_NRO AS categoria,
		A.CLAED_GRUPD_NRO AS grupo,
		A.CLAED_MODAD_NRO AS modalidade,
		A.CLAED_ELEMD_NRO AS elemento,
		A.CLAED_NRO AS desdobro,
		--A.FLG_ORD_GLO_EST AS "GLOBAL",
		A.LIC_NRO AS numlic,
		A.DEPSEC_NRO AS centrocusto,
		A.ATC_NRO AS idcadped_licit,
		to_char(SECEST_NRO, 'fm000000000') AS destino,
		--A.EMPE_NRO AS id_empenho,
		case when CTRATO_NRO is null then null
				else lpad(CTRATO_EX_ANO,
				4,
				'0') || '/' || substr(CTRATO_NRO, 3, 2)
		end contrato,
		A.EX_ANO,
		I.NROSEQ AS item,
		i.MTSV_NRO as material,
		i.quant AS qtd,
		I.VLR_UNIT as prcunt,
		I.VLR_TOTAL AS prctot,
		COALESCE(G.QUANT,0) qtdanu,
		COALESCE(G.VALOR,0) AS prctotanu,
		2 AS empresa,
		'S' AS entrou
		FROM
		system.D_AUTCOMPR A
		JOIN SYSTEM.D_ATC_ITENS I ON
		I.ATC_NRO = A.NRO
		JOIN sequenciado S ON
		S.id_cadped = A.NRO AND S.numped = TO_CHAR(A.NRO_DOC, 'fm00000') || '/' || SUBSTR(A.ex_ano, 3, 2)
		LEFT JOIN D_ATCIT_ANULA G ON a.NRO = g.ATCIT_ATC_NRO AND g.ATCIT_NROSEQ = I.nroseq
		WHERE
		A.ex_ano < %v
		ORDER BY S.numped, S.sequencia, cabecalho, item) qr WHERE idcadped_licit IS null and numlic is not null
		GROUP BY numlic, material`, modules.Cache.Ano-1, modules.Cache.Ano)
	totalRows, err := modules.CountRows(query)
	if err != nil {
		panic(fmt.Sprintf("erro ao contar linhas: %v", err.Error()))
	}
	bar := modules.NewProgressBar(p, totalRows, "CadproSaldoAnterior")

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

	cacheCadprolic := make(map[string]int)
	queryCadprolic, err := cnxFdb.Query("select numlic||'-'||cadpro key, item from cadprolic")
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

	rows, err := cnxOra.Queryx(query)
	if err != nil {
		panic(fmt.Sprintf("erro ao executar consulta: %v", err.Error()))
	}
	defer rows.Close()

	for rows.Next() {
		var registro ModelCadproSaldoAnterior
		if err := rows.StructScan(&registro); err != nil {
			panic(fmt.Sprintf("erro ao ler registro: %v", err.Error()))
		}

		registro.Codreduz = cacheCadest[registro.Codreduz]
		if registro.Codreduz == "" {
			panic(fmt.Sprintf("codreduz não encontrado para material: %s", registro.Codreduz))
		}

		registro.Item = cacheCadprolic[fmt.Sprintf("%v-%v", registro.Numlic, registro.Codreduz)]
		if registro.Item == 0 {
			//fmt.Printf("item não encontrado para numlic e codreduz: %v-%v", registro.Numlic, registro.Codreduz)
			continue
		}

		if _, err := tx.Exec(`INSERT INTO CADPRO_SALDO_ANT (ano, numlic, item, cadpro, qtdped, vatoped) VALUES (?, ?, ?, ?, ?, ?)`,
			registro.Ano,
			registro.Numlic,
			registro.Item,
			registro.Codreduz,
			registro.QtdPed,
			registro.VlrPed); err != nil {
			panic(fmt.Sprintf("erro ao inserir CADPRO_SALDO_ANT: %v", err.Error()))
		}

		bar.Increment()
	}
	tx.Commit()

	cnxFdb.Exec(`insert into regpreco_saldo_ant (ano, item, cadpro, qtdent, numlic, vatoent) 
	select ano, item, cadpro, qtdped, numlic, vatoped from cadpro_saldo_ant where numlic in (select numlic from cadlic where registropreco = 'S')`)
}