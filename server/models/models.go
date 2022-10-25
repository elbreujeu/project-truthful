package models

type RegisterInfos struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Email     string `json:"email_address"`
	Birthdate string `json:"birthdate"`
}

type LoginInfos struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserPreviewInfos struct {
	Id          int    `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
}

type Answer struct {
	Author       UserPreviewInfos `json:"author"`
	QuestionText string           `json:"question_text"`
	AnswerText   string           `json:"answer_text"`
	AnswerDate   string           `json:"answer_date"`
	LikeCount    int              `json:"like_count"`
}

type UserProfileInfos struct {
	Id             int      `json:"id"`
	Username       string   `json:"username"`
	DisplayName    string   `json:"display_name"`
	FollowerCount  int      `json:"follower_count"`
	FollowingCount int      `json:"following_count"`
	AnswerCount    int      `json:"answer_count"`
	Answers        []Answer `json:"answers"`
}

type FollowUserInfos struct {
	UserId int  `json:"user_id"`
	Follow bool `json:"follow"`
}

type AskQuestionInfos struct {
	UserId       int    `json:"user_id"`
	QuestionText string `json:"question_text"`
}
