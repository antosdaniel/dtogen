package simple_struct

func ToDst(src Src) Dst {
	return Dst{
		Bool:   src.Bool,
		Int:    src.Int,
		String: src.String,
		Time:   src.Time,
	}
}
