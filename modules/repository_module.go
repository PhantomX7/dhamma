package modules

import (
	"go.uber.org/fx"

	cardRepo "github.com/PhantomX7/dhamma/modules/card/repository"
	chatTemplateRepo "github.com/PhantomX7/dhamma/modules/chat_template/repository"
	domainRepo "github.com/PhantomX7/dhamma/modules/domain/repository"
	eventRepo "github.com/PhantomX7/dhamma/modules/event/repository"
	eventAttendanceRepo "github.com/PhantomX7/dhamma/modules/event_attendance/repository"
	followerRepo "github.com/PhantomX7/dhamma/modules/follower/repository"
	permissionRepo "github.com/PhantomX7/dhamma/modules/permission/repository"
	pointMutationRepo "github.com/PhantomX7/dhamma/modules/point_mutation/repository"
	refreshTokenRepo "github.com/PhantomX7/dhamma/modules/refresh_token/repository"
	roleRepo "github.com/PhantomX7/dhamma/modules/role/repository"
	userRepo "github.com/PhantomX7/dhamma/modules/user/repository"
	userDomainRepo "github.com/PhantomX7/dhamma/modules/user_domain/repository"
	userRoleRepo "github.com/PhantomX7/dhamma/modules/user_role/repository"
)

var RepositoryModule = fx.Options(
	fx.Provide(
		cardRepo.New,
		chatTemplateRepo.New,
		domainRepo.New,
		eventRepo.New,
		eventAttendanceRepo.New,
		followerRepo.New,
		permissionRepo.New,
		pointMutationRepo.New,
		refreshTokenRepo.New,
		roleRepo.New,
		userRepo.New,
		userDomainRepo.New,
		userRoleRepo.New,
	),
)
