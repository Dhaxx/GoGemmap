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
}

type ModelPedidos struct {
	Anoreduz  string       `db:"ANOREDUZ"`
	Numped   nulls.String `db:"NUMPED"`
	UltimoPedido nulls.String `db:"ULTIMO_PEDIDO"`
	Ano        nulls.String `db:"ANO"`
	IdCadped   int64        `db:"ID_CADPED"`
	Sequencia  int64        `db:"SEQUENCIA"`
	Npedlicit  nulls.String `db:"NPEDLICIT"`
	IdCadpedlicit nulls.Int    `db:"ID_CADPEDLIT"`
	Codif      nulls.String `db:"CODIF"`
	Codccusto  nulls.Int    `db:"CODCCUSTO"`
	Datped    nulls.Time   `db:"DATPED"`
	Ficha      nulls.Int    `db:"FICHA"`
	Entrou     nulls.String `db:"ENTROU"`
	Numlic     nulls.Int    `db:"NUMLIC"`
	Proclic    nulls.String `db:"PROCLIC"`
	Localentg  nulls.String `db:"localentg"`
	Condpgto   nulls.String `db:"condpgto"`
	Prozoentrega nulls.String `db:"prozoentrega"`
	Obs        nulls.String `db:"obs"`
	Aditamento  nulls.String `db:"aditamento"`
	Contrato  nulls.String `db:"contrato"`
	Item       int64        `db:"item"`
	Material   nulls.String `db:"material"`
	Qtd        nulls.Float64 `db:"qtd"`
	Prcunt    nulls.Float64 `db:"prcunt"`
	Prctot    nulls.Float64 `db:"prctot"`
	Qtdanu    nulls.Float64 `db:"qtdanu"`
	Prctotanu nulls.Float64 `db:"prctotanu"`
	Categoria  nulls.String `db:"categoria"`
	Grupo      nulls.String `db:"grupo"`
	Modalidade nulls.String `db:"modalidade"`
	Elemento    nulls.String `db:"elemento"`
	Desdobro    nulls.String `db:"desdobro"`
	Vingrupo  nulls.String `db:"vingrupo"`
	Vincodigo nulls.String `db:"vincodigo"`
	Destino   nulls.String `db:"destino"`
	Pkemp  nulls.Int    `db:"pkemp"`
	Empresa   int64        `db:"empresa"`
	Cadpro string
}

type ModelIcadreq struct {
	Id_requi   int           `db:"id_requi"`
	Requi      string        `db:"requi"`
	Codccusto  string        `db:"codccusto"`
	Empresa   string        `db:"empresa"`
	Item       int           `db:"item"`
	Destino    string        `db:"destino"`
	Cadpro     string        `db:"codreduz"`
	Quan1      nulls.Float64 `db:"quan1"`
	Quan2      nulls.Float64 `db:"quan2"`
	Vaun1      nulls.Float64 `db:"vaun1"`
	Vaun2      nulls.Float64 `db:"vaun2"`
	Vato1      nulls.Float64 `db:"vato1"`
	Vato2      nulls.Float64 `db:"vato2"`
}

type ModelRequi struct {
	Id_requi  int          `db:"id_requi"`
	Requi     string       `db:"requi"`
	Num       int          `db:"num"`
	Ano       int          `db:"exercicio"`
	Destino   string       `db:"destino"`
	Codccusto string       `db:"codccusto"`
	Datae     nulls.Time    `db:"datae"`
	Dtlan     nulls.Time    `db:"dtlan"`
	Dtpag     nulls.Time    `db:"dtpag"`
	Comp      string       `db:"comp"`
	Codif     nulls.String `db:"fornecedor"`
	Docum     nulls.String `db:"docum"`
	Tipo      string       `db:"tipomov"`
	Entr      string
	Said      string
	Empresa   string      	 `db:"entidade"`
	Item       int           `db:"item"`
	Cadpro     string        `db:"codreduz"`
	Quan1      float64 
	Quan2      float64 
	Vaun1      float64 
	Vaun2      float64
	Vato1      float64
	Vato2      float64
	Quantidade float64
	ValorUnit  float64
	ValorTotal float64
	Numped     nulls.String `db:"numped"`
}

type ModelMotor struct {
	Cod       string       `db:"cod"`
	Nome      string       `db:"nome"`
	Cnh       nulls.String `db:"cnh"`
	Categcnh  nulls.String `db:"categcnh"`
	Dtvenccnh nulls.Time   `db:"dtvenccnh"`
}

type ModelVeiculo struct {
	Placa 	string       `db:"placa"`
	Sequencia int64        `db:"sequencia"`
	Modelo string       `db:"modelo"`
	Chassi string       `db:"chassi"`
	Cor    string       `db:"cor"`
	Ano    nulls.Int    `db:"ano"`
	Anomod nulls.Int    `db:"anomod"`
	Renavam string       `db:"renavam"`
	Aquisicao nulls.Time   `db:"aquisicao"`
	Motorista nulls.String `db:"motorista"`
	CodigoMarcaVeiculo nulls.String `db:"codigo_marca_vei"`
	Kminicial nulls.Float64 `db:"kminicial"`
	Obs     nulls.String `db:"obs"`
	Combustivel nulls.String `db:"combustivel"`
	Alienacao nulls.Time   `db:"alienacao"`
	Licenca nulls.Time   `db:"licenca"`
	Trocaoleo nulls.Float64 `db:"trocaoleo"`
}

type ModelAbastecimento struct {
	Id_requi  int          `db:"id_requi"`
	Requi     string       `db:"requi"`
	Num       int          `db:"num"`
	Ano       int          `db:"exercicio"`
	Datae     nulls.Time    `db:"datae"`
	Dtlan     nulls.Time    `db:"dtlan"`
	Dtpag     nulls.Time    `db:"dtpag"`
	Codif     nulls.String `db:"fornecedor"`
	Docum     nulls.String `db:"docum"`
	Codccusto string       
	Destino  string
	Comp      string       `db:"comp"`
	Tipo      string       `db:"tipomov"`
	Entr      string
	Said      string
	Empresa   string      	 `db:"entidade"`
	Item       int           `db:"item"`
	Cadpro     string        `db:"codreduz"`
	Quan1      float64 
	Quan2      float64 
	Vaun1      float64 
	Vaun2      float64
	Vato1      float64
	Vato2      float64
	Quantidade float64
	ValorUnit  float64
	ValorTotal float64
	Motorista nulls.String `db:"motorista"`
	Placa     nulls.String `db:"placa"`
	Km 	   nulls.Float64 `db:"km"`
}