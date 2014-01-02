package utils

type Arr []string

// Add separator before each element of array except the first
//	"", sep => ""
//	"a", sep => "a"
//	"a, b", sep => "a sep b"
func (a Arr) Mult(separator string) string {
	res := ""
	for index, element := range a {
		if index > 0 {
			res = res + separator
		}
		res = res + element
	}
	return res
}
