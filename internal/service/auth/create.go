package auth

import (
	"context"
	"encoding/json"
	"log"
	"strconv"

	"github.com/en7ka/auth/internal/models"
	repoconv "github.com/en7ka/auth/internal/repository/auth/converter"
	"github.com/segmentio/kafka-go"
)

type userRegisteredEvent struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

func (s *serv) Create(ctx context.Context, info *models.UserInfo) (int64, error) {
	if info == nil {
		return 0, nil
	}

	var userID int64
	err := s.txManager.ReadCommited(ctx, func(ctx context.Context) error {
		var txErr error
		userID, txErr = s.userRepository.Create(ctx, repoconv.ToRepoUserInfo(info))
		if txErr != nil {
			return txErr
		}
		return nil
	})
	if err != nil {
		return 0, err
	}

	log.Printf("--- User created with ID: %d. Checking if producer exists... ---", userID)

	if s.producer != nil {

		log.Println("--- Producer exists. Attempting to write message... ---")

		evt := userRegisteredEvent{
			ID:    userID,
			Name:  info.Username,
			Email: info.Email,
			Role:  info.Role,
		}
		payload, mErr := json.Marshal(evt)
		if mErr != nil {

			log.Printf("marshal user-registered event: %v", mErr)
		} else {
			if wErr := s.producer.WriteMessages(ctx, kafka.Message{
				Key:   []byte(strconv.FormatInt(userID, 10)),
				Value: payload,
			}); wErr != nil {
				log.Printf("publish user-registered failed: %v", wErr)
			}
		}
	}

	return userID, nil
}
