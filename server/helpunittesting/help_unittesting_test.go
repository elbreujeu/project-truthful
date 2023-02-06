package helpunittesting

import (
	"fmt"
	"testing"
	"time"
)

func TestGenerateTestQuestions(t *testing.T) {
	count := 5
	receiverId := 1
	creationTime := time.Now()
	questions := GenerateTestQuestions(count, receiverId, creationTime)
	if len(questions) != count {
		t.Errorf("Expected number of questions to be %d, but got %d", count, len(questions))
	}

	for i, q := range questions {
		expectedText := "question" + fmt.Sprintf("%d", i)
		if q.Text != expectedText {
			t.Errorf("Expected text to be %s, but got %s", expectedText, q.Text)
		}
		if q.IsAuthorAnonymous != false {
			t.Errorf("Expected IsAuthorAnonymous to be %v, but got %v", false, q.IsAuthorAnonymous)
		}
		if q.Author.Id != int64(i) {
			t.Errorf("Expected Author.Id to be %d, but got %d", i, q.Author.Id)
		}
		expectedUsername := "username" + fmt.Sprintf("%d", i)
		if q.Author.Username != expectedUsername {
			t.Errorf("Expected Author.Username to be %s, but got %s", expectedUsername, q.Author.Username)
		}
		expectedDisplayName := "display_name" + fmt.Sprintf("%d", i)
		if q.Author.DisplayName != expectedDisplayName {
			t.Errorf("Expected Author.DisplayName to be %s, but got %s", expectedDisplayName, q.Author.DisplayName)
		}
		if q.ReceiverId != receiverId {
			t.Errorf("Expected ReceiverId to be %d, but got %d", receiverId, q.ReceiverId)
		}
		if q.CreatedAt != creationTime {
			t.Errorf("Expected CreatedAt to be %v, but got %v", creationTime, q.CreatedAt)
		}
	}
}
