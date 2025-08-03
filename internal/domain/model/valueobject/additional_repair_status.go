package valueobject

const (
	StatusAberta      Status = "ABERTA"
	StatusEmAndamento Status = "EM_ANDAMENTO"
	StatusConcluida   Status = "CONCLUIDA"
	StatusCancelada   Status = "CANCELADA"
)

func ParseAdditionalRepairStatus(desc string) Status {
	switch desc {
	case string(StatusAberta):
		return StatusAberta
	case string(StatusEmAndamento):
		return StatusEmAndamento
	case string(StatusConcluida):
		return StatusConcluida
	case string(StatusCancelada):
		return StatusCancelada
	default:
		return StatusAberta
	}
}
