package Symbol

type Symbol struct {
	Stype  string
	Skind  string
	Sindex int
}

func New(St string, Sk string, Si int) Symbol {
	mofa := Symbol{St, Sk, Si}
	return mofa

}
func GetType(s Symbol) string {
	return s.Stype
}

func GetKind(s Symbol) string {
	return s.Skind
}

func GetIndex(s Symbol) int {
	return s.Sindex
}
