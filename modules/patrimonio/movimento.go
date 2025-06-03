package patrimonio

import (
	"GoGemmap/connection"
	"GoGemmap/modules"
	"fmt"

	"github.com/vbauerster/mpb"
)

func Movbem(p *mpb.Progress) {
	modules.LimpaTabela([]string{"pt_movbem"})

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

	query := fmt.Sprintf(`SELECT
		row_number() OVER (ORDER BY data_mov) codigo_mov,
		%v empresa_mov,
		qr.*
	FROM
		(
		SELECT
			a.BEPA_NRO codigo_pat_mov,
			a.DATA_INICIO data_mov,
			a.FLG_ADQTRANSF tipo_mov,
			CASE WHEN flg_adqtransf = 'T' THEN 0 ELSE b.valor_mov END valor_mov,
			a.DEPSEC_NRO codigo_set_mov,
			CASE WHEN a.FLG_ADQTRANSF = 'T' THEN 'TRANSFERENCIA' ELSE 'AQUISICAO' end historico_mov,
			'N' depreciacao_mov
		FROM
			system.D3_BP_DS a
		JOIN (
			SELECT
				bepa_nro,
				max(dbv.VLR_BASE_CALC) valor_mov
			FROM
				system.D3_BP_VALOR dbv
			WHERE
				flg_reavaliado = 'N'
			GROUP BY
				bepa_nro,
				DATA_AVAL
			ORDER BY
				dbv.DATA_AVAL ASC) b ON
			a.BEPA_NRO = b.BEPA_NRO
	UNION ALL
		SELECT codigo_pat_mov, data_mov, tipo_mov, valor_mov, codigo_set_mov, historico_mov, rn.DEPRECIACAO_MOV FROM (
		SELECT
			ddc.BPVLR_BEPA_NRO codigo_pat_mov,
			ddc.BPVLR_DATA_AVAL data_mov,
			'R' tipo_mov,
			-valor valor_mov,
			NULL codigo_set_mov,
			'DEPRECIACAO ' || mes || '-' || ano historico_mov,
			'S' DEPRECIACAO_MOV
		FROM
			system.D3_DEPR_CORR ddc
		ORDER BY ano, mes) rn
	UNION ALL
		SELECT
			bepa_nro,
			MAX(DATA_AVAL) AS data_reavaliacao,
			'R',
			MAX(VLR_BASE_CALC) - MIN(VLR_BASE_CALC) AS valor_mov,
			NULL,
			'REAVALIACAO',
			'N' DEPRECIACAO_MOV
		FROM
			system.D3_BP_VALOR
		GROUP BY
			bepa_nro
		HAVING
			COUNT(*) > 1
		ORDER BY
			data_mov)  QR`, modules.Cache.Empresa)
	
	totalRows, _ := modules.CountRows(query)
	barMovbem := modules.NewProgressBar(p, totalRows, "MOVBEM")
	if err != nil {
		panic(fmt.Sprintf("Erro ao contar linhas: %v", err))
	}
	
	insert, err := tx.Prepare(`insert
		into
		pt_movbem (codigo_mov,
		empresa_mov,
		codigo_pat_mov,
		data_mov,
		tipo_mov,
		codigo_set_mov,
		historico_mov,
		valor_mov,
		depreciacao_mov,
		dt_contabil)
	values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		panic(fmt.Sprintf("Erro ao preparar insert: %v", err))
	}

	rows, err := cnxOra.Queryx(query)
	if err != nil {
		panic(fmt.Sprintf("Erro ao executar query: %v", err))
	}

	for rows.Next() {
		var registro ModelMovbem
		if err := rows.StructScan(&registro); err != nil {
			panic(fmt.Sprintf("Erro ao escanear valores: %v", err))
		}

		if _, err := insert.Exec(
			registro.CodigoMov,
			registro.EmpresaMov,
			registro.CodigoPatMov,
			registro.DataMov,
			registro.TipoMov,
			registro.CodigoSetMov,
			registro.HistoricoMov,
			registro.ValorMov,
			registro.DepreciacaoMov,
			registro.DtContabil,
		); err != nil {
			panic(fmt.Sprintf("Erro ao executar insert: %v", err))
		}

		barMovbem.Increment()
	}
	tx.Commit()

	if _, err = cnxFdb.Exec(`MERGE INTO pt_cadpat a USING (SELECT codigo_pat_mov, MAX(data_mov) AS ultima_data
		FROM pt_movbem
		WHERE tipo_mov = 'R' AND depreciacao_mov <> 'S'
		GROUP BY codigo_pat_mov) b 
	ON b.codigo_pat_mov = a.codigo_pat
	WHEN MATCHED THEN UPDATE SET a.dtlan_pat = b.ultima_data, a.dt_contabil = b.ultima_data
	`); err != nil {
		panic(fmt.Sprintf("Erro ao executar merge: %v", err))
	}
}