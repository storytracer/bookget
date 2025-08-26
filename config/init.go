package config

import (
	"strconv"
	"strings"
)

var Conf Input

// initSeq    false = minimum value <= current page number <= maximum value
func initSeqRange() {
	if Conf.Seq == "" || !strings.Contains(Conf.Seq, ":") {
		return
	}
	m := strings.Split(Conf.Seq, ":")
	if len(m) == 1 {
		Conf.SeqStart, _ = strconv.Atoi(m[0])
		Conf.SeqEnd = Conf.SeqStart
	} else {
		Conf.SeqStart, _ = strconv.Atoi(m[0])
		Conf.SeqEnd, _ = strconv.Atoi(m[1])
	}
	return
}

// initVolumeRange    false = minimum value <= current page number <= maximum value
func initVolumeRange() {
	m := strings.Split(Conf.Volume, ":")
	if len(m) == 1 {
		Conf.VolStart, _ = strconv.Atoi(m[0])
		Conf.VolEnd = Conf.VolStart
	} else {
		Conf.VolStart, _ = strconv.Atoi(m[0])
		Conf.VolEnd, _ = strconv.Atoi(m[1])
	}
	return
}

// PageRange    return true (minimum value <= current page number <= maximum value)
func PageRange(index, size int) bool {
	//not set
	if Conf.SeqStart <= 0 {
		return true
	}
	//negative end page
	if Conf.SeqEnd < 0 && (index-size >= Conf.SeqEnd) {
		return false
	}
	//end page
	if Conf.SeqEnd > 0 {
		//finished
		if index >= Conf.SeqEnd {
			return false
		}
		//start page
		if index+1 >= Conf.SeqStart {
			return true
		}
	} else if index+1 >= Conf.SeqStart { //after start page
		return true
	}
	return false
}

// VolumeRange    return true (minimum value <= current page number <= maximum value)
func VolumeRange(index int) bool {
	//not set
	if Conf.VolStart <= 0 {
		return true
	}
	//negative end page
	if Conf.VolEnd < 0 && index > Conf.VolStart {
		return false
	}
	//end page
	if Conf.VolEnd > 0 {
		//finished
		if index >= Conf.VolEnd {
			return false
		}
		//start page
		if index+1 >= Conf.VolStart {
			return true
		}
	} else if index+1 >= Conf.VolStart { //after start page
		return true
	}
	return false
}
