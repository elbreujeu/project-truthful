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
		if q.ReceiverId != receiverId {
			t.Errorf("Expected ReceiverId to be %d, but got %d", receiverId, q.ReceiverId)
		}
		if q.CreatedAt != creationTime {
			t.Errorf("Expected CreatedAt to be %v, but got %v", creationTime, q.CreatedAt)
		}
	}
}
