package compras

import (
	// "time"

	"github.com/gobuffalo/nulls"
)

type ModelCadunimedida struct {
	Sigla     string       `db:"SIGLA"`
	Descricao string       `db:"DESCRICAO"`
}

type ModelCadgrupo struct {
	Grupo          string       `db:"GRUPO"`
	Nome           string       `db:"NOME"`
	BalcoTce       string       `db:"BALCO_TCE"`
	Ocultar        string       `db:"OCULTAR"`
	BalcoTceSaida  string       `db:"BALCO_TCE_SAIDA"`
	ConvTipo       string       `db:"TIPO"`
	ConvNro        string       `db:"NRO"`
}

type ModelCadsubgr struct {
	Grupo     string       `db:"KEY_GRUPO"`
	Subgrupo  string       `db:"SUBGRUPO"`
	Nome      nulls.String `db:"NOME"`
	Ocultar   nulls.String `db:"OCULTAR"`
	KeySubgrupo nulls.String `db:"KEY_SUBGRUPO"`
}

type ModelCadest struct {
	Cadpro    string       
	Grupo     string       
	Subgrupo  string       `db:"KEY_SUBGRUPO"`
	Codigo    int       `db:"CODIGO"`
	Disc1     nulls.String `db:"DISC1"`
	Tipopro   nulls.String `db:"TIPOPRO"`
	Unid1     nulls.String `db:"UNID1"`
	Quanmin   nulls.Float64 `db:"QUANMIN"`
	Codreduz  nulls.String `db:"CODREDUZ"`
	Quanmax   nulls.Float64 `db:"QUANMAX"`
	Ocultar   nulls.String `db:"OCULTAR"`
	Usopro    nulls.String `db:"USOPRO"`
	SubgrupoNome nulls.String `db:"NOMESUBGRUPO"`
}

type ModelCentroCusto struct {
	Poder     string
	Orgao     string
	Destino   string       `db:"DESTINO"`
	Ccusto    string       `db:"CCUSTO"`
	Descr     string       `db:"DESCR"`
	Codccusto string       `db:"CODCCUSTO"`
	Empresa   string       `db:"EMPRESA"`
	Responsa   nulls.String `db:"RESPONSA"`
	Ocultar   string       `db:"OCULTAR"`
}

type ModelDestino struct {
	Destino string `db:"COD"`
	Descr   string `db:"DESTI"`
	Empresa int    `db:"EMPRESA"`
}

type ModelCadorc struct {
	Sequencia int64	   `db:"SEQUENCIA"` // Sequência única para cada registro
	// Campos comuns
	IdCadorc   float64      `db:"ID_CADORC"`
	Ano           nulls.String `db:"ANO"`
	Num           nulls.String `db:"NUM"`
	Numorc     string       `db:"NUMORC"`
	DtOrc         nulls.Time   `db:"DTORC"`
	NumLic        nulls.Int    `db:"NUMLIC"`
	Descr         nulls.String `db:"DESCR"`
	Obs           nulls.String `db:"OBS"`
	CodCcusto  nulls.Int    `db:"CODCCUSTO"`
	Status        nulls.String `db:"STATUS"`
	Liberado      nulls.String `db:"LIBERADO"`
	ProcLic       nulls.String `db:"PROCLIC"`
	RegistroPreco nulls.String `db:"REGISTROPRECO"`
	Ficha      nulls.String    `db:"FICHA"`
	Desdobro      nulls.String `db:"DESDOBRO"`
	Codreduz 	nulls.String `db:"MATERIAL"`
	Qtd      nulls.Float64 `db:"QTD"`
	Marca   nulls.String  `db:"MARCA"`
	Valor    nulls.Float64 `db:"VALOR"`
	Item     nulls.Int     `db:"ITEM"`
	ItemOrc  nulls.Int     `db:"ITEMORC"`
	Empresa       int64        `db:"EMPRESA"`
	Cadpro   nulls.String
}

type ModelFcadorc struct {
	IdCadorc int64    `db:"ID_CADORC"`
	Numorc   string   `db:"NUMORC"`
	Codif   string   `db:"CODIF"`
	Nome   string   `db:"NOME"`
	Valorc float64  `db:"VALORC"`
}

type ModelVcadorc struct {
	Numorc  string   `db:"NUMORC"`
	IdCadorc int64    `db:"ID_CADORC"`
	Codif  string   `db:"CODIF"`
	Item	int64    `db:"ITEM"`
	VlrUni nulls.Float64  `db:"VLRUNI"`
	VlrTot nulls.Float64  `db:"VLRTOT"`
	Classe nulls.String   `db:"CLASSE"`
	Marca nulls.String   `db:"MARCA"`
	Ganhou nulls.String   `db:"GANHOU"`
	VlrGanhou nulls.Float64  `db:"VLRGANHOU"`
}

type ModelCadlic struct {
	Licit           string       
	Numpro          nulls.Int    `db:"NUMPRO"`
	Datae           nulls.Time   `db:"DATAE"`
	Dtpub           nulls.Time   `db:"DTPUB"`
	Dtenc           nulls.Time   `db:"DTENC"`
	Horenc          nulls.String `db:"HORENC"`
	Horabe          nulls.String `db:"HORABE"`
	Discr           nulls.String `db:"DISCR"`
	Discr7          nulls.String `db:"DISCR7"`
	Modlic          string       `db:"MODLIC"`
	Dthom           nulls.Time   `db:"DTHOM"`
	Dtadj           nulls.Time   `db:"DTADJ"`
	Comp            nulls.String `db:"COMP"`
	Numero          nulls.String `db:"NUMERO"`
	Ano             nulls.String `db:"ANO"`
	Valor           nulls.Float64 `db:"VALOR"`
	Tipopubl        nulls.String `db:"TIPOPUBL"`
	Detalhe         nulls.String `db:"DETALHE"`
	Horreal         nulls.String `db:"HORREAL"`
	Local           nulls.String `db:"LOCAL"`
	Proclic         string       `db:"SEQUENCIA"`
	Numlic          int64        `db:"NUMLIC"`
	Liberacompra    nulls.String `db:"LIBERACOMPRA"`
	Microempresa    nulls.String `db:"MICROEMPRESA"`
	Licnova         nulls.String `db:"LICNOVA"`
	Codtce          nulls.String `db:"CODTCE"`
	ProcessoData    nulls.Time   `db:"PROCESSO_DATA"`
	Codmod          int64        
	Anomod          nulls.String `db:"ANOMOD"`
	Registropreco   nulls.String `db:"REGISTROPRECO"`
	Empresa         int64        `db:"EMPRESA"`
	Modalidade      int        `db:"MODALIDADE"`
	Processo 		nulls.String `db:"PROCLIC"`
	Dtreal 	   nulls.Time   `db:"DTREAL"`
}

type ModelProlics struct {
	Sessao  int64        `db:"SESSAO"`
	Codif  string       `db:"CODIF"`
	Numlic int64        `db:"NUMLIC"`
	Habilitado nulls.String `db:"HABILITADO"`
	Status nulls.String `db:"STATUS"`
	Nome  nulls.String `db:"NOME"`
	Representante nulls.String `db:"REPRESENTANTE"`
	Cpf nulls.String `db:"CPF"`
}

type ModelCadprolic struct {
	Item    int64        `db:"ITEM"`
	Cadpro  string       
	Codreduz nulls.String `db:"MATERIAL"`
	Numlic int64        `db:"NUMLIC"`
	Quan1   nulls.Float64 `db:"QUAN1"`
	Vamed1 nulls.Float64 `db:"VAMED1"`
	Valor  nulls.Float64 `db:"VALOR"`
	Tipo   nulls.String        `db:"TIPO"`
	Codccusto nulls.Int    `db:"CODCCUSTO"`
	Reduz nulls.String `db:"REDUZ"`
	Lotelic nulls.String `db:"LOTELIC"`
	Ficha  nulls.String    `db:"FICHA"`
	QtdItem nulls.Float64 `db:"QUANTIDADE"`
}

type ModelProposta struct {
	Sessao   int64        `db:"SESSAO"`
	Codif  string       `db:"CODIF"`
	Item    int64        `db:"ITEM"`
	Quan1  nulls.Float64 `db:"QUAN1"`
	QtdAdt nulls.Float64 `db:"QUAN1"`
	Vaun1 nulls.Float64 `db:"VAUN1"`
	VaunAdt nulls.Float64 `db:"VAUN1"`
	Vato1 nulls.Float64 `db:"VATO1"`
	Numlic int64        `db:"NUMLIC"`
	Status  nulls.String `db:"STATUS"`
	Subem nulls.String `db:"SUBEM"`
	Marca nulls.String `db:"MARCA"`
	ItemLance string        `db:"ITEM_LANCE"`
	Lotelic nulls.String `db:"LOTELIC"`
	Codreduz nulls.String `db:"MATERIAL"`
	Cadpro string       
	Tpcontrole nulls.String `db:"TPCONTROLE_SALDO"`
	Rodada nulls.Int    `db:"NRO_RODADA"`
}

type ModelPedidos struct {
	Anoreduz  string       `db:"ANOREDUZ"`
	Numped   nulls.String `db:"NUMPED"`
	Ano        nulls.String `db:"EX_ANO"`
	IdCadped   int64        `db:"ID_CADPED"`
	Sequencia  int64        `db:"SEQUENCIA"`
	Cabecalho  int64        `db:"CABECALHO"`
	IdCadpedlicit nulls.Int    `db:"IDCADPED_LICIT"`
	Codif      nulls.String `db:"CODIF"`
	Codccusto  nulls.Int    `db:"CENTROCUSTO"`
	Datped    nulls.Time   `db:"DT_EMISSAO"`
	Ficha      nulls.Int    `db:"FICHA"`
	Entrou     nulls.String `db:"ENTROU"`
	Numlic     nulls.Int    `db:"NUMLIC"`
	Localentg  nulls.String `db:"LOCAL_ENTREGA"`
	Condpgto   nulls.String `db:"COND_PAGTO"`
	Obs        nulls.String `db:"OBS"`
	Contrato  nulls.String `db:"CONTRATO"`
	Item       int64        `db:"ITEM"`
	Material   nulls.String `db:"MATERIAL"`
	Qtd        nulls.Float64 `db:"QTD"`
	Prcunt    nulls.Float64 `db:"PRCUNT"`
	Prctot    nulls.Float64 `db:"PRCTOT"`
	Qtdanu    nulls.Float64 `db:"QTDANU"`
	Prctotanu nulls.Float64 `db:"PRCTOTANU"`
	Categoria  nulls.String `db:"CATEGORIA"`
	Grupo      nulls.String `db:"GRUPO"`
	Modalidade nulls.String `db:"MODALIDADE"`
	Elemento    nulls.String `db:"ELEMENTO"`
	Desdobro    nulls.String `db:"DESDOBRO"`
	Vingrupo  nulls.String `db:"VINGRUPO"`
	Vincodigo nulls.String `db:"VINCODIGO"`
	Destino   nulls.String `db:"DESTINO"`
	Pkemp  nulls.Int    `db:"PKEMP"`
	Empresa   int64        `db:"EMPRESA"`
	Cadpro string
}

type ModelIcadreq struct {
	Id_requi   int           `db:"ID_REQUI"`
	Requi      string        `db:"REQUI"`
	Codccusto  string        `db:"CODCCUSTO"`
	Empresa   string        `db:"EMPRESA"`
	Item       int           `db:"ITEM"`
	Destino    string        `db:"DESTINO"`
	Cadpro     string        `db:"CODREDUZ"`
	Quan1      nulls.Float64 `db:"QUAN1"`
	Quan2      nulls.Float64 `db:"QUAN2"`
	Vaun1      nulls.Float64 `db:"VAUN1"`
	Vaun2      nulls.Float64 `db:"VAUN2"`
	Vato1      nulls.Float64 `db:"VATO1"`
	Vato2      nulls.Float64 `db:"VATO2"`
}

type ModelRequi struct {
	Id_requi  int          `db:"ID_REQUI"`
	Requi     string       
	Num       string          
	Ano       int          `db:"ANO"`
	Destino   string       `db:"DESTINO"`
	Codccusto string       `db:"CODCCUSTO"`
	Datae     nulls.Time    `db:"DATAE"`
	Dtlan     nulls.Time    `db:"DTLAN"`
	Dtpag     nulls.Time    `db:"DTPAG"`
	Comp      string       `db:"COMP"`
	Codif     nulls.String `db:"CODIF"`
	Docum     nulls.String `db:"DOCUM"`
	Tipo      string       `db:"TIPOMOV"`
	Entr      string
	Said      string
	Empresa   string      	 `db:"ENTIDADE"`
	Item       int           `db:"ITEM"`
	Cadpro     string        `db:"CODREDUZ"`
	Quan1      float64 
	Quan2      float64 
	Vaun1      float64 
	Vaun2      float64
	Vato1      float64
	Vato2      float64
	Quantidade float64       `db:"QUANTIDADE"`
	ValorUnit  float64       `db:"VALORUNITARIO"`
	ValorTotal float64       `db:"VALORTOTAL"`
	Numped     nulls.String `db:"NUMPED"`
}

type ModelMotor struct {
	Cod       string       `db:"COD"`
	Nome      string       `db:"NOME"`
	Cnh       nulls.String `db:"CNH"`
	Categcnh  nulls.String `db:"CATEGCNH"`
	Dtvenccnh nulls.Time   `db:"DTVENCCNH"`
}

type ModelVeiculo struct {
	Placa 	string       `db:"PLACA"`
	Sequencia int64        `db:"SEQUENCIA"`
	Modelo string       `db:"MODELO"`
	Chassi string       `db:"CHASSI"`
	Cor    string       `db:"COR"`
	Ano    nulls.Int    `db:"ANO"`
	Anomod nulls.Int    `db:"ANOMOD"`
	Renavam string       `db:"RENAVAM"`
	Aquisicao nulls.Time   `db:"AQUISICAO"`
	Motorista nulls.String `db:"MOTORISTA"`
	CodigoMarcaVeiculo nulls.String `db:"CODIGO_MARCA_VEI"`
	Kminicial nulls.Float64 `db:"KMINICIAL"`
	Obs     nulls.String `db:"OBS"`
	Combustivel nulls.String `db:"COMBUSTIVEL"`
	Alienacao nulls.Time   `db:"ALIENACAO"`
	Licenca nulls.Time   `db:"LICENCA"`
	Trocaoleo nulls.Float64 `db:"TROCAOLEO"`
	Ordem nulls.String `db:"ORDEM"`
}

type ModelAbastecimento struct {
	Id_requi  int          `db:"ID_REQUI"`
	Requi     string       `db:"REQUI"`
	Num       string          `db:"NUM"`
	Ano       int          `db:"ANO"`
	Datae     nulls.Time    `db:"DATAE"`
	Dtlan     nulls.Time    `db:"DTLAN"`
	Dtpag     nulls.Time    `db:"DTPAG"`
	Codif     nulls.String `db:"CODIF"`
	Docum     nulls.String `db:"DOCUM"`
	Codccusto string       
	Destino  string
	Comp      string       `db:"COMP"`
	Tipo      string       `db:"TIPOMOV"`
	Entr      string
	Said      string
	Empresa   string      	 `db:"ENTIDADE"`
	Item       int           `db:"ITEM"`
	Cadpro     string        `db:"CODREDUZ"`
	Quan1      float64       `db:"QUAN1"`
	Quan2      float64       `db:"QUAN2"`
	Vaun1      float64       `db:"VAUN1"`
	Vaun2      float64        `db:"VAUN2"`
	Vato1      float64    `db:"VATO1"`
	Vato2      float64  `db:"VATO2"`
	Motorista nulls.String `db:"MOTORISTA"`
	Placa     nulls.String `db:"PLACA"`
	Km 	   nulls.Float64 `db:"KM"`
}