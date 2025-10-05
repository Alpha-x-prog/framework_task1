package core

// canTransit проверяет допустимость перехода статуса (по id).
// Пример: 1=new, 2=in_work, 3=review, 4=closed, 5=canceled (проверь соответствие в БД).
func CanTransit(from, to int) bool {
	if from == to {
		return true
	}
	switch from {
	case 1: // new
		return to == 2 || to == 5
	case 2: // in_work
		return to == 3 || to == 5
	case 3: // review
		return to == 2 || to == 4 || to == 5
	case 4: // closed
		return false
	case 5: // canceled
		return false
	default:
		return false
	}
}
