package jsnm

func (j *Jsnm) Range(do func(i int, ji *Jsnm) bool) *Jsnm {
	arr := j.Arr()
	for idx, it := range arr {
		if do(idx, it) {
			break
		}
	}
	return j
}
