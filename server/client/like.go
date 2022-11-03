package client

import (
	"errors"
	"net/http"
	"project_truthful/client/database"
)

func LikeAnswer(userId int, postId int) (int, error) {
	userExists, err := database.CheckUserIdExists(userId, database.DB)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if !userExists {
		return http.StatusNotFound, errors.New("user not found")
	}

	postExists, err := database.CheckAnswerIdExists(postId, database.DB)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if !postExists {
		return http.StatusNotFound, errors.New("post not found")
	}

	likeExists, err := database.CheckLikeExists(userId, postId, database.DB)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if likeExists {
		return http.StatusBadRequest, errors.New("user already likes this post")
	}

	err = database.AddLike(userId, postId, database.DB)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return 0, nil
}

func UnlikeAnswer(userId int, postId int) (int, error) {
	// no need to check whether user or post exists, since we can just delete the like
	likeExists, err := database.CheckLikeExists(userId, postId, database.DB)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	if !likeExists {
		return http.StatusBadRequest, errors.New("user does not like this post")
	}

	err = database.RemoveLike(userId, postId, database.DB)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return 0, nil
}
