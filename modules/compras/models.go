package compras

import (
	// "time"

	"github.com/gobuffalo/nulls"
)

type ModelCadunimedida struct {
	Sigla     string       `db:"sigla"`
	Descricao string       `db:"descricao"`
}

type ModelCadgrupo struct {
	Grupo          string       
	Nome           string       `db:"nome"`
	BalcoTce       string       `db:"balco_tce"`
	Ocultar        string       `db:"ocultar"`
	BalcoTceSaida  string       `db:"balco_tce_saida"`
	ConvTipo       string       `db:"tipo"`
	ConvNro        string       `db:"nro"`
}

type ModelCadsubgr struct {
	Grupo     string       
	Subgrupo  string       `db:"sequencia"`
	Nome      nulls.String `db:"nome"`
	Ocultar   nulls.String `db:"ocultar"`
	ConvTipoSubgr  string       `db:"tipo_subgrupo"`
	ConvNroSubgr   string       `db:"nro_subgrupo"`
	ConvTipoGrupo  string       `db:"tipo_grupo"`
	ConvNroGrupo   string       `db:"nro_grupo"`
	SubgrupoOrig nulls.String `db:"subgrupo"`
}

type ModelCadest struct {
	Cadpro    string       
	Grupo     string       
	Subgrupo  string       `db:"seqcont"`
	Codigo    string       `db:"codigo"`
	Disc1     nulls.String `db:"disc1"`
	Tipopro   nulls.String `db:"tipopro"`
	Unid1     nulls.String `db:"unid1"`
	Quanmin   nulls.Float64 `db:"quanmin"`
	Codreduz  nulls.String `db:"codreduz"`
	Quanmax   nulls.Float64 `db:"quanmax"`
	Ocultar   nulls.String `db:"ocultar"`
	Usopro    nulls.String `db:"usopro"`
	ConvNro   nulls.String `db:"nro_subgrupo"`
	ConvTipo  nulls.String `db:"tipo_subgrupo"`
}