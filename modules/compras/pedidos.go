package compras

import (
	"GoGemmap/connection"
	"GoGemmap/modules"
	"fmt"
	"strconv"
	"strings"

	"github.com/vbauerster/mpb"
)

func Cadped(p *mpb.Progress) {
	modules.LimpaTabela([]string{"cadlic"})

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

	query := fmt.Sprintf(`select
		substr(ano,3,2) anoreduz,
		case when npedlicit is null then lpad(numped,5,'0') || '/' || substr(ano,3,2) else null end numped,
		ultimo_pedido,
		ano,
		ID_CADPED,
		sequencia,
		npedlicit,
		id_cadpedlicit,
		codif,
		codccusto,
		datped,
		ficha,
		entrou,
		numlic,
		proclic,
		localentg,
		condpgto,
		prozoentrega,
		obs,
		aditamento,
		contrato,
		item,
		material,
		qtd,
		prcunt,
		prctot,
		qtdanu,
		prctotanu,
		categoria,
		grupo,
		modalidade,
		ELEMENTO,
		DESDOBRO,
		VINGRUPO,
		vincodigo,
		destino,
		pkemp,
		empresa
		from (
				select
							row_number() over (partition by ID_CADPED order by ano, ID_CADPED, item) sequencia,
							NRO_DOC numped,
							ultimo_pedido,
							ano,
							ID_CADPED,
							case when nro_org is null then null else lpad(nro_org,5,'0') || '/' || substr(ano_org,3,2) end npedlicit,
							IDCADPED_LICIT id_cadpedlicit,
							codif,
							DT_EMISSAO datped,
							ficha,
							CENTROCUSTO codccusto,
							'N' entrou,
							numlic,
							null proclic,
							LOCAL_ENTREGA localentg,
							COND_PAGTO condpgto,
							null prozoentrega,
							obs,
							null aditamento,
							case when CONTRATO_NRO is null then null else lpad(contrato_nro,4,'0')  || '/' || substr(CONTRATO_ANO,3,2) end   contrato,
							item,
							material,
							quant qtd,
							VLR_UNIT prcunt,
							VLR_TOTAL prctot,
							qtdanu,
							prctotanu,
							categoria,
							grupo,
							modalidade,
							ELEMENTO,
							DESDOBRO,
							VINGRUPO,
							lpad(VINCODIGO,3,'0') vincodigo,
							lpad(ALMOXARIFADO,9,'0') destino,
							case when subem = 0 then ID_EMPENHO else null end pkemp,
					%v empresa
				From (
						SELECT
							A.NRO AS ID_CADPED,
							A.EX_ANO AS ANO,
							A.NRO_DOC,
							a.ATC_NRO,
							ultimo_pedido,
							A.NUMERO_SUB AS SUBEM,
							A.DT_EMISSAO,
							A.FO_PES_NRO AS CODIF,
							A.COND_PAGTO,
							A.OBS,
							A.LOCAL_ENTREGA,
							A.PLDESP_NRO ficha,
							A.FONREC_NRO AS FONTE,
							A.CDAPLVA_CDAPLFX_NRO AS VINGRUPO,
							A.CDAPLVA_NRO AS VINCODIGO,
							A.CLAED_CATED_NRO AS CATEGORIA ,
							A.CLAED_GRUPD_NRO AS GRUPO ,
							A.CLAED_MODAD_NRO AS MODALIDADE ,
							A.CLAED_ELEMD_NRO AS ELEMENTO,
							a.CLAED_NRO AS DESDOBRO,
							A.FLG_ORD_GLO_EST AS GLOBAL,
							A.LIC_NRO AS NUMLIC,
							A.DEPSEC_NRO AS CENTROCUSTO,
							A.ATC_NRO AS IDCADPED_LICIT,
							A.SECEST_NRO AS ALMOXARIFADO,
							A.EMPE_NRO AS ID_EMPENHO,
							A.CTRATO_EX_ANO AS CONTRATO_ANO,
							A.CTRATO_NRO AS CONTRATO_NRO,
							L.NRO_PROC ,
							A.EX_ANO,
							I.NROSEQ AS ITEM ,
							I.QUANT ,
							I.MTSV_NRO AS MATERIAL,
							I.VLR_UNIT,
							I.VLR_TOTAL,
							i.QUANT_ESTOR qtdanu,
							i.VALOR_ESTOR prctotanu,
							org.NRO_DOC nro_org,
							org.EX_ANO ano_org
						FROM
							SYSTEM.D_AUTCOMPR A
								join (SELECT MAX(NRO_DOC) ultimo_pedido, EX_ANO FROM SYSTEM.D_AUTCOMPR group by EX_ANO) m on m.ex_ano = a.ex_Ano
								left join (select nro, NRO_DOC, EX_ANO from system.D_AUTCOMPR where ATC_NRO is null) org on org.nro = a.ATC_NRO
								INNER JOIN SYSTEM.D_ATC_ITENS I ON
									I.ATC_NRO = A.NRO
								LEFT JOIN SYSTEM.D_LICITACAO L ON
									L.NRO = A.LIC_NRO
						WHERE
								A.EX_ANO >= %v
						ORDER BY
							ano,id_cadped, subem) pedidos) itens`, modules.Cache.Empresa, modules.Cache.Ano)
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
	
	rows, err := cnxFdb.Query(query)
	if err != nil {
		panic(fmt.Sprintf("erro ao executar query: %v", err.Error()))
	}
	defer rows.Close()

	for rows.Next() {
		var registro ModelPedidos
		if err := rows.Scan(&registro); err != nil {
			panic(fmt.Sprintf("erro ao ler registro: %v", err.Error()))
		}

		if registro.Sequencia == 1 {
			if registro.Numped.String == "" {
				ultimo, err := strconv.ParseInt(registro.UltimoPedido.String, 10, 64)
				if err != nil {
					panic(fmt.Sprintf("erro ao converter ultimo pedido: %v", err.Error()))
				}
				registro.Sequencia += 1
				novoNumero := registro.Sequencia + ultimo
				registro.Numped.String = fmt.Sprintf("%d05/%s", novoNumero, registro.Anoreduz)
			}

			num := strings.Split(registro.Numped.String, "/")[0]
			if _, err := insertCadped.Exec(
				registro.Numped,
				num,
				registro.Ano,
				registro.Datped,
				registro.Codif,
				registro.Entrou,
				registro.Ficha,
				registro.Condpgto,
				registro.Localentg,
				registro.Npedlicit,
				registro.Codccusto,
				registro.Obs,
				registro.Prozoentrega,
				registro.Contrato,
				registro.Aditamento,
				registro.Proclic,
				registro.Numlic,
				registro.IdCadped,
				registro.IdCadpedlicit,
				registro.Empresa,
			); err != nil {
				panic(fmt.Sprintf("erro ao inserir cadped: %v", err.Error()))
			}

			if registro.Pkemp.Valid {
				valor := empenhos[registro.Pkemp.Int]

				if valor != 0 {
					if _, err := insertFcadped.Exec(
						registro.Numped,
						registro.Ficha,
						registro.Pkemp,
						valor,
						registro.Categoria,
						registro.Grupo,
						registro.Modalidade,
						registro.Elemento,
						registro.Desdobro,
						0, // codfcadped
						registro.IdCadped,
					); err != nil {
						panic(fmt.Sprintf("erro ao inserir fcadped: %v", err.Error()))
					}
				}
			}
		}

		registro.Cadpro = modules.Cache.Cadpros[registro.Material.String]
		if registro.Cadpro == "" {
			panic(fmt.Sprintf("cadpro não encontrado para material: %s", registro.Material.String))
		}

		if _, err := insertIcadped.Exec(
			registro.Numped,
			registro.Item,
			registro.Cadpro,
			registro.Qtd,
			registro.Prcunt,
			registro.Prctot,
			registro.Destino,
			registro.Codccusto,
			registro.Ficha,
			registro.Categoria,
			registro.Grupo,
			registro.Modalidade,
			registro.Elemento,
			registro.Desdobro,
			registro.Vingrupo,
			registro.Vincodigo,
			registro.IdCadped,
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

	cnxFdb.Exec(`INSERT
		INTO
		regprecodoc (numlic,
		codatualizacao,
		dtprazo,
		ultima)
	SELECT
		DISTINCT a.numlic,
		0,
		dateadd(1 YEAR TO a.dthom),
		'S'
	FROM
		cadlic a
	WHERE
		a.registropreco = 'S'
		AND a.dthom IS NOT NULL
		AND NOT EXISTS(
		SELECT
			1
		FROM
			regprecodoc x
		WHERE
			x.numlic = a.numlic);

	INSERT
		INTO
		regpreco (cod,
		dtprazo,
		numlic,
		codif,
		cadpro,
		codccusto,
		item,
		codatualizacao,
		quan1,
		vaun1,
		vato1,
		qtdent,
		subem,
		status,
		ultima,
		tpcontrole_saldo)
	SELECT
		b.item,
		dateadd(1 YEAR TO a.dthom),
		b.numlic,
		b.codif,
		b.cadpro,
		b.codccusto,
		b.item,
		0,
		b.quan1,
		b.vaun1,
		b.vato1,
		0,
		b.subem,
		b.status,
		'S',
		'Q'
	FROM
		cadlic a
	INNER JOIN cadpro b ON
		(a.numlic = b.numlic)
	WHERE
		a.registropreco = 'S'
		AND a.dthom IS NOT NULL
		AND NOT EXISTS(
		SELECT
			1
		FROM
			regpreco x
		WHERE
			x.numlic = b.numlic
			AND x.codif = b.codif
			AND x.cadpro = b.cadpro
			AND x.codccusto = b.codccusto
			AND x.item = b.item);

	INSERT
		INTO
		regprecohis (
		numlic,
		codif,
		cadpro,
		codccusto,
		item,
		codatualizacao,
		quan1,
		vaun1,
		vato1,
		subem,
		status,
		motivo,
		marca,
		numorc,
		ultima)
	SELECT
		b.numlic,
		b.codif,
		b.cadpro,
		b.codccusto,
		b.item,
		0,
		b.quan1,
		b.vaun1,
		b.vato1,
		b.subem,
		b.status,
		b.motivo,
		b.marca,
		b.numorc,
		'S'
	FROM
		cadlic a
	INNER JOIN cadpro b ON
		(a.numlic = b.numlic)
	WHERE
		a.registropreco = 'S'
		AND a.dthom IS NOT NULL
		AND NOT EXISTS(
		SELECT
			1
		FROM
			regprecohis x
		WHERE
			x.numlic = b.numlic
			AND x.codif = b.codif
			AND x.cadpro = b.cadpro
			AND x.codccusto = b.codccusto
			AND x.item = b.item);

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