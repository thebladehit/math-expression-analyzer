package tree

func Unique(forms []*Node) []*Node {
	seen := map[string]bool{}
	var result []*Node
	for _, f := range forms {
		if key := ToInfix(f); !seen[key] {
			seen[key] = true
			result = append(result, f)
		}
	}
	return result
}
