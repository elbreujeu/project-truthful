package models

import (
	"time"
)

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

type UserPreview struct {
	Id          int64  `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
}

type Question struct {
	Id                int         `json:"id"`
	Text              string      `json:"text"`
	IsAuthorAnonymous bool        `json:"is_author_anonymous"`
	Author            UserPreview `json:"author"`
	ReceiverId        int         `json:"receiver_id"`
	CreatedAt         time.Time   `json:"created_at"`
}

type Answer struct {
	Id                int         `json:"id"`
	IsAuthorAnonymous bool        `json:"is_author_anonymous"`
	Author            UserPreview `json:"author"`
	QuestionText      string      `json:"question_text"`
	AnswerText        string      `json:"answer_text"`
	AnswerDate        string      `json:"answer_date"`
	CreatedAt         time.Time   `json:"date_answered"`
	LikeCount         int         `json:"like_count"`
	LikedByRequester  bool        `json:"liked_by_requester"`
}

type UserProfileInfos struct {
	Id                    int      `json:"id"`
	Username              string   `json:"username"`
	DisplayName           string   `json:"display_name"`
	FollowerCount         int      `json:"follower_count"`
	FollowingCount        int      `json:"following_count"`
	AnswerCount           int      `json:"answer_count"`
	IsFollowedByRequester bool     `json:"followed_by_requester"`
	IsRequestingSelf      bool     `json:"requesting_user"`
	Answers               []Answer `json:"answers"`
}

type FollowUserInfos struct {
	UserId int  `json:"user_id"`
	Follow bool `json:"follow"`
}

type AskQuestionInfos struct {
	UserId            int    `json:"user_id"`
	QuestionText      string `json:"text"`
	IsAuthorAnonymous bool   `json:"is_author_anonymous"`
}

type AnswerQuestionInfos struct {
	QuestionId int    `json:"question_id"`
	AnswerText string `json:"text"`
}

type LikeAnswerInfos struct {
	AnswerId int  `json:"answer_id"`
	Like     bool `json:"like"`
}

type DeleteAnswerInfos struct {
	AnswerId int `json:"answer_id"`
}

type DeleteQuestionInfos struct {
	QuestionId int `json:"question_id"`
}

type UpdateUserInfos struct {
	DisplayName string `json:"display_name"`
	Email       string `json:"email_address"`
}

type PromoteUserInfos struct {
	UserId      int    `json:"user_id"`
	PromoteType string `json:"promote_type"`
}

type BanUserInfos struct {
	UserId   int    `json:"user_id"`
	Duration int    `json:"duration"` // in hours, if 0, then it's permanent
	Reason   string `json:"reason"`
}

type PardonUserInfos struct {
	BanId int `json:"ban_id"`
}

type RateLimit struct {
	IpAddress       string
	RequestCount    int
	LastRequestTime time.Time
}

type OauthLoginInfos struct {
	Provider string `json:"provider"`
	Token    string `json:"token"`
}

type GoogleInfos struct {
	Issuer          string `json:"iss"`
	AuthorizedParty string `json:"azp"`
	Audience        string `json:"aud"`
	Subject         string `json:"sub"`
	Email           string `json:"email"`
	EmailVerified   bool   `json:"email_verified"`
	NotValidBefore  int64  `json:"nbf"`
	Name            string `json:"name"`
	Picture         string `json:"picture"`
	GivenName       string `json:"given_name"`
	FamilyName      string `json:"family_name"`
	IssuedAt        int64  `json:"iat"`
	ExpiryDate      int64  `json:"exp"`
	JwtId           string `json:"jti"`
}
