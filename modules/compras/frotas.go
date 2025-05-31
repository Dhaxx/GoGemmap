package compras

import (
	"GoGemmap/connection"
	"GoGemmap/modules"
	"fmt"
	"github.com/vbauerster/mpb"
)

func Motor(p *mpb.Progress) {
	modules.LimpaTabela([]string{"motor"})

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

	insert, err := tx.Prepare(`INSERT
		INTO
		motor(cod,
		nome,
		cnh,
		categcnh,
		dtvenccnh)
	VALUES(?, ?, ?, ?, ?)`)
	if err != nil {
		panic(fmt.Sprintf("erro ao preparar insert: %v", err.Error()))
	}
	defer insert.Close()

	query := `SELECT
		e.PES_NRO cod,
		p.nome nome,
		CNH_NRO cnh,
		CNH_CATEG categcnh,
		CNH_VENCTO dtvenccnh
	FROM
		system.E_MOTORISTA e
	INNER JOIN system.PESSOA p ON
		p.nro = e.PES_NRO`

	totalLinhas, err := modules.CountRows(query)
	if err != nil {
		panic(fmt.Sprintf("erro ao contar linhas: %v", err.Error()))
	}
	bar := modules.NewProgressBar(p, totalLinhas, "Motoristas")

	rows, err := cnxOra.Queryx(query)
	if err != nil {
		panic(fmt.Sprintf("erro ao executar consulta: %v", err.Error()))
	}
	defer rows.Close()

	for rows.Next() {
		var registro ModelMotor

		if err := rows.StructScan(&registro); err != nil {
			panic(fmt.Sprintf("erro ao escanear registro: %v", err.Error()))
		}

		if registro.Nome, err = modules.DecodeToWin1252(registro.Nome); err != nil {
			panic(fmt.Sprintf("erro ao decodificar nome: %v", err.Error()))
		}

		if _, err := insert.Exec(
			registro.Cod,
			registro.Nome,
			registro.Cnh,
			registro.Categcnh,
			registro.Dtvenccnh,
		); err != nil {
			panic(fmt.Sprintf("erro ao inserir registro: %v", err.Error()))
		}
		bar.Increment()
	}
}

func VeiculoTipo(p *mpb.Progress) {
	modules.LimpaTabela([]string{"veiculotipo"})

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

	query := `SELECT
		'insert into veiculo_tipo (codigo_tip, descricao_tip) values ('||nro||', '''||nome||''');'
	FROM
		system.E_tipo_veic`

	rows, err := cnxOra.Queryx(query)
	if err != nil {
		panic(fmt.Sprintf("erro ao executar consulta: %v", err.Error()))
	}
	defer rows.Close()

	for rows.Next() {
		var query string
		if err := rows.Scan(&query); err != nil {
			panic(fmt.Sprintf("erro ao escanear registro: %v", err.Error()))
		}
		if _, err := tx.Exec(query); err != nil {
			panic(fmt.Sprintf("erro ao executar insert: %v", err.Error()))
		}
	}
}

func VeiculoMarca(p *mpb.Progress) {
	modules.LimpaTabela([]string{"veiculomarca"})

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

	query := `SELECT
		'insert into veiculo_marca(codigo_mar,descricao_mar,codigo_tip_mar) values (' || codigo_mar || ', ''' || descricao_mar || ''', ' || codigo_tip_mar || ');'
	FROM
		(
		SELECT
			nro codigo_mar,
			nome descricao_mar,
			(
			SELECT
				v.etipveic_nro
			FROM
				system.e_veiculo v
			WHERE
				v.EMARCVE_NRO = e.nro
				AND ROWNUM = 1) codigo_tip_mar
		FROM
			system.e_marca_veic e) qr`

	rows, err := cnxOra.Queryx(query)
	if err != nil {
		panic(fmt.Sprintf("erro ao executar consulta: %v", err.Error()))
	}
	defer rows.Close()

	for rows.Next() {
		var query string
		if err := rows.Scan(&query); err != nil {
			panic(fmt.Sprintf("erro ao escanear registro: %v", err.Error()))
		}
		if _, err := tx.Exec(query); err != nil {
			panic(fmt.Sprintf("erro ao executar insert: %v", err.Error()))
		}
	}
}

func Veiculo(p *mpb.Progress) {
	modules.LimpaTabela([]string{"veiculo"})

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

	insert, err := tx.Prepare(`INSERT
		INTO
		veiculo(placa,
		sequencia,
		modelo,
		chassi,
		cor,
		ano,
		anomod,
		renavam,
		aquisicao,
		motorista,
		codigo_marca_vei,
		kminicial,
		obs,
		combustivel,
		alienacao,
		licenca,
		trocaoleo)
	VALUES(?,?,?,?,?,?,?,?,
	?,?,?,?,?,?,?,?,?)`)
	if err != nil {
		panic(fmt.Sprintf("erro ao preparar insert: %v", err.Error()))
	}
	defer insert.Close()

	query := `select
			PLACA_LETRAS || PLACA_NUMEROS placa,
			v.nro sequencia,
			substr(v.nome,1,45) modelo,
			nro_chassi chassi,
			cor,
			ano_fabr ano,
			ano_mod anomod,
			nro_renavam renavam,
			DT_AQUISICAO aquisicao,
			EMOT_PES_NRO motorista,
			EMARCVE_NRO codigo_marca_vei,
			km kminicial,
			obs,
			substr(c.nome,1,1) combustivel,
			DT_VENDA alienacao,
			VENCTO_LICENC licenca,
			TROLEO_KM trocaoleo
	from system.E_VEICULO V
			left join system.E_TIPO_COMB c on c.nro = v.ETPCOMB_NRO`

	totalLinhas, err := modules.CountRows(query)
	if err != nil {	
		panic(fmt.Sprintf("erro ao contar linhas: %v", err.Error()))
	}

	bar := modules.NewProgressBar(p, totalLinhas, "Veículos")
	rows, err := cnxOra.Queryx(query)
	if err != nil {
		panic(fmt.Sprintf("erro ao executar consulta: %v", err.Error()))
	}
	defer rows.Close()

	for rows.Next() {
		var registro ModelVeiculo

		if err := rows.StructScan(&registro); err != nil {
			panic(fmt.Sprintf("erro ao escanear registro: %v", err.Error()))
		}

		if registro.Modelo, err = modules.DecodeToWin1252(registro.Modelo); err != nil {
			panic(fmt.Sprintf("erro ao decodificar modelo: %v", err.Error()))
		}

		if registro.Cor, err = modules.DecodeToWin1252(registro.Cor); err != nil {
			panic(fmt.Sprintf("erro ao decodificar cor: %v", err.Error()))
		}

		if _, err := insert.Exec(
			registro.Placa,
			registro.Sequencia,
			registro.Modelo,
			registro.Chassi,
			registro.Cor,
			registro.Ano,
			registro.Anomod,
			registro.Renavam,
			registro.Aquisicao,
			registro.Motorista,
			registro.CodigoMarcaVeiculo,
			registro.Kminicial,
			registro.Obs,
			registro.Combustivel,
			registro.Alienacao,
			registro.Licenca,
			registro.Trocaoleo,
		); err != nil {
			panic(fmt.Sprintf("erro ao inserir registro: %v", err.Error()))
		}
		bar.Increment()
	}
}

func Abastecimento(p *mpb.Progress) {
	modules.LimpaTabela([]string{"abastecimento"})

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
		entr,
		said,
		comp,
		codif,
		entr_said)
	VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)`)
	if err != nil {
		panic("Erro ao preparar insert: " + err.Error())
	}
	defer insertRequi.Close()

	insertIcadreq, err := tx.Prepare(`INSERT
		INTO
		icadreq (id_requi,
		requi,
		codccusto,
		empresa,
		item,
		quan1,
		quan2,
		vaun1,
		vaun2,
		vato1,
		vato2,
		cadpro,
		destino,
		km,
		placa)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		panic("Erro ao preparar insert: " + err.Error())
	}
	defer insertIcadreq.Close()

	query := fmt.Sprintf(`select
		%v entidade,
		r.NRO id_requi,
		lpad(r.nro,6,'0') || '/' || substr(extract(year from dt_emissao),3,2) requi,
		lpad(r.nro,6,'0') num,
		extract(year from dt_emissao) ano,
		r.DT_EMISSAO dtlan,
		r.DT_EMISSAO datae,
		r.DT_EMISSAO dtpag,
		r.FO_PES_NRO codif,
		r.NUMERO_NF docum,
		'X' tipomov,
		'P' comp,
		r.EMOT_PES_NRO motorista,
		i.NROSEQ item,
		r.KM_ATUAL km,
		PLACA_LETRAS || PLACA_NUMEROS placa,
		i.EPRODSERV_NRO codreduz,
		i.QUANT quan1,
		i.QUANT quan2,
		i.VLR_UNIT vaun1,
		i.VLR_UNIT vaun2,
		i.VLR_TOTAL vato1,
		i.VLR_TOTAL vato2
	from system.E_REQUISICAO r
			join system.E_REQ_ITENS i on i.EREQ_NRO = r.NRO
			join system.E_VEICULO v on v.NRO = r.EVEIC_NRO
	where extract(year from DT_EMISSAO) = %v
	order by DT_EMISSAO, r.nro`, modules.Cache.Empresa, modules.Cache.Ano)

	cacheCentrocusto := make(map[string][]string)
	queryCentroCusto := `SELECT placa, codccusto, destino from centrocusto where placa is not null`
	rowsCentroCusto, err := cnxFdb.Query(queryCentroCusto)
	if err != nil {
		panic(fmt.Sprintf("erro ao executar consulta de centro de custo: %v", err.Error()))
	}
	defer rowsCentroCusto.Close()

	for rowsCentroCusto.Next() {
		var placa, codccusto, destino string
		if err := rowsCentroCusto.Scan(&placa, &codccusto, &destino); err != nil {
			panic(fmt.Sprintf("erro ao escanear registro de centro de custo: %v", err.Error()))
		}
		cacheCentrocusto[placa] = []string{codccusto, destino}
	}

	rows, err := cnxOra.Queryx(query)
	if err != nil {
		panic(fmt.Sprintf("erro ao executar consulta: %v", err.Error()))
	}
	defer rows.Close()
	
	for rows.Next() {
		var registro ModelAbastecimento
		var entrSaid = "S"
		if err := rows.StructScan(&registro); err != nil {
			panic(fmt.Sprintf("erro ao escanear registro: %v", err.Error()))
		}

		if registro.Cadpro, err = modules.DecodeToWin1252(registro.Cadpro); err != nil {
			panic(fmt.Sprintf("erro ao decodificar cadpro: %v", err.Error()))
		}

		CentroCustoInfo := cacheCentrocusto[registro.Placa.String]
		if CentroCustoInfo == nil {
			panic(fmt.Sprintf("centro de custo não encontrado para placa: %s", registro.Placa.String))
		}
		registro.Codccusto = CentroCustoInfo[0]
		registro.Destino = CentroCustoInfo[1]

		if _, err := insertRequi.Exec(
			modules.Cache.Empresa,
			registro.Id_requi,
			registro.Requi,
			registro.Num,
			registro.Ano,
			registro.Destino,
			registro.Codccusto,
			registro.Datae,
			registro.Dtlan,
			registro.Entr,
			registro.Said,
			registro.Comp,
			registro.Codif,
			entrSaid,
		); err != nil {
			panic(fmt.Sprintf("erro ao inserir requisição: %v", err.Error()))
		}

		if _, err := insertIcadreq.Exec(
			registro.Id_requi,
			registro.Requi,
			registro.Codccusto,
			registro.Empresa,
			registro.Item,
			registro.Quan1,
			registro.Quan2,
			registro.Vaun1,
			registro.Vaun2,
			registro.Vato1,
			registro.Vato2,
			registro.Cadpro,
			registro.Destino,
			registro.Km,
			registro.Placa,
		); err != nil {
			panic(fmt.Sprintf("erro ao inserir item de requisição: %v", err.Error()))
		}
	}
}