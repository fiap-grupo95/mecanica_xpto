package valueobject

type Status string

const (
	StatusRecebida            Status = "Recebida"
	StatusEmDiagnostico       Status = "Em diagnóstico"
	StatusAguardandoAprovacao Status = "Aguardando aprovação"
	StatusEmExecucao          Status = "Em execução"
	StatusFinalizada          Status = "Finalizada"
	StatusEntregue            Status = "Entregue"
)

func ParseServiceOrderStatus(desc string) Status {
	switch desc {
	case string(StatusRecebida):
		return StatusRecebida
	case string(StatusEmDiagnostico):
		return StatusEmDiagnostico
	case string(StatusAguardandoAprovacao):
		return StatusAguardandoAprovacao
	case string(StatusEmExecucao):
		return StatusEmExecucao
	case string(StatusFinalizada):
		return StatusFinalizada
	case string(StatusEntregue):
		return StatusEntregue
	default:
		return StatusRecebida // ou algum valor padrão/erro
	}
}

func (s Status) String() string {
	return string(s)
}
