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

func (s ServiceOrderStatus) IsValid() bool {
	switch s {
	case StatusRecebida, StatusEmDiagnostico, StatusAguardandoAprovacao,
		StatusAprovada, StatusRejeitada, StatusEmExecucao,
		StatusFinalizada, StatusEntregue, StatusCancelada:
		return true
	default:
		return false
	}
}

func (s ServiceOrderStatus) IsSame(c ServiceOrderStatus) bool {
	return s == c
}

func (s ServiceOrderStatus) IsRecebida() bool {
	return s == StatusRecebida
}

func (s ServiceOrderStatus) IsEmDiagnostico() bool {
	return s == StatusEmDiagnostico
}
func (s ServiceOrderStatus) IsAguardandoAprovacao() bool {
	return s == StatusAguardandoAprovacao
}
func (s ServiceOrderStatus) IsAprovada() bool {
	return s == StatusAprovada
}
func (s ServiceOrderStatus) IsRejeitada() bool {
	return s == StatusRejeitada
}
func (s ServiceOrderStatus) IsEmExecucao() bool {
	return s == StatusEmExecucao
}
func (s ServiceOrderStatus) IsFinalizada() bool {
	return s == StatusFinalizada
}
func (s ServiceOrderStatus) IsEntregue() bool {
	return s == StatusEntregue
}
func (s ServiceOrderStatus) IsCancelada() bool {
	return s == StatusCancelada
}

func (s ServiceOrderStatus) String() string {
	return string(s)
}
