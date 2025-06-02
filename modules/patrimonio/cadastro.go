package patrimonio

import (
	"GoGemmap/connection"
	"GoGemmap/modules"
	"fmt"

	"github.com/vbauerster/mpb"
)

func Cadpat(p *mpb.Progress) {
	modules.LimpaTabela([]string{"pt_cadpat"})

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

	query := fmt.Sprintf(` SELECT
        DISTINCT(P.NRO) codigo_pat,
		%v empresa_pat,
        '001' codigo_gru_pat,
        to_char(ROW_number() OVER (ORDER BY p.nro),'fm000000') chapa_pat,
        E.PLCTRZ_CDCONT_NRO cpl,
        P1.DEPSEC_NRO codigo_set_pat,
        F.depsec_nro codigo_set_atu_pat,
        'C' orig_pat,
        P.BPESP_NRO codigo_tip_pat,
        P.ESTCON_NRO codigo_sit_pat,
        P.NOME discr_pat,
        P.OBS obs_pat,
        P.AQUIS_DATA datae_pat,
        p.DATA_INC dtlan_pat,
        P.AQUIS_VALOR valaqu_pat,
        0 valatu_pat,
        P.FO_PES_NRO codigo_for_pat,
        P.AQUIS_NRO_DOC docum,
        NULL percenqtd_pat,
        'V' dae_pat,
        coalesce(P3.vlr_residual, p2.vlr_residual) valres_pat,
        'M' percentemp_pat,
        P.NRO_EMPENHO nempg_pat,
        EXTRACT(YEAR FROM P.DT_EMPENHO) anoemp_pat, 
        p.BX_DOC_DATA dtpag_pat,
		P.NRO || P.NRO_PATR  hash_sinc,
        cast(P.FLG_MODAL_BX as int) codigo_bai_pat,
        p.nro_patr chapa_pat_alt
    FROM
        SYSTEM.D3_BEM_PATR P
    LEFT JOIN SYSTEM.D3_BP_DS P1 ON
        P1.BEPA_NRO = P.NRO
    LEFT JOIN SYSTEM.D3_BP_VALOR P2 ON
        P2.BEPA_NRO = P.NRO
        AND (COALESCE(P2.FLG_REAVALIADO, 'N') = 'N')
    LEFT JOIN SYSTEM.D3_BP_VALOR P3 ON
        P3.BEPA_NRO = P.NRO
        AND (COALESCE(P3.FLG_REAVALIADO, 'N') = 'S')
    LEFT JOIN SYSTEM.D3_BP_ESPECIE E ON
        E.NRO = P.BPESP_NRO
    LEFT JOIN (SELECT b.BEPA_NRO, b.DEPSEC_NRO, b.NROSEQ, b.DATA_INICIO, b.FLG_ADQTRANSF FROM system.D3_BP_DS b WHERE (b.BEPA_NRO, b.NROSEQ) IN ( SELECT BEPA_NRO, MAX(NROSEQ) FROM system.D3_BP_DS GROUP BY BEPA_NRO )) F
       ON f.bepa_nro = p.nro
    WHERE
        (COALESCE(P1.FLG_ADQTRANSF, 'A') = 'A') order by P.NRO`, modules.Cache.Empresa)
	
	totalRows, _ := modules.CountRows(query)
	barCadPat := modules.NewProgressBar(p, totalRows, "CADPAT")

	insert, err := tx.Prepare(`INSERT
		INTO
		pt_cadpat (codigo_pat,
		empresa_pat,
		codigo_gru_pat,
		chapa_pat,
		codigo_cpl_pat,
		codigo_set_pat,
		codigo_set_atu_pat,
		orig_pat,
		codigo_tip_pat,
		codigo_sit_pat,
		discr_pat,
		obs_pat,
		datae_pat,
		dtlan_pat,
		valaqu_pat,
		valatu_pat,
		codigo_for_pat,
		nota_pat,
		percenqtd_pat,
		dae_pat,
		valres_pat,
		percentemp_pat,
		nempg_pat,
		anoemp_pat,
		dtpag_pat,
		hash_sinc,
		codigo_bai_pat,
		chapa_pat_alt)
	VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`)
	if err != nil {
		panic(fmt.Sprintf("Erro ao preparar insert: %v", err))
	}

	rows, err := cnxOra.Queryx(query)
	if err != nil {
		panic(fmt.Sprintf("Erro ao executar query: %v", err))
	}

	for rows.Next() {
		var registro ModelCadpat

		err := rows.StructScan(&registro)
		if err != nil {
			panic(fmt.Sprintf("Erro ao escanear valores: %v", err))
		}

		if registro.DiscrPat.String, err = modules.DecodeToWin1252(registro.DiscrPat.String); err != nil {
			panic(fmt.Sprintf("Erro ao decodificar DiscrPat: %v", err))
		}

		if _, err := insert.Exec(
			registro.CodigoPat,
			registro.EmpresaPat,
			registro.CodigoGruPat,
			registro.ChapaPat,
			registro.CodigoCplPat,
			registro.CodigoSetPat,
			registro.CodigoSetAtuPat,
			registro.OrigPat,
			registro.CodigoTipPat,
			registro.CodigoSitPat,
			registro.DiscrPat,
			registro.ObsPat,
			registro.DataePat,
			registro.DtlanPat,
			registro.ValaquPat,
			registro.ValatuPat,
			registro.CodigoForPat,
			registro.NotaPat,
			nil, //registro.PercenqtdPat,
			registro.DaePat,
			registro.ValresPat,
			registro.PercentempPat,
			registro.NempgPat,
			registro.AnoempPat,
			registro.DtpagPat,
			registro.HashSinc,
			registro.CodigoBaiPat,
			registro.ChapaPatAlt,
		); err != nil {
			panic(fmt.Sprintf("Erro ao inserir valores: %v", err))
		}

		barCadPat.Increment()
	}
}