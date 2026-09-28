package an

type StateType int

const (
	StState StateType = iota
	OPState
	CONSTState
	OBState
	CBState
	VARState
	FNState
	UMState
	ENDState
)

type Transition struct {
	Next StateType
	Err  string
}

var StateTable = map[StateType]map[TokenType]Transition{
	StState: {
		TokenCONST: Transition{Next: CONSTState},
		TokenVAR:   Transition{Next: VARState},
		TokenFN:    Transition{Next: FNState},
		TokenOB:    Transition{Next: OBState},
		TokenUM:    Transition{Next: UMState},
		TokenEND:   Transition{Next: ENDState, Err: "Пустий вираз"},
		TokenOP:    Transition{Next: OPState, Err: "Вираз не може починатися з операції '{t}'"},
		TokenCB:    Transition{Next: CBState, Err: "Вираз не може починатися із закриваючої дужки '{t}'"},
	},
	OPState: {
		TokenCONST: Transition{Next: CONSTState},
		TokenVAR:   Transition{Next: VARState},
		TokenFN:    Transition{Next: FNState},
		TokenOB:    Transition{Next: OBState},
		TokenUM:    Transition{Next: UMState},
		TokenOP:    Transition{Next: OPState, Err: "Повторна операція '{t}'"},
		TokenCB:    Transition{Next: CBState, Err: "Закриваюча дужка після операцією ')'"},
		TokenEND:   Transition{Next: ENDState, Err: "Закінчення операцією"},
	},
	CONSTState: {
		TokenCONST: Transition{Next: CONSTState, Err: "Дві консанти підряд '{t}'"},
		TokenVAR:   Transition{Next: VARState, Err: "Змінна після константи '{t}'"},
		TokenFN:    Transition{Next: FNState, Err: "Функція після константи '{t}'"},
		TokenOB:    Transition{Next: OBState, Err: "Відкриваюча дужка після константи '('"},
		TokenUM:    Transition{Next: UMState, Err: "Унарний мінус після константи"},
		TokenOP:    Transition{Next: OPState},
		TokenCB:    Transition{Next: CBState},
		TokenEND:   Transition{Next: ENDState},
	},
	VARState: {
		TokenCONST: Transition{Next: CONSTState, Err: "Константа після змінної '{t}'"},
		TokenVAR:   Transition{Next: VARState, Err: "Дві змінні підряд '{t}'"},
		TokenFN:    Transition{Next: FNState, Err: "Функція після змінної '{t}'"},
		TokenOB:    Transition{Next: OBState, Err: "Відкриваюча дужка після змінної '('"},
		TokenUM:    Transition{Next: UMState, Err: "Унарний мінус після змінної"},
		TokenOP:    Transition{Next: OPState},
		TokenCB:    Transition{Next: CBState},
		TokenEND:   Transition{Next: ENDState},
	},
	OBState: {
		TokenCONST: Transition{Next: CONSTState},
		TokenVAR:   Transition{Next: VARState},
		TokenFN:    Transition{Next: FNState},
		TokenOB:    Transition{Next: OBState},
		TokenUM:    Transition{Next: UMState},
		TokenOP:    Transition{Next: OPState, Err: "Операція після відкриваючої дужки '{t}'"},
		TokenCB:    Transition{Next: CBState, Err: "Пусті дужки"},
		TokenEND:   Transition{Next: ENDState, Err: "Не закрита дужка"},
	},
	CBState: {
		TokenCONST: Transition{Next: CONSTState, Err: "Константа після закриваючої дужки '{t}'"},
		TokenVAR:   Transition{Next: VARState, Err: "Змінна після закриваючої дужки '{t}'"},
		TokenFN:    Transition{Next: FNState, Err: "Функція після закриваючої дужки '{t}'"},
		TokenOB:    Transition{Next: OBState, Err: "Пропущена операція після закриваючої дужки"},
		TokenUM:    Transition{Next: UMState, Err: "Пропущена операція після закриваючої дужки"},
		TokenOP:    Transition{Next: OPState},
		TokenCB:    Transition{Next: CBState},
		TokenEND:   Transition{Next: ENDState},
	},
	FNState: {
		TokenCONST: Transition{Next: CONSTState, Err: "Константа після функції '{t}'"},
		TokenVAR:   Transition{Next: VARState, Err: "Змінна після функції '{t}'"},
		TokenFN:    Transition{Next: FNState, Err: "Функція після функції '{t}'"},
		TokenOB:    Transition{Next: OBState},
		TokenUM:    Transition{Next: UMState, Err: "Унарний мінус після функції '-'"},
		TokenOP:    Transition{Next: OPState, Err: "Операція після функції '{t}'"},
		TokenCB:    Transition{Next: CBState, Err: "Закриваюча дужка після функції ')'"},
		TokenEND:   Transition{Next: ENDState, Err: "Не закінчена функція"},
	},
	UMState: {
		TokenCONST: Transition{Next: CONSTState},
		TokenVAR:   Transition{Next: VARState},
		TokenFN:    Transition{Next: FNState},
		TokenOB:    Transition{Next: OBState},
		TokenUM:    Transition{Next: UMState, Err: "Унарний мінус після унарного мінусу '-'"},
		TokenOP:    Transition{Next: OPState, Err: "Операція після унарного мінусу '{t}'"},
		TokenCB:    Transition{Next: CBState, Err: "Закриваюча дужка після унарного мінусу ')'"},
		TokenEND:   Transition{Next: ENDState, Err: "Унарний мінус вкінці виразу"},
	},
}

//-------------
//
//type State interface {
//	GetNext(sData *SharedData, token Token) (State, *Error)
//	ValidateEnd(sData *SharedData, token Token) *Error
//}
//
//func checkOpenBrackets(sData *SharedData, token Token) (State, *Error) {
//	sData.brackets = append(sData.brackets, token.Pos)
//	return &OBState{}, nil
//}
//
//func checkCloseBracket(sData *SharedData, token Token) (State, *Error) {
//	if len(sData.brackets) == 0 {
//		return nil, &Error{Pos: token.Pos, Msg: "Зайва закриваюча дужка"}
//	}
//	sData.brackets = sData.brackets[:len(sData.brackets)-1]
//	return &CBState{}, nil
//}
//
//type StartState struct{}
//
//func (s *StartState) GetNext(sData *SharedData, token Token) (State, *Error) {
//	switch token.Type {
//	case TokenCONST:
//		return &ConstState{}, nil
//	case TokenUM:
//		return &UMinusState{}, nil
//	case TokenVAR:
//		return &VarState{}, nil
//	case TokenFN:
//		return &FnState{}, nil
//	case TokenOB:
//		return checkOpenBrackets(sData, token)
//	case TokenOP:
//		return nil, &Error{Pos: token.Pos, Msg: "Неможна використовувати + на початку"}
//	case TokenCB:
//		return nil, &Error{Pos: token.Pos, Msg: "Не можна використовувати закриваючу дужку на початку"}
//	default:
//		return nil, &Error{Pos: token.Pos, Msg: "Неможлива операція"}
//	}
//}
//
//func (s *StartState) ValidateEnd(sData *SharedData, token Token) *Error {
//	return &Error{Pos: 0, Msg: "Пустий рядок"}
//}
//
//type ConstState struct{}
//
//func (c *ConstState) GetNext(sData *SharedData, token Token) (State, *Error) {
//	switch token.Type {
//	case TokenCONST:
//		return nil, &Error{Pos: token.Pos, Msg: "Два числа підряд"}
//	case TokenCB:
//		return checkCloseBracket(sData, token)
//	case TokenOP:
//		return &OpState{}, nil
//	default:
//		return nil, &Error{Pos: token.Pos, Msg: "Неможлива операція"}
//	}
//}
//func (c *ConstState) ValidateEnd(sData *SharedData, token Token) *Error {
//	return nil
//}
//
//type OpState struct{}
//
//func (s *OpState) GetNext(sData *SharedData, token Token) (State, *Error) {
//	switch token.Type {
//	case TokenOP:
//		return nil, &Error{Pos: token.Pos, Msg: "Дві операції підряд"}
//	case TokenCONST:
//		return &ConstState{}, nil
//	case TokenVAR:
//		return &VarState{}, nil
//	case TokenFN:
//		return &FnState{}, nil
//	case TokenOB:
//		return checkOpenBrackets(sData, token)
//	case TokenUM:
//		return &UMinusState{}, nil
//	default:
//		return nil, &Error{Pos: token.Pos, Msg: "Неможлива операція"}
//	}
//}
//func (s *OpState) ValidateEnd(sData *SharedData, token Token) *Error {
//	return &Error{Pos: token.Pos, Msg: "Рядок не може закінчуватися операцією"}
//}
//
//type UMinusState struct{}
//
//func (s *UMinusState) GetNext(sData *SharedData, token Token) (State, *Error) {
//	switch token.Type {
//	case TokenOB:
//		return checkOpenBrackets(sData, token)
//	case TokenVAR:
//		return &VarState{}, nil
//	case TokenFN:
//		return &FnState{}, nil
//	case TokenCONST:
//		return &ConstState{}, nil
//	default:
//		return nil, &Error{Pos: token.Pos, Msg: "Неможлива операція"}
//	}
//}
//func (s *UMinusState) ValidateEnd(sData *SharedData, token Token) *Error {
//	return &Error{Pos: token.Pos, Msg: "Рядок не може закінчуватися унарним мінусом"}
//}
//
//type VarState struct{}
//
//func (s *VarState) GetNext(sData *SharedData, token Token) (State, *Error) {
//	switch token.Type {
//	case TokenVAR:
//		return nil, &Error{Pos: token.Pos, Msg: "Дві змінні підряд"}
//	case TokenCB:
//		return checkCloseBracket(sData, token)
//	case TokenOP:
//		return &OpState{}, nil
//	default:
//		return nil, &Error{Pos: token.Pos, Msg: "Неможлива операція"}
//	}
//}
//func (s *VarState) ValidateEnd(sData *SharedData, token Token) *Error {
//	return nil
//}
//
//type FnState struct{}
//
//func (s *FnState) GetNext(sData *SharedData, token Token) (State, *Error) {
//	switch token.Type {
//	case TokenOB:
//		return checkOpenBrackets(sData, token)
//	default:
//		return nil, &Error{Pos: token.Pos, Msg: "Неможлива операція"}
//	}
//}
//func (s *FnState) ValidateEnd(sData *SharedData, token Token) *Error {
//	return &Error{Pos: token.Pos, Msg: "Рядок не може закінчуватися функцією"}
//}
//
//type OBState struct{}
//
//func (s *OBState) GetNext(sData *SharedData, token Token) (State, *Error) {
//	switch token.Type {
//	case TokenCONST:
//		return &ConstState{}, nil
//	case TokenUM:
//		return &UMinusState{}, nil
//	case TokenVAR:
//		return &VarState{}, nil
//	case TokenFN:
//		return &FnState{}, nil
//	case TokenOB:
//		return checkOpenBrackets(sData, token)
//	case TokenCB:
//		return nil, &Error{Pos: token.Pos, Msg: "Не можна використовувати пусті дужки"}
//	default:
//		return nil, &Error{Pos: token.Pos, Msg: "Неможлива операція"}
//	}
//}
//func (s *OBState) ValidateEnd(sData *SharedData, token Token) *Error {
//	return &Error{Pos: token.Pos, Msg: "Рядок не може закінчуватися відкритою дужкою"}
//}
//
//type CBState struct{}
//
//func (s *CBState) GetNext(sData *SharedData, token Token) (State, *Error) {
//	switch token.Type {
//	case TokenOP:
//		return &OpState{}, nil
//	case TokenCB:
//		return checkCloseBracket(sData, token)
//	default:
//		return nil, &Error{Pos: token.Pos, Msg: "Неможлива операція"}
//	}
//}
//func (s *CBState) ValidateEnd(sData *SharedData, token Token) *Error {
//	return nil
//}
