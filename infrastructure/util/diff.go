package util

import "github.com/sergi/go-diff/diffmatchpatch"

var dmp = diffmatchpatch.New()

func Diff(before, after string) string {
	return dmp.DiffPrettyText(dmp.DiffMain(before, after, true))
}

func MakePatch(before, after string) string {
	return dmp.PatchToText(dmp.PatchMake(before, after))
}

func ApplyPatch()  {
	
}
