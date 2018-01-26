package jsnm

func (j *Jsnm) Range(do func(i int, ji *Jsnm)) *Jsnm {
	arr := j.Arr()
	for idx, it := range arr {
		do(idx, it)
	}
	return j
}
