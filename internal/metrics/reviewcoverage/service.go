package reviewcoverage

import (
	"strings"
	"swe-dashboard/internal/models"
)

type SCM interface {
	ListMergeRequest(state, scope string, createdafterday int) (mergerequests models.MergeRequests, err error)
	ListMergeRequestNotes(projectID int, mergeRequestID int) (comments []*models.Comment, err error)
	GetMergeRequestChanges(projectID int, mergeRequestID int) (mergerequest models.MergeRequest, err error)
	GetRepository(projectID int) (repository models.Repo, err error)
}

type ReviewCoverageService interface {
	List() (coverages []models.ItemCount, err error)
}

type reviewCoverage struct {
	scm SCM
}

func NewReviewCoverageService(scm SCM) ReviewCoverageService {
	return &reviewCoverage{
		scm: scm,
	}
}

func (r *reviewCoverage) List() (coverages []models.ItemCount, err error) {
	mergerequests, err := r.scm.ListMergeRequest("merged", "all", 7)
	if err != nil {
		return coverages, err
	}

	coverages = []models.ItemCount{}
	repositories := mergerequests.GroupByRepositories()
	for i := 0; i < len(repositories); i++ {
		repo, err := r.scm.GetRepository(repositories[i].ID)
		if err != nil {
			return coverages, err
		}
		repo.MRs = repositories[i].MRs
		calculates, err := r.calculateCoverage(repo.Name, repo.MRs)
		if err != nil {
			return coverages, err
		}
		coverages = append(coverages, calculates...)
	}

	return coverages, nil
}

func (r *reviewCoverage) calculateCoverage(reponame string, mrs []models.MergeRequest) (coverages []models.ItemCount, err error) {
	coverages = []models.ItemCount{}
	for i := 0; i < len(mrs); i++ {
		mr := mrs[i]
		comments, err := r.scm.ListMergeRequestNotes(mr.ProjectID, mr.IID)
		if err != nil {
			return coverages, err
		}

		mrwithchanges, err := r.scm.GetMergeRequestChanges(mr.ProjectID, mr.IID)
		if err != nil {
			return coverages, err
		}

		commentedfilecount := r.countOfCommentedFiles(comments)
		changesFileCount := r.countOfChangedFiles(mrwithchanges.Changes)
		coverage := float64(commentedfilecount) / float64(changesFileCount)
		coverages = append(coverages, models.ItemCount{
			Name:  reponame,
			Name1: r.cleanTitle(mr.Title),
			Count: coverage,
		})
	}

	return coverages, err
}

func (r *reviewCoverage) countOfCommentedFiles(comments []*models.Comment) int {
	count := 0
	for i := 0; i < len(comments); i++ {
		if comments[i].System {
			continue
		}

		if comments[i].FileName == "" {
			continue
		}
		count++
	}

	return count
}

func (r *reviewCoverage) countOfChangedFiles(changes []*models.MergeRequestChanges) int {
	filenames := map[string]int{}
	for i := 0; i < len(changes); i++ {
		change := changes[i]
		oldpathcount, ok := filenames[change.OldPath]
		if !ok {
			filenames[change.OldPath] = 1
		} else {
			filenames[change.OldPath] = oldpathcount + 1
		}

		newpathcount, ok := filenames[change.NewPath]
		if !ok {
			filenames[change.NewPath] = 1
		} else {
			filenames[change.NewPath] = newpathcount + 1
		}
	}

	return len(filenames)
}

func (r *reviewCoverage) cleanTitle(title string) string {
	title = strings.ReplaceAll(title, "\"", "")
	title = strings.ReplaceAll(title, "/", "-")
	title = strings.ReplaceAll(title, "{", "")
	title = strings.ReplaceAll(title, "}", "")
	return title
}
