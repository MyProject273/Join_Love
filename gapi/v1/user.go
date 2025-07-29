package v1

import (
	"context"
	"errors"
	"fmt"

	"github.com/MyProject273/Join_Love/gapi/helper"
	db "github.com/MyProject273/Join_Love/internal/db/sqlc"
	"github.com/MyProject273/Join_Love/pb/user"
	"github.com/MyProject273/Join_Love/pb/user/profile"
	"github.com/MyProject273/Join_Love/pkg/config"
	"github.com/MyProject273/Join_Love/pkg/i18n"
	"github.com/MyProject273/Join_Love/pkg/utils"
	"github.com/jackc/pgx"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServer struct {
	user.UnimplementedUserServiceServer
	store  db.Store
	config config.Config
	logger zerolog.Logger
}

func NewUserServer(config config.Config, store db.Store, logger zerolog.Logger) *UserServer {
	return &UserServer{
		store:  store,
		config: config,
		logger: logger,
	}
}

// Profile
func (u *UserServer) GetUserProfile(ctx context.Context, req *profile.GetUserProfileRequest) (res *profile.UserProfileResponse, err error) {
	payload, ok := helper.GetAuthPayload(ctx)
	fmt.Println(payload)
	lang := helper.ExtractMetadata(ctx).Lang
	if !ok {
		return nil, status.Error(codes.Internal, "unauthenticated context")

	}
	if !helper.HasPermission(payload.Role, []string{"admin", "user"}) {
		return nil, status.Error(codes.PermissionDenied, "access denied")
	}

	uuid, err := utils.StringToPgUUID(req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}
	profileData, err := u.store.GetUserProfile(ctx, uuid)
	if err != nil {
		var pgErr *pgx.PgError
		if errors.As(err, &pgErr) && pgErr.Message == pgx.ErrNoRows.Error() {
			u.logger.Error().Err(err).Str("email", req.GetUserId()).Msgf("failed to check existing user (type: %T)", err)
			return nil, status.Errorf(codes.Internal, "%s", i18n.GetI18nMessage("internal_error", lang))
		}
		return nil, status.Errorf(codes.Internal, "%s", i18n.GetI18nMessage("internal_error", lang))
	}

	res = &profile.UserProfileResponse{
		UserId:           req.UserId,
		JobTitle:         *profileData.JobTitle,
		Education:        *profileData.Education,
		Interests:        *profileData.Interests,
		RelationshipGoal: *profileData.RelationshipGoal,
		Religion:         *profileData.Religion,
		Smoking:          *profileData.Smoking,
		Drinking:         *profileData.Drinking,
		HeightCm:         *profileData.HeightCm,
		WeightKg:         *profileData.WeightKg,
	}
	return res, nil

}
