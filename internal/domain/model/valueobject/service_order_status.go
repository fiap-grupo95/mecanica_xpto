package valueobject

type ServiceOrderStatus string

const (
	StatusRecebida            ServiceOrderStatus = "Recebida"
	StatusEmDiagnostico       ServiceOrderStatus = "Em Diagnóstico"
	StatusAguardandoAprovacao ServiceOrderStatus = "Aguardando Aprovação"
	StatusAprovada            ServiceOrderStatus = "Aprovada"
	StatusRejeitada           ServiceOrderStatus = "Rejeitada"
	StatusEmExecucao          ServiceOrderStatus = "Em Execução"
	StatusFinalizada          ServiceOrderStatus = "Finalizada"
	StatusEntregue            ServiceOrderStatus = "Entregue"
	StatusCancelada           ServiceOrderStatus = "Cancelada"
)

func ParseServiceOrderStatus(status string) ServiceOrderStatus {
	switch status {
	case "Recebida":
		return StatusRecebida
	case "Em Diagnóstico":
		return StatusEmDiagnostico
	case "Aguardando Aprovação":
		return StatusAguardandoAprovacao
	case "Aprovada":
		return StatusAprovada
	case "Rejeitada":
		return StatusRejeitada
	case "Em Execução":
		return StatusEmExecucao
	case "Finalizada":
		return StatusFinalizada
	case "Entregue":
		return StatusEntregue
	case "Cancelada":
		return StatusCancelada
	default:
		return ServiceOrderStatus(status)
	}
}

func (s ServiceOrderStatus) String() string {
	return string(s)
}
