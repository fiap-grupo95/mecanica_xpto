package valueobject

type ServiceOrderStatus string

const (
	StatusRecebida            ServiceOrderStatus = "RECEBIDA"
	StatusEmDiagnostico       ServiceOrderStatus = "EM DIAGNÓSTICO"
	StatusAguardandoAprovacao ServiceOrderStatus = "AGUARDANDO APROVAÇÃO"
	StatusAprovada            ServiceOrderStatus = "APROVADA"
	StatusRejeitada           ServiceOrderStatus = "REJEITADA"
	StatusEmExecucao          ServiceOrderStatus = "EM EXECUÇÃO"
	StatusFinalizada          ServiceOrderStatus = "FINALIZADA"
	StatusEntregue            ServiceOrderStatus = "ENTREGUE"
	StatusCancelada           ServiceOrderStatus = "CANCELADA"
)

func ParseServiceOrderStatus(status string) ServiceOrderStatus {
	switch status {
	case "RECEBIDA":
		return StatusRecebida
	case "EM DIAGNÓSTICO":
		return StatusEmDiagnostico
	case "AGUARDANDO APROVAÇÃO":
		return StatusAguardandoAprovacao
	case "APROVADA":
		return StatusAprovada
	case "REJEITADA":
		return StatusRejeitada
	case "EM EXECUÇÃO":
		return StatusEmExecucao
	case "FINALIZADA":
		return StatusFinalizada
	case "ENTREGUE":
		return StatusEntregue
	case "CANCELADA":
		return StatusCancelada
	default:
		return ServiceOrderStatus(status)
	}
}

func (s ServiceOrderStatus) String() string {
	return string(s)
}
