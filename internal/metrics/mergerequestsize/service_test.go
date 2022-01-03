package mergerequestsize

import (
	"swe-dashboard/internal/models"
	"testing"
	"time"
)

func TestCalculateChanges(t *testing.T) {
	diffline := 6.0
	changes := []*models.MergeRequestChanges{
		{Diff: "@@ -172,6 +172,16 @@ class Announcement\n      */\n     private $is_promoted;\n \n+"},
		{Diff: "@@ -0,0 +1,27 @@\n+<?php\n+\n+namespace ConfigBundle\\Enum;"},
		{Diff: "@@ -4,6 +4,7 @@ namespace ConfigBundle\\Form\\Type\\Announcement;\n \n-"},
		{Diff: "@@ -29,6 +29,8 @@\n                             {{ form_row(form.user_device) }}\n                             <hr>\n                             {{ form_row(form.is_promoted) }}\n+"},
	}

	mrsize := mergeRequestSizes{}
	count := mrsize.calculateChanges(time.Now(), changes)

	if count.Count != diffline {
		t.Fatalf("wrong calculation count value. vgot: %f, want: %f", count.Count, diffline)
	}
}
