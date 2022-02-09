package mergerequestparticipants

import (
	"fmt"
	"swe-dashboard/internal/models"
)

type engageItem struct {
	Author      *models.User
	Participant *models.User
	Count       float64
}

type SCM interface {
	GetRepository(projectID int) (repository models.Repo, err error)
	ListMergeRequest(state, scope string, createdafterday int) (mergerequests models.MergeRequests, err error)
	GetMergeRequestParticipants(projectID int, mergeRequestID int) (users []*models.User, err error)
}

type MergeRequestParticipantsService interface {
	List() ([]models.ItemCount, error)
	EngageParticipants() []models.ItemCount
	Engagements() []models.ItemCount
}

func NewMergeRequestParticipantsService(scm SCM) MergeRequestParticipantsService {
	return &mergeRequestParticipants{
		scm:                    scm,
		engagementParticipants: []models.ItemCount{},
		engagements:            []models.ItemCount{},
	}
}

type mergeRequestParticipants struct {
	scm                    SCM
	engagementParticipants []models.ItemCount
	engagements            []models.ItemCount
}

func (mrp *mergeRequestParticipants) EngageParticipants() []models.ItemCount {
	return mrp.engagementParticipants
}

func (mrp *mergeRequestParticipants) Engagements() []models.ItemCount {
	return mrp.engagements
}

func (mrp *mergeRequestParticipants) List() (participants []models.ItemCount, err error) {
	mergerequests, err := mrp.scm.ListMergeRequest("merged", "all", 30)
	if err != nil {
		return participants, err
	}

	participants = []models.ItemCount{}
	repositories := mergerequests.GroupByRepositories()

	for i := 0; i < len(repositories); i++ {
		repo, err := mrp.scm.GetRepository(repositories[i].ID)
		if err != nil {
			return participants, err
		}

		repo.MRs = repositories[i].MRs
		participantCount := map[int]*models.ItemCount{}
		engagementParticipantCount := map[int]map[int]*engageItem{}
		engagementsCount := map[string]*engageItem{}

		for i := 0; i < len(repo.MRs); i++ {
			mr := repo.MRs[i]
			if mr.Author.ID == mr.MergedBy.ID {
				continue
			}

			users, err := mrp.scm.GetMergeRequestParticipants(mr.ProjectID, mr.IID)
			if err != nil {
				return participants, err
			}
			mrp.countParticipants(mr.Author, participantCount, users)
			mrp.engagementParticipantCounts(mr.Author, engagementParticipantCount, users)
			mrp.engagementCounts(mr, engagementsCount)
		}

		for _, v := range participantCount {
			participants = append(participants, models.ItemCount{
				Name:  repo.Name,
				Name1: v.Name1,
				Name2: v.Name2,
				Count: v.Count,
			})
		}

		for _, v := range engagementParticipantCount {
			for _, vv := range v {
				mrp.engagementParticipants = append(mrp.engagementParticipants, models.ItemCount{
					Name:  repo.Name,
					Name1: vv.Author.Name,
					Name2: vv.Participant.Name,
					Count: vv.Count,
				})
			}
		}

		for _, v := range engagementsCount {
			mrp.engagements = append(mrp.engagements, models.ItemCount{
				Name:  repo.Name,
				Name1: v.Author.Name,
				Name2: v.Participant.Name,
				Count: v.Count,
			})
		}
	}

	return participants, nil
}

func (mrp *mergeRequestParticipants) countParticipants(author *models.User, counts map[int]*models.ItemCount, mrparticipants []*models.User) {
	for p := 0; p < len(mrparticipants); p++ {
		user := mrparticipants[p]
		if user.ID == author.ID {
			continue
		}

		v, ok := counts[user.ID]
		if !ok {
			counts[user.ID] = &models.ItemCount{
				Name:  "", // for repository name
				Name1: user.Username,
				Name2: user.Name,
				Count: 1,
			}
		} else {
			v.Count = v.Count + 1
			counts[user.ID] = v
		}
	}
}

func (mrp *mergeRequestParticipants) engagementParticipantCounts(author *models.User, counts map[int]map[int]*engageItem, participants []*models.User) {
	for i := 0; i < len(participants); i++ {
		participant := participants[i]
		if participant.ID == author.ID {
			continue
		}
		v, ok := counts[author.ID]
		if !ok {
			count := map[int]*engageItem{}
			count[participant.ID] = &engageItem{
				Author:      author,
				Participant: participant,
				Count:       1,
			}
			counts[author.ID] = count
		} else {
			vv, ok := v[participant.ID]
			if !ok {
				counts[author.ID][participant.ID] = &engageItem{
					Author:      author,
					Participant: participant,
					Count:       1,
				}
			} else {
				vv.Count = vv.Count + 1
				counts[author.ID][participant.ID] = vv
			}
		}
	}
}

func (mrp *mergeRequestParticipants) engagementCounts(mr models.MergeRequest, counts map[string]*engageItem) {
	pair := fmt.Sprintf("%d-%d", mr.Author.ID, mr.MergedBy.ID)
	v, ok := counts[pair]
	if !ok {
		counts[pair] = &engageItem{
			Author:      mr.Author,
			Participant: mr.MergedBy,
			Count:       1,
		}
	} else {
		v.Count = v.Count + 1
		counts[pair] = v
	}
}
