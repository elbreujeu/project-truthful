package helpunittesting

import (
	"fmt"
	"project_truthful/models"
	"time"
)

func GenerateTestQuestions(count int, receiverId int, creationTime time.Time) []models.Question {
	questions := make([]models.Question, count)
	for i := 0; i < count; i++ {
		questions[i] = models.Question{
			Id:                i,
			Text:              "question" + fmt.Sprintf("%d", i),
			Author:            models.UserPreview{Id: int64(i), Username: "username" + fmt.Sprintf("%d", i), DisplayName: "display_name" + fmt.Sprintf("%d", i)},
			IsAuthorAnonymous: false,
			ReceiverId:        receiverId,
			CreatedAt:         creationTime,
		}
	}
	return questions
}
