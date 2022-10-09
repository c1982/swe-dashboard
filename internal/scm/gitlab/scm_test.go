package gitlab

import (
	"testing"
)

func TestCommitDiffWeightCalculate(t *testing.T) {
	scm := &SCM{}
	list := []struct {
		weight int
		input  string
	}{
		{4, "test\ntest\ntest\n\n"},
		{25, "--- a/doc/update/5.4-to-6.0.md\n+++ b/doc/update/5.4-to-6.0.md\n@@ -71,6 +71,8 @@\n"},
		{69, "sudo -u git -H bundle exec rake migrate_keys RAILS_ENV=production\n sudo -u git -H bundle exec rake migrate_inline_notes RAILS_ENV=production\n"},
		{38, "\n+sudo -u git -H bundle exec rake gitlab:assets:compile RAILS_ENV=production\n+\n "},
		{14, "```\n \n ### 6. Update config files"},
		{29, "```### 6. Update config files"},
		{0, ""},
		{55, "Binary files a/chr_poses.psb and b/chr_poses.psb differ\n"},
	}

	for _, v := range list {
		weight := scm.calculateCommitDiffWeight(v.input)
		if weight != v.weight {
			t.Errorf("weight not correct. got %d, want: %d", weight, v.weight)
		}
	}
}
