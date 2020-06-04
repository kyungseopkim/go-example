package arxml

type ISignal struct {
    Name 		string 		`json:"name"`
    Length		int32		`json:"length"`
    Desc		string		`json:"desc"`
    Ref 		string 		`json:"ref"`
    Init 		float32		`json:"init"`
    IsSigned	bool		`json:"isSigned"`
}

func NewISignal(name string, length int32, desc string, ref string, init float32, isSigned bool) ISignal {
    return ISignal{name, length, desc, ref, init, isSigned}
}

func (isignal ISignal) String() string {
    return ToJson(isignal)
}
