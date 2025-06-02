package main

import (
	"GoGemmap/modules"
	"GoGemmap/modules/compras"
	// "GoGemmap/modules/patrimonio"
	"sync"

	"github.com/vbauerster/mpb"
)

func main() {
	var wg sync.WaitGroup

	// Rotina de Compras
	wg.Add(1)
	pc := mpb.New()
	go func() {
		defer wg.Done()
		modules.LimpaCompras()
		compras.Cadunimedida(pc)
		compras.Grupo(pc)
		compras.Subgrupo(pc)
		compras.Cadest(pc)
		compras.Destino(pc)
		compras.CentroCusto(pc)

		compras.Cadorc(pc)
		compras.Fcadorc(pc)
		compras.Vcadorc(pc)

		modules.LimpaLicitacoes()
		compras.Cadlic(pc)
		compras.Prolics(pc)
		compras.Cadprolic(pc)
		compras.CadproProposta(pc)
		compras.Cadped(pc)

		compras.Motor(pc)
		compras.VeiculoTipo(pc)
		compras.VeiculoMarca(pc)
		compras.Veiculo(pc)
		compras.Abastecimento(pc)
		
		compras.SaldoInicial(pc)
		compras.Requi(pc)
	}()

	// Rotina de PATRIMÔNIO
	// wg.Add(1)
	// pp := mpb.New()
	// go func() {
		// defer wg.Done()

	// 	modules.LimpaPatrimonio()
	// 	patrimonio.Cadresponsavel(pp)
	// 	patrimonio.TipoMov(pp)
	// 	patrimonio.Cadajuste(pp)
	// 	patrimonio.Cadbai(pp)
	// 	patrimonio.Cadsit(pp)
	// 	patrimonio.Cadtip(pp)
	// 	patrimonio.Cadpatg(pp)
	// 	patrimonio.CadpatdCadpats(pp)
	// 	patrimonio.Cadpat(pp)
		// patrimonio.Aquisicao(pp)
		// patrimonio.Movbem(pp)
	// }()

	wg.Wait()
}