package sliceutils //can also use build in slicing packages to do this job 


func RemoveByIndex[T any](s []T, i int) []T {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}