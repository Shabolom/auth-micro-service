package user

import (
	"context"
	"fmt"
)

func (s *Service) DeleteUser(ctx context.Context, userID string) error {
	err := s.userRepo.DeleteUser(ctx, userID)
	if err != nil {
		s.logger.Info("delete error")
		return fmt.Errorf("DeleteUser: %w", err)
	}

	return nil
}
