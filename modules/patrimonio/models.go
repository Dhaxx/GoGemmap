package patrimonio

import (
	// "time"

	"github.com/gobuffalo/nulls"
)

type ModelCadpat struct {
	CodigoPat       int64         `db:"CODIGO_PAT"`
	EmpresaPat      int64         `db:"EMPRESA_PAT"`
	CodigoGruPat    int64         `db:"CODIGO_GRU_PAT"`
	ChapaPat        nulls.String  `db:"CHAPA_PAT"`
	CodigoCplPat    nulls.String  `db:"CPL"`
	CodigoSetPat    nulls.String  `db:"CODIGO_SET_PAT"`
	CodigoSetAtuPat nulls.String  `db:"CODIGO_SET_ATU_PAT"`
	OrigPat         nulls.String  `db:"ORIG_PAT"`
	CodigoTipPat    int64         `db:"CODIGO_TIP_PAT"`
	CodigoSitPat    int64         `db:"CODIGO_SIT_PAT"`
	DiscrPat        nulls.String  `db:"DISCR_PAT"`
	ObsPat          nulls.String  `db:"OBS_PAT"`
	DataePat        nulls.Time    `db:"DATAE_PAT"`
	DtlanPat        nulls.Time    `db:"DTLAN_PAT"`
	ValaquPat       nulls.Float64 `db:"VALAQU_PAT"`
	ValatuPat       nulls.Float64 `db:"VALATU_PAT"`
	CodigoForPat    nulls.String  `db:"CODIGO_FOR_PAT"`
	NotaPat		 	nulls.String  `db:"DOCUM"`
	PercenqtdPat    nulls.String  `db:"PERCENQTD_PAT"`
	DaePat          nulls.String  `db:"DAE_PAT"`
	ValresPat       nulls.Float64 `db:"VALRES_PAT"`
	PercentempPat   nulls.String  `db:"PERCENTEMP_PAT"`
	NempgPat        nulls.String  `db:"NEMPG_PAT"`
	AnoempPat       nulls.String  `db:"ANOEMP_PAT"`
	DtpagPat        nulls.Time    `db:"DTPAG_PAT"`
	HashSinc        nulls.String  `db:"HASH_SINC"`
	CodigoBaiPat    nulls.Int     `db:"CODIGO_BAI_PAT"`
	ChapaPatAlt     nulls.String  `db:"CHAPA_PAT_ALT"`
}

type ModelMovbem struct {
	CodigoMov       int64         `db:"CODIGO_MOV"`
	EmpresaMov      int64         `db:"EMPRESA_MOV"`
	CodigoPatMov    int64         `db:"CODIGO_PAT_MOV"`
	DataMov         nulls.Time    `db:"DATA_MOV"`
	TipoMov         nulls.String  `db:"TIPO_MOV"`
	CodigoSetMov    nulls.String  `db:"CODIGO_SET_MOV"`
	HistoricoMov    nulls.String  `db:"HISTORICO_MOV"`
	ValorMov        nulls.Float64 `db:"VALOR_MOV"`
	DepreciacaoMov  nulls.String  `db:"DEPRECIACAO_MOV"`
	DtContabil	  nulls.Time    `db:"DATA_MOV"`
}