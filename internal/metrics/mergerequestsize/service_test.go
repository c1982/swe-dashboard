package mergerequestsize

import (
	"swe-dashboard/internal/models"
	"testing"
)

func TestCalculateChanges(t *testing.T) {
	wantnewline := 5.0
	wantdeletedline := 1.0
	changes := []*models.MergeRequestChanges{
		{Diff: "@@ -172,6 +172,16 @@ class Announcement\n      */\n     private $is_promoted;\n \n+"},
		{Diff: "@@ -0,0 +1,27 @@\n+<?php\n+\n+namespace ConfigBundle\\Enum;"},
		{Diff: "@@ -4,6 +4,7 @@ namespace ConfigBundle\\Form\\Type\\Announcement;\n \n-"},
		{Diff: "@@ -29,6 +29,8 @@\n                             {{ form_row(form.user_device) }}\n                             <hr>\n                             {{ form_row(form.is_promoted) }}\n+"},
	}

	mrsize := mergeRequestSizes{}
	newline, deletedline := mrsize.calculateChanges(changes)

	if newline != wantnewline {
		t.Fatalf("wrong newline count value. got: %f, want: %f", newline, wantnewline)
	}

	if deletedline != wantdeletedline {
		t.Fatalf("wrong deletedline count value. got: %f, want: %f", deletedline, wantdeletedline)
	}
}
