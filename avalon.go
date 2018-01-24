package main

import (
	"math/rand"
	"strconv"
)

type Role struct {
	roleType int
	info     []int
}

func (r *Role) String() string {
	return strconv.Itoa(r.roleType)
}

const (
	rtLoyal = iota
	rtMinion
	rtPercival
	rtMerlin
	rtOberon
	rtMordred
	rtAssassin
	rtMorgana
)

var roleNames = map[int]string{
	rtLoyal:    "Loyal servant of Arthur",
	rtMinion:   "Minion of Mordred",
	rtPercival: "Percival",
	rtMerlin:   "Merlin",
	rtOberon:   "Oberon",
	rtMordred:  "Mordred",
	rtAssassin: "Assassin",
	rtMorgana:  "Morgana",
}

var (
	infoLoyal    = []int{}
	infoMinion   = []int{rtMinion, rtMordred, rtAssassin, rtMorgana}
	infoMerlin   = []int{rtMinion, rtOberon, rtAssassin, rtMorgana}
	infoPercival = []int{rtMerlin, rtMorgana}
	infoOberon   = []int{rtMinion, rtMordred, rtAssassin, rtMorgana}
	infoMordred  = []int{rtMinion, rtAssassin, rtMorgana}
	infoAssassin = []int{rtMinion, rtMordred, rtMorgana}
	infoMorgana  = []int{rtMinion, rtMordred, rtAssassin}
)

var info map[int][]int

func init() {
	info = make(map[int][]int)
	info[rtLoyal] = infoLoyal
	info[rtMinion] = infoMinion
	info[rtPercival] = infoPercival
	info[rtMerlin] = infoMerlin
	info[rtOberon] = infoOberon
	info[rtMordred] = infoMordred
	info[rtAssassin] = infoAssassin
	info[rtMorgana] = infoMorgana
}

func numBad(nplayers int) int {
	off := 1
	if nplayers%2 == 1 {
		off = 0
	}
	if nplayers == 9 {
		return 3
	}
	return nplayers/2 - off
}

func getInfo(index, role int, roles []*Role) []int {
	result := make([]int, 0, len(roles))

	for _, role := range info[role] {
		for i, r := range roles {
			if r.roleType == role && index != i {
				result = append(result, i)
			}
		}
	}
	for i := range result {
		j := rand.Intn(i + 1)
		result[i], result[j] = result[j], result[i]
	}

	return result
}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func AssignRoles(nplayers int, special []int) []*Role {
	roles := make([]*Role, 0, nplayers)
	nbad := numBad(nplayers)
	ngood := nplayers - nbad
	done := make([]int, 0, len(special))
	for i := 0; i < ngood; i++ {
		r := new(Role)

		if !contains(done, rtPercival) && contains(special, rtPercival) {
			done = append(done, rtPercival)
			r.roleType = rtPercival
		} else if !contains(done, rtMerlin) && contains(special, rtMerlin) {
			done = append(done, rtMerlin)
			r.roleType = rtMerlin
		} else {
			r.roleType = rtLoyal
		}
		roles = append(roles, r)
	}
	for i := 0; i < nbad; i++ {
		r := new(Role)
		if !contains(done, rtOberon) && contains(special, rtOberon) {
			done = append(done, rtOberon)
			r.roleType = rtOberon
		} else if !contains(done, rtMordred) && contains(special, rtMordred) {
			done = append(done, rtMordred)
			r.roleType = rtMordred
		} else if !contains(done, rtAssassin) && contains(special, rtAssassin) {
			done = append(done, rtAssassin)
			r.roleType = rtAssassin
		} else if !contains(done, rtMorgana) && contains(special, rtMorgana) {
			done = append(done, rtMorgana)
			r.roleType = rtMorgana
		} else {
			r.roleType = rtMinion
		}
		roles = append(roles, r)
	}

	for i := range roles {
		j := rand.Intn(i + 1)
		roles[i], roles[j] = roles[j], roles[i]
	}

	for i, r := range roles {
		r.info = getInfo(i, r.roleType, roles)
	}

	return roles
}
